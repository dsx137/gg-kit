package concurrent

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dsx137/gg-kit/internal/lang"
)

func TestReentrantRWLock_Basic(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 基本读锁测试
	rwlock.RLock()
	lang.Useless(nil)
	rwlock.RUnlock()

	// 基本写锁测试
	rwlock.Lock()
	lang.Useless(nil)
	rwlock.Unlock()
}

func TestReentrantRWLock_ReadReentrancy(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 读锁重入
	rwlock.RLock()
	rwlock.RLock()
	rwlock.RLock()

	lang.Useless(nil)

	// 必须调用相应次数的RUnlock
	rwlock.RUnlock()
	rwlock.RUnlock()
	rwlock.RUnlock()
}

func TestReentrantRWLock_WriteReentrancy(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 写锁重入
	rwlock.Lock()
	rwlock.Lock()
	rwlock.Lock()

	lang.Useless(nil)

	// 必须调用相应次数的Unlock
	rwlock.Unlock()
	rwlock.Unlock()
	rwlock.Unlock()
}

func TestReentrantRWLock_LockUpgrade(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 读锁升级到写锁
	rwlock.RLock()
	rwlock.Lock()

	lang.Useless(nil)
	// 按照设计，需要显式解锁
	rwlock.Unlock()
	rwlock.RUnlock()
}

func TestReentrantRWLock_LockDowngrade(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 写锁降级到读锁
	rwlock.Lock()
	rwlock.RLock()

	lang.Useless(nil)

	// 按照设计，需要显式解锁
	rwlock.RUnlock()
	rwlock.Unlock()
}

func TestReentrantRWLock_MultipleReaders(t *testing.T) {
	rwlock := NewReentrantRWLock()
	var wg sync.WaitGroup
	readCount := int32(0)

	// 启动多个读者
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rwlock.RLock()
			atomic.AddInt32(&readCount, 1)
			time.Sleep(100 * time.Millisecond)
			atomic.AddInt32(&readCount, -1)
			rwlock.RUnlock()
		}()
	}

	// 等待所有读者启动
	time.Sleep(50 * time.Millisecond)

	// 验证多个读者可以同时持有锁
	if atomic.LoadInt32(&readCount) <= 1 {
		t.Error("Multiple readers should be able to hold the lock simultaneously")
	}

	wg.Wait()
}

func TestReentrantRWLock_WriterExclusion(t *testing.T) {
	rwlock := NewReentrantRWLock()
	var wg sync.WaitGroup
	writeCount := int32(0)
	maxConcurrentWrites := int32(0)

	// 启动多个写者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rwlock.Lock()
			current := atomic.AddInt32(&writeCount, 1)

			// 更新最大并发写者数
			for {
				max := atomic.LoadInt32(&maxConcurrentWrites)
				if current <= max || atomic.CompareAndSwapInt32(&maxConcurrentWrites, max, current) {
					break
				}
			}

			time.Sleep(50 * time.Millisecond)
			atomic.AddInt32(&writeCount, -1)
			rwlock.Unlock()
		}()
	}

	wg.Wait()

	// 验证写者互斥
	if atomic.LoadInt32(&maxConcurrentWrites) > 1 {
		t.Errorf("Writers should be mutually exclusive, but found %d concurrent writers", maxConcurrentWrites)
	}
}

func TestReentrantRWLock_ReaderWriterExclusion(t *testing.T) {
	rwlock := NewReentrantRWLock()
	var wg sync.WaitGroup
	readCount := int32(0)
	writeCount := int32(0)
	violations := int32(0)

	// 启动读者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				rwlock.RLock()
				atomic.AddInt32(&readCount, 1)

				// 检查是否有写者同时存在
				if atomic.LoadInt32(&writeCount) > 0 {
					atomic.AddInt32(&violations, 1)
				}

				time.Sleep(10 * time.Millisecond)
				atomic.AddInt32(&readCount, -1)
				rwlock.RUnlock()
				runtime.Gosched()
			}
		}()
	}

	// 启动写者
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				rwlock.Lock()
				atomic.AddInt32(&writeCount, 1)

				// 检查是否有读者同时存在
				if atomic.LoadInt32(&readCount) > 0 {
					atomic.AddInt32(&violations, 1)
				}

				time.Sleep(20 * time.Millisecond)
				atomic.AddInt32(&writeCount, -1)
				rwlock.Unlock()
				runtime.Gosched()
			}
		}()
	}

	wg.Wait()

	// 验证读写互斥
	if atomic.LoadInt32(&violations) > 0 {
		t.Errorf("Found %d reader-writer exclusion violations", violations)
	}
}

func TestReentrantRWLock_ComplexScenario(t *testing.T) {
	rwlock := NewReentrantRWLock()
	var wg sync.WaitGroup

	// 复杂场景：读锁重入 + 锁升级
	wg.Add(1)
	go func() {
		defer wg.Done()

		// 获取读锁
		rwlock.RLock()
		rwlock.RLock() // 读锁重入

		// 升级到写锁
		rwlock.Lock()

		// 写锁重入
		rwlock.Lock()

		// 降级到读锁
		rwlock.RLock()

		lang.Useless(nil)

		// 按顺序释放
		rwlock.RUnlock() // 释放降级的读锁
		rwlock.Unlock()  // 释放重入的写锁
		rwlock.Unlock()  // 释放原始写锁
		rwlock.RUnlock() // 释放重入的读锁
		rwlock.RUnlock() // 释放原始读锁
	}()

	wg.Wait()
}

func TestReentrantRWLock_PanicOnInvalidUnlock(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 测试解锁未持有的读锁
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when unlocking unheld read lock")
		}
	}()
	rwlock.RUnlock()
}

func TestReentrantRWLock_PanicOnInvalidWriteUnlock(t *testing.T) {
	rwlock := NewReentrantRWLock()

	// 测试解锁未持有的写锁
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when unlocking unheld write lock")
		}
	}()
	rwlock.Unlock()
}

func TestReentrantRWLock_Fairness(t *testing.T) {
	rwlock := NewReentrantRWLock()
	var wg sync.WaitGroup

	results := make([]int, 0, 100)
	var resultMu sync.Mutex

	// 启动一个写者
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(50 * time.Millisecond) // 让读者先启动

		rwlock.Lock()
		resultMu.Lock()
		results = append(results, -1) // -1 表示写者
		resultMu.Unlock()
		time.Sleep(10 * time.Millisecond)
		rwlock.Unlock()
	}()

	// 启动多个读者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		readerID := i
		go func() {
			defer wg.Done()

			rwlock.RLock()
			resultMu.Lock()
			results = append(results, readerID)
			resultMu.Unlock()
			time.Sleep(100 * time.Millisecond)
			rwlock.RUnlock()
		}()
	}

	wg.Wait()

	// 基本验证：确保写者最终能获得锁
	hasWriter := false
	for _, result := range results {
		if result == -1 {
			hasWriter = true
			break
		}
	}

	if !hasWriter {
		t.Error("Writer should eventually acquire the lock")
	}
}

// 压力测试
func TestReentrantRWLock_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	rwlock := NewReentrantRWLock()
	var wg sync.WaitGroup
	counter := int64(0)

	// 多个读者
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				rwlock.RLock()
				_ = atomic.LoadInt64(&counter)
				rwlock.RUnlock()
				runtime.Gosched()
			}
		}()
	}

	// 多个写者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				rwlock.Lock()
				atomic.AddInt64(&counter, 1)
				rwlock.Unlock()
				runtime.Gosched()
			}
		}()
	}

	wg.Wait()

	expectedCounter := int64(5 * 200)
	if atomic.LoadInt64(&counter) != expectedCounter {
		t.Errorf("Expected counter to be %d, got %d", expectedCounter, counter)
	}
}

// 基准测试
func BenchmarkReentrantRWLock_ReadOnly(b *testing.B) {
	rwlock := NewReentrantRWLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwlock.RLock()
			lang.Useless(nil)
			rwlock.RUnlock()
		}
	})
}

func BenchmarkReentrantRWLock_WriteOnly(b *testing.B) {
	rwlock := NewReentrantRWLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwlock.Lock()
			lang.Useless(nil)
			rwlock.Unlock()
		}
	})
}

func BenchmarkReentrantRWLock_ReadHeavy(b *testing.B) {
	rwlock := NewReentrantRWLock()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if b.N%100 == 0 {
				rwlock.Lock()
				lang.Useless(nil)
				rwlock.Unlock()
			} else {
				rwlock.RLock()
				lang.Useless(nil)
				rwlock.RUnlock()
			}
		}
	})
}
