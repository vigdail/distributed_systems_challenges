package main

import (
	"log"

	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const echoOk = "echo_ok"

type echoRequest struct {
	Echo string `json:"echo"`
}

type echoResponse struct {
	Type string `json:"type"`
	Echo string `json:"echo"`
}

func makeEchoResponse(echo string) echoResponse {
	return echoResponse{
		echoOk,
		echo,
	}
}

type service struct {
	node *maelstrom.Node
}

func (s *service) echo(msg maelstrom.Message) error {
	var request echoRequest
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	response := makeEchoResponse(request.Echo)
	return s.node.Reply(msg, response)
}

func main() {
	n := maelstrom.NewNode()
	s := service{n}
	n.Handle("echo", s.echo)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
