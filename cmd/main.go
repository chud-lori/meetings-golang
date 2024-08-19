package main

import (
	"log"
	"meeting_service/adapters/controllers"
	"meeting_service/adapters/repositories"
	"meeting_service/adapters/web"
	"meeting_service/domain/services"
	"meeting_service/infrastructure"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey != "secret-api-key" {
			log.Println("unahothrized bruh")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type loggingTraffic struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingTraffic(w http.ResponseWriter) *loggingTraffic {
	return &loggingTraffic{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingTraffic) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func TrafficMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := NewLoggingTraffic(w)

		// call the next handler

		duration := time.Since(start)

		log.Printf(
			"Method %s | Path %s | Duration: %v | Status: %d",
			r.Method,
			r.URL.Path,
			duration,
			lrw.statusCode,
		)

		next.ServeHTTP(w, r)
	})
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}

	return f
}

func main() {
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

	var handler http.Handler = router
	handler = APIKeyMiddleware(handler)
	handler = TrafficMiddleware(handler)

	//middlewareChain := []Middleware{
	//    APIKeyMiddleware,
	//    TrafficMiddleware,
	//}

	//wrappedRouter := Chain(router.ServeHTTP, middlewareChain...)

	server := http.Server{
		Addr:    "localhost:1234",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
