package goerr

import (
	"fmt"
	"net/http"
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

	错误码 → HTTP 状态码 → 默认文案 的完整清单见本文件的 codeTable，它是唯一真源：
	Code.Message() 与 status.go 的预构建 Status 都从那里取值。
	新增错误码需要改两处：下方常量分组、codeTable；若要提供 StatusXxx()，再在 status.go 登记。
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
	ErrNotLoggedIn           Code = 10010033
	ErrLoginExpired          Code = 10010034
	ErrUserNotRegistered     Code = 10010035
	ErrUserDeactivated       Code = 10010036
	ErrPasswordIncorrect     Code = 10010037
	ErrCaptchaInvalid        Code = 10010038
	ErrCaptchaExpired        Code = 10010039
	ErrAccountLocked         Code = 10010040
	ErrUserAlreadyRegistered Code = 10010041
	ErrSessionKickedOut      Code = 10010042
	ErrLoginNotPerformed     Code = 10010043
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

// 订单模块（模块代号 07）。
const (
	ErrOrderNotFound      Code = 10010701
	ErrOrderStatusInvalid Code = 10010702
	ErrOrderClosed        Code = 10010703
	ErrOrderCanceled      Code = 10010704
	ErrOrderPaid          Code = 10010705
	ErrOrderCompleted     Code = 10010706
	ErrOrderExpired       Code = 10010707
	ErrOrderRefunded      Code = 10010708
	ErrOrderPending       Code = 10010709
	ErrOrderPaying        Code = 10010710
)

// codeMeta 是错误码的附属信息：语义化 HTTP 状态码与默认文案。
type codeMeta struct {
	httpCode int
	msg      string
}

// codeTable 是错误码元信息的唯一真源。Code.Message() 与 status.go 的预构建
// Status 都从这里取值，避免「码加了但文案或 HTTP 码漏配」这类只能靠人眼发现的偏差。
// 未登记的码在 newStatus 时 panic（包初始化期暴露），Message() 返回「未知错误」。
var codeTable = map[Code]codeMeta{
	// 通用（模块 00）
	ErrNo:                    {http.StatusOK, "成功"},
	ErrVipRights:             {http.StatusOK, "非 VIP 用户"},
	ErrDuplicateRequest:      {http.StatusOK, "重复请求"},
	ErrInsufficientBalance:   {http.StatusOK, "余额不足"},
	ErrInsufficientStock:     {http.StatusOK, "库存不足"},
	ErrBusinessRuleViolation: {http.StatusOK, "业务规则限制"},
	ErrAccountDisabled:       {http.StatusOK, "账号已禁用"},
	ErrOperationNotAllowed:   {http.StatusOK, "当前状态不允许此操作"},
	ErrQuotaExceeded:         {http.StatusOK, "配额已用尽"},
	ErrUserNotRegistered:     {http.StatusOK, "用户未注册"},
	ErrUserDeactivated:       {http.StatusOK, "用户已注销"},
	ErrPasswordIncorrect:     {http.StatusOK, "密码错误"},
	ErrCaptchaInvalid:        {http.StatusOK, "验证码错误"},
	ErrCaptchaExpired:        {http.StatusOK, "验证码已过期"},
	ErrAccountLocked:         {http.StatusOK, "账号已锁定"},
	ErrLoginNotPerformed:     {http.StatusOK, "未执行登录操作"},

	ErrRequestFail: {http.StatusBadRequest, "请求失败"},
	ErrParams:      {http.StatusBadRequest, "参数不合法"},
	ErrInvalidJson: {http.StatusBadRequest, "JSON 格式错误"},

	ErrValidateParams: {http.StatusUnprocessableEntity, "参数校验失败"},

	ErrAuthentication:       {http.StatusUnauthorized, "认证失败"},
	ErrAuthenticationHeader: {http.StatusUnauthorized, "认证请求头不合法"},
	ErrAppKey:               {http.StatusUnauthorized, "应用密钥无效"},
	ErrSign:                 {http.StatusUnauthorized, "签名无效"},
	ErrExpired:              {http.StatusUnauthorized, "令牌已过期"},
	ErrTokenMissing:         {http.StatusUnauthorized, "缺少令牌"},
	ErrTokenInvalid:         {http.StatusUnauthorized, "令牌无效"},
	ErrTokenRevoked:         {http.StatusUnauthorized, "令牌已吊销"},
	ErrNotLoggedIn:          {http.StatusUnauthorized, "未登录"},
	ErrLoginExpired:         {http.StatusUnauthorized, "登录已过期"},
	ErrSessionKickedOut:     {http.StatusUnauthorized, "会话已被踢下线"},

	ErrPermission: {http.StatusForbidden, "权限不足"},

	ErrNotFound:       {http.StatusNotFound, "资源不存在"},
	ErrRecordNotFound: {http.StatusNotFound, "记录不存在"},

	ErrConflict:              {http.StatusConflict, "资源冲突"},
	ErrAlreadyExists:         {http.StatusConflict, "资源已存在"},
	ErrUserAlreadyRegistered: {http.StatusConflict, "用户已注册"},

	ErrTooManyRequests: {http.StatusTooManyRequests, "请求过于频繁"},

	ErrInternalServer:   {http.StatusInternalServerError, "服务器内部错误"},
	ErrTimeout:          {http.StatusInternalServerError, "服务器响应超时"},
	ErrCanceled:         {http.StatusInternalServerError, "请求已取消"},
	ErrDeadlineExceeded: {http.StatusInternalServerError, "请求超时"},

	ErrServiceUnavailable:    {http.StatusServiceUnavailable, "服务不可用"},
	ErrDependencyUnavailable: {http.StatusServiceUnavailable, "依赖服务不可用"},

	// Elasticsearch（模块 01）
	ErrElasticsearchServer: {http.StatusInternalServerError, "Elasticsearch 服务异常"},
	ErrElasticsearchDSL:    {http.StatusInternalServerError, "Elasticsearch 查询语句错误"},

	// MySQL（模块 02）
	ErrMysqlServer: {http.StatusInternalServerError, "MySQL 服务异常"},
	ErrMysqlSQL:    {http.StatusInternalServerError, "SQL 语句错误"},
	ErrMysqlQuery:  {http.StatusInternalServerError, "MySQL 查询失败"},

	// MongoDB（模块 03）
	ErrMongoServer: {http.StatusInternalServerError, "MongoDB 服务异常"},
	ErrMongoDSL:    {http.StatusInternalServerError, "MongoDB 查询语句错误"},
	ErrMongoQuery:  {http.StatusInternalServerError, "MongoDB 查询失败"},

	// Redis（模块 04）
	ErrRedisServer: {http.StatusInternalServerError, "Redis 服务异常"},
	ErrRedisQuery:  {http.StatusInternalServerError, "Redis 查询失败"},

	// Kafka（模块 05）
	ErrKafkaServer:   {http.StatusInternalServerError, "Kafka 服务异常"},
	ErrKafkaProducer: {http.StatusInternalServerError, "Kafka 生产者异常"},
	ErrKafkaConsumer: {http.StatusInternalServerError, "Kafka 消费者异常"},

	// RabbitMQ（模块 06）
	ErrRabbitMQServer:   {http.StatusInternalServerError, "RabbitMQ 服务异常"},
	ErrRabbitMQProducer: {http.StatusInternalServerError, "RabbitMQ 生产者异常"},
	ErrRabbitMQConsumer: {http.StatusInternalServerError, "RabbitMQ 消费者异常"},

	// 订单（模块 07）
	ErrOrderNotFound:      {http.StatusOK, "订单不存在"},
	ErrOrderStatusInvalid: {http.StatusOK, "订单状态异常"},
	ErrOrderClosed:        {http.StatusOK, "订单已关闭"},
	ErrOrderCanceled:      {http.StatusOK, "订单已取消"},
	ErrOrderPaid:          {http.StatusOK, "订单已支付"},
	ErrOrderCompleted:     {http.StatusOK, "订单已完成"},
	ErrOrderExpired:       {http.StatusOK, "订单已过期"},
	ErrOrderRefunded:      {http.StatusOK, "订单已退款"},
	ErrOrderPending:       {http.StatusOK, "订单待支付"},
	ErrOrderPaying:        {http.StatusOK, "订单支付中"},
}

// String 转换错误码为字符串。
func (c Code) String() string {
	return strconv.FormatInt(int64(c), 10)
}

// Message 返回错误码对应的默认消息；未在 codeTable 登记的码返回「未知错误」。
func (c Code) Message() string {
	if m, ok := codeTable[c]; ok {
		return m.msg
	}
	return "未知错误"
}
