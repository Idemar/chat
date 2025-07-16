package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// videresend er en kanal som inneholder innkommende meldinger
	// som skal videresendes til de andre klientene.
	forward chan []byte

	// join  er en kanal for klienter som ønsker å bli med i rommet.
	join chan *client

	// Leave er en kanal for klienter som ønsker å forlate rommet.
	leave chan *client

	// clients holder alle nåværende klienter i dette rommet.
	clients map[*client]bool
}

// newRoom lager et nytt rom
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// blir med
			r.clients[client] = true
		case client := <-r.leave:
			// forlater
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			// videresend melding til alle klienter
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
