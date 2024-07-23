package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/shubh/distributed_kv_go/config"
	"github.com/shubh/distributed_kv_go/db"
)

// Server will have http method handlers to be used for the db
type Server struct {
	db *db.Database
	// shardIdx, shardCount int
	// addrs                map[int]string
	shards *config.Shards
}

// NewServer creates a new server with the given db and http handlers to be used to get and set the key value pair
func NewServer(db *db.Database, s *config.Shards) *Server {
	return &Server{
		db: db,
		// shardIdx:   shardIdx,
		// shardCount: shardCount, this has been moved to config
		// addrs:      addrs,
		shards: s,
	}
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shards.CurrIdx, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error redirecting the request: %v", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	shard := s.shards.Index(key)
	value, err := s.db.GetKey(key)

	if shard != s.shards.CurrIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v", shard, s.shards.CurrIdx, s.shards.Addrs[shard], value, err)
}

// handling write request

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.shards.Index(key)
	if shard != s.shards.CurrIdx {
		s.redirect(shard, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIdx = %d, current shard = %d", err, shard, s.shards.CurrIdx)
}

// Delete extraKeys deletes keys that dont belong to the current shard

func (s *Server) DeleteExtraKeysHandler(w http.ResponseWriter, r *http.Request) {
	err := s.db.DeleteExtraKeys(func(name string) bool {
		return s.shards.Index(name) != s.shards.CurrIdx
	})

	fmt.Fprintf(w, "Error = %v", err)
}
