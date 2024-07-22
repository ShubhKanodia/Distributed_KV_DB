package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/shubh/distributed_kv_go/config"
	"github.com/shubh/distributed_kv_go/db"
	"github.com/shubh/distributed_kv_go/web"
)

var (
	dbLocation = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile = flag.String("config", "sharding.toml", "static sharding ")
	shard      = flag.String("shard", "", "Name of the shard for the data")
)

func parseFlags() { //function to validate the flags
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("db-location flag is required")
	}
	if *httpAddr == "" {
		log.Fatal("http-addr flag is required")
	}
	if *configFile == "" {
		log.Fatal("config flag is required")
	}
	if *shard == "" {
		log.Fatal("shard flag is required")
	}

}
func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	parseFlags()

	c, err := config.ParseFile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config %q: %v", *configFile, err)
	}

	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config: %v", err)
	}

	log.Printf("Shard count is %d, current shard: %d", shards.Count, shards.CurrIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("Error creating %q: %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(db, shards)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))

}
