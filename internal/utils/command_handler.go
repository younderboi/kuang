package utils

import "fmt"

type CommandHandler func(args []string) (string, error)

type CommandManager struct {
	handlers map[string]CommandHandler
}

func NewCommandManager() *CommandManager {
	return &CommandManager{
		handlers: make(map[string]CommandHandler),
	}
}

// Register a command handler
func (cm *CommandManager) RegisterHandler(cmd string, handler CommandHandler) {
	cm.handlers[cmd] = handler
}

// Handle the command by calling the appropriate handler
func (cm *CommandManager) HandleCommand(cmd string, args []string) (string, error) {
	if handler, exists := cm.handlers[cmd]; exists {
		return handler(args)
	}
	return "", fmt.Errorf("Command not found: %s", cmd)
}

// TODO: add default handler??
