// go-server/protocol.go
package main

// Message defines the common JSON envelope for all inter-process communications
type Message struct {
    Type      string      `json:"type"`
    FromID    string      `json:"from_id"`
    ToID      string      `json:"to_id"`
    Lamport   int         `json:"lamport"`
    Physical  float64     `json:"physical"`
    Payload   interface{} `json:"payload"`
}

// PostPayload carries data for a new post
type PostPayload struct {
    PostID    string   `json:"post_id"`
    Text      string   `json:"text"`
    Followers []string `json:"followers"`
}

// FollowPayload carries data for a follow action
type FollowPayload struct {
    FollowerID string `json:"follower_id"`
    FollowedID string `json:"followed_id"`
}

// PrivateMsgPayload carries data for a private message
type PrivateMsgPayload struct {
    MsgID string `json:"msg_id"`
    Text  string `json:"text"`
}

// SyncReplyPayload carries offset for Berkeley
type SyncReplyPayload struct {
    Offset float64 `json:"offset"`
}
