#!/bin/zsh
set -e

trap 'kill $(jobs -p)' SIGINT SIGTERM EXIT

# Change to the script's directory
cd $(dirname $0)

# Kill any running instances
pkill -f distributed_kv_go || true
sleep 0.1

# Install the latest version
go install

# Launch three instances

$GOPATH/bin/distributed_kv_go -db-location=$PWD/shard1.db -config=$PWD/sharding.toml -shard=shard1 -http-addr=127.0.0.1:8080 &
$GOPATH/bin/distributed_kv_go -db-location=$PWD/shard2.db -config=$PWD/sharding.toml -shard=shard2 -http-addr=127.0.0.1:8081 &
$GOPATH/bin/distributed_kv_go -db-location=$PWD/shard3.db -config=$PWD/sharding.toml -shard=shard3 -http-addr=127.0.0.1:8082 &
$GOPATH/bin/distributed_kv_go -db-location=$PWD/shard4.db -config=$PWD/sharding.toml -shard=shard4 -http-addr=127.0.0.1:8083 &


wait