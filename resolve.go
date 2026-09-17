package goerr

import (
	"errors"
	"fmt"
)

// ResolveStatus 将 *Status、func() *Status 或 nil 解析为具体 *Status，失败时返回错误，不 panic。
// nil 或 nil *Status 时使用 fallback。
// func() *Status 为 nil 时返回错误；若工厂返回 nil *Status，则使用 fallback。
//
// Deprecated: 用 any 承载 *Status / func() *Status 放弃了编译期类型检查。请直接传 *Status 调用 [New] / [Newf] / [WrapStatus]。
func ResolveStatus(status any, fallback *Status) (*Status, error) {
	switch s := status.(type) {
	case nil:
		return fallback, nil
	case *Status:
		if s == nil {
			return fallback, nil
		}
		return s, nil
	case func() *Status:
		if s == nil {
			return nil, errors.New("goerr: status factory is nil")
		}
		r := s()
		if r == nil {
			return fallback, nil
		}
		return r, nil
	default:
		return nil, fmt.Errorf("goerr: unsupported status type %T", status)
	}
}

// NewFrom 等价于 [New]，但 status 可为 *Status、func() *Status 或 nil，类型错误时返回 error。
// 状态缺失时与 [New] 一致，按 [StatusInternalServer] 兜底。
//
// Deprecated: 请直接调用 [New]，status 传 StatusXxx() 的返回值。
func NewFrom(err error, status any, msgs ...string) (*Item, error) {
	st, rerr := ResolveStatus(status, StatusInternalServer())
	if rerr != nil {
		return nil, rerr
	}
	return newItem(err, st, msgs, skipViaBuilder), nil
}

// NewfFrom 等价于 [Newf]，status 解析规则同 [ResolveStatus]（fallback 为 [StatusInternalServer]）。
//
// Deprecated: 请直接调用 [Newf]，status 传 StatusXxx() 的返回值。
func NewfFrom(status any, format string, args ...any) (*Item, error) {
	st, err := ResolveStatus(status, StatusInternalServer())
	if err != nil {
		return nil, err
	}
	return newItem(nil, st, []string{fmt.Sprintf(format, args...)}, skipViaBuilder), nil
}

// WrapStatusFrom 等价于 [WrapStatus]，status 解析失败时由第二个返回值报告；
// 状态缺失时 fallback 为 [StatusInternalServer]。
// 第一个返回值是 error 接口而非 *Item：err 为 nil 时它是真正的 nil，可直接向上返回。
//
// Deprecated: 请直接调用 [WrapStatus]，status 传 StatusXxx() 的返回值。
func WrapStatusFrom(err error, status any) (error, error) {
	if err == nil {
		return nil, nil
	}
	st, rerr := ResolveStatus(status, StatusInternalServer())
	if rerr != nil {
		return nil, rerr
	}
	return wrapStatusItem(err, st, skipViaBuilder), nil
}
