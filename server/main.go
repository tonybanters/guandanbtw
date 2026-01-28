package main

import (
	"guandanbtw/room"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	hub := room.New_Hub()
	go hub.Run()

	http.HandleFunc("/ws", hub.Handle_Websocket)

	http.Handle("/", http.FileServer(http.Dir("../client/dist")))

	log.Println("server starting on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
