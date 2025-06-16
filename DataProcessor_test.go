package main

import "testing"

func TestDataProcessor(t *testing.T) {
	in := make(chan []byte, 100)
	out := make(chan Result, 100)
	go DataProcessor(in, out)
	in <- []byte("1\n+\n3\n4")
	in <- []byte("2\n/\n9\n3")
	close(in)
	expected := map[string]int{"1": 7, "2": 3}
	for result := range out {
		if expected[result.Id] != result.Value {
			t.Errorf("error for ID %s: got %d, want %d", result.Id, result.Value, expected[result.Id])
		}
	}
}
