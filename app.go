package main

import (
	"github.com/gorilla/mux"
	"github.com/snikiforov55/back-session/session"
	"log"
	"net/http"
)

func StartWebServer(port string, router *mux.Router) {
	log.Println("Starting HTTP service at " + port)
	http.Handle("/", router)
	err := http.ListenAndServe(":"+port, nil) // Goroutine will block here
	if err != nil {
		log.Println("An error occured starting HTTP listener at port " + port)
		log.Println("Error: " + err.Error())
	}
}

func main() {
	server := session.NewServer()
	StartWebServer("8090", server.Router)
}
