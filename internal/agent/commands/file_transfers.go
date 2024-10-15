package agent

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

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
