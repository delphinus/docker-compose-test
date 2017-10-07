package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
)

func main() {
	rd, err := redis.Dial("tcp", "redis:6379")
	if err != nil {
		log.Fatalf("error on Dial: %v", err)
	}

	http.HandleFunc("/", handler(rd))
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func handler(rd redis.Conn) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := rd.Do("incr", "counter")
		if err != nil {
			writeError(w, err)
			return
		}
		res, ok := res.(int64)
		if !ok {
			writeError(w, errors.New("invalid value"))
			return
		}
		w.Write([]byte(fmt.Sprintf("counter: %d", res)))
	})
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
