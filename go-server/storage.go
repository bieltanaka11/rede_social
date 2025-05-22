package main

import (
    "log"
    "sync"
    "time"
)

type Storage struct {
    mu          sync.Mutex
    Posts       []Message             
    Follows     map[string]map[string]bool
    PrivateMsgs map[string][]Message  
}

func NewStorage() *Storage {
    return &Storage{
        Posts:       make([]Message, 0),
        Follows:     make(map[string]map[string]bool),
        PrivateMsgs: make(map[string][]Message),
    }
}

func (s *Storage) LogMessage(msg Message) {
    s.mu.Lock()
    defer s.mu.Unlock()
    switch msg.Type {
    case "POST":
        s.Posts = append(s.Posts, msg)
    case "FOLLOW":
        payload := msg.Payload.(FollowPayload)
        if _, ok := s.Follows[payload.FollowerID]; !ok {
            s.Follows[payload.FollowerID] = make(map[string]bool)
        }
        s.Follows[payload.FollowerID][payload.FollowedID] = true
    case "MSG_PRIVATE":
        s.PrivateMsgs[msg.ToID] = append(s.PrivateMsgs[msg.ToID], msg)
    }
    log.Printf("[%s] [%s] Stored message: %+v\n", time.Now().Format(time.RFC3339), msg.Type, msg)
}
