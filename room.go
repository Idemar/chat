package main

import (
	"chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	// videresend er en kanal som inneholder innkommende meldinger
	// som skal videresendes til de andre klientene.
	forward chan *message

	// join  er en kanal for klienter som ønsker å bli med i rommet.
	join chan *client

	// Leave er en kanal for klienter som ønsker å forlate rommet.
	leave chan *client

	// clients holder alle nåværende klienter i dette rommet.
	clients map[*client]bool

	// tracer vil motta trace informasjon om aktiviter i rommet
	tracer trace.Tracer

	// avatar er hvordan avatarinformasjon innhentes.
	avatar Avatar
}

// newRoom lager et nytt rom som er klart til bruk
func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
		avatar:  avatar,
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// blir med
			r.clients[client] = true
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			// forlater
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			// videresend melding til alle klienter
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- sent to client")
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
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
