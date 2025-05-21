package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
)

var driftOffset float64 = 0.0

func main() {
	serverID := os.Getenv("SERVER_ID")
	if serverID == "" {
		log.Fatal("SERVER_ID env var required, e.g. 'server0'")
	}
	peers := os.Getenv("PEER_ADDRS")
	peerList := []string{}
	if peers != "" {
		peerList = append(peerList, peers)
	}

	store := NewStorage()

	// ZeroMQ sockets
	rep, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatalf("Failed to create REP socket: %v", err)
	}
	rep.Bind("tcp://*:5555")

	pub, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatalf("Failed to create PUB socket: %v", err)
	}
	pub.Bind("tcp://*:5556")

	push, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		log.Fatalf("Failed to create PUSH socket: %v", err)
	}
	for _, addr := range peerList {
		push.Connect(addr)
	}

	pull, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		log.Fatalf("Failed to create PULL socket: %v", err)
	}
	pull.Bind("tcp://*:5560")

	// Calcule um índice de servidor (0,1,2) a partir do SERVER_ID:
	idx := serverIndex(serverID) // se SERVER_ID é "server1", retorna 1
	syncPort := 7000 + idx

	syncRep, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		log.Fatalf("sync socket: %v", err)
	}
	syncRep.Bind(fmt.Sprintf("tcp://*:%d", syncPort))

	go func() {
		for {
			data, err := pull.RecvBytes(0)
			if err != nil {
				log.Printf("Error receiving from pull: %v", err)
				continue
			}
			var msg Message
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Printf("Invalid JSON from peer: %v", err)
				continue
			}
			store.LogMessage(msg)
			if msg.Type == "POST" {
				bytes, _ := json.Marshal(msg.Payload)
				var p PostPayload
				json.Unmarshal(bytes, &p)
				notifyFollowers(pub, p.Followers, p)
			}
		}
	}()

	go func() {
		for {
			raw, err := syncRep.RecvBytes(0)
			if err != nil {
				log.Printf("Error receiving sync request: %v", err)
				continue
			}
			var msg Message
			if err := json.Unmarshal(raw, &msg); err != nil {
				log.Printf("Invalid sync JSON: %v", err)
				errResp := map[string]interface{}{"status": 400, "error": "invalid sync JSON"}
				b, _ := json.Marshal(errResp)
				syncRep.SendBytes(b, 0)
				continue
			}

			switch msg.Type {
			case "SYNC_REQUEST":
				resp := Message{
					Type:     "SYNC_REPLY",
					FromID:   serverID,
					ToID:     msg.FromID,
					Lamport:  0,
					Physical: float64(time.Now().UnixNano())/1e9 + driftOffset,
					Payload:  map[string]float64{"offset": float64(time.Now().UnixNano())/1e9 + driftOffset},
				}
				b, _ := json.Marshal(resp)
				syncRep.SendBytes(b, 0)

			case "SYNC_ADJUST":
				payloadBytes, _ := json.Marshal(msg.Payload)
				var adj struct{ Adjust float64 }
				json.Unmarshal(payloadBytes, &adj)
				driftOffset += adj.Adjust
				ack := map[string]interface{}{"status": 200}
				b, _ := json.Marshal(ack)
				syncRep.SendBytes(b, 0)

			default:
				ack := map[string]interface{}{"status": 400, "error": "unknown sync type"}
				b, _ := json.Marshal(ack)
				syncRep.SendBytes(b, 0)
			}
		}
	}()

	for {
		raw, err := rep.RecvBytes(0)
		if err != nil {
			log.Printf("Error receiving request: %v", err)
			continue
		}
		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			errResp := map[string]interface{}{"status": 400, "error": "invalid JSON"}
			errData, _ := json.Marshal(errResp)
			rep.Send(string(errData), 0)
			continue
		}

		msg.Physical = float64(time.Now().UnixNano())/1e9 + driftOffset

		store.LogMessage(msg)
		if push != nil {
			out, _ := json.Marshal(msg)
			push.SendBytes(out, 0)
		}

		if msg.Type == "POST" {
			bytes, _ := json.Marshal(msg.Payload)
			var p PostPayload
			json.Unmarshal(bytes, &p)
			notifyFollowers(pub, p.Followers, p)
		}

		okResp := map[string]interface{}{"status": 200, "message": "OK"}
		okData, _ := json.Marshal(okResp)
		rep.Send(string(okData), 0)
	}
}

func notifyFollowers(pub *zmq.Socket, followers []string, p PostPayload) {
	for _, f := range followers {
		notif := Message{
			Type:     "NOTIFY",
			FromID:   p.PostID,
			ToID:     f,
			Lamport:  0,
			Physical: float64(time.Now().UnixNano())/1e9 + driftOffset,
			Payload: map[string]string{
				"post_id": p.PostID,
				"text":    p.Text,
			},
		}
		data, _ := json.Marshal(notif)
		pub.Send(f+" "+string(data), 0)
	}
}

func serverIndex(id string) int {
	numPart := id[len(id)-1:]
	i, _ := strconv.Atoi(numPart)
	return i
}
