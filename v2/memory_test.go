package goerr

import (
	"runtime"
	"testing"
	"time"
)

// TestNoMemoryLeak 测试无内存泄漏
func TestNoMemoryLeak(t *testing.T) {
	// 强制 GC
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 创建大量错误对象
	const iterations = 10000
	for i := 0; i < iterations; i++ {
		err := New("LEAK_TEST", "test error").
			WithDetail("iteration", i).
			WithDetail("timestamp", time.Now().Unix())

		// 确保错误被使用
		_ = err.Error()
		_ = err.Code()
		_ = err.Details()

		// 模拟错误链
		if i%10 == 0 {
			err = Wrap(err, "WRAPPED", "wrapped error")
		}
	}

	// 强制 GC
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	// 计算内存增长
	allocDiff := m2.Alloc - m1.Alloc

	// 理论上应该接近 0,因为所有对象都应该被 GC
	// 允许一些合理的增长(如 1MB)
	maxAllowed := uint64(1 * 1024 * 1024)

	if allocDiff > maxAllowed {
		t.Errorf("Memory leak detected: allocated %d bytes (max allowed: %d)",
			allocDiff, maxAllowed)
	}

	t.Logf("Memory stats: before=%d, after=%d, diff=%d", m1.Alloc, m2.Alloc, allocDiff)
}

// TestCodeRegistryNoLeak 测试错误码注册表无内存泄漏
func TestCodeRegistryNoLeak(t *testing.T) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 大量读取错误码配置
	const iterations = 100000
	for i := 0; i < iterations; i++ {
		_ = GetCodeConfig(ErrNotFound)
		_ = GetCodeConfig(ErrInvalidArgument)
		_ = GetCodeConfig(ErrInternal)
	}

	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	// GetCodeConfig 应该是零分配的
	allocDiff := m2.Alloc - m1.Alloc
	maxAllowed := uint64(100 * 1024) // 允许 100KB

	if allocDiff > maxAllowed {
		t.Errorf("Code registry leak: allocated %d bytes", allocDiff)
	}

	t.Logf("Code registry memory: before=%d, after=%d, diff=%d",
		m1.Alloc, m2.Alloc, allocDiff)
}

// TestDetailsNoLeak 测试详细信息无内存泄漏
func TestDetailsNoLeak(t *testing.T) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 创建基础错误
	baseErr := New("BASE", "base error")

	// 大量添加详细信息
	const iterations = 10000
	for i := 0; i < iterations; i++ {
		err := baseErr.WithDetail("key", i)
		_ = err.Details()

		// err 应该在循环结束时被 GC
	}

	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	allocDiff := m2.Alloc - m1.Alloc
	maxAllowed := uint64(1 * 1024 * 1024)

	if allocDiff > maxAllowed {
		t.Errorf("Details leak: allocated %d bytes", allocDiff)
	}

	t.Logf("Details memory: before=%d, after=%d, diff=%d",
		m1.Alloc, m2.Alloc, allocDiff)
}

// TestErrorChainNoLeak 测试错误链无内存泄漏
func TestErrorChainNoLeak(t *testing.T) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 创建深层错误链
	const iterations = 1000
	const chainDepth = 10

	for i := 0; i < iterations; i++ {
		var err error = New("ROOT", "root error")

		for j := 0; j < chainDepth; j++ {
			err = Wrap(err, "WRAP", "wrapped error")
		}

		// 使用错误链
		_ = err.Error()

		// err 应该在循环结束时被 GC(包括整个链)
	}

	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	allocDiff := m2.Alloc - m1.Alloc
	maxAllowed := uint64(1 * 1024 * 1024)

	if allocDiff > maxAllowed {
		t.Errorf("Error chain leak: allocated %d bytes", allocDiff)
	}

	t.Logf("Error chain memory: before=%d, after=%d, diff=%d",
		m1.Alloc, m2.Alloc, allocDiff)
}

// TestStackNoLeak 测试堆栈信息无内存泄漏
func TestStackNoLeak(t *testing.T) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 创建大量错误并访问堆栈
	const iterations = 10000
	for i := 0; i < iterations; i++ {
		err := New("STACK", "error with stack")
		stack := err.Stack()

		// 访问堆栈信息
		for _, frame := range stack {
			_ = frame.String()
		}

		// err 和 stack 应该被 GC
	}

	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	allocDiff := m2.Alloc - m1.Alloc
	maxAllowed := uint64(1 * 1024 * 1024)

	if allocDiff > maxAllowed {
		t.Errorf("Stack leak: allocated %d bytes", allocDiff)
	}

	t.Logf("Stack memory: before=%d, after=%d, diff=%d",
		m1.Alloc, m2.Alloc, allocDiff)
}

// TestConcurrentNoLeak 测试并发场景无内存泄漏
func TestConcurrentNoLeak(t *testing.T) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)

	// 并发创建错误
	const goroutines = 100
	const iterations = 100

	done := make(chan bool, goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer func() { done <- true }()

			for i := 0; i < iterations; i++ {
				err := New("CONCURRENT", "concurrent error").
					WithDetail("goroutine", id).
					WithDetail("iteration", i)

				_ = err.Error()
				_ = err.Code()
				_ = err.Details()
			}
		}(g)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < goroutines; i++ {
		<-done
	}

	runtime.GC()
	time.Sleep(100 * time.Millisecond)
	runtime.ReadMemStats(&m2)

	allocDiff := m2.Alloc - m1.Alloc
	maxAllowed := uint64(2 * 1024 * 1024)

	if allocDiff > maxAllowed {
		t.Errorf("Concurrent leak: allocated %d bytes", allocDiff)
	}

	t.Logf("Concurrent memory: before=%d, after=%d, diff=%d",
		m1.Alloc, m2.Alloc, allocDiff)
}

// BenchmarkMemoryAllocation 内存分配基准测试
func BenchmarkMemoryAllocation(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = New("BENCH", "benchmark error")
		}
	})

	b.Run("WithDetail", func(b *testing.B) {
		err := New("BENCH", "benchmark error")
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = err.WithDetail("key", i)
		}
	})

	b.Run("GetCodeConfig", func(b *testing.B) {
		RegisterCode("BENCH_CODE", "benchmark", 500)
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = GetCodeConfig("BENCH_CODE")
		}
	})

	b.Run("Details", func(b *testing.B) {
		err := New("BENCH", "benchmark error").
			WithDetail("key1", "value1").
			WithDetail("key2", "value2")
		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			_ = err.Details()
		}
	})
}

// TestGCPressure 测试 GC 压力
func TestGCPressure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping GC pressure test in short mode")
	}

	// 记录初始 GC 统计
	var stats1, stats2 runtime.MemStats
	runtime.ReadMemStats(&stats1)
	initialGC := stats1.NumGC

	// 创建大量短生命周期的错误对象
	const duration = 2 * time.Second
	start := time.Now()
	count := 0

	for time.Since(start) < duration {
		err := New("GC_TEST", "gc pressure test").
			WithDetail("count", count).
			WithDetail("time", time.Now().Unix())

		_ = err.Error()
		count++

		// 每 1000 次手动触发 GC
		if count%1000 == 0 {
			runtime.GC()
		}
	}

	// 最后一次 GC
	runtime.GC()
	runtime.ReadMemStats(&stats2)

	gcCount := stats2.NumGC - initialGC

	t.Logf("Created %d errors in %v", count, duration)
	t.Logf("GC runs: %d", gcCount)
	t.Logf("Final alloc: %d bytes", stats2.Alloc)
	t.Logf("Total alloc: %d bytes", stats2.TotalAlloc)

	// 检查最终内存使用是否合理
	if stats2.Alloc > 10*1024*1024 {
		t.Errorf("High memory usage after GC: %d bytes", stats2.Alloc)
	}
}
