package nodes

import (
	"fmt"
	"github.com/cloustone/pandas/kuiper/xsql"
	"github.com/cloustone/pandas/kuiper/xstream/api"
	"sync"
)

// UnOperation interface represents unary operations (i.e. Map, Filter, etc)
type UnOperation interface {
	Apply(ctx api.StreamContext, data interface{}, fv *xsql.FunctionValuer, afv *xsql.AggregateFunctionValuer) interface{}
}

// UnFunc implements UnOperation as type func (context.Context, interface{})
type UnFunc func(api.StreamContext, interface{}) interface{}

// Apply implements UnOperation.Apply method
func (f UnFunc) Apply(ctx api.StreamContext, data interface{}) interface{} {
	return f(ctx, data)
}

type UnaryOperator struct {
	*defaultSinkNode
	op        UnOperation
	mutex     sync.RWMutex
	cancelled bool
}

// NewUnary creates *UnaryOperator value
func New(name string, bufferLength int) *UnaryOperator {
	return &UnaryOperator{
		defaultSinkNode: &defaultSinkNode{
			input: make(chan interface{}, bufferLength),
			defaultNode: &defaultNode{
				name:        name,
				outputs:     make(map[string]chan<- interface{}),
				concurrency: 1,
			},
		},
	}
}

// SetOperation sets the executor operation
func (o *UnaryOperator) SetOperation(op UnOperation) {
	o.op = op
}

// Exec is the entry point for the executor
func (o *UnaryOperator) Exec(ctx api.StreamContext, errCh chan<- error) {
	o.ctx = ctx
	log := ctx.GetLogger()
	log.Debugf("Unary operator %s is started", o.name)

	if len(o.outputs) <= 0 {
		go func() { errCh <- fmt.Errorf("no output channel found") }()
		return
	}

	// validate p
	if o.concurrency < 1 {
		o.concurrency = 1
	}
	//reset status
	o.statManagers = nil

	for i := 0; i < o.concurrency; i++ { // workers
		instance := i
		go o.doOp(ctx.WithInstance(instance), errCh)
	}
}

func (o *UnaryOperator) doOp(ctx api.StreamContext, errCh chan<- error) {
	logger := ctx.GetLogger()
	if o.op == nil {
		logger.Infoln("Unary operator missing operation")
		return
	}
	exeCtx, cancel := ctx.WithCancel()

	defer func() {
		logger.Infof("unary operator %s instance %d done, cancelling future items", o.name, ctx.GetInstanceId())
		cancel()
	}()

	stats, err := NewStatManager("op", ctx)
	if err != nil {
		select {
		case errCh <- err:
			logger.Errorf("unary operator %s error %s", o.name, err)
		case <-ctx.Done():
			logger.Infof("unary operator %s cancelling....", o.name)
			o.mutex.Lock()
			cancel()
			o.cancelled = true
			o.mutex.Unlock()
		}
		return
	}
	o.mutex.Lock()
	o.statManagers = append(o.statManagers, stats)
	o.mutex.Unlock()
	fv, afv := xsql.NewFunctionValuersForOp(exeCtx)

	for {
		select {
		// process incoming item
		case item := <-o.input:
			processed := false
			if item, processed = o.preprocess(item); processed {
				break
			}
			stats.IncTotalRecordsIn()
			stats.ProcessTimeStart()
			result := o.op.Apply(exeCtx, item, fv, afv)

			switch val := result.(type) {
			case nil:
				continue
			case error:
				logger.Errorf("Operation %s error: %s", ctx.GetOpId(), val)
				o.Broadcast(val)
				stats.IncTotalExceptions()
				continue
			default:
				stats.ProcessTimeEnd()
				o.Broadcast(val)
				stats.IncTotalRecordsOut()
				stats.SetBufferLength(int64(len(o.input)))
			}
		// is cancelling
		case <-ctx.Done():
			logger.Infof("unary operator %s instance %d cancelling....", o.name, ctx.GetInstanceId())
			o.mutex.Lock()
			cancel()
			o.cancelled = true
			o.mutex.Unlock()
			return
		}
	}
}
