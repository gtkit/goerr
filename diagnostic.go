package goerr

import (
	"errors"
	"reflect"
	"runtime"
	"slices"
)

// Diagnostic 是供内部排障使用的独立快照，不是客户端响应或序列化协议。
// Causes 及其中的 Frames 由本次 Inspect 独占；普通值复制仍会共享切片，
// 同一快照的并发修改需要调用方同步。不同 Inspect 调用互不共享可写存储。
type Diagnostic struct {
	_ struct{} // 守卫：强制跨包调用方带键初始化。
	// HasError 表示输入接口非 nil，不能用业务成功码代替此判断。
	HasError bool
	// Code 遵守 Convert 对首个匹配 Item 的业务码选择。
	Code Code
	// HTTPCode 是该主状态的 HTTP 状态码。
	HTTPCode int
	// Message 是该主状态的公开文案，与根错误全文分开保存。
	Message string
	// Error 是根错误的完整原文，可能含内部信息，未经处理不得返回客户端。
	Error string
	// Causes 是实际 Unwrap 树的前序列表，不包含自定义 As 提供的虚拟节点。
	Causes []Cause
}

// Cause 表示实际错误树中的一个节点，不持有原始 error 或可写 PC 存储。
// 并发修改同一个 Cause 或其 Frames 需要调用方同步。
type Cause struct {
	_ struct{} // 守卫：强制跨包调用方带键初始化。
	// Parent 是直接父节点在 Diagnostic.Causes 中的索引，根为 -1。
	Parent int
	// HasCode 表示实际非 nil 节点实现了 Error 接口；否则 Code 为零。
	HasCode bool
	// Code 是该实际节点的业务码，不改变 Diagnostic 的主状态。
	Code Code
	// Frames 只包含该路径首次携带的已有栈，不捕获读取位置的栈。
	// 切片独立；runtime.Frame.Func 仅引用运行时只读元数据。
	Frames []runtime.Frame
}

const (
	inspectVisitLimit = 256
	inspectTextLimit  = 16384
)

// Inspect 返回 err 的完整诊断快照；失败返回 Diagnostic{} 和非 nil 诊断错误，
// 不修改或替换原业务错误。诊断错误的类型和文字不作为分类契约。
// nil 输入得到 HasError=false、成功状态和空原因列表；非 nil 成功 Item 仍有 HasError=true。
// 已知 typed-nil *Item 可安全读取，其他 typed-nil 节点在调用其方法前使读取失败。
//
// 最多访问 256 个节点/子槽（多值 Unwrap 的 nil 槽也计数），最多保存 16384 个 rune，
// 文本累计 Error、Message 和帧的 File/Function；循环或超限均失败，不返回截断结果。
// 同一对象在不同分支分别列出，沿同一路径共享的同一 PC 切片只保存一次栈。
//
// 输入须在读取期间保持稳定，外部 Error/Unwrap/As/Code 方法须有限、无阻塞且不 panic；
// 预算不隔离这些方法内部的阻塞、分配或副作用。在输入支持并发读取时可并发调用 Inspect。
// 该函数不输出日志、不编码、不抓新栈；调用方负责输出、字段选择及脱敏。
func Inspect(err error) (Diagnostic, error) {
	if err == nil {
		item := Convert(nil)
		return Diagnostic{Code: item.Code(), HTTPCode: item.HTTPCode(), Message: item.Message()}, nil
	}
	reader := diagnosticReader{textLeft: inspectTextLimit}
	if readErr := reader.walk(err); readErr != nil {
		return Diagnostic{}, readErr
	}
	status := StatusInternalServer()
	diagnostic := Diagnostic{
		HasError: true, Code: status.Code(), HTTPCode: status.HTTPCode(), Message: status.Msg(),
		Error: err.Error(), Causes: reader.causes,
	}
	if reader.matched && reader.item != nil {
		diagnostic.Code = reader.item.Code()
		diagnostic.HTTPCode = reader.item.HTTPCode()
		diagnostic.Message = reader.item.Message()
	}
	if readErr := reader.reserve(diagnostic.Error, diagnostic.Message); readErr != nil {
		return Diagnostic{}, readErr
	}
	return diagnostic, nil
}

// 读取状态仅存活于一次 Inspect；不接入旧构造/转换路径。
type diagnosticReader struct {
	causes   []Cause
	item     *Item
	matched  bool
	textLeft int
}

type diagnosticStep struct {
	single   error
	children []error
	next     int
	parent   int
	key      error // 仅可比较节点进入当前路径集合；nil 表示无需删除。
	pcs      []uintptr
}

func (r *diagnosticReader) walk(root error) error {
	path := []diagnosticStep{{single: root, parent: -1}}
	ancestors := make(map[error]struct{})
	visits := 0
	for len(path) > 0 {
		top := &path[len(path)-1]
		node, parent := top.single, top.parent
		if node != nil {
			top.single = nil
		} else if top.next < len(top.children) {
			node = top.children[top.next]
			top.next++
		} else {
			delete(ancestors, top.key)
			path = path[:len(path)-1]
			continue
		}
		visits++
		if visits > inspectVisitLimit {
			return errors.New("goerr: inspect visit limit exceeded")
		}
		if node == nil {
			continue
		}
		item, isItem := node.(*Item)
		value := reflect.ValueOf(node)
		if !isItem {
			switch value.Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
				if value.IsNil() {
					return errors.New("goerr: inspect unknown typed-nil error")
				}
			}
		}
		step := diagnosticStep{parent: len(r.causes)}
		if value.Comparable() {
			if _, exists := ancestors[node]; exists {
				return errors.New("goerr: inspect cyclic error tree")
			}
			ancestors[node] = struct{}{}
			step.key = node
		}
		// 按标准库 AsType 的顺序匹配，但不再次无界遍历整棵树。
		if !r.matched {
			if isItem {
				r.item, r.matched = item, true
			} else if custom, ok := node.(interface{ As(any) bool }); ok {
				r.matched = custom.As(&r.item)
			}
		}
		cause := Cause{Parent: parent}
		if !isItem || item != nil {
			if coded, ok := node.(Error); ok {
				cause.HasCode, cause.Code = true, coded.Code()
			}
		}
		if item != nil && len(item.stack) > 0 {
			step.pcs = item.stack
			shared := slices.ContainsFunc(path, func(ancestor diagnosticStep) bool {
				return len(ancestor.pcs) == len(step.pcs) && &ancestor.pcs[0] == &step.pcs[0]
			})
			if !shared {
				frames, readErr := r.frames(item.stack)
				if readErr != nil {
					return readErr
				}
				cause.Frames = frames
			}
		}
		r.causes = append(r.causes, cause)
		switch unwrapped := node.(type) {
		case interface{ Unwrap() error }:
			step.single = unwrapped.Unwrap()
		case interface{ Unwrap() []error }:
			step.children = unwrapped.Unwrap()
		}
		path = append(path, step)
	}
	return nil
}

func (r *diagnosticReader) frames(pcs []uintptr) ([]runtime.Frame, error) {
	var result []runtime.Frame
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if frame.PC != 0 {
			if readErr := r.reserve(frame.File, frame.Function); readErr != nil {
				return nil, readErr
			}
			result = append(result, frame)
		}
		if !more {
			return result, nil
		}
	}
}

func (r *diagnosticReader) reserve(texts ...string) error {
	for _, text := range texts {
		for range text {
			if r.textLeft == 0 {
				return errors.New("goerr: inspect text limit exceeded")
			}
			r.textLeft--
		}
	}
	return nil
}
