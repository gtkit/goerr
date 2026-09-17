package goerr

import "regexp"

// DefaultMaxMessageRunes 网关未显式配置时的默认 message 截断长度（可按业务调整）。
const DefaultMaxMessageRunes = 512

var (
	// 前导 (^|[^0-9]) 是必需的：没有它，订单号、流水号等长数字串中任意一段 11 位子串
	// 都会被当成手机号打码，反而妨碍排障。Go 正则没有 lookbehind，只能捕获后用 ${1} 放回。
	phoneCNRegex = regexp.MustCompile(`(^|[^0-9])((?:\+?86)?1[3-9]\d{9})\b`)
	emailRegex   = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
)

func redactPhone(s string) string { return phoneCNRegex.ReplaceAllString(s, "${1}[phone]") }
func redactEmail(s string) string { return emailRegex.ReplaceAllString(s, "[email]") }

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
		out = redactPhone(out)
	}
	if opts.RedactEmail {
		out = redactEmail(out)
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
	n := 0
	for i := range msg { // range string 按 rune 迭代，i 是当前 rune 的起始字节位
		n++
		if n > maxRunes {
			return msg[:i] + "…"
		}
	}
	return msg
}

// RedactForLog 对整段字符串做常见 PII 占位替换，便于在落日志前对 Error() 等全文做一层脱敏。
// 仍可能遗漏业务自定义敏感格式，需与日志平台脱敏、权限与保留策略配合使用。
func RedactForLog(s string) string {
	return redactEmail(redactPhone(s))
}
