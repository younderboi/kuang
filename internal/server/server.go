package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"konstantinovitz.com/poopitmaster/internal/utils"
)

type Server struct {
	address        string
	port           string
	currentAgent   AgentProxy
	agents         map[AgentID]AgentProxy
	commandManager utils.CommandManager
}

func NewServer() *Server {
	return &Server{
		currentAgent: nil,
		agents:       make(map[AgentID]AgentProxy),
	}
}

func (server *Server) Start() error {
	commandManager := utils.NewCommandManager()

	commandManager.RegisterHandler("listen", server.handleListen)
	commandManager.RegisterHandler("attach", server.handleAttach)
	commandManager.RegisterHandler("list", server.handleListAgents)
	commandManager.RegisterHandler("kill", server.handleKill)
	commandManager.RegisterHandler("clear", server.handleClear)
	commandManager.RegisterHandler("detach", server.handleDetach)
	commandManager.RegisterHandler("ex", server.handleAgentExec)
	// TODO: exec
	// TODO: exit

	server.REPL(commandManager)

	return nil
}

func (server *Server) REPL(commandManager *utils.CommandManager) {
	// Command and response loop
	for {
		// Step 1: READ command from Operator
		reader := bufio.NewReader(os.Stdin)

		if server.currentAgent != nil {
			fmt.Printf("[operator] [%s] -> ", server.currentAgent.GetAgentID())
		} else {
			fmt.Printf("[operator] [nil] -> ")
		}

		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		// Split command into base and arguments
		cmdSlice := strings.Split(command, " ")
		cmdBase := cmdSlice[0]
		cmdArgs := cmdSlice[1:]

		// NOTE: special snowflake case here
		if cmdBase == "exit" {
			server.Stop()
			fmt.Printf("Exiting\n")
			break
		}

		// Handle the command using the command manager
		resp, err := commandManager.HandleCommand(cmdBase, cmdArgs)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println(resp)
		}

	}
}

func (server *Server) Stop() {
	for agentID, agent := range server.agents {
		err := agent.Stop()
		if err != nil {
			fmt.Printf("Error stopping agent %s: %v\n", agentID, err)
		} else {
			fmt.Printf("Agent %s stopped successfully\n", agentID)
		}
	}
	fmt.Println("Server stopped.")
}
