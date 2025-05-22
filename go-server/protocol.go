package main


type Message struct {
    Type      string      `json:"type"`
    FromID    string      `json:"from_id"`
    ToID      string      `json:"to_id"`
    Lamport   int         `json:"lamport"`
    Physical  float64     `json:"physical"`
    Payload   interface{} `json:"payload"`
}

type PostPayload struct {
    PostID    string   `json:"post_id"`
    Text      string   `json:"text"`
    Followers []string `json:"followers"`
}

type FollowPayload struct {
    FollowerID string `json:"follower_id"`
    FollowedID string `json:"followed_id"`
}

type PrivateMsgPayload struct {
    MsgID string `json:"msg_id"`
    Text  string `json:"text"`
}

type SyncReplyPayload struct {
    Offset float64 `json:"offset"`
}
