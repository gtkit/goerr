package goerr

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestClampMessage(t *testing.T) {
	t.Parallel()
	maxInt := int(^uint(0) >> 1)
	tests := []struct {
		name string
		in   string
		max  int
		want string
	}{
		{name: "empty", in: "", max: 5, want: ""},
		{name: "below limit", in: "hello", max: 10, want: "hello"},
		{name: "at limit", in: "你好", max: 2, want: "你好"},
		{name: "Chinese", in: "你好世界", max: 2, want: "你好…"},
		{name: "emoji", in: "🙂🚀x", max: 2, want: "🙂🚀…"},
		{name: "combining runes", in: "e\u0301!", max: 1, want: "e…"},
		{name: "existing ellipsis", in: "ab…", max: 1, want: "a…"},
		{name: "zero", in: "ab", max: 0, want: ""},
		{name: "negative", in: "你好", max: -1, want: ""},
		{name: "minimum int", in: "hello", max: -maxInt - 1, want: ""},
		{name: "maximum int", in: "你好🙂", max: maxInt, want: "你好🙂"},
		{name: "invalid bytes kept", in: "\xffa\xfe", max: 3, want: "\xffa\xfe"},
		{name: "invalid prefix", in: "\xffa\xfe", max: 1, want: "\xff…"},
		{name: "valid rune after invalid", in: "\xff🙂x", max: 2, want: "\xff🙂…"},
		{name: "encoded RuneError", in: "\ufffdx", max: 1, want: "\ufffd…"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ClampMessage(tt.in, tt.max)
			if got != tt.want {
				t.Errorf("ClampMessage(%q, %d) = %q, want %q", tt.in, tt.max, got, tt.want)
			}
			assertClampProperties(t, tt.in, tt.max, got)
		})
	}
}

func assertClampProperties(t *testing.T, input string, limit int, got string) {
	t.Helper()
	if limit <= 0 {
		if got != "" {
			t.Fatalf("non-positive limit %d must return empty text, got %q", limit, got)
		}
		return
	}
	runes := utf8.RuneCountInString(input)
	if runes <= limit {
		if got != input {
			t.Fatalf("untruncated text must preserve its original bytes: %q -> %q", input, got)
		}
	} else {
		if !strings.HasSuffix(got, "…") {
			t.Fatalf("truncated text must have the existing ellipsis suffix, got %q", got)
		}
		prefix := strings.TrimSuffix(got, "…")
		if !strings.HasPrefix(input, prefix) || utf8.RuneCountInString(prefix) != limit {
			t.Fatalf("truncated text must retain %d runes from the original byte prefix: %q -> %q", limit, input, got)
		}
		if utf8.RuneCountInString(input[len(prefix):]) != runes-limit {
			t.Fatalf("truncation must end at an original rune boundary: %q -> %q", input, got)
		}
	}
	if utf8.ValidString(input) && !utf8.ValidString(got) {
		t.Fatalf("valid UTF-8 must remain valid: %q -> %q", input, got)
	}
	if again := ClampMessage(got, limit); again != got {
		t.Fatalf("repeated truncation must be idempotent: %q -> %q", got, again)
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
	f.Add("", 0)
	f.Add("你好🙂🚀", 2)
	f.Add("e\u0301!", 1)
	f.Add("ab…", 1)
	f.Add("\xff🙂x\xfe", 2)
	maxInt := int(^uint(0) >> 1)
	f.Add("hello", maxInt)
	f.Add("hello", -maxInt-1)
	f.Add("你好", -1)
	f.Fuzz(func(t *testing.T, s string, n int) {
		assertClampProperties(t, s, n, ClampMessage(s, n))
	})
}

func TestRedactForLog(t *testing.T) {
	s := RedactForLog("call 13900001111 fail, mail=a@b.co")
	if strings.Contains(s, "13900001111") || strings.Contains(s, "a@b.co") {
		t.Errorf("RedactForLog: %q", s)
	}
}

func TestRedactPhone_boundary(t *testing.T) {
	tests := []struct{ in, want string }{
		{"order_no=2024113812345678", "order_no=2024113812345678"}, // 长数字串里的 11 位子串不能被当成手机号
		{"tel 13900001111", "tel [phone]"},
		{"tel:+8613900001111.", "tel:[phone]."},
		{"tel:8613900001111", "tel:[phone]"},
		{"13900001111", "[phone]"},
		{"13800000000,13900000000", "[phone],[phone]"},
	}
	for _, tt := range tests {
		if got := RedactForLog(tt.in); got != tt.want {
			t.Errorf("RedactForLog(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
