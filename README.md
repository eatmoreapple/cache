## cache

This is a simple cache which provides type constraints and a few other niceties.


### Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/eatmoreapple/cache"
)

func main() {
    c := cache.New[string](10 * time.Minute, time.Minute)

    c.Set("foo", "bar")
    c.Set("baz", "42")

    fmt.Println(c.Get("foo"))
    fmt.Println(c.Get("baz"))
	
	c2 := cache.NewNumericCache[int64](10 * time.Minute, time.Minute)
	c2.Set("foo", 42)
	c2.Increment("foo", 1)
	fmt.Println(c2.Get("foo"))
}
```