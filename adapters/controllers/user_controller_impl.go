package controllers

import (
	"encoding/json"
	"fmt"
	"meeting_service/adapters/transport"
	"meeting_service/domain/ports"
	"meeting_service/infrastructure"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

func (c *UserController) FindById(w http.ResponseWriter, r *http.Request) {

    _, span := otel.Tracer("example-tracer").Start(r.Context(), "example-span")
    defer span.End()

    userId := r.PathValue("userId")

    user, err := c.UserService.FindById(r.Context(), userId)

    if err != nil {
        logger, _ := r.Context().Value("logger").(*logrus.Entry)
        logger.Info("Error find by id controller: ", err)

        WriteResponse(w, WebResponse{
            Message: "Failed get user id",
            Status: 0,
            Data: nil,
        }, http.StatusNotFound)
        return
    }

    response := WebResponse{
        Message: "success get user by id",
        Status: 1,
        Data: &user,
    }

    WriteResponse(w, &response, http.StatusOK)
}

func (controller *UserController) FindAll(w http.ResponseWriter, r *http.Request) {
    logger, _ := r.Context().Value("logger").(*logrus.Entry)

    ctx, endSpan := infrastructure.TraceFunction(r.Context(), "FindAllHandler",
        attribute.String("functionNAMA", "FindAll"),
    )

    users, err := controller.UserService.FindAll(ctx)

    if err != nil {
        logger.Info("Error Find All user: ", err)
        panic(err)
    }

    response := WebResponse{
        Message: "success get all users",
        Status: 1,
        Data: users,
    }
    endSpan()
    WriteResponse(w, &response, http.StatusOK)
}

