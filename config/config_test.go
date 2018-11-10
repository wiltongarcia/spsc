package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	file := &os.File{}
	ttr := New(file)
	if reflect.TypeOf(ttr).String() != "*config.cfg" {
		t.Error("Wrong type")
	}
}

func TestGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "running")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	filepath := fmt.Sprintf("%s/json", dir)

	cf := &Config{Debug: true}
	j, err := json.Marshal(cf)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath, j, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("Unable to read the config file")
	}

	c := &cfg{file}
	cf, err = c.Get()
	if !cf.Debug || err != nil {
		t.Error("Error in getting the config values")
	}
}
