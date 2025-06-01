package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-while/go-loggedrwmutex"
)

func main() {
	mux := &loggedrwmutex.LoggedSyncRWMutex{Name: "ResourceMutex"}
	mux.DebugAll = true // Enable all debug messages
	tmax := 10000
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		i := 0
		for {
			if i > tmax {
				break
			}
			mux.Lock()
			fmt.Println("Writer: Acquired lock")
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			fmt.Println("Writer: Releasing lock")
			mux.Unlock()
			i++
		}
	}()

	go func() {
		defer wg.Done()
		i := 0
		for {
			if i > tmax {
				break
			}
			mux.RLock()
			fmt.Println("Reader 1: Acquired read lock")
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			fmt.Println("Reader 1: Releasing read lock")
			mux.RUnlock()
			i++
		}
	}()

	go func() {
		defer wg.Done()
		i := 0
		for {
			if i > tmax {
				break
			}
			mux.RLock()
			fmt.Println("Reader 2: Acquired read lock")
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			fmt.Println("Reader 2: Releasing read lock")
			mux.RUnlock()
			i++
		}
	}()

	go func() {
		defer wg.Done()
		i := 0
		for {
			if i > tmax {
				break
			}
			mux.PrintStatus(true)
			i++
		}
	}()
	wg.Wait()
}
