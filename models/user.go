package models

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type User struct {
	Username string
	Conn     *websocket.Conn
	Global   *Room
}

func (u *User) Read() {
	for {
		if _, message, err := u.Conn.ReadMessage(); err != nil {
			fmt.Println("Error on read message:", err.Error())

			break
		} else {
			u.Global.Messages <- NewMessage(string(message), u.Username)
		}
	}

	u.Global.Leave <- u
}

func (u *User) Write(message *Message) {
	b, _ := json.Marshal(message)

	if err := u.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		fmt.Println("Error on write message:", err.Error())
	}
}