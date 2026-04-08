package goerr

import (
	"testing"
)

func TestValidateCode_range(t *testing.T) {
	if err := ValidateCode(Code(9999999)); err == nil {
		t.Error("below min should fail")
	}
	if err := ValidateCode(Code(11000000)); err == nil {
		t.Error("above max should fail")
	}
	if err := ValidateCode(ErrNo); err != nil {
		t.Errorf("ErrNo should validate: %v", err)
	}
}

func TestCode_Message_knownAndDefault(t *testing.T) {
	if got := ErrParams.Message(); got != "参数不合法" {
		t.Errorf("ErrParams.Message() = %q, want %q", got, "参数不合法")
	}
	if got := (Code(10999998)).Message(); got != "未知错误" {
		t.Errorf("unknown code should return default, got %q", got)
	}
}
