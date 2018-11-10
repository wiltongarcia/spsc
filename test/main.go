package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/wiltongarcia/spsc/config"
	"github.com/wiltongarcia/spsc/timetracker"
)

func main() {
	ms := flag.Int("ms", 1000, "Interval of request in milliseconds")
	flag.Parse()
	for {
		// Config
		cfg := config.New("../config/config.json")
		configData, err := cfg.Get()
		if err != nil {
			log.Fatal("Unable to read the config file")
		}

		// Request
		URL := fmt.Sprintf("http://%s:%d/", configData.Host, configData.Port)
		request, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			log.Fatal("Unable to create request")
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			// File
			file, err := os.OpenFile(configData.FilePath, os.O_RDWR, 0644)
			if err != nil {
				log.Fatal("Unable read the data file")
			}

			// Time Tracker
			tt := timetracker.New(file)
			tt.Load()

			log.Printf("Total from file: %d", tt.Get())
		} else {
			defer response.Body.Close()
			total, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Fatal("Unable to read response")
			}
			log.Print(string(total[:]))
		}

		time.Sleep(time.Duration(*ms) * time.Millisecond)
	}
}
