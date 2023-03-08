package main

import (
	"fmt"
	"labs/lab2/cache"
)

func main() {

	Cache1 := cache.NewCache(5)
	//CacheMap := make(map[string]string)
	//CacheQueue := make([]string, 0)

	//Cache1 := cache.LruCache{Size: 5, Remaining: 5, Cache: CacheMap, Queue: CacheQueue}

	fmt.Println("\nEmpty Cache of Size 5    ---->    ", Cache1)

	Cache1.Put("Key1", "Value1")

	fmt.Println("\nPuttin 1 new data in the cache(Key1)")
	fmt.Println("Current Cache            ---->    ", Cache1)

	Cache1.Put("Key2", "Value2")
	Cache1.Put("Key2", "Value2")
	Cache1.Put("Key2", "Value2")
	Cache1.Put("Key3", "Value3")

	fmt.Println("\nPuttin 4 new data in the cache(Key2, Key3, Key4, Key5)")
	fmt.Println("Current Cache            ---->    ", Cache1)

	Cache1.Put("Key2", "Value2")

	fmt.Println("\nPuttin 1 new data in the cache(Key3)")
	fmt.Println("Current Cache            ---->    ", Cache1)

	fmt.Print("\nGet the Value associated with Key3   ---->    ")
	fmt.Println(Cache1.Get("Key3"))

	fmt.Println("Current Cache            ---->    ", Cache1)

}
