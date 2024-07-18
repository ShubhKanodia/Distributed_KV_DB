package web

import (
	"fmt"
	"net/http"

	"github.com/shubh/distributed_kv_go/db"
)

// Server will have http method handlers to be used for the db
type Server struct {
	db *db.Database
}

// NewServer creates a new server with the given db and http handlers to be used to get and set the key value pair
func NewServer(db *db.Database) *Server {
	return &Server{db: db}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "Value: %q, Error: %v", value, err)
}

// handling write request
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error: %v", err)
}
