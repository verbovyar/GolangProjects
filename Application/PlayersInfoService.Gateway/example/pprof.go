package main

import (
	"expvar"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"sync"
)

func testFunc(c *Counter) {
	ch := make(chan interface{})

	for i := 0; i < 10; i++ {
		go func() {
			c.Inc()
			ch <- 1
		}()
	}
}

type Counter struct {
	cnt int
	m   *sync.RWMutex
}

func (c *Counter) Inc() {
	c.m.Lock()
	defer c.m.Unlock()
	c.cnt++
}

func (c *Counter) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return strconv.FormatInt(int64(c.cnt), 10)
}

type Goroutines struct {
}

func (g *Goroutines) String() string {
	return strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
}

func main() {
	g := &Goroutines{}
	c := &Counter{m: &sync.RWMutex{}}
	expvar.Publish("Goroutines", g)
	expvar.Publish("Counter", c)
	testFunc(c)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return
	}
}
