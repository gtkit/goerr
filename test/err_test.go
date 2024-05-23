package test

import (
	"testing"

	"github.com/gtkit/goerr"
)

func TestName(t *testing.T) {
	err := goerr.Err("test error")
	errwithmsg := goerr.WithMsg(err, "with message")

	myerr := goerr.New(errwithmsg, goerr.InvalidJson(), "http error")
	errcode := myerr.Status().ErrCode()
	t.Log("errcode:", errcode)

	httpcode := myerr.Status().HTTPCode()
	t.Log("httpcode:", httpcode)

	msg := myerr.Status().Msg()
	t.Log("msg:", msg)

	errmsg := myerr.Error()
	t.Log("errmsg:", errmsg)
}
