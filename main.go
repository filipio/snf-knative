// package main

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/allegro/bigcache/v3"
// )

// func main() {
// 	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

// 	cache.Set("my-unique-key", []byte("value"))

// 	entry, _ := cache.Get("my-unique-key")
// 	fmt.Println(string(entry))
// }