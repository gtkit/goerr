package goerr

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestNew_nilCause(t *testing.T) {
	item := New(nil, StatusNotFound(), "订单不存在")
	if got := item.Error(); got != "订单不存在" {
		t.Errorf("Error() = %q, want %q", got, "订单不存在")
	}
	if got := item.Message(); got != "订单不存在" {
		t.Errorf("Message() = %q, want custom client message", got)
	}
	if item.Code() != ErrNotFound {
		t.Errorf("Code() = %v", item.Code())
	}
	if item.Unwrap() != nil {
		t.Error("Unwrap should be nil")
	}
}

func TestNew_withCause_messageVsError(t *testing.T) {
	underlying := errors.New("dial tcp: refused")
	item := New(underlying, StatusMysqlServer(), "数据库不可用")

	// Message()：给客户端/统一响应体 —— 不含底层实现细节
	if got := item.Message(); got != "数据库不可用" {
		t.Errorf("Message() = %q", got)
	}
	// Error()：给日志 —— 含 cause 全文
	if want := "数据库不可用: dial tcp: refused"; item.Error() != want {
		t.Errorf("Error() = %q, want %q", item.Error(), want)
	}
	if !errors.Is(item, underlying) {
		t.Error("errors.Is should find underlying")
	}
}

func TestNew_defaultMsgWhenNoOverride(t *testing.T) {
	item := New(nil, StatusParams())
	if got := item.Message(); got != StatusParams().Msg() {
		t.Errorf("Message() = %q, want status default %q", got, StatusParams().Msg())
	}
	if item.Error() != StatusParams().Msg() {
		t.Errorf("Error() should match when no cause and no override")
	}
}

func TestNewFn(t *testing.T) {
	item := NewFn(nil, StatusAuth, "请先登录")
	if item.Message() != "请先登录" {
		t.Errorf("Message() = %q", item.Message())
	}
	if item.Code() != ErrAuthentication {
		t.Errorf("Code = %v", item.Code())
	}
}

func TestNewf(t *testing.T) {
	item := Newf(StatusInvalidJson(), "invalid field %q", "age")
	want := `invalid field "age"`
	if item.Error() != want {
		t.Errorf("Error() = %q, want %q", item.Error(), want)
	}
	if item.Message() != want {
		t.Errorf("Message() = %q, want %q", item.Message(), want)
	}
}

func TestNewf_nilUsesOK(t *testing.T) {
	item := Newf(nil, "only %s", "text")
	if item.StatusInfo().Code() != ErrNo {
		t.Errorf("nil status should be StatusOK: %v", item.Code())
	}
}

func TestWrap_innerItem_preservesClientMessage(t *testing.T) {
	inner := New(io.EOF, StatusRecordNotFound(), "记录不存在")
	w := Wrap(inner, "查询用户")
	if w.Message() != inner.Message() {
		t.Errorf("Wrap should preserve inner Message() for API; got %q want %q", w.Message(), inner.Message())
	}
	if !strings.Contains(w.Error(), "查询用户:") {
		t.Errorf("Error() should contain wrap prefix: %q", w.Error())
	}
	if !errors.Is(w, io.EOF) {
		t.Error("errors.Is(w, EOF) should hold")
	}
}

func TestWrap_plainError(t *testing.T) {
	base := errors.New("root")
	w := Wrap(base, "load config")
	if w.Message() != "load config" {
		t.Errorf("Message() = %q", w.Message())
	}
	if !strings.Contains(w.Error(), "load config") || !strings.Contains(w.Error(), "root") {
		t.Errorf("Error() = %q", w.Error())
	}
}

func TestWrapf(t *testing.T) {
	inner := Newf(StatusParams(), "bad")
	w := Wrapf(inner, "step %d", 1)
	if w.Message() != inner.Message() {
		t.Errorf("Message = %q", w.Message())
	}
}

func TestWrap_nil(t *testing.T) {
	if Wrap(nil, "x") != nil {
		t.Fatal("Wrap(nil) should be nil")
	}
}

func TestWrapStatus(t *testing.T) {
	e := errors.New("boom")
	w := WrapStatus(e, StatusRedisServer())
	if w.Message() != StatusRedisServer().Msg() {
		t.Errorf("Message = %q", w.Message())
	}
	if w.Code() != ErrRedisServer {
		t.Errorf("Code = %v", w.Code())
	}
	if !errors.Is(w, e) {
		t.Error("unwrap chain")
	}
}

func TestWithStack_plain(t *testing.T) {
	e := errors.New("x")
	w := WithStack(e)
	if w.Message() != "x" || w.Error() != "x" {
		t.Errorf("WithStack plain: Message=%q Error=%q", w.Message(), w.Error())
	}
}

func TestWithStack_returnsSameItem(t *testing.T) {
	inner := New(nil, StatusOK())
	w := WithStack(inner)
	if w != inner {
		t.Error("WithStack(*Item) should return same pointer")
	}
}

func TestAsItem(t *testing.T) {
	inner := New(nil, StatusNotFound())
	wrapped := fmt.Errorf("outer: %w", inner)
	got, ok := AsItem(wrapped)
	if !ok || got.Code() != ErrNotFound {
		t.Fatalf("AsItem: ok=%v code=%v", ok, got.Code())
	}
}

func TestItem_FormatVerbose(t *testing.T) {
	item := New(errors.New("cause"), StatusParams(), "bad")
	s := fmt.Sprintf("%+v", item)
	if !strings.Contains(s, "bad") || !strings.Contains(s, "Caused by") {
		t.Errorf("verbose format: %s", s)
	}
}

func TestItem_StatusInfo_nilStatus(t *testing.T) {
	item := &Item{msg: "m", clientMsg: "m"}
	if item.StatusInfo().Code() != ErrNo {
		t.Errorf("nil status should fall back to OK")
	}
}
