package test

import (
	"testing"

	"github.com/gtkit/goerr"
)

func TestItemStatus(t *testing.T) {
	cause := goerr.Newf(goerr.StatusParams, "bad field")
	item := goerr.New(cause, goerr.StatusInvalidJson, "http error")

	if code := item.Code(); code != goerr.ErrInvalidJson {
		t.Errorf("Code() = %v, want ErrInvalidJson", code)
	}
	if http := item.HTTPStatus(); http != goerr.StatusInvalidJson().HTTPCode() {
		t.Errorf("HTTPStatus() = %d", http)
	}
	if msg := item.Message(); msg == "" {
		t.Error("Message() empty")
	}
	if item.Error() == "" {
		t.Error("Error() empty")
	}
	t.Log("code:", item.Code(), "http:", item.HTTPStatus(), "msg:", item.Message(), "err:", item.Error())
}
