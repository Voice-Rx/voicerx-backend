package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"voicerx-backend/auth"
	"voicerx-backend/uploads"
	"voicerx-backend/consumers"
)

func main() {
	m := http.NewServeMux()

	const addr = ":8000"

	m.HandleFunc("/token", auth.HandleLogin)
	m.HandleFunc("/uploadAudio", uploads.HandleAudioUpload)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumers.StartConsumers(ctx)
	}()

	srv := http.Server{
		Handler:      m,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm

		fmt.Println("Received termination signal. Shutting down server")
		cancel() // Cancel the context to stop the consumer Goroutine
		wg.Wait() // Wait for all goroutines to finish
		srv.Shutdown(ctx)
	}()

	fmt.Println("Server started on port ", addr)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
