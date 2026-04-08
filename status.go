package goerr

import "net/http"

// Status 错误状态，包含业务错误码、HTTP 状态码和消息。
// Status 是不可变值类型，所有实例在包初始化时预构建，运行时零分配。
type Status struct {
	code     Code
	httpCode int
	msg      string
}

// NewStatus 创建自定义 Status，用于业务方扩展错误码。
// code 不做强制校验，业务方可自由定义编码规范。
// 如需校验是否符合本库默认编码规范，可先调用 ValidateCode。
// 统一响应体内的业务错误建议使用 http.StatusOK；仅协议级场景再使用非 200。
func NewStatus(code Code, httpCode int, msg string) *Status {
	return &Status{code: code, httpCode: httpCode, msg: msg}
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

// newStatus 构造 Status，msg 取 Code 默认消息。
func newStatus(code Code, httpCode int) Status {
	mustValidCode(code)
	return Status{code: code, httpCode: httpCode, msg: code.Message()}
}

// 预构建的不可变 Status 值。
// HTTP 状态码语义化：业务逻辑错误 200，参数错误 400/422，认证 401，
// 权限 403，资源不存在 404，冲突 409，限流 429，服务端错误 500。
// 调用方通过 StatusXxx() 获取指针，指向包级别变量，零堆分配。
var (
	// 成功 & 业务逻辑错误 → HTTP 200
	statusOK                    = newStatus(ErrNo, http.StatusOK)
	statusVip                   = newStatus(ErrVipRights, http.StatusOK)
	statusDuplicateRequest      = newStatus(ErrDuplicateRequest, http.StatusOK)
	statusInsufficientBalance   = newStatus(ErrInsufficientBalance, http.StatusOK)
	statusInsufficientStock     = newStatus(ErrInsufficientStock, http.StatusOK)
	statusBusinessRuleViolation = newStatus(ErrBusinessRuleViolation, http.StatusOK)
	statusAccountDisabled       = newStatus(ErrAccountDisabled, http.StatusOK)
	statusOperationNotAllowed   = newStatus(ErrOperationNotAllowed, http.StatusOK)
	statusQuotaExceeded         = newStatus(ErrQuotaExceeded, http.StatusOK)

	// 参数/请求错误 → HTTP 400 Bad Request
	statusRequestFail = newStatus(ErrRequestFail, http.StatusBadRequest)
	statusParams      = newStatus(ErrParams, http.StatusBadRequest)
	statusInvalidJson = newStatus(ErrInvalidJson, http.StatusBadRequest)

	// 参数校验失败 → HTTP 422 Unprocessable Entity
	statusValidateParams = newStatus(ErrValidateParams, http.StatusUnprocessableEntity)

	// 认证错误 → HTTP 401 Unauthorized
	statusAuth         = newStatus(ErrAuthentication, http.StatusUnauthorized)
	statusAuthHeader   = newStatus(ErrAuthenticationHeader, http.StatusUnauthorized)
	statusAppKey       = newStatus(ErrAppKey, http.StatusUnauthorized)
	statusSign         = newStatus(ErrSign, http.StatusUnauthorized)
	statusAuthExpired  = newStatus(ErrExpired, http.StatusUnauthorized)
	statusTokenMissing = newStatus(ErrTokenMissing, http.StatusUnauthorized)
	statusTokenInvalid = newStatus(ErrTokenInvalid, http.StatusUnauthorized)
	statusTokenRevoked = newStatus(ErrTokenRevoked, http.StatusUnauthorized)

	// 权限不足 → HTTP 403 Forbidden
	statusPermission = newStatus(ErrPermission, http.StatusForbidden)

	// 资源不存在 → HTTP 404 Not Found
	statusNotFound       = newStatus(ErrNotFound, http.StatusNotFound)
	statusRecordNotFound = newStatus(ErrRecordNotFound, http.StatusNotFound)

	// 冲突 → HTTP 409 Conflict
	statusConflict      = newStatus(ErrConflict, http.StatusConflict)
	statusAlreadyExists = newStatus(ErrAlreadyExists, http.StatusConflict)

	// 限流 → HTTP 429 Too Many Requests
	statusTooManyRequests = newStatus(ErrTooManyRequests, http.StatusTooManyRequests)

	// 服务端错误 → HTTP 500 Internal Server Error
	statusInternalServer   = newStatus(ErrInternalServer, http.StatusInternalServerError)
	statusTimeout          = newStatus(ErrTimeout, http.StatusInternalServerError)
	statusCanceled         = newStatus(ErrCanceled, http.StatusInternalServerError)
	statusDeadlineExceeded = newStatus(ErrDeadlineExceeded, http.StatusInternalServerError)

	// 服务不可用 → HTTP 503 Service Unavailable
	statusServiceUnavailable    = newStatus(ErrServiceUnavailable, http.StatusServiceUnavailable)
	statusDependencyUnavailable = newStatus(ErrDependencyUnavailable, http.StatusServiceUnavailable)

	// 中间件/存储错误 → HTTP 500 Internal Server Error
	statusElasticsearchServer = newStatus(ErrElasticsearchServer, http.StatusInternalServerError)
	statusElasticsearchDSL    = newStatus(ErrElasticsearchDSL, http.StatusInternalServerError)

	statusMysqlServer = newStatus(ErrMysqlServer, http.StatusInternalServerError)
	statusMysqlSQL    = newStatus(ErrMysqlSQL, http.StatusInternalServerError)
	statusMysqlQuery  = newStatus(ErrMysqlQuery, http.StatusInternalServerError)

	statusMongoServer = newStatus(ErrMongoServer, http.StatusInternalServerError)
	statusMongoDSL    = newStatus(ErrMongoDSL, http.StatusInternalServerError)
	statusMongoQuery  = newStatus(ErrMongoQuery, http.StatusInternalServerError)

	statusRedisServer = newStatus(ErrRedisServer, http.StatusInternalServerError)
	statusRedisQuery  = newStatus(ErrRedisQuery, http.StatusInternalServerError)

	statusKafkaServer   = newStatus(ErrKafkaServer, http.StatusInternalServerError)
	statusKafkaProducer = newStatus(ErrKafkaProducer, http.StatusInternalServerError)
	statusKafkaConsumer = newStatus(ErrKafkaConsumer, http.StatusInternalServerError)

	statusRabbitMQServer   = newStatus(ErrRabbitMQServer, http.StatusInternalServerError)
	statusRabbitMQProducer = newStatus(ErrRabbitMQProducer, http.StatusInternalServerError)
	statusRabbitMQConsumer = newStatus(ErrRabbitMQConsumer, http.StatusInternalServerError)
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

// 参数/请求错误 → HTTP 400 Bad Request
func StatusRequestFail() *Status { return &statusRequestFail }
func StatusParams() *Status      { return &statusParams }
func StatusInvalidJson() *Status { return &statusInvalidJson }

// 参数校验失败 → HTTP 422 Unprocessable Entity
func StatusValidateParams() *Status { return &statusValidateParams }

// 认证错误 → HTTP 401 Unauthorized
func StatusAuth() *Status         { return &statusAuth }
func StatusAuthHeader() *Status   { return &statusAuthHeader }
func StatusAppKey() *Status       { return &statusAppKey }
func StatusSign() *Status         { return &statusSign }
func StatusAuthExpired() *Status  { return &statusAuthExpired }
func StatusTokenMissing() *Status { return &statusTokenMissing }
func StatusTokenInvalid() *Status { return &statusTokenInvalid }
func StatusTokenRevoked() *Status { return &statusTokenRevoked }

// 权限不足 → HTTP 403 Forbidden
func StatusPermission() *Status { return &statusPermission }

// 资源不存在 → HTTP 404 Not Found
func StatusNotFound() *Status       { return &statusNotFound }
func StatusRecordNotFound() *Status { return &statusRecordNotFound }

// 冲突 → HTTP 409 Conflict
func StatusConflict() *Status      { return &statusConflict }
func StatusAlreadyExists() *Status { return &statusAlreadyExists }

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
