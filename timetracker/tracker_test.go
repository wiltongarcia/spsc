package timetracker

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	file := &os.File{}
	ttr := New(file)
	if reflect.TypeOf(ttr).String() != "*timetracker.tt" {
		t.Error("Wrong type")
	}
}

func TestLoad(t *testing.T) {
	dir, err := ioutil.TempDir("", "running")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	filepath := fmt.Sprintf("%s/data", dir)

	b := make([]byte, 8)
	now := time.Now().Unix()
	binary.PutVarint(b, now)
	err = ioutil.WriteFile(filepath, b, 0644)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	tt := &tt{file: file}
	tt.Load()
	if len(tt.dataRange) != 1 || tt.dataRange[0] != now {
		t.Error("Error on load data")
	}
}

func TestIncrement(t *testing.T) {
	tt := &tt{}
	now := time.Now().Unix()
	tt.dataRange = append(tt.dataRange, now)
	if tt.Increment() != 2 {
		t.Error("Wrong total")
	}
}

func TestGet(t *testing.T) {
	tt := &tt{}
	now := time.Now().Unix()
	tt.dataRange = append(tt.dataRange, now)
	if tt.Get() != 1 {
		t.Error("Wrong total")
	}
}

func TestWrite(t *testing.T) {
	dir, err := ioutil.TempDir("", "running")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	filepath := fmt.Sprintf("%s/data", dir)
	file, err := os.Create(filepath)
	tt := &tt{}
	tt.file = file
	now := time.Now().Unix()
	tt.dataRange = append(tt.dataRange, now)
	tt.Write()

	file, err = os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	b := make([]byte, 8)
	_, err = file.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	ts, n := binary.Varint(b)
	if n <= 0 {
		log.Fatal("Unable to convert data")
	}

	if ts != now {
		t.Error("Error in write file")
	}
}
