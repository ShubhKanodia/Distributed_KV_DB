package db_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/shubh/distributed_kv_go/db"
)

// helper function to create a temporary file with the given contents
func TestGetSet(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "kvdb")

	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}

	name := f.Name()
	f.Close()
	defer os.Remove(name) //removing the file after the test
	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Error creating db: %v", err)
	}

	defer closeFunc()

	if err := db.SetKey("First", []byte("one")); err != nil {
		t.Fatalf("Error setting key: %v", err)
	}

	value, err := db.GetKey("First")
	if err != nil {
		t.Fatalf("Error getting key: %v", err)
	}

	if !bytes.Equal(value, []byte("one")) {
		t.Fatalf("Expected value to be %q, but got %q", "one", value)
	}

}
