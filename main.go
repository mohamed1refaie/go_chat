package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type User struct {
	Name   string
	Output chan Message
}

type Message struct {
	Username string
	Text     string
}

type ChatServer struct {
	Users map[string]User
	Join  chan User
	Leave chan User
	Input chan Message
}

func (cs *ChatServer) Run() {
	for {
		select {
		case user := <-cs.Join:
			cs.Users[user.Name] = user
			go func() {
				cs.Input <- Message{
					Username: "SYSTEM",
					Text:     fmt.Sprintf("%s Joined", user.Name),
				}
			}()
		case user := <-cs.Leave:
			delete(cs.Users, user.Name)
			go func() {
				cs.Input <- Message{
					Username: "SYSTEM",
					Text:     fmt.Sprintf("%s Left", user.Name),
				}
			}()
		case msg := <-cs.Input:
			for _, user := range cs.Users {
				user.Output <- msg
			}
		}
	}
}

var chatServer = &ChatServer{
	Users: make(map[string]User),
	Join:  make(chan User),
	Leave: make(chan User),
	Input: make(chan Message),
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	var username string
	_, usr, err := conn.ReadMessage()
	if err != nil {
		return
	}
	username = string(usr)
	user := User{
		Name:   string(username),
		Output: make(chan Message),
	}
	chatServer.Join <- user
	defer func() {
		chatServer.Leave <- user
	}()

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			chatServer.Input <- Message{
				Username: user.Name,
				Text:     string(msg),
			}
		}
	}()
	for msg := range user.Output {

		if err := conn.WriteMessage(1, []byte(msg.Text)); err != nil {
			return
		}
	}
}

func main() {

	go chatServer.Run()
	http.HandleFunc("/ws", handleRequest)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.ListenAndServe(":8080", nil)

}
