# goerr

统一业务错误码处理包，提供结构化的业务错误码、HTTP 响应策略和错误消息封装。

## 错误码规范

```
项目组代号(10) + 服务代号(01) + 模块代号(0~99) + 错误码(0~99)
```

- 统一响应体内默认返回 HTTP `200`
- 客户端通过业务错误码判断成功或失败
- 仅在请求未进入统一响应体，或网关/框架必须表达协议级语义时，才返回非 `200`
- 错误码一经发布不得复用；废弃码保留占位

## 快速使用

```go
import "github.com/gtkit/goerr"

// 创建带状态码的错误
err := goerr.New(err, goerr.StatusMysqlServer)

// 自定义消息
err := goerr.New(err, goerr.StatusParams, "user_id is required")

// 纯状态错误（无底层 cause）
err := goerr.New(nil, goerr.StatusNotFound, "order not found")

// 格式化消息
err := goerr.Newf(goerr.StatusParams, "field %q is required", "name")

// 包装已有错误
err := goerr.Wrap(err, "query users")

// 为错误补充状态码
err := goerr.WrapStatus(err, goerr.StatusRedisServer)

// 提取错误信息
if item, ok := goerr.AsItem(err); ok {
    code := item.Code()         // 业务错误码
    http := item.HTTPStatus()   // HTTP 状态码
    msg  := item.Message()      // 状态消息
}

// 自定义业务错误码
var StatusOrderExpired = goerr.NewStatus(goerr.Code(10010701), 200, "Order expired")
```

`StatusXxx()` 写法仍然兼容；`New` / `Newf` / `WrapStatus` 现在也支持直接传 `StatusXxx`。

绝大多数业务错误建议沿用预构建的 `StatusXxx()`，它们默认使用 HTTP `200`。只有协议级或网关级场景，才建议通过 `NewStatus` 指定非 `200` HTTP 状态码。

MySQL 查询“没有这条数据”时，建议使用 `StatusRecordNotFound()`；`StatusMysqlQuery()` 只用于 SQL 已执行但查询过程本身出错。

如果你要注册自定义错误码，可以先调用 `ValidateCode` 做预检查；`NewStatus` 遇到非法错误码会直接 panic。

```go
code := goerr.Code(10010701)
if err := goerr.ValidateCode(code); err != nil {
	return err
}
status := goerr.NewStatus(code, 200, "Order expired")
_ = status
```

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
		return goerr.New(err, goerr.StatusRecordNotFound, "user not found")
	default:
		return goerr.New(err, goerr.StatusMysqlQuery)
	}
}
```

## 设计要点

- **零外部依赖** — 仅使用标准库
- **Status 预构建单例** — `StatusXxx()` 返回包级别指针，零堆分配；预构建业务错误默认 HTTP `200`
- **Item 不可变** — `Wrap` 不修改原始错误，并发安全无需锁
- **Unwrap 支持** — 完整兼容 `errors.Is` / `errors.As` 错误链
- **`%+v` 调用栈** — 配合 zap 等日志库输出完整堆栈
