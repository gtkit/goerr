package goerr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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

func TestWrap_innerItem_preservesClientMessage(t *testing.T) {
	inner := New(io.EOF, StatusRecordNotFound(), "记录不存在")
	w := Wrap(inner, "查询用户")
	wi, _ := AsItem(w)
	if wi.Message() != inner.Message() {
		t.Errorf("Wrap should preserve inner Message() for API; got %q want %q", wi.Message(), inner.Message())
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
	wi, _ := AsItem(w)
	if wi.Message() != StatusInternalServer().Msg() {
		t.Errorf("Message() = %q", wi.Message())
	}
	if wi.Code() != ErrInternalServer {
		t.Errorf("Code() = %v", wi.Code())
	}
	if !strings.Contains(w.Error(), "load config") || !strings.Contains(w.Error(), "root") {
		t.Errorf("Error() = %q", w.Error())
	}
}

func TestWrapf(t *testing.T) {
	inner := Newf(StatusParams(), "bad")
	w := Wrapf(inner, "step %d", 1)
	wi, _ := AsItem(w)
	if wi.Message() != inner.Message() {
		t.Errorf("Message = %q", wi.Message())
	}
}

func TestWrapf_plainErrorDefaultsInternalServer(t *testing.T) {
	base := errors.New("root")
	w := Wrapf(base, "step %d", 1)
	wi, _ := AsItem(w)
	if wi.Message() != StatusInternalServer().Msg() {
		t.Errorf("Message() = %q", wi.Message())
	}
	if wi.Code() != ErrInternalServer {
		t.Errorf("Code() = %v", wi.Code())
	}
	if !strings.Contains(w.Error(), "step 1") || !strings.Contains(w.Error(), "root") {
		t.Errorf("Error() = %q", w.Error())
	}
}

func TestWrapStatus(t *testing.T) {
	e := errors.New("boom")
	w := WrapStatus(e, StatusRedisServer())
	wi, _ := AsItem(w)
	if wi.Message() != StatusRedisServer().Msg() {
		t.Errorf("Message = %q", wi.Message())
	}
	if wi.Code() != ErrRedisServer {
		t.Errorf("Code = %v", wi.Code())
	}
	if !errors.Is(w, e) {
		t.Error("unwrap chain")
	}
}

func TestWithStack_plain(t *testing.T) {
	e := errors.New("x")
	w := WithStack(e)
	wi, _ := AsItem(w)
	if wi.Message() != StatusInternalServer().Msg() || w.Error() != "x" {
		t.Errorf("WithStack plain: Message=%q Error=%q", wi.Message(), w.Error())
	}
	if wi.Code() != ErrInternalServer {
		t.Errorf("Code() = %v", wi.Code())
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
	if !strings.Contains(s, "bad: cause") {
		t.Errorf("verbose format should carry full chain text: %s", s)
	}
	if !strings.Contains(s, "error_test.go") {
		t.Errorf("verbose format should carry caller frame: %s", s)
	}
}

func TestItem_StatusInfo_nilStatus(t *testing.T) {
	item := &Item{msg: "m", clientMsg: "m"}
	if item.StatusInfo().Code() != ErrInternalServer {
		t.Errorf("nil status should fall back to InternalServer, got %v", item.Code())
	}
}

// --- nil 语义 ---

func TestItem_nilReceiverSafe(t *testing.T) {
	var i *Item
	if i.Error() != "" || i.Message() != "" || i.Unwrap() != nil {
		t.Fatal("nil receiver accessors should return zero values")
	}
	if i.Code() != ErrInternalServer || i.HTTPCode() != http.StatusInternalServerError {
		t.Fatalf("nil receiver status should fall back to InternalServer, got %v/%d", i.Code(), i.HTTPCode())
	}
	if s := fmt.Sprintf("%+v|%s|%q", i, i, i); s != `||""` {
		t.Fatalf("nil receiver format = %q", s)
	}
}

func TestWrapFamily_nilReturnsUntypedNil(t *testing.T) {
	wsf, rerr := WrapStatusFrom(nil, StatusParams())
	if rerr != nil {
		t.Fatal(rerr)
	}
	for name, got := range map[string]error{
		"Wrap":           Wrap(nil, "x"),
		"Wrapf":          Wrapf(nil, "x%d", 1),
		"WrapStatus":     WrapStatus(nil, StatusParams()),
		"WithStack":      WithStack(nil),
		"WrapStatusFrom": wsf,
	} {
		if got != nil { // 接口比较：类型化 nil 会在这里暴露
			t.Errorf("%s(nil) = %#v, want untyped nil", name, got)
		}
	}
}

func TestTypedNilCause_noPanic(t *testing.T) {
	var cause error = (*Item)(nil)
	for name, err := range map[string]error{
		"New":        New(cause, StatusParams()),
		"Wrap":       Wrap(cause, "x"),
		"WrapStatus": WrapStatus(cause, StatusParams()),
		"WithStack":  WithStack(cause),
		"Convert":    Convert(cause),
	} {
		if err == nil {
			t.Errorf("%s: typed-nil cause is a non-nil error, result must not be nil", name)
		}
	}
}

// --- 状态兜底 ---

func TestNilStatus_fallsBackToInternalServer(t *testing.T) {
	nf, err := NewFrom(io.EOF, nil)
	if err != nil {
		t.Fatal(err)
	}
	for name, it := range map[string]*Item{
		"New":     New(io.EOF, nil),
		"Newf":    Newf(nil, "x"),
		"NewFrom": nf,
		"zero":    {msg: "m"},
	} {
		if it.Code() != ErrInternalServer {
			t.Errorf("%s: nil status should fall back to InternalServer, got %v", name, it.Code())
		}
	}
}

func TestWithStack_keepsOuterWrapping(t *testing.T) {
	inner := New(errors.New("root"), StatusNotFound(), "找不到")
	outer := fmt.Errorf("外层: %w", inner)
	if got := WithStack(outer); got != outer {
		t.Fatalf("chain already has *Item: should return the original error, got %v", got)
	}
}

// --- Convert ---

func TestConvert(t *testing.T) {
	if got := Convert(nil); got.Code() != ErrNo || got.HTTPCode() != http.StatusOK || got.Message() != StatusOK().Msg() {
		t.Errorf("Convert(nil) = %v/%d/%q", got.Code(), got.HTTPCode(), got.Message())
	}
	item := New(io.EOF, StatusNotFound(), "找不到")
	if got := Convert(fmt.Errorf("ctx: %w", item)); got != item {
		t.Errorf("Convert should return the *Item in chain, got %v", got)
	}
	got := Convert(io.EOF)
	if got.Code() != ErrInternalServer || got.Message() != StatusInternalServer().Msg() ||
		got.Error() != io.EOF.Error() || !errors.Is(got, io.EOF) {
		t.Errorf("Convert(plain) = %v/%q/%q", got.Code(), got.Message(), got.Error())
	}
}

// --- 调用栈 ---

// stackFrames 返回 %+v 里的栈帧行。
func stackFrames(err error) []string {
	var out []string
	for _, l := range strings.Split(fmt.Sprintf("%+v", err), "\n") {
		if strings.HasPrefix(l, "\t") {
			out = append(out, l)
		}
	}
	return out
}

func TestStack_firstFrameIsCaller(t *testing.T) {
	nf, _ := NewFrom(nil, StatusParams)
	nff, _ := NewfFrom(StatusParams(), "x")
	wsf, _ := WrapStatusFrom(io.EOF, StatusParams())
	for name, err := range map[string]error{
		"New":            New(nil, StatusParams()),
		"NewFn":          NewFn(nil, StatusParams),
		"Newf":           Newf(StatusParams(), "x"),
		"NewFrom":        nf,
		"NewfFrom":       nff,
		"Wrap":           Wrap(io.EOF, "x"),
		"Wrapf":          Wrapf(io.EOF, "x%d", 1),
		"WrapStatus":     WrapStatus(io.EOF, StatusParams()),
		"WrapStatusFrom": wsf,
		"WithStack":      WithStack(io.EOF),
	} {
		frames := stackFrames(err)
		if len(frames) == 0 || !strings.Contains(frames[0], "error_test.go") {
			t.Errorf("%s: first frame should be the caller in error_test.go, got %q", name, frames)
		}
		for _, f := range frames {
			if strings.Contains(f, "/error.go:") || strings.Contains(f, "/resolve.go:") {
				t.Errorf("%s: internal frame leaked into stack: %s", name, f)
			}
		}
	}
}

func TestStack_sharedAcrossChain(t *testing.T) {
	base := New(io.EOF, StatusParams(), "a")
	w, _ := AsItem(Wrapf(Wrap(base, "b"), "c%d", 1))
	if len(base.stack) == 0 || &w.stack[0] != &base.stack[0] {
		t.Fatal("wrapped items should share the innermost stack, not capture a new one")
	}
	// 若 Format 仍递归打印 cause，三层 Item 会打出三份栈
	if got, want := len(stackFrames(w)), len(stackFrames(base)); got != want {
		t.Fatalf("%%+v of a 3-level chain printed %d frames, want exactly one stack of %d", got, want)
	}
}

func TestItem_Format_unknownVerb(t *testing.T) {
	if s := fmt.Sprintf("%d", New(nil, StatusParams())); s != "%!d(*goerr.Item=参数不合法)" {
		t.Errorf("unknown verb = %q", s)
	}
}
