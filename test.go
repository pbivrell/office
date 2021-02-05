// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type dataPair struct {
	message []byte
	mt      int
}

func (c *coordinator) subscribe(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

}

type subscriber struct {
	chans map[string]chan dataPair
	*sync.Mutex
}

const (
	BufferSize = 100
)

func (s *subscriber) subscribe() (chan dataPair, func()) {
	s.Lock()
	defer s.Unlock()

	index := betterguid.New()

	c := make(chan dataPair, BufferSize)

	s.chans[index] = c

	fmt.Println("inserted", index)

	for _, neighbor := range s.chans {
		data, _ := json.Marshal(struct {
			Message string `json:"message"`
		}{
			Message: "solicit",
		})
		neighbor <- dataPair{
			message: data,
			mt:      websocket.TextMessage,
		}
		break
	}

	return c, func() {
		s.Lock()
		defer s.Unlock()

		fmt.Println("removing", index)
		delete(s.chans, index)
	}
}

type coordinator struct {
	b chan dataPair
	s *subscriber
}

func (c *coordinator) process() {
	for {
		select {
		case data := <-c.b:
			c.s.Lock()
			for _, dChan := range c.s.chans {
				dChan <- data
			}
			c.s.Unlock()
		}
	}

}

func (c *coordinator) broadcast(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	readClose := make(chan struct{}, 0)
	defer conn.Close()
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			c.b <- dataPair{
				mt:      mt,
				message: message,
			}
		}
		readClose <- struct{}{}
	}()
	dChan, unsubscribe := c.s.subscribe()
	defer unsubscribe()

	for {
		select {
		case data := <-dChan:
			err = conn.WriteMessage(data.mt, data.message)
			if err != nil {
				log.Println("write:", err)
				if errors.Is(err, websocket.ErrCloseSent) {
					return
				}
			}
		case <-readClose:
			return

		}
	}

}

func main() {

	bChan := make(chan dataPair, 2000)

	coord := coordinator{
		b: bChan,
		s: &subscriber{
			chans: make(map[string]chan dataPair),
			Mutex: &sync.Mutex{},
		},
	}

	go coord.process()

	flag.Parse()
	log.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/echo", coord.broadcast)
	r.HandleFunc("/listen", coord.subscribe)
	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir("./html/"))))
	log.Fatal(http.ListenAndServe(*addr, r))
}
