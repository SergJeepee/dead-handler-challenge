package main

import (
	"net/http"
	"math/rand"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	fmt.Println("Dummy starting on 8080...")
	http.HandleFunc("/", yummy)
	http.ListenAndServe(":8080", nil)
}

func yummy(wr http.ResponseWriter, rq *http.Request) {
	fmt.Println("Received", rq.Method, "; payload:", readPayloadInfo(rq), "content type:", rq.Header.Get("content-type"))
	time.Sleep(time.Millisecond*10 + time.Millisecond*time.Duration(rand.Intn(100)))

	rnd := rand.Int()
	//set randomly other than 200
	if rnd%15 == 0 {
		wr.WriteHeader(http.StatusInternalServerError)
	} else if rnd%16 == 0 {
		wr.WriteHeader(http.StatusBadRequest)
	} else if rnd%17 == 0 {
		wr.WriteHeader(http.StatusUseProxy)
	}
	wr.Write([]byte("pong"))
}

func readPayloadInfo(rq *http.Request) string {
	body, err := ioutil.ReadAll(rq.Body)
	defer rq.Body.Close()
	if err != nil {
		return "Read payload error:" + err.Error()
	}
	payload := string(body)
	if payload == "" {
		return "[empty]"
	}
	return payload
}
