package net

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

type NodeConfig struct {
	Version    string
	ListenAddr string
}

type Node struct {
	NodeConfig

	handler Handler

	transport *TCPTransport
	mtx       sync.RWMutex
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
}

func NewNode(config NodeConfig) *Node {
	n := &Node{
		NodeConfig: config,
		handler:    &DefaultHandler{},
		peers:      make(map[net.Addr]*Peer),
		addPeer:    make(chan *Peer),
		delPeer:    make(chan *Peer),
		msgCh:      make(chan *Message),
	}

	tr := NewTCPTransport(n.ListenAddr)
	n.transport = tr

	tr.AddPeer = n.addPeer
	tr.DelPeer = n.delPeer

	return n
}

func (n *Node) Start() {
	go n.loop()

	logrus.WithFields(logrus.Fields{
		"port": n.ListenAddr,
	}).Info("started new game server")

	err := n.transport.ListenAndAccept()
	if err != nil {
		panic(err)
	}
}

func (n *Node) handleConn(p *Peer) {
	buf := make([]byte, 1024)
	for {
		count, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		n.msgCh <- &Message{
			Payload: bytes.NewReader(buf[:count]),
			From:    p.conn.RemoteAddr(),
		}
	}

	n.delPeer <- p
}

func (n *Node) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	n.addPeer <- peer
	return peer.Send([]byte(n.Version))
}

func (n *Node) loop() {
	for {
		select {
		case peer := <-n.addPeer:
			go peer.ReadLoop(n.msgCh)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Player was connected")

			n.peers[peer.conn.RemoteAddr()] = peer
		case peer := <-n.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Player was disconnected")

			delete(n.peers, peer.conn.RemoteAddr())
		case msg := <-n.msgCh:
			if err := n.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
