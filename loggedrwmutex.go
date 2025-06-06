package loggedrwmutex

import (
	"fmt"
	"sync"
)

// Any Debug flags can only be set on boot time / before initializing any mutexes
// or you may get data race errors. this is a non fix as this package is only for debugging purposes.
var GlobalDebug = false    // global debug flag for all mutexes
var DisableLogging = false // global flag to disable logging and bypass directly to original mutexes without counting

// LoggedSyncRWMutex is a mutex that logs its actions.
// It wraps sync.Mutex and sync.RWMutex to provide logging for lock and unlock actions.
// This is useful for debugging and tracking mutex usage in concurrent applications.
// It can be used in place of sync.RWMutex.
// Usage:
// import : "github.com/go-while/go-loggedrwmutex"
//
//		func init() {
//			loggedrwmutex.GlobalDebug = true // enables global debug messages for all mutexes}
//		}
//	    var mux *loggedrwmutex.LoggedSyncRWMutex
//		mux := &loggedrwmutexLoggedSyncRWMutex{Name: "XXYYZZ" }'
//		item.mux = mux
//		item.mux.DebugAll = true // enables all debug messages
//		item.mux.DebugLock = true // enables debug messages for Lock
//		item.mux.DebugUnlock = true // enables debug messages for Unlock
//		item.mux.DebugRLock = true // enables debug messages for RLock
//		item.mux.DebugRUnlock = true // enables debug messages for RUnlock
//		item.mux.Lock()           // locks the mutex
//		item.mux.Unlock()         // unlocks the mutex
//		item.mux.RLock()          // acquires a read lock
//		item.mux.RUnlock()        // releases a read lock
//		locked, rlocked := item.mux.Status(true) // checks the status of the mutex
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

// Status prints the current status of the mutex, including whether it is locked or read-locked.
func (m *LoggedSyncRWMutex) PrintStatus(forceprint bool) (locked bool, rlocked bool) {
	if !DisableLogging {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.lockedCount > 0 || m.rLockedCount > 0 || forceprint {
		fmt.Printf("?? [loggedMUTEX] Status '%s' locked=%d, rLocked=%d totalLocked/totalUnlocked=%d/%d totalrLocked/totalrUnlocked=%d/%d\n", m.Name, m.lockedCount, m.rLockedCount, m.totalLocked, m.totalUnlocked, m.totalrLocked, m.totalrUnlocked)
	}
	return
}

func (m *LoggedSyncRWMutex) Lock() {
	if !DisableLogging {
		m.mu.Lock()
		m.lockedCount++
		m.totalLocked++
		if m.DebugLock || m.DebugAll || GlobalDebug {
			fmt.Printf("[loggedMUTEX] Lock '%s' locked=%d/%d\n", m.Name, m.lockedCount, m.totalLocked)
		}
		m.mu.Unlock()
	}

	m.RWMutex.Lock()
}

func (m *LoggedSyncRWMutex) Unlock() {
	m.RWMutex.Unlock()

	if !DisableLogging {
		m.mu.Lock()
		m.lockedCount--
		m.totalUnlocked++
		if m.DebugUnlock || m.DebugAll || GlobalDebug {
			fmt.Printf("[loggedMUTEX] Unlock '%s' locked=%d/%d\n", m.Name, m.lockedCount, m.totalUnlocked)
		}
		m.mu.Unlock()
	}
}

func (m *LoggedSyncRWMutex) RLock() {
	if !DisableLogging {

		m.mu.Lock()
		m.rLockedCount++
		m.totalrLocked++
		if m.DebugRLock || m.DebugAll || GlobalDebug {
			fmt.Printf("[loggedMUTEX] RLock '%s' rLocked=%d/%d\n", m.Name, m.rLockedCount, m.totalrLocked)
		}
		m.mu.Unlock()
	}
	m.RWMutex.RLock()
}

func (m *LoggedSyncRWMutex) RUnlock() {
	m.RWMutex.RUnlock()

	if !DisableLogging {
		m.mu.Lock()
		m.rLockedCount--
		m.totalrUnlocked++
		if m.DebugRUnlock || m.DebugAll || GlobalDebug {
			fmt.Printf("[loggedMUTEX] RUnlock '%s' rLockedCount=%d/%d\n", m.Name, m.rLockedCount, m.totalrUnlocked)
		}
		m.mu.Unlock()
	}
}
