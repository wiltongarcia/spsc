package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var err error

// Model struct
type Config struct {
	Debug    bool   `json:"debug"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	FilePath string `json:"filepath"`
}

// Struct of module
type cfg struct {
	file *os.File
}

// Server interface is used for defining all methods
type Configuration interface {
	Get() (*Config, error)
}

// Contructor
func New(file *os.File) Configuration {
	return &cfg{file}
}

// Return the Model with config data
func (c *cfg) Get() (*Config, error) {
	// Read the file
	b, err := ioutil.ReadAll(c.file)
	if err != nil {
		return nil, err
	}
	// Parse the json data
	cg := &Config{}
	err = json.Unmarshal([]byte(b), cg)
	if err != nil {
		return cg, err
	}

	return cg, nil
}
