package ws

import (
	"Scholar/internal/pkg/client/model"
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/stretchr/testify/assert"
)

var (
	wg       sync.WaitGroup
	upgrader = websocket.Upgrader{}
)

const (
	testPort = "8001"
)

func newReceiverFunc(t *testing.T, shouldReceive []model.Message, done chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		for _, expected := range shouldReceive {
			mt, rawMsg, err := c.ReadMessage()
			assert.Nil(t, err, "read message should return nil in err field")
			assert.Equal(t, websocket.TextMessage, mt, "message should be text")

			var msg model.Message
			_ = json.Unmarshal(rawMsg, &msg)

			assert.Equal(t, expected, msg)
		}

		close(done)
	}
}

func newTestReceiverServer(t *testing.T, shouldReceive []model.Message, done chan struct{}) (*httptest.Server, func()) {
	s := httptest.NewServer(newReceiverFunc(t, shouldReceive, done))
	return s, func() {
		s.Close()
	}
}

func newTestConn(testServerUrl string) *websocket.Conn {
	conn, _, _ := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(testServerUrl, "http"), nil)

	return conn
}

func Test_Writing_PositiveCase(t *testing.T) {
	source := gofakeit.UUID()
	messages := []model.Message{
		{
			MsgType: gofakeit.UUID(),
			Source:  source,
			Payload: []byte(gofakeit.UUID()),
		},
		{
			MsgType: gofakeit.UUID(),
			Source:  source,
			Payload: []byte(gofakeit.UUID()),
		},
		{
			MsgType: gofakeit.UUID(),
			Source:  source,
			Payload: []byte(gofakeit.UUID()),
		},
	}

	done := make(chan struct{})
	testServer, closer := newTestReceiverServer(t, messages, done)
	defer closer()

	conn := newTestConn(testServer.URL)
	defer conn.Close()

	msgChan := make(chan *model.Message)

	client := NewClient(source, conn, msgChan)

	for _, msg := range messages {
		err := client.SendMessage(model.MessageType(msg.MsgType), msg.Payload)
		assert.Nil(t, err, "error should be nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fail()
	case <-done:
	}
}
