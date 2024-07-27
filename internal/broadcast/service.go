package broadcast

import (
	"encoding/json"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
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

func MakeBroadcastRequest(message int) BroadcastRequest {
	return BroadcastRequest{Broadcast, message}
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

type GossipMessage struct {
	Dst     string
	Message int
}

type Service struct {
	node      *maelstrom.Node
	neighbors []string

	messagesMu sync.RWMutex
	messages   map[int]bool

	Broadcasts chan GossipMessage
}

func MakeService(node *maelstrom.Node) Service {
	return Service{
		node, make([]string, 0),
		sync.RWMutex{},
		make(map[int]bool),
		make(chan GossipMessage, 100),
	}
}

func (s *Service) BroadcastHandler(msg maelstrom.Message) error {
	var request BroadcastRequest
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	s.messagesMu.Lock()
	s.messages[request.Message] = true
	s.messagesMu.Unlock()

	s.Gossip(msg.Src, request.Message)

	return s.node.Reply(msg, MakeBroadcastResponse())
}

func (s *Service) GossipHandler(msg maelstrom.Message) error {
	var request GossipRequest
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	s.messagesMu.Lock()
	if _, exist := s.messages[request.Message]; exist {
		s.messagesMu.Unlock()
		return s.node.Reply(msg, MakeGossipResponse())
	}
	s.messages[request.Message] = true
	s.messagesMu.Unlock()

	s.Gossip(msg.Src, request.Message)

	return s.node.Reply(msg, MakeGossipResponse())
}

func (s *Service) ReadHandler(msg maelstrom.Message) error {
	messages := s.GetMessages()
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

func (s *Service) Gossip(src string, msg int) {
	for _, dst := range s.neighbors {
		if dst == src || dst == s.node.ID() {
			continue
		}

		s.Broadcasts <- GossipMessage{dst, msg}
	}
}

func (s *Service) GetMessages() []int {
	s.messagesMu.RLock()
	messages := make([]int, 0, len(s.messages))
	for m := range s.messages {
		messages = append(messages, m)
	}
	s.messagesMu.RUnlock()

	return messages
}

type GossipRequest struct {
	Type    string
	Message int
}

func MakeGossipRequest(message int) GossipRequest {
	return GossipRequest{"gossip", message}
}

type GossipResponse struct {
	Type string
}

func MakeGossipResponse() GossipResponse {
	return GossipResponse{"gossip_ok"}
}
