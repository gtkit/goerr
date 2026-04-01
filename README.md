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

## New、WrapStatus、Wrap：内部规范（必读）

| API | 何时使用 |
| --- | --- |
| **`New(err, *Status, msgs...)`** | **业务代码里新建「带业务码」的错误**：有/无底层 `err` 均可；需要自定义对外文案时用 `msgs`。 |
| **`WrapStatus(err, *Status)`** | **已有任意 `error`，在中间件 / 边界统一补上业务状态**（例如把第三方库错误映射到本服务的 `Status`）。`status == nil` 时等价于 `StatusInternalServer()`。 |
| **`Wrap` / `Wrapf`** | **仅追加英文/中文上下文**，不参与「赋业务码」。若内层已是 `*Item`，会保留其 `Status` 与对外 `Message()`。 |

**禁止（除非经过网关纠正）：**

- **`Wrap` / `Wrapf` 作用在「非 `*Item`」的普通 `error` 上时**，得到的结果 **没有业务 `Status`**（`Code()` 会落在默认 OK 语义上），**不得直接作为「业务错误」交给统一响应层返回给客户端**。
- 正确做法：先用 **`WrapStatus(err, goerr.StatusXxx())`** 或 **`New(err, goerr.StatusXxx())`** 赋予业务码，再视需要 **`Wrap`** 追加调用链说明；或在网关根据错误类型分支映射到合法 `*Item`。

**动态类型（配置、反射、`*Status` 与 `func() *Status` 混用）**：使用 **`ResolveStatus`**、**`NewFrom`**、**`NewfFrom`**、**`WrapStatusFrom`**，通过 **`error`** 表达失败，**避免**向 `Newf` 等传入错误类型。

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

- 统一响应体内默认返回 HTTP `200`
- 客户端通过业务错误码判断成功或失败
- 仅在请求未进入统一响应体，或网关/框架必须表达协议级语义时，才返回非 `200`
- 错误码一经发布不得复用；废弃码保留占位

---

## 类型与常用 API

- **`Code`**：业务错误码，见 `errcode.go` 中常量与 `Code.Message()` 默认文案。
- **`Status`**：不可变状态（码 + HTTP 状态 + 默认消息），通过 **`StatusXxx()`** 取预构建指针，或 **`NewStatus`** 自定义。
- **`Item`**：实现了 `error` 与 `goerr.Error`，是实际在函数间传递的值。

构造方式简述：

- **`New(err, status *Status, msgs...)`**：最常见；`status` 传 **`StatusMysqlServer()`** 的返回值。`status == nil` 时视为 **`StatusOK()`**。
- **`NewFn(err, statusFn, msgs...)`**：传入 **`StatusMysqlServer`** 函数本身，等价于在内部调用 `statusFn()`。
- **`Newf(status *Status, format, args...)`**：无底层 `cause`，仅格式化一条消息；`status == nil` 时为 **`StatusOK()`**。
- **`NewFrom` / `NewfFrom`**：需要 **`ResolveStatus`** 语义时使用（返回 **`error`**）。
- **`Wrap` / `Wrapf`**：在已有 `error` 上追加上下文；见上文「禁止」说明。
- **`WrapStatus(err, *Status)`**：给任意错误补上业务状态；`nil` 时用 **`StatusInternalServer()`**。
- **`WrapStatusFrom`**：与 `WrapStatus` 相同语义，但 `status` 可为 `any` 并由 **`ResolveStatus`** 解析。

---

## 快速使用

```go
import "github.com/gtkit/goerr"

// 创建带状态码的错误（第二个参数为 *Status，请调用 StatusXxx()）
err := goerr.New(err, goerr.StatusMysqlServer())

// 自定义对外文案（写入 Message()；若存在 cause，Error() 仍会附带底层错误）
err = goerr.New(err, goerr.StatusParams(), "user_id is required")

// 纯状态错误（无底层 cause）
err = goerr.New(nil, goerr.StatusNotFound(), "order not found")

// 仅格式化消息、无 cause
err = goerr.Newf(goerr.StatusParams(), "field %q is required", "name")

// 包装已有错误（日志更详细；若内层是 *Item，对外 Message 仍读内层）
err = goerr.Wrap(err, "query users")

// 为错误补充业务状态码（适合「任意 error → 业务码」）
err = goerr.WrapStatus(err, goerr.StatusRedisServer())

// 组装统一响应（示意：用 Message，不要用 Error 给前端）
if item, ok := goerr.AsItem(err); ok {
	code := item.Code()       // 业务错误码 → JSON `code`
	http := item.HTTPCode() // 协议层 HTTP 状态码（多数业务场景仍为 200）
	msg := item.Message()     // → JSON `message`，给客户端看（建议再经 SanitizeForClient）
	_ = code
	_ = http
	_ = msg
}

// 自定义业务错误码（需先保证 Code 合法）
st := goerr.NewStatus(goerr.Code(10010701), 200, "Order expired")
_ = st
```

绝大多数业务错误建议沿用预构建的 `StatusXxx()`，它们默认使用 HTTP `200`。只有协议级或网关级场景，才建议通过 `NewStatus` 指定非 `200` HTTP 状态码。

MySQL 查询「没有这条数据」时，建议使用 `StatusRecordNotFound()`；`StatusMysqlQuery()` 只用于 SQL 已执行但查询过程本身出错。

若注册自定义错误码，可先调用 `ValidateCode` 做预检查；`NewStatus` 遇到非法错误码会直接 panic。

```go
code := goerr.Code(10010701)
if err := goerr.ValidateCode(code); err != nil {
	return err
}
status := goerr.NewStatus(code, 200, "Order expired")
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

## 测试矩阵与发布

详见 **[CONTRIBUTING.md](./CONTRIBUTING.md)**（`go test`、`race`、fuzz、`go vet`）。  
版本历史与迁移说明见 **[CHANGELOG.md](./CHANGELOG.md)**。

---

## 设计要点

- **零外部依赖** — 仅使用标准库
- **Status 预构建单例** — `StatusXxx()` 返回包级别指针，零堆分配；预构建业务错误默认 HTTP `200`
- **Item 不可变** — `Wrap` 不修改原始错误，并发安全无需锁
- **Message / Error 分工** — 见上文「Message() 与 Error()」约定
- **严格类型** — `Newf` / `WrapStatus` 使用 `*Status`；动态解析用 `*From` / `ResolveStatus`
- **Unwrap 支持** — 完整兼容 `errors.Is` / `errors.As` 错误链
- **`%+v` 调用栈** — 配合 zap 等日志库输出完整堆栈与 cause
