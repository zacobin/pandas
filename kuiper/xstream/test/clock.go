package test

import (
	"github.com/benbjohnson/clock"
	"github.com/cloustone/pandas/kuiper"
)

func ResetClock(t int64) {
	mock := clock.NewMock()
	mock.Set(util.TimeFromUnixMilli(t))
	util.Clock = mock
}

func GetMockClock() *clock.Mock {
	return util.Clock.(*clock.Mock)
}
