package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	pool        *redis.Pool
	redisServer = flag.String("redisServer", ":6379", "")
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	conn := pool.Get()
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
	defer conn.Close()
	result, err := conn.Do("get", "content_1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hello, %q", result)
}
func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	pool = newPool(*redisServer)

	http.HandleFunc("/redis", indexHandler)
	err := http.ListenAndServe(":8880", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
