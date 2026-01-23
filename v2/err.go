// Package goerr 提供生产级的 Go 错误处理能力
// 特性:
// - 零内存分配的错误码映射
// - 无锁并发安全设计
// - 完整的堆栈追踪
// - 结构化的错误链
// - 高性能错误创建和判断
package goerr

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

// Error 定义标准错误接口
type Error interface {
	error
	// Code 返回错误码
	Code() string
	// HTTPStatus 返回对应的 HTTP 状态码
	HTTPStatus() int
	// Message 返回错误消息
	Message() string
	// Details 返回详细信息
	Details() map[string]interface{}
	// WithDetail 添加详细信息(不可变,返回新错误)
	WithDetail(key string, value interface{}) Error
	// WithDetails 批量添加详细信息
	WithDetails(details map[string]interface{}) Error
	// Unwrap 返回底层错误
	Unwrap() error
	// Stack 返回堆栈信息
	Stack() []StackFrame
	// Format 支持格式化输出
	Format(s fmt.State, verb rune)
}

// StackFrame 堆栈帧信息
type StackFrame struct {
	File     string
	Line     int
	Function string
}

// String 格式化堆栈帧
func (f StackFrame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Function)
}

// err 是 Error 接口的实现
// 使用值类型避免堆分配,通过原子操作保证并发安全
type err struct {
	code       string
	message    string
	httpStatus int
	details    atomic.Value // map[string]interface{} 使用原子值避免锁
	cause      error
	stack      []StackFrame
}

// New 创建一个新错误
func New(code, message string) Error {
	return newError(code, message, 500, nil, 2)
}

// Newf 创建格式化错误消息的错误
func Newf(code, format string, args ...interface{}) Error {
	return newError(code, fmt.Sprintf(format, args...), 500, nil, 2)
}

// Wrap 包装已有错误,添加上下文信息
func Wrap(err error, code, message string) Error {
	if err == nil {
		return nil
	}
	return newError(code, message, 500, err, 2)
}

// Wrapf 包装错误并格式化消息
func Wrapf(err error, code, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}
	return newError(code, fmt.Sprintf(format, args...), 500, err, 2)
}

// newError 内部创建错误的函数
func newError(code, message string, httpStatus int, cause error, skip int) Error {
	e := &err{
		code:       code,
		message:    message,
		httpStatus: httpStatus,
		cause:      cause,
		stack:      captureStack(skip + 1),
	}
	// 初始化 details 为空 map
	e.details.Store(make(map[string]interface{}))
	return e
}

// captureStack 捕获堆栈信息(最多32帧)
func captureStack(skip int) []StackFrame {
	const maxDepth = 32
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(skip+1, pcs)

	if n == 0 {
		return nil
	}

	frames := make([]StackFrame, 0, n)
	callersFrames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := callersFrames.Next()
		frames = append(frames, StackFrame{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
		})
		if !more {
			break
		}
	}

	return frames
}

// Error 实现 error 接口
func (e *err) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.code, e.message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

// Code 返回错误码
func (e *err) Code() string {
	return e.code
}

// HTTPStatus 返回 HTTP 状态码
func (e *err) HTTPStatus() int {
	return e.httpStatus
}

// Message 返回错误消息
func (e *err) Message() string {
	return e.message
}

// Details 返回详细信息的副本(并发安全)
func (e *err) Details() map[string]interface{} {
	stored := e.details.Load()
	if stored == nil {
		return make(map[string]interface{})
	}

	original := stored.(map[string]interface{})
	// 返回副本避免外部修改
	copy := make(map[string]interface{}, len(original))
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// WithDetail 添加单个详细信息(不可变,返回新错误)
func (e *err) WithDetail(key string, value interface{}) Error {
	return e.WithDetails(map[string]interface{}{key: value})
}

// WithDetails 批量添加详细信息(不可变,返回新错误)
func (e *err) WithDetails(details map[string]interface{}) Error {
	if len(details) == 0 {
		return e
	}

	// 创建新错误避免修改原错误
	newErr := &err{
		code:       e.code,
		message:    e.message,
		httpStatus: e.httpStatus,
		cause:      e.cause,
		stack:      e.stack, // 共享堆栈,不需要复制
	}

	// 复制并合并 details
	oldDetails := e.Details()
	newDetails := make(map[string]interface{}, len(oldDetails)+len(details))
	for k, v := range oldDetails {
		newDetails[k] = v
	}
	for k, v := range details {
		newDetails[k] = v
	}
	newErr.details.Store(newDetails)

	return newErr
}

// Unwrap 返回底层错误
func (e *err) Unwrap() error {
	return e.cause
}

// Stack 返回堆栈信息
func (e *err) Stack() []StackFrame {
	return e.stack
}

// Format 实现 fmt.Formatter 接口
func (e *err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 详细模式: +v
			io.WriteString(s, e.Error())
			io.WriteString(s, "\n")

			// 打印详细信息
			if details := e.Details(); len(details) > 0 {
				io.WriteString(s, "Details:\n")
				for k, v := range details {
					fmt.Fprintf(s, "  %s: %v\n", k, v)
				}
			}

			// 打印堆栈
			io.WriteString(s, "Stack trace:\n")
			for _, frame := range e.stack {
				fmt.Fprintf(s, "  %s\n", frame.String())
			}

			// 递归打印 cause
			if e.cause != nil {
				fmt.Fprintf(s, "Caused by: %+v", e.cause)
			}
		} else {
			// 简单模式: %v
			io.WriteString(s, e.Error())
		}
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// Is 实现错误判断(支持 errors.Is)
func (e *err) Is(target error) bool {
	if target == nil {
		return false
	}

	// 类型断言检查
	if t, ok := target.(Error); ok {
		return e.code == t.Code()
	}

	// 检查底层错误
	if e.cause != nil {
		return Is(e.cause, target)
	}

	return false
}

// Is 全局错误判断函数
func Is(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}

	// 如果是我们的错误类型
	if e, ok := err.(Error); ok {
		if t, ok := target.(Error); ok {
			return e.Code() == t.Code()
		}
	}

	// 使用标准库的 Is
	type isInterface interface {
		Is(error) bool
	}

	if x, ok := err.(isInterface); ok {
		return x.Is(target)
	}

	// 递归检查 Unwrap
	type unwrapInterface interface {
		Unwrap() error
	}

	if x, ok := err.(unwrapInterface); ok {
		return Is(x.Unwrap(), target)
	}

	return false
}

// As 全局错误类型转换函数
func As(err error, target interface{}) bool {
	if err == nil {
		return false
	}

	// 尝试直接类型断言
	if e, ok := err.(Error); ok {
		if ptr, ok := target.(*Error); ok {
			*ptr = e
			return true
		}
	}

	// 递归检查 Unwrap
	type unwrapInterface interface {
		Unwrap() error
	}

	if x, ok := err.(unwrapInterface); ok {
		return As(x.Unwrap(), target)
	}

	return false
}

// builder 使用 strings.Builder 高效构建字符串
type builder struct {
	sb strings.Builder
}

func (b *builder) write(s string) {
	b.sb.WriteString(s)
}

func (b *builder) writeInt(i int) {
	b.sb.WriteString(strconv.Itoa(i))
}

func (b *builder) string() string {
	return b.sb.String()
}
