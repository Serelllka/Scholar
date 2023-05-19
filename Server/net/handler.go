package net

import (
	"fmt"
	"io"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct {
}

func (h *DefaultHandler) HandleMessage(msg *Message) error {
	out, _ := io.ReadAll(msg.Payload)
	fmt.Printf("from: %s, content: %s\n", msg.From, string(out))
	return nil
}
