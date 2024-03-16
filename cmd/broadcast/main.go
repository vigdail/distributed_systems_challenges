package main

import (
	"dist_sys/internal/broadcast"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	s := broadcast.MakeService(n)

	n.Handle(broadcast.Broadcast, s.BroadcastHandler)
	n.Handle(broadcast.Read, s.ReadHandler)
	n.Handle(broadcast.Topology, s.TopologyHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
