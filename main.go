package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	//	"github.com/gorilla/websocket"
)

var clientUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 30000,
}

var boatUpgrader = websocket.Upgrader{
	ReadBufferSize:  30000,
	WriteBufferSize: 1024,
}

func serveBroadcaster(w http.ResponseWriter, r *http.Request) *Hub {
	conn, err := clientUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return nil
	}

	h := newHub()
	b := &Broadcaster{}
	b.Make(conn, h)
	h.broadcaster = b

	return h
}

func serveViewer(w http.ResponseWriter, r *http.Request, h *Hub) *Viewer {
	conn, err := clientUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return nil
	}

	v := &Viewer{}
	v.Make(conn, h)
	v.hub.register <- v

	return v
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.Handle("/static/", http.FileServer(http.Dir("static/")))

	//	http.Handle("/", server)
	var h *Hub

	http.HandleFunc("/viewer", func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			v := serveViewer(w, r, h)

			go readMessages(v)
			go writeMessages(v)
		}
	})

	http.HandleFunc("/broadcaster", func(w http.ResponseWriter, r *http.Request) {
		if h != nil {
			h = serveBroadcaster(w, r)

			go h.run()
			go readMessages(h.broadcaster)
			go writeMessages(h.broadcaster)
		}
	})

	http.ListenAndServe(":3000", nil)
}
