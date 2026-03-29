package goerr

import (
	"regexp"
	"unicode/utf8"
)

// DefaultMaxMessageRunes 网关未显式配置时的默认 message 截断长度（可按业务调整）。
const DefaultMaxMessageRunes = 512

var (
	phoneCNRegex = regexp.MustCompile(`(?i)(?:\+?86)?1[3-9]\d{9}\b`)
	emailRegex   = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
)

// SanitizeOptions 控制对外 message 的截断与简单脱敏（网关/BFF 层使用）。
type SanitizeOptions struct {
	// MaxRunes 最大 Unicode 字符数；<=0 时使用 DefaultMaxMessageRunes。
	MaxRunes int
	// RedactPhone 将常见中国大陆手机号模式替换为占位符。
	RedactPhone bool
	// RedactEmail 将邮箱模式替换为占位符。
	RedactEmail bool
}

// SanitizeForClient 对即将写入 HTTP 响应体的文案做长度限制与可选脱敏。
// 本包只做通用、无外部依赖的保守处理；更复杂策略（身份证、银行卡、自定义 token）请在业务网关实现。
func SanitizeForClient(msg string, opts SanitizeOptions) string {
	out := msg
	if opts.RedactPhone {
		out = phoneCNRegex.ReplaceAllString(out, "[phone]")
	}
	if opts.RedactEmail {
		out = emailRegex.ReplaceAllString(out, "[email]")
	}
	max := opts.MaxRunes
	if max <= 0 {
		max = DefaultMaxMessageRunes
	}
	return ClampMessage(out, max)
}

// ClampMessage 按 Unicode 字符数截断，超出部分追加省略号「…」。
func ClampMessage(msg string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	n := utf8.RuneCountInString(msg)
	if n <= maxRunes {
		return msg
	}
	rs := []rune(msg)
	return string(rs[:maxRunes]) + "…"
}

// RedactForLog 对整段字符串做常见 PII 占位替换，便于在落日志前对 Error() 等全文做一层脱敏。
// 仍可能遗漏业务自定义敏感格式，需与日志平台脱敏、权限与保留策略配合使用。
func RedactForLog(s string) string {
	s = phoneCNRegex.ReplaceAllString(s, "[phone]")
	s = emailRegex.ReplaceAllString(s, "[email]")
	return s
}
