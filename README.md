# Dummkopf

A stupid HTTP test server.

## Features

- Dynamic responses on any path (except `/metrics`).

The following query string parameters are supported:

| Parameter | Description | Example |
|-----------|-------------|---------|
| delay     | Delays the response. CSV list of `duration:percent`. | `delay=100us:10,50ms:50,0.5s:5` |
| status    | Status code to return instead of 200. CSV list of `code:percent`. | `status=410:5,500:20` |

- `/metrics` returns Prometheus metrics.

In addition to the default Prometheus metrics, the following additional metrics
are returned:

| Metric                        | Description |
|-------------------------------|-------------|
| http_requests_total           | Total number of HTTP requests processed, partitioned by status code. |
| http_request_duration_seconds | HTTP request duration distribution. |

## Getting Started

Requires Go 1.13+ and Go module support.

Build and run:
```
$ make
$ ./dummkopf
```

Get help:
```
$ ./dummkopf -help
```

## Example Requests

Return a 200 status with no delay:
```
$ curl localhost:9000
```

Return a 502 status with a 50 ms delay:
```
$ curl 'localhost:9000/?status=502:100&delay=50ms:100'
```

Return a 200 status 75% of the time and a 502 status 25% of the time, with an
80% chance of a 25 ms delay and a 20% chance of a 200 ms delay:
```
$ curl 'localhost:9000/?status=502:25&delay=25ms:80,200ms:20'
```

Beat it up!
```
$ fortio load -n 1000 -qps 0 -c 10 -allow-initial-errors -jitter 'localhost:9000/?status=401:3,500:10&delay=50ms:50,100ms:25,250ms:10,0.5s:10,1s:5'
```

Return Prometheus metrics:
```
$ curl localhost:9000/metrics
```
