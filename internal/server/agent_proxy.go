package server

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
	_, err := fmt.Fprintf(agent.Conn, fullCommand+"\n")
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

func (agent *TCPAgentProxy) Stop() error {
	err := agent.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// =========================
// === Local Agent Proxy ===
// =========================
// TODO: untested
type LocalAgent struct {
	AgentID  AgentID
	CmdChan  chan string
	RespChan chan string
}

func NewAgent(agentID AgentID) (agent *LocalAgent) {
	return &LocalAgent{
		AgentID:  agentID,
		CmdChan:  make(chan string),
		RespChan: make(chan string),
	}
}

func (agent *LocalAgent) Send(CMD string, args []string) error {
	agent.CmdChan <- CMD

	return nil
}

func (agent *LocalAgent) Receive() (string, error) {
	resp := <-agent.RespChan
	return resp, nil
}

func (agent *LocalAgent) Start() {
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

func (agent *LocalAgent) Stop() {
	println("Killing Agent", agent.AgentID)

	agent.RespChan <- fmt.Sprintf("Agent %s killed", agent.AgentID)

	return
}
