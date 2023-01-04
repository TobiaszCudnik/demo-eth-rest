# Demo ETH REST API

Demo of a distributed REST facade fetching data from Ethereum via JSON-RPC over HTTP.

- horizontally scalable
- request grouping
- redis cache
- redis locks
- load balancer
- docker-compose
- separate test container (delve)

## Usage

Terminal 1:
```
$ go install github.com/go-task/task/v3/cmd/task@latest
$ task docker-start
```

Terminal 2:
```
$ http :8080/v1/block/0x5BAD55
$ task benchmark-blocks
```

## Endpoints

- `GET /v1/block/:number`

## Dependencies

- docker
- httpie
- siege
- nodejs

## Taskfile

```
$ task --list
task: Available tasks for this project:
* benchmark-blocks:             Benchmark 100 random blocks
* build:                        Build a binary
* clean:                        Stop and remove all the data
* debug:                        Run in debug mode
* docker-rebuild-start:         Rebuild and start docker compose
* docker-start:                 Start docker compose
* docker-stop:                  Stop docker compose
* flush-cache:                  Flush redis cache
* kill:                         Kill the service process
* rebuild-test-instance:        Rebuild the test instance
* restart-load-balancer:        Restart the load balancer with a new config
* start:                        Start from source
* test:                         Run go test
* test-load-blocks:             Load test 10k random blocks for 1 minute
* test-locks:                   Test if redis locking works
```

## Tests

### Locks

Semi-automated test verifying lock synchronization, executed in a separate container. After running the test, check the logs for "locked, waiting for".

`$ task test-locks`

Output:
```
eth-rest-test-1  | 2022/11/27 11:51:14 [service.getBlock] get 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:14 [service.getBlock] check cache for 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:14 [service.requestBlock] request 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:14 [service.getBlock] get 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:14 [service.getBlock] check cache for 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:14 [service.getBlock] locked, waiting for 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:15 [eth.GetBlock] sleeping for 0.500000
eth-rest-test-1  | 2022/11/27 11:51:15 [service.dispatchCompleted] sending block-completed for 0x5BAD55
eth-rest-test-1  | 2022/11/27 11:51:15 [service.GetBlock] block-completed received for 0x5BAD55
```

### Benchmarks

Benchmarks can be found in [./BENCHMARK.md](BENCHMARK.md).

One example is *full cache, 3 workers, redis 1g, http-rpc* with the following results:

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   0.63 secs
Data transferred:              10.91 MB
Response time:                  0.01 secs
Transaction rate:            3968.25 trans/sec
Throughput:                    17.31 MB/sec
Concurrency:                   23.22
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            0.05
Shortest transaction:           0.00
```

### Load test

`$ task test-load-blocks`

- **130k** hits / minute (10k random blocks)
- **0.04s** avg response time

```
Transactions:                 130736 hits
Availability:                 100.00 %
Elapsed time:                  60.04 secs
Data transferred:             942.20 MB
Response time:                  0.04 secs
Transaction rate:            2177.48 trans/sec
Throughput:                    15.69 MB/sec
Concurrency:                   79.49
Successful transactions:      130833
Failed transactions:               0
Longest transaction:            1.24
Shortest transaction:           0.00
```

## TODO

- retry logic for ETH calls
  - lock support
- throttling for ETH calls
- transactions support the way blocks are implemented
- transactions extracted directly from blocks data
  - kept separately in cache
- benchmarks for transactions
- JSON-RPC via websocket
- unit tests
- fully automated load tests
  - predefined docker-compose files
- more documentation
  - OpenAPI & swagger
