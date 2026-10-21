package goerr

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"slices"
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

func TestItem_FormatTextContract(t *testing.T) {
	t.Parallel()
	item := New(io.EOF, StatusParams(), "bad")
	for _, tt := range []struct {
		name   string
		item   *Item
		format string
		want   string
	}{
		{name: "value", item: item, format: "%v", want: "bad: EOF"},
		{name: "string", item: item, format: "%s", want: "bad: EOF"},
		{name: "quoted", item: item, format: "%q", want: `"bad: EOF"`},
		{name: "unknown", item: item, format: "%d", want: "%!d(*goerr.Item=bad: EOF)"},
		{name: "width and precision", item: item, format: "%8.2s", want: "bad: EOF"},
		{name: "alternate", item: item, format: "%#v", want: "bad: EOF"},
		{name: "nil verbose", format: "%+v", want: ""},
		{name: "nil quoted", format: "%q", want: `""`},
		{name: "zero verbose", item: &Item{}, format: "%+v", want: ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := fmt.Sprintf(tt.format, tt.item); got != tt.want {
				t.Fatalf("format %q = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}

type formatTestState struct {
	verbose bool
	failAt  int
	limit   int
	calls   int
	failure error
	written []byte
}

func (s *formatTestState) Width() (int, bool)     { return 0, false }
func (s *formatTestState) Precision() (int, bool) { return 0, false }
func (s *formatTestState) Flag(c int) bool        { return c == '+' && s.verbose }

func (s *formatTestState) Write(data []byte) (int, error) {
	s.calls++
	if s.calls == s.failAt {
		n := min(s.limit, len(data))
		s.written = append(s.written, data[:n]...)
		s.failure = io.ErrClosedPipe
		return n, s.failure
	}
	s.written = append(s.written, data...)
	return len(data), nil
}

func TestItem_FormatStopsAfterWriteError(t *testing.T) {
	t.Parallel()
	item := New(io.EOF, StatusParams(), "bad")
	stack := slices.Clone(item.stack)
	if len(stack) < 2 {
		t.Fatal("the fixture must have caller frames to detect writes after failure")
	}
	for _, tt := range []struct {
		name    string
		verb    rune
		verbose bool
		failAt  int
		limit   int
		want    string
	}{
		{name: "before verbose text", verb: 'v', verbose: true, failAt: 1},
		{name: "partial verbose text", verb: 'v', verbose: true, failAt: 1, limit: 3, want: "bad"},
		{name: "before frame", verb: 'v', verbose: true, failAt: 2, want: "bad: EOF"},
		{name: "partial frame", verb: 'v', verbose: true, failAt: 2, limit: 2, want: "bad: EOF\n\t"},
		{name: "value", verb: 'v', failAt: 1, limit: 3, want: "bad"},
		{name: "string", verb: 's', failAt: 1, limit: 3, want: "bad"},
		{name: "quoted", verb: 'q', failAt: 1, limit: 4, want: `"bad`},
		{name: "unknown", verb: 'd', failAt: 1, limit: 3, want: "%!d"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			state := &formatTestState{verbose: tt.verbose, failAt: tt.failAt, limit: tt.limit}
			item.Format(state, tt.verb)
			if !errors.Is(state.failure, io.ErrClosedPipe) || state.calls != tt.failAt || string(state.written) != tt.want {
				t.Fatalf("Format must stop after the failed write: calls=%d bytes=%q failure=%v", state.calls, state.written, state.failure)
			}
			if item.Error() != "bad: EOF" || item.Message() != "bad" || item.Code() != ErrParams ||
				item.HTTPCode() != http.StatusBadRequest || item.Unwrap() != io.EOF || !slices.Equal(item.stack, stack) {
				t.Fatal("formatting failure must preserve the original business error, cause and stack")
			}
		})
	}
}

func TestItem_FprintfWriterFailure(t *testing.T) {
	t.Parallel()
	item := New(io.EOF, StatusParams(), "bad")
	want := fmt.Sprintf("%+v", item)
	for _, limit := range []int{0, 7} {
		t.Run(fmt.Sprintf("accepted_%d_bytes", limit), func(t *testing.T) {
			t.Parallel()
			writer := &formatTestState{failAt: 1, limit: limit}
			n, err := fmt.Fprintf(writer, "%+v", item)
			if n != limit || !errors.Is(err, io.ErrClosedPipe) || writer.calls != 1 || string(writer.written) != want[:limit] {
				t.Fatalf("Fprintf must return the writer's partial result: n=%d bytes=%q err=%v", n, writer.written, err)
			}
			if item.Error() != "bad: EOF" || item.Message() != "bad" || !errors.Is(item, io.EOF) || errors.Is(item, io.ErrClosedPipe) {
				t.Fatal("the writer failure must remain separate from the business error")
			}
		})
	}
}

func TestIsCompatibility(t *testing.T) {
	t.Parallel()
	var typedNil *Item
	item := New(io.EOF, StatusParams(), "bad")
	otherItem := New(io.EOF, StatusParams(), "bad")
	for _, tt := range []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{name: "both nil", want: true},
		{name: "nil error", target: io.EOF},
		{name: "nil target", err: io.EOF},
		{name: "same sentinel", err: io.EOF, target: io.EOF, want: true},
		{name: "different sentinel", err: io.EOF, target: io.ErrUnexpectedEOF},
		{name: "same text distinct errors", err: errors.New("same"), target: errors.New("same")},
		{name: "standard wrapper", err: fmt.Errorf("load: %w", io.EOF), target: io.EOF, want: true},
		{name: "Item cause", err: item, target: io.EOF, want: true},
		{name: "Item identity", err: item, target: item, want: true},
		{name: "same code distinct Items", err: item, target: otherItem},
		{name: "Join branch", err: errors.Join(io.ErrUnexpectedEOF, item), target: io.EOF, want: true},
		{name: "typed nil target", err: typedNil, target: typedNil, want: true},
		{name: "typed nil is not nil", err: typedNil},
		{name: "typed nil error absent", target: typedNil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Is(tt.err, tt.target); got != tt.want || got != errors.Is(tt.err, tt.target) {
				t.Fatalf("Is(%v, %v) = %t, want %t and standard errors.Is compatibility", tt.err, tt.target, got, tt.want)
			}
		})
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
	for l := range strings.SplitSeq(fmt.Sprintf("%+v", err), "\n") {
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

func TestCompatibilityFirstItem(t *testing.T) {
	t.Parallel()
	first := New(io.EOF, StatusParams(), "first public message")
	second := New(io.ErrClosedPipe, StatusNotFound(), "second public message")
	override := WrapStatus(second, StatusRedisServer())
	overrideItem, ok := AsItem(override)
	if !ok || overrideItem == nil {
		t.Fatal("WrapStatus must create an Item")
	}
	for _, tt := range []struct {
		name   string
		input  error
		want   *Item
		hidden *Item
	}{
		{name: "standard wrapper", input: fmt.Errorf("outer: %w", first), want: first},
		{name: "join", input: errors.Join(io.EOF, first, second), want: first, hidden: second},
		{name: "reversed join", input: errors.Join(second, first), want: second, hidden: first},
		{name: "multiple wraps", input: fmt.Errorf("%w | %w", first, second), want: first, hidden: second},
		{name: "indexed wraps use argument order", input: fmt.Errorf("%[2]w | %[1]w", first, second), want: first, hidden: second},
		{name: "status override", input: override, want: overrideItem, hidden: second},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Convert(tt.input); got != tt.want {
				t.Errorf("Convert must return the first Item, got %p want %p", got, tt.want)
			}
			if !HasCode(tt.input, tt.want.Code()) || !errors.Is(tt.input, tt.want) {
				t.Error("first Item must supply the code and remain reachable")
			}
			if tt.hidden != nil && (HasCode(tt.input, tt.hidden.Code()) || !errors.Is(tt.input, tt.hidden)) {
				t.Error("other Item must remain reachable without changing the selected code")
			}
			if got := WithStack(tt.input); got != tt.input {
				t.Error("an existing stack must preserve the original outer error")
			}
		})
	}
	if other := New(io.EOF, StatusParams()); errors.Is(first, other) {
		t.Error("equal business codes must not imply error identity")
	}
}

type compatibilityAsError struct{ item *Item }

func (*compatibilityAsError) Error() string { return "outer custom As" }

func (e *compatibilityAsError) As(target any) bool {
	if item, ok := target.(**Item); ok {
		*item = e.item
		return true
	}
	return false
}

func TestCompatibilityCustomAs(t *testing.T) {
	t.Parallel()
	inner := New(io.EOF, StatusParams(), "public message")
	input := &compatibilityAsError{item: inner}
	if got, ok := AsItem(input); !ok || got != inner || Convert(input) != inner || !HasCode(input, ErrParams) {
		t.Fatal("custom As must keep the existing Item selection")
	}
	if WithStack(input) != input || errors.Is(input, inner) {
		t.Error("As must not become an Unwrap edge or discard the outer identity")
	}
}

type compatibilityForeignError struct{ status *Status }

func (*compatibilityForeignError) Error() string         { return "foreign error" }
func (e *compatibilityForeignError) StatusInfo() *Status { return e.status }
func (e *compatibilityForeignError) Code() Code          { return e.status.Code() }
func (*compatibilityForeignError) Message() string       { return "foreign public message" }
func (e *compatibilityForeignError) HTTPCode() int       { return e.status.HTTPCode() }

func TestCompatibilityForeignError(t *testing.T) {
	t.Parallel()
	var input Error = &compatibilityForeignError{status: StatusParams()}
	if item, ok := AsItem(input); ok || item != nil || HasCode(input, ErrParams) {
		t.Fatal("the legacy query must not treat every Error implementation as an Item")
	}
	for _, got := range []*Item{Convert(input), Convert(Wrap(input, "outer"))} {
		if got.Code() != ErrInternalServer || got.HTTPCode() != http.StatusInternalServerError ||
			got.Message() != StatusInternalServer().Msg() || !errors.Is(got, input) {
			t.Errorf("non-Item error must keep the internal fallback and cause: %+v", got)
		}
	}
}

func TestCompatibilityTypedNil(t *testing.T) {
	t.Parallel()
	var item *Item
	for _, tt := range []struct {
		name  string
		input error
	}{
		{name: "direct", input: item},
		{name: "first join branch", input: errors.Join(item, New(io.EOF, StatusNotFound()))},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got, ok := AsItem(tt.input); !ok || got != nil {
				t.Fatal("the first typed-nil Item must remain a nil/true match")
			}
			got := Convert(tt.input)
			if got == nil || got.Code() != ErrInternalServer || got.HTTPCode() != http.StatusInternalServerError ||
				got.Message() != StatusInternalServer().Msg() || !errors.Is(got, tt.input) {
				t.Fatal("Convert must preserve the non-nil failure fallback and cause")
			}
			if !HasCode(tt.input, ErrInternalServer) || HasCode(tt.input, ErrNotFound) || WithStack(tt.input) != tt.input {
				t.Error("typed-nil selection and WithStack identity must remain unchanged")
			}
		})
	}
}

type compatibilityMutableError struct{ text string }

func (e *compatibilityMutableError) Error() string { return e.text }

func TestCompatibilityTextSnapshot(t *testing.T) {
	t.Parallel()
	cause := &compatibilityMutableError{text: "original cause"}
	inner := New(cause, StatusParams(), "public message")
	outer := Wrap(inner, "outer")
	cause.text = "changed cause"
	// 这里核对明确的全文缓存契约，不用文案来判定错误类别。
	if inner.Error() != "public message: original cause" || outer.Error() != "outer: public message: original cause" {
		t.Error("construction must preserve its original text snapshot")
	}
	if Convert(outer).Message() != "public message" || !errors.Is(outer, cause) || !errors.Is(outer, inner) {
		t.Error("cached text must not change the public message or sever cause identity")
	}
}

// 阶段 1 建立的失败反例，作为缺栈输入的补栈回归保留。
func TestWithStackMissingStack(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name  string
		input error
	}{
		{name: "converted_plain", input: Convert(io.EOF)},
		{name: "zero_item", input: &Item{}},
		{name: "empty_non_nil_stack", input: &Item{stack: []uintptr{}}},
		{name: "converted_nil", input: Convert(nil)},
		{name: "wrapped_plain", input: fmt.Errorf("outer: %w", Convert(io.EOF))},
		{name: "nil_status", input: &Item{msg: "plain: EOF", clientMsg: "plain public", cause: io.EOF}},
		{name: "custom_status", input: &Item{msg: "custom: EOF", clientMsg: "custom public", status: NewStatus(1, 299, "custom"), cause: io.EOF}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			original := Convert(tt.input)
			if len(original.stack) != 0 {
				t.Fatal("fixture must have no existing stack")
			}
			wantText, wantMessage := tt.input.Error(), original.Message()
			wantCode, wantHTTP := original.Code(), original.HTTPCode()
			got := WithStack(tt.input)
			item, ok := AsItem(got)
			if !ok || item == nil {
				t.Fatal("WithStack must produce a non-nil Item")
			}
			if got == tt.input || errors.Unwrap(got) != tt.input || !errors.Is(got, tt.input) {
				t.Error("missing stack requires a new wrapper with the original input as its direct cause")
			}
			if got.Error() != wantText || item.Message() != wantMessage || item.Code() != wantCode || item.HTTPCode() != wantHTTP {
				t.Error("adding a stack must preserve the full text, public message and status")
			}
			if len(item.stack) == 0 {
				t.Error("WithStack returned an Item without a stack")
			} else {
				frame, _ := runtime.CallersFrames(item.stack).Next()
				if !strings.HasSuffix(frame.File, "/error_test.go") || !strings.Contains(frame.Function, "TestWithStackMissingStack") {
					t.Errorf("first frame must be the caller, got %s:%d %s", frame.File, frame.Line, frame.Function)
				}
			}
			if len(original.stack) != 0 || WithStack(got) != got {
				t.Error("original input must stay unchanged and the new wrapper must be reused")
			}
		})
	}
}

func TestWithStackChains(t *testing.T) {
	t.Parallel()
	missing := Convert(io.EOF)
	stacked := New(io.ErrClosedPipe, StatusNotFound(), "public message")
	for _, tt := range []struct {
		name        string
		input       error
		wantCapture bool
	}{
		{name: "join missing first", input: errors.Join(missing, stacked), wantCapture: true},
		{name: "join stacked first", input: errors.Join(stacked, missing)},
		{name: "multiple wraps", input: fmt.Errorf("outer: %w | %w", missing, stacked), wantCapture: true},
		{name: "shared missing branch", input: errors.Join(missing, missing), wantCapture: true},
		{name: "custom As missing", input: &compatibilityAsError{item: missing}, wantCapture: true},
		{name: "custom As typed nil", input: &compatibilityAsError{}},
		{name: "foreign Error", input: &compatibilityForeignError{status: StatusParams()}, wantCapture: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			original := Convert(tt.input)
			wantText, wantMessage := tt.input.Error(), original.Message()
			wantCode, wantHTTP := original.Code(), original.HTTPCode()
			got := WithStack(tt.input)
			if !tt.wantCapture {
				if got != tt.input {
					t.Error("an existing stack or typed-nil match must preserve the original input")
				}
				return
			}
			item, ok := AsItem(got)
			if !ok || item == nil || len(item.stack) == 0 {
				t.Fatal("the first missing stack must be captured")
			}
			if got == tt.input || errors.Unwrap(got) != tt.input || !errors.Is(got, tt.input) {
				t.Error("new wrapper must retain the original outer input as its direct cause")
			}
			if got.Error() != wantText || item.Message() != wantMessage || item.Code() != wantCode || item.HTTPCode() != wantHTTP {
				t.Error("branch structure or custom As must not change the text or selected status")
			}
			for _, cause := range []error{io.EOF, io.ErrClosedPipe} {
				if errors.Is(got, cause) != errors.Is(tt.input, cause) {
					t.Errorf("cause reachability changed for %v", cause)
				}
			}
			frame, _ := runtime.CallersFrames(item.stack).Next()
			if !strings.HasSuffix(frame.File, "/error_test.go") || !strings.Contains(frame.Function, "TestWithStackChains") {
				t.Errorf("first frame must be the caller, got %s:%d %s", frame.File, frame.Line, frame.Function)
			}
			downstream, ok := AsItem(Wrap(got, "downstream"))
			if !ok || downstream == nil || len(downstream.stack) == 0 || &downstream.stack[0] != &item.stack[0] {
				t.Error("downstream wrappers must reuse the captured stack")
			}
			if len(original.stack) != 0 || WithStack(got) != got {
				t.Error("the original must stay unchanged and a second WithStack must reuse the wrapper")
			}
		})
	}
}
