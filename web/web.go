package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"

	"github.com/shubh/distributed_kv_go/db"
)

// Server will have http method handlers to be used for the db
type Server struct {
	db                   *db.Database
	shardIdx, shardCount int
	addrs                map[int]string
	isReplica            bool
	primaryAddress       string
	replicaAddress       string
}

// NewServer creates a new server with the given db and http handlers to be used to get and set the key value pair
func NewServer(db *db.Database, shardIdx, shardCount int, addrs map[int]string, isReplica bool, primaryAddress string, replicaAddress string) *Server {
	return &Server{
		db:             db,
		shardIdx:       shardIdx,
		shardCount:     shardCount,
		addrs:          addrs,
		isReplica:      isReplica,
		primaryAddress: primaryAddress,
		replicaAddress: replicaAddress,
	}
}
func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shardCount))
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shardIdx, shard, url)

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

	shard := s.getShard(key)
	value, err := s.db.GetKey(key)

	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v", shard, s.shardIdx, s.addrs[shard], value, err)
}

// handling write request
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	if s.isReplica {
		http.Error(w, "This is a read-only replica", http.StatusForbidden)
		return
	}
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.getShard(key)
	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIdx = %d, current shard = %d", err, shard, s.shardIdx)
	go s.notifyReplica(key, value)
}
func (s *Server) notifyReplica(key, value string) {
	if s.replicaAddress == "" {
		return
	}
	url := fmt.Sprintf("http://%s/sync?key=%s&value=%s", s.replicaAddress, key, value)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Error notifying replica: %v", err)
	}
}

func (s *Server) SyncHandler(w http.ResponseWriter, r *http.Request) {
	if !s.isReplica {
		http.Error(w, "This is not a replica", http.StatusForbidden)
		return
	}

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	err := s.db.SetKey(key, []byte(value))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error syncing: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// to handle writing to replicas and display appropriate error log

func (s *Server) ReplicaHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/set" {
		http.Error(w, "This is a read-only replica. Set operations are not allowed.", http.StatusForbidden)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
