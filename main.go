package main

import (
	"fmt"
	"time"

	"github.com/stas-raranetskyi/cache"
)

func main() {
	cache := cache.New()

	cache.Set("userId", 42, time.Second)
	time.Sleep(time.Second * 3)
	userId := cache.Get("userId")

	fmt.Println("print", userId)

}
