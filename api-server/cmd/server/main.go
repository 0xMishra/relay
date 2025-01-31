package main

import (
	"fmt"
	"net/http"

	"github.com/0xMishra/relay/api-server/internal/handlers"
	"github.com/0xMishra/relay/api-server/internal/middlewares"
	"github.com/0xMishra/relay/api-server/internal/utils"
)

func main() {
	mux := http.NewServeMux()

	// To run ECS task
	mux.HandleFunc("POST /project", handlers.RunEcsTaskHandler)

	// To send realltime build logs to the user
	mux.Handle("/ws/{id}", middlewares.RedisSetup(http.HandlerFunc(handlers.SocketLogsHandler)))

	// To proxy the users request to static files in S3 bucket
	mux.HandleFunc("/", handlers.ReverseProxy)

	fmt.Println("api server running on PORT:3000")
	err := http.ListenAndServe(":3000", middlewares.SetCorsHeaders(mux))
	utils.CheckErr(err, true)
}
