package server

import (
	"fmt"
)

func (server *Server) handleClear(args ...string) (string, error) {
	fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	return "", nil
}

func (server *Server) handleAgentPassthrough(args ...string) (string, error) {
	fmt.Println("Agent command passthrough with args:", args)

	err := server.agentProxy.Send(args)
	if err != nil {
		return "", fmt.Errorf("Failed to send command to agent %s: %v", server.agentProxy.GetAgentID(), err)
	}

	resp, err := server.agentProxy.Receive()
	if err != nil {
		return "", fmt.Errorf("Receiving response from agent %s failed: %v", server.agentProxy.GetAgentID(), err)
	}

	// fmt.Println("Go response from agent:", resp)

	return resp, nil
}
