# goerr

统一业务错误码处理包，提供结构化的业务错误码、HTTP 响应策略和错误消息封装。  
与标准库 `errors` 配合，支持错误链、`errors.Is` / `errors.As` 以及 `%+v` 输出调用栈。

**版本说明与破坏性变更见 [CHANGELOG.md](./CHANGELOG.md)。**  
**合并与发版前测试清单见 [CONTRIBUTING.md](./CONTRIBUTING.md)。**

---

## Message() 与 Error()：团队约定（必读）

`goerr` 在 `*Item` 上同时实现了面向「对外展示」与「对内排障」的两套字符串，**用途不同，不要混用**。

| 方法 | 用途 | 典型使用位置 |
| --- | --- | --- |
| **`Message() string`** | **给客户端 / 调用方看的短文案** | 统一 JSON 响应里的 `message` 字段、移动端提示、OpenAPI 文档中的用户可见说明 |
| **`Error() string`** | **给服务端日志、告警、链路追踪看的完整描述** | `log.Print(err)`、`zap.Error(err)`、`fmt.Errorf("...: %w", err)`、可含底层 `cause` 拼接后的全文 |

**约定说明：**

1. **`Message()`** 应尽量稳定、可读、**不暴露内部实现细节**（例如不要把数据库驱动返回的英文原句直接当作主提示）。有底层错误时，`New` 会把「自定义或 Status 默认文案」放在 `Message()`，把「文案 + 底层错误」放在 `Error()` 供日志使用。
2. **`Error()`** 实现标准库的 `error` 接口，适合**排障**：通常包含更多上下文，例如 `New(err, ...)` 时会把 `err.Error()` 拼进全文。
3. 业务网关或 BFF 组装 HTTP 响应时，请用 **`Code()` + `Message()` + `HTTPCode()`**（以及你们协议里约定的字段名）；**不要**把 `Error()` 整段直接回给前端，除非你们明确允许暴露内部信息。
4. 开发排查问题时，请打印 **`Error()`** 或使用 **`fmt.Sprintf("%+v", item)`**（若需栈与 cause 信息），**不要**只看 `Message()`，否则会丢掉底层原因。

下面示例演示「响应体用 Message，日志用 Error」：

```go
item := goerr.New(someErr, goerr.StatusMysqlServer(), "数据库繁忙")

// 返回给客户端的 JSON（示意）
// { "code": 10010201, "message": "数据库繁忙" }  ← 使用 item.Message()
// 不要误用 item.Error() 作为对外 message

// 服务端日志
log.Printf("request failed: %v", item.Error())
// 输出中会包含底层 someErr，便于排查
```

---

## 调用方式

按代码所在的层，只需要记三个入口：

| 层 | 写法 | 说明 |
| --- | --- | --- |
| **业务层报错** | `return goerr.New(err, goerr.StatusXxx(), "给客户端看的文案")` | `err` 可为 nil；文案可省，省略时用状态默认文案 |
| **中间层加上下文** | `return goerr.Wrap(err, "load profile")` | `err == nil` 时返回真 nil，成功路径直接透传；内层 `*Item` 的状态与对外文案原样保留 |
| **响应层组装 JSON** | `e := goerr.Convert(err)` | 永不返回 nil，成功与失败走同一条路径 |

从 repository 到 handler 的完整链路：

```go
// repository：把底层错误翻译成业务语义
func (r *userRepo) Get(ctx context.Context, id int64) (*User, error) {
	u, err := r.db.Get(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, goerr.New(err, goerr.StatusRecordNotFound(), "用户不存在")
	}
	if err != nil {
		return nil, goerr.New(err, goerr.StatusMysqlQuery())
	}
	return u, nil
}

// service：只加上下文、不改状态。err 为 nil 时 Wrap 返回 nil，不必先判
func (s *userService) Profile(ctx context.Context, id int64) (*Profile, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, goerr.Wrap(err, "load profile")
	}
	return toProfile(u), nil
}

// handler：需要时按业务码分支
func (h *handler) Profile(c *gin.Context) {
	p, err := h.svc.Profile(c, id)
	if goerr.HasCode(err, goerr.ErrRecordNotFound) {
		// 例如：埋点、走兜底数据
	}
	respond(c, p, err)
}

// 统一响应：nil / *Item / 被 fmt.Errorf 包过的 *Item / 第三方 error，一条路径
func respond(c *gin.Context, data any, err error) {
	e := goerr.Convert(err)
	c.JSON(e.HTTPCode(), gin.H{"code": e.Code(), "message": e.Message(), "data": data})
	// 日志请打原始 err（%+v 带调用栈），不要打 Convert 的结果
}
```

这条链上，客户端收到 `{"code":10010026,"message":"用户不存在"}`；日志里 `%+v` 打出
`load profile: 用户不存在: sql: no rows in result set`，后面跟 repository 那一行起的调用栈——
整条链只有一份栈，就是最接近错误起源的那份。

### 其他入口

| 入口 | 何时用 |
| --- | --- |
| `Newf(status, format, args...)` | 只想格式化一条文案、没有底层 `err` |
| `Wrapf(err, format, args...)` | `Wrap` 的格式化版 |
| `WrapStatus(err, status)` | 中间件 / 边界把第三方 `error` 映射成本服务的业务状态 |
| `WithStack(err)` | 只想给普通 `error` 补一份栈、不加文案 |
| `HasCode(err, code)` | 按业务码分支。同一失败语义在不同位置构造出的是不同 `*Item` 实例，`errors.Is` 判不出相等 |
| `AsItem(err) (*Item, bool)` | 需要区分「链上有没有 `*Item`」时用；只是取字段请用 `Convert` |

### 三条规则

1. **nil 就是 nil。** `Wrap` / `Wrapf` / `WrapStatus` / `WithStack` 在 `err == nil` 时返回接口值真 nil；`Convert(nil)` 返回成功状态；`*Item` 的方法对 nil receiver 也安全。
2. **状态缺失按失败处理。** 任何入口 `status == nil`、或包装的是不带状态的普通 `error`，一律按 `StatusInternalServer()` 兜底——漏传状态的错误不会被读成成功。
3. **`Message()` 给客户端，`Error()` 给日志。** `New(err, st, "文案")` 把「文案」放进 `Message()`，把「文案: err 原文」放进 `Error()`；`Wrap` 只叠 `Error()`，`Message()` 保持内层的。

### Convert 的三种输入

| 输入 | 返回 |
| --- | --- |
| `nil` | 成功状态（`StatusOK()`） |
| 链上有 `*Item` | 该 `*Item`，自定义对外文案原样保留 |
| 普通 `error` | 按 `StatusInternalServer()` 包装；`Message()` 是状态默认文案，`Error()` 保留原文 |

---

## 网关：对 Message() 做截断与脱敏

下游传入的自定义 `msgs`、或未来扩展的文案，**不得未经处理**直接进公网响应。建议在 **BFF / API 网关** 对 **`Message()`** 的结果调用：

```go
msg := item.Message()
safe := goerr.SanitizeForClient(msg, goerr.SanitizeOptions{
	MaxRunes:    256,              // 按产品约定调整；0 或缺省时用 DefaultMaxMessageRunes
	RedactPhone: true,
	RedactEmail: true,
})
// 将 safe 写入 JSON `message` 字段
```

本包提供 **手机号、邮箱** 的简单占位替换与 **Unicode 安全截断**；身份证、银行卡、SQL 片段、内部 token 等请在业务层或专用网关策略中扩展。

---

## 日志：Error() 与 PII

`Error()` 常包含 **cause 链全文**，容易带上 **手机号、邮箱、SQL、路径** 等敏感信息。建议：

1. 落库/落盘前对字符串调用 **`goerr.RedactForLog(err.Error())`**（或接入日志平台的脱敏管道）；  
2. 与 **日志分级、留存周期、访问审计** 策略一起使用，**不要**把原始 `Error()` 直接暴露给非运维角色。

---

## 错误码规范

```
项目组代号(10) + 服务代号(01) + 模块代号(0~99) + 错误码(0~99)
```

- `HTTPCode()` 按业务定制返回语义化 HTTP 状态码：业务态可为 `200`，参数/认证/权限/资源/限流/服务端错误分别映射到对应 HTTP 大类
- 客户端通过 HTTP 状态码判断错误大类，通过响应体中的业务错误码判断具体错误
- 若某个服务采用“所有统一响应体都返回 HTTP 200”的协议，可在网关层忽略 `HTTPCode()`，只使用 `Code()` 与 `Message()`
- 错误码一经发布不得复用；废弃码保留占位

已分配的模块：

| 模块 | 用途 | 模块 | 用途 |
| --- | --- | --- | --- |
| `00` | 通用 | `07` | 订单 |
| `01` | Elasticsearch | `08` | 支付（含退款） |
| `02` | MySQL | `09` | 账号与第三方登录 |
| `03` | MongoDB | `10` | AI 大模型 |
| `04` | Redis | `11` | 文件 |
| `05` | Kafka | `12` | 邮件 |
| `06` | RabbitMQ | | |

完整清单（码 → HTTP → 默认文案）见 `errcode.go` 的 `codeTable`，它是唯一真源：`Code.Message()`
与预构建的 `StatusXxx()` 都从那里取值。

**选码原则**：基础设施模块（`01`~`06`）的码描述「哪个组件的哪种故障」，对外文案由调用方在 `Message()`
里补充业务含义；业务模块（`07`~`12`）的码直接承载业务语义，可以原样返给客户端。

---

## 类型与常用 API

- **`Code`**：业务错误码，见 `errcode.go` 中常量与 `Code.Message()` 默认文案。
- **`Status`**：错误状态（码 + HTTP 状态 + 默认消息），字段私有、构造后不再写入，通过 **`StatusXxx()`** 取预构建指针（全进程共享，勿解引用赋值），或通过 **`NewStatus`** / **`MustNewStatus`** 自定义。
- **`Item`**：实现了 `error` 与 `goerr.Error`，是实际在函数间传递的值。

---

## 自定义错误码与状态

绝大多数业务错误建议沿用预构建的 `StatusXxx()`；它们的 HTTP 状态码按本包业务约定语义化映射。需要自定义状态时，用 `NewStatus` 或 `MustNewStatus`：

```go
// NewStatus 不校验 code，适合业务方自由编码
st := goerr.NewStatus(goerr.Code(900001), 200, "Custom business error")
err := goerr.New(nil, st, "自定义文案")
```

MySQL 查询「没有这条数据」时，建议使用 `StatusRecordNotFound()`；`StatusMysqlQuery()` 只用于 SQL 已执行但查询过程本身出错。

若注册自定义错误码，可先调用 `ValidateCode` 做预检查；若希望包级变量初始化时强制校验，使用 `MustNewStatus`，非法 code 会直接 panic。

```go
code := goerr.Code(10010701)
if err := goerr.ValidateCode(code); err != nil {
	return err
}
status := goerr.MustNewStatus(code, 200, "Order expired")
_ = status
```

---

## 数据库错误映射

推荐把数据库错误先翻译成业务语义，再交给 `goerr`：

| 场景 | 推荐状态码 |
| --- | --- |
| `sql.ErrNoRows` / `gorm.ErrRecordNotFound` | `StatusRecordNotFound()` |
| 唯一键冲突、重复创建 | `StatusAlreadyExists()` 或 `StatusConflict()` |
| SQL 拼接错误、语句非法 | `StatusMysqlSQL()` |
| 查询执行失败、扫描失败 | `StatusMysqlQuery()` |
| 连接不可用、驱动异常、实例不可达 | `StatusMysqlServer()` |

`StatusNotFound()` 更适合路由、页面或资源入口不存在；数据库查询为空请优先使用 `StatusRecordNotFound()`。

```go
import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

func mapDBError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, gorm.ErrRecordNotFound):
		return goerr.New(err, goerr.StatusRecordNotFound(), "user not found")
	default:
		return goerr.New(err, goerr.StatusMysqlQuery(), "query failed")
	}
}
```

---

## AI 大模型错误映射

`goerr` 零依赖，不引入 provider 包，映射写在业务侧的 adapter 里。模块 `10` 的码按
**「谁能解决」** 划分 HTTP：用户能自行修正的给 `200`，我方配置或实现问题给 `500`，供应商侧故障给 `503`。

```go
import "github.com/gtkit/go-llm-provider/v2/provider"

func fromLLM(err error) error {
	if err == nil {
		return nil
	}
	if pe, ok := errors.AsType[*provider.ProviderError](err); ok {
		switch pe.Code {
		case provider.ErrorCodeAuth:
			return goerr.New(err, goerr.StatusLLMAuth()) // 我们的 key 配错了，对用户是内部错误
		case provider.ErrorCodeRateLimit:
			return goerr.New(err, goerr.StatusLLMRateLimited())
		case provider.ErrorCodeTimeout:
			return goerr.New(err, goerr.StatusLLMTimeout())
		case provider.ErrorCodeContextLength:
			return goerr.New(err, goerr.StatusLLMContextTooLong()) // 用户能改，给 200
		case provider.ErrorCodeContentFilter:
			return goerr.New(err, goerr.StatusLLMContentFiltered())
		case provider.ErrorCodeInvalidRequest:
			return goerr.New(err, goerr.StatusLLMInvalidRequest()) // 我们拼错了请求
		}
	}
	switch {
	case errors.Is(err, provider.ErrQuotaExceeded):
		return goerr.New(err, goerr.StatusQuotaExceeded()) // 用户配额，复用通用码
	case errors.Is(err, provider.ErrStructuredDecode), errors.Is(err, provider.ErrStructuredValidation):
		return goerr.New(err, goerr.StatusLLMStructuredOutput())
	case errors.Is(err, provider.ErrModelNotPriced), errors.Is(err, provider.ErrInvalidPricing):
		return goerr.New(err, goerr.StatusLLMPricingMissing())
	}
	return goerr.New(err, goerr.StatusLLMProviderError())
}
```

`ProviderError` 的 `Retryable` / `RetryAfter` 是重试策略的输入，留在 adapter 里决定是否重试，不进错误码。

---

## 支付回调：用 HTTPCode() 控制渠道重试

微信 / 支付宝的异步通知响应体格式由渠道规定，不走统一 JSON；但 `HTTPCode()` 仍然有用——
**验签失败返 4xx 让渠道停止重试，我方处理失败返 5xx 让渠道继续重试**：

```go
func notify(c *gin.Context) {
	err := svc.HandlePayNotify(c.Request)
	e := goerr.Convert(err)
	if err != nil {
		// 金额不符、验签失败是安全事件，日志级别拉高并接告警
		if goerr.HasCode(err, goerr.ErrPaymentSignInvalid) || goerr.HasCode(err, goerr.ErrPaymentAmountMismatch) {
			logger.Error("payment security event", zap.Error(err))
		}
	}
	c.String(e.HTTPCode(), channelAck(e.Code())) // 4xx 停止重试，5xx 继续重试
}
```

---

## 消息队列：用错误码作为死信判据

`ErrXxxMessageDecode`（反序列化失败，重投多少次都不会成功）与 `ErrXxxRetryExhausted`
（重试次数耗尽）是**进死信队列的两个判据**，有独立码才能在死信里按码归类：

```go
func (c *consumer) handle(ctx context.Context, msg []byte) error {
	var evt OrderEvent
	if err := json.Unmarshal(msg, &evt); err != nil {
		return goerr.New(err, goerr.StatusKafkaMessageDecode()) // 不可重试，直接进死信
	}
	if err := c.svc.Apply(ctx, evt); err != nil {
		if attempts >= maxRetry {
			return goerr.Wrap(err, "retry exhausted") // 保留原始业务码，外层按 attempts 决定进死信
		}
		return err
	}
	return nil
}
```

---

## 测试矩阵与发布

详见 **[CONTRIBUTING.md](./CONTRIBUTING.md)**（`go test`、`race`、fuzz、`go vet`）。  
版本历史与迁移说明见 **[CHANGELOG.md](./CHANGELOG.md)**。

---

## 设计要点

- **零外部依赖** — 仅使用标准库
- **错误码单一真源** — 业务码的 HTTP 状态码与默认文案都登记在 `errcode.go` 的 `codeTable`，`Code.Message()` 与预构建 `Status` 共用；漏配元信息的码在包初始化期直接 panic
- **Status 预构建单例** — `StatusXxx()` 返回包级别指针，零堆分配。该指针被全进程共享，**请勿对其解引用赋值**（`*goerr.StatusOK() = ...` 会污染所有调用方看到的状态）
- **Item 构造后不再写入** — `Wrap` 返回新错误，不修改原始错误，可跨 goroutine 共享
- **nil 语义干净** — `Wrap` / `Wrapf` / `WrapStatus` / `WithStack` 返回 `error`，`err == nil` 时是接口值真 nil；`*Item` 的方法对 nil receiver 也安全，不会 panic
- **Message / Error 分工** — 见上文「Message() 与 Error()」约定
- **状态缺失一律按失败处理** — 所有入口在 `status == nil` 或包装普通 `error` 时使用 `StatusInternalServer()`，避免成功码误判
- **严格类型** — 所有入口只接受 `*Status`，状态来源在编译期就确定
- **Unwrap 支持** — 完整兼容 `errors.Is` / `errors.As` 错误链；按业务码判定用 `HasCode`，响应层一行 `Convert`
- **`%+v` 调用栈** — 一条错误链共享一份栈（最内层、最接近起源的那份），用 `runtime.CallersFrames` 解析以正确还原内联帧与行号；`Error()` 文本本身已含整条 cause 链，因此不重复递归打印
