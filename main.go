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
	pongWait = 5 * time.Second

	// This is our zmq source
	publisherSource "tcp://localhost:5000"
)

var upgrader = websocket.Upgrader{}

type subscription struct {
	subscribe string
}

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	defer c.Close()

	err = c.SetReadDeadline(time.Now().Local().Add(pongWait))
	if err != nil {
		log.Fatalf("Could not set read deadline %v", err)
	}

	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Local().Add(pongWait))
		return nil
	})

	sub := subscription{}

	err = c.ReadJSON(&sub)
	if err != nil {
		log.Fatalf("Could not read from client %v", err)
		return
	}
	log.Printf("subscribing to <<%s>>\n", sub.subscribe)

	subscriber, _ := zmq.NewSocket(zmq.SUB)
	subscriber.Connect(publisherSource)
	subscriber.SetSubscribe(sub.subscribe)

	defer subscriber.Close()

	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			log.Fatalf("Trouble reading from zeromq %v", err)
		}
		log.Printf("Received %v", msg)
		data := []byte(msg[1])
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Print(err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	http.HandleFunc("/ws", serve)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
