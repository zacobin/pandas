package xstream

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloustone/pandas/kuiper/util"
	"github.com/cloustone/pandas/kuiper/xstream/api"
	"github.com/cloustone/pandas/kuiper/xstream/checkpoints"
	"github.com/cloustone/pandas/kuiper/xstream/contexts"
	"github.com/cloustone/pandas/kuiper/xstream/nodes"
	"github.com/cloustone/pandas/kuiper/xstream/states"
)

type PrintableTopo struct {
	Sources []string            `json:"sources"`
	Edges   map[string][]string `json:"edges"`
}

type TopologyNew struct {
	sources            []*nodes.SourceNode
	sinks              []*nodes.SinkNode
	ctx                api.StreamContext
	cancel             context.CancelFunc
	drain              chan error
	ops                []nodes.OperatorNode
	name               string
	qos                api.Qos
	checkpointInterval int
	store              api.Store
	coordinator        *checkpoints.Coordinator
	topo               *PrintableTopo
}

func NewWithNameAndQos(name string, qos api.Qos, checkpointInterval int) (*TopologyNew, error) {
	tp := &TopologyNew{
		name:               name,
		drain:              make(chan error),
		qos:                qos,
		checkpointInterval: checkpointInterval,
		topo: &PrintableTopo{
			Sources: make([]string, 0),
			Edges:   make(map[string][]string),
		},
	}
	return tp, nil
}

func (s *TopologyNew) GetContext() api.StreamContext {
	return s.ctx
}

func (s *TopologyNew) Cancel() {
	s.cancel()
	s.store = nil
	s.coordinator = nil
}

func (s *TopologyNew) AddSrc(src *nodes.SourceNode) *TopologyNew {
	s.sources = append(s.sources, src)
	s.topo.Sources = append(s.topo.Sources, fmt.Sprintf("source_%s", src.GetName()))
	return s
}

func (s *TopologyNew) AddSink(inputs []api.Emitter, snk *nodes.SinkNode) *TopologyNew {
	for _, input := range inputs {
		input.AddOutput(snk.GetInput())
		snk.AddInputCount()
		s.addEdge(input.(api.TopNode), snk, "sink")
	}
	s.sinks = append(s.sinks, snk)
	return s
}

func (s *TopologyNew) AddOperator(inputs []api.Emitter, operator nodes.OperatorNode) *TopologyNew {
	for _, input := range inputs {
		input.AddOutput(operator.GetInput())
		operator.AddInputCount()
		s.addEdge(input.(api.TopNode), operator, "op")
	}
	s.ops = append(s.ops, operator)
	return s
}

func (s *TopologyNew) addEdge(from api.TopNode, to api.TopNode, toType string) {
	fromType := "op"
	if _, ok := from.(*nodes.SourceNode); ok {
		fromType = "source"
	}
	f := fmt.Sprintf("%s_%s", fromType, from.GetName())
	t := fmt.Sprintf("%s_%s", toType, to.GetName())
	e, ok := s.topo.Edges[f]
	if !ok {
		e = make([]string, 0)
	}
	s.topo.Edges[f] = append(e, t)
}

func Transform(op nodes.UnOperation, name string, bufferLength int) *nodes.UnaryOperator {
	operator := nodes.New(name, bufferLength)
	operator.SetOperation(op)
	return operator
}

// prepareContext setups internal context before
// stream starts execution.
func (s *TopologyNew) prepareContext() {
	if s.ctx == nil || s.ctx.Err() != nil {
		contextLogger := util.Log.WithField("rule", s.name)
		ctx := contexts.WithValue(contexts.Background(), contexts.LoggerKey, contextLogger)
		s.ctx, s.cancel = ctx.WithCancel()
	}
}

func (s *TopologyNew) drainErr(err error) {
	go func() { s.drain <- err }()
}

func (s *TopologyNew) Open() <-chan error {

	//if stream has opened, do nothing
	if s.ctx != nil && s.ctx.Err() == nil {
		s.ctx.GetLogger().Infoln("rule is already running, do nothing")
		return s.drain
	}
	s.prepareContext() // ensure context is set
	var err error
	if s.store, err = states.CreateStore(s.name, s.qos); err != nil {
		s.drainErr(err)
		return s.drain
	}
	s.enableCheckpoint()
	log := s.ctx.GetLogger()
	log.Infoln("Opening stream")
	// open stream
	go func() {
		// open stream sink, after log sink is ready.
		for _, snk := range s.sinks {
			snk.Open(s.ctx.WithMeta(s.name, snk.GetName(), s.store), s.drain)
		}

		//apply operators, if err bail
		for _, op := range s.ops {
			op.Exec(s.ctx.WithMeta(s.name, op.GetName(), s.store), s.drain)
		}

		// open source, if err bail
		for _, node := range s.sources {
			node.Open(s.ctx.WithMeta(s.name, node.GetName(), s.store), s.drain)
		}

		// activate checkpoint
		if s.coordinator != nil {
			s.coordinator.Activate()
		}
	}()

	return s.drain
}

func (s *TopologyNew) enableCheckpoint() error {
	if s.qos >= api.AtLeastOnce {
		var sources []checkpoints.StreamTask
		for _, r := range s.sources {
			sources = append(sources, r)
		}
		var ops []checkpoints.NonSourceTask
		for _, r := range s.ops {
			ops = append(ops, r)
		}
		var sinks []checkpoints.SinkTask
		for _, r := range s.sinks {
			sinks = append(sinks, r)
		}
		c := checkpoints.NewCoordinator(s.name, sources, ops, sinks, s.qos, s.store, s.checkpointInterval, s.ctx)
		s.coordinator = c
	}
	return nil
}

func (s *TopologyNew) GetCoordinator() *checkpoints.Coordinator {
	return s.coordinator
}

func (s *TopologyNew) GetMetrics() (keys []string, values []interface{}) {
	for _, node := range s.sources {
		for ins, metrics := range node.GetMetrics() {
			for i, v := range metrics {
				keys = append(keys, "source_"+node.GetName()+"_"+strconv.Itoa(ins)+"_"+nodes.MetricNames[i])
				values = append(values, v)
			}
		}
	}
	for _, node := range s.ops {
		for ins, metrics := range node.GetMetrics() {
			for i, v := range metrics {
				keys = append(keys, "op_"+node.GetName()+"_"+strconv.Itoa(ins)+"_"+nodes.MetricNames[i])
				values = append(values, v)
			}
		}
	}
	for _, node := range s.sinks {
		for ins, metrics := range node.GetMetrics() {
			for i, v := range metrics {
				keys = append(keys, "sink_"+node.GetName()+"_"+strconv.Itoa(ins)+"_"+nodes.MetricNames[i])
				values = append(values, v)
			}
		}
	}
	return
}

func (s *TopologyNew) GetTopo() *PrintableTopo {
	return s.topo
}
