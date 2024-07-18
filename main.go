package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/shubh/distributed_kv_go/config"
	"github.com/shubh/distributed_kv_go/db"
	"github.com/shubh/distributed_kv_go/web"
)

var (
	dbLocation     = flag.String("db-location", "", "The path to the bolt db database")
	httpAddr       = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile     = flag.String("config", "sharding.toml", "static sharding ")
	shard          = flag.String("shard", "", "Name of the shard for the data")
	isReplica      = flag.Bool("replica", false, "Whether this instance is a read-only replica")
	primaryAddress = flag.String("primary", "", "Address of the primary shard (required for replicas)")
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
	if *isReplica && *primaryAddress == "" {
		log.Fatal("primary flag needed for replicas")
	}

}
func main() {
	parseFlags()
	content, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var c config.Config

	if _, err := toml.Decode(string(content), &c); err != nil {
		log.Fatalf("toml.Decode(%q): %v", *configFile, err)
	}
	var shardCount int
	var shardIdx int = -1
	var addrs = make(map[int]string)

	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address

		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	var replicaAddress string
	for _, s := range c.Shards {
		if s.Name == *shard {
			replicaAddress = s.Replica
			break
		}
	}

	if shardIdx < 0 {
		log.Fatalf("Invalid shard name %q", *shard)
	}

	log.Printf("Shard count is %d, shard index is %d", shardCount, shardIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(db, shardIdx, shardCount, addrs, *isReplica, *primaryAddress, replicaAddress)

	if *isReplica {
		http.HandleFunc("/get", srv.GetHandler)
		http.HandleFunc("/sync", srv.SyncHandler)
		//for writing to replicas
		http.HandleFunc("/", srv.ReplicaHandler)
	} else {
		http.HandleFunc("/get", srv.GetHandler)
		http.HandleFunc("/set", srv.SetHandler)
	}

	// for storing data - hash(key)%count = shard index
	log.Fatal(http.ListenAndServe(*httpAddr, nil))

}
