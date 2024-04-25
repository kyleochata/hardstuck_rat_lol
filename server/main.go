package main

import (
	"fmt"
	"hardstuck_rat_lol_server/routes"
	"log"
	"net/http"
	//_ "github.com/joho/godotenv/autoload"
)

func main() {
	router := http.NewServeMux()
	routes.RouteSetup(router)
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
	}
	log.Println("Listening on :8080")

}
