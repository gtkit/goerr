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
	Status() ErrStatuser
	ErrCode() int
	ErrMsg() string
	HttpCode() int
}

type item struct {
	msg    string
	status ErrStatuser
	stack  []uintptr
}

func (i *item) Error() string {
	return i.msg
}

func (i *item) Status() ErrStatuser {
	return i.status
}

func (i *item) ErrCode() int {
	return i.status.ErrCode()
}

func (i *item) ErrMsg() string {
	return i.status.Msg()
}

func (i *item) HttpCode() int {
	return i.status.HTTPCode()
}

// Format used by go.uber.org/zap in Verbose
func (i *item) Format(s fmt.State, verb rune) {
	io.WriteString(s, i.msg)
	io.WriteString(s, "\n")

	for _, pc := range i.stack {
		fmt.Fprintf(s, "%+v\n", errors.Frame(pc))
	}
}

// New create a new error.
// err is the original error.
// status is the error status function defined in errstatuser.go .
// emsg is the extra message.
func New(err error, status func() ErrStatuser, errmsg ...string) Error {
	var msg string
	if len(errmsg) > 0 {
		msg = errmsg[0]
	}

	if err != nil {
		if msg == "" {
			return &item{msg: fmt.Sprintf("%s; %s", status().Msg(), err.Error()), status: status(), stack: callers()}
		}
		return &item{msg: fmt.Sprintf("%s; %s; %s", msg, status().Msg(), err.Error()), status: status(), stack: callers()}
	}
	if msg == "" {
		return &item{msg: fmt.Sprintf("%s", status().Msg()), status: status(), stack: callers()}
	}
	return &item{msg: fmt.Sprintf("%s; %s", msg, status().Msg()), status: status(), stack: callers()}
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

// Err create a new error with message
func Err(msg string) error {
	return errors.New(msg)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func WithMessage(err error, msg string) error {
	return errors.WithMessage(err, msg)
}

func WithMsg(err error, msg string) error {
	return errors.WithMessage(err, msg)
}
