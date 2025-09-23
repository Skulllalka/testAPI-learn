package main

import (
	"API/http"
	"API/todo"
	"fmt"
)

func main() {
	todoList := todo.NewList()
	httpHandlers := http.NewHTTPHandlers(todoList)
	httpServer := http.NewHTTPServer(httpHandlers)
	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start server", err)
	}
}
