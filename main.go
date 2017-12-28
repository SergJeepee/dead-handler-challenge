package main

import (
	"net/http"
	"time"
	"sync"
	"github.com/sergjeepee/dead-handler-challenge/utils"
)

var (
	conf   = utils.ParseConfigs()
	result = utils.InitResult()
	wg     = new(sync.WaitGroup)
	mtx    = new(sync.Mutex)
)

func main() {
	utils.FancyIntro(conf)

	start := time.Now()

	// One extra for progress bar
	wg.Add(int(conf.Iterations) + 1)

	go utils.ProgressMonitor(conf, &result, wg)
	for i := 0; i < conf.Iterations; i++ {
		go sendAndHandle()
	}
	wg.Wait()

	result.TotalAsyncElapsed = utils.Millis(time.Now().Sub(start))

	utils.PrintResults(conf, result)
}

func sendAndHandle() {
	start := time.Now()
	resp, err := http.Get(conf.Url)
	if err != nil {
		panic(err)
	}
	elapsed := time.Now().Sub(start)

	mtx.Lock()
	result.HandleAnswerDuration(elapsed)
	httpGenCode := []rune(resp.Status)[0]
	switch string(httpGenCode) {
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
