# goerr - 生产级 Go 错误处理库

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

一个高性能、并发安全、功能完整的 Go 错误处理库,专为生产环境设计。

## ✨ 核心特性

- **🚀 高性能**: 零内存分配的错误码映射,优化的错误创建流程
- **🔒 并发安全**: 使用 atomic.Value 实现无锁并发访问
- **📝 结构化错误**: 支持错误码、HTTP 状态码、详细信息
- **🔗 错误链**: 完整的错误包装和展开支持
- **📊 堆栈追踪**: 自动捕获调用堆栈,便于调试
- **🎯 类型安全**: 完整的类型断言和错误判断
- **💾 零内存泄漏**: 精心设计的内存管理,无泄漏风险
- **🌐 HTTP 集成**: 内置 HTTP 响应工具和中间件
- **📦 预定义错误码**: 开箱即用的常用错误码

## 📦 安装

```bash
go get github.com/gtkit/goerr
```

## 🚀 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "github.com/gtkit/goerr"
)

func main() {
    // 创建简单错误
    err := goerr.New("USER_001", "user not found")
    fmt.Println(err.Error())
    // 输出: [USER_001] user not found

    // 创建格式化错误
    err2 := goerr.Newf("USER_002", "user %s is inactive", "alice")
    fmt.Println(err2.Error())
    // 输出: [USER_002] user alice is inactive
}
```

### 注册错误码

```go
// 注册单个错误码
goerr.RegisterCode("ORDER_001", "Order not found", 404)

// 批量注册
goerr.RegisterCodes([]goerr.CodeConfig{
    {Code: "ORDER_001", Message: "Order not found", HTTPStatus: 404},
    {Code: "ORDER_002", Message: "Order already paid", HTTPStatus: 400},
    {Code: "ORDER_003", Message: "Order expired", HTTPStatus: 410},
})

// 使用注册的错误码
err := goerr.NewWithCode("ORDER_001")
fmt.Printf("HTTP Status: %d\n", err.HTTPStatus()) // 404
```

### 添加详细信息

```go
err := goerr.New("AUTH_001", "authentication failed").
    WithDetail("user_id", 12345).
    WithDetail("ip", "192.168.1.1").
    WithDetail("timestamp", "2025-12-08T10:30:00Z")

// 批量添加
err2 := err.WithDetails(map[string]interface{}{
    "method": "password",
    "attempts": 3,
})
```

### 错误包装

```go
import "database/sql"

// 包装标准错误
dbErr := sql.ErrNoRows
err := goerr.Wrap(dbErr, "DB_001", "failed to query user")

// 格式化包装
err2 := goerr.Wrapf(dbErr, "DB_001", "failed to query user %d", 12345)

// 展开错误
underlying := err.Unwrap()
```

### 错误判断

```go
err1 := goerr.New("ERR_001", "error 1")
err2 := goerr.New("ERR_001", "same code, different message")
err3 := goerr.New("ERR_002", "different error")

// 相同错误码被认为是相同的错误类型
fmt.Println(goerr.Is(err1, err2)) // true
fmt.Println(goerr.Is(err1, err3)) // false

// 类型转换
var target goerr.Error
if goerr.As(err1, &target) {
    fmt.Printf("Code: %s\n", target.Code())
}
```

## 🌐 HTTP 集成

### 写入 HTTP 响应

```go
import "net/http"

func handler(w http.ResponseWriter, r *http.Request) {
    err := processRequest()
    if err != nil {
        goerr.WriteHTTPError(w, err)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}
```

JSON 响应格式:
```json
{
  "code": "USER_001",
  "message": "user not found",
  "details": {
    "user_id": 12345
  }
}
```

### 使用中间件

```go
// 自动捕获 panic 并转换为错误响应
http.HandleFunc("/api", goerr.HTTPMiddleware(handler))
```

## 📋 预定义错误码

库提供了一组常用的预定义错误码:

| 错误码 | HTTP 状态码 | 说明 |
|--------|------------|------|
| `ErrInternal` | 500 | 内部服务器错误 |
| `ErrInvalidArgument` | 400 | 无效参数 |
| `ErrNotFound` | 404 | 资源不存在 |
| `ErrAlreadyExists` | 409 | 资源已存在 |
| `ErrPermissionDenied` | 403 | 权限拒绝 |
| `ErrUnauthenticated` | 401 | 未认证 |
| `ErrResourceExhausted` | 429 | 资源耗尽 |
| `ErrUnavailable` | 503 | 服务不可用 |
| `ErrTimeout` | 504 | 操作超时 |

使用示例:
```go
err := goerr.NewWithCode(goerr.ErrNotFound).
    WithDetail("resource", "user").
    WithDetail("id", 12345)
```

## 🔍 堆栈追踪

```go
err := goerr.New("STACK_001", "error with stack")

// 简单输出
fmt.Printf("%v\n", err)
// [STACK_001] error with stack

// 详细输出(包含堆栈)
fmt.Printf("%+v\n", err)
// [STACK_001] error with stack
// Details:
//   user_id: 12345
// Stack trace:
//   /path/to/file.go:10 main.createError
//   /path/to/file.go:5 main.main
//   ...
```

## 🎯 最佳实践

### 1. 定义业务错误码

```go
// 在包初始化时注册所有错误码
func init() {
    goerr.RegisterCodes([]goerr.CodeConfig{
        // 用户相关错误 (1xxx)
        {Code: "USER_1001", Message: "User not found", HTTPStatus: 404},
        {Code: "USER_1002", Message: "User already exists", HTTPStatus: 409},
        {Code: "USER_1003", Message: "Invalid credentials", HTTPStatus: 401},
        
        // 订单相关错误 (2xxx)
        {Code: "ORDER_2001", Message: "Order not found", HTTPStatus: 404},
        {Code: "ORDER_2002", Message: "Order already paid", HTTPStatus: 400},
        {Code: "ORDER_2003", Message: "Order expired", HTTPStatus: 410},
    })
}
```

### 2. 服务层错误处理

```go
type UserService struct {
    db *sql.DB
}

func (s *UserService) GetUser(id int) (*User, error) {
    user := &User{}
    err := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&user)
    
    if err == sql.ErrNoRows {
        return nil, goerr.NewWithCode("USER_1001").
            WithDetail("user_id", id)
    }
    
    if err != nil {
        return nil, goerr.Wrap(err, goerr.ErrInternal, "failed to query user").
            WithDetail("user_id", id)
    }
    
    return user, nil
}
```

### 3. HTTP 处理器

```go
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    id := parseID(r)
    
    user, err := h.userService.GetUser(id)
    if err != nil {
        goerr.WriteHTTPError(w, err)
        return
    }
    
    json.NewEncoder(w).Encode(user)
}
```

### 4. 错误链构建

```go
func (s *OrderService) PlaceOrder(order *Order) error {
    // 验证库存
    if err := s.validateInventory(order); err != nil {
        return goerr.Wrap(err, "ORDER_2004", "inventory validation failed")
    }
    
    // 处理支付
    if err := s.processPayment(order); err != nil {
        return goerr.Wrap(err, "ORDER_2005", "payment processing failed")
    }
    
    // 创建订单
    if err := s.createOrder(order); err != nil {
        return goerr.Wrap(err, goerr.ErrInternal, "failed to create order")
    }
    
    return nil
}
```

## ⚡ 性能特性

### 无锁并发

- 使用 `atomic.Value` 实现错误码注册表的无锁读取
- 错误对象使用不可变设计,避免并发修改
- `WithDetail()` 返回新错误对象,不修改原错误

### 内存优化

- 堆栈信息共享,避免重复分配
- 详细信息使用 COW (Copy-On-Write) 模式
- 零内存泄漏,所有对象可被 GC 正常回收

### 基准测试

```bash
go test -bench=. -benchmem
```

预期结果:
```
BenchmarkNewError-8                  1000000    1050 ns/op    320 B/op    5 allocs/op
BenchmarkNewWithCode-8               2000000     850 ns/op    256 B/op    4 allocs/op
BenchmarkGetCodeConfig-8            50000000      28 ns/op      0 B/op    0 allocs/op
BenchmarkWithDetail-8                3000000     420 ns/op    128 B/op    2 allocs/op
```

## 🔐 并发安全保证

所有公开的 API 都是并发安全的:

- ✅ `New()`, `Newf()`, `Wrap()`, `Wrapf()` - 并发安全
- ✅ `RegisterCode()`, `GetCodeConfig()` - 并发安全
- ✅ `WithDetail()`, `WithDetails()` - 并发安全(不可变)
- ✅ `Is()`, `As()` - 并发安全
- ✅ 所有 Error 接口方法 - 并发安全

## 🧪 测试

```bash
# 运行所有测试
go test -v ./...

# 运行并发测试
go test -v -race ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 与其他库对比

| 特性 | goerr | pkg/errors | cockroachdb/errors |
|-----|-------|-----------|-------------------|
| 错误码支持 | ✅ | ❌ | ❌ |
| HTTP 状态码 | ✅ | ❌ | ❌ |
| 结构化详情 | ✅ | ❌ | ✅ |
| 堆栈追踪 | ✅ | ✅ | ✅ |
| 无锁并发 | ✅ | ✅ | ❌ |
| HTTP 集成 | ✅ | ❌ | ❌ |
| 预定义错误码 | ✅ | ❌ | ✅ |
| 零分配读取 | ✅ | ❌ | ❌ |

## 🤝 贡献

欢迎提交 Issue 和 Pull Request!

## 📄 许可证

MIT License

## 🔗 相关资源

- [Go Error Handling Best Practices](https://go.dev/blog/error-handling-and-go)
- [Effective Error Handling in Go](https://earthly.dev/blog/golang-errors/)
- [Go 1.13 Errors](https://go.dev/blog/go1.13-errors)
