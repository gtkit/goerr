package goerr

import (
	"fmt"
	"testing"
)

func TestNewStatus_allowsCustomCode(t *testing.T) {
	st := NewStatus(Code(1), 299, "custom")
	if st.Code() != Code(1) || st.HTTPCode() != 299 || st.Msg() != "custom" {
		t.Fatalf("NewStatus() = code %v http %d msg %q", st.Code(), st.HTTPCode(), st.Msg())
	}
}

func TestMustNewStatus(t *testing.T) {
	st := MustNewStatus(ErrOrderExpired, 200, "Order expired")
	if st.Code() != ErrOrderExpired || st.HTTPCode() != 200 || st.Msg() != "Order expired" {
		t.Fatalf("MustNewStatus() = code %v http %d msg %q", st.Code(), st.HTTPCode(), st.Msg())
	}
}

func TestMustNewStatus_panicsOnInvalidCode(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = MustNewStatus(Code(1), 200, "invalid")
}

func ExampleNewStatus() {
	// 自定义错误码，不校验编码规范，业务方可自由编码
	st := NewStatus(Code(900001), 200, "Custom business error")
	fmt.Println(st.Code())
	fmt.Println(st.HTTPCode())
	fmt.Println(st.Msg())
	// Output:
	// 900001
	// 200
	// Custom business error
}

func ExampleMustNewStatus() {
	// 自定义错误码，强制校验编码规范，非法 code 会 panic
	st := MustNewStatus(ErrOrderExpired, 200, "Order expired")
	fmt.Println(st.Code())
	fmt.Println(st.HTTPCode())
	fmt.Println(st.Msg())
	// Output:
	// 10010707
	// 200
	// Order expired
}
