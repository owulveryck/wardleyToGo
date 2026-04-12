package collab

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/coder/websocket"
)

// Client represents a single WebSocket connection to a session.
type Client struct {
	id        string
	name      string
	color     string
	mode      string // "rw" or "ro"
	session   *Session
	conn      *websocket.Conn
	send      chan []byte
	ip        string       // source IP for connection tracking
	onClose   func()       // called once when ReadPump exits
	msgBucket tokenBucket  // per-client message rate limiter
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 54 * time.Second
	maxMsgSize = 64 * 1024
	sendBufLen = 256
)

// ReadPump reads messages from the WebSocket and dispatches them to the session.
func (c *Client) ReadPump(ctx context.Context) {
	defer func() {
		if c.onClose != nil {
			c.onClose()
		}
		c.session.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	c.conn.SetReadLimit(maxMsgSize)

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("client %s: bad message: %v", c.id, err)
			continue
		}

		// Per-client message rate limiting
		if !c.msgBucket.allow() {
			errMsg, _ := MarshalMessage(MsgError, ErrorPayload{Message: "rate limited"})
			select {
			case c.send <- errMsg:
			default:
			}
			continue
		}

		switch msg.Type {
		case MsgHello:
			var hello HelloPayload
			if err := json.Unmarshal(msg.Payload, &hello); err == nil {
				if len(hello.Name) > maxNameLen {
					hello.Name = hello.Name[:maxNameLen]
				}
				c.name = hello.Name
			}
		case MsgOp:
			var op OpPayload
			if err := json.Unmarshal(msg.Payload, &op); err == nil {
				c.session.HandleOp(c, &op)
			}
		case MsgCursor:
			var cursor CursorPayload
			if err := json.Unmarshal(msg.Payload, &cursor); err == nil {
				c.session.HandleCursor(c, &cursor)
			}
		}
	}
}

// WritePump writes messages from the send channel to the WebSocket.
func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				return
			}
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, websocket.MessageText, message)
			cancel()
			if err != nil {
				return
			}
		case <-ticker.C:
			pingCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(pingCtx)
			cancel()
			if err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
