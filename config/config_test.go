package config_test

// go test - go function for testing functionality
// go test -v - go function for testing functionality with verbose output

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/shubh/distributed_kv_go/config"
)

// helper function to create a temporary file with the given contents
func createConfig(t *testing.T, contents string) config.Config {

	f, err := os.CreateTemp(os.TempDir(), "config.toml") //making  a temporary file
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer f.Close()

	name := f.Name()

	defer os.Remove(f.Name()) //removing the file after the test

	if _, err := f.WriteString(contents); err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}

	c, err := config.ParseFile(name)
	fmt.Println(c)
	if err != nil {
		t.Fatalf("Error parsing config: %v", err)
	}
	return *c
}

func TestConfigParse(t *testing.T) {

	got := createConfig(t, `[[shards]]
	name = "shard1"
	idx = 0
	address = "localhost:8080"`)

	want := config.Config{ //what we expect from the parsed file
		Shards: []config.Shard{
			{
				Name:    "shard1",
				Idx:     0,
				Address: "localhost:8080",
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Config doesnt match : Expected %v, but got %v", want, got)
	}
}

func TestParseShards(t *testing.T) {
	c := createConfig(t, `[[shards]] 
	name = "shard1"
	idx = 0
	address = "localhost:8080"
	[[shards]]
	name = "shard2"
	idx = 1
	address = "localhost:8081"`)

	got, err := config.ParseShards(c.Shards, "shard2")

	if err != nil {
		t.Fatalf("Error parsing shards: %v", err)
	}

	want := &config.Shards{
		Count:   2,
		CurrIdx: 1,
		Addrs: map[int]string{
			0: "localhost:8080",
			1: "localhost:8081",
		},
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Shards dont match: Expected %v, but got %v", want, got)
	}

}
