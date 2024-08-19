package entities

import "time"


type Meet struct {
    Id          string      `json:"id"`
    User_1      string      `json:"user_1"`
    User_2      string      `json:"user_2"`
    Created_at  time.Time   `json:"created_at"`
}

