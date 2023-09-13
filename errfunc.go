package goerr

// ErrCode error code
func SetError(code, httpcode int, msg string) ErrCode {
	return ErrCode{
		Code:     code,
		HTTPCode: httpcode,
		Desc:     msg,
	}
}

// Ok 请求成功
func Ok(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrNo
}

// RequestFail 请求失败
func RequestFail(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRequestFail
}

// InternalServer 服务器内部错误
func InternalServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrInternalServer
}

// Params 非法参数
func Params(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrParams
}

// ValidateParams 验证参数
func ValidateParams(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrValidateParams
}

// Authentication 身份验证失败
func Authentication(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAuthentication
}

func VipRights(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrVipRights
}

// NotFound 找不到
func NotFound(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrNotFound
}

// AuthenticationHeader 认证头非法
func AuthenticationHeader(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAuthenticationHeader
}

// AppKey 无效应用程序密钥
func AppKey(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAppKey
}

// Sign 无效应用签名
func Sign(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrSign
}

// Permission 权限不足
func Permission(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrPermission
}

// TooManyRequests 请求过多
func TooManyRequests(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrTooManyRequests
}

// InvalidJson 无效的JSON
func InvalidJson(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrInvalidJson
}

// Timeout 请求超时
func Timeout(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrTimeout
}

// AuthExpired 认证过期
func AuthExpired(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAuthExpired
}

// ElasticsearchServer elasticsearch 服务器错误
func ElasticsearchServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrElasticsearchServer
}

// ElasticsearchDSL elasticsearch DSL 错误
func ElasticsearchDSL(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrElasticsearchDSL
}

// MysqlServer mysql 服务器错误
func MysqlServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMysqlServer
}

// MysqlSQL mysql SQL 错误
func MysqlSQL(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMysqlSQL
}

// MongoServer mongo 服务器错误
func MongoServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMongoServer
}

// MongoDSL mongo DSL 错误
func MongoDSL(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMongoDSL
}

// RedisServer redis 服务器错误
func RedisServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRedisServer
}

// KafkaServer kafka 服务器错误
func KafkaServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrKafkaServer
}

// KafkaProducer kafka 生产者错误
func KafkaProducer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrKafkaProducer
}

// KafkaConsumer kafka 消费者错误
func KafkaConsumer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrKafkaConsumer
}

// RabbitMQServer rabbitmq 服务器错误
func RabbitMQServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRabbitMQServer
}

// RabbitMQProducer rabbitmq 生产者错误
func RabbitMQProducer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRabbitMQProducer
}

// RabbitMQConsumer rabbitmq 消费者错误
func RabbitMQConsumer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRabbitMQConsumer
}
