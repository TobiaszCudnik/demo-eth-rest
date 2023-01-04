# Benchmark

**siege** - https://github.com/JoeDog/siege

## Specs

- Intel Xeon D-1531 @ 12x 2.7GHz
- 32gb ram

## Test: 10 blocks

`$ siege -b -c 25 -r 100 -i -f dev/last-10-blocks --no-parser --no-follow`

### no cache, 1 worker, redis 1g, http-rpc

5.72 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   5.72 secs
Data transferred:              10.86 MB
Response time:                  0.05 secs
Transaction rate:             437.06 trans/sec
Throughput:                     1.90 MB/sec
Concurrency:                   21.64
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            5.04
Shortest transaction:           0.00
```

### full cache, 1 worker, redis 1g, http-rpc

0.63 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   0.63 secs
Data transferred:              10.93 MB
Response time:                  0.01 secs
Transaction rate:            3968.25 trans/sec
Throughput:                    17.35 MB/sec
Concurrency:                   23.19
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            0.03
Shortest transaction:           0.00
```

### no cache, 3 workers, redis 1g, http-rpc

1.05 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   1.05 secs
Data transferred:              10.95 MB
Response time:                  0.01 secs
Transaction rate:            2380.95 trans/sec
Throughput:                    10.42 MB/sec
Concurrency:                   24.30
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            0.45
Shortest transaction:           0.00
```

### full cache, 3 workers, redis 1g, http-rpc

0.63 secs

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

## Test: 100 blocks

`$ siege -b -c 25 -r 100 -i -f dev/last-100-blocks --no-parser --no-follow`

### no cache, 1 worker, redis 1g, http-rpc

6.05 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   6.05 secs
Data transferred:              11.82 MB
Response time:                  0.05 secs
Transaction rate:             413.22 trans/sec
Throughput:                     1.95 MB/sec
Concurrency:                   20.92
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            5.02
Shortest transaction:           0.00
```

### full cache, 1 worker, redis 1g, http-rpc

0.63 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   0.63 secs
Data transferred:              11.87 MB
Response time:                  0.01 secs
Transaction rate:            3968.25 trans/sec
Throughput:                    18.84 MB/sec
Concurrency:                   23.48
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            0.05
Shortest transaction:           0.00
```

### no cache, 3 workers, redis 1g, http-rpc

1.43 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   1.43 secs
Data transferred:              11.95 MB
Response time:                  0.01 secs
Transaction rate:            1748.25 trans/sec
Throughput:                     8.35 MB/sec
Concurrency:                   23.97
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            0.52
Shortest transaction:           0.00
```

### full cache, 3 workers, redis 1g, http-rpc

0.61 secs

```
Transactions:                   2500 hits
Availability:                 100.00 %
Elapsed time:                   0.61 secs
Data transferred:              11.96 MB
Response time:                  0.01 secs
Transaction rate:            4098.36 trans/sec
Throughput:                    19.61 MB/sec
Concurrency:                   23.79
Successful transactions:        2500
Failed transactions:               0
Longest transaction:            0.03
Shortest transaction:           0.00
```
