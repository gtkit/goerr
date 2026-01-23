package goerr

import (
	"strconv"
	"sync/atomic"
)

// CodeConfig 错误码配置
type CodeConfig struct {
	Code       string
	Message    string
	HTTPStatus int
}

// codeRegistry 全局错误码注册表
// 使用 atomic.Value 实现无锁读取,只在初始化时写入
var codeRegistry atomic.Value // map[string]*CodeConfig

func init() {
	// 初始化空注册表
	codeRegistry.Store(make(map[string]*CodeConfig))
}

// RegisterCode 注册错误码(仅在初始化阶段调用)
func RegisterCode(code string, message string, httpStatus int) {
	if code == "" {
		panic("error code cannot be empty")
	}

	// 获取当前注册表
	oldRegistry := codeRegistry.Load().(map[string]*CodeConfig)

	// 检查是否已存在
	if _, exists := oldRegistry[code]; exists {
		panic("error code already registered: " + code)
	}

	// 创建新注册表(复制旧的,添加新的)
	newRegistry := make(map[string]*CodeConfig, len(oldRegistry)+1)
	for k, v := range oldRegistry {
		newRegistry[k] = v
	}
	newRegistry[code] = &CodeConfig{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}

	// 原子性替换注册表
	codeRegistry.Store(newRegistry)
}

// RegisterCodes 批量注册错误码
func RegisterCodes(codes []CodeConfig) {
	for _, cfg := range codes {
		RegisterCode(cfg.Code, cfg.Message, cfg.HTTPStatus)
	}
}

// GetCodeConfig 获取错误码配置(并发安全)
func GetCodeConfig(code string) *CodeConfig {
	registry := codeRegistry.Load().(map[string]*CodeConfig)
	return registry[code]
}

// NewWithCode 使用预定义错误码创建错误
func NewWithCode(code string) Error {
	cfg := GetCodeConfig(code)
	if cfg == nil {
		// 如果错误码未注册,使用默认配置
		return newError(code, "Unknown error", 500, nil, 2)
	}
	return newError(cfg.Code, cfg.Message, cfg.HTTPStatus, nil, 2)
}

// NewWithCodef 使用预定义错误码创建错误,可格式化消息
func NewWithCodef(code string, format string, args ...interface{}) Error {
	cfg := GetCodeConfig(code)
	httpStatus := 500
	if cfg != nil {
		httpStatus = cfg.HTTPStatus
	}
	return newError(code, formatMessage(format, args...), httpStatus, nil, 2)
}

// WrapWithCode 使用预定义错误码包装错误
func WrapWithCode(err error, code string) Error {
	if err == nil {
		return nil
	}
	cfg := GetCodeConfig(code)
	if cfg == nil {
		return newError(code, "Unknown error", 500, err, 2)
	}
	return newError(cfg.Code, cfg.Message, cfg.HTTPStatus, err, 2)
}

// WrapWithCodef 使用预定义错误码包装错误,可格式化消息
func WrapWithCodef(err error, code string, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}
	cfg := GetCodeConfig(code)
	httpStatus := 500
	if cfg != nil {
		httpStatus = cfg.HTTPStatus
	}
	return newError(code, formatMessage(format, args...), httpStatus, err, 2)
}

// formatMessage 格式化消息,避免重复分配
func formatMessage(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	// 使用 fmt.Sprintf 但尽量减少使用
	return formatSprintf(format, args...)
}

// formatSprintf 简单的格式化实现
func formatSprintf(format string, args ...interface{}) string {
	// 对于简单情况直接返回
	if len(args) == 0 {
		return format
	}

	// 这里使用标准库,在生产环境可以优化为自定义实现
	var b builder
	argIndex := 0
	inVerb := false

	for i := 0; i < len(format); i++ {
		c := format[i]

		if inVerb {
			if c == '%' {
				b.sb.WriteByte('%')
			} else if argIndex < len(args) {
				// 简单处理 %v, %s, %d
				switch c {
				case 'v', 's':
					b.sb.WriteString(toString(args[argIndex]))
					argIndex++
				case 'd':
					if n, ok := args[argIndex].(int); ok {
						b.writeInt(n)
						argIndex++
					}
				default:
					b.sb.WriteByte('%')
					b.sb.WriteByte(c)
				}
			}
			inVerb = false
		} else if c == '%' {
			inVerb = true
		} else {
			b.sb.WriteByte(c)
		}
	}

	return b.string()
}

// toString 将任意类型转换为字符串
func toString(v interface{}) string {
	if v == nil {
		return "<nil>"
	}

	switch val := v.(type) {
	case string:
		return val
	case error:
		return val.Error()
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		// 使用 fmt.Sprint 作为后备
		return formatSprint(v)
	}
}

// formatSprint 后备格式化
func formatSprint(v interface{}) string {
	// 这里可以实现更高效的类型转换
	// 当前使用简单实现
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return "<unknown>"
}

// 预定义常用错误码
const (
	// 通用错误码
	ErrInternal           = "INTERNAL_ERROR"
	ErrInvalidArgument    = "INVALID_ARGUMENT"
	ErrNotFound           = "NOT_FOUND"
	ErrAlreadyExists      = "ALREADY_EXISTS"
	ErrPermissionDenied   = "PERMISSION_DENIED"
	ErrUnauthenticated    = "UNAUTHENTICATED"
	ErrResourceExhausted  = "RESOURCE_EXHAUSTED"
	ErrFailedPrecondition = "FAILED_PRECONDITION"
	ErrAborted            = "ABORTED"
	ErrOutOfRange         = "OUT_OF_RANGE"
	ErrUnimplemented      = "UNIMPLEMENTED"
	ErrUnavailable        = "UNAVAILABLE"
	ErrDataLoss           = "DATA_LOSS"

	// 业务错误码
	ErrTimeout  = "TIMEOUT"
	ErrCanceled = "CANCELED"
	ErrUnknown  = "UNKNOWN"
)

func init() {
	// 注册预定义错误码
	RegisterCodes([]CodeConfig{
		{ErrInternal, "Internal server error", 500},
		{ErrInvalidArgument, "Invalid argument", 400},
		{ErrNotFound, "Resource not found", 404},
		{ErrAlreadyExists, "Resource already exists", 409},
		{ErrPermissionDenied, "Permission denied", 403},
		{ErrUnauthenticated, "Unauthenticated", 401},
		{ErrResourceExhausted, "Resource exhausted", 429},
		{ErrFailedPrecondition, "Failed precondition", 400},
		{ErrAborted, "Operation aborted", 409},
		{ErrOutOfRange, "Out of range", 400},
		{ErrUnimplemented, "Unimplemented", 501},
		{ErrUnavailable, "Service unavailable", 503},
		{ErrDataLoss, "Data loss", 500},
		{ErrTimeout, "Operation timeout", 504},
		{ErrCanceled, "Operation canceled", 499},
		{ErrUnknown, "Unknown error", 500},
	})
}
