package controllers

import (
	"fmt"
	"meeting_service/adapters/transport"
	"meeting_service/domain/ports"
	"net/http"
)

type UserEngageController struct {
    ports.UserEngageService
}

func NewUserEngageController(service ports.UserEngageService) *UserEngageController {
    return &UserEngageController{UserEngageService: service}
}


func (controller *UserEngageController) CheckInteraction(writer http.ResponseWriter, request *http.Request) {
    userRequest := transport.UserEngageRequest{}
    GetPayload(request, &userRequest)

    userResponse, err := controller.UserEngageService.CheckInteractionStatus(request.Context(), &userRequest)

    if err != nil {
        fmt.Println("Error check interaction status controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success get interaction status",
        Status: 1,
        Data: userResponse,
    }

    WriteResponse(writer, &response, http.StatusOK)
}

func (controller *UserEngageController) Request(writer http.ResponseWriter, request *http.Request) {
    userRequest := transport.UserEngageRequest{}
    GetPayload(request, &userRequest)

    err := controller.UserEngageService.RequestAndNotify(request.Context(), &userRequest)

    if err != nil {
        fmt.Println("Error request controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success request meeting",
        Status: 1,
        Data: nil,
    }

    WriteResponse(writer, &response, http.StatusOK)
}

func (controller *UserEngageController) Receive(writer http.ResponseWriter, request *http.Request) {
    userRequest := transport.UserEngageRequest{}
    GetPayload(request, &userRequest)

    err := controller.UserEngageService.ReceiveAndNotify(request.Context(), &userRequest)

    if err != nil {
        fmt.Println("Error receive controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success receive meeting request",
        Status: 1,
        Data: nil,
    }

    WriteResponse(writer, &response, http.StatusOK)
}

func (controller *UserEngageController) Decline(writer http.ResponseWriter, request *http.Request) {
    userRequest := transport.UserEngageRequest{}
    GetPayload(request, &userRequest)

    err := controller.UserEngageService.DeclineAndNotify(request.Context(), &userRequest)

    if err != nil {
        fmt.Println("Error decline controller")
        panic(err)
    }

    response := WebResponse{
        Message: "success decline meeting request",
        Status: 1,
        Data: nil,
    }

    WriteResponse(writer, &response, http.StatusOK)
}
