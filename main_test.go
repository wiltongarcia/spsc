package main

import (
	"bytes"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/wiltongarcia/spsc/config"
	"github.com/wiltongarcia/spsc/timetracker"
)

func TestDefaultHandler(t *testing.T) {

	file := &os.File{}

	tt := timetracker.New(file)

	cfg := &config.Config{}
	handler := defaultHandler(tt, cfg)

	r := httptest.NewRequest("GET", "/", bytes.NewReader([]byte("")))
	w := httptest.NewRecorder()

	handler(w, r)

	if w.Code != 200 {
		t.Error("Wrong response!")
	}
}
