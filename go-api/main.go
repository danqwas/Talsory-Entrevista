package main

import (
	"go-api/config"
	"log"
)

func main() {

	config.LoadConfig()

	app := SetupApp()
	log.Fatal(app.Listen(":" + config.PORT))
}
