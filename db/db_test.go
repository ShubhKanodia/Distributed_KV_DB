package db_test

import (
	"os"
	"testing"

	"github.com/shubh/distributed_kv_go/db"
)

func setKey(t *testing.T, d *db.Database, key, value string) {
	t.Helper()
	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("Error setting key %q: %v", key, err)
	}
}

func getKey(t *testing.T, d *db.Database, key string) string {
	t.Helper()
	value, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("Error getting key %q: %v", key, err)
	}
	return string(value)
}

// helper function to create a temporary file with the given contents
func TestDeleteExtraKeys(t *testing.T) {
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

	setKey(t, db, "First", "one")
	setKey(t, db, "Third", "three")

	value := getKey(t, db, "First")

	if value != "one" {
		t.Fatalf("Expected value to be %q, but got %q", "one", value)
	}
	if err := db.DeleteExtraKeys(func(name string) bool { return name == "Third" }); err != nil {
		t.Fatalf("Error deleting extra keys: %v", err)
	}

	if value = getKey(t, db, "First"); value != "one" {
		t.Fatalf("Expected value to be %q, but got %q", "", value)
	}

	//we want to remove 3
	if value = getKey(t, db, "Third"); value != "" {
		t.Fatalf("Expected value to be %q, but got %q", "", value)
	}

}
