package agent

import (
	"fmt"
	"strings"
	"time"

	"konstantinovitz.com/kuang/internal/commands"
)

type Agent struct {
	Transport      AgentTransport
	MaxRetries     int
	BaseDelay      time.Duration
	CommandManager *commands.CommandManager
}

func (agent *Agent) Start() error {
	retries := 0
	for {
		// TODO: migrate to Transport.Connnect
		err := agent.Transport.Connect()
		if err != nil {
			if retries >= agent.MaxRetries {
				return fmt.Errorf("Max retries reached. Exiting.")
			}
			retries++
			time.Sleep(agent.backoff(retries))
		} else {
			fmt.Println("Connected successfully.")
			agent.REPL()
			retries = 0
		}
	}
}

func (agent *Agent) REPL() {
	for {
		// Read command from the transport
		command, err := agent.Transport.Read()
		if err != nil {
			fmt.Println("Error reading command:", err)
			// Attempt to reconnect or exit if needed
			return
		}
		command = strings.TrimSpace(command)

		fmt.Println("Received command:", command)

		cmdSlice := strings.Split(command, " ")

		res, err := agent.CommandManager.HandleCommand(cmdSlice[0], cmdSlice[1:]...) // TODO: what if command slice length < 2?
		// TODO: make the agent correctly write erorr messages back to the server

		// Send the response back using the transport
		if err := agent.Transport.Write(res + "\nEND_OF_RESPONSE"); err != nil {
			fmt.Println("Error sending response:", err)
			return
		}
	}
}

func (agent *Agent) Stop() error {
	err := agent.Transport.Close()
	return err
}
