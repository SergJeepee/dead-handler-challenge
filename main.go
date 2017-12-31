package main

import (
	"github.com/sergjeepee/dead-handler-challenge/utils"
	"github.com/sergjeepee/dead-handler-challenge/sender"
	"github.com/sergjeepee/dead-handler-challenge/model"
	"github.com/sergjeepee/dead-handler-challenge/constants"
	"github.com/spf13/viper"
	"strings"
	"math"
	"fmt"
	"strconv"
)

var (
	conf   = parseAndValidateConfigs()
	result = initResult()
)

func main() {
	fancyIntro(conf)

	sender.MakeRun(conf, &result)

	printResults(conf, result)
}

func fancyIntro(conf model.Conf) {
	fmt.Println("Welcome to Dead Handler Chanllenge! Now we're going to burn your http handler")
	fmt.Println(`Configs to be used:
	url: ` + conf.Url + `
	iterations: ` + strconv.Itoa(conf.Iterations) + `
	method: ` + conf.Method + `
	payload: ` + conf.Payload + `
	content-type: ` + conf.ContentType)
	fmt.Println()
}

func printResults(conf model.Conf, result model.Result) {
	fmt.Println("=============== RESULTS ===============")
	fmt.Print(`Total requests sent: ` + strconv.Itoa(conf.Iterations) + `
	Total elapsed time: ` + strconv.Itoa(result.TotalAsyncElapsed) + `
	Min response time: ` + strconv.Itoa(result.Min) + `
	Max response time: ` + strconv.Itoa(result.Max) + `
	Average response time: ` + strconv.Itoa(result.Average) + `
	Total responses: ` + strconv.Itoa(result.Responses) + `
	2** responses: ` + strconv.Itoa(result.OkCount) + `
	3** responses: ` + strconv.Itoa(result.RedirectCount) + `
	4** responses: ` + strconv.Itoa(result.ClientErrCount) + `
	5** responses: ` + strconv.Itoa(result.ServerErrorCount))

	fmt.Println("\n\nPress any key to exit")
	var holder string
	fmt.Scanln(&holder)
}

func parseAndValidateConfigs() model.Conf {
	utils.InitViper()
	conf := model.Conf{
		Iterations:  viper.GetInt(constants.IterationConfName),
		PoolSize:    viper.GetInt(constants.PoolSizeConfName) + 2, // extra for main & progress monitor
		Url:         utils.AddHttpPrefix(viper.GetString(constants.UrlConfName)),
		Method:      strings.ToUpper(viper.GetString(constants.MethodConfName)),
		Payload:     viper.GetString(constants.PayloadConfName),
		ContentType: viper.GetString(constants.ContentTypeConfName),
	}
	if conf.Iterations <= 0 {
		panic("Negative '" + constants.IterationConfName + "' value")
	}
	if !utils.InList(conf.Method, constants.AllowedHttpMethods) {
		panic("Unknown, or not allowed HTTP method. Use one of: " +
			strings.Join(constants.AllowedHttpMethods, ", "))
	}
	return conf
}

func initResult() model.Result {
	return model.Result{Min: math.MaxInt64}
}
