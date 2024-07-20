package ws

import (
	"Scholar/internal/pkg/client"
	"Scholar/internal/pkg/client/model"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	// sync
	recvStarter sync.Once
	mtx         sync.Mutex

	hash string

	conn    *websocket.Conn
	msgChan chan *model.Message
}

var _ client.IWebClient = &Client{}

// NewClient ...
func NewClient(hash string, conn *websocket.Conn, msgChan chan *model.Message) *Client {
	return &Client{
		hash:    hash,
		conn:    conn,
		msgChan: msgChan,
	}
}

func (c *Client) SendMessage(ctx context.Context, msgType model.MessageType, payload []byte) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	msg := model.Message{
		MsgType: string(msgType),
		Source:  c.hash,
		Payload: payload,
	}

	msgRaw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, msgRaw)
}

// StartReceiver can be started once
func (c *Client) StartReceiver(ctx context.Context) {
	c.recvStarter.Do(func() {
		c.recvLoop(ctx)
	})
}

// GetId ...
func (c *Client) GetId() string {
	return c.hash
}

func (c *Client) recvLoop(ctx context.Context) {
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				_ = c.conn.Close()
				close(c.msgChan)
			default:
				_, payload, err := c.conn.ReadMessage()
				if err != nil {
					break loop
				}

				msg := &model.Message{}
				_ = json.Unmarshal(payload, &msg)
				c.msgChan <- msg
			}
		}
	}()
}
