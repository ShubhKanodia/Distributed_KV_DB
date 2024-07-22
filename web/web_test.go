package web_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/shubh/distributed_kv_go/config"
	"github.com/shubh/distributed_kv_go/db"
	"github.com/shubh/distributed_kv_go/web"
)

func createShardDb(t *testing.T, idx int) *db.Database {
	t.Helper()

	tmpFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("db%d", idx))
	if err != nil {
		t.Fatalf("Could not create a temp db %d: %v", idx, err) //will tell creatin of which shard failed(first or second)
	}

	tmpFile.Close()

	name := tmpFile.Name()
	t.Cleanup(func() { os.Remove(name) }) //cleanup func will be called when test finishes

	db, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create new database %q: %v", name, err)
	}
	t.Cleanup(func() { closeFunc() })

	return db
}

func createShardServer(t *testing.T, idx int, addrs map[int]string) (*db.Database, *web.Server) {
	t.Helper()

	db := createShardDb(t, idx)

	cfg := &config.Shards{
		Addrs:   addrs,
		Count:   len(addrs),
		CurrIdx: idx,
	}

	s := web.NewServer(db, cfg)
	return db, s
}

func TestWebServer(t *testing.T) {
	var ts1GetHandler, ts1SetHandler func(w http.ResponseWriter, r *http.Request)
	var ts2GetHandler, ts2SetHandler func(w http.ResponseWriter, r *http.Request)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts1GetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts1SetHandler(w, r)
		}
	}))
	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts2GetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts2SetHandler(w, r)
		}
	}))
	defer ts2.Close()

	addrs := map[int]string{
		0: strings.TrimPrefix(ts1.URL, "http://"),
		1: strings.TrimPrefix(ts2.URL, "http://"),
	}

	db1, web1 := createShardServer(t, 0, addrs)
	db2, web2 := createShardServer(t, 1, addrs)

	keys := map[string]int{
		"First": 1,
		"Third": 0,
	}
	//just to check which key goes to which shard before writing the test
	// s := &config.Shards{
	// 	Addrs:   addrs,
	// 	Count:   len(addrs),
	// 	CurrIdx: 0,
	// }

	// log.Printf("First: %d", s.Index("First"))
	// log.Printf("Third: %d", s.Index("Third"))

	ts1GetHandler = web1.GetHandler
	ts1SetHandler = web1.SetHandler
	ts2GetHandler = web2.GetHandler
	ts2SetHandler = web2.SetHandler

	for key := range keys {
		// log.Printf("Setting key %s\n", key)
		resp, err := http.Get(fmt.Sprintf(ts1.URL+"/set?key=%s&value=value-%s", key, key))
		if err != nil {
			t.Fatalf("Could not set key %q: %v", key, err)
		}
		resp.Body.Close()
	}

	for key := range keys {
		// log.Printf("Getting key %s\n", key)
		resp, err := http.Get(fmt.Sprintf(ts1.URL+"/get?key=%s", key))
		if err != nil {
			t.Fatalf("Get key %q error: %v", key, err)
		}

		contents, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Fatalf("Could read contents of the key %q: %v", key, err)
		}

		want := []byte("value-" + key)
		if !bytes.Contains(contents, want) {
			t.Errorf("Unexpected contents of the key %q: got %q, want the result to contain %q", key, contents, want)
		}

		log.Printf("Contents of key %q: %s", key, contents)
	}

	value1, err := db1.GetKey("Third")
	if err != nil {
		t.Fatalf("Third key error: %v", err)
	}

	want1 := "value-Third"
	if !bytes.Equal(value1, []byte(want1)) {
		t.Errorf("Unexpected value of Third key: got %q, want %q", value1, want1)
	}

	value2, err := db2.GetKey("First")
	if err != nil {
		t.Fatalf("First key error: %v", err)
	}

	want2 := "value-First"
	if !bytes.Equal(value2, []byte(want2)) {
		t.Errorf("Unexpected value of First key: got %q, want %q", value2, want2)
	}
}
