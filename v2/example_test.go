package goerr

import (
	"database/sql"
	"fmt"
	"net/http"
)

// Example_basic 基础使用示例
func Example_basic() {
	// 创建简单错误
	err := New("USER_001", "user not found")
	fmt.Println(err.Error())

	// 创建格式化错误
	err2 := Newf("USER_002", "user %s is inactive", "alice")
	fmt.Println(err2.Error())

	// Output:
	// [USER_001] user not found
	// [USER_002] user alice is inactive
}

// Example_wrap 包装错误示例
func Example_wrap() {
	// 模拟数据库错误
	dbErr := sql.ErrNoRows

	// 包装错误添加上下文
	err := Wrap(dbErr, "DB_001", "failed to query user")

	fmt.Println(err.Error())
	fmt.Println("Underlying error:", err.Unwrap())

	// Output:
	// [DB_001] failed to query user: sql: no rows in result set
	// Underlying error: sql: no rows in result set
}

// Example_withDetails 添加详细信息示例
func Example_withDetails() {
	err := New("AUTH_001", "authentication failed").
		WithDetail("user_id", 12345).
		WithDetail("ip", "192.168.1.1").
		WithDetail("timestamp", "2025-12-08T10:30:00Z")

	fmt.Printf("Code: %s\n", err.Code())
	fmt.Printf("Details: %+v\n", err.Details())

	// Output:
	// Code: AUTH_001
	// Details: map[ip:192.168.1.1 timestamp:2025-12-08T10:30:00Z user_id:12345]
}

// Example_registerCode 注册错误码示例
func Example_registerCode() {
	// 注册自定义错误码
	RegisterCode("ORDER_001", "Order not found", 404)
	RegisterCode("ORDER_002", "Order already paid", 400)
	RegisterCode("ORDER_003", "Order expired", 410)

	// 使用注册的错误码
	err := NewWithCode("ORDER_001")

	fmt.Printf("Code: %s\n", err.Code())
	fmt.Printf("Message: %s\n", err.Message())
	fmt.Printf("HTTP Status: %d\n", err.HTTPStatus())

	// Output:
	// Code: ORDER_001
	// Message: Order not found
	// HTTP Status: 404
}

// Example_errorChain 错误链示例
func Example_errorChain() {
	// 构建错误链
	baseErr := fmt.Errorf("connection timeout")
	layer1 := Wrap(baseErr, "NET_001", "network error")
	layer2 := Wrap(layer1, "DB_001", "database connection failed")
	layer3 := Wrap(layer2, "API_001", "API request failed")

	// 检查错误链
	fmt.Printf("Top error: %s\n", layer3.Code())
	fmt.Printf("Is NET_001: %v\n", Is(layer3, layer1))

	// 展开错误链
	current := layer3
	for current != nil {
		if e, ok := current.(Error); ok {
			fmt.Printf("- [%s] %s\n", e.Code(), e.Message())
			current = e.Unwrap()
		} else {
			fmt.Printf("- %v\n", current)
			break
		}
	}

	// Output:
	// Top error: API_001
	// Is NET_001: true
	// - [API_001] API request failed
	// - [DB_001] database connection failed
	// - [NET_001] network error
	// - connection timeout
}

// Example_httpIntegration HTTP 集成示例
func Example_httpIntegration() {
	// 注册业务错误码
	RegisterCode("PAYMENT_001", "Payment failed", 402)

	handler := func(w http.ResponseWriter, r *http.Request) {
		// 业务逻辑
		err := processPayment()
		if err != nil {
			// 统一错误响应
			WriteHTTPError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	// 使用中间件处理 panic
	http.HandleFunc("/payment", HTTPMiddleware(handler))
}

func processPayment() error {
	// 模拟支付失败
	return NewWithCode("PAYMENT_001").
		WithDetail("transaction_id", "txn_12345").
		WithDetail("amount", 99.99)
}

// Example_predefinedErrors 预定义错误码示例
func Example_predefinedErrors() {
	// 使用预定义错误码
	notFoundErr := NewWithCode(ErrNotFound)
	invalidErr := NewWithCode(ErrInvalidArgument)
	authErr := NewWithCode(ErrUnauthenticated)

	fmt.Printf("Not Found: %s (HTTP %d)\n", notFoundErr.Message(), notFoundErr.HTTPStatus())
	fmt.Printf("Invalid: %s (HTTP %d)\n", invalidErr.Message(), invalidErr.HTTPStatus())
	fmt.Printf("Auth: %s (HTTP %d)\n", authErr.Message(), authErr.HTTPStatus())

	// Output:
	// Not Found: Resource not found (HTTP 404)
	// Invalid: Invalid argument (HTTP 400)
	// Auth: Unauthenticated (HTTP 401)
}

// Example_serviceLayer 服务层使用示例
func Example_serviceLayer() {
	// 模拟服务层
	userService := &UserService{}

	// 获取用户
	user, err := userService.GetUser(999)
	if err != nil {
		if Is(err, NewWithCode(ErrNotFound)) {
			fmt.Println("User not found, creating new user...")
		} else {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	fmt.Printf("User: %v\n", user)

	// Output:
	// User not found, creating new user...
}

type UserService struct{}

func (s *UserService) GetUser(id int) (*User, error) {
	// 模拟数据库查询
	if id == 999 {
		return nil, NewWithCode(ErrNotFound).
			WithDetail("user_id", id).
			WithDetail("table", "users")
	}

	return &User{ID: id, Name: "Alice"}, nil
}

type User struct {
	ID   int
	Name string
}

// Example_stackTrace 堆栈追踪示例
func Example_stackTrace() {
	err := createError()

	// 简单输出
	fmt.Printf("Simple: %v\n\n", err)

	// 详细输出(包含堆栈)
	fmt.Printf("Detailed: %+v\n", err)
}

func createError() Error {
	return New("STACK_001", "error with stack trace")
}

// Example_batchRegister 批量注册错误码示例
func Example_batchRegister() {
	// 批量注册错误码
	codes := []CodeConfig{
		{Code: "PRODUCT_001", Message: "Product not found", HTTPStatus: 404},
		{Code: "PRODUCT_002", Message: "Product out of stock", HTTPStatus: 409},
		{Code: "PRODUCT_003", Message: "Product price invalid", HTTPStatus: 400},
	}

	RegisterCodes(codes)

	// 使用注册的错误码
	err1 := NewWithCode("PRODUCT_001")
	err2 := NewWithCode("PRODUCT_002")

	fmt.Println(err1.Error())
	fmt.Println(err2.Error())

	// Output:
	// [PRODUCT_001] Product not found
	// [PRODUCT_002] Product out of stock
}

// Example_errorComparison 错误比较示例
func Example_errorComparison() {
	// 创建相同错误码的错误
	err1 := New("SAME_CODE", "first error")
	err2 := New("SAME_CODE", "second error")
	err3 := New("DIFF_CODE", "different error")

	// 相同错误码被认为是相同的错误类型
	fmt.Printf("err1 == err2: %v\n", Is(err1, err2))
	fmt.Printf("err1 == err3: %v\n", Is(err1, err3))

	// 类型转换
	var target Error
	if As(err1, &target) {
		fmt.Printf("Target code: %s\n", target.Code())
	}

	// Output:
	// err1 == err2: true
	// err1 == err3: false
	// Target code: SAME_CODE
}

// Example_complexScenario 复杂场景示例
func Example_complexScenario() {
	// 注册业务错误码
	RegisterCodes([]CodeConfig{
		{Code: "WITHDRAW_001", Message: "Insufficient balance", HTTPStatus: 400},
		{Code: "WITHDRAW_002", Message: "Account locked", HTTPStatus: 403},
		{Code: "WITHDRAW_003", Message: "Daily limit exceeded", HTTPStatus: 429},
	})

	// 模拟取款操作
	err := withdraw(1000.00)
	if err != nil {
		// 根据错误类型处理
		if e, ok := err.(Error); ok {
			fmt.Printf("Error Code: %s\n", e.Code())
			fmt.Printf("Message: %s\n", e.Message())
			fmt.Printf("HTTP Status: %d\n", e.HTTPStatus())

			// 获取详细信息
			if details := e.Details(); len(details) > 0 {
				fmt.Println("Details:")
				for k, v := range details {
					fmt.Printf("  %s: %v\n", k, v)
				}
			}
		}
		return
	}

	fmt.Println("Withdrawal successful")

	// Output:
	// Error Code: WITHDRAW_001
	// Message: Insufficient balance
	// HTTP Status: 400
	// Details:
	//   amount: 1000
	//   available: 500
	//   currency: USD
	//   user_id: 12345
}

func withdraw(amount float64) error {
	// 模拟余额不足
	balance := 500.00

	if amount > balance {
		return NewWithCode("WITHDRAW_001").
			WithDetails(map[string]interface{}{
				"user_id":   12345,
				"amount":    amount,
				"available": balance,
				"currency":  "USD",
			})
	}

	return nil
}
