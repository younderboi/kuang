package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// === NOP Agent implementation
type Agent struct {
	AgentID  string
	CmdChan  chan string
	RespChan chan string
}

func NewAgent(agentID string) (agent *Agent) {
	return &Agent{
		AgentID:  agentID,
		CmdChan:  make(chan string),
		RespChan: make(chan string),
	}
}

func (agent *Agent) Start() {
	println("Agent", agent.AgentID, "listening")

	for {
		cmd := <-agent.CmdChan
		println("Received command:", cmd)

		if cmd == "kill" {
			agent.Stop()
			break
		}

		agent.RespChan <- fmt.Sprintf("Agent %s received command: %s", agent.AgentID, cmd)
	}

	return
}

func (agent *Agent) Stop() {
	println("Killing Agent", agent.AgentID)

	agent.RespChan <- fmt.Sprintf("Agent %s killed", agent.AgentID)

	return
}

// === Handler Server code
type Server struct {
	agents       map[string]*Agent
	currentAgent *Agent
}

func NewServer() (server *Server) {
	return &Server{
		agents:       make(map[string]*Agent),
		currentAgent: nil,
	}
}

func (server *Server) Start() {
	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("$=: ")

		command, err := reader.ReadString('\n')
		if err != nil {
			println("Error reading input:", err)
			// TODO:
			break
		}

		command = strings.TrimSpace(command)

		// Split command into base and arguments
		cmdSlice := strings.Split(command, " ")
		cmdBase := cmdSlice[0]
		cmdArgs := cmdSlice[1:]

		println("Received input:", command)

		switch cmdBase {

		case "help":
			println("todo")

		case "clear":
			println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")

		case "exit":
			// TODO: add stop logic
			// loop through server.agents and call .Stop()
			for _, agent := range server.agents {
				agent.Stop()
			}
			break

		case "add":
			if len(cmdArgs) < 1 {
				println("add <agentID : string>")
				//  TODO: break somehow
			}

			agentID := cmdArgs[0]

			// TODO: Check agentID collision

			//_, exists? := server.agents[agentID]

			server.agents[agentID] = NewAgent(agentID)

			if server.currentAgent == nil {
				server.currentAgent = server.agents[agentID]
			}

			go server.agents[agentID].Start()

			println("Added agent", agentID)

		case "attach":
			server.currentAgent = server.agents[cmdArgs[0]]
			println("Agent set to:", server.currentAgent)

		case "kill":
			agentID := cmdArgs[0]
			delete(server.agents, agentID)

		case "list":
			for agentID, agent := range server.agents {
				fmt.Printf("=== Agent ID: %s ===\n%+v\n==================\n", agentID, agent)
			}

		// Default to command passthrough
		default:
			if server.currentAgent != nil {
				server.currentAgent.CmdChan <- command
				resp := <-server.currentAgent.RespChan
				println("Server got read response:", resp)
			} else {
				println("No agent set, use 'attach <agentID>'")
			}

		}

	}
}

func (server *Server) Stop() error {
	return nil
}

func main() {
	server := NewServer()
	server.Start()
	println("Exiting")
}
