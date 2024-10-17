package operatorcli

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type AgentID string

type AgentProxy interface {
	GetAgentID() AgentID
	Send(args []string) error
	Receive() (string, error)
	IsAlive() bool
	// Start() (AgentID, error)
	Stop() error
}

// =======================
// === TCP AGENT PROXY ===
// =======================
type TCPAgentProxy struct {
	AgentID AgentID
	Conn    net.Conn
}

func (agent *TCPAgentProxy) GetAgentID() AgentID {
	return agent.AgentID
}

func (agent *TCPAgentProxy) Send(args []string) error {
	// Join the command and arguments
	fullCommand := strings.Join(args, " ") + "\n"

	_, err := fmt.Fprintf(agent.Conn, fullCommand)
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}
	return nil
}

func (agent *TCPAgentProxy) Receive() (string, error) {
	agentReader := bufio.NewReader(agent.Conn)
	var responseBuilder strings.Builder

	// Read and build the response line by line
	for {
		response, err := agentReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return "", err
		}

		if response == "END_OF_RESPONSE\n" {
			break
		}

		// Append each line to the responseBuilder
		responseBuilder.WriteString(response)

	}

	// Return the complete response
	return responseBuilder.String(), nil
}

func (agent *TCPAgentProxy) IsAlive() bool {
	// TODO: slightly oversimplistic??
	err := agent.Send([]string{"ping"})
	if err != nil {
		return false
	}
	_, err = agent.Receive()
	if err != nil {
		return false
	}

	// TODO: check the actual response for pong

	return true
}

func (agent *TCPAgentProxy) Stop() error {
	return agent.Conn.Close()
}
