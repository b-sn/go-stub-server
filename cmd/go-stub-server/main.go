package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/b-sn/go-stub-server/internal/config"
	"github.com/b-sn/go-stub-server/internal/handler"
)

func main() {
	cfg := config.MustLoadConfig()

	servers := handler.PrepareHandleFuncs(cfg)

	var wg sync.WaitGroup
	for port, mux := range servers {
		wg.Add(1)

		go func(port string, mux *http.ServeMux) {
			defer wg.Done()
			server := &http.Server{
				Addr:    fmt.Sprintf(":%s", port),
				Handler: mux,
			}
			log.Printf("Server started on port %s", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("Error starting server on port %s: %v", port, err)
			}
		}(port, mux)
	}

	wg.Wait()
}
