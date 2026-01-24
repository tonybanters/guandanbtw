package room

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type Hub struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New_Hub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			_ = client
		case client := <-h.unregister:
			if client.room != nil {
				client.room.leave <- client
			}
		}
	}
}

func (h *Hub) Handle_Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := new_client(generate_id(), conn)

	h.register <- client

	go client.write_pump()
	go client.read_pump(h)
}

func (h *Hub) create_room() *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room := new_room(generate_room_code())
	h.rooms[room.id] = room
	go room.run()

	return room
}

func (h *Hub) get_room(id string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.rooms[id]
}

func (h *Hub) delete_room(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.rooms, id)
}

func generate_id() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generate_room_code() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)
}
