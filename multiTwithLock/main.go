package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"
)
func cldemote(addr unsafe.Pointer)
func shared_cacheline_demote(start unsafe.Pointer, size uintptr){
	//cache_line_base := start & ~0x3F;
	//fmt.Println(start)
	cache_line_base := unsafe.Pointer(uintptr(start) &^ 0x3F)
	// fmt.Println("after align: ", cache_line_base)
	// fmt.Println("size: ", size)
	// fmt.Println("limit: ", unsafe.Add(start, size))

	i := 1
    for {
		cache_line_base = unsafe.Pointer(uintptr(cache_line_base) &^ 0x3F)

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
//var sharedData [1024]int
var mutex   sync.Mutex
var shareint int

func main() {
	cldemote_enabled := flag.Bool("enable_cldemote", false, "a bool")
	flag.Parse()
	fmt.Println(*cldemote_enabled)
	
	nThreads := runtime.GOMAXPROCS(0)
	fmt.Println(nThreads)
	
	var wg sync.WaitGroup
	wg.Add(nThreads)

	for i := 0; i < nThreads; i++ {
		go func() {
			for start:=time.Now(); time.Since(start) < 10*time.Second; {
				mutex.Lock()
				shareint++
				//fmt.Println(shareint)
				
				if *cldemote_enabled {
					//fmt.Println("demote")
					shared_cacheline_demote(unsafe.Pointer(&shareint), unsafe.Sizeof(shareint))
				}
				mutex.Unlock()
				// osyield
				
			}
			wg.Done()
		}()
	}
	
	wg.Wait()

	fmt.Println(shareint)
	//fmt.Println(sharedData)
    //fmt.Println("done")

}

