package commands

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// === Shell exec
func HandleRunShellCommand(args ...string) (string, error) {
	switch len(args) {
	case 0:
		return "", fmt.Errorf("usage: run [command] <args ...> ")
	case 1:
		return runShellCommand(args[0])
	default:
		return runShellCommand(args[0], args[1:]...)
	}
}

func runShellCommand(command string, args ...string) (string, error) {
	// TODO: how will this behave on windows platform with for example powershell?

	// Create a context with a timeout (e.g., 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the command with the given context (to allow for a timeout)

	cmd := exec.CommandContext(ctx, command, args...)

	// Capture the standard output and error
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("failed to start command: %v", err)
	}

	// Wait for the command to finish or timeout
	err := cmd.Wait()

	// If the context's deadline was exceeded (i.e., a timeout occurred)
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command timed out")
	}

	// If there was an error, return the error message and stderr
	if err != nil {
		return fmt.Sprintf("Command error: %v\nStderr: %s", err, stderr.String()), err
	}

	// Return the output, or a default success message if there’s no output
	if stdout.Len() == 0 && stderr.Len() == 0 {
		return "Command executed successfully (no output)\n", nil
	}

	return stdout.String(), nil
}
