#!/bin/zsh
set -e

trap 'kill $(jobs -p)' SIGINT SIGTERM EXIT

# Change to the script's directory
cd $(dirname $0)

# Kill any running instances
pkill -f distributed_kv_go || true
sleep 0.1

# Build the project
go build -o distributed_kv_go

# Launch three primary instances with their replicas
./distributed_kv_go -db-location=shard1.db -config=sharding.toml -shard=shard1 -http-addr=127.0.0.1:8080 &
./distributed_kv_go -db-location=shard1_replica.db -config=sharding.toml -shard=shard1 -http-addr=127.0.0.1:8083 -replica -primary=127.0.0.1:8080 &

./distributed_kv_go -db-location=shard2.db -config=sharding.toml -shard=shard2 -http-addr=127.0.0.1:8081 &
./distributed_kv_go -db-location=shard2_replica.db -config=sharding.toml -shard=shard2 -http-addr=127.0.0.1:8084 -replica -primary=127.0.0.1:8081 &

./distributed_kv_go -db-location=shard3.db -config=sharding.toml -shard=shard3 -http-addr=127.0.0.1:8082 &
./distributed_kv_go -db-location=shard3_replica.db -config=sharding.toml -shard=shard3 -http-addr=127.0.0.1:8085 -replica -primary=127.0.0.1:8082 &

wait