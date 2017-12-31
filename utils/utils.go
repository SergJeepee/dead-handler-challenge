package utils

import (
	"strings"
	"time"
	"github.com/spf13/viper"
	"github.com/sergjeepee/dead-handler-challenge/constants"
)

//since time.Duration doesn't have Millis converter
func Millis(d time.Duration) int {
	return int(d.Nanoseconds() / constants.MillisMultilayer)
}

// since math.Min deals with float64 only
func Min(first int, second int) int {
	//where the heck is ternary operator in Go?!
	if first <= second {
		return first
	}
	return second
}

// since math.Max also deals with float64 only
func Max(first int, second int) int {
	if first >= second {
		return first
	}
	return second
}

func InList(v string, list []string) bool {
	for _, val := range list {
		if v == val {
			return true
		}
	}
	return false
}

func AddHttpPrefix(s string) string {
	if !strings.HasPrefix(s, "http://") {
		return "http://" + s
	}
	return s
}

func InitViper() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./src/github.com/sergjeepee/dead-handler-challenge/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	HandleError(err)
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
