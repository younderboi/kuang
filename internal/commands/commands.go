package commands

import "fmt"

// CommandHandler defines the function type for handling commands
type CommandHandler func(args ...string) (string, error)

// CommandManager manages command handlers and executes them
type CommandManager struct {
	handlers map[string]CommandHandler
}

// NewCommandManager initializes a new CommandManager
func NewCommandManager() *CommandManager {
	return &CommandManager{
		handlers: make(map[string]CommandHandler),
	}
}

// RegisterHandler registers a command handler for a specific command
func (cm *CommandManager) RegisterHandler(cmd string, handler CommandHandler) {
	cm.handlers[cmd] = handler
}

// RegisterDefaultHandler registers a default handler for unmatched commands
func (cm *CommandManager) RegisterDefaultHandler(handler CommandHandler) {
	cm.handlers["default"] = handler
}

// HandleCommand executes the appropriate command handler, based on the command
func (cm *CommandManager) HandleCommand(cmd string, args ...string) (string, error) {
	fmt.Println("Handling command:", cmd, "with args:", args)
	if handler, exists := cm.handlers[cmd]; exists {
		// Call the specific command handler if it exists
		fmt.Println("Found handler", handler)
		return handler(args...)
	}

	// Call the default handler if no command matches
	if defaultHandler, exists := cm.handlers["default"]; exists {
		return defaultHandler(args...)
	}

	// If no match and no default handler, return an error
	return "", fmt.Errorf("unrecognized command: %s", cmd)
}
