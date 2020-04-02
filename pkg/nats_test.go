package broadcast

import (
	"fmt"
	"testing"
	"time"

	"github.com/cloustone/pandas/models"
	"github.com/cloustone/pandas/pkg/broadcast"
	broadcast_util "github.com/cloustone/pandas/pkg/broadcast/util"
	. "github.com/smartystreets/goconvey/convey"
)

type Abc struct{}

func (a *Abc) Onbroadcast(b broadcast.Broadcast, notify broadcast.Notification) {
	fmt.Println("onbroadcast")
	fmt.Println(notify.ObjectPath)
	fmt.Println("onbroadcast")
	switch notify.ObjectPath {
	case "123":
		fmt.Println("abcok")
	}
}
func Newabc() *Abc {
	return &Abc{}
}

func TestCacheSet(t *testing.T) {

	Convey("TestCacheSet should return ok when two cache items are same", t, func() {
		fmt.Println("hello")
		a := Newabc()
		broadcast_util.InitializeBroadcast(broadcast.NewServingOptions(), "123")
		broadcast_util.RegisterObserver(a, "123")
		broadcast_util.Notify("123", "456", &models.Project{})
		time.Sleep(time.Duration(2) * time.Second)
	})

}
