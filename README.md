# go-loggedrwmutex

`go-loggedrwmutex` is a Go package that provides a logging wrapper around `sync.RWMutex`. It's designed to help debug and track mutex usage in concurrent applications by logging lock and unlock actions.

## Features

- Wraps `sync.RWMutex` to provide logging.
- Logs lock and unlock actions for debugging.
- Can be used as a direct replacement for `sync.RWMutex`.
- Provides a `PrintStatus` method to check the current status of the mutex (locked or read-locked) and print detailed statistics.
- Individual control over debug logging for `Lock`, `Unlock`, `RLock`, and `RUnlock` operations.
- Global debug flag to enable/disable logging for all mutexes.

## Installation

```bash
go get github.com/go-while/go-loggedrwmutex
```

## Usage

```go
import "github.com/go-while/go-loggedrwmutex"

func init() {
	// Any Debug flags can only be set on boot time
	//  before initializing any mutexes!
	// or you may get data race errors.
	// this is a non fix as this package is only for debugging purposes.

	// Enable global debug messages for all mutexes
	loggedrwmutex.GlobalDebug = true

	// global flag to disable logging
	// and bypass directly to original mutexes without counting
	loggedrwmutex.DisableLogging = false
}

func main() {
    // Initialize the logged RWMutex with a name
    mux := &loggedrwmutex.LoggedSyncRWMutex{Name: "MyMutex"}

    // Enable specific debug flags
    mux.DebugAll = true       // Enables all debug messages
    mux.DebugLock = true      // Enables debug messages for Lock
    mux.DebugUnlock = true    // Enables debug messages for Unlock
    mux.DebugRLock = true     // Enables debug messages for RLock
    mux.DebugRUnlock = true   // Enables debug messages for RUnlock

    // Lock and unlock for exclusive access
    mux.Lock()
    // Critical section
    mux.Unlock()

    // RLock and RUnlock for read access
    mux.RLock()
    // Read-only section
    mux.RUnlock()

    // Print the status of the mutex
    mux.PrintStatus(false) // Set to true to force print even if not locked
}
```

## API Reference

### `LoggedSyncRWMutex`

A struct that wraps `sync.RWMutex` and adds logging capabilities.

```go
type LoggedSyncRWMutex struct {
	mu             sync.RWMutex // internal mutex to protect the state of the LoggedSyncRWMutex
	Name           string
	DebugAll       bool   // if true, will print debug messages
	DebugLock      bool   // if true, will print debug messages
	DebugUnlock    bool   // if true, will print debug messages
	DebugRLock     bool   // if true, will print debug messages
	DebugRUnlock   bool   // if true, will print debug messages
	lockedCount    uint64 // number of active locks
	rLockedCount   uint64 // number of active readers
	totalLocked    uint64
	totalUnlocked  uint64
	totalrLocked   uint64
	totalrUnlocked uint64
	sync.RWMutex   // the actual mutex that will be used for locking
}
```

- `Name`: A descriptive name for the mutex, useful for identifying it in logs.
- `DebugAll`: If `true`, enables debug messages for all lock/unlock actions.
- `DebugLock`: If `true`, enables debug messages for `Lock` actions.
- `DebugUnlock`: If `true`, enables debug messages for `Unlock` actions.
- `DebugRLock`: If `true`, enables debug messages for `RLock` actions.
- `DebugRUnlock`: If `true`, enables debug messages for `RUnlock` actions.
- `lockedCount`: Number of active write locks.
- `rLockedCount`: Number of active read locks.
- `totalLocked`: Total number of times the mutex has been locked.
- `totalUnlocked`: Total number of times the mutex has been unlocked.
- `totalrLocked`: Total number of times the mutex has been read-locked.
- `totalrUnlocked`: Total number of times the mutex has been read-unlocked.
- `sync.RWMutex`: The embedded `sync.RWMutex` that performs the actual locking.

### Methods

#### `Lock()`

Locks the mutex for exclusive access.  Increments internal counters to track locking statistics.

#### `Unlock()`

Unlocks the mutex, releasing exclusive access.  Increments internal counters to track unlocking statistics.

#### `RLock()`

Acquires a read lock, allowing concurrent read access. Increments internal read lock counters.

#### `RUnlock()`

Releases a read lock. Increments internal read unlock counters.

#### `PrintStatus(forceprint bool) (locked bool, rlocked bool)`

Prints the current status of the mutex, including lock counts and total lock/unlock statistics.

- `forceprint`: If `true`, forces the status to be printed even if the mutex is not currently locked.
- Returns:
    - `locked`: `true` if the mutex is currently locked for exclusive access, `false` otherwise. (Always false, as the status is checked under a separate mutex)
    - `rlocked`: `true` if the mutex is currently read-locked, `false` otherwise. (Always false, as the status is checked under a separate mutex)

## Example

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-while/go-loggedrwmutex"
)

func main() {
	mux := &loggedrwmutex.LoggedSyncRWMutex{
		Name:     "ResourceMutex",
		DebugAll: true,
	}
	tmax := 10000
	var wg sync.WaitGroup
	wg.Add(4)
	go func() {
		defer wg.Done()
		i := 0
		for {
			if i >= tmax {
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
			if i >= tmax {
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
			if i >= tmax {
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
			if i >= tmax {
				break
			}
			mux.Lock()
			time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			mux.Unlock()

			mux.PrintStatus(true)
			i++
		}
	}()

	wg.Wait() // Wait for all goroutines to finish
	mux.PrintStatus(true)
}


```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues to discuss potential improvements or report bugs.

## License

This project is licensed under the [MIT License](LICENSE).