package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"meeting_service/adapters/controllers"
	"meeting_service/adapters/repositories"
	"meeting_service/adapters/utils"
	"meeting_service/adapters/web"
	"meeting_service/domain/services"
	"meeting_service/infrastructure"
	"meeting_service/pkg/logger"
	"net/http"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

//type Middleware func(http.HandlerFunc) http.HandlerFunc

//func APIKeyMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		apiKey := r.Header.Get("x-api-key")
//		if apiKey != "secret-api-key" {
//            //logger, _ := r.Context().Value("logger").(*logrus.Entry)
//            //logger.Warn("Unauth bnruhhh")
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Failed load keys")
        return
    }

    // Handle SIGINT
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    _, otelShutdown, err := infrastructure.SetupOTelSDK(ctx)
    if err != nil {
        log.Println("Failed otel")
        return
    }

    //tracer := tp.Tracer("my-app")

    // handle shutdown properly avoid leaking
    defer func() {
        err = errors.Join(err, otelShutdown(context.Background()))
    }()

	postgredb := infrastructure.NewPostgreDB()
	defer postgredb.Close()

	userRepository, _ := repositories.NewUserRepositoryPostgre(postgredb)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	userEngageService := services.NewUserEngageService(userRepository, userRepository)
	userEngageController := controllers.NewUserEngageController(userEngageService)

	router := http.NewServeMux()

	web.UserRouter(userController, router)
	web.UserEngageRouter(userEngageController, router)

    httpSpanName := func(operation string, r *http.Request) string {
        return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
    }

    var handler = otelhttp.NewHandler(
        router,
        "/",
        otelhttp.WithSpanNameFormatter(httpSpanName))

    handler = logger.LogTrafficMiddleware(handler)
    handler = infrastructure.TraceMiddleware(handler)
	handler = utils.APIKeyMiddleware(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
        BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler: handler,
	}

    srvErr := make(chan error, 1)
    go func() {
        log.Println("App running on port ", os.Getenv("APP_PORT"))
        srvErr <- server.ListenAndServe()
    }()

    // Wait interuption
    select {
    case err = <-srvErr:
        // error start server
        return
    case <-ctx.Done():
        // wait first ctrl c
        // stop receive signal notif asap
        stop()
    }

    err = server.Shutdown(context.Background())
    return

}

