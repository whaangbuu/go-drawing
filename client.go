package main

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Client is our handler
// for client socket connection
type Client struct {
	id       string
	hub      *Hub
	color    string
	socket   *websocket.Conn
	outbound chan []byte
}

// NewClient is our constructor
// that returns an instance of Client
func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	uuID, _ := uuid.NewV4()
	uuIDStr := uuID.String()
	return &Client{
		id:       uuIDStr,
		color:    generateColor(),
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte),
	}
}

func (client *Client) read() {
	defer func() {
		client.hub.unregister <- client
	}()
	for {
		_, data, err := client.socket.ReadMessage()
		if err != nil {
			break
		}
		client.hub.onMessage(data, client)
	}
}

func (client *Client) write() {
	for {
		select {
		case data, ok := <-client.outbound:
			if !ok {
				client.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			client.socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}

func (client Client) run() {
	go client.read()
	go client.write()
}

func (client Client) close() {
	client.socket.Close()
	close(client.outbound)
}