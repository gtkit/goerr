package test

import (
	"errors"
	"testing"

	"github.com/gtkit/goerr"
)

// 外部包引用路径下的冒烟测试，确保导出 API 与文档一致。
func TestExportedAPI(t *testing.T) {
	item := goerr.New(errors.New("db down"), goerr.StatusMysqlServer(), "服务繁忙")
	if item.Message() == item.Error() {
		t.Error("Message() should differ from Error() when cause is present")
	}
	if item.Code() != goerr.ErrMysqlServer {
		t.Errorf("Code = %v", item.Code())
	}
	_, ok := goerr.AsItem(item)
	if !ok {
		t.Fatal("AsItem")
	}
}
