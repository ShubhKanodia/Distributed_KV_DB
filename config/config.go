package config

import (
	"fmt"
	"hash/fnv"

	"github.com/BurntSushi/toml"
)

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

// Shards describes the sharding configuration - shards count, current index and addresses of all other shards
type Shards struct {
	Count   int
	CurrIdx int
	Addrs   map[int]string
}

//Parseshards function - converts and verifies list of shards

//inside config - specied into a form that can be used

// for routing
func ParseFile(filePath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard index %d", s.Idx)
		}
		addrs[s.Idx] = s.Address
		if s.Name == curShardName {
			shardIdx = s.Idx
		}
	}

	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard index %d not found", i)
		}
	}

	if shardIdx < 0 {
		return nil, fmt.Errorf("Shard %q not found", curShardName)
	}

	return &Shards{
		Count:   shardCount,
		CurrIdx: shardIdx,
		Addrs:   addrs,
	}, nil
}

//move this from server to config(server need not know about shards)

func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}
