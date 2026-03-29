package goerr

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestClampMessage(t *testing.T) {
	tests := []struct {
		in   string
		max  int
		want string
	}{
		{"", 5, ""},
		{"hello", 10, "hello"},
		{"你好世界", 2, "你好…"},
		{"ab", 0, ""},
	}
	for _, tt := range tests {
		if got := ClampMessage(tt.in, tt.max); got != tt.want {
			t.Errorf("ClampMessage(%q, %d) = %q, want %q", tt.in, tt.max, got, tt.want)
		}
	}
}

func TestSanitizeForClient(t *testing.T) {
	msg := "联系 13812345678 或 user@test.com 谢谢"
	out := SanitizeForClient(msg, SanitizeOptions{MaxRunes: 100, RedactPhone: true, RedactEmail: true})
	if strings.Contains(out, "13812345678") || strings.Contains(out, "user@test.com") {
		t.Errorf("expected redaction: %q", out)
	}
	if !strings.Contains(out, "[phone]") || !strings.Contains(out, "[email]") {
		t.Errorf("placeholders: %q", out)
	}
}

func TestSanitizeForClient_defaultMax(t *testing.T) {
	long := strings.Repeat("x", DefaultMaxMessageRunes+50)
	out := SanitizeForClient(long, SanitizeOptions{})
	if n := utf8.RuneCountInString(out); n != DefaultMaxMessageRunes+1 { // 正文截断 + 省略号
		t.Errorf("rune count = %d, want %d", n, DefaultMaxMessageRunes+1)
	}
	if !strings.HasSuffix(out, "…") {
		t.Error("expected ellipsis")
	}
}

func FuzzClampMessage(f *testing.F) {
	f.Add("hello", 3)
	f.Fuzz(func(t *testing.T, s string, n int) {
		n = n % 2048
		if n < 0 {
			n = -n % 2048
		}
		_ = ClampMessage(s, n)
	})
}

func TestRedactForLog(t *testing.T) {
	s := RedactForLog("call 13900001111 fail, mail=a@b.co")
	if strings.Contains(s, "13900001111") || strings.Contains(s, "a@b.co") {
		t.Errorf("RedactForLog: %q", s)
	}
}
