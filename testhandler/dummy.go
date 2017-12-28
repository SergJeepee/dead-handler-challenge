package main

import (
	"net/http"
	"time"
	"math/rand"
	"fmt"
)

func main() {
	fmt.Println("Dummy starting on 8080...")
	http.HandleFunc("/", get)
	http.ListenAndServe(":8080", nil)
}

func get(wr http.ResponseWriter, rq *http.Request) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)))
	rnd := rand.Int()
	//set randomly other than 200
	if rnd%5 == 0 {
		wr.WriteHeader(http.StatusInternalServerError)
	} else if rnd%6 == 0 {
		wr.WriteHeader(http.StatusBadRequest)
	} else if rnd%7 == 0 {
		wr.WriteHeader(http.StatusUseProxy)
	}
	wr.Write([]byte("pong"))
}
