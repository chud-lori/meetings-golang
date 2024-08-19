package controllers

import (
	"encoding/json"
	"fmt"
	"meeting_service/adapters/transport"
	"meeting_service/domain/ports"
	"net/http"
)

type UserController struct {
    ports.UserService
}

func NewUserController(service ports.UserService) *UserController {
    return &UserController{UserService: service}
}

func GetPayload(request *http.Request, result interface{}) {
    decoder := json.NewDecoder(request.Body)
    err := decoder.Decode(result)

    if err != nil {
        panic(err)
    }
}

func WriteResponse(writer http.ResponseWriter, response interface{}, httpCode int64) {

    writer.Header().Add("Content-Type", "application/json")
    writer.WriteHeader(int(httpCode))
    encoder := json.NewEncoder(writer)
    err := encoder.Encode(response)

    if err != nil {
        panic(err)
    }
}

type WebResponse struct {
    Message string   `json:"message"`
    Status  int      `json:"status"`
    Data interface{} `json:"data"`
}

func (controller *UserController) Create(writer http.ResponseWriter, request *http.Request) {
    userRequest := transport.UserRequest{}
    GetPayload(request, &userRequest)

    userResponse, err := controller.UserService.Save(request.Context(), &userRequest)

    if err != nil {
        fmt.Println("Error create controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success save user",
        Status: 1,
        Data: userResponse,
    }

    WriteResponse(writer, &response, http.StatusCreated)
}

func (controller *UserController) Update(writer http.ResponseWriter, request *http.Request) {
    userRequest := transport.UserRequest{}
    GetPayload(request, &userRequest)

    userResponse, err := controller.UserService.Update(request.Context(), &userRequest)

    if err != nil {
        fmt.Println("Error update controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success update user",
        Status: 1,
        Data: userResponse,
    }

    WriteResponse(writer, &response, http.StatusOK)
}

func (controller *UserController) Delete(writer http.ResponseWriter, request *http.Request) {
    userId := request.PathValue("userId")

    err := controller.UserService.Delete(request.Context(), userId)

    if err != nil {
        fmt.Println("Error delete controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success delete user",
        Status: 1,
        Data: "sucess",
    }

    WriteResponse(writer, &response, http.StatusOK)
}

func (controller *UserController) FindById(writer http.ResponseWriter, request *http.Request) {
    userId := request.PathValue("userId")

    user, err := controller.UserService.FindById(request.Context(), userId)

    if err != nil {
        fmt.Println("Error Find by id controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success get user by id",
        Status: 1,
        Data: &user,
    }

    WriteResponse(writer, &response, http.StatusOK)
}

func (controller *UserController) FindAll(writer http.ResponseWriter, request *http.Request) {
    users, err := controller.UserService.FindAll(request.Context())

    if err != nil {
        fmt.Println("Error Find by id controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success get all users",
        Status: 1,
        Data: users,
    }

    WriteResponse(writer, &response, http.StatusOK)
}
