package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"math"
	"github.com/spf13/viper"
	"gopkg.in/cheggaaa/pb.v1"
	"sync"
)

const (
	millisMultilayer  = 1e6
	iterationConfName = "iterations"
	urlConfName       = "url"
	methodConfName    = "method"
	payloadConfName   = "payload"
)

type Conf struct {
	Iterations int
	Url        string
	Method     string
	Payload    string
}

type Result struct {
	Responses         uint
	OkCount           uint
	RedirectCount     uint
	ClientErrCount    uint
	ServerErrorCount  uint
	Min               uint
	Max               uint
	Average           uint
	TotalSyncElapsed  uint // needed for average value calculation
	TotalAsyncElapsed uint // total time elapsed to send and handle requests through all iterations
}

type TimeStatisticOwner interface {
	HandleAnswerDuration(d time.Duration)
}

func (r *Result) HandleAnswerDuration(d time.Duration) {
	r.Responses++
	r.TotalSyncElapsed += Millis(d)
	r.Average = r.TotalSyncElapsed / uint(r.Responses) //it's ok to lose floating tail
	r.Max = max(r.Max, Millis(d))
	r.Min = min(r.Min, Millis(d))
}

func FancyIntro(conf Conf) {
	fmt.Println("Welcome to Dead Handler Chanllenge! Now we're going to burn your http handler")
	fmt.Println(`Configs to be used:
	url: ` + conf.Url + `
	iterations: ` + strconv.Itoa(conf.Iterations) + `
	method: ` + conf.Method + `
	payload: ` + conf.Payload)
	fmt.Println()
}

func ProgressMonitor(conf Conf, result *Result, wg *sync.WaitGroup) {
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

func PrintResults(conf Conf, result Result) {
	fmt.Println("=============== RESULTS ===============")
	fmt.Print(`Total requests sent: ` + strconv.Itoa(conf.Iterations) + `
	Total elapsed time: ` + uint2string(result.TotalAsyncElapsed) + `
	Min response time: ` + uint2string(result.Min) + `
	Max response time: ` + uint2string(result.Max) + `
	Average response time: ` + uint2string(result.Average) + `
	Total responses: ` + uint2string(result.Responses) + `
	2** responses: ` + uint2string(result.OkCount) + `
	3** responses: ` + uint2string(result.RedirectCount) + `
	4** responses: ` + uint2string(result.ClientErrCount) + `
	5** responses: ` + uint2string(result.ServerErrorCount))

	fmt.Println("\n\nPress any key to exit")
	var holder string
	fmt.Scanln(&holder)
}

func ParseConfigs() Conf {
	initViper()
	return Conf{
		Iterations: viper.GetInt(iterationConfName),
		Url:        handleUrl(viper.GetString(urlConfName)),
		Method:     viper.GetString(methodConfName),
		Payload:    viper.GetString(payloadConfName),
	}
}

func InitResult() Result {
	return Result{Min: math.MaxInt64}
}

//since time.Duration doesn't have Millis converter
func Millis(d time.Duration) uint {
	return uint(d.Nanoseconds() / millisMultilayer)
}

// since math.Min deals with float64 only
func min(first uint, second uint) uint {
	//where the heck is ternary operator in Go?!
	if first <= second {
		return first
	}
	return second
}

// since math.Max also deals with float64 only
func max(first uint, second uint) uint {
	if first >= second {
		return first
	}
	return second
}

func uint2string(i uint) string {
	return strconv.Itoa(int(i))
}

func handleUrl(s string) string {
	if !strings.HasPrefix(s, "http://") {
		return "http://" + s
	}
	return s
}

func initViper() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./src/github.com/sergjeepee/dead-handler-challenge/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
