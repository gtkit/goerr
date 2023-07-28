package goerr

/*
ErrCode
*

	项目组代号:10
	服务代号:01
	模块代号:0~99
	错误码：0~99

	| 错误标识	              | 错误码	 | HTTP状态码  | 描述							  |
	| ----------------------- | -------- | ---------- | ----------------------------- |
	| ErrNo                   | 10010000 | 200        | OK                            | 请求成功
	| ErrInternalServer       | 10010001 | 500        | Internal server error         | 服务器内部错误
	| ErrParams               | 10010002 | 400        | Illegal params                | 非法参数
	| ErrValidateParams       | 10010012 | 400        | validator.ValidationErrors    | 验证器。验证错误
	| ErrAuthentication       | 10010003 | 401        | Authentication failed         | 身份验证失败
	| ErrNotFound             | 10010004 | 404        | Page not found                | 找不到
	| ErrAuthenticationHeader | 10010005 | 401        | Authentication header Illegal | 认证头非法
	| ErrAppKey               | 10010006 | 401        | Invalid app key               | 无效应用程序密钥
	| ErrSecretKey            | 10010007 | 401        | Invalid secret key            | 无效应用程序密钥
	| ErrPermission           | 10010008 | 403        | Permission denied             | 没有权限
	| ErrTooManyRequests      | 10010013 | 429        | Too Many Requests             | 请求过多
	| ErrRequestFail	      | 10010014 | 200        | ErrRequestFail	              | 请求失败或者响应为空
	| ErrInvalidJson          | 10010009 | 500        | Invalid Json                  | 无效的JSON
	| ErrTimeout              | 10010010 | 504        | Server response timeout       | 服务器响应超时
	| ErrElasticsearchServer  | 10010101 | 500        | Elasticsearch server error    | Elasticsearch 服务器错误
	| ErrElasticsearchDSL     | 10010102 | 500        | Elasticsearch  DSL error      | Elasticsearch DSL 错误
	| ErrMysqlServer          | 10010201 | 500        | Mysql server error            | Mysql 服务器错误
	| ErrMysqlSQL             | 10010202 | 500        | Illegal SQL                   | 非法SQL
	| ErrMongoServer          | 10010301 | 500        | MongoDB server error          | MongoDB 服务器错误
	| ErrMongoDSL             | 10010302 | 500        | MongoDB DSL error             | MongoDB DSL 错误
	| ErrRedisServer          | 10010401 | 500        | Redis server error            | Redis 服务器错误
	| ErrRedisClient          | 10010402 | 500        | Redis client error            | Redis 客户端错误
	| ErrKafkaServer          | 10010501 | 500        | Kafka server error            | Kafka 服务器错误
	| ErrKafkaProducer        | 10010502 | 500        | Kafka Producer error          | Kafka Producer 错误
	| ErrKafkaConsumer        | 10010503 | 500        | Kafka Consumer error          | Kafka Consumer 错误
	| ErrRabbitMQServer       | 10010601 | 500        | RabbitMQ server error         | RabbitMq 服务器错误
	| ErrRabbitMQProducer     | 10010602 | 500        | RabbitMQ Producer error       | RabbitMq Producer 错误
	| ErrRabbitMQConsumer     | 10010603 | 500        | RabbitMQ Consumer error       | RabbitMq Consumer 错误
*/
type ErrCode struct {
	Code     int
	HTTPCode int
	Desc     string
}

var (
	// ErrNo 请求成功
	ErrNo = ErrCode{
		Code:     10010000,
		HTTPCode: 200,
		Desc:     "OK",
	}

	// ErrInternalServer 服务器内部错误
	ErrInternalServer = ErrCode{
		Code:     10010001,
		HTTPCode: 500,
		Desc:     "Internal server error",
	}

	// ErrParams 非法参数
	ErrParams = ErrCode{
		Code:     10010002,
		HTTPCode: 422,
		Desc:     "Illegal params",
	}

	// ErrAuthentication 身份验证失败
	ErrAuthentication = ErrCode{
		Code:     10010003,
		HTTPCode: 401,
		Desc:     "Authentication failed",
	}

	// ErrNotFound 找不到
	ErrNotFound = ErrCode{
		Code:     10010004,
		HTTPCode: 404,
		Desc:     "Page not found",
	}

	// ErrAuthenticationHeader 认证头非法
	ErrAuthenticationHeader = ErrCode{
		Code:     10010005,
		HTTPCode: 403,
		Desc:     "Authentication header Illegal",
	}

	// ErrAppKey 无效应用程序密钥
	ErrAppKey = ErrCode{
		Code:     10010006,
		HTTPCode: 403,
		Desc:     "Invalid app key",
	}

	// ErrSign 无效应用签名
	ErrSign = ErrCode{
		Code:     10010007,
		HTTPCode: 401,
		Desc:     "Invalid secret key",
	}

	// ErrPermission 没有权限
	ErrPermission = ErrCode{
		Code:     10010008,
		HTTPCode: 403,
		Desc:     "Permission denied",
	}

	// ErrInvalidJson 无效的JSON
	ErrInvalidJson = ErrCode{
		Code:     10010009,
		HTTPCode: 500,
		Desc:     "Invalid Json",
	}

	// ErrTimeout 服务器响应超时
	ErrTimeout = ErrCode{
		Code:     10010010,
		HTTPCode: 408,
		Desc:     "Server response timeout",
	}

	// ErrAuthExpired Authentication expired
	ErrAuthExpired = ErrCode{
		Code:     10010011,
		HTTPCode: 504,
		Desc:     "Authentication expired",
	}

	// ErrValidateParams 验证器。验证错误
	ErrValidateParams = ErrCode{
		Code:     10010012,
		HTTPCode: 422,
		Desc:     "validator.ValidationErrors",
	}

	// ErrTooManyRequests 请求过多
	ErrTooManyRequests = ErrCode{
		Code:     10010013,
		HTTPCode: 429,
		Desc:     "Too Many Requests",
	}

	// ErrRequestFail 请求失败或者响应为空
	ErrRequestFail = ErrCode{
		Code:     10010014,
		HTTPCode: 200,
		Desc:     "请求失败",
	}

	// ErrElasticsearchServer Elasticsearch 服务器错误
	ErrElasticsearchServer = ErrCode{
		Code:     10010101,
		HTTPCode: 500,
		Desc:     "Elasticsearch server error",
	}

	// ErrElasticsearchDSL Elasticsearch DSL 错误
	ErrElasticsearchDSL = ErrCode{
		Code:     10010102,
		HTTPCode: 500,
		Desc:     "Elasticsearch  DSL error",
	}

	// ErrMysqlServer Mysql 服务器错误
	ErrMysqlServer = ErrCode{
		Code:     10010201,
		HTTPCode: 500,
		Desc:     "Mysql server error",
	}

	// ErrMysqlSQL 非法SQL
	ErrMysqlSQL = ErrCode{
		Code:     10010202,
		HTTPCode: 500,
		Desc:     "Illegal SQL",
	}

	// ErrMongoServer MongoDB 服务器错误
	ErrMongoServer = ErrCode{
		Code:     10010301,
		HTTPCode: 500,
		Desc:     "MongoDB server error",
	}

	// ErrMongoDSL MongoDB DSL 错误
	ErrMongoDSL = ErrCode{
		Code:     10010302,
		HTTPCode: 500,
		Desc:     "MongoDB DSL error",
	}

	// ErrRedisServer Redis 服务器错误
	ErrRedisServer = ErrCode{
		Code:     10010401,
		HTTPCode: 500,
		Desc:     "Redis server error",
	}
	// ErrRedisClient Redis 客户端错误
	ErrRedisClient = ErrCode{
		Code:     10010402,
		HTTPCode: 500,
		Desc:     "Redis key error",
	}

	// ErrKafkaServer Kafka 服务器错误
	ErrKafkaServer = ErrCode{
		Code:     10010501,
		HTTPCode: 500,
		Desc:     "Kafka server error",
	}
	ErrKafkaProducer = ErrCode{
		Code:     10010502,
		HTTPCode: 500,
		Desc:     "Kafka producer error",
	}
	ErrKafkaConsumer = ErrCode{
		Code:     10010503,
		HTTPCode: 500,
		Desc:     "Kafka consumer error",
	}

	// ErrRabbitMQServer RabbitMq 服务器错误
	ErrRabbitMQServer = ErrCode{
		Code:     10010601,
		HTTPCode: 500,
		Desc:     "RabbitMq server error",
	}
	ErrRabbitMQProducer = ErrCode{
		Code:     10010602,
		HTTPCode: 500,
		Desc:     "RabbitMq producer error",
	}
	ErrRabbitMQConsumer = ErrCode{
		Code:     10010603,
		HTTPCode: 500,
		Desc:     "RabbitMq consumer error",
	}
)
