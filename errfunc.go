package goerr

import (
	"net/http"
)

// setError error code.
// code 自定义错误码.
// httpCode HTTP状态码.
// msg 自定义错误信息.
func setError(code, httpCode int) ErrStatuser {
	return &ErrStatus{
		errCode:  code,
		httpCode: httpCode,
		msg:      statusText(code),
	}
}

// Ok 请求成功.
func Ok() ErrStatuser {
	return setError(ErrNo, http.StatusOK)
}

// RequestFail 请求失败.
func RequestFail() ErrStatuser {
	return setError(ErrRequestFail, http.StatusBadRequest)
}

// InternalServer 服务器内部错误.
func InternalServer() ErrStatuser {
	return setError(ErrInternalServer, http.StatusInternalServerError)
}

// Params 非法参数.
func Params() ErrStatuser {
	return setError(ErrParams, http.StatusUnprocessableEntity)
}

// ValidateParams 验证参数.
func ValidateParams() ErrStatuser {
	return setError(ErrValidateParams, http.StatusUnprocessableEntity)
}

// Auth 身份验证失败.
func Auth() ErrStatuser {
	return setError(ErrAuthentication, http.StatusUnauthorized)
}

func Vip() ErrStatuser {
	return setError(ErrVipRights, http.StatusUnauthorized)
}

// NotFound 找不到.
func NotFound() ErrStatuser {
	return setError(ErrNotFound, http.StatusNotFound)
}

// AuthHeader 认证头非法.
func AuthHeader() ErrStatuser {
	return setError(ErrAuthenticationHeader, http.StatusForbidden)
}

// AppKey 无效应用程序密钥.
func AppKey() ErrStatuser {
	return setError(ErrAppKey, http.StatusForbidden)
}

// Sign 无效应用签名.
func Sign() ErrStatuser {
	return setError(ErrSign, http.StatusUnauthorized)
}

// Permission 权限不足.
func Permission() ErrStatuser {
	return setError(ErrPermission, http.StatusForbidden)
}

// TooManyRequests 请求过多.
func TooManyRequests() ErrStatuser {
	return setError(ErrTooManyRequests, http.StatusTooManyRequests)
}

// InvalidJson 无效的JSON.
func InvalidJson() ErrStatuser {
	return setError(ErrInvalidJson, http.StatusUnprocessableEntity)
}

// Timeout 请求超时.
func Timeout() ErrStatuser {
	return setError(ErrTimeout, http.StatusRequestTimeout)
}

// AuthExpired 认证过期.
func AuthExpired() ErrStatuser {
	return setError(ErrExpired, http.StatusGatewayTimeout)
}

// ElasticsearchServer elasticsearch 服务器错误.
func ElasticsearchServer() ErrStatuser {
	return setError(ErrElasticsearchServer, http.StatusInternalServerError)
}

// ElasticsearchDSL elasticsearch DSL 错误.
func ElasticsearchDSL() ErrStatuser {
	return setError(ErrElasticsearchDSL, http.StatusInternalServerError)
}

// MysqlServer mysql 服务器错误.
func MysqlServer() ErrStatuser {
	return setError(ErrMysqlServer, http.StatusInternalServerError)
}

// MysqlSQL mysql SQL 错误.
func MysqlSQL() ErrStatuser {
	return setError(ErrMysqlSQL, http.StatusInternalServerError)
}

// MongoServer mongo 服务器错误.
func MongoServer() ErrStatuser {
	return setError(ErrMongoServer, http.StatusInternalServerError)
}

// MongoDSL mongo DSL 错误.
func MongoDSL() ErrStatuser {
	return setError(ErrMongoDSL, http.StatusInternalServerError)
}

// RedisServer redis 服务器错误.
func RedisServer() ErrStatuser {
	return setError(ErrRedisServer, http.StatusInternalServerError)
}

// KafkaServer kafka 服务器错误.
func KafkaServer() ErrStatuser {
	return setError(ErrKafkaServer, http.StatusInternalServerError)
}

// KafkaProducer kafka 生产者错误.
func KafkaProducer() ErrStatuser {
	return setError(ErrKafkaProducer, http.StatusInternalServerError)
}

// KafkaConsumer kafka 消费者错误.
func KafkaConsumer() ErrStatuser {
	return setError(ErrKafkaConsumer, http.StatusInternalServerError)
}

// RabbitMQServer rabbitmq 服务器错误.
func RabbitMQServer() ErrStatuser {
	return setError(ErrRabbitMQServer, http.StatusInternalServerError)
}

// RabbitMQProducer rabbitmq 生产者错误
func RabbitMQProducer() ErrStatuser {
	return setError(ErrRabbitMQProducer, http.StatusInternalServerError)
}

// RabbitMQConsumer rabbitmq 消费者错误
func RabbitMQConsumer() ErrStatuser {
	return setError(ErrRabbitMQConsumer, http.StatusInternalServerError)
}
