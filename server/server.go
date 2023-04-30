package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func do_somejob() {
	// 生成一个随机种子
	rand.Seed(time.Now().UnixNano())
	randomSleep := rand.Intn(3) + 3

	// 暂停当前 goroutine
	time.Sleep(time.Duration(randomSleep) * time.Second)
}

var finalStock = 1

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if finalStock > 0 {
			do_somejob() // some other job
			finalStock--
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	resetHandler := func(w http.ResponseWriter, r *http.Request) {
		finalStock = 1
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/reset", resetHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
