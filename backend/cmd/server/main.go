package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jjeejj/todolist/backend/internal/repository"
	"github.com/jjeejj/todolist/backend/internal/service"
	"github.com/jjeejj/todolist/backend/proto/todolist/v1/todolistv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version, Connect-Timeout-Ms")
		w.Header().Set("Access-Control-Expose-Headers", "Connect-Protocol-Version")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize repository and service
	repo := repository.NewTaskRepository()
	todoService := service.NewTodoService(repo)

	// Create HTTP mux
	mux := http.NewServeMux()

	// Register the TodoService handler
	path, handler := todolistv1connect.NewTodoServiceHandler(todoService)
	mux.Handle(path, handler)

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Wrap with CORS middleware
	corsHandler := corsMiddleware(mux)

	// Use h2c for HTTP/2 without TLS (required for Connect-RPC)
	h2cHandler := h2c.NewHandler(corsHandler, &http2.Server{})

	port := ":8080"
	fmt.Printf("Todo service starting on port %s\n", port)
	fmt.Println("Health check available at: http://localhost:8080/health")
	fmt.Println("TodoService available at: http://localhost:8080/todolist.v1.TodoService/")

	if err := http.ListenAndServe(port, h2cHandler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}