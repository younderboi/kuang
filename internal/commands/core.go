package commands

import (
	"fmt"
)

func HandleKill(args ...string) (string, error) {
	return "not implemented", nil
}

func HandleClear(args ...string) (string, error) {
	fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	return "", nil
}

func HandlePing(args ...string) (string, error) {
	return "pong", nil
}

func HandleNOP(args ...string) (string, error) {
	return "NOP", nil
}
