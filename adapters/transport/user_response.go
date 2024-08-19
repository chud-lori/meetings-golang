package transport

import "time"

type UserResponse struct {
    Id string            `json:"id"`
    Email string         `json:"email"`
    Created_at time.Time `json:"created_at"`
}

type UserEngageResponse struct {
    Status uint8 `json:"status"`
}
