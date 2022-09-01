package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	asyncProducer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("asyn kafka: %v", err)
	}

	go func() {
		for msg := range asyncProducer.Errors() {
			log.Printf("error: %v", msg)
		}
	}()

	go func() {
		for msg := range asyncProducer.Successes() {
			log.Printf("success: %v", msg)
		}
	}()

	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			for {
				asyncProducer.Input() <- &sarama.ProducerMessage{
					Topic: "test2",
					Key:   sarama.StringEncoder(key),
					Value: sarama.ByteEncoder([]byte(fmt.Sprintf("%v -> %v", key, time.Now()))),
				}
				time.Sleep(time.Second * 5)
			}
		}(fmt.Sprintf("%v", i))
	}

	go http.ListenAndServe("localhost:8090", nil)
	log.Println("all started")
	wg.Wait()

}
