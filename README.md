# go-loggedrwmutex

`go-loggedrwmutex` is a Go package that provides a logging wrapper around `sync.RWMutex`. It's designed to help debug and track mutex usage in concurrent applications by logging lock and unlock actions.

## Features

- Wraps `sync.RWMutex` to provide logging.
- Logs lock and unlock actions for debugging.
- Can be used as a direct replacement for `sync.RWMutex`.
- Provides a `PrintStatus` method to check the current status of the mutex (locked or read-locked) and print detailed statistics.

## Installation

```bash
go get github.com/go-while/go-loggedrwmutex
```

## Usage

```go
import "github.com/go-while/go-loggedrwmutex"

func main() {
    // Initialize the logged RWMutex with a name
    mux := &loggedrwmutex.loggedSyncRWMutex{Name: "MyMutex"}

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

### `loggedSyncRWMutex`

A struct that wraps `sync.RWMutex` and adds logging capabilities.

```go
type loggedSyncRWMutex struct {
	mu             sync.RWMutex // internal mutex to protect the state of the loggedSyncRWMutex
	Name           string
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
    - `locked`: `true` if the mutex is currently locked for exclusive access, `false` otherwise.
    - `rlocked`: `true` if the mutex is currently read-locked, `false` otherwise.

## Example

```go
package main

import (
	"fmt"
	"github.com/go-while/go-loggedrwmutex"
	"time"
)

func main() {
	mux := &loggedrwmutex.loggedSyncRWMutex{Name: "ResourceMutex"}

	go func() {
		mux.Lock()
		fmt.Println("Writer: Acquired lock")
		time.Sleep(2 * time.Second)
		fmt.Println("Writer: Releasing lock")
		mux.Unlock()
	}()

	go func() {
		mux.RLock()
		fmt.Println("Reader 1: Acquired read lock")
		time.Sleep(1 * time.Second)
		fmt.Println("Reader 1: Releasing read lock")
		mux.RUnlock()
	}()

	go func() {
		mux.RLock()
		fmt.Println("Reader 2: Acquired read lock")
		time.Sleep(1 * time.Second)
		fmt.Println("Reader 2: Releasing read lock")
		mux.RUnlock()
	}()

	time.Sleep(3 * time.Second)
	mux.PrintStatus(true)
}
```

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues to discuss potential improvements or report bugs.

## License

This project is licensed under the [MIT License](LICENSE).