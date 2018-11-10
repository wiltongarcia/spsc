package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"net/http"

	"github.com/wiltongarcia/spsc/config"
	"github.com/wiltongarcia/spsc/timetracker"
)

func main() {
	// Config
	jsonFile, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatal("Unable to read the config file")
	}
	cfg := config.New(jsonFile)
	configData, err := cfg.Get()
	if err != nil {
		log.Fatal(err)
	}

	// File
	file, err := os.OpenFile(configData.FilePath, os.O_RDWR, 0644)
	if err != nil {
		file, err = os.Create(configData.FilePath)
		if err != nil {
			log.Fatal("Unable create the data file")
		}
	}

	//Time Tracker
	tt := timetracker.New(file)
	tt.Load()

	// Capture the the interrupt signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	// Wait the interrupt signal
	go func(tt timetracker.TimeTracker) {
		select {
		case <-c:
			err := tt.Write()
			if err != nil {
				log.Fatal(err)
			}
			log.Print("Application closed")
			os.Exit(0)
		}
	}(tt)

	// Set Default Handler
	handler := http.NewServeMux()
	handler.HandleFunc("/", defaultHandler(tt, configData))

	// HTTP Server
	hostPort := fmt.Sprintf("%s:%d", configData.Host, configData.Port)
	log.Printf("Initiating server listening at [%s]", hostPort)
	http.ListenAndServe(hostPort, handler)
}

// Default Handler
func defaultHandler(tt timetracker.TimeTracker, configData *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		total := tt.Increment()
		// Set HTTP Response
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, fmt.Sprintf("%d", total))
		if configData.Debug {
			log.Printf("[200] / total: %d", total)
		}
	}
}
