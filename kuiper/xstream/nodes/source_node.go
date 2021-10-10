package nodes

import (
	"sync"

	"github.com/cloustone/pandas/kuiper/plugins"
	"github.com/cloustone/pandas/kuiper/util"
	"github.com/cloustone/pandas/kuiper/xsql"
	"github.com/cloustone/pandas/kuiper/xstream/api"
	"github.com/cloustone/pandas/kuiper/xstream/extensions"
	"github.com/go-yaml/yaml"
)

type SourceNode struct {
	*defaultNode
	sourceType string
	options    map[string]string
	isMock     bool

	mutex   sync.RWMutex
	sources []api.Source
}

func NewSourceNode(name string, options map[string]string) *SourceNode {
	t, ok := options["TYPE"]
	if !ok {
		t = "mqtt"
	}
	return &SourceNode{
		sourceType: t,
		defaultNode: &defaultNode{
			name:        name,
			outputs:     make(map[string]chan<- interface{}),
			concurrency: 1,
		},
		options: options,
	}
}

const OFFSET_KEY = "$$offset"

//Only for mock source, do not use it in production
func NewSourceNodeWithSource(name string, source api.Source, options map[string]string) *SourceNode {
	return &SourceNode{
		sources: []api.Source{source},
		defaultNode: &defaultNode{
			name:        name,
			outputs:     make(map[string]chan<- interface{}),
			concurrency: 1,
		},
		options: options,
		isMock:  true,
	}
}

func (m *SourceNode) Open(ctx api.StreamContext, errCh chan<- error) {
	m.ctx = ctx
	logger := ctx.GetLogger()
	logger.Infof("open source node %s with option %v", m.name, m.options)
	go func() {
		props := m.getConf(ctx)
		if c, ok := props["concurrency"]; ok {
			if t, err := util.ToInt(c); err != nil || t <= 0 {
				logger.Warnf("invalid type for concurrency property, should be positive integer but found %t", c)
			} else {
				m.concurrency = t
			}
		}
		bl := 102400
		if c, ok := props["bufferLength"]; ok {
			if t, err := util.ToInt(c); err != nil || t <= 0 {
				logger.Warnf("invalid type for bufferLength property, should be positive integer but found %t", c)
			} else {
				bl = t
			}
		}
		m.reset()
		logger.Infof("open source node %d instances", m.concurrency)
		for i := 0; i < m.concurrency; i++ { // workers
			go func(instance int) {
				//Do open source instances
				var source api.Source
				var err error
				if !m.isMock {
					source, err = getSource(m.sourceType)
					if err != nil {
						m.drainError(errCh, err, ctx, logger)
						return
					}
					err = source.Configure(m.options["DATASOURCE"], props)
					if err != nil {
						m.drainError(errCh, err, ctx, logger)
						return
					}
					m.mutex.Lock()
					m.sources = append(m.sources, source)
					m.mutex.Unlock()
				} else {
					logger.Debugf("get source instance %d from %d sources", instance, len(m.sources))
					source = m.sources[instance]
				}
				stats, err := NewStatManager("source", ctx)
				if err != nil {
					m.drainError(errCh, err, ctx, logger)
					return
				}
				m.mutex.Lock()
				m.statManagers = append(m.statManagers, stats)
				m.mutex.Unlock()

				if rw, ok := source.(api.Rewindable); ok {
					if offset, err := ctx.GetState(OFFSET_KEY); err != nil {
						m.drainError(errCh, err, ctx, logger)
					} else if offset != nil {
						logger.Infof("Source rewind from %v", offset)
						err = rw.Rewind(offset)
						if err != nil {
							m.drainError(errCh, err, ctx, logger)
						}
					}
				}

				buffer := NewDynamicChannelBuffer()
				buffer.SetLimit(bl)
				sourceErrCh := make(chan error)
				go source.Open(ctx.WithInstance(instance), buffer.In, sourceErrCh)
				logger.Infof("Start source %s instance %d successfully", m.name, instance)
				for {
					select {
					case <-ctx.Done():
						logger.Infof("source %s done", m.name)
						m.close(ctx, logger)
						return
					case err := <-sourceErrCh:
						m.drainError(errCh, err, ctx, logger)
						return
					case data := <-buffer.Out:
						stats.IncTotalRecordsIn()
						stats.ProcessTimeStart()
						tuple := &xsql.Tuple{Emitter: m.name, Message: data.Message(), Timestamp: util.GetNowInMilli(), Metadata: data.Meta()}
						stats.ProcessTimeEnd()
						logger.Debugf("source node %s is sending tuple %+v of timestamp %d", m.name, tuple, tuple.Timestamp)
						//blocking
						m.Broadcast(tuple)
						stats.IncTotalRecordsOut()
						stats.SetBufferLength(int64(buffer.GetLength()))
						if rw, ok := source.(api.Rewindable); ok {
							if offset, err := rw.GetOffset(); err != nil {
								m.drainError(errCh, err, ctx, logger)
							} else {
								err = ctx.PutState(OFFSET_KEY, offset)
								if err != nil {
									m.drainError(errCh, err, ctx, logger)
								}
								logger.Debugf("Source save offset %v", offset)
							}
						}
						logger.Debugf("source node %s has consumed tuple of timestamp %d", m.name, tuple.Timestamp)
					}
				}
			}(i)
		}
	}()
}

func (m *SourceNode) reset() {
	if !m.isMock {
		m.sources = nil
	}
	m.statManagers = nil
}

func doGetSource(t string) (api.Source, error) {
	var (
		s   api.Source
		err error
	)
	switch t {
	case "mqtt":
		s = &extensions.MQTTSource{}
	case "httppull":
		s = &extensions.HTTPPullSource{}
	default:
		s, err = plugins.GetSource(t)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (m *SourceNode) drainError(errCh chan<- error, err error, ctx api.StreamContext, logger api.Logger) {
	select {
	case errCh <- err:
	case <-ctx.Done():
		m.close(ctx, logger)
	}
	return
}

func (m *SourceNode) close(ctx api.StreamContext, logger api.Logger) {
	for _, s := range m.sources {
		if err := s.Close(ctx); err != nil {
			logger.Warnf("close source fails: %v", err)
		}
	}
}

func (m *SourceNode) getConf(ctx api.StreamContext) map[string]interface{} {
	confkey := m.options["CONF_KEY"]
	logger := ctx.GetLogger()
	confPath := "sources/" + m.sourceType + ".yaml"
	if m.sourceType == "mqtt" {
		confPath = "mqtt_source.yaml"
	}
	conf, err := util.LoadConf(confPath)
	props := make(map[string]interface{})
	if err == nil {
		cfg := make(map[interface{}]interface{})
		if err := yaml.Unmarshal(conf, &cfg); err != nil {
			logger.Warnf("fail to parse yaml for source %s. Return an empty configuration", m.sourceType)
		} else {
			def, ok := cfg["default"]
			if !ok {
				logger.Warnf("default conf is not found", confkey)
			} else {
				if def1, ok1 := def.(map[interface{}]interface{}); ok1 {
					props = util.ConvertMap(def1)
				}
				if c, ok := cfg[confkey]; ok {
					if c1, ok := c.(map[interface{}]interface{}); ok {
						c2 := util.ConvertMap(c1)
						for k, v := range c2 {
							props[k] = v
						}
					}
				}
			}
		}
	} else {
		logger.Warnf("config file %s.yaml is not loaded properly. Return an empty configuration", m.sourceType)
	}
	logger.Debugf("get conf for %s with conf key %s: %v", m.sourceType, confkey, props)
	return props
}
