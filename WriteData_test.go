package main

import (
	"errors"
	"os"
	"sync"
	"testing"
)

func createFile(t *testing.T) (_ *os.File, err error) {
	f, err := os.Create("tempFile")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = errors.Join(err)
	}()
	t.Cleanup(func() {
		f.Close()
		os.Remove(f.Name())
	})
	return f, nil
}

func TestWriteData(t *testing.T) {
	ch2 := make(chan Result, 2)
	f, err := createFile(t)
	if err != nil {
		t.Fatalf("failed to read create temp file: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		WriteData(ch2, f)
	}()

	ch2 <- Result{Id: "1", Value: 5}
	ch2 <- Result{Id: "2", Value: 6}
	close(ch2)
	wg.Wait()
	f.Sync()
	content, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	expected := "1:5\n2:6\n"
	if string(content) != expected {
		t.Errorf("Expected: %s , got: %s", expected, string(content))
	}

}
