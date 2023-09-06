# Galt Stack - a Golang Salt Stack implemenation attempt

This repo is just an attempt to learn how to make a distributed system in Go using etcd and Salt Stack is a perfect example to do so as it has a lot of clients and server (that should be horizontally scalable) to communicate with back and forth (client sends grains and settings to the server "all the time" and the server can issue commands to the clients)

# This repository contains the following

1. This README
2. A `server/` directory which contains the server part
3. A `client/` directory which contains the client

## Development

1. Get Go 1.21+
2. Clone this repo `git clone https://github.com/bbruun/galt`
3. Change to `galt/server` and run `go mod tidy`
4. Change to `galt/client` and run `go mod tidy`

Dependencies that are used:
* [etcd embed](https://pkg.go.dev/github.com/coreos/etcd/embed)

### To run the server

```
cd galt/server
go run .
```

### To run the client

Note: Start the server first
```
cd galt/client
go run .
```



# Libs to look at

* "github.com/cameronnewman/go-flatten" search for grains easily using dot notation
