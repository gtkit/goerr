package goerr

func SetError(code, httpcode int, msg string) ErrCode {
	return ErrCode{
		Code:     code,
		HTTPCode: httpcode,
		Desc:     msg,
	}
}

func Ok(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrNo
}
func RequestFail(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRequestFail
}
func InternalServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrInternalServer
}

func Params(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrParams
}

func ValidateParams(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrValidateParams
}
func Authentication(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAuthentication
}
func NotFound(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrNotFound
}
func AuthenticationHeader(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAuthenticationHeader
}
func AppKey(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAppKey
}
func Sign(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrSign
}
func Permission(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrPermission
}
func TooManyRequests(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrTooManyRequests
}
func InvalidJson(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrInvalidJson
}
func Timeout(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrTimeout
}
func AuthExpired(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrAuthExpired
}
func ElasticsearchServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrElasticsearchServer
}
func ElasticsearchDSL(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrElasticsearchDSL
}
func MysqlServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMysqlServer
}
func MysqlSQL(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMysqlSQL
}

func MongoServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMongoServer
}
func MongoDSL(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrMongoDSL
}

func RedisServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRedisServer
}
func KafkaServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrKafkaServer
}
func RabbitMQServer(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrRabbitMQServer
}

func Search(err ...ErrCode) ErrCode {
	if len(err) > 0 {
		return err[0]
	}
	return ErrNotFound
}
