package watcher

import (
	"log"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func Watch(path string, fileChan chan<-string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					fmt.Printf("Nuovo file rilevato: %s\n", event.Name)
					fileChan <- event.Name
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Errore del watcher:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

}