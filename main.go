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
	// 1. Creiamo un contesto principale che scadrà se non riusciamo a connetterci
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Avvio del sistema Black-Box...")

	// 2. Inizializziamo la connessione al Database
	db, err := database.Connect(ctx)
	if err != nil {
		fmt.Printf("Errore critico durante l'avvio del database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Meccanismo per tenere vivo il software (Graceful Shutdown)
	// Creiamo un canale che ascolta i segnali del Sistema Operativo
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Sistema pronto. In attesa di file CSV nella cartella 'inbox'...")
	fmt.Println("Premi Ctrl+C per arrestare il servizio.")

	//Invocazione watcher
	fileQueue := make(chan string, 100)
	fmt.Println("Monitoraggio directory /inbox avviato...")
	watcher.Watch("./inbox", fileQueue)

	// 4. Avviamo la logica di business in una goroutine separata
	go func() {
		for filePath := range fileQueue {
			fmt.Printf("Inizio elaborazione: %s\n", filePath)
			_, err := parser.ProcessFile(filePath)
			if err != nil {
				fmt.Printf("Errore nel parsing di %s: %v\n", filePath, err)
			}

			os.Remove(filePath)
		}
	}()

	// Il programma resta "bloccato" qui finché non riceve un segnale su 'stop'
	<-stop

	fmt.Println("\nSegnale di arresto ricevuto. Chiusura in corso...")

	// Qui potremmo aggiungere logiche per finire di processare i file pendenti
	fmt.Println("Black-Box graceful shutdown completato con successo!")
}
