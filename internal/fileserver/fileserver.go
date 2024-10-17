package fileserver

/*
* TODO: currently vulnerable to file traversal attacks
* TODO: no limit to file upload size
 */

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// FileServer struct to encapsulate the server state
type FileServer struct {
	serveDir string
	port     string
	server   *http.Server
}

// Create a new instance of FileServer
func NewFileServer(serveDir, port string) *FileServer {
	return &FileServer{
		serveDir: serveDir,
		port:     port,
	}
}

// Start the file server
func (fs *FileServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/list", fs.logRequest(fs.handleFileList))
	mux.HandleFunc("/download", fs.logRequest(fs.handleFileDownload))
	mux.HandleFunc("/upload", fs.logRequest(fs.handleFileUpload))

	fs.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", fs.port),
		Handler: mux,
	}

	fmt.Printf("File server running on :%s\n", fs.port)
	err := fs.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v\n", err)
		return err
	}
	return nil
}

// Stop the file server gracefully
func (fs *FileServer) Stop() {
	if fs.server != nil {
		fmt.Println("Shutting down file server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := fs.server.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down server: %v\n", err)
		} else {
			fmt.Println("Server shutdown completed.")
		}
	}
}

// Logging middleware to log incoming requests
func (fs *FileServer) logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("[%s] %s %s from %s\n", time.Now().Format(time.RFC3339), method, path, remoteIP)
		next(w, r)
	}
}
