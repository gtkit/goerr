package goerr

import (
	"fmt"
	"strconv"
)

/*
错误码规范

	编码格式

		项目组代号(10) + 服务代号(01) + 模块代号(00~99) + 错误序号(00~99)

	约定

		1. 模块 00 为通用错误；01~99 为业务或基础设施模块。
		2. 进入统一响应体后，HTTP 默认返回 200；客户端以业务错误码判断成功或失败。
		3. 仅当请求未进入统一响应体，或网关/框架必须表达协议级语义时，才返回非 200 HTTP 状态码。
		4. 错误码一经发布不得复用；废弃码保留占位，不允许回填新语义。
		5. 新增错误码优先复用已有语义，避免为同一失败场景创建多个近义码。
		6. 认证错误优先细分 token missing / invalid / expired / revoked。
		7. 超时错误优先细分 canceled / deadline exceeded；ErrTimeout 仅作兼容兜底。

	| 错误标识                 | 错误码    | 默认HTTP状态码 | 描述                           |
	| ----------------------- | -------- | -------------- | ------------------------------ |
	| ErrNo                   | 10010000 | 200            | OK                             |
	| ErrInternalServer       | 10010001 | 200            | Internal server error          |
	| ErrParams               | 10010002 | 200            | Illegal params                 |
	| ErrValidateParams       | 10010003 | 200            | Validation Errors              |
	| ErrInvalidJson          | 10010004 | 200            | Invalid json                   |
	| ErrAuthentication       | 10010005 | 200            | Authentication failed          |
	| ErrAuthenticationHeader | 10010006 | 200            | Authentication header Illegal  |
	| ErrVipRights            | 10010007 | 200            | Not Vip Rights                 |
	| ErrPermission           | 10010008 | 200            | Permission denied              |
	| ErrAppKey               | 10010009 | 200            | Invalid app key                |
	| ErrSign                 | 10010010 | 200            | Invalid sign                   |
	| ErrExpired              | 10010011 | 200            | Token expired                  |
	| ErrTimeout              | 10010012 | 200            | Server response timeout        |
	| ErrNotFound             | 10010013 | 200            | Page not found                 |
	| ErrTooManyRequests      | 10010014 | 200            | Too Many Requests              |
	| ErrRequestFail          | 10010015 | 200            | Request Fail                   |
	| ErrConflict             | 10010016 | 200            | Conflict                       |
	| ErrAlreadyExists        | 10010017 | 200            | Resource already exists        |
	| ErrServiceUnavailable   | 10010018 | 200            | Service unavailable            |
	| ErrDependencyUnavailable| 10010019 | 200            | Dependency unavailable         |
	| ErrCanceled             | 10010020 | 200            | Request canceled               |
	| ErrDeadlineExceeded     | 10010021 | 200            | Deadline exceeded              |
	| ErrTokenMissing         | 10010022 | 200            | Token missing                  |
	| ErrTokenInvalid         | 10010023 | 200            | Token invalid                  |
	| ErrTokenRevoked         | 10010024 | 200            | Token revoked                  |
	| ErrDuplicateRequest     | 10010025 | 200            | Duplicate request              |
	| ErrRecordNotFound       | 10010026 | 200            | Record not found               |
	|                         |          |                |                                |
	| ErrElasticsearchServer  | 10010101 | 200            | Elasticsearch server error     |
	| ErrElasticsearchDSL     | 10010102 | 200            | Elasticsearch DSL error        |
	| ErrMysqlServer          | 10010201 | 200            | Mysql server error             |
	| ErrMysqlSQL             | 10010202 | 200            | Illegal SQL                    |
	| ErrMysqlQuery           | 10010203 | 200            | Mysql query error              |
	| ErrMongoServer          | 10010301 | 200            | MongoDB server error           |
	| ErrMongoDSL             | 10010302 | 200            | MongoDB DSL error              |
	| ErrMongoQuery           | 10010303 | 200            | MongoDB query error            |
	| ErrRedisServer          | 10010401 | 200            | Redis server error             |
	| ErrRedisQuery           | 10010402 | 200            | Redis query error              |
	| ErrKafkaServer          | 10010501 | 200            | Kafka server error             |
	| ErrKafkaProducer        | 10010502 | 200            | Kafka Producer error           |
	| ErrKafkaConsumer        | 10010503 | 200            | Kafka Consumer error           |
	| ErrRabbitMQServer       | 10010601 | 200            | RabbitMQ server error          |
	| ErrRabbitMQProducer     | 10010602 | 200            | RabbitMQ Producer error        |
	| ErrRabbitMQConsumer     | 10010603 | 200            | RabbitMQ Consumer error        |
*/

// Code 错误码类型，承载 项目组+服务+模块+错误码 的完整标识。
type Code int32

const (
	minCode Code = 10000000
	maxCode Code = 10999999
)

// ValidateCode 校验错误码是否符合固定的编码规范。
func ValidateCode(code Code) error {
	if code < minCode || code > maxCode {
		return fmt.Errorf(
			"goerr: invalid code %d: must follow 项目组(10)+服务(01)+模块(00~99)+错误序号(00~99), valid range [%d,%d]",
			code, minCode, maxCode,
		)
	}
	return nil
}

// 通用错误码（模块代号 00）。
const (
	ErrNo                    Code = 10010000
	ErrInternalServer        Code = 10010001
	ErrParams                Code = 10010002
	ErrValidateParams        Code = 10010003
	ErrInvalidJson           Code = 10010004
	ErrAuthentication        Code = 10010005
	ErrAuthenticationHeader  Code = 10010006
	ErrVipRights             Code = 10010007
	ErrPermission            Code = 10010008
	ErrAppKey                Code = 10010009
	ErrSign                  Code = 10010010
	ErrExpired               Code = 10010011
	ErrTimeout               Code = 10010012
	ErrNotFound              Code = 10010013
	ErrTooManyRequests       Code = 10010014
	ErrRequestFail           Code = 10010015
	ErrConflict              Code = 10010016
	ErrAlreadyExists         Code = 10010017
	ErrServiceUnavailable    Code = 10010018
	ErrDependencyUnavailable Code = 10010019
	ErrCanceled              Code = 10010020
	ErrDeadlineExceeded      Code = 10010021
	ErrTokenMissing          Code = 10010022
	ErrTokenInvalid          Code = 10010023
	ErrTokenRevoked          Code = 10010024
	ErrDuplicateRequest      Code = 10010025
	ErrRecordNotFound        Code = 10010026
)

// Elasticsearch 模块（模块代号 01）。
const (
	ErrElasticsearchServer Code = 10010101
	ErrElasticsearchDSL    Code = 10010102
)

// MySQL 模块（模块代号 02）。
const (
	ErrMysqlServer Code = 10010201
	ErrMysqlSQL    Code = 10010202
	ErrMysqlQuery  Code = 10010203
)

// MongoDB 模块（模块代号 03）。
const (
	ErrMongoServer Code = 10010301
	ErrMongoDSL    Code = 10010302
	ErrMongoQuery  Code = 10010303
)

// Redis 模块（模块代号 04）。
const (
	ErrRedisServer Code = 10010401
	ErrRedisQuery  Code = 10010402
)

// Kafka 模块（模块代号 05）。
const (
	ErrKafkaServer   Code = 10010501
	ErrKafkaProducer Code = 10010502
	ErrKafkaConsumer Code = 10010503
)

// RabbitMQ 模块（模块代号 06）。
const (
	ErrRabbitMQServer   Code = 10010601
	ErrRabbitMQProducer Code = 10010602
	ErrRabbitMQConsumer Code = 10010603
)

// String 转换错误码为字符串。
func (c Code) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// Message 返回错误码对应的默认消息。
// 使用 switch 替代 map 查表：只读场景下 switch 零分配、无 hash 开销、编译器可优化为跳转表。
func (c Code) Message() string {
	switch c {
	case ErrNo:
		return "OK"
	case ErrInternalServer:
		return "Internal Server Error"
	case ErrParams:
		return "Illegal params"
	case ErrValidateParams:
		return "Validation Errors"
	case ErrInvalidJson:
		return "Invalid json"
	case ErrAuthentication:
		return "Authentication failed"
	case ErrAuthenticationHeader:
		return "Authentication header Illegal"
	case ErrVipRights:
		return "Not Vip Rights"
	case ErrPermission:
		return "Permission denied"
	case ErrAppKey:
		return "Invalid app key"
	case ErrSign:
		return "Invalid sign"
	case ErrExpired:
		return "Token expired"
	case ErrTimeout:
		return "Server response timeout"
	case ErrNotFound:
		return "Page not found"
	case ErrTooManyRequests:
		return "Too Many Requests"
	case ErrRequestFail:
		return "Request Fail"
	case ErrConflict:
		return "Conflict"
	case ErrAlreadyExists:
		return "Resource already exists"
	case ErrServiceUnavailable:
		return "Service unavailable"
	case ErrDependencyUnavailable:
		return "Dependency unavailable"
	case ErrCanceled:
		return "Request canceled"
	case ErrDeadlineExceeded:
		return "Deadline exceeded"
	case ErrTokenMissing:
		return "Token missing"
	case ErrTokenInvalid:
		return "Token invalid"
	case ErrTokenRevoked:
		return "Token revoked"
	case ErrDuplicateRequest:
		return "Duplicate request"
	case ErrRecordNotFound:
		return "Record not found"
	case ErrElasticsearchServer:
		return "Elasticsearch server error"
	case ErrElasticsearchDSL:
		return "Elasticsearch DSL error"
	case ErrMysqlServer:
		return "Mysql server error"
	case ErrMysqlSQL:
		return "Illegal SQL"
	case ErrMysqlQuery:
		return "Mysql query error"
	case ErrMongoServer:
		return "MongoDB server error"
	case ErrMongoDSL:
		return "MongoDB DSL error"
	case ErrMongoQuery:
		return "MongoDB query error"
	case ErrRedisServer:
		return "Redis server error"
	case ErrRedisQuery:
		return "Redis query error"
	case ErrKafkaServer:
		return "Kafka server error"
	case ErrKafkaProducer:
		return "Kafka Producer error"
	case ErrKafkaConsumer:
		return "Kafka Consumer error"
	case ErrRabbitMQServer:
		return "RabbitMQ server error"
	case ErrRabbitMQProducer:
		return "RabbitMQ Producer error"
	case ErrRabbitMQConsumer:
		return "RabbitMQ Consumer error"
	default:
		return "Unknown error"
	}
}
