package handler

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/b-sn/go-stub-server/internal/config"
)

func PrepareHandleFuncs(cfg config.Config) map[string]*http.ServeMux {
	servers := make(map[string]*http.ServeMux)

	for port, endpoints := range cfg.Endpoints {
		mux, exists := servers[port]
		if !exists {
			mux = http.NewServeMux()
			servers[port] = mux
		}

		for _, endpoint := range endpoints {
			mux.HandleFunc(endpoint.URL, handleRequest(port, filepath.Join(cfg.ResponsesDir, endpoint.File)))
		}

		// Add 404 handler
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s :%s%s => 404", r.Method, port, r.URL.Path)
			http.Error(w, "Not Found", http.StatusNotFound)
		})
	}

	return servers
}

// handleRequest provides a handler function for the given file path
func handleRequest(port string, filePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Error opening file %s: %v", filePath, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create a new response writer
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			log.Printf("Error hijacking connection: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		var size int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			chunk := scanner.Text()
			size += len(chunk)
			fmt.Fprintln(conn, chunk)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
		}

		log.Printf("%s :%s%s => %s [%d bytes]", r.Method, port, r.URL.Path, filePath, size)
	}
}
