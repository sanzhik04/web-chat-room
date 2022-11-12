package models


import (
	"fmt"
	"net/http"
	"strings"
	"math/rand"
	"time"
	"log"

	"github.com/gorilla/websocket"
)

type Room struct {
	Users    map[string]*User
	Messages chan *Message
	Join     chan *User
	Leave    chan *User
}


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println("%s %s%s %v\n", r.Method, r.Host, r.RequestURI, r.Proto)
		return r.Method == http.MethodGet
	},
}

func (c *Room) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error While Connecting: ", err)
	}

	keys := r.URL.Query()
	username := keys.Get("username")
	if strings.TrimSpace(username) == "" {
		username = fmt.Sprintf("Anonymous+%d",rand.New(rand.NewSource(time.Now().UnixNano())).Int63())
	}

	user := &User{
		Username: username,
		Conn:     conn,
		Global:   c,
	}

	c.Join <- user

	user.Read()
}

func (c *Room) Run() {
	for {
		select {
		case user := <-c.Join:
			c.add(user)
		case message := <-c.Messages:
			c.broadcast(message)
		case user := <-c.Leave:
			c.disconnect(user)
		}
	}
}

func (c *Room) add(user *User) {
	if _, ok := c.Users[user.Username]; !ok {
		c.Users[user.Username] = user

		body := fmt.Sprintf("%s connected to the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))
	}
}

func (c *Room) broadcast(message *Message) {
	fmt.Println("Broadcast message: %v\n", message)
	for _, user := range c.Users {
		user.Write(message)
	}
}

func (c *Room) disconnect(user *User) {
	if _, ok := c.Users[user.Username]; ok {
		defer user.Conn.Close()
		delete(c.Users, user.Username)

		body := fmt.Sprintf("%s disconnected from the chat", user.Username)
		c.broadcast(NewMessage(body, "Server"))
	}
}

func Start(port string) {

	fmt.Println("Chat listening on http://localhost%s\n", port)

	c := &Room{
		Users:    make(map[string]*User),
		Messages: make(chan *Message),
		Join:     make(chan *User),
		Leave:    make(chan *User),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Go WebChat!"))
	})

	http.HandleFunc("/chat", c.Handler)

	go c.Run()

	log.Fatal(http.ListenAndServe(port, nil))
}