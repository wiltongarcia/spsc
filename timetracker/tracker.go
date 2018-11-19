package timetracker

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"time"
)

type TimeTracker interface {
	Load() error
	Increment() int
	Get(param ...int64) int
	Write() error
}

type tt struct {
	file      *os.File
	dataRange []int64
	semaphore chan struct{}
}

// Constructor
func New(file *os.File) TimeTracker {
	return &tt{file: file, semaphore: make(chan struct{}, 1)}
}

// Read the binary data from file
func (t *tt) Load() error {
	b := make([]byte, 8)
	for {
		if _, err := t.file.Read(b); err == io.EOF {
			break
		}
		ts, n := binary.Varint(b)
		if n <= 0 {
			return errors.New("Unable to convert data")
		}
		t.dataRange = append(t.dataRange, ts)
	}
	return nil
}

// Increment the total
func (t *tt) Increment() int {
	// wait a slot
	t.semaphore <- struct{}{}
	defer func() {
		<-t.semaphore // read to release a slot
	}()

	// Actual timestamp
	now := time.Now().Unix()
	// Append the current Request
	t.dataRange = append(t.dataRange, now)

	return t.Get(now)
}

// Increment the total
func (t *tt) Get(param ...int64) int {
	var now int64
	var total int

	if len(param) > 0 {
		now = param[0]
	} else {
		now = time.Now().Unix()
	}
	// Set limit to 60 seconds
	limit := now - 60
	// Loop (DESC) the time tracker and increment the total
	for i := len(t.dataRange) - 1; i >= 0; i-- {
		if t.dataRange[i] > limit {
			total++
		} else {
			break
		}
	}
	return total
}

// Write the range in the file
func (t *tt) Write() error {
	t.file.Truncate(0)
	for i := 0; i < len(t.dataRange); i++ {
		b := make([]byte, 8)
		binary.PutVarint(b, t.dataRange[i])
		_, err := t.file.Write(b)
		if err != nil {
			return err
		}
	}
	t.file.Close()
	return nil
}
