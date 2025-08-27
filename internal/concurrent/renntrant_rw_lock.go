package concurrent

import (
	"sync"

	"github.com/dsx137/gg-kit/internal/lang"
)

type ReentrantRWLock struct {
	mu         *sync.Mutex
	rCond      *sync.Cond
	rOwners    map[int]int
	wCond      *sync.Cond
	wOwner     int
	wHoldCount int
}

func NewReentrantRWLock() *ReentrantRWLock {
	mu := &sync.Mutex{}
	return &ReentrantRWLock{
		mu:         mu,
		rCond:      sync.NewCond(mu),
		rOwners:    make(map[int]int),
		wCond:      sync.NewCond(mu),
		wOwner:     0,
		wHoldCount: 0,
	}
}

func (rw *ReentrantRWLock) RLock() {
	me := lang.GetGoroutineId()

	rw.mu.Lock()
	defer rw.mu.Unlock()

	for rw.wHoldCount > 0 && rw.wOwner != me {
		rw.rCond.Wait()
	}

	rw.rOwners[me]++
}

func (rw *ReentrantRWLock) RUnlock() {
	me := lang.GetGoroutineId()

	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.rOwners[me] == 0 {
		panic("unlock of unlocked lock")
	}

	rw.rOwners[me]--
	if rw.rOwners[me] == 0 {
		delete(rw.rOwners, me)
	}
	if len(rw.rOwners) == 0 && rw.wHoldCount == 0 {
		rw.wCond.Signal()
	}
}

func (rw *ReentrantRWLock) Lock() {
	me := lang.GetGoroutineId()

	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.wOwner == me {
		rw.wHoldCount++
		return
	}

	for (rw.wHoldCount > 0 && rw.wOwner != me) || (len(rw.rOwners) > 0 && !(len(rw.rOwners) == 1 && rw.rOwners[me] > 0)) {
		rw.wCond.Wait()
	}

	rw.wOwner = me
	rw.wHoldCount = 1
}

func (rw *ReentrantRWLock) Unlock() {
	me := lang.GetGoroutineId()

	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.wOwner != me {
		panic("unlock of unlocked lock")
	}

	rw.wHoldCount--
	if rw.wHoldCount == 0 {
		rw.wOwner = 0
		rw.wCond.Signal()
		rw.rCond.Broadcast()
	}
}
