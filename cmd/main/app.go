package main

import (
	"log"
	"net"
	"net/http"
	"restapi-lesson/internal/user"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	log.Println("create roter")
	router := httprouter.New()
	log.Println("register user handler")
	handler := user.NewHandler()
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	log.Println("start application")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("server is listening port 0.0.0.0:8080")
	log.Fatalln(server.Serve(listener))
}
