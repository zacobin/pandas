package contexts

import (
	"github.com/cloustone/pandas/kuiper"
	"github.com/cloustone/pandas/kuiper/xstream/api"
	"github.com/cloustone/pandas/kuiper/xstream/states"
	"log"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestState(t *testing.T) {
	var (
		i      = 0
		ruleId = "testStateRule"
		value1 = 21
		value2 = "hello"
		value3 = "world"
		s      = map[string]interface{}{
			"key1": 21,
			"key3": "world",
		}
	)
	//initialization
	store, err := states.CreateStore(ruleId, api.AtLeastOnce)
	if err != nil {
		t.Errorf("Get store for rule %s error: %s", ruleId, err)
		return
	}
	ctx := Background().WithMeta("testStateRule", "op1", store).(*DefaultContext)
	defer cleanStateData()
	// Do state function
	_ = ctx.IncrCounter("key1", 20)
	_ = ctx.IncrCounter("key1", 1)
	v, err := ctx.GetCounter("key1")
	if err != nil {
		t.Errorf("%d.Get counter error: %s", i, err)
		return
	}
	if !reflect.DeepEqual(value1, v) {
		t.Errorf("%d.Get counter\n\nresult mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, value1, v)
	}
	err = ctx.PutState("key2", value2)
	if err != nil {
		t.Errorf("%d.Put state key2 error: %s", i, err)
		return
	}
	err = ctx.PutState("key3", value3)
	if err != nil {
		t.Errorf("%d.Put state key3 error: %s", i, err)
		return
	}
	v2, err := ctx.GetState("key2")
	if err != nil {
		t.Errorf("%d.Get state key2 error: %s", i, err)
		return
	}
	if !reflect.DeepEqual(value2, v2) {
		t.Errorf("%d.Get state\n\nresult mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, value2, v2)
	}
	err = ctx.DeleteState("key2")
	if err != nil {
		t.Errorf("%d.Delete state key2 error: %s", i, err)
		return
	}
	err = ctx.Snapshot()
	if err != nil {
		t.Errorf("%d.Snapshot error: %s", i, err)
		return
	}
	rs := ctx.snapshot
	if !reflect.DeepEqual(s, rs) {
		t.Errorf("%d.Snapshot\n\nresult mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, s, rs)
	}
}

func cleanStateData() {
	dbDir, err := util.GetDataLoc()
	if err != nil {
		log.Panic(err)
	}
	c := path.Join(dbDir, "checkpoints")
	err = os.RemoveAll(c)
	if err != nil {
		util.Log.Error(err)
	}
}
