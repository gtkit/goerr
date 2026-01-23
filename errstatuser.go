package goerr

import (
	"net/http"
)

type ErrStatus struct {
	errCode  string
	httpCode int
	msg      string
}

func (e *ErrStatus) ErrCode() string {
	return e.errCode
}

func (e *ErrStatus) HTTPCode() int {
	return e.httpCode
}
func (e *ErrStatus) Msg() string {
	return e.msg
}

// setError error code.
// code 自定义错误码.
// httpCode HTTP状态码.
// msg 自定义错误信息.
func setError(code string, httpCode int) *ErrStatus {
	return &ErrStatus{
		errCode:  code,
		httpCode: httpCode,
		msg:      codeMessages[code],
	}
}

// Ok 请求成功.
func Ok() *ErrStatus {
	return setError(ErrNo, http.StatusOK)
}

// RequestFail 请求失败.
func RequestFail() *ErrStatus {
	return setError(ErrRequestFail, http.StatusBadRequest)
}

// InternalServer 服务器内部错误.
func InternalServer() *ErrStatus {
	return setError(ErrInternalServer, http.StatusInternalServerError)
}

// Params 非法参数.
func Params() *ErrStatus {
	return setError(ErrParams, http.StatusUnprocessableEntity)
}

// ValidateParams 验证参数.
func ValidateParams() *ErrStatus {
	return setError(ErrValidateParams, http.StatusUnprocessableEntity)
}

// Auth 身份验证失败.
func Auth() *ErrStatus {
	return setError(ErrAuthentication, http.StatusUnauthorized)
}

func Vip() *ErrStatus {
	return setError(ErrVipRights, http.StatusUnauthorized)
}

// NotFound 找不到.
func NotFound() *ErrStatus {
	return setError(ErrNotFound, http.StatusNotFound)
}

// AuthHeader 认证头非法.
func AuthHeader() *ErrStatus {
	return setError(ErrAuthenticationHeader, http.StatusForbidden)
}

// AppKey 无效应用程序密钥.
func AppKey() *ErrStatus {
	return setError(ErrAppKey, http.StatusForbidden)
}

// Sign 无效应用签名.
func Sign() *ErrStatus {
	return setError(ErrSign, http.StatusUnauthorized)
}

// Permission 权限不足.
func Permission() *ErrStatus {
	return setError(ErrPermission, http.StatusForbidden)
}

// TooManyRequests 请求过多.
func TooManyRequests() *ErrStatus {
	return setError(ErrTooManyRequests, http.StatusTooManyRequests)
}

// InvalidJson 无效的JSON.
func InvalidJson() *ErrStatus {
	return setError(ErrInvalidJson, http.StatusUnprocessableEntity)
}

// Timeout 请求超时.
func Timeout() *ErrStatus {
	return setError(ErrTimeout, http.StatusRequestTimeout)
}

// AuthExpired 认证过期.
func AuthExpired() *ErrStatus {
	return setError(ErrExpired, http.StatusGatewayTimeout)
}

// ElasticsearchServer elasticsearch 服务器错误.
func ElasticsearchServer() *ErrStatus {
	return setError(ErrElasticsearchServer, http.StatusInternalServerError)
}

// ElasticsearchDSL elasticsearch DSL 错误.
func ElasticsearchDSL() *ErrStatus {
	return setError(ErrElasticsearchDSL, http.StatusInternalServerError)
}

// MysqlServer mysql 服务器错误.
func MysqlServer() *ErrStatus {
	return setError(ErrMysqlServer, http.StatusInternalServerError)
}

// MysqlSQL mysql SQL 错误.
func MysqlSQL() *ErrStatus {
	return setError(ErrMysqlSQL, http.StatusInternalServerError)
}

// MongoServer mongo 服务器错误.
func MongoServer() *ErrStatus {
	return setError(ErrMongoServer, http.StatusInternalServerError)
}

// MongoDSL mongo DSL 错误.
func MongoDSL() *ErrStatus {
	return setError(ErrMongoDSL, http.StatusInternalServerError)
}

// RedisServer redis 服务器错误.
func RedisServer() *ErrStatus {
	return setError(ErrRedisServer, http.StatusInternalServerError)
}

// KafkaServer kafka 服务器错误.
func KafkaServer() *ErrStatus {
	return setError(ErrKafkaServer, http.StatusInternalServerError)
}

// KafkaProducer kafka 生产者错误.
func KafkaProducer() *ErrStatus {
	return setError(ErrKafkaProducer, http.StatusInternalServerError)
}

// KafkaConsumer kafka 消费者错误.
func KafkaConsumer() *ErrStatus {
	return setError(ErrKafkaConsumer, http.StatusInternalServerError)
}

// RabbitMQServer rabbitmq 服务器错误.
func RabbitMQServer() *ErrStatus {
	return setError(ErrRabbitMQServer, http.StatusInternalServerError)
}

// RabbitMQProducer rabbitmq 生产者错误
func RabbitMQProducer() *ErrStatus {
	return setError(ErrRabbitMQProducer, http.StatusInternalServerError)
}

// RabbitMQConsumer rabbitmq 消费者错误
func RabbitMQConsumer() *ErrStatus {
	return setError(ErrRabbitMQConsumer, http.StatusInternalServerError)
}
