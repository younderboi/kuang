package main

import (
	"fmt"

	"konstantinovitz.com/poopitmaster/internal/agent"
)

func main() {
	// TODO:address and agentID as address as input
	fmt.Println("Starting agent")

	agent := agent.NewTCPAgent("0.0.0.0", "9002")

	agent.Start()
}
