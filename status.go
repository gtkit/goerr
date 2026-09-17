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
	statusRequestFail = newStatus(ErrRequestFail)
	statusParams      = newStatus(ErrParams)
	statusInvalidJson = newStatus(ErrInvalidJson)

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

	// 中间件/存储错误 → HTTP 500 Internal Server Error
	statusElasticsearchServer = newStatus(ErrElasticsearchServer)
	statusElasticsearchDSL    = newStatus(ErrElasticsearchDSL)

	statusMysqlServer = newStatus(ErrMysqlServer)
	statusMysqlSQL    = newStatus(ErrMysqlSQL)
	statusMysqlQuery  = newStatus(ErrMysqlQuery)

	statusMongoServer = newStatus(ErrMongoServer)
	statusMongoDSL    = newStatus(ErrMongoDSL)
	statusMongoQuery  = newStatus(ErrMongoQuery)

	statusRedisServer = newStatus(ErrRedisServer)
	statusRedisQuery  = newStatus(ErrRedisQuery)

	statusKafkaServer   = newStatus(ErrKafkaServer)
	statusKafkaProducer = newStatus(ErrKafkaProducer)
	statusKafkaConsumer = newStatus(ErrKafkaConsumer)

	statusRabbitMQServer   = newStatus(ErrRabbitMQServer)
	statusRabbitMQProducer = newStatus(ErrRabbitMQProducer)
	statusRabbitMQConsumer = newStatus(ErrRabbitMQConsumer)
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
func StatusRequestFail() *Status { return &statusRequestFail }
func StatusParams() *Status      { return &statusParams }
func StatusInvalidJson() *Status { return &statusInvalidJson }

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

// 中间件/存储错误 → HTTP 500 Internal Server Error
func StatusElasticsearchServer() *Status { return &statusElasticsearchServer }
func StatusElasticsearchDSL() *Status    { return &statusElasticsearchDSL }

func StatusMysqlServer() *Status { return &statusMysqlServer }
func StatusMysqlSQL() *Status    { return &statusMysqlSQL }
func StatusMysqlQuery() *Status  { return &statusMysqlQuery }

func StatusMongoServer() *Status { return &statusMongoServer }
func StatusMongoDSL() *Status    { return &statusMongoDSL }
func StatusMongoQuery() *Status  { return &statusMongoQuery }

func StatusRedisServer() *Status { return &statusRedisServer }
func StatusRedisQuery() *Status  { return &statusRedisQuery }

func StatusKafkaServer() *Status   { return &statusKafkaServer }
func StatusKafkaProducer() *Status { return &statusKafkaProducer }
func StatusKafkaConsumer() *Status { return &statusKafkaConsumer }

func StatusRabbitMQServer() *Status   { return &statusRabbitMQServer }
func StatusRabbitMQProducer() *Status { return &statusRabbitMQProducer }
func StatusRabbitMQConsumer() *Status { return &statusRabbitMQConsumer }
