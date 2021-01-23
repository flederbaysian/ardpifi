package web

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Web struct {
	hub *Hub
}

func InitWeb(webAddr string) *Web {
	log.Debugf("InitWeb")

	hub := newHub()
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/web/index.html")
	})
	http.HandleFunc("/jquery-1.11.1.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/web/jquery-1.11.1.js")
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	go func() {
		err := http.ListenAndServe(webAddr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	return &Web{
		hub: hub,
	}
}

func (w *Web) Render(msg string) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Print(err)
		return
	}

	go func() { w.hub.render <- b }()
}