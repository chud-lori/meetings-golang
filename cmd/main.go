package main

import (
	//"context"
	"context"
	"fmt"
	"os"

	//"log"
	//"math/rand"
	"meeting_service/adapters/controllers"
	"meeting_service/adapters/repositories"
	"meeting_service/adapters/web"
	"meeting_service/domain/services"
	"meeting_service/infrastructure"
	"meeting_service/pkg/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

//type Middleware func(http.HandlerFunc) http.HandlerFunc


func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey != "secret-api-key" {
            log := logger.InitiateLogger(r.Context())
			log.Warn("unahothrized bruh")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

    ctx := context.WithValue(context.Background(), "RequestID", uuid.New().String())
    logger.Init(ctx)

    err := godotenv.Load()
    if err != nil {
        logger.Log.Fatal("Failed load keys")
    }

	postgredb := infrastructure.NewPostgreDB()
	defer postgredb.Close()


	userRepository, _ := repositories.NewUserRepositoryPostgre(postgredb, ctx)
	userService := services.NewUserService(userRepository, ctx)
	userController := controllers.NewUserController(userService, ctx)

	userEngageService := services.NewUserEngageService(userRepository, userRepository)
	userEngageController := controllers.NewUserEngageController(userEngageService)

	router := http.NewServeMux()

	web.UserRouter(userController, router)
	web.UserEngageRouter(userEngageController, router)

	var handler http.Handler = router
    handler = logger.LogTrafficMiddleware(handler, ctx)
	handler = APIKeyMiddleware(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: handler,
	}

	logger.Log.Info("App running on port ", os.Getenv("APP_PORT"))

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

