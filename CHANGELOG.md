# Changelog

## v1.4.0

纯新增，无破坏性变更：补齐支付、账号与第三方登录、AI 大模型、文件、邮件五个业务模块，
并为通用模块与现有中间件模块补上缺失的码。错误码总数 70 → 136。

### 新增模块

- **模块 08 支付**（16 个）——支付单状态、渠道故障、回调验签、退款。
  金额不符 `ErrPaymentAmountMismatch`(400) 与验签失败 `ErrPaymentSignInvalid`(401) 单独建码：
  它们是安全事件，需要独立的日志级别与告警，不能并入通用参数错误。
  异步通知场景下 `HTTPCode()` 决定渠道是否继续重试（4xx 停止 / 5xx 重试）。
- **模块 09 账号与第三方登录**（8 个）——手机号与第三方账号绑定冲突、授权失败、短信验证码频控。
  通用登录态错误（未登录、登录过期、验证码错误等）继续用模块 00 已有的码。
- **模块 10 AI 大模型**（11 个）——与 `go-llm-provider` 的 `ErrorCode` 枚举一一对应。
  HTTP 码按「谁能解决」划分：用户可自行修正的给 200，我方配置/实现问题给 500，供应商侧故障给 503。
  用户配额继续复用通用的 `ErrQuotaExceeded`。
- **模块 11 文件**（8 个）——业务侧文件校验与存储故障。
- **模块 12 邮件**（6 个）——邮箱绑定、发送频控与服务故障。

### 通用模块（00）新增 7 个

`ErrRequestBodyTooLarge`(413)、`ErrUnsupportedMediaType`(415)、`ErrVersionConflict`(409)、
`ErrOperationInProgress`(200)、`ErrRiskControlRejected`(200)、`ErrFeatureDisabled`(200)、
`ErrSystemMaintenance`(503)。

其中 `ErrVersionConflict` 用于乐观锁更新影响 0 行；`ErrOperationInProgress` 用于分布式锁未抢到，
与表示「已处理过」的 `ErrDuplicateRequest` 语义不同。

### 中间件模块各补 2 个

- **01 Elasticsearch**：`ErrElasticsearchIndexNotFound`、`ErrElasticsearchBulkPartial`
  （bulk 部分失败是静默丢数据的常见来源）
- **02 MySQL**：`ErrMysqlDeadlock`（需单独统计并重试）、`ErrMysqlTransaction`
- **04 Redis**：`ErrRedisLockAcquireFailed`、`ErrRedisLockLost`
  （续期时发现锁已不属于自己，是数据一致性事故信号）
- **05 Kafka / 06 RabbitMQ**：`ErrXxxMessageDecode`、`ErrXxxRetryExhausted`
  （二者是进死信队列的判据，独立建码才能在死信里按码归类）

### 其他

- 新增码全部提供 `StatusXxx()`，与既有码保持一致，不需要调用方自己 `NewStatus`。
- `errcode.go` 的包注释补上已分配模块一览与两条选码原则；README 增加模块表、
  AI 大模型映射、支付回调重试控制、消息队列死信判据三节示例。

## v1.3.0

本版含**破坏性变更**，全部为 fail-closed 形态：要么编译期报错，要么把原先「失败被读成成功」的行为改成按失败处理。

### ⚠ 破坏性变更

- **`Wrap` / `Wrapf` / `WrapStatus` / `WithStack` 的返回类型由 `*Item` 改为 `error`。**
  原先 `err == nil` 时返回的是类型化 nil（`(*Item)(nil)`），装进 `error` 接口后接口值非 nil，
  导致 `return goerr.Wrap(err, "ctx")` 这种最自然的写法把**成功路径判成失败**，后续调用 `err.Error()` 还会 panic。
  改为返回接口后，nil 就是真 nil。
  *迁移*：需要业务码的地方改用 `goerr.AsItem(err)` 或新增的 `goerr.HasCode(err, code)`；
  直接 `return` 给 `error` 返回值的调用点不受影响。编译期即可发现全部需要改的位置。

- **`WrapStatusFrom` 的第一个返回值由 `*Item` 改为 `error`**，原因同上。

- **`status == nil` 的兜底方向由 `StatusOK()` 改为 `StatusInternalServer()`**，影响
  `New` / `Newf` / `NewFrom` / `NewfFrom` 与 `Item.StatusInfo()`。
  `goerr.New(err, nil)` 原先会产出业务码 `10010000`（成功），被统一响应层读成请求成功。
  状态缺失说明调用方漏传，按失败处理才安全。
  *迁移*：这是**静默的语义变化**，请检查是否有依赖「nil status 等于成功」的代码。

- **删除 `Err(string) error`。** 它只是 `errors.New` 的别名，却因为挂在 `goerr` 包下而具有误导性——
  返回的是不带业务状态的普通 `error`，会被统一响应层按内部错误兜底。请直接用 `errors.New`。

### 修复

- **`*Item` 的方法对 nil receiver 不再 panic。** `Error()` / `Message()` / `StatusInfo()` /
  `Code()` / `HTTPCode()` / `Unwrap()` / `Format()` 在 receiver 为 nil 时按「状态未知」兜底返回。
- **`%+v` 的调用栈改用 `runtime.CallersFrames` 解析。** 原先直接对 `runtime.Callers` 返回的
  返回地址调用 `FuncForPC` / `FileLine`，既不做 pc 调整也无法还原内联帧，会打印出根本不在调用链上的函数。
- **修复 `%+v` 输出按 Wrap 层数平方膨胀。** `Error()` 文本本身已含整条 cause 链，
  原先又递归打印每层 cause 的全文与各自的栈，n 层 Wrap 就重复 n 份。现在只打印一次文本 + 一份栈。
- **修复 `NewFn` / `NewFrom` / `NewfFrom` / `WrapStatusFrom` 的栈首帧落在 `goerr` 包内部**，
  而不是调用方代码。这几个入口原先经由 `New` / `WrapStatus` 转调，多出一层栈帧。
- **修复 `WithStack` 丢失外层包装上下文。** 对 `fmt.Errorf("ctx: %w", item)` 这类错误，
  原先返回从链上挖出的那个 `*Item`，`ctx:` 前缀被丢掉；现在原样返回 `err`。
- **修复手机号脱敏误伤长数字串。** `RedactForLog` / `SanitizeForClient` 原先会把
  `order_no=2024113812345678` 中的 11 位子串打码成 `[phone]`，反而妨碍排障。
  现在只有前面是非数字字符（或串首）的 11 位号码才算手机号，`+86` / `86` 前缀连同号码一起打码。
- **`Format` 补齐未知动词的处理**，`fmt.Sprintf("%d", item)` 从返回空串改为标准的 `%!d(...)`。

### 新增

- **`Convert(err error) *Item`** — 统一响应层一行拿齐 `Code` / `Message` / `HTTPCode`，永不返回 nil：
  `nil` → 成功状态；链上有 `*Item` → 原样返回；普通 `error` → 按 `StatusInternalServer()` 包装。
  取代原先「`AsItem` + `ok` 判断 + 兜底分支」三段式写法。
- **`HasCode(err error, code Code) bool`** — 按业务码判定错误链。同一失败语义在不同位置构造出的是
  不同 `*Item` 实例，`errors.Is` 判不出相等。

### 弃用

以下 API 打上 `// Deprecated:` 标记（编译不受影响，gopls / staticcheck 会提示替代写法）：

- `NewFn` — 与 `New(err, StatusXxx(), ...)` 只差一对括号。
- `ResolveStatus` / `NewFrom` / `NewfFrom` / `WrapStatusFrom` — 用 `any` 承载 `*Status` / `func() *Status`
  放弃了编译期类型检查；请直接传 `*Status` 调用 `New` / `Newf` / `WrapStatus`。
- `Is` — 请直接用标准库 `errors.Is`。

### 其他

- **错误码元信息收敛为单一真源。** 新增 `errcode.go` 的 `codeTable`，登记每个码的 HTTP 状态码与默认文案；
  `Code.Message()` 与 `status.go` 的预构建 `Status` 都从它取值，替换掉原先 146 行的 `switch` 与
  各自独立的 HTTP 码映射。新增错误码时漏配元信息会在包初始化期 panic，而不是留到运行时才发现。
- **一条错误链共享一份调用栈**（最内层、最接近错误起源的那份），不再每层 Wrap 各存一份。
- `callers` 只分配实际用到的帧数；`ClampMessage` 截断时不再分配整条 rune slice。
- `Wrap` / `Wrapf` 合并到同一个内部构造器，消除逐行重复的实现。

## v1.2.8

- `Wrap` / `Wrapf` / `WithStack` 包装不带 `goerr.Item` 的普通错误时改用 `StatusInternalServer()`，
  避免在统一响应层被误判为成功。
- 新增 `MustNewStatus`，用于必须校验编码规范的包级自定义状态定义。
- 补充 README 中关于 `New`、`WrapStatus`、`Wrap`、自定义状态与业务 HTTP 状态映射的说明。

## v1.2.7

- 新增订单模块（模块代号 07）错误码与默认文案。
