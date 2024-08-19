package services

import (
	"context"
	"fmt"
	"meeting_service/adapters/transport"
	"meeting_service/domain/entities"
	"meeting_service/domain/ports"
	"meeting_service/grpc_service"
	"time"
)

type UserServiceImpl struct {
    ports.UserRepository
}

// provider or constructor
func NewUserService(userRepository ports.UserRepository) *UserServiceImpl {
    return &UserServiceImpl{
        UserRepository: userRepository,
    }
}

func (service *UserServiceImpl) Save(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
    user := entities.User{
        Id: "",
        Email: request.Email,
        Created_at: time.Now(),
    }
    user_result, error := service.UserRepository.Save(ctx, &user)

    if error != nil {

       panic(error)
    }

    user_response := &transport.UserResponse{
        Id: user_result.Id,
        Email: user_result.Email,
        Created_at: user_result.Created_at,
    }

    return user_response, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
    user := entities.User{
        Id: "",
        Email: request.Email,
        Created_at: time.Now(),
    }

    user_result, error := service.UserRepository.Update(ctx, &user)
    if error != nil {
        panic(error)
    }

    user_response := &transport.UserResponse{
        Id: user_result.Id,
        Email: user_result.Email,
        Created_at: user_result.Created_at,
    }

    return user_response, nil
}

func (service *UserServiceImpl) Delete(ctx context.Context, id string) error {

    err := service.UserRepository.Delete(ctx, id)

    if err != nil {
        return err
    }

    return nil
}


func (service *UserServiceImpl) FindById(ctx context.Context, id string) (*transport.UserResponse, error) {
    user := entities.User{}

    user_result, err := service.UserRepository.FindById(ctx, id)

    user.Id = user_result.Id
    user.Email = user_result.Email
    user.Created_at = user_result.Created_at

    if err != nil {
        return nil, err
    }

    user_response := &transport.UserResponse{
        Id: user_result.Id,
        Email: user_result.Email,
        Created_at: user_result.Created_at,
    }

    return user_response, nil
}

func (service *UserServiceImpl) FindAll(ctx context.Context) ([]*transport.UserResponse, error) {

    users_result, err := service.UserRepository.FindAll(ctx)

    if err != nil {
        return nil, err
    }

    //var users_response []transport.UserResponse
    users_response := make([]*transport.UserResponse, len(users_result))

    for i, user := range users_result {
        users_response[i] = &transport.UserResponse{
            Id: user.Id,
            Email: user.Email,
            Created_at: user.Created_at,
        }
    }

    return users_response, nil
}

type UserEngageServiceImpl struct {
    ports.UserEngageRepository
    ports.UserRepository
}

// provider or constructor
func NewUserEngageService(userEngageRepository ports.UserEngageRepository, userRepository ports.UserRepository) *UserEngageServiceImpl{
    return &UserEngageServiceImpl{
        UserEngageRepository: userEngageRepository,
        UserRepository: userRepository,
    }
}

func (service *UserEngageServiceImpl) CheckInteractionStatus(ctx context.Context, request *transport.UserEngageRequest) (*transport.UserEngageResponse, error) {
	status, err := service.UserEngageRepository.CheckInteractionStatus(ctx, request.FromUser, request.ToUser)

	if len(status) == 0 {
		return nil, fmt.Error("Status not found")
	}
    userResponse := &transport.UserEngageResponse{Status: status[0]}

    if err != nil {
        return nil, err
    }

    return userResponse, nil
}

func (service *UserEngageServiceImpl) RequestAndNotify(ctx context.Context, request *transport.UserEngageRequest) error {
    err := service.UserEngageRepository.Request(ctx, request.FromUser, request.ToUser)

    if err != nil {
        return err
    }

    user, err := service.UserRepository.FindById(ctx, request.ToUser)

    if err != nil {
        return err
    }

    go grpc_service.SendGrpcMail(user.Email, fmt.Sprintf("Meeting request from %s", user.Email))

    return nil
}

func (service *UserEngageServiceImpl) ReceiveAndNotify(ctx context.Context, request *transport.UserEngageRequest) error {
    err := service.UserEngageRepository.Receive(ctx, request.FromUser, request.ToUser)

    if err != nil {
        return err
    }

    user, err := service.UserRepository.FindById(ctx, request.ToUser)

    if err != nil {
        return err
    }

    go grpc_service.SendGrpcMail(user.Email, fmt.Sprintf("Meeting accept from %s", user.Email))

    return nil
}

func (service *UserEngageServiceImpl) DeclineAndNotify(ctx context.Context, request *transport.UserEngageRequest) error {
    err := service.UserEngageRepository.Decline(ctx, request.FromUser, request.ToUser)

    if err != nil {
        return err
    }

    user, err := service.UserRepository.FindById(ctx, request.ToUser)

    if err != nil {
        return err
    }

    go grpc_service.SendGrpcMail(user.Email, fmt.Sprintf("Meeting request decline from %s", user.Email))

    return nil
}
