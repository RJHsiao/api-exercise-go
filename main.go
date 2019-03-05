package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RJHsiao/api-exercise-go/config"
	"github.com/RJHsiao/api-exercise-go/database"
	"github.com/RJHsiao/api-exercise-go/routes"
)

func main() {
	config := config.GetConfig()

	err := database.Connect()
	if err != nil {
		log.Fatal("Connect database failed:", err)
	}
	log.Println("Database connected.")
	defer database.Disconnect()

	log.Printf("Listening http://%s:%v\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%v", config.Host, config.Port), routes.GetRouter()))
}
