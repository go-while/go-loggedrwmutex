package loggedrwmutex

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLoggedSyncRWMutex(t *testing.T) {
	mux := &LoggedSyncRWMutex{Name: "TestMutex"}

	// Test Lock/Unlock
	mux.Lock()
	if mux.lockedCount != 1 {
		t.Errorf("lockedCount should be 1 after Lock, got %d", mux.lockedCount)
	}
	mux.Unlock()
	if mux.lockedCount != 0 {
		t.Errorf("lockedCount should be 0 after Unlock, got %d", mux.lockedCount)
	}

	// Test RLock/RUnlock
	mux.RLock()
	if mux.rLockedCount != 1 {
		t.Errorf("rLockedCount should be 1 after RLock, got %d", mux.rLockedCount)
	}
	mux.RUnlock()
	if mux.rLockedCount != 0 {
		t.Errorf("rLockedCount should be 0 after RUnlock, got %d", mux.rLockedCount)
	}

	// Test concurrent RLock
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mux.RLock()
			time.Sleep(10 * time.Millisecond)
			mux.RUnlock()
		}()
	}
	wg.Wait()

	// Test Lock and RLock interaction
	var lockHeld bool
	wg.Add(1)
	go func() {
		defer wg.Done()
		mux.Lock()
		lockHeld = true
		time.Sleep(50 * time.Millisecond)
		mux.Unlock()
		lockHeld = false
	}()

	time.Sleep(10 * time.Millisecond) // Give the lock goroutine a chance to start

	wg.Add(1)
	go func() {
		defer wg.Done()
		mux.RLock() // This should block until the Lock is released
		if lockHeld {
			t.Error("RLock acquired while Lock was held")
		}
		mux.RUnlock()
	}()

	wg.Wait()

	// Test PrintStatus (without capturing output) - just ensure it doesn't panic
	mux.PrintStatus(true)

	// Test total counts
	if mux.totalLocked != 2 {
		t.Errorf("totalLocked should be 2, got %d", mux.totalLocked)
	}
	if mux.totalUnlocked != 2 {
		t.Errorf("totalUnlocked should be 2, got %d", mux.totalUnlocked)
	}
	if mux.totalrLocked != 5 {
		t.Errorf("totalrLocked should be 5, got %d", mux.totalrLocked)
	}
	if mux.totalrUnlocked != 5 {
		t.Errorf("totalrUnlocked should be 5, got %d", mux.totalrUnlocked)
	}

	// Test internal mutex protection
	mux.mu.Lock()
	mux.lockedCount = 100 // Simulate race condition
	mux.mu.Unlock()

	mux.Lock()
	fmt.Println("Lock/Unlock")
	mux.Unlock()

	mux.RLock()
	fmt.Println("RLock/Unlock")
	mux.RUnlock()

	fmt.Println("quit")
}
