package operatorcli

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Server struct {
	address    string
	port       string
	agentProxy AgentProxy
}

func NewServer(bindAddress string) *Server {
	return &Server{
		address:    bindAddress,
		agentProxy: nil,
	}
}

func (server *Server) Start() error {
	server.REPL()
	return nil
}

func (server *Server) REPL() {
	// TODO: Move to .Start()
	// TODO: Separate out the listen logic
	// TODO:
	var listener net.Listener
	var err error

	// Start listening for new connections
	listener, err = net.Listen("tcp", server.address)
	if err != nil {
		fmt.Printf("Failed to start listener: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		// Check if the agent is connected and alive
		// TODO: move IsAlive check to en of loop?
		if server.agentProxy == nil || !server.agentProxy.IsAlive() {
			if server.agentProxy != nil {
				fmt.Println("Agent disconnected.")
				server.agentProxy.Stop() // Clean up the old connection
				server.agentProxy = nil
			}

			// Wait for a new agent connection
			fmt.Println("Waiting for a new agent...")
			agentConn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection: %v\n", err)
				continue
			}

			// Create a new AgentProxy for the new connection
			server.agentProxy = &TCPAgentProxy{
				AgentID: AgentID(agentConn.RemoteAddr().String()),
				Conn:    agentConn,
			}

			fmt.Printf("New agent connected from %s.\n", agentConn.RemoteAddr().String())

			continue
		}

		// === Read from Operator
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("[%s] -> ", server.agentProxy.GetAgentID())

		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			server.Stop()
			fmt.Println("Exiting.")
			break
		}

		if command == "" {
			continue
		}

		// === Evaluate command
		cmdSlice := strings.Split(command, " ")

		err := server.agentProxy.Send(cmdSlice)
		if err != nil {
			fmt.Println("Failed to send command to agent %s: %v\n", server.agentProxy.GetAgentID(), err)
		}

		resp, err := server.agentProxy.Receive()
		if err != nil {
			fmt.Printf("Receiving response from agent %s failed: %v\n", server.agentProxy.GetAgentID(), err)
		}

		fmt.Println(resp)

	}
}

func (server *Server) Stop() {
	server.agentProxy.Stop()

	fmt.Println("Server stopped.")
}
