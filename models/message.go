package models

import(
	"time"
	"math/rand"
)


type Message struct {
	ID     int64  `json:"id"`
	Body   string `json:"body"`
	Sender string `json:"sender"`
}

func NewMessage(body string, sender string) *Message {
	return &Message{
		ID:     rand.New(rand.NewSource(time.Now().UnixNano())).Int63(),
		Body:   body,
		Sender: sender,
	}
}