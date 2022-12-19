package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"nhooyr.io/websocket"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	t := time.NewTicker(time.Second * 15)
	ws := &wsWriter{
		C: t.C,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.handler)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	log.Println("listening on " + port + ". Use the PORT env var to change it")
	err := http.ListenAndServe(":"+port, mux)
	log.Fatal(err)
}

type wsWriter struct {
	C <-chan time.Time
}

// This handler demonstrates how to correctly handle a write only WebSocket connection.
// i.e you only expect to write messages and do not expect to read any messages.
func (ws *wsWriter) handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
	defer cancel()

	ctx = c.CloseRead(ctx)

	svg, err := ioutil.ReadFile("sample.svg")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			c.Close(websocket.StatusNormalClosure, "")
			return
		case <-ws.C:
			log.Println("tick")
			w, err := c.Writer(ctx, websocket.MessageText)
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Fprintf(w, "%s", svg)
			log.Println("sent")
			w.Close()
		}
	}
}
