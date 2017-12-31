package model

import (
	"time"
	"github.com/sergjeepee/dead-handler-challenge/utils"
)

type Conf struct {
	Iterations  int
	PoolSize    int
	Url         string
	Method      string
	Payload     string
	ContentType string
}

type Result struct {
	Responses         int
	OkCount           int
	RedirectCount     int
	ClientErrCount    int
	ServerErrorCount  int
	Min               int
	Max               int
	Average           int
	TotalSyncElapsed  int // needed for average value calculation
	TotalAsyncElapsed int // total time elapsed to send and handle requests through all iterations
}

type TimeStatisticOwner interface {
	HandleAnswerDuration(d time.Duration)
}

func (r *Result) HandleAnswerDuration(d time.Duration) {
	r.Responses++
	millis := utils.Millis(d)
	r.TotalSyncElapsed += millis
	r.Average = r.TotalSyncElapsed / r.Responses //it's ok to lose floating tail
	r.Max = utils.Max(r.Max, millis)
	r.Min = utils.Min(r.Min, millis)
}
