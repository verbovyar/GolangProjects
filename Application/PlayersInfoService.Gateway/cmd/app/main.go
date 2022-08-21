package main

import (
	"fmt"
	"log"
	"modules/config"
	"modules/internal/app"
)

func main() {
	conf, err := config.LoadConfig("././config")
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Println("Started cmd")

	app.Run(conf)
}
