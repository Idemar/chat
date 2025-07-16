package main

import "github.com/gorilla/websocket"

// klienten representerer en enkelt chatbruker
type client struct {
	// socket er web-socketen for denne klienten
	socket *websocket.Conn
	// send er kanalen som meldinger sendes på
	send chan []byte
	// room er rommet denne klienten chatter i
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
