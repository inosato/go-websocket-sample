package main

import (
	"log"
	"net/http"

	"github.com/gollira/websocket"
	"github.com/inosato/go-websocket-sample/trace"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
	avatar  Avatar
}

func newRoom(avatar Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.tracer.Trace("join!")
			r.clients[client] = true

		case client := <-r.leave:
			r.tracer.Trace("leave!")
			r.clients[client] = false

		case msg := <-r.forward:
			r.tracer.Trace("forward!")
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					r.tracer.Trace("forward ERROR")
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("cannot get authCookie: ", err)
	}

	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		log.Print("defer func()")
		r.leave <- client
	}()
	go client.write()
	client.read()
}
