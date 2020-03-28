package main

import (
	"log"
	"time"

	"github.com/gollira/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()

			msg.Name = "anonymous"
			if nickname, ok := c.userData["nickname"].(string); ok {
				msg.Name = nickname
			}
			if avatarURL, ok := c.userData["avatar_url"].(string); ok {
				msg.AvatarURL = avatarURL
			}

			c.room.forward <- msg
		} else {
			log.Fatal("cannot read json ", err)
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
