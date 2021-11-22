Pubsub
==============

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/zhong-my/pubsub/master/LICENSE)

# What is Pubsub

Pubsub is prove of concept for Redis "Pub/Sub" messaging management feature. SUBSCRIBE, UNSUBSCRIBE and PUBLISH implement the [Publish/Subscribe messaging paradigm](http://en.wikipedia.org/wiki/Publish/subscribe) where (citing Wikipedia) senders (publishers) are not programmed to send their messages to specific receivers (subscribeds). Rather, published messages are characterized into channels, without knowledge of what (if any) subscriber there may be. Subscribeds express interest in one or more channels, and only receive messages that are of interest, without knowledge of what (if any) publishers there are. This decoupling of publishers and subscribers can allow for greater scalability and a more dynamic network toplogy.

## Install

```
    go get github.com/zhong-my/pubsub
```

## Usage

```golang
package main

import (
    "fmt"
    . "github.com/kkdai/pubsub"
)

func main() {
    ser := NewPubsub(1)
    c1 := ser.Subscribe("A")
    c2 := ser.Subscribe("B")
    ser.Publish("test1", "A")
    ser.Publish("test2", "B")
    fmt.Println(<-c1) // "test1"
    fmt.Println(<-c2) // "test2"

    // Add subscription "B" for c1.          
    ser.AddSubscription(c1, "B")

    // Publish new content in B
    ser.Publish("test3", "B")

    fmt.Println(<-c1) // "test3"

    // Remove subscription "B" in c1
    ser.RemoveSubscription(c1, "B")
        	
    // Publish new content in B
    ser.Publish("test4", "B")
        
    select {
    case val := <-c1:
    fmt.Printf("Should not get %v notify on remove topic\n", val)
        break
    case <-time.After(time.Second):
        // Will go here, because we remove subscription B in c1.         
        break
    }
}
```

## Benchmark

Benchmark include menmory usage.

Run all benchmark for test command:

```
    go test -bench=.
```

command output:

```
    goos: darwin
    goarch: amd64
    pkg: pubsub
    cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz

    BenchmarkAddSub-4                   3366            749178 ns/op
    BenchmarkRemoveSub-4               10000            190355 ns/op
    BenchmarkBasicFunction-4        21844840            54.95 ns/op
```

## Inspired By

- [Redis: Pubsub](http://redis.io/topics/pubsub)
- [kkdai/pubsub](https://github.com/kkdai/pubsub)

## License

This package is licensed under MIT license. See LICENSE for details.