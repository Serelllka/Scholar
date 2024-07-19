package model

type Message struct {
	MsgType string `json:"type"`
	Source  string `json:"source"`
	Payload []byte `json:"payload"`
}
