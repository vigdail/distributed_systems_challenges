package broadcast

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const Broadcast = "broadcast"
const BroadcastOk = "broadcast_ok"
const Read = "read"
const ReadOk = "read_ok"
const Topology = "topology"
const TopologyOk = "topology_ok"

type BroadcastRequest struct {
	Message int `json:"message"`
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

	messages map[int]bool
}

func MakeService(node *maelstrom.Node) Service {
	return Service{
		node, make([]string, 0),
		make(map[int]bool),
	}
}

func (s *Service) BroadcastHandler(msg maelstrom.Message) error {
	var request BroadcastRequest
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	s.messages[request.Message] = true

	return s.node.Reply(msg, MakeBroadcastResponse())
}

func (s *Service) ReadHandler(msg maelstrom.Message) error {
	messages := make([]int, 0, len(s.messages))
	for m := range s.messages {
		messages = append(messages, m)
	}
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
