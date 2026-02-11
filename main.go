package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/m4dmood/black-box/internal/database"
	"github.com/m4dmood/black-box/internal/parser"
	"github.com/m4dmood/black-box/internal/watcher"
)

func main() {
	// 1. Context creation for managing timeouts and cancellations
	ctx := context.Background()
	//defer cancel()

	fmt.Println("[MAIN] Starting Black-Box...")

	// 2. Database connection
	db, err := database.Connect(ctx)
	if err != nil {
		fmt.Printf("[MAIN] Critical error during database startup: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Keep software alive (Graceful Shutdown)
	// Channel to listen for OS signals (e.g., Ctrl+C)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("[MAIN] System is ready. Waiting for CSV files in the 'inbox' folder...")
	fmt.Println("[MAIN] Press Ctrl+C to stop the service.")

	// Start watcher and pass incoming data to channel
	fileQueue := make(chan string, 100)
	fmt.Println("[MAIN] Monitoring directory /inbox...")
	watcher.Watch("./inbox", fileQueue)

	// 4. Elaboration loop: process incoming data from channel
	go func() {
		for filePath := range fileQueue {
			fmt.Printf("[MAIN] File process started: %s\n", filePath)
			entries, err := parser.ProcessFile(filePath)

			if err != nil {
				fmt.Printf("[MAIN] Error while parsing file %s: %v\n", filePath, err)
			}

			for _, entry := range entries {
				err := db.InsertWithFallback(ctx, entry)
				if err != nil {
					fmt.Printf("[MAIN] Error while inserting entry into database: %v\n", err)
				}
			}

			os.Remove(filePath)
		}
	}()

	// Application will keep running until it receives an interrupt signal
	<-stop

	fmt.Println("\n[MAIN] Stop signal received. Shutting down gracefully...")

	// Successful shutdown message
	fmt.Println("[MAIN] Black-Box graceful shutdown completed successfully!")
}
