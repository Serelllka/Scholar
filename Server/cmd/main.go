package main

import (
	Node "Scholar/net"
	"time"
)

func main() {
	cfg := Node.NodeConfig{
		Version:    "Scholar1.0-alpha",
		ListenAddr: ":3000",
	}
	node := Node.NewNode(cfg)
	go node.Start()

	time.Sleep(time.Millisecond * 250)
	remoteCfg := Node.NodeConfig{
		Version:    "Scholar1.0-alpha",
		ListenAddr: ":4000",
	}
	remoteNode := Node.NewNode(remoteCfg)
	go remoteNode.Start()
	_ = remoteNode.Connect(":3000")

	select {}
}
