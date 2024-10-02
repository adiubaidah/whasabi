package app

import (
	"adiubaidah/adi-bot/helper"
	"adiubaidah/adi-bot/model"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WsClient struct {
	User string
	Conn *websocket.Conn
}

type WebSocketHub struct {
	Clients map[string]*WsClient
	mu      sync.Mutex
}

var WsHub = &WebSocketHub{
	Clients: make(map[string]*WsClient),
}

func (hub *WebSocketHub) ServeWebSocket(w http.ResponseWriter, r *http.Request, phone string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	helper.PanicIfError("Error while upgrade connection", err)
	defer conn.Close()

	hub.mu.Lock()
	hub.Clients[phone] = &WsClient{
		User: phone,
		Conn: conn,
	}
	conn.WriteJSON(&model.WebResponse{
		Code:   200,
		Status: "success",
		Data:   "Connected to WebSocket",
	})
	hub.mu.Unlock()

	// Maintain connection and listen for incoming messages
	for {
		_, _, err := conn.ReadMessage() //will return error if connection is closed
		if err != nil {
			fmt.Println("WebSocket Read Error:", err)
			hub.mu.Lock()
			delete(hub.Clients, phone)
			hub.mu.Unlock()
			break
		}
	}
}

func (hub *WebSocketHub) SendMessage(phone string, response model.WebResponse) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if client, exists := hub.Clients[phone]; exists {
		err := client.Conn.WriteJSON(&response)
		if err != nil {
			fmt.Println("WebSocket Write Error:", err)
			delete(hub.Clients, phone)
		}
	}
}
