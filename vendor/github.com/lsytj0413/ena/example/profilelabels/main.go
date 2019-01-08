package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"
	"time"
)

func main() {
	ctx := context.Background()
	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 10000000; i++ {
			labels := pprof.Labels("handler", "hello")
			pprof.Do(ctx, labels, func(ctx context.Context) {
				generate(1, 10)
			})
		}
		fmt.Println("handler:hello done")
	}()

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 10000000; i++ {
			labels := pprof.Labels("handler2", "hello")
			pprof.Do(ctx, labels, func(ctx context.Context) {
				generate(1, 10)
			})
		}
		fmt.Println("handler2:hello done")
	}()

	log.Fatal(http.ListenAndServe("localhost:5555", nil))
}

func generate(duration int, usage int) string {
	s := fmt.Sprintf("duration: %d", duration)
	for i := 0; i < usage; i++ {
		s += s
	}
	return s
}
