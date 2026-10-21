package goerr

import (
	"fmt"
	"testing"
)

func TestNewStatus_allowsCustomCode(t *testing.T) {
	st := NewStatus(Code(1), 299, "custom")
	if st.Code() != Code(1) || st.HTTPCode() != 299 || st.Msg() != "custom" {
		t.Fatalf("NewStatus() = code %v http %d msg %q", st.Code(), st.HTTPCode(), st.Msg())
	}
}

func TestMustNewStatus(t *testing.T) {
	st := MustNewStatus(ErrOrderExpired, 200, "Order expired")
	if st.Code() != ErrOrderExpired || st.HTTPCode() != 200 || st.Msg() != "Order expired" {
		t.Fatalf("MustNewStatus() = code %v http %d msg %q", st.Code(), st.HTTPCode(), st.Msg())
	}
}

func TestMustNewStatus_panicsOnInvalidCode(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = MustNewStatus(Code(1), 200, "invalid")
}

func TestNewStatus_panicsOnUnregisteredCode(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: code in valid range but missing from codeTable")
		}
	}()
	_ = newStatus(Code(10999999))
}

func ExampleNewStatus() {
	// 自定义错误码，不校验编码规范，业务方可自由编码
	st := NewStatus(Code(900001), 200, "Custom business error")
	fmt.Println(st.Code())
	fmt.Println(st.HTTPCode())
	fmt.Println(st.Msg())
	// Output:
	// 900001
	// 200
	// Custom business error
}

func ExampleMustNewStatus() {
	// 自定义错误码，强制校验编码规范，非法 code 会 panic
	st := MustNewStatus(ErrOrderExpired, 200, "Order expired")
	fmt.Println(st.Code())
	fmt.Println(st.HTTPCode())
	fmt.Println(st.Msg())
	// Output:
	// 10010707
	// 200
	// Order expired
}

// 固定已有公开映射；期望值不在运行时从 codeTable 或 getter 反推。
func TestStatusCatalog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		get      func() *Status
		code     Code
		httpCode int
		message  string
	}{
		{"StatusOK", StatusOK, 10010000, 200, "成功"},
		{"StatusVip", StatusVip, 10010007, 200, "非 VIP 用户"},
		{"StatusDuplicateRequest", StatusDuplicateRequest, 10010025, 200, "重复请求"},
		{"StatusInsufficientBalance", StatusInsufficientBalance, 10010027, 200, "余额不足"},
		{"StatusInsufficientStock", StatusInsufficientStock, 10010028, 200, "库存不足"},
		{"StatusBusinessRuleViolation", StatusBusinessRuleViolation, 10010029, 200, "业务规则限制"},
		{"StatusAccountDisabled", StatusAccountDisabled, 10010030, 200, "账号已禁用"},
		{"StatusOperationNotAllowed", StatusOperationNotAllowed, 10010031, 200, "当前状态不允许此操作"},
		{"StatusQuotaExceeded", StatusQuotaExceeded, 10010032, 200, "配额已用尽"},
		{"StatusUserNotRegistered", StatusUserNotRegistered, 10010035, 200, "用户未注册"},
		{"StatusUserDeactivated", StatusUserDeactivated, 10010036, 200, "用户已注销"},
		{"StatusPasswordIncorrect", StatusPasswordIncorrect, 10010037, 200, "密码错误"},
		{"StatusCaptchaInvalid", StatusCaptchaInvalid, 10010038, 200, "验证码错误"},
		{"StatusCaptchaExpired", StatusCaptchaExpired, 10010039, 200, "验证码已过期"},
		{"StatusAccountLocked", StatusAccountLocked, 10010040, 200, "账号已锁定"},
		{"StatusLoginNotPerformed", StatusLoginNotPerformed, 10010043, 200, "未执行登录操作"},
		{"StatusOperationInProgress", StatusOperationInProgress, 10010047, 200, "操作正在处理中，请稍后重试"},
		{"StatusRiskControlRejected", StatusRiskControlRejected, 10010048, 200, "操作被风控拦截"},
		{"StatusFeatureDisabled", StatusFeatureDisabled, 10010049, 200, "功能暂未开放"},
		{"StatusOrderNotFound", StatusOrderNotFound, 10010701, 200, "订单不存在"},
		{"StatusOrderStatusInvalid", StatusOrderStatusInvalid, 10010702, 200, "订单状态异常"},
		{"StatusOrderClosed", StatusOrderClosed, 10010703, 200, "订单已关闭"},
		{"StatusOrderCanceled", StatusOrderCanceled, 10010704, 200, "订单已取消"},
		{"StatusOrderPaid", StatusOrderPaid, 10010705, 200, "订单已支付"},
		{"StatusOrderCompleted", StatusOrderCompleted, 10010706, 200, "订单已完成"},
		{"StatusOrderExpired", StatusOrderExpired, 10010707, 200, "订单已过期"},
		{"StatusOrderRefunded", StatusOrderRefunded, 10010708, 200, "订单已退款"},
		{"StatusOrderPending", StatusOrderPending, 10010709, 200, "订单待支付"},
		{"StatusOrderPaying", StatusOrderPaying, 10010710, 200, "订单支付中"},
		{"StatusRequestFail", StatusRequestFail, 10010015, 400, "请求失败"},
		{"StatusParams", StatusParams, 10010002, 400, "参数不合法"},
		{"StatusInvalidJson", StatusInvalidJson, 10010004, 400, "JSON 格式错误"},
		{"StatusRequestBodyTooLarge", StatusRequestBodyTooLarge, 10010044, 413, "请求体过大"},
		{"StatusUnsupportedMediaType", StatusUnsupportedMediaType, 10010045, 415, "不支持的内容类型"},
		{"StatusValidateParams", StatusValidateParams, 10010003, 422, "参数校验失败"},
		{"StatusAuth", StatusAuth, 10010005, 401, "认证失败"},
		{"StatusAuthHeader", StatusAuthHeader, 10010006, 401, "认证请求头不合法"},
		{"StatusAppKey", StatusAppKey, 10010009, 401, "应用密钥无效"},
		{"StatusSign", StatusSign, 10010010, 401, "签名无效"},
		{"StatusAuthExpired", StatusAuthExpired, 10010011, 401, "令牌已过期"},
		{"StatusTokenMissing", StatusTokenMissing, 10010022, 401, "缺少令牌"},
		{"StatusTokenInvalid", StatusTokenInvalid, 10010023, 401, "令牌无效"},
		{"StatusTokenRevoked", StatusTokenRevoked, 10010024, 401, "令牌已吊销"},
		{"StatusNotLoggedIn", StatusNotLoggedIn, 10010033, 401, "未登录"},
		{"StatusLoginExpired", StatusLoginExpired, 10010034, 401, "登录已过期"},
		{"StatusSessionKickedOut", StatusSessionKickedOut, 10010042, 401, "会话已被踢下线"},
		{"StatusPermission", StatusPermission, 10010008, 403, "权限不足"},
		{"StatusNotFound", StatusNotFound, 10010013, 404, "资源不存在"},
		{"StatusRecordNotFound", StatusRecordNotFound, 10010026, 404, "记录不存在"},
		{"StatusConflict", StatusConflict, 10010016, 409, "资源冲突"},
		{"StatusAlreadyExists", StatusAlreadyExists, 10010017, 409, "资源已存在"},
		{"StatusUserAlreadyRegistered", StatusUserAlreadyRegistered, 10010041, 409, "用户已注册"},
		{"StatusVersionConflict", StatusVersionConflict, 10010046, 409, "数据已被修改，请刷新后重试"},
		{"StatusTooManyRequests", StatusTooManyRequests, 10010014, 429, "请求过于频繁"},
		{"StatusInternalServer", StatusInternalServer, 10010001, 500, "服务器内部错误"},
		{"StatusTimeout", StatusTimeout, 10010012, 500, "服务器响应超时"},
		{"StatusCanceled", StatusCanceled, 10010020, 500, "请求已取消"},
		{"StatusDeadlineExceeded", StatusDeadlineExceeded, 10010021, 500, "请求超时"},
		{"StatusServiceUnavailable", StatusServiceUnavailable, 10010018, 503, "服务不可用"},
		{"StatusDependencyUnavailable", StatusDependencyUnavailable, 10010019, 503, "依赖服务不可用"},
		{"StatusSystemMaintenance", StatusSystemMaintenance, 10010050, 503, "系统维护中"},
		{"StatusElasticsearchServer", StatusElasticsearchServer, 10010101, 500, "Elasticsearch 服务异常"},
		{"StatusElasticsearchDSL", StatusElasticsearchDSL, 10010102, 500, "Elasticsearch 查询语句错误"},
		{"StatusElasticsearchIndexNotFound", StatusElasticsearchIndexNotFound, 10010103, 500, "Elasticsearch 索引不存在"},
		{"StatusElasticsearchBulkPartial", StatusElasticsearchBulkPartial, 10010104, 500, "Elasticsearch 批量操作部分失败"},
		{"StatusMysqlServer", StatusMysqlServer, 10010201, 500, "MySQL 服务异常"},
		{"StatusMysqlSQL", StatusMysqlSQL, 10010202, 500, "SQL 语句错误"},
		{"StatusMysqlQuery", StatusMysqlQuery, 10010203, 500, "MySQL 查询失败"},
		{"StatusMysqlDeadlock", StatusMysqlDeadlock, 10010204, 500, "MySQL 死锁"},
		{"StatusMysqlTransaction", StatusMysqlTransaction, 10010205, 500, "MySQL 事务执行失败"},
		{"StatusMongoServer", StatusMongoServer, 10010301, 500, "MongoDB 服务异常"},
		{"StatusMongoDSL", StatusMongoDSL, 10010302, 500, "MongoDB 查询语句错误"},
		{"StatusMongoQuery", StatusMongoQuery, 10010303, 500, "MongoDB 查询失败"},
		{"StatusRedisServer", StatusRedisServer, 10010401, 500, "Redis 服务异常"},
		{"StatusRedisQuery", StatusRedisQuery, 10010402, 500, "Redis 查询失败"},
		{"StatusRedisLockAcquireFailed", StatusRedisLockAcquireFailed, 10010403, 500, "Redis 分布式锁获取失败"},
		{"StatusRedisLockLost", StatusRedisLockLost, 10010404, 500, "Redis 分布式锁已丢失"},
		{"StatusKafkaServer", StatusKafkaServer, 10010501, 500, "Kafka 服务异常"},
		{"StatusKafkaProducer", StatusKafkaProducer, 10010502, 500, "Kafka 生产者异常"},
		{"StatusKafkaConsumer", StatusKafkaConsumer, 10010503, 500, "Kafka 消费者异常"},
		{"StatusKafkaMessageDecode", StatusKafkaMessageDecode, 10010504, 500, "Kafka 消息解析失败"},
		{"StatusKafkaRetryExhausted", StatusKafkaRetryExhausted, 10010505, 500, "Kafka 消息重试次数已耗尽"},
		{"StatusRabbitMQServer", StatusRabbitMQServer, 10010601, 500, "RabbitMQ 服务异常"},
		{"StatusRabbitMQProducer", StatusRabbitMQProducer, 10010602, 500, "RabbitMQ 生产者异常"},
		{"StatusRabbitMQConsumer", StatusRabbitMQConsumer, 10010603, 500, "RabbitMQ 消费者异常"},
		{"StatusRabbitMQMessageDecode", StatusRabbitMQMessageDecode, 10010604, 500, "RabbitMQ 消息解析失败"},
		{"StatusRabbitMQRetryExhausted", StatusRabbitMQRetryExhausted, 10010605, 500, "RabbitMQ 消息重试次数已耗尽"},
		{"StatusPaymentNotFound", StatusPaymentNotFound, 10010801, 200, "支付单不存在"},
		{"StatusPaymentAlreadyPaid", StatusPaymentAlreadyPaid, 10010802, 200, "订单已支付，请勿重复支付"},
		{"StatusPaymentClosed", StatusPaymentClosed, 10010803, 200, "支付单已关闭"},
		{"StatusPaymentExpired", StatusPaymentExpired, 10010804, 200, "支付已超时，请重新下单"},
		{"StatusPaymentMethodNotSupported", StatusPaymentMethodNotSupported, 10010805, 200, "不支持的支付方式"},
		{"StatusPayerIdentityMissing", StatusPayerIdentityMissing, 10010806, 200, "缺少支付方身份信息"},
		{"StatusPaymentAmountMismatch", StatusPaymentAmountMismatch, 10010807, 400, "支付金额与订单不一致"},
		{"StatusPaymentSignInvalid", StatusPaymentSignInvalid, 10010808, 401, "支付回调验签失败"},
		{"StatusPaymentCallbackInvalid", StatusPaymentCallbackInvalid, 10010809, 400, "支付回调参数非法"},
		{"StatusPaymentCreateFailed", StatusPaymentCreateFailed, 10010810, 500, "发起支付失败"},
		{"StatusPaymentQueryFailed", StatusPaymentQueryFailed, 10010811, 500, "查询支付结果失败"},
		{"StatusPaymentChannelUnavailable", StatusPaymentChannelUnavailable, 10010812, 503, "支付渠道暂不可用"},
		{"StatusRefundNotFound", StatusRefundNotFound, 10010813, 200, "退款单不存在"},
		{"StatusRefundAmountExceeded", StatusRefundAmountExceeded, 10010814, 200, "退款金额超过可退金额"},
		{"StatusRefundInProgress", StatusRefundInProgress, 10010815, 200, "退款处理中"},
		{"StatusRefundFailed", StatusRefundFailed, 10010816, 500, "退款失败"},
		{"StatusPhoneAlreadyBound", StatusPhoneAlreadyBound, 10010901, 409, "手机号已被其他账号绑定"},
		{"StatusThirdPartyAlreadyBound", StatusThirdPartyAlreadyBound, 10010902, 409, "该第三方账号已绑定其他用户"},
		{"StatusThirdPartyNotBound", StatusThirdPartyNotBound, 10010903, 200, "请先绑定手机号"},
		{"StatusThirdPartyAuthFailed", StatusThirdPartyAuthFailed, 10010904, 200, "第三方授权失败"},
		{"StatusThirdPartyUnavailable", StatusThirdPartyUnavailable, 10010905, 503, "第三方登录服务暂不可用"},
		{"StatusSmsSendTooFrequent", StatusSmsSendTooFrequent, 10010906, 200, "验证码发送过于频繁"},
		{"StatusSmsDailyLimitExceeded", StatusSmsDailyLimitExceeded, 10010907, 200, "今日验证码次数已用完"},
		{"StatusSmsSendFailed", StatusSmsSendFailed, 10010908, 500, "验证码发送失败"},
		{"StatusLLMAuth", StatusLLMAuth, 10011001, 500, "AI 服务配置错误"},
		{"StatusLLMRateLimited", StatusLLMRateLimited, 10011002, 503, "AI 服务繁忙，请稍后重试"},
		{"StatusLLMTimeout", StatusLLMTimeout, 10011003, 500, "AI 服务响应超时"},
		{"StatusLLMContextTooLong", StatusLLMContextTooLong, 10011004, 200, "输入内容过长"},
		{"StatusLLMContentFiltered", StatusLLMContentFiltered, 10011005, 200, "内容不符合规范"},
		{"StatusLLMInvalidRequest", StatusLLMInvalidRequest, 10011006, 500, "AI 请求参数错误"},
		{"StatusLLMProviderError", StatusLLMProviderError, 10011007, 503, "AI 服务异常"},
		{"StatusLLMModelNotAvailable", StatusLLMModelNotAvailable, 10011008, 200, "所选模型暂不可用"},
		{"StatusLLMStructuredOutput", StatusLLMStructuredOutput, 10011009, 500, "AI 返回解析失败"},
		{"StatusLLMStreamInterrupted", StatusLLMStreamInterrupted, 10011010, 500, "AI 输出中断"},
		{"StatusLLMPricingMissing", StatusLLMPricingMissing, 10011011, 500, "模型计费配置缺失"},
		{"StatusFileEmpty", StatusFileEmpty, 10011101, 400, "文件为空"},
		{"StatusFileTooLarge", StatusFileTooLarge, 10011102, 400, "文件超过大小限制"},
		{"StatusFileTypeNotAllowed", StatusFileTypeNotAllowed, 10011103, 400, "不支持的文件类型"},
		{"StatusFileCountExceeded", StatusFileCountExceeded, 10011104, 400, "文件数量超过限制"},
		{"StatusFileNotFound", StatusFileNotFound, 10011105, 404, "文件不存在"},
		{"StatusFileChecksumMismatch", StatusFileChecksumMismatch, 10011106, 400, "文件校验失败"},
		{"StatusFileUploadFailed", StatusFileUploadFailed, 10011107, 500, "文件上传失败"},
		{"StatusFileStorageUnavailable", StatusFileStorageUnavailable, 10011108, 503, "存储服务暂不可用"},
		{"StatusEmailAlreadyBound", StatusEmailAlreadyBound, 10011201, 409, "邮箱已被绑定"},
		{"StatusEmailNotVerified", StatusEmailNotVerified, 10011202, 200, "邮箱未验证"},
		{"StatusEmailSendTooFrequent", StatusEmailSendTooFrequent, 10011203, 200, "邮件发送过于频繁"},
		{"StatusEmailTemplateNotFound", StatusEmailTemplateNotFound, 10011204, 500, "邮件模板不存在"},
		{"StatusEmailSendFailed", StatusEmailSendFailed, 10011205, 500, "邮件发送失败"},
		{"StatusEmailServiceUnavailable", StatusEmailServiceUnavailable, 10011206, 503, "邮件服务暂不可用"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.get()
			if got == nil {
				t.Fatal("a predefined status must not be nil")
			}
			if got.Code() != tt.code || got.HTTPCode() != tt.httpCode || got.Msg() != tt.message {
				t.Fatalf("status = (%d, %d, %q), want (%d, %d, %q)", got.Code(), got.HTTPCode(), got.Msg(), tt.code, tt.httpCode, tt.message)
			}
			if got != tt.get() {
				t.Error("repeated calls must return the same predefined read-only instance")
			}
			if got.Code().Message() != tt.message {
				t.Errorf("Code.Message() = %q, want %q", got.Code().Message(), tt.message)
			}
		})
	}
}
