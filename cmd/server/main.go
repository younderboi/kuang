package main

import "konstantinovitz.com/poopitmaster/internal/server"

func main() {
	handler := server.NewServer()

	err := handler.Start()
	if err != nil {
		println("Failed to start handler:", err)
	}

	defer handler.Stop()
}
