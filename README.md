# Memcache Client with Consistent Hashing

## Overview
This is memcache client library for Golang, with consistent hashing implemented. We don't try hacking the memcache server code with Chord, since we want it to stay original.

## Installing
```bash
$ go get github.com/liyinhgqw/memcache
```

## Example
```go
import (
    "fmt"
    "github.com/liyinhgqw/memcache"
)

func main() {
    // register server addresses and weights
    cl := memcache.NewChClient([]string{"localhost:11211", "localhost:11212"}, []uint32{1, 1})
    cl.Set(&memcache.Item{
        Key:   "a",
        Value: []byte("AAA"),
    })
    cl.Set(&memcache.Item{
        Key:   "b",
        Value: []byte("BBB"),
    })
    items, err := cl.GetMulti([]string{"a", "b"})
    if err != nil {
        fmt.Println(err)
    }
    for k, v := range items {
        fmt.Println(k, v)
    }
}
```

You can see `memcache_test.go` for more usage.

## Limitation
Since we don't mess with the memcache server code using Chord algorithm, there are some limitations:

* All clients **Must** use the same server list and weights.
* Add or delete a server is okay, cache miss may happen for some keys.
* Never ever delete an added server, otherwise inconsistency may happen.

## Reference
* [simple hash memcache client](https://github.com/bradfitz/gomemcache)
* [APE memcache client usage](https://developers.google.com/appengine/docs/go/memcache/reference?csw=1)
* [command protocol](http://blog.elijaa.org/?post/2010/05/21/Memcached-telnet-command-summary)