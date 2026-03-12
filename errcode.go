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
	ErrNo                   int32 = 1001000
	ErrInternalServer       int32 = 10010001
	ErrParams               int32 = 10010002
	ErrValidateParams       int32 = 10010003
	ErrInvalidJson          int32 = 10010004
	ErrAuthentication       int32 = 10010005
	ErrAuthenticationHeader int32 = 10010006
	ErrVipRights            int32 = 10010007
	ErrPermission           int32 = 10010008
	ErrAppKey               int32 = 10010009
	ErrSign                 int32 = 10010010
	ErrExpired              int32 = 10010011
	ErrTimeout              int32 = 10010012
	ErrNotFound             int32 = 10010013
	ErrTooManyRequests      int32 = 10010014
	ErrRequestFail          int32 = 10010015

	ErrElasticsearchServer int32 = 10010101
	ErrElasticsearchDSL    int32 = 10010102

	ErrMysqlServer int32 = 10010201
	ErrMysqlSQL    int32 = 10010202
	ErrMysqlQuery  int32 = 10010203

	ErrMongoServer int32 = 10010301
	ErrMongoDSL    int32 = 10010302
	ErrMongoQuery  int32 = 10010303

	ErrRedisServer int32 = 10010401
	ErrRedisQuery  int32 = 10010402

	ErrKafkaServer   int32 = 10010501
	ErrKafkaProducer int32 = 10010502
	ErrKafkaConsumer int32 = 10010503

	ErrRabbitMQServer   int32 = 10010601
	ErrRabbitMQProducer int32 = 10010602
	ErrRabbitMQConsumer int32 = 10010603
)

var codeMessages = map[int32]string{
	ErrNo:                   "No Error",
	ErrInternalServer:       "Internal Server Error",
	ErrParams:               "Illegal params",
	ErrValidateParams:       "Validation Errors",
	ErrInvalidJson:          "Invalid json",
	ErrAuthentication:       "Authentication failed",
	ErrAuthenticationHeader: "Authentication header Illegal",
	ErrVipRights:            "Not Vip Rights",
	ErrPermission:           "Permission denied",
	ErrAppKey:               "Invalid app key",
	ErrSign:                 "Invalid sign",
	ErrExpired:              "Token expired",
	ErrTimeout:              "Server response timeout",
	ErrNotFound:             "Page not found",
	ErrTooManyRequests:      "Too Many Requests",
	ErrRequestFail:          "Request Fail",
	ErrElasticsearchServer:  "Elasticsearch server error",
	ErrElasticsearchDSL:     "Elasticsearch DSL error",
	ErrMysqlServer:          "Mysql server error",
	ErrMysqlSQL:             "Illegal SQL",
	ErrMysqlQuery:           "Mysql query error",
	ErrMongoServer:          "MongoDB server error",
	ErrMongoDSL:             "MongoDB DSL error",
	ErrMongoQuery:           "MongoDB query error",
	ErrRedisServer:          "Redis server error",
	ErrRedisQuery:           "Redis query error",
	ErrKafkaServer:          "Kafka server error",
	ErrKafkaProducer:        "Kafka Producer error",
	ErrKafkaConsumer:        "Kafka Consumer error",
	ErrRabbitMQServer:       "RabbitMQ server error",
	ErrRabbitMQProducer:     "RabbitMQ Producer error",
	ErrRabbitMQConsumer:     "RabbitMQ Consumer error",
}
