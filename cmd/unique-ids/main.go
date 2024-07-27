package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const generate = "generate"
const generateOk = "generate_ok"

type generateResponse struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

func makeGenerateResponse(id string) generateResponse {
	return generateResponse{
		generateOk,
		id,
	}
}

type service struct {
	node *maelstrom.Node
}

func (s *service) generate(msg maelstrom.Message) error {
	var request maelstrom.MessageBody
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	id := fmt.Sprintf("%s-%d", s.node.ID(), request.MsgID)
	response := makeGenerateResponse(id)
	return s.node.Reply(msg, response)
}

func main() {
	n := maelstrom.NewNode()
	s := service{n}
	n.Handle(generate, s.generate)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
