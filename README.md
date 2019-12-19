# redis-cui

[![Build Status](https://travis-ci.com/kyai/redis-cui.svg?branch=master)](https://travis-ci.com/kyai/redis-cui)
[![GolangCI](https://golangci.com/badges/github.com/kyai/redis-cui.svg)](https://golangci.com)
[![Go Report Card](https://goreportcard.com/badge/github.com/kyai/redis-cui)](https://goreportcard.com/report/github.com/kyai/redis-cui)

Simple, visual command line tool for redis.

![demo](/docs/demo.gif)

## Feature

* Simple and visual
* Server friendly
* Supported vim keys
* Same arguments as `redis-cli`

## Installation

```go
go get github.com/kyai/redis-cui
```

Or download the binary and add it to your `PATH`

[Latest release](https://github.com/kyai/redis-cui/releases/latest)

## Usage

```
$ redis-cui
```

```
Usage: redis-cui [OPTIONS]

  -h <hostname>   Server hostname (default: 127.0.0.1)
  -p <port>       Server port (default: 6379)
  -a <password>   Password to use when connecting to the server
  -n <database>   Database number
  -q              Default redis query condition (default: *)
  --help          Output this help and exit
  --version       Output version and exit
```

## LICENSE

[MIT](https://github.com/kyai/redis-cui/blob/master/LICENSE)
