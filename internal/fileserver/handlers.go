package fileserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Handle file uploads
func (fs *FileServer) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get the file name from the header
	remoteFileName := r.Header.Get("File-Name")
	if remoteFileName == "" {
		http.Error(w, "File name not provided in header", http.StatusBadRequest)
		return
	}

	// Create the file on the server
	dst, err := os.Create(filepath.Join(fs.serveDir, remoteFileName))
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file's content to the new file
	_, err = io.Copy(dst, r.Body)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	// Log the successful upload
	log.Printf("File uploaded successfully: %s from %s", remoteFileName, r.RemoteAddr)

	// Respond to the client
	fmt.Fprintf(w, "File uploaded successfully: %s\n", remoteFileName)
}

func (fs *FileServer) handleFileDownload(w http.ResponseWriter, r *http.Request) {
	// Check for the file name parameter
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	// Construct the full file path
	filePath := filepath.Join(fs.serveDir, fileName)

	// Serve the file if it exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		log.Printf("File not found: %s from %s", fileName, r.RemoteAddr)
		return
	}

	log.Printf("Serving file: %s to %s", fileName, r.RemoteAddr)
	http.ServeFile(w, r, filePath)
}

// List files in the directory
func (fs *FileServer) handleFileList(w http.ResponseWriter, r *http.Request) {
	files, err := fs.listFilesInDir()
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		log.Printf("Failed to list files: %v from %s", err, r.RemoteAddr)
		return
	}

	// Respond with the file list
	for _, file := range files {
		fmt.Fprintln(w, file)
	}

	log.Printf("Listed files to %s", r.RemoteAddr)
}

// List the files in the directory
func (fs *FileServer) listFilesInDir() ([]string, error) {
	entries, err := os.ReadDir(fs.serveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var fileNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			fileNames = append(fileNames, entry.Name())
		}
	}

	return fileNames, nil
}
