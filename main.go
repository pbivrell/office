package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kjk/betterguid"
	"github.com/pbivrell/office/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

	logger := util.NewLogrusLogger()

	fields := util.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
		"ua":     r.Header.Get("User-Agent"),
	}

	websocketStatus.WithLabelValues("open").Inc()
	defer websocketStatus.WithLabelValues("close").Inc()

	logger.WithFields(fields).Infof("web socket created")
	defer logger.WithFields(fields).Infof("web socket closed")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		websocketStatus.WithLabelValues("upgrade").Inc()
		return
	}

	readClose := make(chan struct{}, 0)
	defer conn.Close()
	go func() {
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				websocketStatus.WithLabelValues("read").Inc()
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
				websocketStatus.WithLabelValues("write").Inc()
				if errors.Is(err, websocket.ErrCloseSent) {
					websocketStatus.WithLabelValues("write-close").Inc()
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
	r.Handle("/metrics", promhttp.Handler())
	r.PathPrefix("/html/").Handler(http.StripPrefix("/html/", http.FileServer(http.Dir(htmlDir))))
	log.Fatal(http.ListenAndServe(addr, r))
}

var addr, htmlDir string

var websocketStatus *prometheus.CounterVec

func init() {

	flag.StringVar(&addr, "addr", "localhost:8080", "http service address")
	flag.StringVar(&htmlDir, "html", "./html", "path to html dir")
	flag.Parse()

	websocketStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "websocket_status",
		Help:        "web socket status history",
		ConstLabels: prometheus.Labels{"addr": addr},
	}, []string{
		"status",
	})

	prometheus.MustRegister(websocketStatus)

}
