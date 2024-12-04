package services

import (
	"context"
	"fmt"
	"math/rand"
	"meeting_service/adapters/transport"
	"meeting_service/domain/entities"
	"meeting_service/domain/ports"
	"meeting_service/grpc_service"
	"meeting_service/pkg/logger"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	ports.UserRepository
    logger *logrus.Entry
}

// provider or constructor
func NewUserService(userRepository ports.UserRepository, ctx context.Context) *UserServiceImpl {
	return &UserServiceImpl{
		UserRepository: userRepository,
        logger: logger.InitiateLogger(ctx),
	}
}

func generatePasscode() string {
    // get current ms
    curMs := time.Now().Nanosecond() / 1000

    // convert ms to str and get first 4 char
    msStr := strconv.Itoa(curMs)[:4]

    // generate random char between A and Z
    var alphb []int
    for i := 0; i < 4; i++ {
        alphb = append(alphb, rand.Intn(26)+65)
    }

    // Convert ascii values to character and join them
    var alphChar []string
    for _, a := range alphb {
        alphChar = append(alphChar, string(rune(a)))
    }
    alphStr := strings.Join(alphChar, "")

    // combine alphabet string and ms string
    return alphStr + msStr
}

func (service *UserServiceImpl) Save(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
	user := entities.User{
		Id:         "",
		Email:      request.Email,
        Passcode: generatePasscode(),
		Created_at: time.Now(),
	}
	user_result, error := service.UserRepository.Save(ctx, &user)

	if error != nil {

		panic(error)
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
		Created_at: user_result.Created_at,
	}

	return user_response, nil
}

func (service *UserServiceImpl) Update(ctx context.Context, request *transport.UserRequest) (*transport.UserResponse, error) {
	user := entities.User{
		Id:         "",
		Email:      request.Email,
		Created_at: time.Now(),
	}

	user_result, error := service.UserRepository.Update(ctx, &user)
	if error != nil {
		panic(error)
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
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

func (s *UserServiceImpl) FindById(ctx context.Context, id string) (*transport.UserResponse, error) {
	user := entities.User{}

	user_result, err := s.UserRepository.FindById(ctx, id)

	user.Id = user_result.Id
	user.Email = user_result.Email
	user.Created_at = user_result.Created_at

	if err != nil {
		return nil, err
	}

	user_response := &transport.UserResponse{
		Id:         user_result.Id,
		Email:      user_result.Email,
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
			Id:         user.Id,
			Email:      user.Email,
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
func NewUserEngageService(userEngageRepository ports.UserEngageRepository, userRepository ports.UserRepository) *UserEngageServiceImpl {
	return &UserEngageServiceImpl{
		UserEngageRepository: userEngageRepository,
		UserRepository:       userRepository,
	}
}

func (service *UserEngageServiceImpl) CheckInteractionStatus(ctx context.Context, request *transport.UserEngageRequest) (*transport.UserEngageResponse, error) {
	status, err := service.UserEngageRepository.CheckInteractionStatus(ctx, request.FromUser, request.ToUser)

	if len(status) == 0 {
		return nil, fmt.Errorf("Status not found")
	}
	userResponse := &transport.UserEngageResponse{Status: status[0]}

	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (service *UserEngageServiceImpl) RequestAndNotify(ctx context.Context, request *transport.UserEngageRequest) error {
	statuses, err := service.UserEngageRepository.CheckInteractionStatus(ctx, request.FromUser, request.ToUser)

	if len(statuses) > 0 && statuses[0] == 1 {
		return fmt.Errorf("Failed because previous request not proceed")
	}
	err = service.UserEngageRepository.Request(ctx, request.FromUser, request.ToUser)

	if err != nil {
		return err
	}

	user, err := service.UserRepository.FindById(ctx, request.ToUser)
    fromUser, err := service.UserRepository.FindById(ctx, request.FromUser)

	if err != nil {
		return err
	}

	go grpc_service.SendGrpcMail(user.Email, fmt.Sprintf("Meeting request from %s", fromUser.Email))

	return nil
}

func (service *UserEngageServiceImpl) ReceiveAndNotify(ctx context.Context, request *transport.UserEngageRequest) error {
	statuses, err := service.UserEngageRepository.CheckInteractionStatus(ctx, request.FromUser, request.ToUser)
	if len(statuses) > 0 && statuses[0] != 1 {
		return fmt.Errorf("Failed no request detected")
	}
	err = service.UserEngageRepository.Receive(ctx, request.FromUser, request.ToUser)
	err = service.UserEngageRepository.StoreMeet(ctx, request.FromUser, request.ToUser)

	if err != nil {
		return err
	}

	user1, err := service.UserRepository.FindById(ctx, request.ToUser)
	user2, err := service.UserRepository.FindById(ctx, request.FromUser)

	if err != nil {
		return err
	}

	go grpc_service.SendGrpcMail(fmt.Sprintf("%s,%s", user1.Email, user2.Email), fmt.Sprintf("Meeting set up between %s and %s", user1.Email, user2.Email))

	return nil
}

func (service *UserEngageServiceImpl) DeclineAndNotify(ctx context.Context, request *transport.UserEngageRequest) error {
	statuses, err := service.UserEngageRepository.CheckInteractionStatus(ctx, request.FromUser, request.ToUser)
	// log.Printf("Status %#v\n", statuses)
	// return nil
	if len(statuses) > 0 && statuses[0] != 1 {
		return fmt.Errorf("Failed no request detected")
	}

	err = service.UserEngageRepository.Decline(ctx, request.FromUser, request.ToUser)

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
