package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"konstantinovitz.com/kuang/internal/fileserver"
)

func main() {
	bindPort := flag.String("port", "8000", "Local bind port")
	serveDir := flag.String("servedir", "./www/", "File path to serve files from")

	flag.Parse()

	fileserver := fileserver.NewFileServer(*serveDir, *bindPort)
	fileserver.Start()
	go func() {
		err := fileserver.Start()
		if err != nil {
			fmt.Println("Error starting fileserver", err)
		}
	}()

	// Catch OS signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	fileserver.Stop()
	fmt.Println("Server stopped")
}
