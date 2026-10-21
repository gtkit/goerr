package goerr

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
	"unicode/utf8"
)

func TestInspectCompatibility(t *testing.T) {
	t.Parallel()
	first := New(io.EOF, StatusParams(), "public message")
	second := New(nil, StatusNotFound())
	var typedNil *Item
	for _, tt := range []struct {
		name   string
		input  error
		causes []Cause
	}{
		{name: "nil"},
		{name: "plain", input: io.EOF, causes: []Cause{{Parent: -1}}},
		{name: "zero Item", input: &Item{}, causes: []Cause{{Parent: -1, HasCode: true, Code: ErrInternalServer}}},
		{name: "non-nil success", input: Convert(nil), causes: []Cause{{Parent: -1, HasCode: true, Code: ErrNo}}},
		{name: "known typed nil", input: typedNil, causes: []Cause{{Parent: -1}}},
		{name: "standard prefix", input: fmt.Errorf("outer: %w", first), causes: []Cause{
			{Parent: -1}, {Parent: 0, HasCode: true, Code: ErrParams}, {Parent: 1},
		}},
		{name: "join", input: errors.Join(first, second), causes: []Cause{
			{Parent: -1}, {Parent: 0, HasCode: true, Code: ErrParams}, {Parent: 1}, {Parent: 0, HasCode: true, Code: ErrNotFound},
		}},
		{name: "reversed join", input: errors.Join(second, first), causes: []Cause{
			{Parent: -1}, {Parent: 0, HasCode: true, Code: ErrNotFound}, {Parent: 0, HasCode: true, Code: ErrParams}, {Parent: 2},
		}},
		{name: "multiple wraps", input: fmt.Errorf("both: %w | %w", first, second), causes: []Cause{
			{Parent: -1}, {Parent: 0, HasCode: true, Code: ErrParams}, {Parent: 1}, {Parent: 0, HasCode: true, Code: ErrNotFound},
		}},
		{name: "indexed argument order", input: fmt.Errorf("both: %[2]w | %[1]w", first, second), causes: []Cause{
			{Parent: -1}, {Parent: 0, HasCode: true, Code: ErrParams}, {Parent: 1}, {Parent: 0, HasCode: true, Code: ErrNotFound},
		}},
		{name: "overridden status", input: WrapStatus(first, StatusNotFound()), causes: []Cause{
			{Parent: -1, HasCode: true, Code: ErrNotFound}, {Parent: 0, HasCode: true, Code: ErrParams}, {Parent: 1},
		}},
		{name: "custom As virtual Item", input: &compatibilityAsError{item: first}, causes: []Cause{{Parent: -1}}},
		{name: "custom As typed nil", input: &compatibilityAsError{}, causes: []Cause{{Parent: -1}}},
		{name: "foreign Error", input: &compatibilityForeignError{status: StatusParams()}, causes: []Cause{
			{Parent: -1, HasCode: true, Code: ErrParams},
		}},
		{name: "zero business code", input: New(nil, NewStatus(0, 299, "custom")), causes: []Cause{
			{Parent: -1, HasCode: true, Code: 0},
		}},
		{name: "typed nil before later Item", input: errors.Join(typedNil, second), causes: []Cause{
			{Parent: -1}, {Parent: 0}, {Parent: 0, HasCode: true, Code: ErrNotFound},
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Inspect(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			want := Convert(tt.input)
			if got.HasError != (tt.input != nil) || got.Code != want.Code() ||
				got.HTTPCode != want.HTTPCode() || got.Message != want.Message() {
				t.Error("snapshot must preserve the legacy main state and input presence")
			}
			// 比较根全文是本 API 的契约断言，不按错误文案分类。
			wantText := ""
			if tt.input != nil {
				wantText = tt.input.Error()
			}
			if got.Error != wantText || len(got.Causes) != len(tt.causes) {
				t.Fatalf("text or tree changed: text=%q causes=%d, want %q / %d", got.Error, len(got.Causes), wantText, len(tt.causes))
			}
			for i, cause := range got.Causes {
				want := tt.causes[i]
				if cause.Parent != want.Parent || cause.HasCode != want.HasCode || cause.Code != want.Code {
					t.Errorf("cause %d state = (%d,%t,%d), want (%d,%t,%d)", i,
						cause.Parent, cause.HasCode, cause.Code, want.Parent, want.HasCode, want.Code)
				}
			}
			if tt.input == nil && got.Causes != nil {
				t.Error("nil input must have a nil cause list")
			}
			if strings.HasPrefix(tt.name, "custom As") && got.Causes[0].Frames != nil {
				t.Error("virtual As Item must not become an actual stack carrier")
			}
		})
	}
}

type inspectLink struct {
	cause error
	text  string
}

func (e *inspectLink) Error() string { return e.text }
func (e *inspectLink) Unwrap() error { return e.cause }

type inspectMany []error

func (e inspectMany) Error() string   { return "branches" }
func (e inspectMany) Unwrap() []error { return []error(e) }

type inspectBox struct {
	value any
	cause error
}

func (e inspectBox) Error() string { return "box" }
func (e inspectBox) Unwrap() error { return e.cause }

type inspectLoop []int

func (e inspectLoop) Error() string { return "loop" }
func (e inspectLoop) Unwrap() error { return e }

type inspectAsError struct {
	cause error
	as    func(any) bool
}

func (e *inspectAsError) Error() string      { return "custom As" }
func (e *inspectAsError) Unwrap() error      { return e.cause }
func (e *inspectAsError) As(target any) bool { return e.as(target) }

func TestInspectCustomAsOrder(t *testing.T) {
	t.Parallel()
	virtual := New(nil, StatusNotFound(), "virtual public")
	actual := New(nil, StatusParams())
	firstCalls, laterCalls := 0, 0
	later := &inspectAsError{cause: actual, as: func(any) bool { laterCalls++; return false }}
	root := &inspectAsError{cause: later, as: func(target any) bool {
		firstCalls++
		if item, ok := target.(**Item); ok {
			*item = virtual
			return true
		}
		return false
	}}
	got, err := Inspect(root)
	if err != nil {
		t.Fatal(err)
	}
	if firstCalls != 1 || laterCalls != 0 || got.Code != ErrNotFound || got.Message != "virtual public" {
		t.Error("only custom As calls up to the first match are permitted")
	}
	if len(got.Causes) != 3 || got.Causes[0].HasCode || got.Causes[1].HasCode || got.Causes[2].Code != ErrParams {
		t.Error("virtual main state must remain separate from the actual unwrap tree")
	}
	// AsType 复用同一个 target；即使返回 false 的 As 改写 target，也保持这一细节。
	root.as = func(target any) bool {
		if item, ok := target.(**Item); ok {
			*item = virtual
		}
		return false
	}
	later.as = func(target any) bool {
		item, ok := target.(**Item)
		return ok && *item == virtual
	}
	got, err = Inspect(root)
	if err != nil || got.Code != Convert(root).Code() || got.Code != ErrNotFound {
		t.Fatal("custom As target reuse must match the legacy conversion", err)
	}
	root.cause = io.EOF
	got, err = Inspect(root)
	if err != nil || got.Code != ErrInternalServer {
		t.Fatal("a false As result must not select its modified target", err)
	}
}

func assertInspectFailure(t *testing.T, input error) {
	t.Helper()
	got, err := Inspect(input)
	if err == nil || !reflect.DeepEqual(got, Diagnostic{}) {
		t.Fatalf("expected atomic diagnostic failure, got snapshot=%#v error=%v", got, err)
	}
}

type inspectNilPointer struct{}

func (*inspectNilPointer) Error() string { return "unknown nil pointer" }

type inspectNilMap map[string]int

func (inspectNilMap) Error() string { return "unknown nil map" }

type inspectNilFunc func() string

func (e inspectNilFunc) Error() string { return e() }

type inspectNilChan chan int

func (inspectNilChan) Error() string { return "unknown nil channel" }

func TestInspectUnknownTypedNil(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name  string
		input error
	}{
		{name: "pointer", input: (*inspectNilPointer)(nil)},
		{name: "map", input: inspectNilMap(nil)},
		{name: "slice", input: inspectMany(nil)},
		{name: "function", input: inspectNilFunc(nil)},
		{name: "channel", input: inspectNilChan(nil)},
		{name: "nested after matched Item", input: &Item{cause: inspectNilFunc(nil), status: StatusParams()}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertInspectFailure(t, tt.input)
		})
	}
}

func TestInspectStructureBoundaries(t *testing.T) {
	t.Parallel()
	self := &inspectLink{text: "self"}
	self.cause = self
	left, right := &inspectLink{text: "left"}, &inspectLink{text: "right"}
	left.cause, right.cause = right, left
	for _, tt := range []struct {
		name  string
		input error
	}{
		{name: "self cycle", input: self},
		{name: "indirect cycle", input: left},
		{name: "incomparable cycle", input: inspectLoop{1}},
		{name: "wide nil over limit", input: make(inspectMany, 256)},
		{name: "very wide nil", input: make(inspectMany, 65536)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertInspectFailure(t, tt.input)
		})
	}
	chain := io.EOF
	for range 255 {
		chain = &inspectLink{cause: chain, text: "node"}
	}
	got, err := Inspect(chain)
	if err != nil || len(got.Causes) != 256 {
		t.Fatalf("exactly 256 visits must succeed, got %d / %v", len(got.Causes), err)
	}
	assertInspectFailure(t, &inspectLink{cause: chain, text: "one more"})
	for _, input := range []error{make(inspectMany, 255), make(inspectMany, 0)} {
		got, err := Inspect(input)
		if err != nil || len(got.Causes) != 1 {
			t.Fatal("nil slots count toward the limit but do not create causes", err)
		}
	}
	// struct 的类型可比较，但内嵌接口的动态值不可比较；不能只查 Type.Comparable。
	boxed := inspectBox{value: []int{1}, cause: io.EOF}
	shared := &inspectLink{cause: boxed, text: "shared"}
	got, err = Inspect(inspectMany{nil, shared, nil, shared})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Causes) != 7 {
		t.Fatalf("shared branches must be repeated, got %d causes", len(got.Causes))
	}
	for i, parent := range []int{-1, 0, 1, 2, 0, 4, 5} {
		if got.Causes[i].Parent != parent {
			t.Errorf("cause %d parent = %d, want %d", i, got.Causes[i].Parent, parent)
		}
	}
}

type inspectCountedCycle struct {
	errorCalls  int
	unwrapCalls int
	asCalls     int
}

func (e *inspectCountedCycle) Error() string { e.errorCalls++; return "cycle" }
func (e *inspectCountedCycle) Unwrap() error { e.unwrapCalls++; return e }
func (e *inspectCountedCycle) As(any) bool   { e.asCalls++; return false }

func TestInspectCycleStopsBeforeText(t *testing.T) {
	t.Parallel()
	input := new(inspectCountedCycle)
	assertInspectFailure(t, input)
	if input.errorCalls != 0 || input.unwrapCalls != 1 || input.asCalls != 1 {
		t.Fatalf("cyclic input must not be re-read: Error=%d Unwrap=%d As=%d", input.errorCalls, input.unwrapCalls, input.asCalls)
	}
}

func TestInspectTextBoundaries(t *testing.T) {
	t.Parallel()
	plainLimit := 16384 - utf8.RuneCountInString(StatusInternalServer().Msg())
	for _, tt := range []struct {
		name    string
		input   error
		wantErr bool
	}{
		{name: "root exact", input: errors.New(strings.Repeat("界", plainLimit))},
		{name: "root over", input: errors.New(strings.Repeat("界", plainLimit+1)), wantErr: true},
		{name: "message exact", input: &Item{clientMsg: strings.Repeat("好", 16384), status: StatusParams()}},
		{name: "message over", input: &Item{clientMsg: strings.Repeat("好", 16385), status: StatusParams()}, wantErr: true},
		{name: "combined exact", input: &Item{msg: strings.Repeat("😊", 16383), clientMsg: "m", status: StatusParams()}},
		{name: "combined over", input: &Item{msg: strings.Repeat("😊", 16384), clientMsg: "m", status: StatusParams()}, wantErr: true},
		{name: "invalid UTF8 exact", input: &Item{msg: strings.Repeat("\xff", 16383), clientMsg: "m", status: StatusParams()}},
		{name: "invalid UTF8 over", input: &Item{msg: strings.Repeat("\xff", 16384), clientMsg: "m", status: StatusParams()}, wantErr: true},
		{name: "large root", input: errors.New(strings.Repeat("a", 1<<20)), wantErr: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.wantErr {
				assertInspectFailure(t, tt.input)
				return
			}
			got, err := Inspect(tt.input)
			if err != nil || got.Error != tt.input.Error() || got.Message != Convert(tt.input).Message() {
				t.Fatal("text at the rune limit must remain complete and byte-for-byte unchanged", err)
			}
		})
	}
}

func TestInspectFrameTextBoundary(t *testing.T) {
	t.Parallel()
	origin := New(nil, StatusParams())
	parsed := runtime.CallersFrames(origin.stack)
	frameRunes := 0
	for {
		frame, more := parsed.Next()
		frameRunes += utf8.RuneCountInString(frame.File) + utf8.RuneCountInString(frame.Function)
		if !more {
			break
		}
	}
	padding := 16384 - frameRunes - 1 // 再加公开文案 m 的一个 rune。
	if padding < 0 {
		t.Fatal("fixture stack already exceeds the text budget")
	}
	input := &Item{msg: strings.Repeat("界", padding), clientMsg: "m", status: StatusParams(), stack: origin.stack}
	got, err := Inspect(input)
	if err != nil || len(got.Causes[0].Frames) == 0 {
		t.Fatal("frame text and root text exactly at the combined limit must succeed", err)
	}
	input.msg += "x"
	assertInspectFailure(t, input)
	// 仅帧文本超限同样原子失败；不把部分帧作为完整结果。
	input.msg, input.clientMsg = "", "m"
	input.stack = slices.Repeat(origin.stack, 100)
	assertInspectFailure(t, input)
}

func TestInspectExistingStacks(t *testing.T) {
	t.Parallel()
	origin := New(io.EOF, StatusParams(), "origin")
	wrapped := Wrap(origin, "load")
	got, err := Inspect(fmt.Errorf("standard prefix: %w", wrapped))
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Causes) != 4 || got.Causes[0].Frames != nil || got.Causes[2].Frames != nil || got.Causes[3].Frames != nil {
		t.Fatal("only the first actual carrier along this path must retain the shared stack")
	}
	wantFrame, _ := runtime.CallersFrames(origin.stack).Next()
	if len(got.Causes[1].Frames) == 0 || !reflect.DeepEqual(got.Causes[1].Frames[0], wantFrame) {
		t.Fatal("snapshot must parse the existing origin stack, without capturing Inspect")
	}
	got, err = Inspect(errors.Join(wrapped, wrapped))
	if err != nil || len(got.Causes) != 7 {
		t.Fatal("shared objects in sibling branches are valid", err)
	}
	if len(got.Causes[1].Frames) == 0 || !reflect.DeepEqual(got.Causes[1].Frames, got.Causes[4].Frames) ||
		got.Causes[2].Frames != nil || got.Causes[5].Frames != nil || &got.Causes[1].Frames[0] == &got.Causes[4].Frames[0] {
		t.Error("sibling branches must each own their stack data without duplicating descendants")
	}
	other := New(nil, StatusNotFound())
	combined := New(errors.Join(origin, other), StatusInternalServer())
	got, err = Inspect(combined)
	if err != nil || len(got.Causes) != 5 || len(got.Causes[0].Frames) == 0 ||
		got.Causes[2].Frames != nil || len(got.Causes[4].Frames) == 0 {
		t.Fatal("different branch origins must remain associated with their actual carriers", err)
	}
	for _, input := range []error{io.EOF, Convert(io.EOF), &Item{}, &Item{stack: []uintptr{}}, &Item{stack: []uintptr{0, 1}}} {
		got, err := Inspect(input)
		if err != nil {
			t.Fatal(err)
		}
		for _, cause := range got.Causes {
			if cause.Frames != nil {
				t.Error("reading must not add frames when no valid existing stack is available")
			}
		}
	}
}

func TestInspectOwnership(t *testing.T) {
	t.Parallel()
	origin := New(io.EOF, StatusParams(), "origin")
	input := fmt.Errorf("outer: %w", Wrap(origin, "load"))
	pcs := slices.Clone(origin.stack)
	first, err := Inspect(input)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Inspect(input)
	if err != nil || !reflect.DeepEqual(first, second) || len(first.Causes[1].Frames) == 0 {
		t.Fatal("repeated reads must return equivalent independent data", err)
	}
	first.Error, first.Message = "changed", "changed"
	first.Causes[1].Code = ErrNotFound
	first.Causes[1].Frames[0].File = "changed"
	first.Causes[1].Frames[0].PC = 0
	first.Causes[1].Frames = append(first.Causes[1].Frames, runtime.Frame{})
	first.Causes[0] = Cause{Parent: 99}
	third, err := Inspect(input)
	if err != nil || !reflect.DeepEqual(second, third) || !slices.Equal(origin.stack, pcs) ||
		origin.Code() != ErrParams || !errors.Is(input, io.EOF) {
		t.Fatal("returned mutations must not affect another snapshot or the original error", err)
	}
}

func TestInspectConcurrentOwnership(t *testing.T) {
	t.Parallel()
	origin := New(io.EOF, StatusParams(), "origin")
	input := errors.Join(Wrap(origin, "left"), Wrap(origin, "right"))
	want, err := Inspect(input)
	if err != nil {
		t.Fatal(err)
	}
	var workers sync.WaitGroup
	for range 32 {
		workers.Go(func() {
			for range 16 {
				got, err := Inspect(input)
				if err != nil || !reflect.DeepEqual(got, want) {
					t.Error("concurrent read changed the snapshot", err)
					return
				}
				got.Causes[1].Code = ErrNotFound
				got.Causes[1].Frames[0].File = "caller mutation"
				got.Causes[4].Frames[0].Function = "caller mutation"
			}
		})
	}
	workers.Wait()
	got, err := Inspect(input)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatal("parallel mutations must remain local to each returned snapshot", err)
	}
}
