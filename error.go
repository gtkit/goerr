package goerr

import (
	"errors"
	"fmt"
	"runtime"
)

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

// Error an error with caller stack information
type Error interface {
	error
	Status() *ErrStatus
	ErrCode() int32
	ErrMsg() string
	HttpCode() int
}

type Item struct {
	msg    string
	status *ErrStatus
	stack  []uintptr
}

func (i *Item) Error() string {
	return i.msg
}

func (i *Item) Status() *ErrStatus {
	return i.status
}

func (i *Item) ErrCode() int32 {
	if i.status == nil {
		return 0
	}
	return i.status.ErrCode()
}

func (i *Item) ErrMsg() string {
	if i.status == nil {
		return "~ok~"
	}
	return i.status.Msg()
}

func (i *Item) HttpCode() int {
	if i.status == nil {
		return 200
	}
	return i.status.HTTPCode()
}

// Format used by go.uber.org/zap in Verbose
//func (i *Item) Format(s fmt.State, verb rune) {
//	io.WriteString(s, i.msg)
//	io.WriteString(s, "\n")
//
//	if i.status != nil {
//		for _, pc := range i.stack {
//			fmt.Fprintf(s, "%+v\n", errors.Frame(pc))
//		}
//	}
//
//}

// New create a new error.
// err is the original error.
// status is the error status function defined in errstatuser.go .
// emsg is the extra message, if it is not empty, it will replace the original error message.
func New(err error, status func() *ErrStatus, errmsg ...string) *Item {
	var (
		msg       string
		errStatus *ErrStatus
	)

	if status == nil || status() == nil {
		errStatus = &ErrStatus{
			errCode:  0,
			httpCode: 200,
			msg:      "~ok~",
		}
	} else {
		errStatus = status()
	}

	if len(errmsg) > 0 {
		msg = errmsg[0]
	}

	if err != nil {
		if msg == "" {
			return &Item{msg: fmt.Sprintf("%s; %s", errStatus.Msg(), err.Error()), status: errStatus, stack: callers()}
		}
		return &Item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), status: errStatus, stack: callers()}
	}
	if msg == "" {
		return &Item{msg: fmt.Sprintf("%s", errStatus.Msg()), status: errStatus, stack: callers()}
	}
	return &Item{msg: fmt.Sprintf("%s;", msg), status: errStatus, stack: callers()}
}

// Errorf create a new error
func Errorf(format string, args ...interface{}) Error {
	return &Item{msg: fmt.Sprintf(format, args...), stack: callers()}
}

// Wrap with some extra message into err
func Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}

	var e *Item
	if ok := errors.As(err, &e); !ok {
		return &Item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
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

	var e *Item
	if ok := errors.As(err, &e); !ok {
		return &Item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// WithStack add caller stack information
func WithStack(err error) Error {
	if err == nil {
		return nil
	}

	if e, ok := AsType[*Item](err); ok {
		return e
	}

	return &Item{msg: err.Error(), stack: callers()}
}

// Err create a new error with message
func Err(msg string) error {
	return errors.New(msg)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, &target)
}

func AsType[E error](err error) (E, bool) {
	return errors.AsType[E](err)
}

func WithMsg(err error, msg string) error {
	return errors.New(msg + ": " + err.Error())
}
