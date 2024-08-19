package ports

import (
	"net/http"
)

type UserController interface {
    Create(writer http.ResponseWriter, request *http.Request)
    Update(writer http.ResponseWriter, request *http.Request)
    Delete(writer http.ResponseWriter, request *http.Request)
    FindById(writer http.ResponseWriter, request *http.Request)
    FindAll(writer http.ResponseWriter, request *http.Request)
}

type UserEngageController interface {
    CheckInteraction(writer http.ResponseWriter, request *http.Request)
    Request(writer http.ResponseWriter, request *http.Request)
    Receive(writer http.ResponseWriter, request *http.Request)
    Decline(writer http.ResponseWriter, request *http.Request)
}

