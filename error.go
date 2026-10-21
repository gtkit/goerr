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

// Item 是 Error 的默认实现，不可变——所有字段在构造后不再修改，可跨 goroutine 共享。
type Item struct {
	// msg 供 Error() 使用：完整描述，已含整条 cause 链的文本，用于日志。
	msg string
	// clientMsg 供 Message() 使用：对外短文案；空则回退 StatusInfo().Msg()。
	clientMsg string
	status    *Status
	cause     error // 原始错误，支持 Unwrap
	// stack 整条链共享同一份——只在链上第一个 Item 处捕获。
	stack []uintptr
}

// 下列方法对 nil receiver 安全：*Item 可能以类型化 nil 形式出现
// （典型是忽略了 AsItem 的第二个返回值），此时按「状态未知」兜底而非 panic。

// Error 实现 error 接口，返回完整描述（含整条 cause 链），用于日志与 fmt 打印。
func (i *Item) Error() string {
	if i == nil {
		return ""
	}
	return i.msg
}

// StatusInfo 返回错误状态。状态缺失时一律兜底为 [StatusInternalServer]，
// 避免未知状态被统一响应层当成成功。
func (i *Item) StatusInfo() *Status {
	if i == nil || i.status == nil {
		return StatusInternalServer()
	}
	return i.status
}

func (i *Item) Code() Code { return i.StatusInfo().Code() }

func (i *Item) HTTPCode() int { return i.StatusInfo().HTTPCode() }

// Message 返回面向调用方/客户端的短文案。
// 若构造时写入了 clientMsg 则优先返回；否则使用当前 Status 的默认文案。
// 与 Error() 的分工见 README「Message() 与 Error()：团队约定」；网关脱敏见「网关：对 Message() 做截断与脱敏」。
func (i *Item) Message() string {
	if i == nil {
		return ""
	}
	if i.clientMsg != "" {
		return i.clientMsg
	}
	return i.StatusInfo().Msg()
}

// Unwrap 返回底层原始错误，使 errors.Is/errors.As 能穿透 Item。
func (i *Item) Unwrap() error {
	if i == nil {
		return nil
	}
	return i.cause
}

// Format 支持 %+v 输出调用栈（配合 zap 等日志库的 Verbose 模式）。
// msg 本身已含整条 cause 链的文本、stack 整条链共享一份，因此不再递归打印 cause，
// 否则 n 层 Wrap 会让日志量按 n² 膨胀、同一份栈重复 n 遍。
// State 写入报错后立即停止；Format 没有错误返回值，不能代替外层 fmt 写入的错误处理。
func (i *Item) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if _, err := fmt.Fprint(s, i.Error()); err != nil {
				return
			}
			if i != nil {
				frames := runtime.CallersFrames(i.stack)
				for {
					f, more := frames.Next()
					if f.Function != "" {
						if _, err := fmt.Fprintf(s, "\n\t%s:%d %s", f.File, f.Line, f.Function); err != nil {
							return
						}
					}
					if !more {
						break
					}
				}
			}
			return
		}
		fallthrough
	case 's':
		if _, err := fmt.Fprint(s, i.Error()); err != nil {
			return
		}
	case 'q':
		if _, err := fmt.Fprintf(s, "%q", i.Error()); err != nil {
			return
		}
	default:
		if _, err := fmt.Fprintf(s, "%%!%c(*goerr.Item=%s)", verb, i.Error()); err != nil {
			return
		}
	}
}

// runtime.Callers 需跳过的帧数，按各自的调用链路数出。
// 每个公开 API 必须直接调用内部构造器，不得互相转调，
// 否则会多出一层栈帧，栈首帧落在本包内部而不是调用方代码。
const (
	skipViaBuilder = 5 // 公开 API → 内部构造器 → stackFor → callers
	skipDirect     = 3 // 公开 API → callers
)

// callers 抓取调用栈，只分配实际用到的帧数。
func callers(skip int) []uintptr {
	var pcs [32]uintptr
	n := runtime.Callers(skip, pcs[:])
	out := make([]uintptr, n)
	copy(out, pcs[:n])
	return out
}

// stackFor 优先复用 cause 链上已有的栈：一条链只需一份，
// 且最内层那份最接近错误起源。链上没有才在当前位置捕获。
func stackFor(cause error, skip int) []uintptr {
	// e 可能是类型化 nil（调用方传了 var it *Item 之类），必须在取字段前挡住。
	if e, ok := errors.AsType[*Item](cause); ok && e != nil && e.stack != nil {
		return e.stack
	}
	return callers(skip)
}

// orInternal 统一 nil status 的兜底方向：状态缺失一律按内部服务器错误处理，
// 避免失败被统一响应层判成成功。
func orInternal(status *Status) *Status {
	if status == nil {
		return StatusInternalServer()
	}
	return status
}

// newItem 是 New 系列的唯一构造出口。msgs 非空且首项非空时作为对外文案，
// 否则用 status 默认文案；有 cause 时 Error() 额外拼接底层错误全文。
func newItem(cause error, status *Status, msgs []string, skip int) *Item {
	st := orInternal(status)
	cm := st.Msg()
	if len(msgs) > 0 && msgs[0] != "" {
		cm = msgs[0]
	}
	msg := cm
	if cause != nil {
		msg = cm + ": " + cause.Error()
	}
	return &Item{msg: msg, clientMsg: cm, status: st, cause: cause, stack: stackFor(cause, skip)}
}

// New 创建带业务状态码的错误。
//
//	err    — 原始错误（可为 nil，此时构造一个纯状态错误）。
//	status — *Status，一般为 StatusXxx() 的返回值；nil 时按 [StatusInternalServer] 兜底。
//	msgs   — 可选；仅取首项，非空则作为对外文案（Message()），并参与 Error()。
//
// 示例:
//
//	goerr.New(err, goerr.StatusInternalServer())
//	goerr.New(err, goerr.StatusParams(), "user_id is required")
//	goerr.New(nil, goerr.StatusNotFound(), "order not found")
func New(err error, status *Status, msgs ...string) *Item {
	return newItem(err, status, msgs, skipViaBuilder)
}

// NewFn 接受 func() *Status（即 StatusXxx 函数引用本身），其余同 [New]。
//
// Deprecated: 与 New(err, StatusXxx(), msgs...) 只差一对括号，请直接调用 [New]。
func NewFn(err error, statusFn func() *Status, msgs ...string) *Item {
	return newItem(err, statusFn(), msgs, skipViaBuilder)
}

// Newf 创建带格式化消息的状态错误（无底层 cause）。
// status 为 nil 时按 [StatusInternalServer] 兜底。
func Newf(status *Status, format string, args ...any) *Item {
	return newItem(nil, status, []string{fmt.Sprintf(format, args...)}, skipViaBuilder)
}

// wrapItem 是 Wrap 系列的唯一构造出口：内层已是 *Item 时保留其状态与对外文案，
// 否则按内部服务器错误兜底，避免普通 error 被统一响应层误判为成功。
func wrapItem(err error, msg string, skip int) *Item {
	st := StatusInternalServer()
	cm := st.Msg()
	if e, ok := errors.AsType[*Item](err); ok {
		st, cm = e.StatusInfo(), e.Message()
	}
	return &Item{
		msg:       msg + ": " + err.Error(),
		clientMsg: cm,
		status:    st,
		cause:     err,
		stack:     stackFor(err, skip),
	}
}

// Wrap 为已有错误附加上下文消息，err 为 nil 时返回 nil（接口值真 nil，不是类型化 nil）。
// 始终返回新的错误，不修改原始错误——并发安全。
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return wrapItem(err, msg, skipViaBuilder)
}

// Wrapf 同 [Wrap]，支持格式化消息。
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return wrapItem(err, fmt.Sprintf(format, args...), skipViaBuilder)
}

// wrapStatusItem 是 WrapStatus 系列的唯一构造出口。
func wrapStatusItem(err error, status *Status, skip int) *Item {
	st := orInternal(status)
	cm := st.Msg()
	return &Item{
		msg:       cm + ": " + err.Error(),
		clientMsg: cm,
		status:    st,
		cause:     err,
		stack:     stackFor(err, skip),
	}
}

// WrapStatus 为已有错误附加业务状态码，适用于在中间件或 handler 层为底层错误补充 HTTP 状态信息。
// err 为 nil 时返回 nil；status 为 nil 时使用 [StatusInternalServer]。
// 若需在运行时解析 *Status / func() *Status 且避免类型错误，请使用 [WrapStatusFrom]。
func WrapStatus(err error, status *Status) error {
	if err == nil {
		return nil
	}
	return wrapStatusItem(err, status, skipViaBuilder)
}

// WithStack 为 err 补充调用栈，保留原始全文、状态和 cause。err 为 nil 时返回 nil。
// 首个匹配的 *Item 已有非空栈或为 typed-nil 时原样返回 err，保留外层包装身份；
// 非 nil Item 缺栈时返回新包装，在调用方位置抓栈，不修改原 Item。
// 未找到 Item 的普通错误仍按内部服务器错误兜底。
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	item, ok := errors.AsType[*Item](err)
	if ok && (item == nil || len(item.stack) > 0) {
		return err
	}
	st := StatusInternalServer()
	cm := st.Msg()
	if item != nil {
		st, cm = item.StatusInfo(), item.Message()
	}
	return &Item{
		msg:       err.Error(),
		clientMsg: cm,
		status:    st,
		cause:     err,
		stack:     callers(skipDirect), // 首个匹配的 Item 缺栈，或链上没有 Item
	}
}

// --- 从任意 error 取业务信息 ---

// Is 等价于 errors.Is。
//
// Deprecated: 请直接使用标准库 errors.Is。
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// AsItem 从错误链中提取 *Item。只在「链上没有 *Item 时要走另一条分支」的场景用它；
// 统一响应层请用 [Convert]，免去 ok 判断与兜底代码。
func AsItem(err error) (*Item, bool) {
	return errors.AsType[*Item](err)
}

// okItem 是 Convert(nil) 的返回值。没有错误就是成功，所有调用方共享这一个只读实例。
var okItem = &Item{msg: StatusOK().Msg(), clientMsg: StatusOK().Msg(), status: StatusOK()}

// Convert 把任意 error 转成 *Item，供统一响应层一次拿齐 Code / Message / HTTPCode，永不返回 nil：
//
//   - err 为 nil        → 成功状态（[StatusOK]），成功与失败走同一条响应路径；
//   - 链上有 *Item      → 返回该 *Item，自定义对外文案原样保留；
//   - 其他普通 error    → 按 [StatusInternalServer] 包装，对外文案为状态默认文案，Error() 保留原文。
//
// 兜底包装不抓调用栈：响应层的位置对排障没有价值，起源栈在链上的 *Item 里。
// 日志请打印原始 err 而非 Convert 的结果，链上挖出的 *Item 不含外层 fmt.Errorf 加的前缀。
//
//	e := goerr.Convert(err)
//	c.JSON(e.HTTPCode(), gin.H{"code": e.Code(), "message": e.Message(), "data": data})
func Convert(err error) *Item {
	if err == nil {
		return okItem
	}
	if e, ok := errors.AsType[*Item](err); ok && e != nil {
		return e
	}
	st := StatusInternalServer()
	return &Item{msg: err.Error(), clientMsg: st.Msg(), status: st, cause: err}
}

// HasCode 判断错误链上的 *Item 是否携带指定业务错误码。
// 相比 errors.Is，它按业务码而非错误实例比较——同一失败语义在不同位置构造出的
// Item 是不同实例，errors.Is 判不出相等。
func HasCode(err error, code Code) bool {
	e, ok := errors.AsType[*Item](err)
	return ok && e.Code() == code
}
