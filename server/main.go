package main

import (
	"guandanbtw/room"
	"log"
	"net/http"
)

func main() {
	hub := room.New_Hub()
	go hub.Run()

	http.HandleFunc("/ws", hub.Handle_Websocket)

	http.Handle("/", http.FileServer(http.Dir("../client/dist")))

	log.Println("server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
