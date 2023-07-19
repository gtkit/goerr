// @Author xiaozhaofu 2022/11/1 16:05:00
package goerr

import (
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

// Error an error with caller stack information
type Error interface {
	error
	Code() ErrCode
}

type item struct {
	msg   string
	code  ErrCode
	stack []uintptr
}

func (i *item) Error() string {
	return i.msg
}

func (i *item) Code() ErrCode {
	return i.code
}

// Format used by go.uber.org/zap in Verbose
func (i *item) Format(s fmt.State, verb rune) {
	io.WriteString(s, i.msg)
	io.WriteString(s, "\n")

	for _, pc := range i.stack {
		fmt.Fprintf(s, "%+v\n", errors.Frame(pc))
	}
}

// New create a new error
func New(err error, code ErrCode, msg string) Error {
	if err != nil {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), code: code, stack: callers()}
	}
	return &item{msg: fmt.Sprintf("%s;", msg), code: code, stack: callers()}
}

// Errorf create a new error
func Errorf(format string, args ...interface{}) Error {
	return &item{msg: fmt.Sprintf(format, args...), stack: callers()}
}

// Wrap with some extra message into err
func Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}

	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// Wrapf with some extra message into err
func Wrapf(err error, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)

	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// WithStack add caller stack information
func WithStack(err error) Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*item); ok {
		return e
	}

	return &item{msg: err.Error(), stack: callers()}
}

// 返回没有堆栈的自定义错误
func Custom(msg string) error {
	return errors.New(msg)
}
func Err(msg string) error {
	return errors.New(msg)
}
func Is(err, target error) bool {
	return errors.Is(err, target)
}

func WithMessage(err error, msg string) error {
	return errors.WithMessage(err, msg)
}
