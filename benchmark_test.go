package goerr

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"testing"
)

// 基准显式丢弃被测结果；业务语义由契约测试验证，b.Loop 保活调用结果。
func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = New(io.EOF, StatusParams(), "invalid input")
	}
}

func BenchmarkWrap(b *testing.B) {
	inner := New(io.EOF, StatusParams(), "invalid input")
	b.ReportAllocs()
	for b.Loop() {
		_ = Wrap(inner, "load")
	}
}

func BenchmarkWrapf(b *testing.B) {
	inner := New(io.EOF, StatusParams(), "invalid input")
	b.ReportAllocs()
	for b.Loop() {
		_ = Wrapf(inner, "step %d", 1)
	}
}

func BenchmarkWrapStatus(b *testing.B) {
	inner := New(io.EOF, StatusParams(), "invalid input")
	b.ReportAllocs()
	for b.Loop() {
		_ = WrapStatus(inner, StatusNotFound())
	}
}

func BenchmarkConvertNil(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Convert(nil)
	}
}

func BenchmarkConvertPlain(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Convert(io.EOF)
	}
}

func BenchmarkConvertItem(b *testing.B) {
	inner := New(io.EOF, StatusParams(), "invalid input")
	b.ReportAllocs()
	for b.Loop() {
		_ = Convert(inner)
	}
}

func BenchmarkConvertDepth(b *testing.B) {
	for _, depth := range []int{1, 4, 16, 64} {
		b.Run(strconv.Itoa(depth), func(b *testing.B) {
			var input error = New(io.EOF, StatusParams(), "invalid input")
			for range depth {
				input = fmt.Errorf("outer: %w", input)
			}
			b.ReportAllocs()
			for b.Loop() {
				_ = Convert(input)
			}
		})
	}
}

func BenchmarkConvertParallel(b *testing.B) {
	inner := New(io.EOF, StatusParams(), "invalid input")
	input := fmt.Errorf("outer: %w", inner)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if got := Convert(input); got != inner {
				b.Error("Convert lost the original Item")
			}
		}
	})
}

func BenchmarkHasCode(b *testing.B) {
	input := fmt.Errorf("outer: %w", New(io.EOF, StatusParams(), "invalid input"))
	b.ReportAllocs()
	for b.Loop() {
		_ = HasCode(input, ErrParams)
	}
}

func BenchmarkWithStackPlain(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = WithStack(io.EOF)
	}
}

func BenchmarkWithStackExisting(b *testing.B) {
	input := fmt.Errorf("outer: %w", New(io.EOF, StatusParams(), "invalid input"))
	b.ReportAllocs()
	for b.Loop() {
		_ = WithStack(input)
	}
}

// 缺栈路径当前没有捕获成本；后续修复必须与正常旧路径分别比较。
func BenchmarkWithStackMissing(b *testing.B) {
	input := Convert(io.EOF)
	b.ReportAllocs()
	for b.Loop() {
		_ = WithStack(input)
	}
}

func BenchmarkFormatVerbose(b *testing.B) {
	input := New(io.EOF, StatusParams(), "invalid input")
	b.ReportAllocs()
	for b.Loop() {
		_ = fmt.Sprintf("%+v", input)
	}
}

func BenchmarkInspect(b *testing.B) {
	origin := New(io.EOF, StatusParams(), "invalid input")
	type inspectInput struct {
		name  string
		input error
	}
	inputs := []inspectInput{
		{name: "nil"},
		{name: "plain", input: io.EOF},
		{name: "item", input: origin},
		{name: "no_stack", input: Convert(io.EOF)},
		{name: "branches", input: errors.Join(Wrap(origin, "left"), Wrap(origin, "right"))},
	}
	for _, depth := range []int{1, 4, 16, 64} {
		var input error = origin
		for range depth {
			input = fmt.Errorf("outer: %w", input)
		}
		inputs = append(inputs, inspectInput{name: "depth_" + strconv.Itoa(depth), input: input})
	}
	limit := io.EOF
	for range 255 {
		limit = fmt.Errorf("x: %w", limit)
	}
	inputs = append(inputs, inspectInput{name: "visit_limit", input: limit})
	for _, input := range inputs {
		b.Run(input.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := Inspect(input.input); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
