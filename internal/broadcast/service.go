package broadcast

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
	"sync"
)

const Broadcast = "broadcast"
const BroadcastOk = "broadcast_ok"
const Read = "read"
const ReadOk = "read_ok"
const Topology = "topology"
const TopologyOk = "topology_ok"

type BroadcastRequest struct {
	Type    string `json:"type"`
	Message int    `json:"message"`
}

type BroadcastResponse struct {
	Type string `json:"type"`
}

func MakeBroadcastResponse() BroadcastResponse {
	return BroadcastResponse{BroadcastOk}
}

type ReadResponse struct {
	Type     string `json:"type"`
	Messages []int  `json:"messages"`
}

func MakeReadResponse(messages []int) ReadResponse {
	return ReadResponse{ReadOk, messages}
}

type TopologyRequest struct {
	Topology map[string][]string `json:"topology"`
}

type TopologyResponse struct {
	Type string `json:"type"`
}

func MakeTopologyResponse() TopologyResponse {
	return TopologyResponse{TopologyOk}
}

type Service struct {
	node      *maelstrom.Node
	neighbors []string

	messagesMu sync.Mutex
	messages   map[int]bool
}

func MakeService(node *maelstrom.Node) Service {
	return Service{
		node, make([]string, 0),
		sync.Mutex{},
		make(map[int]bool),
	}
}

func (s *Service) BroadcastHandler(msg maelstrom.Message) error {
	var request BroadcastRequest
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}
	go func() {
		if err := s.node.Reply(msg, MakeBroadcastResponse()); err != nil {
			log.Fatal(err)
		}
	}()

	s.messagesMu.Lock()
	if _, exist := s.messages[request.Message]; exist {
		s.messagesMu.Unlock()
		return nil
	}
	s.messages[request.Message] = true
	for _, n := range s.neighbors {
		go func(node string) {
			_ = s.node.Send(node, request)
		}(n)
	}
	s.messagesMu.Unlock()

	return nil
}

func (s *Service) ReadHandler(msg maelstrom.Message) error {
	s.messagesMu.Lock()
	messages := make([]int, 0, len(s.messages))
	for m := range s.messages {
		messages = append(messages, m)
	}
	s.messagesMu.Unlock()
	response := MakeReadResponse(messages)
	return s.node.Reply(msg, response)
}

func (s *Service) TopologyHandler(msg maelstrom.Message) error {
	var request TopologyRequest
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	s.neighbors = request.Topology[s.node.ID()]

	response := MakeTopologyResponse()
	return s.node.Reply(msg, response)
}
