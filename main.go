package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/m4dmood/black-box/internal/database"
	"github.com/m4dmood/black-box/internal/parser"
	"github.com/m4dmood/black-box/internal/watcher"
)

func main() {
	// 1. Context creation for managing timeouts and cancellations
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("[BLACK_BOX] Starting Black-Box...")

	// 2. Database connection
	db, err := database.Connect(ctx)
	if err != nil {
		fmt.Printf("[BLACK_BOX] Critical error during database startup: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Keep software alive (Graceful Shutdown)
	// Channel to listen for OS signals (e.g., Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("[BLACK_BOX] System is ready. Waiting for CSV files in the 'inbox' folder...")
	fmt.Println("[BLACK_BOX] Press Ctrl+C to stop the service.")

	// Start watcher and pass incoming data to channel
	fileQueue := make(chan string, 100)
	fmt.Println("[BLACK_BOX] Monitoring directory /inbox...")
	watcher.Watch("./inbox", fileQueue)

	// 4. Elaboration loop: process incoming data from channel
	go func() {
		for filePath := range fileQueue {
			fmt.Printf("[BLACK_BOX] File process started: %s\n", filePath)
			_, err := parser.ProcessFile(filePath)
			if err != nil {
				fmt.Printf("[BLACK_BOX] Error while parsing file %s: %v\n", filePath, err)
			}

			os.Remove(filePath)
		}
	}()

	// Application will keep running until it receives an interrupt signal
	<-stop

	fmt.Println("\n[BLACK_BOX] Stop signal received. Shutting down gracefully...")

	// Successful shutdown message
	fmt.Println("[BLACK_BOX] Black-Box graceful shutdown completed successfully!")
}
