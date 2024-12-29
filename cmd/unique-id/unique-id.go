package main

import (
	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

// A message that is received by maelstrom
type RpcMessage struct {
	// A string identifying the node this message came from
	Src string `json:"type,omitempty"`

	// A string identifying the node this message is to
	Dest string `json:"dest,omitempty"`

	//An object: the payload of the message
	Body MessageBody `json:"body,omitempty"`
}

type MessageBody struct {
	// A string identifying the type of message this is
	Type string `json:"type,omitempty"`

	// Optional. Message identifier that is unique to the source node.
	MsgID int `json:"msg_id,omitempty"`

	// Optional. For request/response, the msg_id of the request.
	InReplyTo int `json:"in_reply_to,omitempty"`

	// Error code, if an error occurred.
	Code int `json:"code,omitempty"`

	// Error message, if an error occurred.
	Text string `json:"text,omitempty"`
}

type RpcResponse struct {
	// A string identifying the type of message this is
	Type string `json:"type,omitempty"`

	// A globally unique id
	Id uuid.UUID `json:"id"`

	// // Optional. Message identifier that is unique to the source node.
	// MsgID int `json:"msg_id,omitempty"`
	//
	// // Optional. For request/response, the msg_id of the request.
	// InReplyTo int `json:"in_reply_to,omitempty"`
}

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		uid, err := uuid.NewV7()
		if err != nil {
			log.Fatalln(err)
		}
		body := RpcResponse{
			Type: "generate_ok",
			Id:   uid,
		}
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatalln(err)
	}
}
