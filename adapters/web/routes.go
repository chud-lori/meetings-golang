package web

import (
	"meeting_service/domain/ports"
	"net/http"
)

func UserRouter(controller ports.UserController, serve *http.ServeMux){
    serve.HandleFunc("POST /api/user", controller.Create)
    serve.HandleFunc("PUT /api/user", controller.Update)
    serve.HandleFunc("DELETE /api/user/{userId}", controller.Delete)
    serve.HandleFunc("GET /api/user/{userId}", controller.FindById)
    serve.HandleFunc("GET /api/user", controller.FindAll)

    //router.PanicHandler
}

func UserEngageRouter(controller ports.UserEngageController, serve *http.ServeMux){
    serve.HandleFunc("POST /api/engage/checkinteraction", controller.CheckInteraction)
    serve.HandleFunc("POST /api/engage/request", controller.Request)
    serve.HandleFunc("POST /api/engage/receive", controller.Receive)
    serve.HandleFunc("POST /api/engage/decline", controller.Decline)
}

