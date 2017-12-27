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
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	if rand.Int()%5 == 0 {
		wr.WriteHeader(http.StatusInternalServerError)
	} else if rand.Int()%3 == 0 {
		wr.WriteHeader(http.StatusBadRequest)
	}
	wr.Write([]byte("pong"))
}
