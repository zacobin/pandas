package kuiper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/cloustone/pandas/kuiper/kvstore"
	"github.com/cloustone/pandas/kuiper/util"
	"github.com/cloustone/pandas/kuiper/xsql"
	"github.com/cloustone/pandas/kuiper/xsql/plans"
	"github.com/cloustone/pandas/kuiper/xstream"
	"github.com/cloustone/pandas/kuiper/xstream/api"
	"github.com/cloustone/pandas/kuiper/xstream/nodes"
)

var (
	dataDir   string
	rootDbDir string
)

type ruleManager struct {
	registry   ruleRegistry
	repository RuleRepository
}

func newRuleManager(r RuleRepository) *ruleManager {
	return &ruleManager{
		repository: r,
		registry:   ruleRegistry{internal: make(map[string]*ruleState)},
	}
}

func (rm *ruleManager) createRuleState(rule *api.Rule) (*ruleState, error) {
	rs := &ruleState{
		name: rule.Id,
	}
	rm.registry.store(rule.Id, rs)
	if tp, err := rm.execInitRule(rule); err != nil {
		return rs, err
	} else {
		rs.topology = tp
		rs.triggered = true
		return rs, nil
	}
}

func (rm *ruleManager) execInitRule(rule *api.Rule) (*xstream.TopologyNew, error) {
	if tp, inputs, err := rm.createTopo(rule); err != nil {
		return nil, err
	} else {
		for i, m := range rule.Actions {
			for name, action := range m {
				props, ok := action.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("expect map[string]interface{} type for the action properties, but found %v", action)
				}
				tp.AddSink(inputs, nodes.NewSinkNode(fmt.Sprintf("%s_%d", name, i), name, props))
			}
		}
		return tp, nil
	}
}

func getStream(m kvstore.KvStore, name string) (stmt *xsql.StreamStmt, err error) {
	s, f := m.Get(name)
	if !f {
		return nil, fmt.Errorf("Cannot find key %s. ", name)
	}
	s1, _ := s.(string)
	parser := xsql.NewParser(strings.NewReader(s1))
	stream, err := xsql.Language.Parse(parser)
	stmt, ok := stream.(*xsql.StreamStmt)
	if !ok {
		err = fmt.Errorf("Error resolving the stream %s, the data in db may be corrupted.", name)
	}
	return
}

func (rm *ruleManager) createTopo(rule *api.Rule) (*xstream.TopologyNew, []api.Emitter, error) {
	return rm.createTopoWithSources(rule, nil)
}

func (rm *ruleManager) createTopoWithSources(rule *api.Rule, sources []*nodes.SourceNode) (*xstream.TopologyNew, []api.Emitter, error) {
	name := rule.Id
	sql := rule.Sql

	// log.Infof("Init rule with options %+v", rule.Options)
	shouldCreateSource := sources == nil

	if selectStmt, err := getStatementFromSql(sql); err != nil {
		return nil, nil, err
	} else {
		tp, err := xstream.NewWithNameAndQos(name, rule.Options.Qos, rule.Options.CheckpointInterval)
		if err != nil {
			return nil, nil, err
		}
		var inputs []api.Emitter
		streamsFromStmt := xsql.GetStreams(selectStmt)
		dimensions := selectStmt.Dimensions
		if !shouldCreateSource && len(streamsFromStmt) != len(sources) {
			return nil, nil, fmt.Errorf("Invalid parameter sources or streams, the length cannot match the statement, expect %d sources.", len(streamsFromStmt))
		}
		if rule.Options.SendMetaToSink && (len(streamsFromStmt) > 1 || dimensions != nil) {
			return nil, nil, fmt.Errorf("Invalid option sendMetaToSink, it can not be applied to window")
		}
		store := kvstore.GetKvStore(path.Join(rootDbDir, "stream"))
		err = store.Open()
		if err != nil {
			return nil, nil, err
		}
		defer store.Close()

		var alias, aggregateAlias xsql.Fields
		for _, f := range selectStmt.Fields {
			if f.AName != "" {
				if !xsql.HasAggFuncs(f.Expr) {
					alias = append(alias, f)
				} else {
					aggregateAlias = append(aggregateAlias, f)
				}
			}
		}
		for i, s := range streamsFromStmt {
			streamStmt, err := getStream(store, s)
			if err != nil {
				return nil, nil, fmt.Errorf("fail to get stream %s, please check if stream is created", s)
			}
			pp, err := plans.NewPreprocessor(streamStmt, alias, rule.Options.IsEventTime)
			if err != nil {
				return nil, nil, err
			}
			var srcNode *nodes.SourceNode
			if shouldCreateSource {
				node := nodes.NewSourceNode(s, streamStmt.Options)
				srcNode = node
			} else {
				srcNode = sources[i]
			}
			tp.AddSrc(srcNode)
			preprocessorOp := xstream.Transform(pp, "preprocessor_"+s, rule.Options.BufferLength)
			preprocessorOp.SetConcurrency(rule.Options.Concurrency)
			tp.AddOperator([]api.Emitter{srcNode}, preprocessorOp)
			inputs = append(inputs, preprocessorOp)
		}

		var w *xsql.Window
		if dimensions != nil {
			w = dimensions.GetWindow()
			if w != nil {
				if w.Filter != nil {
					wfilterOp := xstream.Transform(&plans.FilterPlan{Condition: w.Filter}, "windowFilter", rule.Options.BufferLength)
					wfilterOp.SetConcurrency(rule.Options.Concurrency)
					tp.AddOperator(inputs, wfilterOp)
					inputs = []api.Emitter{wfilterOp}
				}
				wop, err := nodes.NewWindowOp("window", w, rule.Options.IsEventTime, rule.Options.LateTol, streamsFromStmt, rule.Options.BufferLength)
				if err != nil {
					return nil, nil, err
				}
				tp.AddOperator(inputs, wop)
				inputs = []api.Emitter{wop}
			}
		}

		if w != nil && selectStmt.Joins != nil {
			joinOp := xstream.Transform(&plans.JoinPlan{Joins: selectStmt.Joins, From: selectStmt.Sources[0].(*xsql.Table)}, "join", rule.Options.BufferLength)
			joinOp.SetConcurrency(rule.Options.Concurrency)
			tp.AddOperator(inputs, joinOp)
			inputs = []api.Emitter{joinOp}
		}

		if selectStmt.Condition != nil {
			filterOp := xstream.Transform(&plans.FilterPlan{Condition: selectStmt.Condition}, "filter", rule.Options.BufferLength)
			filterOp.SetConcurrency(rule.Options.Concurrency)
			tp.AddOperator(inputs, filterOp)
			inputs = []api.Emitter{filterOp}
		}

		var ds xsql.Dimensions
		if dimensions != nil || len(aggregateAlias) > 0 {
			ds = dimensions.GetGroups()
			if (ds != nil && len(ds) > 0) || len(aggregateAlias) > 0 {
				aggregateOp := xstream.Transform(&plans.AggregatePlan{Dimensions: ds, Alias: aggregateAlias}, "aggregate", rule.Options.BufferLength)
				aggregateOp.SetConcurrency(rule.Options.Concurrency)
				tp.AddOperator(inputs, aggregateOp)
				inputs = []api.Emitter{aggregateOp}
			}
		}

		if selectStmt.Having != nil {
			havingOp := xstream.Transform(&plans.HavingPlan{selectStmt.Having}, "having", rule.Options.BufferLength)
			havingOp.SetConcurrency(rule.Options.Concurrency)
			tp.AddOperator(inputs, havingOp)
			inputs = []api.Emitter{havingOp}
		}

		if selectStmt.SortFields != nil {
			orderOp := xstream.Transform(&plans.OrderPlan{SortFields: selectStmt.SortFields}, "order", rule.Options.BufferLength)
			orderOp.SetConcurrency(rule.Options.Concurrency)
			tp.AddOperator(inputs, orderOp)
			inputs = []api.Emitter{orderOp}
		}

		if selectStmt.Fields != nil {
			projectOp := xstream.Transform(&plans.ProjectPlan{Fields: selectStmt.Fields, IsAggregate: xsql.IsAggStatement(selectStmt), SendMeta: rule.Options.SendMetaToSink}, "project", rule.Options.BufferLength)
			projectOp.SetConcurrency(rule.Options.Concurrency)
			tp.AddOperator(inputs, projectOp)
			inputs = []api.Emitter{projectOp}
		}
		return tp, inputs, nil
	}
}

func (rm *ruleManager) doStartRule(rs *ruleState) error {
	rs.triggered = true
	//	rm.ruleProcessor.ExecReplaceruleState(rs.name, true)
	go func() {
		tp := rs.topology
		select {
		case err := <-tp.Open():
			tp.GetContext().SetError(err)
			// logger.Printf("closing rule %s for error: %v", rs.Name, err)
			tp.Cancel()
		}
	}()
	return nil
}

func (rm *ruleManager) getAllRules() ([]string, error) {
	return []string{}, nil // TODO
}

func (rm *ruleManager) getAllRulesWithStatus() ([]map[string]interface{}, error) {
	names, err := rm.getAllRules()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	result := make([]map[string]interface{}, len(names))
	for i, name := range names {
		s, err := rm.getRuleState(name)
		if err != nil {
			return nil, err
		}
		result[i] = map[string]interface{}{
			"id":     name,
			"status": s,
		}
	}
	return result, nil
}

func (rm *ruleManager) getRuleState(name string) (string, error) {
	if rs, ok := rm.registry.load(name); ok {
		return rm.doGetruleState(rs)
	} else {
		return "", fmt.Errorf("Rule %s is not found", name)
	}
}

func (rm *ruleManager) doGetruleState(rs *ruleState) (string, error) {
	result := ""
	if !rs.triggered {
		result = "Stopped: canceled manually or by error."
		return result, nil
	}
	c := (*rs.topology).GetContext()
	if c != nil {
		err := c.Err()
		switch err {
		case nil:
			result = "Running"
		case context.Canceled:
			result = "Stopped: canceled by error."
		case context.DeadlineExceeded:
			result = "Stopped: deadline exceed."
		default:
			result = fmt.Sprintf("Stopped: %v.", err)
		}
	} else {
		result = "Stopped: no context found."
	}
	return result, nil
}

func (rm *ruleManager) getRuleStatus(_ context.Context, token string, name string) (string, error) {
	if rs, ok := rm.registry.load(name); ok {
		result, err := rm.doGetruleState(rs)
		if err != nil {
			return "", err
		}
		if result == "Running" {
			keys, values := (*rs.topology).GetMetrics()
			metrics := "{"
			for i, key := range keys {
				value := values[i]
				switch value.(type) {
				case string:
					metrics += fmt.Sprintf("\"%s\":%q,", key, value)
				default:
					metrics += fmt.Sprintf("\"%s\":%v,", key, value)
				}
			}
			metrics = metrics[:len(metrics)-1] + "}"
			dst := &bytes.Buffer{}
			if err = json.Indent(dst, []byte(metrics), "", "  "); err != nil {
				result = metrics
			} else {
				result = dst.String()
			}
		}
		return result, nil
	} else {
		return "", util.NewErrorWithCode(util.NOT_FOUND, fmt.Sprintf("Rule %s is not found", name))
	}
}

func (rm *ruleManager) getRuleTopo(_ context.Context, token string, name string) (string, error) {
	if rs, ok := rm.registry.load(name); ok {
		topo := rs.topology.GetTopo()
		bytes, err := json.Marshal(topo)
		if err != nil {
			return "", util.NewError(fmt.Sprintf("Fail to encode rule %s's topo", name))
		} else {
			return string(bytes), nil
		}
	} else {
		return "", util.NewErrorWithCode(util.NOT_FOUND, fmt.Sprintf("Rule %s is not found", name))
	}
}

func (rm *ruleManager) startRule(r *api.Rule) error {
	var rs *ruleState
	var err error

	rs, ok := rm.registry.load(r.Id)
	if !ok || (!rs.triggered) {
		rs, err = rm.createRuleState(r)
		if err != nil {
			return err
		}
	}
	return rm.doStartRule(rs)
}

func (rm *ruleManager) getRuleByName(name string) (*api.Rule, error) {
	// p.getRuleByJson(name, s1)
	return nil, nil //TODO
}

func (rm *ruleManager) stopRule(r *api.Rule) (err error) {
	if rs, ok := rm.registry.load(r.Id); ok && rs.triggered {
		(*rs.topology).Cancel()
		rs.triggered = false
		//rm.ruleProcessor.ExecReplaceruleState(name, false)
		err = fmt.Errorf("Rule %s was stopped.", r.Id)
	} else {
		err = fmt.Errorf("Rule %s was not found.", r.Id)
	}
	return
}

func (rm *ruleManager) restartRule(r *api.Rule) error {
	rm.stopRule(r)
	return rm.startRule(r)
}

func (rm *ruleManager) recoverRule(name string) string {
	rule, err := rm.getRuleByName(name)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	if !rule.Triggered {
		rs := &ruleState{
			name: name,
		}
		rm.registry.store(name, rs)
		return fmt.Sprintf("Rule %s was stoped.", name)
	}

	// TODO err = rm.startRule(name)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return fmt.Sprintf("Rule %s was started.", name)

}

func getStatementFromSql(sql string) (*xsql.SelectStatement, error) {
	parser := xsql.NewParser(strings.NewReader(sql))
	if stmt, err := xsql.Language.Parse(parser); err != nil {
		return nil, fmt.Errorf("Parse SQL %s error: %s.", sql, err)
	} else {
		if r, ok := stmt.(*xsql.SelectStatement); !ok {
			return nil, fmt.Errorf("SQL %s is not a select statement.", sql)
		} else {
			return r, nil
		}
	}
}

func (rm *ruleManager) buileRuleByJson(name, ruleJson string) (*api.Rule, error) {
	opt := util.Config.Rule
	//set default rule options
	rule := &api.Rule{
		Options: &opt,
	}
	if err := json.Unmarshal([]byte(ruleJson), &rule); err != nil {
		return nil, fmt.Errorf("Parse rule %s error : %s.", ruleJson, err)
	}

	//validation
	if rule.Id == "" && name == "" {
		return nil, fmt.Errorf("Missing rule id.")
	}
	if name != "" && rule.Id != "" && name != rule.Id {
		return nil, fmt.Errorf("Name is not consistent with rule id.")
	}
	if rule.Id == "" {
		rule.Id = name
	}
	if rule.Sql == "" {
		return nil, fmt.Errorf("Missing rule SQL.")
	}
	if _, err := getStatementFromSql(rule.Sql); err != nil {
		return nil, err
	}
	if rule.Actions == nil || len(rule.Actions) == 0 {
		return nil, fmt.Errorf("Missing rule actions.")
	}
	if rule.Options == nil {
		rule.Options = &api.RuleOption{}
	}
	//Set default options
	if rule.Options.CheckpointInterval < 0 {
		return nil, fmt.Errorf("rule option checkpointInterval %d is invalid, require a positive integer", rule.Options.CheckpointInterval)
	}
	if rule.Options.Concurrency < 0 {
		return nil, fmt.Errorf("rule option concurrency %d is invalid, require a positive integer", rule.Options.Concurrency)
	}
	if rule.Options.BufferLength < 0 {
		return nil, fmt.Errorf("rule option bufferLength %d is invalid, require a positive integer", rule.Options.BufferLength)
	}
	if rule.Options.LateTol < 0 {
		return nil, fmt.Errorf("rule option lateTolerance %d is invalid, require a positive integer", rule.Options.LateTol)
	}
	return rule, nil
}

func (rm *ruleManager) execQuery(ruleid, sql string) (*xstream.TopologyNew, error) {
	if tp, inputs, err := rm.createTopo(rm.getDefaultRule(ruleid, sql)); err != nil {
		return nil, err
	} else {
		tp.AddSink(inputs, nodes.NewSinkNode("sink_memory_log", "logToMemory", nil))
		go func() {
			select {
			case err := <-tp.Open():
				//log.Infof("closing query for error: %v", err)
				tp.GetContext().SetError(err)
				tp.Cancel()
			}
		}()
		return tp, nil
	}
}

func (rm *ruleManager) getDefaultRule(name, sql string) *api.Rule {
	return &api.Rule{
		Id:  name,
		Sql: sql,
		Options: &api.RuleOption{
			IsEventTime:        false,
			LateTol:            1000,
			Concurrency:        1,
			BufferLength:       1024,
			SendMetaToSink:     false,
			Qos:                api.AtMostOnce,
			CheckpointInterval: 300000,
		},
	}
}

func (rm *ruleManager) getRuleByJson(name, ruleJson string) (*api.Rule, error) {
	opt := util.Config.Rule
	//set default rule options
	rule := &api.Rule{
		Options: &opt,
	}
	if err := json.Unmarshal([]byte(ruleJson), &rule); err != nil {
		return nil, fmt.Errorf("Parse rule %s error : %s.", ruleJson, err)
	}

	//validation
	if rule.Id == "" && name == "" {
		return nil, fmt.Errorf("Missing rule id.")
	}
	if name != "" && rule.Id != "" && name != rule.Id {
		return nil, fmt.Errorf("Name is not consistent with rule id.")
	}
	if rule.Id == "" {
		rule.Id = name
	}
	if rule.Sql == "" {
		return nil, fmt.Errorf("Missing rule SQL.")
	}
	if _, err := getStatementFromSql(rule.Sql); err != nil {
		return nil, err
	}
	if rule.Actions == nil || len(rule.Actions) == 0 {
		return nil, fmt.Errorf("Missing rule actions.")
	}
	if rule.Options == nil {
		rule.Options = &api.RuleOption{}
	}
	//Set default options
	if rule.Options.CheckpointInterval < 0 {
		return nil, fmt.Errorf("rule option checkpointInterval %d is invalid, require a positive integer", rule.Options.CheckpointInterval)
	}
	if rule.Options.Concurrency < 0 {
		return nil, fmt.Errorf("rule option concurrency %d is invalid, require a positive integer", rule.Options.Concurrency)
	}
	if rule.Options.BufferLength < 0 {
		return nil, fmt.Errorf("rule option bufferLength %d is invalid, require a positive integer", rule.Options.BufferLength)
	}
	if rule.Options.LateTol < 0 {
		return nil, fmt.Errorf("rule option lateTolerance %d is invalid, require a positive integer", rule.Options.LateTol)
	}
	return rule, nil
}

func (rm *ruleManager) deleteRule(rule *api.Rule) (err error) {
	if err := cleanSinkCache(rule); err != nil {
		err = fmt.Errorf("Clean sink cache faile: %s.", err)
	}
	if err := cleanCheckpoint(rule); err != nil {
		err = fmt.Errorf("Clean checkpoint cache faile: %s.", err)
	}
	return
}

func cleanCheckpoint(r *api.Rule) error {
	dbDir, _ := util.GetDataLoc()
	c := path.Join(dbDir, "checkpoints", r.Id)
	return os.RemoveAll(c)
}

func cleanSinkCache(rule *api.Rule) error {
	dbDir, err := util.GetDataLoc()
	if err != nil {
		return err
	}
	store := kvstore.GetKvStore(path.Join(dbDir, "sink"))
	err = store.Open()
	if err != nil {
		return err
	}
	defer store.Close()
	for d, m := range rule.Actions {
		con := 1
		for name, action := range m {
			props, _ := action.(map[string]interface{})
			if c, ok := props["concurrency"]; ok {
				if t, err := util.ToInt(c); err == nil && t > 0 {
					con = t
				}
			}
			for i := 0; i < con; i++ {
				key := fmt.Sprintf("%s%s_%d%d", rule.Id, name, d, i)
				util.Log.Debugf("delete cache key %s", key)
				store.Delete(key)
			}
		}
	}
	return nil
}
