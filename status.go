package goerr

// Status 错误状态，包含业务错误码、HTTP 状态码和消息。
// 所有字段私有且构造后不再写入；StatusXxx() 返回的是包级预构建实例的指针，
// 供全进程共享，请勿对其解引用赋值（那会污染所有调用方看到的状态）。
type Status struct {
	code     Code
	httpCode int
	msg      string
}

// NewStatus 创建自定义 Status，用于业务方扩展错误码。
// code 不做强制校验，业务方可自由定义编码规范。
// 如需强制校验是否符合本库默认编码规范，请使用 [MustNewStatus]。
func NewStatus(code Code, httpCode int, msg string) *Status {
	return &Status{code: code, httpCode: httpCode, msg: msg}
}

// MustNewStatus 创建自定义 Status，并强制校验 code 是否符合本库默认编码规范。
// 校验失败时直接 panic，适合包级变量初始化；运行时动态输入请先调用 [ValidateCode]。
func MustNewStatus(code Code, httpCode int, msg string) *Status {
	mustValidCode(code)
	return NewStatus(code, httpCode, msg)
}

func mustValidCode(code Code) {
	if err := ValidateCode(code); err != nil {
		panic(err)
	}
}

// Code 返回业务错误码。
func (s *Status) Code() Code { return s.code }

// HTTPCode 返回对应的 HTTP 状态码。
func (s *Status) HTTPCode() int { return s.httpCode }

// Msg 返回错误消息。
func (s *Status) Msg() string { return s.msg }

// newStatus 从 codeTable 取 HTTP 状态码与默认文案构造 Status。
// 未登记的码在包初始化期 panic，让「加了码却漏配元信息」立刻暴露而非留到运行时。
func newStatus(code Code) Status {
	mustValidCode(code)
	m, ok := codeTable[code]
	if !ok {
		panic("goerr: code " + code.String() + " is not registered in codeTable")
	}
	return Status{code: code, httpCode: m.httpCode, msg: m.msg}
}

// 预构建的不可变 Status 值，HTTP 状态码与文案来自 codeTable。
// 调用方通过 StatusXxx() 获取指针，指向包级别变量，零堆分配。
var (
	// 成功 & 业务逻辑错误 → HTTP 200
	statusOK                    = newStatus(ErrNo)
	statusVip                   = newStatus(ErrVipRights)
	statusDuplicateRequest      = newStatus(ErrDuplicateRequest)
	statusInsufficientBalance   = newStatus(ErrInsufficientBalance)
	statusInsufficientStock     = newStatus(ErrInsufficientStock)
	statusBusinessRuleViolation = newStatus(ErrBusinessRuleViolation)
	statusAccountDisabled       = newStatus(ErrAccountDisabled)
	statusOperationNotAllowed   = newStatus(ErrOperationNotAllowed)
	statusQuotaExceeded         = newStatus(ErrQuotaExceeded)
	statusUserNotRegistered     = newStatus(ErrUserNotRegistered)
	statusUserDeactivated       = newStatus(ErrUserDeactivated)
	statusPasswordIncorrect     = newStatus(ErrPasswordIncorrect)
	statusCaptchaInvalid        = newStatus(ErrCaptchaInvalid)
	statusCaptchaExpired        = newStatus(ErrCaptchaExpired)
	statusAccountLocked         = newStatus(ErrAccountLocked)
	statusLoginNotPerformed     = newStatus(ErrLoginNotPerformed)
	statusOperationInProgress   = newStatus(ErrOperationInProgress)
	statusRiskControlRejected   = newStatus(ErrRiskControlRejected)
	statusFeatureDisabled       = newStatus(ErrFeatureDisabled)

	// 订单模块（模块代号 07）业务态 → HTTP 200
	statusOrderNotFound      = newStatus(ErrOrderNotFound)
	statusOrderStatusInvalid = newStatus(ErrOrderStatusInvalid)
	statusOrderClosed        = newStatus(ErrOrderClosed)
	statusOrderCanceled      = newStatus(ErrOrderCanceled)
	statusOrderPaid          = newStatus(ErrOrderPaid)
	statusOrderCompleted     = newStatus(ErrOrderCompleted)
	statusOrderExpired       = newStatus(ErrOrderExpired)
	statusOrderRefunded      = newStatus(ErrOrderRefunded)
	statusOrderPending       = newStatus(ErrOrderPending)
	statusOrderPaying        = newStatus(ErrOrderPaying)

	// 参数/请求错误 → HTTP 400 Bad Request
	statusRequestFail          = newStatus(ErrRequestFail)
	statusParams               = newStatus(ErrParams)
	statusInvalidJson          = newStatus(ErrInvalidJson)
	statusRequestBodyTooLarge  = newStatus(ErrRequestBodyTooLarge)
	statusUnsupportedMediaType = newStatus(ErrUnsupportedMediaType)

	// 参数校验失败 → HTTP 422 Unprocessable Entity
	statusValidateParams = newStatus(ErrValidateParams)

	// 认证错误 → HTTP 401 Unauthorized
	statusAuth             = newStatus(ErrAuthentication)
	statusAuthHeader       = newStatus(ErrAuthenticationHeader)
	statusAppKey           = newStatus(ErrAppKey)
	statusSign             = newStatus(ErrSign)
	statusAuthExpired      = newStatus(ErrExpired)
	statusTokenMissing     = newStatus(ErrTokenMissing)
	statusTokenInvalid     = newStatus(ErrTokenInvalid)
	statusTokenRevoked     = newStatus(ErrTokenRevoked)
	statusNotLoggedIn      = newStatus(ErrNotLoggedIn)
	statusLoginExpired     = newStatus(ErrLoginExpired)
	statusSessionKickedOut = newStatus(ErrSessionKickedOut)

	// 权限不足 → HTTP 403 Forbidden
	statusPermission = newStatus(ErrPermission)

	// 资源不存在 → HTTP 404 Not Found
	statusNotFound       = newStatus(ErrNotFound)
	statusRecordNotFound = newStatus(ErrRecordNotFound)

	// 冲突 → HTTP 409 Conflict
	statusConflict              = newStatus(ErrConflict)
	statusAlreadyExists         = newStatus(ErrAlreadyExists)
	statusUserAlreadyRegistered = newStatus(ErrUserAlreadyRegistered)
	statusVersionConflict       = newStatus(ErrVersionConflict)

	// 限流 → HTTP 429 Too Many Requests
	statusTooManyRequests = newStatus(ErrTooManyRequests)

	// 服务端错误 → HTTP 500 Internal Server Error
	statusInternalServer   = newStatus(ErrInternalServer)
	statusTimeout          = newStatus(ErrTimeout)
	statusCanceled         = newStatus(ErrCanceled)
	statusDeadlineExceeded = newStatus(ErrDeadlineExceeded)

	// 服务不可用 → HTTP 503 Service Unavailable
	statusServiceUnavailable    = newStatus(ErrServiceUnavailable)
	statusDependencyUnavailable = newStatus(ErrDependencyUnavailable)
	statusSystemMaintenance     = newStatus(ErrSystemMaintenance)

	// 中间件/存储错误 → HTTP 500 Internal Server Error
	statusElasticsearchServer        = newStatus(ErrElasticsearchServer)
	statusElasticsearchDSL           = newStatus(ErrElasticsearchDSL)
	statusElasticsearchIndexNotFound = newStatus(ErrElasticsearchIndexNotFound)
	statusElasticsearchBulkPartial   = newStatus(ErrElasticsearchBulkPartial)

	statusMysqlServer      = newStatus(ErrMysqlServer)
	statusMysqlSQL         = newStatus(ErrMysqlSQL)
	statusMysqlQuery       = newStatus(ErrMysqlQuery)
	statusMysqlDeadlock    = newStatus(ErrMysqlDeadlock)
	statusMysqlTransaction = newStatus(ErrMysqlTransaction)

	statusMongoServer = newStatus(ErrMongoServer)
	statusMongoDSL    = newStatus(ErrMongoDSL)
	statusMongoQuery  = newStatus(ErrMongoQuery)

	statusRedisServer            = newStatus(ErrRedisServer)
	statusRedisQuery             = newStatus(ErrRedisQuery)
	statusRedisLockAcquireFailed = newStatus(ErrRedisLockAcquireFailed)
	statusRedisLockLost          = newStatus(ErrRedisLockLost)

	statusKafkaServer         = newStatus(ErrKafkaServer)
	statusKafkaProducer       = newStatus(ErrKafkaProducer)
	statusKafkaConsumer       = newStatus(ErrKafkaConsumer)
	statusKafkaMessageDecode  = newStatus(ErrKafkaMessageDecode)
	statusKafkaRetryExhausted = newStatus(ErrKafkaRetryExhausted)

	statusRabbitMQServer         = newStatus(ErrRabbitMQServer)
	statusRabbitMQProducer       = newStatus(ErrRabbitMQProducer)
	statusRabbitMQConsumer       = newStatus(ErrRabbitMQConsumer)
	statusRabbitMQMessageDecode  = newStatus(ErrRabbitMQMessageDecode)
	statusRabbitMQRetryExhausted = newStatus(ErrRabbitMQRetryExhausted)

	// 支付模块（模块代号 08）
	statusPaymentNotFound           = newStatus(ErrPaymentNotFound)
	statusPaymentAlreadyPaid        = newStatus(ErrPaymentAlreadyPaid)
	statusPaymentClosed             = newStatus(ErrPaymentClosed)
	statusPaymentExpired            = newStatus(ErrPaymentExpired)
	statusPaymentMethodNotSupported = newStatus(ErrPaymentMethodNotSupported)
	statusPayerIdentityMissing      = newStatus(ErrPayerIdentityMissing)
	statusPaymentAmountMismatch     = newStatus(ErrPaymentAmountMismatch)
	statusPaymentSignInvalid        = newStatus(ErrPaymentSignInvalid)
	statusPaymentCallbackInvalid    = newStatus(ErrPaymentCallbackInvalid)
	statusPaymentCreateFailed       = newStatus(ErrPaymentCreateFailed)
	statusPaymentQueryFailed        = newStatus(ErrPaymentQueryFailed)
	statusPaymentChannelUnavailable = newStatus(ErrPaymentChannelUnavailable)
	statusRefundNotFound            = newStatus(ErrRefundNotFound)
	statusRefundAmountExceeded      = newStatus(ErrRefundAmountExceeded)
	statusRefundInProgress          = newStatus(ErrRefundInProgress)
	statusRefundFailed              = newStatus(ErrRefundFailed)

	// 账号与第三方登录模块（模块代号 09）
	statusPhoneAlreadyBound      = newStatus(ErrPhoneAlreadyBound)
	statusThirdPartyAlreadyBound = newStatus(ErrThirdPartyAlreadyBound)
	statusThirdPartyNotBound     = newStatus(ErrThirdPartyNotBound)
	statusThirdPartyAuthFailed   = newStatus(ErrThirdPartyAuthFailed)
	statusThirdPartyUnavailable  = newStatus(ErrThirdPartyUnavailable)
	statusSmsSendTooFrequent     = newStatus(ErrSmsSendTooFrequent)
	statusSmsDailyLimitExceeded  = newStatus(ErrSmsDailyLimitExceeded)
	statusSmsSendFailed          = newStatus(ErrSmsSendFailed)

	// AI 大模型模块（模块代号 10）
	statusLLMAuth              = newStatus(ErrLLMAuth)
	statusLLMRateLimited       = newStatus(ErrLLMRateLimited)
	statusLLMTimeout           = newStatus(ErrLLMTimeout)
	statusLLMContextTooLong    = newStatus(ErrLLMContextTooLong)
	statusLLMContentFiltered   = newStatus(ErrLLMContentFiltered)
	statusLLMInvalidRequest    = newStatus(ErrLLMInvalidRequest)
	statusLLMProviderError     = newStatus(ErrLLMProviderError)
	statusLLMModelNotAvailable = newStatus(ErrLLMModelNotAvailable)
	statusLLMStructuredOutput  = newStatus(ErrLLMStructuredOutput)
	statusLLMStreamInterrupted = newStatus(ErrLLMStreamInterrupted)
	statusLLMPricingMissing    = newStatus(ErrLLMPricingMissing)

	// 文件模块（模块代号 11）
	statusFileEmpty              = newStatus(ErrFileEmpty)
	statusFileTooLarge           = newStatus(ErrFileTooLarge)
	statusFileTypeNotAllowed     = newStatus(ErrFileTypeNotAllowed)
	statusFileCountExceeded      = newStatus(ErrFileCountExceeded)
	statusFileNotFound           = newStatus(ErrFileNotFound)
	statusFileChecksumMismatch   = newStatus(ErrFileChecksumMismatch)
	statusFileUploadFailed       = newStatus(ErrFileUploadFailed)
	statusFileStorageUnavailable = newStatus(ErrFileStorageUnavailable)

	// 邮件模块（模块代号 12）
	statusEmailAlreadyBound       = newStatus(ErrEmailAlreadyBound)
	statusEmailNotVerified        = newStatus(ErrEmailNotVerified)
	statusEmailSendTooFrequent    = newStatus(ErrEmailSendTooFrequent)
	statusEmailTemplateNotFound   = newStatus(ErrEmailTemplateNotFound)
	statusEmailSendFailed         = newStatus(ErrEmailSendFailed)
	statusEmailServiceUnavailable = newStatus(ErrEmailServiceUnavailable)
)

// --- 公开访问函数 ---
// 返回 *Status 指针，指向包级别预构建值，调用零分配。
// 函数分组顺序与上方 var 声明一致。

// 成功 & 业务逻辑错误 → HTTP 200
func StatusOK() *Status                    { return &statusOK }
func StatusVip() *Status                   { return &statusVip }
func StatusDuplicateRequest() *Status      { return &statusDuplicateRequest }
func StatusInsufficientBalance() *Status   { return &statusInsufficientBalance }
func StatusInsufficientStock() *Status     { return &statusInsufficientStock }
func StatusBusinessRuleViolation() *Status { return &statusBusinessRuleViolation }
func StatusAccountDisabled() *Status       { return &statusAccountDisabled }
func StatusOperationNotAllowed() *Status   { return &statusOperationNotAllowed }
func StatusQuotaExceeded() *Status         { return &statusQuotaExceeded }
func StatusUserNotRegistered() *Status     { return &statusUserNotRegistered }
func StatusUserDeactivated() *Status       { return &statusUserDeactivated }
func StatusPasswordIncorrect() *Status     { return &statusPasswordIncorrect }
func StatusCaptchaInvalid() *Status        { return &statusCaptchaInvalid }
func StatusCaptchaExpired() *Status        { return &statusCaptchaExpired }
func StatusAccountLocked() *Status         { return &statusAccountLocked }
func StatusLoginNotPerformed() *Status     { return &statusLoginNotPerformed }
func StatusOperationInProgress() *Status   { return &statusOperationInProgress }
func StatusRiskControlRejected() *Status   { return &statusRiskControlRejected }
func StatusFeatureDisabled() *Status       { return &statusFeatureDisabled }

// 订单模块（模块代号 07）业务态 → HTTP 200
func StatusOrderNotFound() *Status      { return &statusOrderNotFound }
func StatusOrderStatusInvalid() *Status { return &statusOrderStatusInvalid }
func StatusOrderClosed() *Status        { return &statusOrderClosed }
func StatusOrderCanceled() *Status      { return &statusOrderCanceled }
func StatusOrderPaid() *Status          { return &statusOrderPaid }
func StatusOrderCompleted() *Status     { return &statusOrderCompleted }
func StatusOrderExpired() *Status       { return &statusOrderExpired }
func StatusOrderRefunded() *Status      { return &statusOrderRefunded }
func StatusOrderPending() *Status       { return &statusOrderPending }
func StatusOrderPaying() *Status        { return &statusOrderPaying }

// 参数/请求错误 → HTTP 400 Bad Request
func StatusRequestFail() *Status          { return &statusRequestFail }
func StatusParams() *Status               { return &statusParams }
func StatusInvalidJson() *Status          { return &statusInvalidJson }
func StatusRequestBodyTooLarge() *Status  { return &statusRequestBodyTooLarge }
func StatusUnsupportedMediaType() *Status { return &statusUnsupportedMediaType }

// 参数校验失败 → HTTP 422 Unprocessable Entity
func StatusValidateParams() *Status { return &statusValidateParams }

// 认证错误 → HTTP 401 Unauthorized
func StatusAuth() *Status             { return &statusAuth }
func StatusAuthHeader() *Status       { return &statusAuthHeader }
func StatusAppKey() *Status           { return &statusAppKey }
func StatusSign() *Status             { return &statusSign }
func StatusAuthExpired() *Status      { return &statusAuthExpired }
func StatusTokenMissing() *Status     { return &statusTokenMissing }
func StatusTokenInvalid() *Status     { return &statusTokenInvalid }
func StatusTokenRevoked() *Status     { return &statusTokenRevoked }
func StatusNotLoggedIn() *Status      { return &statusNotLoggedIn }
func StatusLoginExpired() *Status     { return &statusLoginExpired }
func StatusSessionKickedOut() *Status { return &statusSessionKickedOut }

// 权限不足 → HTTP 403 Forbidden
func StatusPermission() *Status { return &statusPermission }

// 资源不存在 → HTTP 404 Not Found
func StatusNotFound() *Status       { return &statusNotFound }
func StatusRecordNotFound() *Status { return &statusRecordNotFound }

// 冲突 → HTTP 409 Conflict
func StatusConflict() *Status              { return &statusConflict }
func StatusAlreadyExists() *Status         { return &statusAlreadyExists }
func StatusUserAlreadyRegistered() *Status { return &statusUserAlreadyRegistered }
func StatusVersionConflict() *Status       { return &statusVersionConflict }

// 限流 → HTTP 429 Too Many Requests
func StatusTooManyRequests() *Status { return &statusTooManyRequests }

// 服务端错误 → HTTP 500 Internal Server Error
func StatusInternalServer() *Status   { return &statusInternalServer }
func StatusTimeout() *Status          { return &statusTimeout }
func StatusCanceled() *Status         { return &statusCanceled }
func StatusDeadlineExceeded() *Status { return &statusDeadlineExceeded }

// 服务不可用 → HTTP 503 Service Unavailable
func StatusServiceUnavailable() *Status    { return &statusServiceUnavailable }
func StatusDependencyUnavailable() *Status { return &statusDependencyUnavailable }
func StatusSystemMaintenance() *Status     { return &statusSystemMaintenance }

// 中间件/存储错误 → HTTP 500 Internal Server Error
func StatusElasticsearchServer() *Status        { return &statusElasticsearchServer }
func StatusElasticsearchDSL() *Status           { return &statusElasticsearchDSL }
func StatusElasticsearchIndexNotFound() *Status { return &statusElasticsearchIndexNotFound }
func StatusElasticsearchBulkPartial() *Status   { return &statusElasticsearchBulkPartial }

func StatusMysqlServer() *Status      { return &statusMysqlServer }
func StatusMysqlSQL() *Status         { return &statusMysqlSQL }
func StatusMysqlQuery() *Status       { return &statusMysqlQuery }
func StatusMysqlDeadlock() *Status    { return &statusMysqlDeadlock }
func StatusMysqlTransaction() *Status { return &statusMysqlTransaction }

func StatusMongoServer() *Status { return &statusMongoServer }
func StatusMongoDSL() *Status    { return &statusMongoDSL }
func StatusMongoQuery() *Status  { return &statusMongoQuery }

func StatusRedisServer() *Status            { return &statusRedisServer }
func StatusRedisQuery() *Status             { return &statusRedisQuery }
func StatusRedisLockAcquireFailed() *Status { return &statusRedisLockAcquireFailed }
func StatusRedisLockLost() *Status          { return &statusRedisLockLost }

func StatusKafkaServer() *Status         { return &statusKafkaServer }
func StatusKafkaProducer() *Status       { return &statusKafkaProducer }
func StatusKafkaConsumer() *Status       { return &statusKafkaConsumer }
func StatusKafkaMessageDecode() *Status  { return &statusKafkaMessageDecode }
func StatusKafkaRetryExhausted() *Status { return &statusKafkaRetryExhausted }

func StatusRabbitMQServer() *Status         { return &statusRabbitMQServer }
func StatusRabbitMQProducer() *Status       { return &statusRabbitMQProducer }
func StatusRabbitMQConsumer() *Status       { return &statusRabbitMQConsumer }
func StatusRabbitMQMessageDecode() *Status  { return &statusRabbitMQMessageDecode }
func StatusRabbitMQRetryExhausted() *Status { return &statusRabbitMQRetryExhausted }

// 支付 → 业务态 200，金额/回调参数 400，验签 401，发起/查询/退款失败 500，渠道不可用 503
func StatusPaymentNotFound() *Status           { return &statusPaymentNotFound }
func StatusPaymentAlreadyPaid() *Status        { return &statusPaymentAlreadyPaid }
func StatusPaymentClosed() *Status             { return &statusPaymentClosed }
func StatusPaymentExpired() *Status            { return &statusPaymentExpired }
func StatusPaymentMethodNotSupported() *Status { return &statusPaymentMethodNotSupported }
func StatusPayerIdentityMissing() *Status      { return &statusPayerIdentityMissing }
func StatusPaymentAmountMismatch() *Status     { return &statusPaymentAmountMismatch }
func StatusPaymentSignInvalid() *Status        { return &statusPaymentSignInvalid }
func StatusPaymentCallbackInvalid() *Status    { return &statusPaymentCallbackInvalid }
func StatusPaymentCreateFailed() *Status       { return &statusPaymentCreateFailed }
func StatusPaymentQueryFailed() *Status        { return &statusPaymentQueryFailed }
func StatusPaymentChannelUnavailable() *Status { return &statusPaymentChannelUnavailable }
func StatusRefundNotFound() *Status            { return &statusRefundNotFound }
func StatusRefundAmountExceeded() *Status      { return &statusRefundAmountExceeded }
func StatusRefundInProgress() *Status          { return &statusRefundInProgress }
func StatusRefundFailed() *Status              { return &statusRefundFailed }

// 账号与第三方登录 → 绑定冲突 409，业务态 200，第三方不可用 503，发送失败 500
func StatusPhoneAlreadyBound() *Status      { return &statusPhoneAlreadyBound }
func StatusThirdPartyAlreadyBound() *Status { return &statusThirdPartyAlreadyBound }
func StatusThirdPartyNotBound() *Status     { return &statusThirdPartyNotBound }
func StatusThirdPartyAuthFailed() *Status   { return &statusThirdPartyAuthFailed }
func StatusThirdPartyUnavailable() *Status  { return &statusThirdPartyUnavailable }
func StatusSmsSendTooFrequent() *Status     { return &statusSmsSendTooFrequent }
func StatusSmsDailyLimitExceeded() *Status  { return &statusSmsDailyLimitExceeded }
func StatusSmsSendFailed() *Status          { return &statusSmsSendFailed }

// AI 大模型 → 用户可修正 200，我方配置/实现问题 500，供应商侧故障 503
func StatusLLMAuth() *Status              { return &statusLLMAuth }
func StatusLLMRateLimited() *Status       { return &statusLLMRateLimited }
func StatusLLMTimeout() *Status           { return &statusLLMTimeout }
func StatusLLMContextTooLong() *Status    { return &statusLLMContextTooLong }
func StatusLLMContentFiltered() *Status   { return &statusLLMContentFiltered }
func StatusLLMInvalidRequest() *Status    { return &statusLLMInvalidRequest }
func StatusLLMProviderError() *Status     { return &statusLLMProviderError }
func StatusLLMModelNotAvailable() *Status { return &statusLLMModelNotAvailable }
func StatusLLMStructuredOutput() *Status  { return &statusLLMStructuredOutput }
func StatusLLMStreamInterrupted() *Status { return &statusLLMStreamInterrupted }
func StatusLLMPricingMissing() *Status    { return &statusLLMPricingMissing }

// 文件 → 校验类 400，不存在 404，上传失败 500，存储不可用 503
func StatusFileEmpty() *Status              { return &statusFileEmpty }
func StatusFileTooLarge() *Status           { return &statusFileTooLarge }
func StatusFileTypeNotAllowed() *Status     { return &statusFileTypeNotAllowed }
func StatusFileCountExceeded() *Status      { return &statusFileCountExceeded }
func StatusFileNotFound() *Status           { return &statusFileNotFound }
func StatusFileChecksumMismatch() *Status   { return &statusFileChecksumMismatch }
func StatusFileUploadFailed() *Status       { return &statusFileUploadFailed }
func StatusFileStorageUnavailable() *Status { return &statusFileStorageUnavailable }

// 邮件 → 绑定冲突 409，业务态 200，发送失败 500，服务不可用 503
func StatusEmailAlreadyBound() *Status       { return &statusEmailAlreadyBound }
func StatusEmailNotVerified() *Status        { return &statusEmailNotVerified }
func StatusEmailSendTooFrequent() *Status    { return &statusEmailSendTooFrequent }
func StatusEmailTemplateNotFound() *Status   { return &statusEmailTemplateNotFound }
func StatusEmailSendFailed() *Status         { return &statusEmailSendFailed }
func StatusEmailServiceUnavailable() *Status { return &statusEmailServiceUnavailable }
