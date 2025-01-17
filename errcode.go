package goerr

/*
ErrStatus
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
	| ErrValidateParams       | 10010003 | 400        | validator.ValidationErrors    | 验证器。验证错误
	| ErrInvalidJson          | 10010004 | 400        | Invalid json                  | 非法json
	| ErrAuthentication       | 10010005 | 401        | Authentication failed         | 身份验证失败
	| ErrAuthenticationHeader | 10010006 | 401        | Authentication header Illegal | 认证头非法
	| ErrVipRights       	  | 10010007 | 401        | Not Vip Rights		          | 非会员
	| ErrPermission           | 10010008 | 403        | Permission denied             | 权限不足,禁止访问
	| ErrAppKey               | 10010009 | 401        | Invalid app key               | 无效应用程序密钥
	| ErrSign                 | 10010010 | 401        | Invalid sign                  | 无效签名
	| ErrExpired              | 10010011 | 504        | Token expired                 | 令牌过期
	| ErrTimeout              | 10010012 | 504        | Server response timeout       | 服务器响应超时
	| ErrNotFound             | 10010013 | 404        | Page not found                | 找不到
	| ErrTooManyRequests      | 10010014 | 429        | Too Many Requests             | 请求过多
	| ErrRequestFail	      | 10010015 | 200        | ErrRequestFail	              | 请求失败或者响应为空

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

const (
	ErrNo                   = 10010000
	ErrInternalServer       = 10010001
	ErrParams               = 10010002
	ErrValidateParams       = 10010003
	ErrInvalidJson          = 10010004
	ErrAuthentication       = 10010005
	ErrAuthenticationHeader = 10010006
	ErrVipRights            = 10010007
	ErrPermission           = 10010008
	ErrAppKey               = 10010009
	ErrSign                 = 10010010
	ErrExpired              = 10010011
	ErrTimeout              = 10010012
	ErrNotFound             = 10010013
	ErrTooManyRequests      = 10010014
	ErrRequestFail          = 10010015

	ErrElasticsearchServer = 10010101
	ErrElasticsearchDSL    = 10010102

	ErrMysqlServer = 10010201
	ErrMysqlSQL    = 10010202

	ErrMongoServer = 10010301
	ErrMongoDSL    = 10010302

	ErrRedisServer = 10010401

	ErrKafkaServer   = 10010501
	ErrKafkaProducer = 10010502
	ErrKafkaConsumer = 10010503

	ErrRabbitMQServer   = 10010601
	ErrRabbitMQProducer = 10010602
	ErrRabbitMQConsumer = 10010603
)

func statusText(code int) string {
	switch code {
	case ErrNo:
		return "OK"
	case ErrRequestFail:
		return "请求失败"
	case ErrInternalServer:
		return "服务器内部错误"
	case ErrPermission:
		return "Permission denied"
	case ErrParams:
		return "Illegal params"
	case ErrValidateParams:
		return "Parameter validation failure"
	case ErrAuthentication:
		return "Authentication failed"
	case ErrVipRights:
		return "Not Vip Rights"
	case ErrNotFound:
		return "Not Found"
	case ErrAuthenticationHeader:
		return "Authentication header Illegal"
	case ErrAppKey:
		return "Invalid app key"
	case ErrSign:
		return "Invalid signature"
	case ErrTooManyRequests:
		return "Too Many Requests"
	case ErrInvalidJson:
		return "Invalid Json"
	case ErrTimeout:
		return "Server response timeout"
	case ErrExpired:
		return "Authentication expired"
	case ErrElasticsearchServer:
		return "Elasticsearch server error"
	case ErrElasticsearchDSL:
		return "Elasticsearch DSL error"
	case ErrMysqlServer:
		return "Mysql server error"
	case ErrMysqlSQL:
		return "Illegal SQL"
	case ErrMongoServer:
		return "MongoDB server error"
	case ErrMongoDSL:
		return "MongoDB DSL error"
	case ErrRedisServer:
		return "Redis server error"
	case ErrKafkaServer:
		return "Kafka server error"
	case ErrKafkaProducer:
		return "Kafka producer error"
	case ErrKafkaConsumer:
		return "Kafka consumer error"
	case ErrRabbitMQServer:
		return "RabbitMQ server error"
	case ErrRabbitMQProducer:
		return "RabbitMQ producer error"
	case ErrRabbitMQConsumer:
		return "RabbitMQ consumer error"
	default:
		return "Unknown Error"
	}
}
