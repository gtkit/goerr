package goerr

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

// TestNewError 测试创建错误
func TestNewError(t *testing.T) {
	err := New("TEST_001", "test error")

	if err == nil {
		t.Fatal("expected non-nil error")
	}

	if err.Code() != "TEST_001" {
		t.Errorf("expected code TEST_001, got %s", err.Code())
	}

	if err.Message() != "test error" {
		t.Errorf("expected message 'test error', got %s", err.Message())
	}

	if err.HTTPStatus() != 500 {
		t.Errorf("expected status 500, got %d", err.HTTPStatus())
	}
}

// TestNewf 测试格式化创建错误
func TestNewf(t *testing.T) {
	err := Newf("TEST_002", "user %s not found", "alice")

	if err.Message() != "user alice not found" {
		t.Errorf("expected formatted message, got %s", err.Message())
	}
}

// TestWrapError 测试包装错误
func TestWrapError(t *testing.T) {
	baseErr := errors.New("base error")
	wrapped := Wrap(baseErr, "WRAP_001", "wrapped error")

	if wrapped.Code() != "WRAP_001" {
		t.Errorf("expected code WRAP_001, got %s", wrapped.Code())
	}

	if wrapped.Unwrap() != baseErr {
		t.Error("expected wrapped error to contain base error")
	}

	// 测试包装 nil
	nilWrapped := Wrap(nil, "WRAP_002", "should be nil")
	if nilWrapped != nil {
		t.Error("wrapping nil should return nil")
	}
}

// TestWithDetails 测试添加详细信息
func TestWithDetails(t *testing.T) {
	err := New("TEST_003", "test error")

	// 添加单个详细信息
	err2 := err.WithDetail("user_id", 12345)
	details := err2.Details()

	if details["user_id"] != 12345 {
		t.Errorf("expected user_id=12345, got %v", details["user_id"])
	}

	// 原始错误不应该被修改
	if len(err.Details()) != 0 {
		t.Error("original error should not be modified")
	}

	// 批量添加
	err3 := err.WithDetails(map[string]interface{}{
		"request_id": "req-123",
		"ip":         "192.168.1.1",
	})

	details3 := err3.Details()
	if len(details3) != 2 {
		t.Errorf("expected 2 details, got %d", len(details3))
	}
}

// TestErrorIs 测试错误判断
func TestErrorIs(t *testing.T) {
	err1 := New("ERR_001", "error 1")
	err2 := New("ERR_001", "error 1 again")
	err3 := New("ERR_002", "error 2")

	// 相同错误码
	if !Is(err1, err2) {
		t.Error("expected errors with same code to be equal")
	}

	// 不同错误码
	if Is(err1, err3) {
		t.Error("expected errors with different codes to be different")
	}

	// nil 检查
	if Is(err1, nil) {
		t.Error("error should not equal nil")
	}

	if Is(nil, err1) {
		t.Error("nil should not equal error")
	}

	if !Is(nil, nil) {
		t.Error("nil should equal nil")
	}
}

// TestErrorAs 测试错误类型转换
func TestErrorAs(t *testing.T) {
	err := New("TEST_004", "test error")

	var target Error
	if !As(err, &target) {
		t.Error("As should succeed for correct type")
	}

	if target.Code() != "TEST_004" {
		t.Errorf("expected code TEST_004, got %s", target.Code())
	}

	// 测试失败情况
	var wrongTarget *struct{}
	if As(err, wrongTarget) {
		t.Error("As should fail for incorrect type")
	}
}

// TestErrorFormat 测试错误格式化
func TestErrorFormat(t *testing.T) {
	err := New("TEST_005", "test error").
		WithDetail("key1", "value1").
		WithDetail("key2", 123)

	// 简单格式化
	simple := fmt.Sprintf("%s", err)
	if !strings.Contains(simple, "TEST_005") {
		t.Errorf("simple format should contain code, got: %s", simple)
	}

	// 详细格式化
	detailed := fmt.Sprintf("%+v", err)
	if !strings.Contains(detailed, "Details:") {
		t.Errorf("detailed format should contain details section, got: %s", detailed)
	}
	if !strings.Contains(detailed, "Stack trace:") {
		t.Errorf("detailed format should contain stack trace, got: %s", detailed)
	}
}

// TestStackTrace 测试堆栈追踪
func TestStackTrace(t *testing.T) {
	err := New("TEST_006", "test error")
	stack := err.Stack()

	if len(stack) == 0 {
		t.Error("expected non-empty stack trace")
	}

	// 检查第一帧包含当前函数
	if !strings.Contains(stack[0].Function, "TestStackTrace") {
		t.Errorf("expected stack to contain TestStackTrace, got: %s", stack[0].Function)
	}
}

// TestRegisterCode 测试错误码注册
func TestRegisterCode(t *testing.T) {
	code := "TEST_REG_001"
	RegisterCode(code, "test registered error", 400)

	cfg := GetCodeConfig(code)
	if cfg == nil {
		t.Fatal("expected config to be registered")
	}

	if cfg.Message != "test registered error" {
		t.Errorf("expected message 'test registered error', got %s", cfg.Message)
	}

	if cfg.HTTPStatus != 400 {
		t.Errorf("expected status 400, got %d", cfg.HTTPStatus)
	}
}

// TestNewWithCode 测试使用错误码创建错误
func TestNewWithCode(t *testing.T) {
	code := "TEST_WITH_CODE_001"
	RegisterCode(code, "registered error message", 404)

	err := NewWithCode(code)
	if err.Code() != code {
		t.Errorf("expected code %s, got %s", code, err.Code())
	}

	if err.Message() != "registered error message" {
		t.Errorf("expected registered message, got %s", err.Message())
	}

	if err.HTTPStatus() != 404 {
		t.Errorf("expected status 404, got %d", err.HTTPStatus())
	}

	// 测试未注册的错误码
	unregistered := NewWithCode("UNREGISTERED_CODE")
	if unregistered.Message() != "Unknown error" {
		t.Errorf("expected default message for unregistered code, got %s", unregistered.Message())
	}
}

// TestConcurrentAccess 测试并发安全
func TestConcurrentAccess(t *testing.T) {
	const goroutines = 100
	const iterations = 1000

	code := "CONCURRENT_TEST_001"
	RegisterCode(code, "concurrent test", 500)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				// 测试创建错误
				err := NewWithCode(code)
				if err.Code() != code {
					t.Errorf("concurrent error creation failed")
					return
				}

				// 测试获取配置
				cfg := GetCodeConfig(code)
				if cfg == nil {
					t.Errorf("concurrent config access failed")
					return
				}

				// 测试添加详细信息
				err2 := err.WithDetail(fmt.Sprintf("key_%d", id), j)
				if err2.Details()[fmt.Sprintf("key_%d", id)] != j {
					t.Errorf("concurrent detail addition failed")
					return
				}
			}
		}(i)
	}

	wg.Wait()
}

// TestErrorChain 测试错误链
func TestErrorChain(t *testing.T) {
	baseErr := errors.New("base error")

	err1 := Wrap(baseErr, "WRAP_1", "first wrap")
	err2 := Wrap(err1, "WRAP_2", "second wrap")
	err3 := Wrap(err2, "WRAP_3", "third wrap")

	// 检查错误链
	if err3.Code() != "WRAP_3" {
		t.Errorf("expected top code WRAP_3, got %s", err3.Code())
	}

	// Unwrap 测试
	unwrapped := err3.Unwrap()
	if unwrapped == nil {
		t.Fatal("expected unwrapped error")
	}

	if e, ok := unwrapped.(Error); ok {
		if e.Code() != "WRAP_2" {
			t.Errorf("expected unwrapped code WRAP_2, got %s", e.Code())
		}
	}
}

// TestPredefinedErrors 测试预定义错误码
func TestPredefinedErrors(t *testing.T) {
	testCases := []struct {
		code       string
		httpStatus int
	}{
		{ErrInternal, 500},
		{ErrInvalidArgument, 400},
		{ErrNotFound, 404},
		{ErrAlreadyExists, 409},
		{ErrPermissionDenied, 403},
		{ErrUnauthenticated, 401},
		{ErrUnavailable, 503},
	}

	for _, tc := range testCases {
		t.Run(tc.code, func(t *testing.T) {
			err := NewWithCode(tc.code)
			if err.HTTPStatus() != tc.httpStatus {
				t.Errorf("expected status %d, got %d", tc.httpStatus, err.HTTPStatus())
			}
		})
	}
}

// TestDetailsImmutability 测试详细信息不可变性
func TestDetailsImmutability(t *testing.T) {
	err1 := New("TEST_007", "test error")
	err2 := err1.WithDetail("key1", "value1")
	err3 := err2.WithDetail("key2", "value2")

	// 检查每个错误的详细信息是独立的
	if len(err1.Details()) != 0 {
		t.Error("err1 should have no details")
	}

	if len(err2.Details()) != 1 {
		t.Errorf("err2 should have 1 detail, got %d", len(err2.Details()))
	}

	if len(err3.Details()) != 2 {
		t.Errorf("err3 should have 2 details, got %d", len(err3.Details()))
	}

	// 修改返回的 details 不应该影响原错误
	details := err3.Details()
	details["key3"] = "value3"

	if len(err3.Details()) != 2 {
		t.Error("modifying returned details should not affect original error")
	}
}

// TestNilError 测试 nil 错误处理
func TestNilError(t *testing.T) {
	var err error = nil

	wrapped := Wrap(err, "WRAP_NIL", "should be nil")
	if wrapped != nil {
		t.Error("wrapping nil should return nil")
	}

	// As 测试
	var target Error
	if As(nil, &target) {
		t.Error("As with nil error should return false")
	}

	// Is 测试
	if Is(nil, New("TEST", "test")) {
		t.Error("nil should not equal any error")
	}
}
