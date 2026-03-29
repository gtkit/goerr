package goerr

import (
	"errors"
	"fmt"
)

// ResolveStatus 将 *Status、func() *Status 或 nil 解析为具体 *Status，失败时返回错误，不 panic。
// nil 或 nil *Status 时使用 fallback。
// func() *Status 为 nil 时返回错误；若工厂返回 nil *Status，则使用 fallback。
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

// NewFrom 等价于 [New]，但 status 可为 *Status、func() *Status 或 nil（使用 StatusOK），类型错误时返回 error。
func NewFrom(err error, status any, msgs ...string) (*Item, error) {
	st, rerr := ResolveStatus(status, StatusOK())
	if rerr != nil {
		return nil, rerr
	}
	return New(err, st, msgs...), nil
}

// NewfFrom 等价于 [Newf]，status 解析规则同 [ResolveStatus]（fallback 为 StatusOK）。
func NewfFrom(status any, format string, args ...any) (*Item, error) {
	st, err := ResolveStatus(status, StatusOK())
	if err != nil {
		return nil, err
	}
	return Newf(st, format, args...), nil
}

// WrapStatusFrom 等价于 [WrapStatus]，status 解析失败时返回 error；nil 或 nil *Status 时 fallback 为 [StatusInternalServer]。
func WrapStatusFrom(err error, status any) (*Item, error) {
	if err == nil {
		return nil, nil
	}
	st, rerr := ResolveStatus(status, StatusInternalServer())
	if rerr != nil {
		return nil, rerr
	}
	return WrapStatus(err, st), nil
}
