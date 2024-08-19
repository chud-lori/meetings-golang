package entities

import "time"


type Interaction struct {
    Id          string      `json:"id"`
    From_id     string      `json:"from_id"`
    To_id       string      `json:"to_id"`
    Status      int         `json:"status"`
    Created_at  time.Time   `json:"created_at"`
}

