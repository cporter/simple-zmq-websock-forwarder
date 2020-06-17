package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	zmq "github.com/pebbe/zmq4"
)

const (
	// Time allowed to read the next pong message from the peer.
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	// This is our zmq source
	publisherSource = "tcp://localhost:5000"
)

var upgrader = websocket.Upgrader{}

type subscription struct {
	subscribe string
}

// Client is a type
type Client struct {
	conn *websocket.Conn
}

func (c *Client) readPump() {
	defer func() {
		log.Printf("read pump closing")
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("Breaking from reader")
			break
		}
		log.Printf("received %v", message)
	}
}

func (c *Client) writePump() {
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect(publisherSource)
	// subscribe to everything
	subscriber.SetSubscribe("")

	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			log.Fatalf("Trouble reading from zeromq %v", err)
		}
		log.Printf("Received %v", msg)
		data := []byte(msg[1])
		err = c.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("%v\n", err)
			break
		}
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	client := Client{c}
	go client.readPump()
	go client.writePump()
}

func main() {
	http.HandleFunc("/ws", serve)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
