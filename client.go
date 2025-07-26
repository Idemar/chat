package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// klienten representerer en enkelt chatbruker
type client struct {
	// socket er web-socketen for denne klienten
	socket *websocket.Conn
	// send er kanalen som meldinger sendes på
	send chan *message
	// room er rommet denne klienten chatter i
	room *room
	// userData inneholder informasjon om brukeren
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
