package main

import (
	"context"
	"dist_sys/internal/broadcast"
	log "log/slog"
	"sort"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	s := broadcast.MakeService(n)

	n.Handle(broadcast.Broadcast, s.BroadcastHandler)
	n.Handle(broadcast.Read, s.ReadHandler)
	n.Handle(broadcast.Topology, s.TopologyHandler)
	n.Handle("gossip", s.GossipHandler)

	worker := func() {
		for {
			msg, ok := <-s.Broadcasts
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
			if _, err := n.SyncRPC(ctx, msg.Dst, broadcast.MakeGossipRequest(msg.Message)); err != nil {
				messages := s.GetMessages()
				sort.Ints(messages)
				log.With("msg", msg).With("messages", messages).Error(err.Error())
				go func() {
					<-time.After(time.Millisecond * 100)
					s.Broadcasts <- msg
				}()
			}
			cancel()
		}
	}

	for i := 0; i < 10; i++ {
		go worker()
	}

	if err := n.Run(); err != nil {
		log.With("func", "Run()").Error(err.Error())
	}
}
