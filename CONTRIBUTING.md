# 参与贡献与发布检查清单

本地开发与检查使用 **Go 1.27.0 或更高版本**，与 go.mod 的最低版本一致；Go 1.26 工具链不满足当前未发布代码要求。

当前未发布开发版本为 **v1.5.0**，module/import 路径仍为 `github.com/gtkit/goerr`，不增加 logger 或 JSON 运行依赖。

## 代码变更的本地检查

以下命令用于后续获准的代码验证，不代表已执行或通过。纯文档、注释与说明修改只编辑指定内容并回读，不运行这些检查；原始命令、结果和遗留项见 [OPTIMIZATION_PLAN.md](./OPTIMIZATION_PLAN.md)。

```bash
go build ./...
go vet ./...
golangci-lint run ./...
go fix -diff ./...
go test -count=1 ./...
go test -race -count=1 -timeout=5m ./...
go test -bench=. -benchmem -count=3 ./...
go test -coverprofile=coverage.out ./...
go test . -run '^$' -fuzz '^FuzzClampMessage$' -fuzztime=30s
```

`go fix` 仅输出差异，目标为空，不自动应用。覆盖率目标为 ≥ 80%；退出码、覆盖率门槛与性能可比性分别判断。基准须保留原始样本与分配数据，不能把编译、模拟输出或单元测试当作真实日志平台接收证据。

## 契约与测试矩阵（现有用例与后续验收目标）

| 区域 | 内容 |
| --- | --- |
| 构造 | `New`、`NewFn`、`Newf`、`NewFrom`、`NewfFrom`、`NewStatus`、`MustNewStatus` |
| 包装 | `Wrap`、`Wrapf`、`WrapStatus`、`WrapStatusFrom`；`WithStack` 的缺栈补齐、有栈/typed-nil 身份、原 cause 与栈复用 |
| 解析 | `ResolveStatus`、`AsItem`、`Convert`、`HasCode`、`Unwrap` / `errors.Is`；保留首个匹配 Item 的主状态选择 |
| 展示 | `Message` vs `Error`、`Format` 正常格式、nil/零值、直接 State 写入失败后停止及外层 Fprintf 部分写入错误；业务 error/cause/栈保持 |
| 诊断 | `Inspect` 的主状态/根全文、Join/多 `%w` 实际原因树、虚拟 As 匹配、已有起源栈；256 节点/子槽与 16384 rune 边界，typed-nil、不可比较节点、循环、失败原子性与重复/并发读取所有权 |
| 安全 | `ClampMessage`、`SanitizeForClient`、`RedactForLog` |
| 错误码 | `ValidateCode`、`Code.Message`；136 个预定义 Status 的固定数字码/HTTP/默认文案、单例身份；弃用 `Is` 的 14 个兼容场景 |
| 外部调用方 | `test/integration_test.go` 的导出 API 冒烟、文本/JSON 普通函数切换与并存、显式投影/脱敏、写出前及部分写出后失败、诊断失败/无错误时跳过出口；16 个 worker 只读共享 Join 快照且各自独占 writer；业务错误和响应字段保持 |
| Example | `ExampleInspect`、`ExampleDiagnostic` 及既有截断/脱敏示例；跨包值使用带键初始化 |
| 性能 | 门禁修复续项使用共同夹具采集旧版/修后各 17 组、每组 20 次（500ms）；当前 `Inspect` 10 组各 10 次（100ms），共 780 样本。保留全部原始值、benchstat、中位数、分配及保守比值界，详见计划第 19 节 |
| 性质 | `ClampMessage` 的 15 个命名边界用例与 9 个 Fuzz 种子；保留真实有符号阈值及 int 极值，断言原字节/前缀、rune 边界、省略号、合法 UTF-8 和重复截断幂等性。阶段 6 的 30s、4 worker Fuzz 通过，实际执行 444,896 次 |

诊断入口仅 `Inspect`、`Diagnostic`、`Cause` 三个导出符号；内部快照不能直接作为客户端 DTO 或稳定 wire 协议。客户端响应继续使用原业务错误的 `Convert` 与 `SanitizeForClient`；日志出口由调用方普通函数承载，字段、敏感信息与已有 gtkit/logger/JSON 版本由调用方负责，不新增库内配置或日志接口。

门禁修复续项已通过 Go 1.27.2 下的 build、vet、lint（0 issues）、fix（无输出）、完整 race（330 条带 Test 字段的 pass 事件、两个包通过，含 6 个 Example）。根包语句覆盖率 **100.0%（384/384）**，183 个函数条目均为 100%；语句覆盖率不等于所有输入或路径证明。新增契约检查只固定现有行为，没有修改码表、导出 API 或依赖。

性能复核中，16 条正常旧路径的分配不增加，样本中位耗时最大回退 **3.80%**；9 组保守区间满足预算，7 组仍跨过 10% 门槛，**稳定耗时门槛未确认**。缺栈实际捕获成本单列；没有因采样不足而改预算或筛选数据。真实命令、环境、退出码、首批/补采结果、调用方边界和遗留项见 [计划第 19 节](./OPTIMIZATION_PLAN.md#19-剩余门禁修复续项本轮获准先方案后实施)。阶段 6 的 30s Fuzz 记录保留于性质行；本轮没有重跑持续 Fuzz 或真实日志/支付平台。

## 发版

只有获得明确发版授权后才执行；当前 v1.5.0 仍保留在 `[Unreleased]`。

1. 核对 CHANGELOG 版本号/日期、`Version`、module/import 路径与拟发布版本一致，再将未发布内容整理为版本区段。
2. 按授权创建 `vMAJOR.MINOR.PATCH` tag 与相应发布记录，不自动提交、推送或发布。
3. 破坏性变更先做兼容评估并给出迁移说明；首次升级为 v2 必须另有明确授权，不因新增功能或发版指令自动修改 module/import 主版本。
