package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

func cldemote(addr unsafe.Pointer)
func shared_cacheline_demote(start unsafe.Pointer, size uintptr) {
	//cache_line_base := start & ~0x3F;
	//fmt.Println(start)
	cache_line_base := unsafe.Pointer(uintptr(start) &^ 0x3F)
	//fmt.Println("after align: ", cache_line_base)
	//fmt.Println("size: ", size)
	//fmt.Println("limit: ", unsafe.Add(start, size))

	i := 1
	for {
		//cache_line_base = unsafe.Pointer(uintptr(cache_line_base) &^ 0x3F)

		cldemote(cache_line_base)
		// next cacheline start size
		cache_line_base = unsafe.Add(cache_line_base, 64)
		//cache_line_base = *(*int)(unsafe.Pointer(uintptr(cache_line_base) + 64))
		//fmt.Println("new cache line base: ", cache_line_base)
		//fmt.Println("demote: ", i)
		i++
		if uintptr(cache_line_base) > uintptr(unsafe.Add(start, size)) {
			break
		}
	}
}

// Shared data
var sharedData [1024]int

func main() {
	cldemote_enabled := flag.Bool("enable_cldemote", false, "a bool")
	flag.Parse()
	//fmt.Println(*cldemote_enabled)

	// Set GOMAXPROCS to the number of available CPUs
	runtime.GOMAXPROCS(56)
	for start := time.Now(); time.Since(start) < 10*time.Second; {
		for i := 0; i < len(sharedData); i++ {
			sharedData[i]++
		}
		//sharedData[1023] = rand.Int()
		if *cldemote_enabled {
			//fmt.Println("demote")
			shared_cacheline_demote(unsafe.Pointer(&sharedData), unsafe.Sizeof(sharedData))
		}

	}
	fmt.Println(sharedData[0])
	// Create a wait group to wait for all goroutines to finish
	// 	var wg sync.WaitGroup

	// 	// Create multiple goroutines to stress different cache lines
	// 	for i := 0; i < 2; i++ {
	// 		wg.Add(1)
	// 		go stressCacheLine(i, &wg)
	// 	}
	// 	// Wait for all goroutines to finish
	// 	wg.Wait()
	// 	fmt.Println("sharedData", sharedData[0])
}
