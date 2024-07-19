# Distributed Key-Value Database

This is a lightweight, distributed key-value database built with Go. It features hash-based sharding, read-only replicas, and HTTP-based communication for efficient data storage and retrieval across multiple nodes.

![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=for-the-badge&logo=go)

## Features

- **Distributed Architecture**: Utilizes hash-based sharding for data distribution across multiple nodes.
- **Read-Only Replicas**: Enhances read scalability with dedicated read-only replica nodes.
- **Eventual Consistency**: Implements asynchronous replication for eventual consistency between primary shards and replicas.
- **HTTP-Based Communication**: Enables seamless inter-shard operations and data retrieval.
- **BoltDB Backend**: Leverages BoltDB for efficient local storage on each node.
- **Concurrent Operations**: Uses Go's goroutines for non-blocking replica updates.

## Architecture

The distributed key-value store consists of multiple shards, each responsible for a subset of the key space. Each shard has a primary node for read and write operations, and a read-only replica for enhanced read performance.

```
Client
  |
  v
Load Balancer
 / | \
v  v  v
Shard1  Shard2  Shard3
  |     |      |
  v     v      v
Replica1 Replica2 Replica3
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- BoltDB

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/ShubhKanodia/Distributed_KV_DB.git
   ```

2. Navigate to the project directory:
   ```sh
   cd Distributed_KV_DB
   ```

3. Build the project:
   ```sh
   go build
   ```

### Running the Distributed Key-Value Database

Use the provided `launch.sh` script to start multiple shards and their replicas:
```sh
./launch.sh
```
This script will launch three primary shards and their corresponding replicas.

## Usage

### Setting a Key-Value Pair
```sh
curl "http://localhost:8080/set?key=exampleKey&value=exampleValue"
```

### Getting a Value
```sh
curl "http://localhost:8080/get?key=exampleKey"
```

## Demo

### Normal Get and Set Operations
<img width="625" alt="Screenshot 2024-07-19 at 12 35 20 AM" src="https://github.com/user-attachments/assets/62a93b8a-e808-4513-b309-c3c5a4b81d7c">


### Shard Redirection
When querying any shard, the database automatically redirects the request to the appropriate shard:

#### On querying shard1 at localhost:8081 for the value at shard0
<img width="635" alt="Screenshot 2024-07-19 at 1 07 08 AM" src="https://github.com/user-attachments/assets/aa77de9c-62dd-412f-86bb-2b2a64b094ea">

#### On querying shard2 at localhost:8082 for the value at shard0
<img width="640" alt="Screenshot 2024-07-19 at 1 10 15 AM" src="https://github.com/user-attachments/assets/8e28c9e9-0223-47e6-ba85-bef318fef589">

#### querying shard0 for the value at shard2
<img width="628" alt="Screenshot 2024-07-19 at 3 55 17 PM" src="https://github.com/user-attachments/assets/a9bc8175-c5d0-482a-8b83-e2ce73c1295b">

### Replication and Read-Only Replicas

Replicas support read operations but reject write attempts:

#### Setting value at primary(shard1) on localhost:8081
<img width="890" alt="Screenshot 2024-07-19 at 3 48 23 AM" src="https://github.com/user-attachments/assets/60e67932-4c1c-4522-9a9a-d895b64b2ee6">

#### Querying its replica on locahost:8084 for read
<img width="857" alt="Screenshot 2024-07-19 at 3 48 30 AM" src="https://github.com/user-attachments/assets/fe01e8aa-31fb-4877-8676-eebb69436241">

#### Trying to write a value to its replica using 'set' operation(rejected)
<img width="838" alt="Screenshot 2024-07-19 at 3 48 37 AM" src="https://github.com/user-attachments/assets/e59f3eed-e911-4318-b743-e899ae301a5f">

## Configuration

The database uses a TOML configuration file (`sharding.toml`) to define shard and replica settings:

```toml
[[shards]]
name = "shard1"
idx = 0
address = "localhost:8080"
replica = "localhost:8083"

[[shards]]
name = "shard2"
idx = 1
address = "localhost:8081"
replica = "localhost:8084"

[[shards]]
name = "shard3"
idx = 2
address = "localhost:8082"
replica = "localhost:8085"
```

---

This version should render correctly in Markdown.
