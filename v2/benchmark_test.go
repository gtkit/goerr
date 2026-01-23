package goerr

import (
	"errors"
	"fmt"
	"testing"
)

// BenchmarkNewError 测试创建错误的性能
func BenchmarkNewError(b *testing.B) {
	b.Run("goerr.New", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = New("TEST_ERROR", "test error message")
		}
	})

	b.Run("errors.New", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = errors.New("test error message")
		}
	})

	b.Run("fmt.Errorf", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fmt.Errorf("test error message")
		}
	})
}

// BenchmarkNewWithCode 测试使用错误码创建错误
func BenchmarkNewWithCode(b *testing.B) {
	RegisterCode("BENCH_ERROR", "benchmark error", 400)

	b.Run("goerr.NewWithCode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewWithCode("BENCH_ERROR")
		}
	})

	b.Run("goerr.NewWithCodef", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = NewWithCodef("BENCH_ERROR", "user_id: %d", 12345)
		}
	})
}

// BenchmarkWrapError 测试包装错误的性能
func BenchmarkWrapError(b *testing.B) {
	baseErr := errors.New("base error")

	b.Run("goerr.Wrap", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Wrap(baseErr, "WRAP_ERROR", "wrapped error")
		}
	})

	b.Run("fmt.Errorf_%w", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fmt.Errorf("wrapped error: %w", baseErr)
		}
	})
}

// BenchmarkWithDetails 测试添加详细信息的性能
func BenchmarkWithDetails(b *testing.B) {
	err := New("TEST_ERROR", "test error")

	b.Run("WithDetail", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = err.WithDetail("key", "value")
		}
	})

	b.Run("WithDetails", func(b *testing.B) {
		details := map[string]interface{}{
			"user_id":    12345,
			"request_id": "req-12345",
			"ip":         "192.168.1.1",
		}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = err.WithDetails(details)
		}
	})
}

// BenchmarkErrorIs 测试错误判断的性能
func BenchmarkErrorIs(b *testing.B) {
	err1 := New("TEST_ERROR", "test error")
	err2 := New("TEST_ERROR", "another test error")
	err3 := New("OTHER_ERROR", "other error")

	b.Run("Is_same_code", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Is(err1, err2)
		}
	})

	b.Run("Is_different_code", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Is(err1, err3)
		}
	})
}

// BenchmarkGetCodeConfig 测试获取错误码配置的性能
func BenchmarkGetCodeConfig(b *testing.B) {
	RegisterCode("PERF_TEST", "performance test", 400)

	b.Run("GetCodeConfig", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = GetCodeConfig("PERF_TEST")
		}
	})

	b.Run("GetCodeConfig_notfound", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = GetCodeConfig("NOT_EXISTS")
		}
	})
}

// BenchmarkErrorFormat 测试错误格式化性能
func BenchmarkErrorFormat(b *testing.B) {
	err := New("TEST_ERROR", "test error").
		WithDetail("user_id", 12345).
		WithDetail("request_id", "req-12345")

	b.Run("Error_string", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = err.Error()
		}
	})

	b.Run("Format_v", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%v", err)
		}
	})

	b.Run("Format_+v", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%+v", err)
		}
	})
}

// BenchmarkConcurrentAccess 测试并发访问性能
func BenchmarkConcurrentAccess(b *testing.B) {
	RegisterCode("CONCURRENT_TEST", "concurrent test", 400)

	b.Run("NewWithCode_concurrent", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = NewWithCode("CONCURRENT_TEST")
			}
		})
	})

	b.Run("GetCodeConfig_concurrent", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = GetCodeConfig("CONCURRENT_TEST")
			}
		})
	})

	b.Run("WithDetail_concurrent", func(b *testing.B) {
		err := New("TEST_ERROR", "test error")
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = err.WithDetail("key", "value")
			}
		})
	})
}

// BenchmarkErrorChain 测试错误链性能
func BenchmarkErrorChain(b *testing.B) {
	b.Run("DeepChain", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := errors.New("root error")
			for j := 0; j < 10; j++ {
				err = Wrap(err, "WRAP_ERROR", "wrapped")
			}
			_ = err
		}
	})
}
