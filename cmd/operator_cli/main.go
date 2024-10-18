package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	operatorcli "konstantinovitz.com/kuang/internal/operator_cli"
)

func main() {
	bindAddress := flag.String("address", "0.0.0.0:9000", "Local bind port for TCP agents")

	flag.Parse()

	server := operatorcli.NewServer(*bindAddress)

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
