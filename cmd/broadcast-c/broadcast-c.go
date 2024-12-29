package main

import (
	"encoding/json"
	"slices"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastRequest struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
	Msg_id  int64  `json:"msg_id"`
	// This is an array of nodes that have received this message.
	// Every time we get a BroadcastRequest,
	// we must check if the len of this array is equal to the number of nodes.
	// If it is, we don't broadcast the message any further.
	// If it isn't, we check if our node's name is already in the array.
	// if it isn't, we append it broadcast it further.
	AckNodes []string `json:"ack_nodes"`
}

type BroadcastResponse struct {
	Type string `json:"type"`
}

type TopologyRequest struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type TopologyResponse struct {
	Type string `json:"type"`
}

type ReadRequest struct {
	Type string `json:"type"`
}

type ReadResponse struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages"`
}

//	 This is the set of all received messages.
//	 Anytime a new message is received we:
//	 1. Check if it has not been received yet,
//		if it has: don't add it to the Set.
//		if it hasn't: add it to the Set.
type MessageSet struct {
	set map[int64]struct{}
	sync.Mutex
}

func NewMessageSet() MessageSet {
	return MessageSet{
		set: make(map[int64]struct{}),
	}
}

func (s *MessageSet) add(newItem int64) {
	s.Lock()
	defer s.Unlock()

	_, found := s.set[newItem]
	if !found {
		s.set[newItem] = struct{}{}
	}
	return
}

func (s *MessageSet) readAll() []int64 {
	s.Lock()
	defer s.Unlock()

	var returnVec []int64
	for key := range s.set {
		returnVec = append(returnVec, key)
	}
	return returnVec
}

func AppendIdToAckNodes(ackNodes *[]string, id string) {
	if ackNodes == nil {
		return
	}

	// declare a flag variable
	found := false

	// search for an occurence of the node name in the node slices
	for i := range *ackNodes {
		if (*ackNodes)[i] == id {
			found = true
			break
		}
	}

	// if not found, append current node name into the slice
	if !found {
		*ackNodes = append(*ackNodes, id)
	}
	return
}

func main() {
	node := maelstrom.NewNode()
	nodeMessageSet := NewMessageSet()

	// this will be initialized by the 'topology' rpc
	var nodeTopology []string

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body TopologyRequest

		json.Unmarshal(msg.Body, &body)

		// we get our neighbour nodes and copy them over as our node's topology
		neighbourNodes, _ := body.Topology[node.ID()]
		nodeTopology = neighbourNodes

		response := TopologyResponse{
			Type: "topology_ok",
		}
		return node.Reply(msg, response)
	})

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastRequest
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// add the message to this node's internal message Set.
		nodeMessageSet.add(body.Message)

		AppendIdToAckNodes(&body.AckNodes, node.ID())

		// I gotta improve this
		// if there are still some nodes to get this message, broadcast it further
		if len(body.AckNodes) <= len(node.NodeIDs()) {
			for _, val := range nodeTopology {
				if slices.Index(body.AckNodes, val) == -1 {
					go node.Send(val, body)
				}
			}
		}

		response := BroadcastResponse{
			Type: "broadcast_ok",
		}
		return node.Reply(msg, response)
	})

	node.Handle("read", func(msg maelstrom.Message) error {

		response := ReadResponse{
			Type:     "read_ok",
			Messages: nodeMessageSet.readAll(),
		}
		return node.Reply(msg, response)
	})

	node.Run()
}
