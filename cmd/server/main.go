package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"konstantinovitz.com/kuang/internal/server"
)

func main() {
	bindAddress := flag.String("address", "0.0.0.0:9000", "Local bind port for TCP agents")

	flag.Parse()

	server := server.NewServer(*bindAddress)

	go func() {
		err := server.Start()
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Catch OS signals to gracefully shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	server.Stop()
	fmt.Println("Server stopped")
}
