package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/shubh/distributed_kv_go/db"
	"github.com/shubh/distributed_kv_go/web"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
)

func parseFlags() { //function to validate the flags
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("db-location flag is required")
	}
}
func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	parseFlags()

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(db)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))

}
