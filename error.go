package goerr

import (
	"errors"
	"fmt"
	"runtime"
)

// Error 带业务状态码的错误接口。
// 注意：Unwrap() 由 *Item 实现以支持 errors.Is/errors.As 链式解包，
// 但不放入接口约束——其他实现者不一定有 cause 链。
type Error interface {
	error
	StatusInfo() *Status
	Code() Code
	Message() string
	HTTPCode() int
}

// Item 是 Error 的默认实现，不可变——所有字段在构造后不再修改。
type Item struct {
	// msg 供 Error() 使用：完整描述，常含 cause，用于日志。
	msg string
	// clientMsg 供 Message() 优先返回：对外短文案；空则回退 StatusInfo().Msg()。
	clientMsg string
	status    *Status
	cause     error // 原始错误，支持 Unwrap
	stack     []uintptr
}

// Error 实现 error 接口，返回完整描述（常含 cause），用于日志与 fmt 打印。
func (i *Item) Error() string { return i.msg }

func (i *Item) StatusInfo() *Status {
	if i.status == nil {
		return StatusOK()
	}
	return i.status
}

func (i *Item) Code() Code {
	return i.StatusInfo().Code()
}

// Message 返回面向调用方/客户端的短文案（不含底层错误细节时由构造逻辑保证）。
// 若构造时写入了 clientMsg 则优先返回；否则使用当前 Status 的默认文案。
// 与 Error() 的分工见 README「Message() 与 Error()：团队约定」；网关脱敏见「网关：对 Message() 做截断与脱敏」。
func (i *Item) Message() string {
	if i.clientMsg != "" {
		return i.clientMsg
	}
	return i.StatusInfo().Msg()
}

func (i *Item) HTTPCode() int {
	return i.StatusInfo().HTTPCode()
}

// Unwrap 返回底层原始错误，使 errors.Is/errors.As 能穿透 Item。
func (i *Item) Unwrap() error { return i.cause }

// Format 支持 %+v 输出调用栈信息（配合 zap 等日志库的 Verbose 模式）。
func (i *Item) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprint(s, i.msg)
			for _, pc := range i.stack {
				fn := runtime.FuncForPC(pc)
				if fn == nil {
					continue
				}
				file, line := fn.FileLine(pc)
				fmt.Fprintf(s, "\n\t%s:%d %s", file, line, fn.Name())
			}
			// 递归打印 cause 链
			if i.cause != nil {
				fmt.Fprintf(s, "\nCaused by: %+v", i.cause)
			}
			return
		}
		fmt.Fprint(s, i.msg)
	case 's':
		fmt.Fprint(s, i.msg)
	case 'q':
		fmt.Fprintf(s, "%q", i.msg)
	}
}

// callers 抓取调用栈。skip 为传给 runtime.Callers 的帧跳过数（通常为 3：Callers + callers + 公开 API）。
func callers(skip int) []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(skip, pcs[:])
	return pcs[:n]
}

// New 创建带业务状态码的错误。
//
//	err    — 原始错误（可为 nil，此时构造一个纯状态错误）。
//	status — *Status，一般为 StatusXxx() 的返回值。
//	msgs   — 可选；非空则作为对外文案（Message()），并参与 Error()；见 README「Message() 与 Error()」约定。
//
// 示例:
//
//	goerr.New(err, goerr.StatusInternalServer())
//	goerr.New(err, goerr.StatusParams(), "user_id is required")
//	goerr.New(nil, goerr.StatusNotFound(), "order not found")
func New(err error, status *Status, msgs ...string) *Item {
	resolvedStatus := status
	if resolvedStatus == nil {
		resolvedStatus = StatusOK()
	}

	msg := resolvedStatus.Msg()
	if len(msgs) > 0 && msgs[0] != "" {
		msg = msgs[0]
	}

	if err != nil {
		return &Item{
			msg:       msg + ": " + err.Error(),
			clientMsg: msg,
			status:    resolvedStatus,
			cause:     err,
			stack:     callers(3),
		}
	}
	return &Item{
		msg:       msg,
		clientMsg: msg,
		status:    resolvedStatus,
		cause:     nil,
		stack:     callers(3),
	}
}

// NewFn 接受 func() *Status（即 StatusXxx 函数引用本身）
func NewFn(err error, statusFn func() *Status, msgs ...string) *Item {
	return New(err, statusFn(), msgs...)
}

// Newf 创建带格式化消息的状态错误（无底层 cause）。
// status 为 nil 时使用 [StatusOK]。
func Newf(status *Status, format string, args ...any) *Item {
	resolvedStatus := status
	if resolvedStatus == nil {
		resolvedStatus = StatusOK()
	}
	out := fmt.Sprintf(format, args...)
	return &Item{
		msg:       out,
		clientMsg: out,
		status:    resolvedStatus,
		stack:     callers(3),
	}
}

// Wrap 为已有错误附加上下文消息。
// 如果 err 已是 *Item，会保留其 status 和 cause 链；否则创建新的无状态 Item。
// 始终返回新的 *Item，不修改原始错误——并发安全。
func Wrap(err error, msg string) *Item {
	if err == nil {
		return nil
	}

	if e, ok := errors.AsType[*Item](err); ok {
		return &Item{
			msg:       msg + ": " + e.msg,
			clientMsg: e.Message(),
			status:    e.status,
			cause:     err, // 保留完整 cause 链
			stack:     callers(3),
		}
	}
	return &Item{
		msg:       msg + ": " + err.Error(),
		clientMsg: msg,
		cause:     err,
		stack:     callers(3),
	}
}

// Wrapf 同 Wrap，支持格式化消息。
// 注意：不能直接调用 Wrap，否则 callers() 栈帧会多一层，导致调用位置不准。
func Wrapf(err error, format string, args ...any) *Item {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)
	if e, ok := errors.AsType[*Item](err); ok {
		return &Item{
			msg:       msg + ": " + e.msg,
			clientMsg: e.Message(),
			status:    e.status,
			cause:     err,
			stack:     callers(3),
		}
	}
	return &Item{
		msg:       msg + ": " + err.Error(),
		clientMsg: msg,
		cause:     err,
		stack:     callers(3),
	}
}

// WrapStatus 为已有错误附加业务状态码。
// 适用于在中间件或 handler 层为底层错误补充 HTTP 状态信息。
// status 为 nil 时使用 [StatusInternalServer]。
// 若需在运行时解析 *Status / func() *Status 且避免类型错误，请使用 [WrapStatusFrom]。
func WrapStatus(err error, status *Status) *Item {
	if err == nil {
		return nil
	}
	resolvedStatus := status
	if resolvedStatus == nil {
		resolvedStatus = StatusInternalServer()
	}
	cm := resolvedStatus.Msg()
	return &Item{
		msg:       cm + ": " + err.Error(),
		clientMsg: cm,
		status:    resolvedStatus,
		cause:     err,
		stack:     callers(3),
	}
}

// WithStack 为普通 error 附加调用栈。若已是 *Item 则直接返回。
func WithStack(err error) *Item {
	if err == nil {
		return nil
	}
	if e, ok := errors.AsType[*Item](err); ok {
		return e
	}
	return &Item{
		msg:       err.Error(),
		clientMsg: err.Error(),
		cause:     err,
		stack:     callers(3),
	}
}

// --- 标准库 errors 包的便捷透传 ---

func Err(e string) error {
	return errors.New(e)
}

// Is 等价于 errors.Is。
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// AsItem 从错误链中提取 *Item（Go 1.26 errors.AsType 泛型写法）。
func AsItem(err error) (*Item, bool) {
	return errors.AsType[*Item](err)
}
