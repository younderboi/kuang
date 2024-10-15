package server

import (
	"fmt"
	"net"
)

func ListenTCP(address string, agentID AgentID) (AgentProxy, error) {
	// NOTE: will need to use a new port for each new agent
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Error starting handler: %v", err)
	}

	println("Listening for TCP agents on address: %s", address)

	agentConn, err := listener.Accept() // Block until agent connects
	if err != nil {
		return nil, err
	}

	println("Agent", agentID, "connected from:", agentConn)

	return &TCPAgentProxy{
		AgentID: agentID,
		Conn:    agentConn,
	}, nil
}
