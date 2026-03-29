package goerr

import "net/http"

// Status 错误状态，包含业务错误码、HTTP 状态码和消息。
// Status 是不可变值类型，所有实例在包初始化时预构建，运行时零分配。
type Status struct {
	code     Code
	httpCode int
	msg      string
}

const defaultBusinessHTTPCode = http.StatusOK

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
// 业务错误默认使用 HTTP 200，客户端通过业务 code 判断失败。
// 调用方通过 StatusXxx() 获取指针，指向包级别变量，零堆分配。
var (
	statusOK                    = newStatus(ErrNo, defaultBusinessHTTPCode)
	statusRequestFail           = newStatus(ErrRequestFail, defaultBusinessHTTPCode)
	statusInternalServer        = newStatus(ErrInternalServer, defaultBusinessHTTPCode)
	statusParams                = newStatus(ErrParams, defaultBusinessHTTPCode)
	statusValidateParams        = newStatus(ErrValidateParams, defaultBusinessHTTPCode)
	statusAuth                  = newStatus(ErrAuthentication, defaultBusinessHTTPCode)
	statusVip                   = newStatus(ErrVipRights, defaultBusinessHTTPCode)
	statusNotFound              = newStatus(ErrNotFound, defaultBusinessHTTPCode)
	statusAuthHeader            = newStatus(ErrAuthenticationHeader, defaultBusinessHTTPCode)
	statusAppKey                = newStatus(ErrAppKey, defaultBusinessHTTPCode)
	statusSign                  = newStatus(ErrSign, defaultBusinessHTTPCode)
	statusPermission            = newStatus(ErrPermission, defaultBusinessHTTPCode)
	statusTooManyRequests       = newStatus(ErrTooManyRequests, defaultBusinessHTTPCode)
	statusInvalidJson           = newStatus(ErrInvalidJson, defaultBusinessHTTPCode)
	statusTimeout               = newStatus(ErrTimeout, defaultBusinessHTTPCode)
	statusAuthExpired           = newStatus(ErrExpired, defaultBusinessHTTPCode)
	statusConflict              = newStatus(ErrConflict, defaultBusinessHTTPCode)
	statusAlreadyExists         = newStatus(ErrAlreadyExists, defaultBusinessHTTPCode)
	statusServiceUnavailable    = newStatus(ErrServiceUnavailable, defaultBusinessHTTPCode)
	statusDependencyUnavailable = newStatus(ErrDependencyUnavailable, defaultBusinessHTTPCode)
	statusCanceled              = newStatus(ErrCanceled, defaultBusinessHTTPCode)
	statusDeadlineExceeded      = newStatus(ErrDeadlineExceeded, defaultBusinessHTTPCode)
	statusTokenMissing          = newStatus(ErrTokenMissing, defaultBusinessHTTPCode)
	statusTokenInvalid          = newStatus(ErrTokenInvalid, defaultBusinessHTTPCode)
	statusTokenRevoked          = newStatus(ErrTokenRevoked, defaultBusinessHTTPCode)
	statusDuplicateRequest      = newStatus(ErrDuplicateRequest, defaultBusinessHTTPCode)
	statusRecordNotFound        = newStatus(ErrRecordNotFound, defaultBusinessHTTPCode)

	statusElasticsearchServer = newStatus(ErrElasticsearchServer, defaultBusinessHTTPCode)
	statusElasticsearchDSL    = newStatus(ErrElasticsearchDSL, defaultBusinessHTTPCode)

	statusMysqlServer = newStatus(ErrMysqlServer, defaultBusinessHTTPCode)
	statusMysqlSQL    = newStatus(ErrMysqlSQL, defaultBusinessHTTPCode)
	statusMysqlQuery  = newStatus(ErrMysqlQuery, defaultBusinessHTTPCode)

	statusMongoServer = newStatus(ErrMongoServer, defaultBusinessHTTPCode)
	statusMongoDSL    = newStatus(ErrMongoDSL, defaultBusinessHTTPCode)
	statusMongoQuery  = newStatus(ErrMongoQuery, defaultBusinessHTTPCode)

	statusRedisServer = newStatus(ErrRedisServer, defaultBusinessHTTPCode)
	statusRedisQuery  = newStatus(ErrRedisQuery, defaultBusinessHTTPCode)

	statusKafkaServer   = newStatus(ErrKafkaServer, defaultBusinessHTTPCode)
	statusKafkaProducer = newStatus(ErrKafkaProducer, defaultBusinessHTTPCode)
	statusKafkaConsumer = newStatus(ErrKafkaConsumer, defaultBusinessHTTPCode)

	statusRabbitMQServer   = newStatus(ErrRabbitMQServer, defaultBusinessHTTPCode)
	statusRabbitMQProducer = newStatus(ErrRabbitMQProducer, defaultBusinessHTTPCode)
	statusRabbitMQConsumer = newStatus(ErrRabbitMQConsumer, defaultBusinessHTTPCode)
)

// --- 公开访问函数 ---
// 返回 *Status 指针，指向包级别预构建值，调用零分配。

func StatusOK() *Status              { return &statusOK }
func StatusRequestFail() *Status     { return &statusRequestFail }
func StatusInternalServer() *Status  { return &statusInternalServer }
func StatusParams() *Status          { return &statusParams }
func StatusValidateParams() *Status  { return &statusValidateParams }
func StatusAuth() *Status            { return &statusAuth }
func StatusVip() *Status             { return &statusVip }
func StatusNotFound() *Status        { return &statusNotFound }
func StatusAuthHeader() *Status      { return &statusAuthHeader }
func StatusAppKey() *Status          { return &statusAppKey }
func StatusSign() *Status            { return &statusSign }
func StatusPermission() *Status      { return &statusPermission }
func StatusTooManyRequests() *Status { return &statusTooManyRequests }
func StatusInvalidJson() *Status     { return &statusInvalidJson }
func StatusTimeout() *Status         { return &statusTimeout }
func StatusAuthExpired() *Status     { return &statusAuthExpired }
func StatusConflict() *Status        { return &statusConflict }
func StatusAlreadyExists() *Status   { return &statusAlreadyExists }
func StatusServiceUnavailable() *Status {
	return &statusServiceUnavailable
}
func StatusDependencyUnavailable() *Status { return &statusDependencyUnavailable }
func StatusCanceled() *Status              { return &statusCanceled }
func StatusDeadlineExceeded() *Status      { return &statusDeadlineExceeded }
func StatusTokenMissing() *Status          { return &statusTokenMissing }
func StatusTokenInvalid() *Status          { return &statusTokenInvalid }
func StatusTokenRevoked() *Status          { return &statusTokenRevoked }
func StatusDuplicateRequest() *Status      { return &statusDuplicateRequest }
func StatusRecordNotFound() *Status        { return &statusRecordNotFound }

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

// NewStatus 创建自定义 Status，用于业务方扩展错误码。
// code 必须遵循 项目组(10)+服务(01)+模块(0~99)+错误码(0~99) 的规范。
// 统一响应体内的业务错误建议使用 http.StatusOK；仅协议级场景再使用非 200。
// 非法 code 会 panic；若需要预检查，可先调用 ValidateCode。
func NewStatus(code Code, httpCode int, msg string) *Status {
	mustValidCode(code)
	return &Status{code: code, httpCode: httpCode, msg: msg}
}
