package config

// Shard describes a a db partition that holds set of keys, each shard has unique keys and values
type Shard struct {
	Name    string
	Idx     int
	Address string
}

// Config describes the sharding configuration
type Config struct {
	Shards []Shard
}
