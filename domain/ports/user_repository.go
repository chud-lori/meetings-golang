package ports

import (
	"context"
	"meeting_service/domain/entities"
)

type UserRepository interface {
    Save(ctx context.Context, user *entities.User) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) (*entities.User, error)
    Delete(ctx context.Context, id string) error
    FindById(ctx context.Context, id string) (*entities.User, error)
    FindAll(ctx context.Context) ([]*entities.User, error)
}

type UserEngageRepository interface {
    CheckInteractionStatus(ctx context.Context, user1 string, user2 string) ([]uint8, error)
    Request(ctx context.Context, fromUser string, toUser string) error
    Receive(ctx context.Context, fromUser string, toUser string) error
    Decline(ctx context.Context, fromUser string, toUser string) error
}

