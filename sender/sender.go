package sender

import (
	"github.com/sergjeepee/dead-handler-challenge/utils"
	"github.com/sergjeepee/dead-handler-challenge/model"
	"time"
	"sync/atomic"
	"strings"
	"io"
	"io/ioutil"
	"sync"
	"net/http"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
)

var (
	wg     = new(sync.WaitGroup)
	mtx    = new(sync.Mutex)
	cli    = &http.Client{}
	conf   model.Conf
	result *model.Result
)

func MakeRun(c model.Conf, r *model.Result) {
	conf = c
	result = r

	start := time.Now()

	// Wait for both request pool and progress monitor are done
	wg.Add(2)

	go progressMonitor(conf, result, wg)

	pool := make(chan byte, conf.PoolSize)
	done := make(chan bool)
	totalJobsCount := new(int64)

	// fill the pool for the first time
	for i := 0; i < utils.Min(conf.PoolSize, conf.Iterations); i++ {
		pool <- 1
	}

	jobLoop:
	for {
		select {
		case <-done:
			if atomic.LoadInt64(totalJobsCount) < int64(conf.Iterations) {
				pool <- 1
			} else {
				wg.Done()
				break jobLoop
			}
		case <-pool:
			atomic.AddInt64(totalJobsCount, 1)
			go func() {
				sendAndHandle()
				done <- true
			}()
		}
	}
	wg.Wait()
	result.TotalAsyncElapsed = utils.Millis(time.Now().Sub(start))
}

func sendAndHandle() {
	start := time.Now()

	req, err := http.NewRequest(conf.Method, conf.Url, strings.NewReader(conf.Payload))
	utils.HandleError(err)
	req.Header.Set("content-type", conf.ContentType)
	resp, respErr := cli.Do(req)
	utils.HandleError(respErr)
	elapsed := time.Now().Sub(start)

	//we have to readout and close response body to ensure http connection going to be freed and can be reused
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()

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
}

func progressMonitor(conf model.Conf, result *model.Result, wg *sync.WaitGroup) {
	fmt.Println("In progress:  ")
	bar := pb.New(conf.Iterations)
	bar.Start()
	for !bar.IsFinished() {
		if int(bar.Get()) == conf.Iterations {
			bar.Finish()
		}
		bar.Add(int(int64(result.Responses) - bar.Get()))
	}
	fmt.Println()
	wg.Done()
}
