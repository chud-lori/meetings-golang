package transport

type UserRequest struct {
    Email string `validate:"required,max=200,min=1" json:"email"`
}

type UserEngageRequest struct {
    FromUser string `json:"from_user"`
    ToUser string   `json:"to_user"`
    Status uint8    `json:"status"`
}

