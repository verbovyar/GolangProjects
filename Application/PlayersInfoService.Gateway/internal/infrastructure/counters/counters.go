package counters

import (
	"expvar"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	startTime = time.Now().UTC()
)

func getGoroutines() interface{} {
	return runtime.NumGoroutine()
}

func getCpu() interface{} {
	return runtime.NumCPU()
}

func getUptime() interface{} {
	return int64(time.Since(startTime))
}

func GetCounters() {
	mux := http.NewServeMux()

	expvar.Publish("Goroutines", expvar.Func(getGoroutines))
	expvar.Publish("Uptime", expvar.Func(getUptime))
	expvar.Publish("Cpu", expvar.Func(getCpu))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err)
	}
}
