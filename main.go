package main

import (
	"net/http"
	"time"
	"sync"
	"github.com/sergjeepee/dead-handler-challenge/utils"
)

var conf = utils.ParseConfigs()
var result = utils.InitResult()
var wg = new(sync.WaitGroup)
var mtx = new(sync.Mutex)

func main() {
	utils.FancyIntro(conf)

	start := time.Now()

	wg.Add(int(conf.Iterations))
	for i := 0; i < conf.Iterations; i++ {
		go sendAndHandle()
	}
	wg.Wait()

	result.TotalSyncElapsed = utils.Millis(time.Now().Sub(start))

	utils.PrintResults(conf, result)
}

func sendAndHandle() {
	start := time.Now()
	resp, err := http.Get(conf.Url)
	if err != nil {
		panic(err.Error())
	}
	elapsed := time.Now().Sub(start)
	result.HandleAnswerDuration(elapsed)

	httpGenCode := resp.Status[:1]
	mtx.Lock()
	switch httpGenCode {
	case "2":
		result.OkCount++
	case "3":
		result.RedirectCount++
	case "4":
		result.ClientErrCount++
	case "5":
		result.ServerErrorCount++
	default:
		panic("Unexpected http response code")
	}
	mtx.Unlock()
	wg.Done()
}
