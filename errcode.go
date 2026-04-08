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
		2. HTTP 状态码语义化：业务逻辑错误返回 200，参数错误 400/422，认证 401，权限 403，
		   资源不存在 404，冲突 409，限流 429，服务端/中间件错误 500，服务不可用 503。
		3. 客户端通过 HTTP 状态码判断错误大类，通过响应体中的业务 code 判断具体错误。
		4. 错误码一经发布不得复用；废弃码保留占位，不允许回填新语义。
		5. 新增错误码优先复用已有语义，避免为同一失败场景创建多个近义码。
		6. 认证错误优先细分 token missing / invalid / expired / revoked。
		7. 超时错误优先细分 canceled / deadline exceeded；ErrTimeout 仅作兼容兜底。

	| 错误标识                 | 错误码    | 默认HTTP状态码 | 描述                           |
	| ----------------------- | -------- | -------------- | ------------------------------ |
	| ErrNo                   | 10010000 | 200            | 成功                           |
	| ErrInternalServer       | 10010001 | 500            | 服务器内部错误                 |
	| ErrParams               | 10010002 | 400            | 参数不合法                     |
	| ErrValidateParams       | 10010003 | 422            | 参数校验失败                   |
	| ErrInvalidJson          | 10010004 | 400            | JSON 格式错误                  |
	| ErrAuthentication       | 10010005 | 401            | 认证失败                       |
	| ErrAuthenticationHeader | 10010006 | 401            | 认证请求头不合法               |
	| ErrVipRights            | 10010007 | 200            | 非 VIP 用户                    |
	| ErrPermission           | 10010008 | 403            | 权限不足                       |
	| ErrAppKey               | 10010009 | 401            | 应用密钥无效                   |
	| ErrSign                 | 10010010 | 401            | 签名无效                       |
	| ErrExpired              | 10010011 | 401            | 令牌已过期                     |
	| ErrTimeout              | 10010012 | 500            | 服务器响应超时                 |
	| ErrNotFound             | 10010013 | 404            | 资源不存在                     |
	| ErrTooManyRequests      | 10010014 | 429            | 请求过于频繁                   |
	| ErrRequestFail          | 10010015 | 400            | 请求失败                       |
	| ErrConflict             | 10010016 | 409            | 资源冲突                       |
	| ErrAlreadyExists        | 10010017 | 409            | 资源已存在                     |
	| ErrServiceUnavailable   | 10010018 | 503            | 服务不可用                     |
	| ErrDependencyUnavailable| 10010019 | 503            | 依赖服务不可用                 |
	| ErrCanceled             | 10010020 | 500            | 请求已取消                     |
	| ErrDeadlineExceeded     | 10010021 | 500            | 请求超时                       |
	| ErrTokenMissing         | 10010022 | 401            | 缺少令牌                       |
	| ErrTokenInvalid         | 10010023 | 401            | 令牌无效                       |
	| ErrTokenRevoked         | 10010024 | 401            | 令牌已吊销                     |
	| ErrDuplicateRequest     | 10010025 | 200            | 重复请求                       |
	| ErrRecordNotFound       | 10010026 | 404            | 记录不存在                     |
	| ErrInsufficientBalance  | 10010027 | 200            | 余额不足                       |
	| ErrInsufficientStock    | 10010028 | 200            | 库存不足                       |
	| ErrBusinessRuleViolation| 10010029 | 200            | 业务规则限制                   |
	| ErrAccountDisabled      | 10010030 | 200            | 账号已禁用                     |
	| ErrOperationNotAllowed  | 10010031 | 200            | 当前状态不允许此操作           |
	| ErrQuotaExceeded        | 10010032 | 200            | 配额已用尽                     |
	|                         |          |                |                                |
	| ErrElasticsearchServer  | 10010101 | 500            | Elasticsearch 服务异常         |
	| ErrElasticsearchDSL     | 10010102 | 500            | Elasticsearch 查询语句错误     |
	| ErrMysqlServer          | 10010201 | 500            | MySQL 服务异常                 |
	| ErrMysqlSQL             | 10010202 | 500            | SQL 语句错误                   |
	| ErrMysqlQuery           | 10010203 | 500            | MySQL 查询失败                 |
	| ErrMongoServer          | 10010301 | 500            | MongoDB 服务异常               |
	| ErrMongoDSL             | 10010302 | 500            | MongoDB 查询语句错误           |
	| ErrMongoQuery           | 10010303 | 500            | MongoDB 查询失败               |
	| ErrRedisServer          | 10010401 | 500            | Redis 服务异常                 |
	| ErrRedisQuery           | 10010402 | 500            | Redis 查询失败                 |
	| ErrKafkaServer          | 10010501 | 500            | Kafka 服务异常                 |
	| ErrKafkaProducer        | 10010502 | 500            | Kafka 生产者异常               |
	| ErrKafkaConsumer        | 10010503 | 500            | Kafka 消费者异常               |
	| ErrRabbitMQServer       | 10010601 | 500            | RabbitMQ 服务异常              |
	| ErrRabbitMQProducer     | 10010602 | 500            | RabbitMQ 生产者异常            |
	| ErrRabbitMQConsumer     | 10010603 | 500            | RabbitMQ 消费者异常            |
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
	ErrInsufficientBalance   Code = 10010027
	ErrInsufficientStock     Code = 10010028
	ErrBusinessRuleViolation Code = 10010029
	ErrAccountDisabled       Code = 10010030
	ErrOperationNotAllowed   Code = 10010031
	ErrQuotaExceeded         Code = 10010032
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
		return "成功"
	case ErrInternalServer:
		return "服务器内部错误"
	case ErrParams:
		return "参数不合法"
	case ErrValidateParams:
		return "参数校验失败"
	case ErrInvalidJson:
		return "JSON 格式错误"
	case ErrAuthentication:
		return "认证失败"
	case ErrAuthenticationHeader:
		return "认证请求头不合法"
	case ErrVipRights:
		return "非 VIP 用户"
	case ErrPermission:
		return "权限不足"
	case ErrAppKey:
		return "应用密钥无效"
	case ErrSign:
		return "签名无效"
	case ErrExpired:
		return "令牌已过期"
	case ErrTimeout:
		return "服务器响应超时"
	case ErrNotFound:
		return "资源不存在"
	case ErrTooManyRequests:
		return "请求过于频繁"
	case ErrRequestFail:
		return "请求失败"
	case ErrConflict:
		return "资源冲突"
	case ErrAlreadyExists:
		return "资源已存在"
	case ErrServiceUnavailable:
		return "服务不可用"
	case ErrDependencyUnavailable:
		return "依赖服务不可用"
	case ErrCanceled:
		return "请求已取消"
	case ErrDeadlineExceeded:
		return "请求超时"
	case ErrTokenMissing:
		return "缺少令牌"
	case ErrTokenInvalid:
		return "令牌无效"
	case ErrTokenRevoked:
		return "令牌已吊销"
	case ErrDuplicateRequest:
		return "重复请求"
	case ErrRecordNotFound:
		return "记录不存在"
	case ErrInsufficientBalance:
		return "余额不足"
	case ErrInsufficientStock:
		return "库存不足"
	case ErrBusinessRuleViolation:
		return "业务规则限制"
	case ErrAccountDisabled:
		return "账号已禁用"
	case ErrOperationNotAllowed:
		return "当前状态不允许此操作"
	case ErrQuotaExceeded:
		return "配额已用尽"
	case ErrElasticsearchServer:
		return "Elasticsearch 服务异常"
	case ErrElasticsearchDSL:
		return "Elasticsearch 查询语句错误"
	case ErrMysqlServer:
		return "MySQL 服务异常"
	case ErrMysqlSQL:
		return "SQL 语句错误"
	case ErrMysqlQuery:
		return "MySQL 查询失败"
	case ErrMongoServer:
		return "MongoDB 服务异常"
	case ErrMongoDSL:
		return "MongoDB 查询语句错误"
	case ErrMongoQuery:
		return "MongoDB 查询失败"
	case ErrRedisServer:
		return "Redis 服务异常"
	case ErrRedisQuery:
		return "Redis 查询失败"
	case ErrKafkaServer:
		return "Kafka 服务异常"
	case ErrKafkaProducer:
		return "Kafka 生产者异常"
	case ErrKafkaConsumer:
		return "Kafka 消费者异常"
	case ErrRabbitMQServer:
		return "RabbitMQ 服务异常"
	case ErrRabbitMQProducer:
		return "RabbitMQ 生产者异常"
	case ErrRabbitMQConsumer:
		return "RabbitMQ 消费者异常"
	default:
		return "未知错误"
	}
}
