package main

import (
	"go-api/config"
	"log"
)

func main() {

	config.LoadConfig()

	app := SetupApp()

	log.Printf("🚀 API de Go (Servicio Principal) escuchando en el puerto :%s", config.PORT)
	log.Fatal(app.Listen(":" + config.PORT))
}
