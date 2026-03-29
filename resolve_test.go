package goerr

import (
	"errors"
	"testing"
)

func TestResolveStatus(t *testing.T) {
	fb := StatusOK()
	st := StatusParams()

	got, err := ResolveStatus(nil, fb)
	if err != nil || got != fb {
		t.Fatalf("nil: got %v err %v", got, err)
	}
	got, err = ResolveStatus((*Status)(nil), fb)
	if err != nil || got != fb {
		t.Fatalf("nil *Status: got %v err %v", got, err)
	}
	got, err = ResolveStatus(st, fb)
	if err != nil || got != st {
		t.Fatalf("*Status: got %v err %v", got, err)
	}
	got, err = ResolveStatus((func() *Status)(nil), fb)
	if err == nil {
		t.Fatal("nil func should error")
	}
	fn := func() *Status { return st }
	got, err = ResolveStatus(fn, fb)
	if err != nil || got != st {
		t.Fatalf("func: got %v err %v", got, err)
	}
	fnNil := func() *Status { return nil }
	got, err = ResolveStatus(fnNil, fb)
	if err != nil || got != fb {
		t.Fatalf("func returns nil: got %v err %v", got, err)
	}
	_, err = ResolveStatus(42, fb)
	if err == nil {
		t.Fatal("unsupported type should error")
	}
}

func TestNewFrom(t *testing.T) {
	item, err := NewFrom(errors.New("e"), StatusNotFound(), "x")
	if err != nil || item.Message() != "x" {
		t.Fatalf("%v %v", item, err)
	}
	_, err = NewFrom(nil, struct{}{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewfFrom(t *testing.T) {
	item, err := NewfFrom(StatusParams(), "n=%d", 1)
	if err != nil || item.Message() != "n=1" {
		t.Fatalf("%v %v", item, err)
	}
	_, err = NewfFrom("bad", "x")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWrapStatusFrom(t *testing.T) {
	item, err := WrapStatusFrom(errors.New("e"), StatusRedisServer())
	if err != nil || item.Code() != ErrRedisServer {
		t.Fatal(item, err)
	}
	item, err = WrapStatusFrom(nil, StatusOK())
	if item != nil || err != nil {
		t.Fatal("nil err")
	}
	_, err = WrapStatusFrom(errors.New("e"), 1)
	if err == nil {
		t.Fatal("expected error")
	}
}
