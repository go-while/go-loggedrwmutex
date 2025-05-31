package loggedrwmutex

import (
	"log"
	"sync"
)

// LoggedSyncRWMutex is a mutex that logs its actions.
// It wraps sync.Mutex and sync.RWMutex to provide logging for lock and unlock actions.
// This is useful for debugging and tracking mutex usage in concurrent applications.
// It can be used in place of sync.RWMutex.
// Usage:
// import : "github.com/go-while/go-loggedrwmutex"
//
//	    var mux *loggedrwmutex.LoggedSyncRWMutex
//		mux := &loggedrwmutexLoggedSyncRWMutex{Name: "XXYYZZ" }'
//		item.mux = mux
//		item.mux.Lock()           // locks the mutex
//		item.mux.Unlock()         // unlocks the mutex
//		item.mux.RLock()          // acquires a read lock
//		item.mux.RUnlock()        // releases a read lock
//		locked, rlocked := item.mux.Status() // checks the status of the mutex
type LoggedSyncRWMutex struct {
	mu             sync.RWMutex // internal mutex to protect the state of the LoggedSyncRWMutex
	Name           string
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
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.lockedCount > 0 || m.rLockedCount > 0 || forceprint {
		log.Printf("?? [loggedMUTEX] Status '%s' locked=%d, rLocked=%d totalLocked/totalUnlocked=%d/%d totalrLocked/totalrUnlocked=%d/%d", m.Name, m.lockedCount, m.rLockedCount, m.totalLocked, m.totalUnlocked, m.totalrLocked, m.totalrUnlocked)
	}
	return
}

func (m *LoggedSyncRWMutex) Lock() {
	m.mu.Lock()
	m.lockedCount++
	m.totalLocked++
	m.mu.Unlock()

	m.RWMutex.Lock()
}

func (m *LoggedSyncRWMutex) Unlock() {
	m.RWMutex.Unlock()

	m.mu.Lock()
	m.lockedCount--
	m.totalUnlocked++
	m.mu.Unlock()
}

func (m *LoggedSyncRWMutex) RLock() {
	m.mu.Lock()
	m.rLockedCount++
	m.totalrLocked++
	m.mu.Unlock()

	m.RWMutex.RLock()
}

func (m *LoggedSyncRWMutex) RUnlock() {
	m.RWMutex.RUnlock()

	m.mu.Lock()
	m.rLockedCount--
	m.totalrUnlocked++
	m.mu.Unlock()
}
