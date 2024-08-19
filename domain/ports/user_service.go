package ports

import (
	"context"
	"meeting_service/adapters/transport"
)

type UserService interface {
    Save(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error)
    Update(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error)
    Delete(ctx context.Context, id string) error
    FindById(ctx context.Context, id string) (*transport.UserResponse, error)
    FindAll(ctx context.Context) ([]*transport.UserResponse, error)
}

type UserEngageService interface {
    CheckInteractionStatus(ctx context.Context, request *transport.UserEngageRequest) (*transport.UserEngageResponse, error)
    RequestAndNotify(ctx context.Context, request *transport.UserEngageRequest) error
    ReceiveAndNotify(ctx context.Context, request *transport.UserEngageRequest) error
    DeclineAndNotify(ctx context.Context, request *transport.UserEngageRequest) error
}

