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

		1. 模块 00 为通用错误；01~99 为业务或基础设施模块，当前已分配：

			00 通用    01 Elasticsearch  02 MySQL   03 MongoDB  04 Redis
			05 Kafka   06 RabbitMQ       07 订单    08 支付     09 账号与第三方登录
			10 AI 大模型                 11 文件    12 邮件

		2. HTTP 状态码语义化：业务逻辑错误返回 200，参数错误 400/422，认证 401，权限 403，
		   资源不存在 404，冲突 409，限流 429，服务端/中间件错误 500，服务不可用 503。
		3. 客户端通过 HTTP 状态码判断错误大类，通过响应体中的业务 code 判断具体错误。
		4. 错误码一经发布不得复用；废弃码保留占位，不允许回填新语义。
		5. 新增错误码优先复用已有语义，避免为同一失败场景创建多个近义码。
		6. 认证错误优先细分 token missing / invalid / expired / revoked。
		7. 超时错误优先细分 canceled / deadline exceeded；ErrTimeout 仅作兼容兜底。
		8. 基础设施模块（01~06）的码描述「哪个组件的哪种故障」，业务含义由调用方在
		   Message() 里补充；业务模块（07~12）的码直接承载业务语义。
		9. 涉及资金与鉴权的失败（金额不符、验签失败）单独建码，不并入通用参数错误——
		   它们是安全事件，需要独立的日志级别与告警。

	错误码 → HTTP 状态码 → 默认文案 的完整清单见本文件的 codeTable，它是唯一真源：
	Code.Message() 与 status.go 的预构建 Status 都从那里取值。
	新增错误码需要改三处：下方常量分组、codeTable、status.go 的 var 与 StatusXxx()。
	前两处漏改会在包初始化期 panic，不会留到运行时。
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
	ErrRequestBodyTooLarge   Code = 10010044
	ErrUnsupportedMediaType  Code = 10010045
	ErrVersionConflict       Code = 10010046
	ErrOperationInProgress   Code = 10010047
	ErrRiskControlRejected   Code = 10010048
	ErrFeatureDisabled       Code = 10010049
	ErrSystemMaintenance     Code = 10010050
)

// Elasticsearch 模块（模块代号 01）。
const (
	ErrElasticsearchServer        Code = 10010101
	ErrElasticsearchDSL           Code = 10010102
	ErrElasticsearchIndexNotFound Code = 10010103
	ErrElasticsearchBulkPartial   Code = 10010104
)

// MySQL 模块（模块代号 02）。
const (
	ErrMysqlServer      Code = 10010201
	ErrMysqlSQL         Code = 10010202
	ErrMysqlQuery       Code = 10010203
	ErrMysqlDeadlock    Code = 10010204
	ErrMysqlTransaction Code = 10010205
)

// MongoDB 模块（模块代号 03）。
const (
	ErrMongoServer Code = 10010301
	ErrMongoDSL    Code = 10010302
	ErrMongoQuery  Code = 10010303
)

// Redis 模块（模块代号 04）。
const (
	ErrRedisServer            Code = 10010401
	ErrRedisQuery             Code = 10010402
	ErrRedisLockAcquireFailed Code = 10010403
	ErrRedisLockLost          Code = 10010404
)

// Kafka 模块（模块代号 05）。
const (
	ErrKafkaServer         Code = 10010501
	ErrKafkaProducer       Code = 10010502
	ErrKafkaConsumer       Code = 10010503
	ErrKafkaMessageDecode  Code = 10010504
	ErrKafkaRetryExhausted Code = 10010505
)

// RabbitMQ 模块（模块代号 06）。
const (
	ErrRabbitMQServer         Code = 10010601
	ErrRabbitMQProducer       Code = 10010602
	ErrRabbitMQConsumer       Code = 10010603
	ErrRabbitMQMessageDecode  Code = 10010604
	ErrRabbitMQRetryExhausted Code = 10010605
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

// 支付模块（模块代号 08）。覆盖支付单、回调验签与退款。
const (
	ErrPaymentNotFound           Code = 10010801
	ErrPaymentAlreadyPaid        Code = 10010802
	ErrPaymentClosed             Code = 10010803
	ErrPaymentExpired            Code = 10010804
	ErrPaymentMethodNotSupported Code = 10010805
	ErrPayerIdentityMissing      Code = 10010806
	ErrPaymentAmountMismatch     Code = 10010807
	ErrPaymentSignInvalid        Code = 10010808
	ErrPaymentCallbackInvalid    Code = 10010809
	ErrPaymentCreateFailed       Code = 10010810
	ErrPaymentQueryFailed        Code = 10010811
	ErrPaymentChannelUnavailable Code = 10010812
	ErrRefundNotFound            Code = 10010813
	ErrRefundAmountExceeded      Code = 10010814
	ErrRefundInProgress          Code = 10010815
	ErrRefundFailed              Code = 10010816
)

// 账号与第三方登录模块（模块代号 09）。通用登录态错误见模块 00。
const (
	ErrPhoneAlreadyBound      Code = 10010901
	ErrThirdPartyAlreadyBound Code = 10010902
	ErrThirdPartyNotBound     Code = 10010903
	ErrThirdPartyAuthFailed   Code = 10010904
	ErrThirdPartyUnavailable  Code = 10010905
	ErrSmsSendTooFrequent     Code = 10010906
	ErrSmsDailyLimitExceeded  Code = 10010907
	ErrSmsSendFailed          Code = 10010908
)

// AI 大模型模块（模块代号 10）。
// HTTP 码按「谁能解决」划分：用户可自行修正的给 200，我方配置/实现问题给 500，供应商侧故障给 503。
const (
	ErrLLMAuth              Code = 10011001
	ErrLLMRateLimited       Code = 10011002
	ErrLLMTimeout           Code = 10011003
	ErrLLMContextTooLong    Code = 10011004
	ErrLLMContentFiltered   Code = 10011005
	ErrLLMInvalidRequest    Code = 10011006
	ErrLLMProviderError     Code = 10011007
	ErrLLMModelNotAvailable Code = 10011008
	ErrLLMStructuredOutput  Code = 10011009
	ErrLLMStreamInterrupted Code = 10011010
	ErrLLMPricingMissing    Code = 10011011
)

// 文件模块（模块代号 11）。业务侧文件校验与存储；网关层请求体上限见 ErrRequestBodyTooLarge。
const (
	ErrFileEmpty              Code = 10011101
	ErrFileTooLarge           Code = 10011102
	ErrFileTypeNotAllowed     Code = 10011103
	ErrFileCountExceeded      Code = 10011104
	ErrFileNotFound           Code = 10011105
	ErrFileChecksumMismatch   Code = 10011106
	ErrFileUploadFailed       Code = 10011107
	ErrFileStorageUnavailable Code = 10011108
)

// 邮件模块（模块代号 12）。
const (
	ErrEmailAlreadyBound       Code = 10011201
	ErrEmailNotVerified        Code = 10011202
	ErrEmailSendTooFrequent    Code = 10011203
	ErrEmailTemplateNotFound   Code = 10011204
	ErrEmailSendFailed         Code = 10011205
	ErrEmailServiceUnavailable Code = 10011206
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
	ErrOperationInProgress:   {http.StatusOK, "操作正在处理中，请稍后重试"},
	ErrRiskControlRejected:   {http.StatusOK, "操作被风控拦截"},
	ErrFeatureDisabled:       {http.StatusOK, "功能暂未开放"},

	ErrRequestFail: {http.StatusBadRequest, "请求失败"},
	ErrParams:      {http.StatusBadRequest, "参数不合法"},
	ErrInvalidJson: {http.StatusBadRequest, "JSON 格式错误"},

	ErrRequestBodyTooLarge:  {http.StatusRequestEntityTooLarge, "请求体过大"},
	ErrUnsupportedMediaType: {http.StatusUnsupportedMediaType, "不支持的内容类型"},

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
	ErrVersionConflict:       {http.StatusConflict, "数据已被修改，请刷新后重试"},

	ErrTooManyRequests: {http.StatusTooManyRequests, "请求过于频繁"},

	ErrInternalServer:   {http.StatusInternalServerError, "服务器内部错误"},
	ErrTimeout:          {http.StatusInternalServerError, "服务器响应超时"},
	ErrCanceled:         {http.StatusInternalServerError, "请求已取消"},
	ErrDeadlineExceeded: {http.StatusInternalServerError, "请求超时"},

	ErrServiceUnavailable:    {http.StatusServiceUnavailable, "服务不可用"},
	ErrDependencyUnavailable: {http.StatusServiceUnavailable, "依赖服务不可用"},
	ErrSystemMaintenance:     {http.StatusServiceUnavailable, "系统维护中"},

	// Elasticsearch（模块 01）
	ErrElasticsearchServer:        {http.StatusInternalServerError, "Elasticsearch 服务异常"},
	ErrElasticsearchDSL:           {http.StatusInternalServerError, "Elasticsearch 查询语句错误"},
	ErrElasticsearchIndexNotFound: {http.StatusInternalServerError, "Elasticsearch 索引不存在"},
	ErrElasticsearchBulkPartial:   {http.StatusInternalServerError, "Elasticsearch 批量操作部分失败"},

	// MySQL（模块 02）
	ErrMysqlServer:      {http.StatusInternalServerError, "MySQL 服务异常"},
	ErrMysqlSQL:         {http.StatusInternalServerError, "SQL 语句错误"},
	ErrMysqlQuery:       {http.StatusInternalServerError, "MySQL 查询失败"},
	ErrMysqlDeadlock:    {http.StatusInternalServerError, "MySQL 死锁"},
	ErrMysqlTransaction: {http.StatusInternalServerError, "MySQL 事务执行失败"},

	// MongoDB（模块 03）
	ErrMongoServer: {http.StatusInternalServerError, "MongoDB 服务异常"},
	ErrMongoDSL:    {http.StatusInternalServerError, "MongoDB 查询语句错误"},
	ErrMongoQuery:  {http.StatusInternalServerError, "MongoDB 查询失败"},

	// Redis（模块 04）
	ErrRedisServer:            {http.StatusInternalServerError, "Redis 服务异常"},
	ErrRedisQuery:             {http.StatusInternalServerError, "Redis 查询失败"},
	ErrRedisLockAcquireFailed: {http.StatusInternalServerError, "Redis 分布式锁获取失败"},
	ErrRedisLockLost:          {http.StatusInternalServerError, "Redis 分布式锁已丢失"},

	// Kafka（模块 05）
	ErrKafkaServer:         {http.StatusInternalServerError, "Kafka 服务异常"},
	ErrKafkaProducer:       {http.StatusInternalServerError, "Kafka 生产者异常"},
	ErrKafkaConsumer:       {http.StatusInternalServerError, "Kafka 消费者异常"},
	ErrKafkaMessageDecode:  {http.StatusInternalServerError, "Kafka 消息解析失败"},
	ErrKafkaRetryExhausted: {http.StatusInternalServerError, "Kafka 消息重试次数已耗尽"},

	// RabbitMQ（模块 06）
	ErrRabbitMQServer:         {http.StatusInternalServerError, "RabbitMQ 服务异常"},
	ErrRabbitMQProducer:       {http.StatusInternalServerError, "RabbitMQ 生产者异常"},
	ErrRabbitMQConsumer:       {http.StatusInternalServerError, "RabbitMQ 消费者异常"},
	ErrRabbitMQMessageDecode:  {http.StatusInternalServerError, "RabbitMQ 消息解析失败"},
	ErrRabbitMQRetryExhausted: {http.StatusInternalServerError, "RabbitMQ 消息重试次数已耗尽"},

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

	// 支付（模块 08）
	ErrPaymentNotFound:           {http.StatusOK, "支付单不存在"},
	ErrPaymentAlreadyPaid:        {http.StatusOK, "订单已支付，请勿重复支付"},
	ErrPaymentClosed:             {http.StatusOK, "支付单已关闭"},
	ErrPaymentExpired:            {http.StatusOK, "支付已超时，请重新下单"},
	ErrPaymentMethodNotSupported: {http.StatusOK, "不支持的支付方式"},
	ErrPayerIdentityMissing:      {http.StatusOK, "缺少支付方身份信息"},
	ErrPaymentAmountMismatch:     {http.StatusBadRequest, "支付金额与订单不一致"},
	ErrPaymentSignInvalid:        {http.StatusUnauthorized, "支付回调验签失败"},
	ErrPaymentCallbackInvalid:    {http.StatusBadRequest, "支付回调参数非法"},
	ErrPaymentCreateFailed:       {http.StatusInternalServerError, "发起支付失败"},
	ErrPaymentQueryFailed:        {http.StatusInternalServerError, "查询支付结果失败"},
	ErrPaymentChannelUnavailable: {http.StatusServiceUnavailable, "支付渠道暂不可用"},
	ErrRefundNotFound:            {http.StatusOK, "退款单不存在"},
	ErrRefundAmountExceeded:      {http.StatusOK, "退款金额超过可退金额"},
	ErrRefundInProgress:          {http.StatusOK, "退款处理中"},
	ErrRefundFailed:              {http.StatusInternalServerError, "退款失败"},

	// 账号与第三方登录（模块 09）
	ErrPhoneAlreadyBound:      {http.StatusConflict, "手机号已被其他账号绑定"},
	ErrThirdPartyAlreadyBound: {http.StatusConflict, "该第三方账号已绑定其他用户"},
	ErrThirdPartyNotBound:     {http.StatusOK, "请先绑定手机号"},
	ErrThirdPartyAuthFailed:   {http.StatusOK, "第三方授权失败"},
	ErrThirdPartyUnavailable:  {http.StatusServiceUnavailable, "第三方登录服务暂不可用"},
	ErrSmsSendTooFrequent:     {http.StatusOK, "验证码发送过于频繁"},
	ErrSmsDailyLimitExceeded:  {http.StatusOK, "今日验证码次数已用完"},
	ErrSmsSendFailed:          {http.StatusInternalServerError, "验证码发送失败"},

	// AI 大模型（模块 10）
	ErrLLMAuth:              {http.StatusInternalServerError, "AI 服务配置错误"},
	ErrLLMRateLimited:       {http.StatusServiceUnavailable, "AI 服务繁忙，请稍后重试"},
	ErrLLMTimeout:           {http.StatusInternalServerError, "AI 服务响应超时"},
	ErrLLMContextTooLong:    {http.StatusOK, "输入内容过长"},
	ErrLLMContentFiltered:   {http.StatusOK, "内容不符合规范"},
	ErrLLMInvalidRequest:    {http.StatusInternalServerError, "AI 请求参数错误"},
	ErrLLMProviderError:     {http.StatusServiceUnavailable, "AI 服务异常"},
	ErrLLMModelNotAvailable: {http.StatusOK, "所选模型暂不可用"},
	ErrLLMStructuredOutput:  {http.StatusInternalServerError, "AI 返回解析失败"},
	ErrLLMStreamInterrupted: {http.StatusInternalServerError, "AI 输出中断"},
	ErrLLMPricingMissing:    {http.StatusInternalServerError, "模型计费配置缺失"},

	// 文件（模块 11）
	ErrFileEmpty:              {http.StatusBadRequest, "文件为空"},
	ErrFileTooLarge:           {http.StatusBadRequest, "文件超过大小限制"},
	ErrFileTypeNotAllowed:     {http.StatusBadRequest, "不支持的文件类型"},
	ErrFileCountExceeded:      {http.StatusBadRequest, "文件数量超过限制"},
	ErrFileNotFound:           {http.StatusNotFound, "文件不存在"},
	ErrFileChecksumMismatch:   {http.StatusBadRequest, "文件校验失败"},
	ErrFileUploadFailed:       {http.StatusInternalServerError, "文件上传失败"},
	ErrFileStorageUnavailable: {http.StatusServiceUnavailable, "存储服务暂不可用"},

	// 邮件（模块 12）
	ErrEmailAlreadyBound:       {http.StatusConflict, "邮箱已被绑定"},
	ErrEmailNotVerified:        {http.StatusOK, "邮箱未验证"},
	ErrEmailSendTooFrequent:    {http.StatusOK, "邮件发送过于频繁"},
	ErrEmailTemplateNotFound:   {http.StatusInternalServerError, "邮件模板不存在"},
	ErrEmailSendFailed:         {http.StatusInternalServerError, "邮件发送失败"},
	ErrEmailServiceUnavailable: {http.StatusServiceUnavailable, "邮件服务暂不可用"},
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
