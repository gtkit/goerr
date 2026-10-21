# goerr 优化方案与分阶段实施计划

日期：2026-10-21（Asia/Taipei）  
状态：阶段 4 的普通函数出口用例、日志接入边界与参考编译已完成，真实结果见第 16 节；阶段 5、6 未授权、未进入。既有 6 项 lint、旧 fix 建议、整体覆盖率 66.2% 和耗时可信度缺口保留，不报告全部门禁通过。  
当前授权和实施细节以第 16 节为准；第 13～15 节保留阶段 1～3 当时的状态与原始证据。  
需求来源：前序 oops 对照分析、方案收敛及阶段 1～3；本轮用户明确要求执行阶段 4、落实可插拔出口与日志接入边界并记录真实结果，禁止自动进入阶段 5。

## 1. 本次授权、规则与修订记录

- 阶段 0 历史授权仅修订本文件并回读；阶段 1 授权固定契约、核对环境、建立原实现基线、添加现有 API 的契约测试与补栈反例并记录真实结果；本轮仅执行阶段 4。
- 阶段 1 写入范围为本计划、`benchmark_test.go`、`error_test.go` 及真实错误记录；阶段 2 写入范围见第 14 节，阶段 3 写入范围见第 15 节，本轮阶段 4 写入范围见第 16 节。检查产物保存在独立临时目录，不修改规则或 OpenSpec，不提交、推送或发布。
- 已实际读取全局 `/Users/xiaozhaofu/.codex/AGENTS.md`、项目 `AGENTS.md`、`.harness/error-journal.md`，以及包结构、API 兼容性和文档 Guide。项目错误记录中的示例不当作本次真实事故。
- 阶段 0 修订原因：原方案的日志类型、接口层和 Go 版本目标与最新项目规则不一致，属性提供器、策略链和大量选项也超出必要设计；当时只修方案，尚无对应运行代码需要回退。本轮回滚定位见第 16.3 节；阶段 2、3 的历史回滚定位保留于第 14.5、15.3 节。
- 本轮按 harness-implement 执行获准的阶段 4，不启动 OpenSpec 或重新开展 research/plan 流程；阶段 1 已有门禁失败不构成扩大修复范围的授权。
- 保留原阶段编号 0～6 和实施先后顺序。阶段 0～3 已执行，阶段 4 本轮明确获准；阶段 5、6 未授权。阶段 1 的验收证据已具备，既有质量缺口不等于已解决，不借进入阶段 2 自动扩大修复范围。
- 后续执行前必须重新读取最新指令、适用规则、本方案与项目错误记录，确认获准阶段、对应步骤和验收标准。只完成获准范围，不自动进入下一阶段。

### 1.1 已确认现状与目标

| 项目 | 已确认事实 | 本方案目标 / 边界 |
| --- | --- | --- |
| 模块 | `github.com/gtkit/goerr` | 保持 module/import 路径与现有主版本，不首次发布 v2 |
| 当前语言声明 | 阶段 2 已将 go.mod 改为 `go 1.27.0`；阶段 1 原值为 1.26 | 已用现有 Go 1.27.2 验证；Go 1.26 不再满足最低声明，不宣称旧工具链兼容 |
| 当前版本常量 | 阶段 3 前为 v1.4.0；阶段 3 已核实当时远端最高 tag v1.4.0、Release 列表为空，并更新本地未发布 Version 为 v1.5.0 | 新增能力同步 MINOR 元数据，不改 module/import，不自动发布；历史证据见第 15 节，阶段 4 未重新查询远端或修改版本 |
| 运行环境 | 阶段 1 实际核实 `/usr/local/bin/go` 为 go1.27.2 darwin/amd64；详细环境和结果见第 13 节 | 基线与后续对比统一采用记录的同一环境，不能将已装工具链写成 go.mod 已升级 |
| 依赖 | 当前 go.mod 无外部依赖声明 | goerr 继续保持零外部运行依赖；不为日志出口引入日志库 |
| 日志规则 | 项目指定 `github.com/gtkit/logger`，沿用调用项目已锁定的路径和版本 | goerr 返回诊断数据，由终结错误的调用方按该约定记录；不注入 Logger，不设计日志 facade |
| 日志版本/API | goerr 未锁定 logger；阶段 4 已核实本机 logger v1.10.0 / zap v1.28.0 的包级 API 并离线编译参考函数；真实消费方版本与运行环境未确认 | 本地参考不成为依赖升级要求；其他路径/版本的真实接入须分别核实并遵守依赖审批，不自动迁移 major |
| 目录 | 现有根目录 `goerr` 包与 `test` 外部测试子包 | 不移动现有文件，不新增适配目录、子模块或依赖方向 |
| 写入前状态 | 阶段 1 写入前 `git status --short` 显示 `.gitignore` 已修改、本方案未跟踪，退出码 0 | 仅执行第 1 节授权范围；保留用户 .gitignore 与其他既有改动 |
| 源码基线 | 阶段 1 重新核实 HEAD 为 `d1f31a90c7080fffce42fc591ca4003a744a0eef`，运行代码无既有差异；最高本地 tag 为 v1.4.0 | 本地事实不冒充远端 tag/Release 核验；新增 API 前仍须核定发布版本 |

本方案没有服务、存储、连接池或分布式数据迁移；相关容量与一致性模型不适用。真实消费者、日志接收字段、接入人员、时间和绝对性能预算未确认，不编造。

## 2. 保留的问题与处理范围

下列结论来自前序源码分析；本轮重新读取相关源码，阶段 1 实测结果单列在第 13 节，其他待实现能力不视为已验证。

| 编号 | 当前事实 | 修订后的处理 |
| --- | --- | --- |
| F1 | HasCode 比较错误链中首个匹配的 Item，不能查询全部 Join 分支或被覆盖的内层码 | 保留旧语义并说明；诊断快照展示实际原因节点，不另增 FindCode 等查询 API |
| F2 | Convert(普通错误) 可生成无栈 Item，WithStack 遇到 Item 就返回 | 修复非 nil 无有效栈 Item 的补栈路径；不借修复新增栈配置入口 |
| F3 | 标准库外层包装不自动转发 Item 的 `%+v` 栈输出 | 提供最小诊断读取入口，同时取得原始错误全文和已有起源栈 |
| F4 | 已导出 Error 接口，旧转换与包装主要识别 Item | 不扩大旧转换契约；新快照可读取真实节点实现的 Error 接口 |
| F5 | 缺少独立结构化诊断数据出口 | 只增加中立的数据快照；业务属性与日志策略由调用方显式合并 |
| F6 | 现有栈已延迟解析，单链共享起源栈 | 保留；本期不增加全局或实例缓存框架 |
| F7 | ClampMessage 的 fuzz 仅检查调用，不断言性质 | 后续针对既有契约补充性质断言 |
| F8 | README/CHANGELOG 将支付回调 4xx 概括为停止重试 | 后续按具体渠道原文纠正，不在错误库新增 ACK 或重试策略 |
| F9 | README 使用已分配的 10010701 表示 Order expired | 后续改用对应 ErrOrderExpired，不改码表语义 |

## 3. 目标、非目标与方案选择

### 3.1 目标

1. 修复补栈缺口，保护已有 nil、状态、文案、错误身份和 cause 契约。
2. 新诊断入口读出原始全文、主状态、原因节点及已有栈，支持标准包装、Join 和多 `%w`。
3. 诊断数据与日志后端分离：调用方更换输出函数即可更换格式或目标，不修改 goerr。
4. 对齐 Go 1.27 和项目日志约定，以最小导出面积实现必要能力。
5. 保留完整兼容、失败、风险、回滚与阶段验收约束，所有实现状态继续标记待验证。

### 3.2 本期不做

- 不设计库内日志方法、日志接口层、Logger Option 或日志库适配模块；不直接依赖 oops。
- 不新增 WithAttrs、属性提供器、脱敏插件链、Render、Encoder/Redactor 接口及函数适配器、schema 框架、Builder 或全局注册中心。
- 不新增 FindCode、WithStackOptions、CallerSkip 等独立导出功能；已有业务码查询语义不变。
- 不新增公有 Options、可调资源限额、输出开关、环境变量或配置文件；不预建独立编码器文件。
- 不捕获请求/context/body，不自动生成追踪 ID，不读取源码片段，不启动后台任务、网络、重试或存储。
- 不改变现有码表、HTTP 映射、默认文案、Status 单例、自定义状态规则或旧错误文本快照语义。
- 不调整目录架构，不进行版本发布，不把 Go 1.27 目标写成已完成或继续兼容 Go 1.26 的承诺。

### 3.3 候选比较与推荐

| 方案 | 收益与代价 | 决策 |
| --- | --- | --- |
| 库内输出框架：快照、编码器、策略接口、提供器和 Options | 控制面多，但当前没有对应的多实现需求；导出符号、配置组合、依赖和失败路径过多 | 移除 |
| 最小快照出口：一个读取函数、必要值类型，调用方定义输出函数 | 能隔离日志依赖与编码格式，库内不需要插件生命周期和配置；调用方负责字段映射与脱敏 | 推荐 |
| 只保留字符串与 `%+v` | 无新增导出面，但仍无法稳定取得标准包装下的栈和原因信息 | 不满足本次结构化出口目标 |

这是方案中的明确推荐，不是代码实现批准。后续如必须增加本清单以外的公共能力，应先说明必要性并修订本文件，不能从“可插拔”推导出任意扩展点。

## 4. 可插拔出口与项目日志边界

```text
原始 error
    → Inspect：独立诊断快照
    → 调用方选择的普通输出函数
    → 调用方合并业务字段、执行已有脱敏、映射项目 logger 字段
    → 调用方终结错误时记录一次
```

- 可插拔边界是中立快照与调用方函数，不是 goerr 导出的插件接口。调用方可用普通 `func(Diagnostic) error` 组合格式或目标；该函数类型不必在 goerr 中命名、注册或管理。
- 多实现接口如有需要，定义在真正使用它的调用方；本库不提供只用一次的 Encoder、Logger、Handler、Redactor 或 Provider 抽象。
- goerr 不主动记录诊断读取失败，不替换原业务 error；调用方分别处理业务错误与诊断/输出错误。
- 库内没有需要日志副作用的流程。本期不为这一出口新增 logger 依赖；未来确有库内日志需求时，只能按项目锁定的 gtkit/logger 路径与版本直接使用包级 API，并先完成依赖核实与相应审批，不设计注入选项。
- operation、request_id、trace_id 等已有业务字段由调用方从自己已持有的数据显式加入日志，不为携带这些字段增加错误包装或 context 扫描。
- 输出函数只取得拥有独立数据的快照，不需要原始 error 引用。日志字段可以裁剪、重命名和脱敏；它不改变原错误的 Code、Message、HTTPCode 或 errors.Is/As。
- JSON 是调用方可选格式，goerr 不实现 JSON 编码器或 Item.MarshalJSON。调用方按其已锁定的 gtkit/json 版本处理；零依赖外部包互操作测试若使用 encoding/json，仅证明标准数据兼容，属于零依赖公共库例外，不据此改变项目生产 JSON 约定。
- 内部诊断快照可能包含敏感错误全文、路径和原因；客户端响应继续使用 Convert 与现有 SanitizeForClient，不把诊断数据直接当公网响应。

## 5. 兼容性与版本约束

| 对象 | 保留 / 变更约束 |
| --- | --- |
| Go 版本 | 阶段 2 go.mod 已声明 1.27.0；最低版本提高影响 Go 1.26 使用者，README/CONTRIBUTING/Unreleased 已说明，不宣称旧工具链兼容 |
| 模块与主版本 | module/import 路径保持；不自动创建 v2 分支、路径、tag 或 Release；如兼容评估认定需要首次 MAJOR 变更，停止相关实施并说明，不擅自执行 |
| 既有 API | 签名、返回类型、Error 方法集和弃用标记保持；不把字符串参数换为 Option，不为改 bug 新增符号 |
| HasCode / Convert | 保留首个 Item 与既有 nil/普通错误选择；不自动按最深原因、最大 HTTP 码或第三方 Error 改主状态 |
| WithStack | nil 仍是真 nil，有有效栈仍保留原 err 身份；非 nil 无栈 Item 修复允许返回新包装，须单列身份影响并记录 |
| cause 与文案 | 新包装保留原始 err，标准外层前缀不丢失；Error/Message 分工、默认状态及 errors.Is/As 可穿透性保持 |
| typed-nil | 不把非 nil 接口当成功；保护已知 Item 的 nil receiver 行为；新读取不直接调用未知 typed-nil 的方法 |
| 抓栈 | 保持默认捕获深度与起源栈共享；Inspect 只读已有栈，不抓自己的读取位置作为起源 |
| 序列化与日志 | 不给 Item 增加隐式编码/日志方法，不改已有日志调用的输出形状；新值类型不引入项目禁止的日志类型 |
| 新导出面积 | 仅为结构化诊断出口保留第 6.1 节清单；接口、字段和 GoDoc 同属公共契约，新增须有 Example、兼容评估与 MINOR 元数据同步 |
| 数据与目录 | 无存储迁移；不移动文件、包或依赖层级；不修改现有码表与状态实现 |

Go 1.27 使用按实际语义需要的现代写法，实施时读取对应现代 Go Guide；不为展示新特性添加 UUID、泛型抽象或无关替换。阶段 1 不升级工具链或修改依赖，只使用已安装的 Go 1.27.2 执行基线；go fix 仅使用只读 `-diff`。

## 6. 最小详细设计（实施结果与未验证项见第 13～15 节）

### 6.1 导出面预算

阶段 3 已新增函数仅 `Inspect(err error) (Diagnostic, error)`；必要数据类型仅 `Diagnostic` 与 `Cause`。名称及字段已按阶段 1 契约固定，若需额外函数、方法、接口、Option、错误类型或导出变量，先给出不可由现有 API/标准库替代的理由并修订计划。

| 值类型 | 必要字段及语义 |
| --- | --- |
| Diagnostic | HasError、Code、HTTPCode、Message、Error、Causes；分别表示输入是否非 nil、按旧 Convert 选择的主状态与公开文案、原始错误全文、原因节点列表 |
| Cause | Parent、HasCode、Code、Frames；Parent 是前序节点索引，根为 -1；HasCode 区分无业务状态节点与有效业务码；Frames 使用拥有独立切片的 runtime.Frame 值 |

- 两个类型只用于内部诊断数据传递，按 Guide 使用首位未导出的空字段守卫跨包位置初始化；不追加专用 getter、构造器、日志接口方法或 JSON 方法。
- 主状态与原始全文分开取得，不能把 Convert 找到的内层 Item.Error 当作原始全文而丢失标准库前缀。
- nil 输入产生成功状态与 HasError=false；非 nil 输入即使业务码为 ErrNo，也不能据码推导成 nil。未知 typed-nil 节点导致诊断读取失败，不把它写成成功。
- 成功时返回完整快照与 nil；失败时返回零快照与非 nil 诊断错误，调用方检查 error 后才能使用数据。取消部分成功、截断标志和允许继续输出的配置组合。
- 诊断错误只表达该读取操作失败，不伪装为原业务失败；本期不为每个内部失败原因增加导出错误类型或 sentinel，没有要求调用方按错误字符串分支。

### 6.2 补栈修复

- WithStack 区分链上存在 Item 与存在有效栈；覆盖 Convert(普通错误)、非 nil 零值 Item 及其标准库包装。
- 需要补栈时构造新包装，保留原始 err、全文、公开文案和状态；已有有效栈的输入继续原样返回，不重复采集。
- typed-nil、空状态与成功状态 Item 分别固定预期，不用归一化 nil 掩盖接口语义。
- 不新增 WithStackOptions、CallerSkip 或 StackDepth 配置；该 bug 修复只涉及内部构造与对应回归。

### 6.3 内部遍历与读取边界

- 私有遍历器仅服务 Inspect；前序深度优先，按原子节点顺序支持 Unwrap() error、Unwrap() []error、Join 与多 `%w`。
- 节点列表展示真实结构：自定义 As 返回的虚拟对象不伪装成原因节点。完整合法树的主状态仍遵守原 Convert 的选择，包括其既有自定义 As 行为。
- Cause 中可读实际实现 goerr.Error 的非 nil 节点码，但不让这些码替换旧主状态；需要对外翻译状态时仍显式调用 New/WrapStatus。
- 不能直接用 map[error] 比较全部 error；共享引用不等于循环。使用显式工作栈并按剩余预算逐个推进宽分支，不先复制整棵子树。
- 仅保留两个内部保护常量：最多访问 256 个节点、诊断文本最多 16384 个 rune；是保护边界，不是实测容量承诺，不导出为配置。超过边界返回诊断错误，不悄悄丢弃内容或回退原文。
- 循环或结构读取失败时，不再调用可能无界遍历同一树的 Convert；已有栈按分支关联，一条分支的共享栈不重复展示，不把不同分支拼成伪造栈。
- 方法调用要求外部 error 实现遵守有限、无阻塞的读取契约。访问预算不能隔离外部 Error/Unwrap/As 内部的阻塞、分配或副作用，不承诺运行时沙箱，也不新增通用回调恢复框架。
- 不保留原始 error、Status 指针、context、HTTP 对象或内部栈存储；每次快照拥有独立数据，调用方修改返回值不改变原错误或其他快照。

### 6.4 输出与性能边界

- 不提供库内 Render、编码器、策略配置或默认日志输出；调用方收到快照后自行选择普通输出函数。
- 脱敏复用调用方已有策略和本包现有字符串辅助；不在本库发明凭据键注册、插件顺序或保护字段回调机制。
- 需要两种输出验证可替换性时，用外部测试包中的两个小函数，不增加生产接口或专用实现文件。具体 gtkit/logger 接入仅在依赖/API 已核实后作为调用方说明，不引入本库运行依赖。
- 输出失败由调用方处理，不改业务 error，不自动重试或二次日志，也不因诊断失败回落输出未经处理的原文。
- 旧 New/Wrap/Convert 等不调用 Inspect，不增加 Item 字段、缓存或输出链开销。新路径成本由主动读取快照的调用承担。
- 后续比较旧路径 allocs/op、B/op 与耗时；拟定正常旧路径分配和字节数不增加，同一 Go 1.27 环境中位耗时回退不超过 10%。补栈修复确实增加捕获的场景单列，不能冒充零成本。
- 新读取测 1/4/16/64 层、分支树、节点预算边界、有栈/无栈和重复读取；不设没有实测基础的绝对微秒或 QPS 指标。本期不实施缓存，测量若显示新需求，只记录事实，不自动扩展任务。

## 7. 文件与依赖范围

| 范围 | 文件 | 后续用途 / 限制 |
| --- | --- | --- |
| 本轮计划 | OPTIMIZATION_PLAN.md | 阶段 1～3 的契约、执行状态、实际证据与遗留项 |
| 阶段 1 测试 | error_test.go、拟新增 benchmark_test.go | 原实现基准、旧行为契约与预期失败反例；不引用尚不存在的 API |
| 阶段 1 错误记录 | .harness/error-journal.md | 仅追加本轮真实失败及处置，遵守 AGENTS.md，不改规则或示例历史 |
| 阶段 2 执行范围 | error.go、error_test.go、go.mod | WithStack 补栈修复、错误链回归和 Go 1.27.0 声明；仅在第 14 节授权范围执行 |
| 阶段 3 已新增 | diagnostic.go、diagnostic_test.go | 读取入口、必要值类型、私有有界遍历及契约测试；未拆目录或扩展框架 |
| 阶段 3 / 后续测试 | example_test.go、benchmark_test.go；后续 test/integration_test.go、sanitize_test.go | 本轮只增加快照 Example/基准；两个出口组合与既有截断性质仍属于后续阶段 |
| 后续文档 | README.md、CHANGELOG.md、CONTRIBUTING.md | 阶段 2/3 随行为/API 同步必要说明，阶段 5 集中核对；不等到阶段 5 才补已有变更的必需文档 |
| 阶段 3 版本元数据 | version.go | 已按阶段 3 核实的远端 tag/Release 和原 Version 准备 v1.5.0 未发布 MINOR 元数据；未执行发版目标 |
| 保持不改 | Makefile、现有码表、Status 实现、其他规则及工作流文件 | 不借本方案变更发布流程、目录或业务状态 |

goerr 运行依赖仍以标准库为限；没有 logger、JSON 库的具体依赖版本选择。诊断读取使用现有 errors/fmt/runtime 和必要的标准库内部辅助，不让外部日志类型进入参数、字段或返回值。

阶段 2 已将 go.mod 提高到 1.27.0，未增加外部依赖，go.sum 未改变。阶段 3 已核对远端 tag/Release 与原 Version，按新增导出能力准备 v1.5.0 未发布元数据；未执行提交、tag 或 Release。

## 8. 风险、失败处理与回滚

| 风险 | 触发与影响 | 缓解 / 停止条件 |
| --- | --- | --- |
| Go 最低版本提高 | 旧工具链不能使用新语言声明或特性 | 单列 1.26→1.27 的兼容影响；实施前确认工具链和版本策略，不能谎称兼容旧工具链 |
| 旧业务分支漂移 | HasCode、Convert 或第三方状态选择被改变 | 保留旧语义，快照展示与主状态分离；未解释差异即停止 |
| 错误身份或文案丢失 | 补栈返回内层对象或断开 cause | 原始 err 为 cause；验证身份、全文、状态及 errors.Is/As |
| 返回数据共享 | 快照引用原栈切片或跨次共享结果 | 返回独立值数据；用修改快照和并发读取反例验收 |
| 内部信息进入公网/日志 | 调用方直接输出诊断原文 | 明确内部用途；调用方按项目日志约定脱敏，客户端继续既有响应路径 |
| 不完整读取被当成功 | 循环、超限、异常节点被静默跳过 | 失败返回零快照与 error；不提供忽略错误或部分输出开关 |
| 导出面重新膨胀 | 为测试、格式或一次性使用添加接口/Option | 只保留第 6.1 节清单；新增项须有独立必要性并先修订计划 |
| 错误依赖或日志用法 | 猜测 logger 版本、字段 API 或增加日志注入 | 库不直接记录；真实适配先核实消费方锁定版本，不绕过项目日志约定 |
| 老路径性能回退 | Item 增大或默认增加诊断处理 | 旧入口不走新流程；同环境基准超预算则停止 |
| 支付说明失实 | 通用 HTTP 映射替代渠道 ACK 协议 | 按具体原文修正文案，保留现有码表与业务边界 |

回滚步骤：

1. 阶段 0 仅修订本文件，当轮没有代码、依赖或存储回滚；阶段 1 测试与计划的撤回范围见第 13.5 节，其他用户改动保持原样。
2. 后续每阶段记录实际文件与可定位差异；未获提交授权时不自动提交，不能用强制重置覆盖用户或后续变更。
3. 出口接入问题：调用方停用自有快照输出函数，保留原日志与响应调用；不需注销全局插件或清理后台资源。
4. 快照实现问题：撤除调用方对新增函数/类型的引用后，才可回退到旧库；直接回退依赖会造成新调用无法编译，不能忽略这一前置条件。
5. 补栈修复问题：仅回退对应差异，重新核对旧契约；不连带撤回无关说明或用户改动。
6. Go 版本回退：先撤回依赖 Go 1.27 的实际代码与对应声明，再使用旧工具链核对。不能只把 go.mod 数字改回 1.26 就声称已恢复兼容。
7. 出现业务码漂移、cause 丢失、竞争或内部信息泄漏即停止推进；修复/回滚只在当时授权范围执行。外部消费者、生产接入和发布仍需各自明确授权。

## 9. 分阶段实施计划

保留阶段 0～6；阶段 0/1 已执行，阶段 2 本轮获准，阶段 3～6 待授权。每阶段结束登记真实结果并停止在授权边界。

### 阶段 0：方案修订与回读（历史已完成）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 0.1 | 用户只授权方案修订 | 读取最新全局/项目规则、错误记录、必要 Guide 和原方案 | 实际文件内容；现状与目标分开 | 规则不可读或需扩展文件范围 | 规则与修订依据；已完成读取 |
| 0.2 | 0.1 完成 | 仅改本文件，移除旧日志设计、对齐 Go 1.27、收敛 API/配置与阶段内容 | 本文件实际写入内容 | 需改源码、依赖、其他文档或错误记录 | 修订方案；已完成 |
| 0.3 | 0.2 保存成功 | 全文回读，核对所有章节与阶段无旧设计冲突 | 实际回读与文件状态 | 仍有不一致或把待实现写成已完成 | 回读记录；已完成，实施未开始 |

### 阶段 1：兼容契约与可复现基线（本轮执行结束，保留基线缺口）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 1.1 | 明确授权阶段 1；重新读规则、本方案和错误记录 | 固定最小导出函数/类型/字段、读取失败与节点语义；核对 Go 最低版本影响、真实版本与新增 MINOR 要求 | 第 13.1 节契约清单、版本依据及兼容影响；测试改动前保存并回读，退出 0 | 必须扩大 API、依赖、目录或首次 MAJOR 才能继续 | 已完成；新 API 仅固定设计，未实现 |
| 1.2 | 1.1 完成；Go 1.27 工具链入口已核实 | 原实现先跑既有测试，再仅增加旧入口基准；首轮采样先于 1.3，夹具修正补采另记 | 真实版本、命令、退出码、原始样本、GOOS/GOARCH/CPU | 环境缺失、已有失败未分清或样本不可比 | 原有检查、首轮和最终各 170 个样本已保存；取证完成，耗时可比性不足，10% 门槛未验证，原有门禁未全过 |
| 1.3 | 基线已保存 | 添加现有 API 可编译的补栈反例与旧状态/身份/文本契约用例，新入口先记录预期 | 预期失败与旧行为分别登记，不把失败写成通过 | 用例依赖不存在的 API 或预期无依据 | 已完成执行与记录；5 组兼容用例 race 通过，6 个补栈子用例预期失败，缺陷未修复 |

阶段验收：规则与最小导出面固定，版本影响透明，原始基线和补栈反例可复核；不修改运行实现，不进入阶段 2。

### 阶段 2：补栈修复与错误链兼容（本轮范围完成，保留门禁缺口）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 2.1 | 阶段 1 已验收且阶段 2 获准 | 按 Guide 将 go.mod 对齐 `go 1.27.0`，实施补栈修复；同步最低版本及修正行为的必要说明 | 反例转绿；旧 nil、状态、身份、首帧、栈共享与 cause 回归 | 旧正确行为漂移或新导出符号混入修复 | 核心修复与说明已完成；见第 14 节实际回归与限制 |
| 2.2 | 2.1 核对完成 | 核对补栈后的错误链、标准包装与 Join 兼容；只固定快照的遍历用例，不提前加入未被调用的私有实现 | 原始错误可穿透、外层文本和主状态保留；未来遍历反例已记录 | 错误身份、分支或状态改变，或预建一次性抽象 | 错误链回归已通过，后续用例仅记录预期；见第 14 节 |

阶段验收：补栈和错误链兼容有真实证据，旧业务判断不变，Go 1.27 影响已说明；记录结果后停止。

### 阶段 3：最小诊断快照（功能已完成，门禁遗留见第 15 节）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 3.1 | 阶段 2 已验收且阶段 3 获准 | 仅新增 Inspect、Diagnostic、Cause，在该功能内实现私有有界读取；固定原始全文、主状态与各分支已有栈 | 单链/Join/多 %w、不可比较节点、循环/超限；标准前缀与起源栈同时可读；旧 Convert 选择保持 | 未读完被当完整、抓诊断位置新栈、丢前缀或需要额外公共入口 | 最小快照及测试已完成；见第 15 节 |
| 3.2 | 3.1 完成 | 完成独立返回数据、失败原子性与外部包用法；同步 GoDoc、Example、README、CHANGELOG 和已核定 MINOR 元数据 | 篡改返回值/并发读取、错误分支、导出清单及版本一致性 | 返回数据共享、隐式编码或额外抽象进入生产 | 新增能力、独立数据、说明及 MINOR 元数据已核对；见第 15 节 |

阶段验收：仅必要导出符号，成功结果完整、失败有 error、旧错误不变；不实施属性注入或提供器，记录后停止。

### 阶段 4：可插拔出口与日志接入边界（功能完成，真实结果见第 16 节）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 4.1 | 阶段 3 已验收且阶段 4 获准 | 外部测试包用两个普通输出函数消费同一快照，证明格式/目标可换；不新增库内生产接口 | 两种输出各自结果、切换/并存、输出失败不影响原业务错误 | 需要日志类型、注册器、配置或回调框架进入 goerr | 已完成；专项和完整 race 通过，含部分写入失败边界；见第 16 节 |
| 4.2 | 4.1 完成；真实接入所需 logger 版本/API 已核实 | 按 gtkit/logger 包级使用约定写必要调用方接入说明，调用方负责字段与脱敏；无法核实则明确未验证 | 实际版本/API 依据与说明回读；替身用例只证明数据适配 | 未核实 API 被写成可运行示例，或需未经批准引入依赖 | 已完成必要说明和真实 SDK 参考编译；未连接日志平台，不宣称真实接入通过；见第 16 节 |

阶段验收：调用方函数可替换，库无日志依赖/隐式日志行为；真实平台未接入不得报接入通过，记录后停止。

### 阶段 5：说明、示例和契约同步（未授权，纯文档阶段）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 5.1 | 相关能力已验收且阶段 5 获准 | 集中核对 Go 1.27、最小 API、主状态/原因区别、内部快照与客户端响应、logger 与 JSON 边界 | 与实际实现和规则逐项回读 | 候选能力被写成已实现或需修改代码 | 相关说明；未开始 |
| 5.2 | 5.1 完成且渠道原文已核对 | 纠正支付重试说明与订单过期示例，不改码表和 HTTP 映射 | 原文与说明回读 | 未核实渠道假设或擅自改变业务状态 | 指定文档修正；未开始 |
| 5.3 | 5.1～5.2 完成 | 同步贡献矩阵和未发布记录，核对已实施阶段的 MINOR 元数据，不改历史版本号/日期 | 文档与版本记录回读 | 文档任务扩展为测试、发布或规则维护 | 文档交付；未开始 |

阶段验收：只修改指定说明并回读，不运行任何测试、构建或 OpenSpec。整体运行验收归阶段 6。

### 阶段 6：集成验收、性能评估与交付（未授权）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 6.1 | 阶段 1～5 已验收且阶段 6 获准 | 补齐外部包、并发和既有截断性质用例，按最新代码规则执行必要检查 | 第 10 节矩阵、实际 Go 1.27 环境与真实结果 | 兼容/数据归属/导出面/日志规则未达标 | 集成证据；未开始 |
| 6.2 | 功能结果完整且基准环境可比较 | 测原路径与快照读取成本，不自动实施缓存或扩展配置 | 原始样本、allocs/op、B/op、中位耗时和差异原因 | 超预算未解决或环境不具可比性 | 性能结果或明确未验证项；未开始 |
| 6.3 | 前述结果可复核 | 核对文件、导出清单、版本、限制、接入与撤回步骤，更新本计划状态 | 每项实际证据与回滚定位 | 存在阻塞问题或拟自动提交/发布 | 最终交付记录；未开始 |

阶段验收：必要条件均有真实证据，未验证范围明确；不包含消费者仓库修改、部署、提交、推送、tag、Release 或首次 v2。

## 10. 验收矩阵与未来检查

本矩阵是完整方案验收目标，阶段 1 的已执行范围与实际结果见第 13 节；尚未实现的能力继续标为未验证。

| 编号 | 场景 | 通过条件 |
| --- | --- | --- |
| M1 兼容 | nil、typed-nil、零值 Item、普通错误、自定义/空状态、旧与弃用 API | 旧签名、状态、公开文案与 cause 保持；仅明确的无栈修正列作差异 |
| M2 补栈 | Convert(普通错误)、无栈 Item、标准包装、有栈输入 | 缺栈有真实调用方帧；有效栈复用；原始前缀、身份约束和 errors.Is/As 保持 |
| M3 分支 | Join、多个 %w、覆盖状态、交换分支、第三方 Error | 主状态与旧 Convert 一致，原因节点结构和码可读；没有额外公有查询函数 |
| M4 边界 | 空节点、typed-nil、不可比较 error、重复引用、循环与预算边界 | 新读取有界且失败显式；零快照不能被当成功，不宣称可终止外部阻塞方法 |
| M5 所有权 | 返回切片修改、重复/并发读取、多个有栈分支 | 快照之间与原错误互不污染，无原对象引用、全局缓存和竞争 |
| M6 出口 | 两个调用方普通函数、替换/并存、输出返回错误 | 不改库即可更换出口，业务错误与响应不受输出失败影响，不自动日志或重试 |
| M7 规则 | 生产参数/字段/返回值、依赖、Options、日志接入说明 | 无项目禁止的日志类型或接口层；无 Logger 注入，无多余扩展点；logger API 有实际依据 |
| M8 信息边界 | 原始错误全文/路径与业务附加字段、客户端响应 | 内部诊断不进入默认公网响应；字段和脱敏留在调用方，不隐式改 Item 序列化 |
| M9 既有截断 | 空串、中文/emoji、零/负阈值、合法与非法 UTF-8 | 保持原省略号语义；合法 UTF-8 输入仍合法，非法输入不增加未经要求的转换规则 |
| M10 性能 | 同一 Go 1.27 下的旧路径与不同深度/分支快照 | 满足第 6.4 节预算，新增捕栈成本单列，不把复杂度分析冒充实测 |
| M11 版本与文档 | Go 1.27 声明、实际工具链、新导出清单、MINOR、支付说明 | 现状与目标区分；GoDoc/Example/README/CHANGELOG 对应；不复用码、不自动发版 |

后续代码阶段按届时已读规则和真实环境选用必要命令；下列是完整门禁候选，阶段 1 实际选择、命令参数与退出码单列第 13 节，不能由此推断全部执行：

```sh
go build ./...
go fix -diff ./...
go vet ./...
golangci-lint run ./...
go test -race -count=1 -timeout=5m ./...
go test -bench=. -benchmem -count=3 ./...
go test -coverprofile=coverage.out ./...
go test . -run '^$' -fuzz '^FuzzClampMessage$' -fuzztime=30s
```

实际执行须先核对 Go 1.27 和工具入口，go fix 的只读差异目标为空；不默认安装工具，不用有格式化或发布副作用的 Makefile 目标代替检查。测试、覆盖率、模拟输出和真实日志平台接入分别报告，不互相替代。

## 11. 本次记录与后续状态

| 日期 | 阶段 | 实际操作与证据 | 状态 / 未验证 |
| --- | --- | --- | --- |
| 2026-10-21 | 首版阶段 0 | 保存并回读过原方案；未实施代码 | 历史记录，不当作本轮已完成回读 |
| 2026-10-21 | 历史阶段 0 修订 | 已读最新规则、错误记录、必要 Guide、原方案、go.mod/version.go 和写入前状态；仅修订本文件；分段全文回读命令退出码 0，修订段落再次回读 | 当轮方案修订与回读已完成；当时阶段 1～6 未授权、未开始；当轮未运行代码检查 |
| 2026-10-21 | 阶段 1（含中断后明确继续） | 兼容契约落盘后增加测试；核实 Go 1.27.2；记录原始/最终命令与采样，实际失败先记录；真实结果见第 13 节 | 执行与取证结束；保留 lint/fix/覆盖率和耗时可比性缺口，6 个补栈反例仍失败；运行代码未改，阶段 2～6 未授权 |
| 2026-10-21 | 阶段 2 | 仅修 WithStack 与最低 Go 声明，补齐回归/说明；完整 race/build/vet 通过，17×3 基准及真实遗留门禁见第 14 节 | 本阶段范围完成并停止；fix/lint/覆盖率与耗时限制未全部通过，阶段 3～6 未授权 |
| 2026-10-21 | 阶段 3 | 按既定三个导出符号实现独立、有界快照，补契约/Example/基准与说明，核定 v1.5.0 未发布元数据；全量 race/build/vet、81 个基准原始样本见第 15 节 | 本阶段功能完成并停止；整体 lint/fix/覆盖率仍未达标，耗时门槛未验证，阶段 4～6 未授权 |
| 2026-10-21 | 阶段 4 | 外部包两个普通函数消费同一快照，验证切换/并存/脱敏及写入失败；README 明确 gtkit/logger 边界，真实缓存 SDK 离线编译；命令、结果和范围核对见第 16 节 | 本阶段功能完成并停止；既有门禁缺口保留，真实平台未接入；阶段 5、6 未授权、未进入 |

阶段 0 未写错误记录或其他规则。本轮 `.gitignore` 的既有改动保持原样；阶段 1 仅按 AGENTS.md 追加真实失败记录，不改规则。各阶段记录授权范围、实际基线、文件与导出差异、命令/环境/退出码、逐项结果、失败与未验证项、回滚定位和停止位置。

## 12. 依据与未确认事项

- 当前约束依据：根目录 AGENTS.md、全局 AGENTS.md，以及已读取的 pkg-structure、pkg-api-compat、pkg-docs Guide；文档任务不触发代码门禁或工作流。
- 前序 oops 核心源码分析仅作设计参考；其 main 未固定 commit，不据此承诺发布版本能力，也不继承其全部接口。
- [Go errors 文档](https://pkg.go.dev/errors)用于错误树背景；前序已读的 Go 1.26 fmt 源码只解释当前现象，不作为 Go 1.27 已验证证据。后续代码阶段以实际 Go 1.27 源码/工具链重新核对相关行为。
- [微信支付回调通知注意事项](https://pay.wechatpay.cn/doc/v3/merchant/4012075420)为前序已核对的渠道依据；说明修正时复核具体原文，不外推其他渠道。
- 未确认：实际消费者与旧工具链需求、消费方 gtkit/logger/JSON 的锁定路径版本与平台字段、接入人员和时间、性能绝对预算，以及第 13 节单列的环境限制。阶段 4 的 logger v1.10.0 / zap v1.28.0 仅为已核实并编译的本地 API 参考，不能据此宣称其他消费者已兼容、平台已接入或供应链验证已完成。
- 阶段 1 未修改运行代码；阶段 2 仅执行补栈、最低 Go 声明和对应回归/说明，不升级已装工具链或增加运行依赖；阶段 3 已实现与验证最小快照，真实结果与未通过项见第 15 节；阶段 4 已完成普通出口和日志边界说明，未改生产源码/API/依赖；阶段 5、6 未开始。

## 13. 阶段 1 执行记录（2026-10-21）

### 13.1 步骤 1.1：固定兼容契约

本节在任何测试改动前保存并回读，命令退出 0；下列新入口为阶段 3 的设计契约，尚未实现。仅保留三个新增符号及列明字段，不再追加字段、错误类型、选项或库内出口接口：

```go
func Inspect(err error) (Diagnostic, error)

type Diagnostic struct {
    _        struct{} // 强制跨包带键初始化；不是返回数据字段
    HasError bool
    Code     Code
    HTTPCode int
    Message  string
    Error    string
    Causes   []Cause
}

type Cause struct {
    _       struct{} // 强制跨包带键初始化；不是返回数据字段
    Parent  int
    HasCode bool
    Code    Code
    Frames  []runtime.Frame
}
```

| 契约 | 固定语义与后续验收 |
| --- | --- |
| 主状态 | 对完整、读取方法有限且无阻塞的错误树，Code/HTTPCode/Message 等于旧 `Convert(err)` 的对应值，包括自定义 As 提供的 Item；不改旧 Convert/HasCode 的首个 Item 选择 |
| 原文 | Error 保存原始根 `err.Error()` 全文，保留标准包装前缀；不保存每个原因的重复全文，不由主状态 Item 的 Error 替代根全文 |
| 真 nil | `Inspect(nil)` 成功：HasError=false，Code=ErrNo，HTTPCode=200，Message=StatusOK().Msg()，Error=""，Causes=nil；区别于全部字段为零的失败快照 |
| 非 nil 成功码 | `Convert(nil)` 返回的成功 Item 作为非 nil 输入时，HasError=true，仍保留成功码与公开文案；不能据 ErrNo 推断无错误输入 |
| typed-nil | 已知 `(*Item)(nil)` 可安全表示真实节点：HasError=true，根原文为空，主状态遵守 Convert 的内部错误兜底，节点 HasCode=false、Code=0、Frames=nil；未知 typed-nil 节点在调用其方法前使读取失败，不统一归一化为 nil |
| 原因节点 | 前序深度优先，按实际 Unwrap 子项顺序；根 Parent=-1，其他 Parent 指向较早的直接父索引；普通节点也入列表，HasCode=false 时 Code=0。只对实际、非 nil、实现现有 goerr.Error 的节点读 Code，不把自定义 As 的虚拟 Item 加入原因树 |
| 空子项与共享引用 | 真 nil 子项不产生 Cause；同一对象出现在不同分支时分别列出，不能用全局去重将共享引用当循环。当前路径的可比较重复节点判为循环；不可比较节点不得因比较而 panic，无法识别的循环仍由预算失败退出 |
| 预算 | 最多 256 次节点/子槽访问（根占 1，子切片中的 nil 槽也消耗预算，避免宽 nil 分支绕过边界）；最多 16384 个 rune，累计返回的 Error、Message 及实际保存帧的 File/Function 文本。恰好达到上限允许，超过即失败。UTF-8 无效字节按 Go range/rune 计数但原字符串不改写；不提前复制整个子切片 |
| 栈 | Frames 为已有 PC 通过 runtime.CallersFrames 解析后的独立值切片，不采集 Inspect 位置；一条路径上同源共享栈归属于首次携带它的实际节点，后代不重复保存；兄弟分支分别关联。空栈为 nil；runtime.Frame.Func 是运行时只读元数据，不是原错误或可写栈存储 |
| 所有权 | Causes、每个 Frames 在不同快照间互不共享可写存储；不保存 error、Status、context 或原 PC 切片，不增加全局缓存；修改返回数据不改变原对象与其他快照 |
| 读取失败 | 循环、超预算或未知 typed-nil 返回 `Diagnostic{}` 与非 nil 诊断错误，不能返回部分数据；诊断错误只供调用方处理本次读取失败，不替代原业务 error，不新增可按类型/字符串分支的公共契约。外部方法须有限、无阻塞且不 panic，本期不提供隔离或通用恢复框架 |
| 可插拔出口 | 调用方用普通 `func(Diagnostic) error` 自行输出、合并字段与脱敏；不命名/导出该函数类型，不增加 Logger、注册器、Provider、Option 或日志依赖 |

既有 API 的阶段 1/2 兼容约束：

1. Error 接口方法集、Item/Status 结构、所有签名与弃用 API 保留；标准包装/Join 的主状态仍是首个匹配 Item。HasCode 不是全树查询，不能因后续分支含目标码而变成 true；同码不同 Item 不因码相同而被 errors.Is 视为同一对象。
2. New/Wrap 的公开文案与完整文本、nil status 的内部错误兜底、自定义 Status 的宽松构造、内建 Status 单例与码表均保持；Wrap/Wrapf/WrapStatus/WithStack(nil) 返回真 nil，Convert(nil) 共享成功实例且不抓栈。
3. WithStack 沿用首个 Item 的选择。首个 Item 非 nil 且 len(stack)>0 时返回原输入，保留标准外层身份；没有 Item 的普通错误按现有方式新包装。首个 Item 非 nil 且缺栈时，阶段 2 才允许返回新包装：Error 等于原输入全文，Code/HTTPCode/Message 等于原状态，直接 cause 为原输入，errors.Is 可找到原输入与根因，首帧为调用方，二次 WithStack 返回同一新包装。
4. 零值 Item 与 Convert(nil) 的非 nil 成功 Item 同属缺栈输入，修复时不得修改原对象/成功单例。已知 typed-nil Item 保留 WithStack 原样返回的既有身份与 nil receiver 行为，不借补栈改变其语义；第一匹配为 typed-nil 时不改为选择后续 Item。
5. 当前 HEAD、Version 和最高本地 tag 已核实为上述基线/v1.4.0；未查询远端 tag/Release，不能把本地结果称为远端最新版。仅新增诊断出口属于拟定 MINOR，最终号须在阶段 3 实施前核对真实远端版本再定；本轮不改元数据，不沿用历史 CHANGELOG 中的破坏性变更作为再次破坏兼容的授权。
6. 安装环境 Go 1.27.2 与 go.mod 的 `go 1.26` 分开记录。阶段 2 将最低声明提高为 `go 1.27.0` 会排除 Go 1.26 工具链，实际消费者需求尚未确认；不得报告兼容旧工具链。若消费者约束要求首次 MAJOR，先停止相关变更，不自动创建 v2。

### 13.2 步骤 1.2：环境与基线执行安排

- 已核实工具入口：`/usr/local/bin/go`（GOROOT=/usr/local/go，go1.27.2，darwin/amd64，GOAMD64=v1，CGO_ENABLED=1）；检查显式使用 GOTOOLCHAIN=local，避免隐式工具链下载；不修改用户配置。
- 已有 lint：`/Users/xiaozhaofu/go/bin/golangci-lint`，版本 2.14.0、由 go1.27.2 构建；根目录没有 `.golangci*` 配置，按现有工具调用、不指定新配置。benchstat 为 `/Users/xiaozhaofu/go/bin/benchstat`，实际构建 Go 1.27.1、模块 `golang.org/x/perf v0.0.0-20260908200009-22c9c6c9d4da`，未安装或升级工具。
- GOMAXPROCS/GOGC/GOMEMLIMIT/GOEXPERIMENT/GODEBUG 环境覆盖均未设置，未发现 PGO 文件；供电、后台负载和沙箱外硬件信息未确认。首轮基准实际输出 CPU 为 Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz，后缀 -12 表示该次基准 GOMAXPROCS=12，不据此冒充 sysctl 查询成功。
- 原始证据目录：`/private/tmp/goerr-stage1.E7FRRL`。先保存运行源码与既有测试指纹、原计划副本、工作区差异；再跑原有 build、只读 fix、vet、lint、race、覆盖率，逐条保留输出和退出码。失败先记录、分清工具限制/已有缺陷，不修改运行代码或全局工具。
- 仅增加 benchmark_test.go，使用 b.Loop、ReportAllocs；分别测 New、Wrap、Wrapf、WrapStatus、Convert（nil/普通/已有 Item）、HasCode、WithStack（普通/有栈/Convert 缺栈），以及格式化和深度 1/4/16/64 的旧链读取。补栈缺失场景单列，避免后续修复新增捕获成本混入正常旧路径；基准 setup 不计时。
- 原实现无现有 Benchmark。本轮补齐基准夹具后采用 `go test -run '^$' -bench . -benchmem -count=10 ./...`，保留全部原始样本和 benchstat 单组统计，不进行无改动的虚假前后性能比较。深度基准仅测现有 Convert，不提前实现 Inspect。
- 步骤 1.3 在基线已保存后仅追加 error_test.go：固定首个 Item、Join/多 %w 顺序、自定义 As、身份/文案/状态和 typed-nil；单独运行无栈补栈反例并如实记录预期失败。不 skip、不弱化断言，不据反例自动修复实现。
- 上述为执行安排，真实命令、退出码、样本及结果在后文区分已执行与待执行。

首轮基准 `baseline-bench.txt/json`：退出 0，共 17 组各 10 个样本，根包耗时 223.578s；`baseline-benchstat.txt/json` 退出 0，部分耗时区间波动达 ±35%，不能用采样次数替代稳定性判断。随后 lint 检出新增夹具裸调用未显式接收结果（ERR-20261021-005）；仅将基准调用改为显式接收/丢弃被测值，`corrected-benchmark-lint` 再次退出 1，但已仅剩原始 6 项，无基准文件问题。首轮夹具副本和样本保留，不冒充修正后样本；最终补采见第 13.4 节。兼容用例及补栈反例在首轮基线保存后加入，中断时尚未执行；用户明确要求继续后，已重新读取最新规则、计划和错误记录，恢复后的实际结果见第 13.3 节。

步骤 1.3 在编码前保存的用例清单（现已执行，结果见第 13.3 节；原有 nil、弃用入口、首帧、栈共享等用例保留，不重复重写）：

| 测试 | 固定输入与预期 |
| --- | --- |
| TestCompatibilityFirstItem | 单包装、Join 正反顺序、多 %w（含索引格式）、覆盖状态；Convert/HasCode 只选择既有首个 Item；errors.Is 仍能到达其他分支，不能把首码选择误写成全树查询 |
| TestCompatibilityCustomAs | 自定义 As 返回真实树之外的 Item，保留主状态/身份选择；未来 Inspect 不为虚拟 Item 造节点，仅在本计划固定预期 |
| TestCompatibilityForeignError | 实现现有 Error 但不是 Item 的第三方错误，Convert/Wrap 继续内部错误兜底；新出口读取原因码不得反改这一旧契约 |
| TestCompatibilityTypedNil | 直接 typed-nil 和其位于 Join 第一分支时，AsItem 为 nil/true，Convert 为非 nil 内部错误兜底，HasCode 按首匹配处理，WithStack 原样返回 |
| TestCompatibilityTextSnapshot | 构造后改变底层测试错误的文本，New/Wrap 已缓存全文不变、cause 身份保留；文本断言只用于明确的公开输出兼容契约 |
| TestWithStackMissingStack | 六个缺栈输入：Convert(EOF)、零值 Item、Convert(nil)、标准包装的 Convert(EOF)、空状态 Item、自定义宽松状态 Item；要求新包装保留状态/全文/cause、调用方首帧、原对象无修改及二次调用幂等。现实现预期失败，阶段 2 才转绿 |

已执行的原始基线（全部先于本轮测试修改）：

| 命令（可执行文件见上文） | 实际退出码 | 结果与证据文件（证据目录内） |
| --- | --- | --- |
| `go build ./...` | 0 | 无输出；original-build.txt/json |
| `go fix -diff ./...` | 1 | 只读差异：error_test.go:289 的 strings.Split 循环建议改 SplitSeq；未应用，original-fix.txt/json |
| `go vet ./...` | 0 | 无输出；original-vet.txt/json |
| `golangci-lint run ./...` | 1 | 既有 6 项：error.go:84/90/101/103/105 的 fmt 写入 errcheck；resolve_test.go:24 的 ineffassign；original-lint.txt/json |
| `go test -race -count=1 -timeout=5m ./...` | 0 | 根包与 test 子包均通过，无 race 报告；original-race.txt/json |
| `go test -count=1 -coverprofile=/private/tmp/goerr-stage1.E7FRRL/original-coverage.out ./...` | 0 | 根包 55.4%，test 子包无可统计语句；命令成功但未达到 Guide 的 80% 门槛；original-coverage.txt/json/out |
| `go tool cover -func=/private/tmp/goerr-stage1.E7FRRL/original-coverage.out` | 0 | 总计 55.4%，大量 StatusXxx getter 和弃用 Is 为 0%，Message 为 80%；original-cover-functions.txt/json；不据此扩展全仓补测 |

原始基线指纹 `before-sha256.json` 覆盖运行文件、既有测试、go.mod/go.sum 和用户 .gitignore，共 16 个文件；添加 benchmark 前复核 `changed_files=[]`（before-fixtures-comparison.json）。原始差异只有用户 .gitignore 和未跟踪计划，未包含运行代码修改。

检查统一由临时 `run.py` 逐条采集原始输出及独立退出码，每条 `.json` 包含命令、cwd、UTC 起止时间与环境覆盖：GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off、GOCACHE 与 GOLANGCI_LINT_CACHE 均指向证据目录子目录；这些是本次进程环境，不是项目配置变更。GOWORK 与 GOEXPERIMENT 为空，未启用额外构建标签或编译参数。原有 fix/lint/覆盖率缺口已单列，不阻止继续采集基线，但不能宣称原有质量门禁通过。

本轮环境/命令问题记录：ERR-20261021-001 为 sysctl 权限失败（合并 shell 未单独捕获其退出码）；ERR-20261021-002 为未经清单确认读取不存在的 .golangci.yml（退出 1）。未改权限、安装工具或创建替代配置。ERR-20261021-003/004 分别记录真实基线门禁与覆盖率缺口；这些问题不在阶段 1 自动修复。

### 13.3 步骤 1.3：兼容用例与补栈反例的真实结果

全部使用同一 Go 1.27.2 与第 13.2 节进程环境；未调用或实现 Inspect。

| 命令 | 实际退出码 | 真实结果 / 证据文件 |
| --- | --- | --- |
| `go version`（恢复后） | 0 | go1.27.2 darwin/amd64；resumed-go-version.txt/json；与中断前一致 |
| `go test -race -count=1 -timeout=5m -run '^TestCompatibility' ./...` | 0 | 新增 5 个顶层兼容测试通过；test 子包没有匹配用例，不能称该子包在此命令中已测试；compatibility-race.txt/json |
| `go test -race -count=1 -timeout=5m -run '^TestWithStackMissingStack$' .` | 1 | 六个子用例均复现缺栈与未新建 cause 包装；无 race 报告；missing-stack-race.txt/json |
| `go test -json -race -count=1 -timeout=5m ./...` | 1 | 完整执行：59 个 pass 事件、7 个 fail 事件（6 个反例子用例及父测试），没有 skip 或 race 报告；根包失败，test 子包通过；final-full-race.txt/json、final-full-race-summary.json |
| `go build ./...`（测试修改后） | 0 | 无输出；final-build.txt/json |
| `go vet ./...`（测试修改后） | 0 | 无输出；final-vet.txt/json |
| `golangci-lint run ./...`（测试修改后） | 1 | 仅原有 6 项，与原始基线一致，无本轮新增用例或基准问题；final-lint.txt/json |
| `go fix -diff ./...`（测试修改后） | 1 | 仅原有 stackFrames 的 Split→SplitSeq 建议；因新增 import 行号移至 290，没有新增现代化差异；final-fix.txt/json；未应用 |

上述 pass/fail 计数由真实 `go test -json` 的测试事件汇总，包含父测试与子用例，不是 66 个独立顶层测试。不因没有 race 报告而声称失败用例通过。

| 兼容矩阵 | 已有/新增证据 | 本阶段结论 |
| --- | --- | --- |
| nil / typed-nil / 空状态 / 弃用入口 | 保留既有 nil、ResolveStatus、From 系列与首帧测试；新增 typed-nil 首匹配分支反例 | 既有测试及新增兼容测试通过；未改变 nil 语义 |
| 主状态与身份 | FirstItem 六个子用例：标准包装、Join 正反顺序、多 %w 与索引参数顺序、显式状态覆盖；同码不同对象不匹配 errors.Is | 通过；HasCode 保留首匹配语义，其他分支仍可被 errors.Is 穿透 |
| 自定义 As 与非 Item Error | CustomAs、ForeignError | 通过；虚拟 Item 不等于 Unwrap 子节点，第三方业务码不替代旧主状态 |
| 文本、公开文案与 cause | 既有 New/Wrap/Convert 测试；新增 TextSnapshot 修改自有测试根因文本 | 通过；已缓存全文不重新序列化，公开文案与原始 cause 身份保留 |
| 有效栈 | 既有有栈身份、首帧、栈共享与新增首 Item 表格 | 通过；原输入身份与原始前缀保留 |
| 缺失栈 | converted_plain、zero_item、converted_nil、wrapped_plain、nil_status、custom_status | 六例均失败；当前 WithStack 仍无栈，身份/首帧的新行为须待阶段 2 修复后验证 |
| 新快照、预算与出口 | 仅第 13.1 节固定设计 | 未实现、未运行；不能宣称有界读取、快照所有权或出口集成已通过 |

ERR-20261021-006 保留为开放的预期缺陷记录；001/002 的环境查询与路径处理问题、005 的基准夹具问题已按真实处置关闭，003/004 原有门禁与覆盖率缺口仍开放。没有通过修改忽略项、跳过测试、增加构建标签或降低断言取得全绿。

### 13.4 最终基准补采：实际结果与原始样本

命令：`go test -run '^$' -bench . -benchmem -count=10 ./...`，退出码 **0**；统计命令 `benchstat /private/tmp/goerr-stage1.E7FRRL/final-bench.txt` 退出码 **0**。实际起止时间（Asia/Taipei）为 `2026-10-21T16:12:34+08:00` 至 `2026-10-21T16:16:11+08:00`，UTC 原始时间保存在对应 json 中；原始输出与独立命令记录为证据目录下 final-bench.txt/json、final-benchstat.txt/json。

共 17 组、170 个样本，均来自未修改的运行实现；CPU 为 Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz，GOOS/GOARCH 为 darwin/amd64，GOMAXPROCS=12。基准使用 b.Loop、ReportAllocs，另保留 RunParallel 读取基准；不带 race，setup 不计入顺序基准的循环时间。`-run '^$'` 只采集基准，其 PASS 不能代替第 13.3 节失败的完整测试集。

最终夹具为 benchmark_test.go，证据副本 final-benchmark_test.go；真实 SHA-256 为 `63a9170c87471187908ec9e95fe7f6a6372253c3d22bb19d11ad5bb832eb3319`。该值已与实际文件、保存副本、final-source-sha256.json 和计划回读一致；此前无依据候选值的文档错误已在 ERR-20261021-007 记录、更正并关闭。采样期间基准与 error_test.go 指纹均未变化。

下表耗时为固定版本 benchstat 对打印样本的单组中位数与其报告区间（n 表示 ns/op，µ 表示 µs/op）；B/op、allocs/op 在每组 10 个样本中均为表列常数，原始迭代次数与样本在下方完整保留。final-samples.json 与 final-summary.json 由真实输出解析，非合成数据。

| 基准（后缀为 GOMAXPROCS） | benchstat 耗时 / 报告区间 | B/op | allocs/op | 样本数 |
| --- | --- | --- | --- | --- |
| New-12 | 1.231µ ± 61% | 136 | 3 | 10 |
| Wrap-12 | 166.7n ± 41% | 104 | 2 | 10 |
| Wrapf-12 | 258.8n ± 33% | 120 | 3 | 10 |
| WrapStatus-12 | 99.80n ± 11% | 128 | 2 | 10 |
| ConvertNil-12 | 2.440n ± 30% | 0 | 0 | 10 |
| ConvertPlain-12 | 48.56n ± 30% | 80 | 1 | 10 |
| ConvertItem-12 | 10.70n ± 43% | 0 | 0 | 10 |
| ConvertDepth/1-12 | 12.61n ± 17% | 0 | 0 | 10 |
| ConvertDepth/4-12 | 29.11n ± 9% | 0 | 0 | 10 |
| ConvertDepth/16-12 | 111.8n ± 14% | 0 | 0 | 10 |
| ConvertDepth/64-12 | 386.4n ± 10% | 0 | 0 | 10 |
| ConvertParallel-12 | 2.889n ± 16% | 0 | 0 | 10 |
| HasCode-12 | 14.65n ± 33% | 0 | 0 | 10 |
| WithStackPlain-12 | 539.9n ± 17% | 112 | 2 | 10 |
| WithStackExisting-12 | 13.38n ± 7% | 0 | 0 | 10 |
| WithStackMissing-12 | 6.951n ± 6% | 0 | 0 | 10 |
| FormatVerbose-12 | 1.849µ ± 6% | 744 | 12 | 10 |

**可比性限制：**最终耗时区间最大为 ±61%（New），多组超过方案的 10% 回退预算；供电和后台负载未控制/未确认，不能将波动归因于某个实现或宣称耗时性能达标。首轮与补采也不能当作优化前后两版：运行实现没有改变、夹具有修正。本轮已完成取证，但未建立足以裁定 10% 耗时变化的稳定比较条件；按停止条件保留未验证，不继续追加性能优化、切换环境或自动进入下一阶段。现有分配/字节样本可作为对应夹具的参考基线，后续任何耗时门槛仍须在可比较条件下重新核实。

WithStackMissing 当前为零分配且没有实际捕获，这是已复现缺陷的成本，不能当成未来修复必须保持零成本的要求。ConvertDepth 的 1/4/16/64 为包在既有 Item 外的标准包装层数，不代表新 Inspect 的实现或性能。并行耗时为 RunParallel 的平均每操作成本，不等于单请求延迟或生产吞吐量。

<details>
<summary>最终基准原始输出（17 组 × 10 次；包含 CPU、迭代次数、耗时与分配）</summary>

```text
goos: darwin
goarch: amd64
pkg: github.com/gtkit/goerr
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkNew-12                  	  566342	      1985 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	  569516	      2204 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	  680902	      1733 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	  970190	      1228 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	  979053	      1239 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1216447	      1030 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1223614	       996.0 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1300850	       907.2 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1353384	       983.4 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	  846408	      1234 ns/op	     136 B/op	       3 allocs/op
BenchmarkWrap-12                 	12277746	       102.4 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	12093657	        97.17 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	13137909	        98.78 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	11486084	       146.9 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	 6119178	       190.5 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	 6771130	       213.4 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	 5805573	       208.9 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	 7500080	       199.7 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	 7725237	       170.8 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	 7346791	       162.5 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrapf-12                	 4304510	       325.4 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 2264643	       444.3 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4388774	       260.7 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4571798	       244.3 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4943055	       256.8 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4012360	       344.0 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4577407	       250.4 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 5732811	       239.9 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 5005785	       267.7 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 5041876	       253.8 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapStatus-12           	12023872	       116.6 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	11267335	       103.8 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	14591094	        91.19 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	12721270	       104.0 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	10910032	       107.6 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	10697875	       103.7 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	12897732	        88.43 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	14921442	        85.57 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	13885178	        90.94 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	11712250	        95.89 ns/op	     128 B/op	       2 allocs/op
BenchmarkConvertNil-12           	522722493	         2.128 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	555292808	         3.189 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	402980470	         3.175 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	399750938	         2.659 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	466212789	         2.261 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	460455235	         2.497 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	505811365	         2.375 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	460066599	         3.026 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	480655722	         2.383 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	542449198	         2.288 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertPlain-12         	28622732	        53.46 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	17388433	        65.42 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	17713053	        61.01 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	17240304	        62.96 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	24169674	        47.70 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	29279409	        45.79 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	24575626	        49.42 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	27900148	        42.06 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	28627536	        43.06 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	35062426	        44.85 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertItem-12          	203011186	         6.082 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	144169401	         7.994 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	182670220	         6.293 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	201962506	         5.912 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	202780912	         9.762 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	100000000	        20.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	100000000	        12.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	100000000	        11.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	100000000	        12.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	98619201	        12.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	62047924	        17.90 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	75862077	        14.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	90452619	        12.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	100000000	        12.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	90864332	        12.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	98047218	        11.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	97064078	        11.85 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	87769174	        13.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	96342037	        14.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	99765483	        12.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	37545448	        32.91 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	40656027	        29.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	39806584	        29.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	39346764	        28.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	51250288	        26.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	40706221	        28.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	46060981	        29.96 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	44796110	        26.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	49629872	        24.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	50517262	        30.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	12199854	        96.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	11550501	        95.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	12495847	       107.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	10512334	       111.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	10894351	       111.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	11989239	       113.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	10339561	       115.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	10377777	       126.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	11767525	       117.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	11499895	       101.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 3267609	       383.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 3069709	       413.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2719460	       439.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 3647018	       374.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2715886	       416.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2802810	       376.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 3093502	       357.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2785249	       426.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 3152767	       358.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 3627974	       389.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	346215249	         2.986 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	422745504	         2.453 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	516056450	         2.294 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	503942550	         2.440 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	516769209	         2.631 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	343987970	         3.198 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	438006661	         3.335 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	459498331	         3.078 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	403545396	         2.793 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	412225752	         3.021 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	75008878	        14.11 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	83183671	        14.14 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	92398584	        14.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	85467090	        15.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	75095835	        19.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	63551242	        15.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	45018482	        24.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	77723388	        16.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	92074798	        13.27 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	88757296	        12.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackPlain-12       	 2612744	       480.4 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2175552	       534.8 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2415961	       505.4 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2375022	       550.8 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2141883	       576.2 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2163340	       477.5 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 1795983	       633.9 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 1819399	       655.2 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2264594	       512.9 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2175063	       544.9 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackExisting-12    	87164408	        14.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	73613823	        14.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	92557197	        16.39 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	78038846	        14.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	93533906	        13.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	90828301	        13.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	85822771	        12.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	94075618	        12.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	100000000	        12.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	92391000	        12.65 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	177658672	         6.716 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	177861632	         6.838 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	167476731	         7.365 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	167207019	         7.008 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	176100217	         6.708 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	186872630	         6.497 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	164114593	         7.106 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	173782960	         6.905 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	166964954	         7.434 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	159570064	         6.997 ns/op	       0 B/op	       0 allocs/op
BenchmarkFormatVerbose-12        	  750339	      1730 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  633590	      1658 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  625969	      1829 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  627627	      1967 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  545386	      1937 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  614248	      1857 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  646534	      1844 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  630134	      1802 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  586699	      2040 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  727568	      1853 ns/op	     744 B/op	      12 allocs/op
PASS
ok  	github.com/gtkit/goerr	214.498s
PASS
ok  	github.com/gtkit/goerr/test	0.676s
```

</details>

### 13.5 阶段验收、未通过项与停止位置

| 验收项 | 实际状态 | 证据与限制 |
| --- | --- | --- |
| 最小出口与兼容契约 | 已固定并回读 | 13.1：只保留拟新增 Inspect/Diagnostic/Cause 及列明字段，调用方普通输出函数可替换；所有新功能仍未实现 |
| Go 1.27 环境核对 | 已完成 | Go 1.27.2 与 GOROOT/平台/工具版本有真实记录；go.mod 仍为 1.26，未升级模块最低版本或安装工具 |
| 原始测试 / 新兼容测试 | 通过各自执行范围 | 原始两个包 race 通过；新增 5 组兼容测试通过；不外推为未实现快照的验收 |
| 缺栈反例 | 已复现，修复未完成 | 六个子用例按计划保持失败；最终全量测试退出 1，无其他失败或 race 报告 |
| build / vet | 通过 | 修改前后均退出 0 |
| 现代化 / lint 门禁 | 未通过 | 既有一个 SplitSeq 建议、6 项 lint 问题；不修改运行代码、不隐藏告警、不把已有问题顺带并入阶段 2 |
| 覆盖率门槛 | 原始基线未达标 | 原始根包 55.4% < 80%；未自动扩展为全仓补测，也没有声称新增用例后的覆盖率已达标 |
| 基准取证 | 已完成；耗时可比性受限 | 最终 17×10 样本与原始输出已落入本计划；分配量稳定；多组耗时波动超过 10%，严格耗时判定未验证 |
| 范围与兼容保护 | 已核对 | final-workspace-comparison.json：现有文件只有 error_test.go 变化；6 个运行 .go 文件、go.mod/go.sum、用户 .gitignore 指纹均保持；采样期间测试输入未变 |
| 后续能力与外部范围 | 未执行 / 未验证 | 新诊断读取、出口/日志实际接入、消费者与 Go 1.26 要求、远端 tag/Release、首次 MAJOR 判断及发布仍未验证；阶段 2～6 未授权 |

实际项目交付仅 OPTIMIZATION_PLAN.md、error_test.go（runtime 测试导入与追加用例）、新增 benchmark_test.go，以及 AGENTS.md 要求的真实错误记录及处置。没有修改运行 API、Item/Status 实现、码表、依赖、模块/版本元数据、其他说明、规则或 OpenSpec；用户已有 .gitignore 改动保留。原始命令日志、profile 与指纹保存在 `/private/tmp/goerr-stage1.E7FRRL`；关键环境、退出码、缺口和全部最终基准样本已保存在本计划。

合规结果：本轮无新增运行导出符号或外部依赖，无库内日志、日志接口或 Options；新增测试遵循现有错误身份与文本契约，基准报告分配。覆盖率、原有现代化/lint、缺陷修复和稳定耗时比较不能勾选通过。

停止位置：阶段 1 执行与取证结束，已知缺口如上；这不是完整质量验收通过或阶段 2 授权。后续只有明确授权后才从对应阶段步骤继续；本轮不自动修复缺陷、扩展计划范围、提交、推送或发版。

回滚定位：撤回本轮新增 benchmark_test.go、error_test.go 的 runtime 测试导入与追加块，以及本计划的本轮记录；先对照实际差异，保护用户现有改动。错误记录保留真实历史与处置，不删除既有条目。本段仅给出回滚方法，未执行回滚。

交付回读：已实际读取本轮授权、阶段状态、兼容契约、结果矩阵与统计；计划内 170 个基准样本与保存的原始输出逐字比较一致，受保护文件指纹一致，阶段 2 仍未授权，失败和未验证项均保留。对应 plan-readback.json 的上述项目均为 true；回读未触发运行代码修改或新的测试流程。

## 14. 阶段 2 实施细化与执行记录（2026-10-21）

### 14.1 本轮授权、现状与约束（代码修改前保存）

- 需求来源：用户在阶段 1 汇报后要求“继续”，按原顺序执行下一阶段 2，仅包括 2.1/2.2；不自动进入阶段 3，不实现 Inspect 或私有遍历器，不扩展导出面、日志接口、配置或依赖。
- 已重新读取最新全局/项目规则、错误记录、原计划、相关 Guide/skill 和现有实现/调用点；当前 Go 仍为 go1.27.2 darwin/amd64，GOROOT=/usr/local/go，GOWORK 为空。阶段 1 的运行输入指纹已重新比较一致，不把历史命令当作本轮执行结果。
- 当前源码的 WithStack 在 AsType 找到任何 Item 时返回原输入，6 个缺栈反例已在阶段 1 真实复现；版本仍为 v1.4.0、go.mod 仍为 1.26、零外部依赖。本轮尚未修改实现。
- 原实现 build/vet 和旧兼容测试可用；原有 6 项 lint、SplitSeq 建议、覆盖率门槛和不稳定耗时已分清来源。只交付阶段 2 范围，不修 Format、ResolveStatus 测试、截断逻辑、码表或为覆盖率扩展全仓测试。
- 文件范围：error.go（仅 WithStack 及其 GoDoc）、go.mod（仅最低版本）、error_test.go（补栈/错误链回归）、README.md、CHANGELOG.md、CONTRIBUTING.md（仅最低版本与补栈说明）、本计划、真实错误记录及关闭。benchmark_test.go/go.sum/version.go 保持；其他用户或阶段 1 改动保留。
- 证据目录 `/private/tmp/goerr-stage2.WY4TCO` 已保存 before 副本、指纹、HEAD 和工作区状态；不重置用户文件。测试进程采用与阶段 1 相同的 local/readonly/offline 约束，缓存和日志写入本轮证据目录，不改用户环境配置。

### 14.2 零决策实施步骤与验收

| 步骤 | 前置条件 | 具体改动 | 检查及可复核证据 | 停止条件 | 交付物 / 当前状态 |
| --- | --- | --- | --- | --- | --- |
| 2.1a | 本节已保存并回读；阶段 1 契约与反例可复核 | WithStack 只做一次 AsType；首匹配 typed-nil 或 len(stack)>0 保留原 err；其他情况新建 Item，用原 err 全文和直接 cause；缺栈 Item 保留 StatusInfo/Message，普通错误仍内部错误兜底；直接 callers(skipDirect) 保持调用方首帧 | 原 6 例转绿，追加空而非 nil 栈切片；nil、typed-nil、有栈身份、成功/自定义状态、原对象不变、幂等、首帧和下游栈共享 | 改旧查询选择、状态/全文/cause 漂移，出现新导出或修改原对象 | 已实施；7 个缺栈子例及旧兼容 race 通过 |
| 2.1b | 2.1a 同一变更范围；Go 1.27.2 已核实 | go.mod 改 `go 1.27.0`；README/CONTRIBUTING 写明最低版本与无栈/有栈/typed-nil 行为；CHANGELOG 只新增 Unreleased 的 Changed/Fixed，不改历史、Version、tag 或 module 路径 | 保存后回读声明与说明；真实 build/vet/race；核对新增运行导出符号为零、go.sum 不变 | 需要依赖安装、首次 MAJOR 或改码表/发布元数据 | 已实施并回读；实际 build/vet/race 通过 |
| 2.2 | 2.1 实施完成，未扩大范围 | 用既有 API 加错误链表格：Join 的缺栈首分支及逆序已有栈、多 %w、共享分支引用、自定义 As 的无栈/typed-nil、第三方 Error；所有包装保留标准外层文本与原 cause，旧主状态选择不变；下游 Wrap 复用新捕获栈 | race 下运行新增表格及完整测试集，按实际命令/退出码记录；只在文档固定阶段 3 遍历反例 | 引入新查询、提前实现遍历器，或兼容用例无依据 | 已执行；7 个链子例通过，未实现遍历器 |
| 阶段收口 | 上述范围已验证 | 按 AGENTS 跑 build、只读 fix、vet、lint、完整 race、覆盖率、现有基准 count=3；仅记录本轮引入的问题并在范围内纠正 | 每条命令独立输出/退出码；分清既有失败；保留原始样本，比较 B/op/allocs，不宣称 10% 耗时门槛通过 | 新失败、scope 不清或同一问题自纠 3 次；不自动修复未授权基线问题 | 已完成取证与记录；遗留门禁如实未通过 |

调用方与适配：当前模块内 WithStack 的引用在 error_test.go 与 benchmark_test.go，按已读源码/文本范围核对，原有有效栈、nil/typed-nil 调用无需改适配；缺栈输入的身份变化是本次明确修复，原对象由 cause 与 errors.Is 保留。仓库外消费者、动态函数值及其他平台未验证，不宣称全调用图或旧 Go 工具链兼容。

未来遍历用例只记录预期：标准外层前缀、Join/多 %w 实际子项顺序、不可比较错误、共享引用与循环的区别、自定义 As 虚拟对象不造真实节点、nil/typed-nil 和 256/16384 预算边界；阶段 2 不为这些用例加入未调用的私有读取实现。

兼容与风险：Go 最低声明提高将排除 1.26 工具链；不把这一变化称作仍兼容旧工具链。仅无栈、非 nil Item 路径允许新包装/新栈和对应分配成本，已带栈与 typed-nil 保持身份；不改变 Item 布局、默认状态、公开文案、错误码或 Error 方法集。没有存储迁移或外部调用。阶段 1 的耗时可比性不足，阶段 2 的小样本只核对基准执行与分配，不以 ns/op 波动证明优化收益。

失败与回滚：新兼容失败先记录并停在本阶段，不通过 skip/忽略项降低门槛；已存在且未扩大授权的门禁问题继续报告。回退时仅撤回本轮 WithStack 差异、go.mod 最低声明、阶段 2 追加测试及必要说明，参照 before 副本与实际差异保护阶段 1 测试和用户 .gitignore，不强制重置；先撤回实际 1.27 代码/声明再核对旧工具链，不能凭数字回退宣称已兼容。回滚方法尚未执行。

### 14.3 实际执行与阶段验收结果

步骤 2.1 已实施 WithStack 最小修复、go.mod 1.27.0 及必要说明；修改后已回读。`go test -race -count=1 -timeout=5m -run 'TestWithStack|TestCompatibility|TestStack|TestWrapFamily|TestTypedNil|TestNilStatus|TestItem_nilReceiver' ./...` 本轮退出 0，根包匹配用例通过；test 子包无匹配用例。原 6 个补栈反例及新增空而非 nil 栈切片已转绿，证据为 step21-race.txt/json；此命令不是完整测试集。步骤 2.2 已按保存清单加入 7 个错误链子用例并回读、完整验证通过；未实现未来遍历器。未执行回滚或发布。

完整验证及基准均已完成；下列是本轮实际结果，不覆盖阶段 1 的历史失败：

| 命令 | 实际退出码 | 结果与证据（证据目录内） |
| --- | --- | --- |
| `go test -json -race -count=1 -timeout=5m ./...` | 0 | 两个包通过，75 个 pass 事件、无 fail/skip/race；7 个缺栈子用例和 7 个错误链子用例均 pass；full-race.txt/json、full-race-summary.json；计数含父/子事件，不是 75 个独立顶层测试 |
| `go build ./...` | 0 | 无输出；build.txt/json |
| `go vet ./...` | 0 | 无输出；vet.txt/json |
| `golangci-lint run ./...` | 1 | 与阶段 1 相同的 6 项：error.go 的 Format 五个 fmt 写入 errcheck，resolve_test.go:24 的 ineffassign；lint.txt/json；没有新增问题 |
| `go fix -diff ./...` | 1 | 与阶段 1 相同的 stackFrames 循环 Split→SplitSeq 建议；fix.txt/json；只读未应用，未借修复扩大格式/现代化范围 |
| `go test -count=1 -coverprofile=/private/tmp/goerr-stage2.WY4TCO/coverage.out ./...` | 0 | 根包 56.3%，test 子包无可统计语句；仍低于 80%，命令成功不代表覆盖率门槛通过；coverage.txt/json/out |

已关闭 ERR-20261021-006：原缺栈反例、本轮新增边界和错误链回归证实修复；ERR-20261021-008 保留本轮既有门禁/覆盖率缺口，003/004 仍开放。go.mod 与说明已回读；source-review.json 的实际结果表明 error.go 除 WithStack/GoDoc 外逐字未变，.gitignore/go.sum/version.go/benchmark_test.go/码表/Status/其他运行文件均与本轮开始时相同。新增运行导出符号或签名变化为零；实际新增的测试函数不属于库运行 API。


### 14.4 基准与成本记录

`go test -bench=. -benchmem -count=3 ./...` 退出 **0**，共 17 组、51 个样本，两个包通过；实际起止时间（Asia/Taipei）为 `2026-10-21T16:38:49+08:00` 至 `2026-10-21T16:39:55+08:00`。本轮未改 benchmark_test.go，CPU/平台/并发度与阶段 1 实际输出一致；原始命令与输出为 bench.txt/json，解析样本为 samples.json，分配比较为 allocation-comparison.json。

`benchstat /private/tmp/goerr-stage2.WY4TCO/bench.txt` 退出 **0**；实际警告为每组至少 6 个样本才能给出当前 95% 区间，当前区间显示 ±∞。此处按 AGENTS 的 count=3 检查既有基准执行和分配样本，不把样本量或单组中位数冒充耗时回退通过；阶段 1 的大幅波动也仍存在于比较前提中，10% 耗时门槛继续未验证。

| 基准（GOMAXPROCS=12） | 每样本 B/op | 每样本 allocs/op | 样本数 |
| --- | --- | --- | --- |
| New-12 | 136 | 3 | 3 |
| Wrap-12 | 104 | 2 | 3 |
| Wrapf-12 | 120 | 3 | 3 |
| WrapStatus-12 | 128 | 2 | 3 |
| ConvertNil-12 | 0 | 0 | 3 |
| ConvertPlain-12 | 80 | 1 | 3 |
| ConvertItem-12 | 0 | 0 | 3 |
| ConvertDepth/1-12 | 0 | 0 | 3 |
| ConvertDepth/4-12 | 0 | 0 | 3 |
| ConvertDepth/16-12 | 0 | 0 | 3 |
| ConvertDepth/64-12 | 0 | 0 | 3 |
| ConvertParallel-12 | 0 | 0 | 3 |
| HasCode-12 | 0 | 0 | 3 |
| WithStackPlain-12 | 112 | 2 | 3 |
| WithStackExisting-12 | 0 | 0 | 3 |
| WithStackMissing-12 | 112 | 2 | 3 |
| FormatVerbose-12 | 744 | 12 | 3 |

与阶段 1 最终样本比较，只有 WithStackMissing 的分配数据变化：0 B/op、0 allocs/op → 112 B/op、2 allocs/op；这是缺陷修复后实际抓栈与新包装的成本，不要求保持错误的零捕获成本。其他 16 组的 B/op/allocs/op 与对应旧样本一致；该结论仅覆盖这些夹具，不证明生产负载或耗时收益。没有实施缓存或为波动追加优化。

<details>
<summary>阶段 2 基准原始输出（17 组 × 3 次）</summary>

```text
goos: darwin
goarch: amd64
pkg: github.com/gtkit/goerr
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkNew-12                  	 1000000	      1021 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1000000	      1004 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1182777	      1019 ns/op	     136 B/op	       3 allocs/op
BenchmarkWrap-12                 	12194799	        97.99 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	11942175	       104.2 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	11469201	       100.1 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrapf-12                	 5240856	       234.8 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 5194674	       236.0 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4921024	       237.0 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapStatus-12           	12467814	        96.80 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	12565153	        97.14 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	12264240	        99.34 ns/op	     128 B/op	       2 allocs/op
BenchmarkConvertNil-12           	398997398	         2.636 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	475905789	         2.504 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	436161472	         2.522 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertPlain-12         	27439004	        44.25 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	26417694	        45.19 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	26302770	        44.41 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertItem-12          	182972235	         6.624 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	182357083	         6.674 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	182037795	         6.605 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	91868652	        13.08 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	93411630	        13.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	88284648	        13.38 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	35458237	        33.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	37148563	        32.28 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	37989681	        32.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	 9987159	       117.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	 9354921	       119.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	10284382	       116.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2869485	       465.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2864619	       421.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2780720	       440.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	239766943	         7.273 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	275779233	         5.843 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	241336758	         4.938 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	70915472	        14.77 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	76671897	        15.17 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	86515180	        15.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackPlain-12       	 2360846	       506.8 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2317681	       506.7 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2456216	       493.6 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackExisting-12    	87016852	        13.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	91991121	        14.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	65864469	        17.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	 2210932	       595.0 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackMissing-12     	 2311236	       503.4 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackMissing-12     	 2422010	       489.2 ns/op	     112 B/op	       2 allocs/op
BenchmarkFormatVerbose-12        	  628971	      1935 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  659121	      1945 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  554967	      2409 ns/op	     744 B/op	      12 allocs/op
PASS
ok  	github.com/gtkit/goerr	64.089s
PASS
ok  	github.com/gtkit/goerr/test	0.670s
```

</details>

### 14.5 阶段验收、回滚定位与停止位置

| 验收项 | 实际状态 | 依据 / 限制 |
| --- | --- | --- |
| 2.1a 补栈修复 | 通过本阶段回归 | 7 个缺栈子用例（含原 6 例）、nil/typed-nil、有栈身份、全文、状态/文案、原 cause、首帧及二次复用均在 race 下通过 |
| 2.1b Go 与说明 | 已实施并回读 | go.mod 为 1.27.0，README/CONTRIBUTING/Unreleased 对应；实际测试 Go 1.27.2，旧工具链不再满足最低声明；Version 仍 v1.4.0，未发版 |
| 2.2 错误链 | 通过本阶段回归 | 7 个链用例通过，Join/多 %w/共享分支/自定义 As/第三方 Error 保留旧主状态与 errors.Is 可达性；下游 Wrap 共享新捕获栈 |
| 整体功能检查 | 通过本轮范围 | 两个包完整 race 通过；build/vet 退出 0；未声明其他平台或仓库外消费者已适配 |
| 现代化 / lint | 未通过，既有缺口保留 | fix 退出 1（一个旧 helper 建议），lint 退出 1（原有 6 项）；不称全部质量门禁已过，不在本阶段扩大修复 |
| 覆盖率 | 命令成功，门槛未达 | 根包 56.3% < 80%，未为凑数字扩展补测；无合成覆盖结果 |
| 性能 | 分配样本已核对，耗时未验证 | 17×3 原始样本如上；正常旧路径的分配未增加，补栈成本单列；不宣称 10% 耗时预算通过 |
| 导出面 / 依赖 / 范围 | 已静态核对 | error.go 只有 WithStack/GoDoc 变化；Error/Item/Status 与运行签名未改；无新依赖、日志接口、Option、目录或私有遍历器；go.sum/version.go/benchmark 与用户 .gitignore 未改 |
| 历史缺陷记录 | 已处置对应条目 | 006 已按真实回归关闭；003/004/008 保留为未解决门禁/覆盖率记录；不将历史失败改写为从未发生 |

交付文件：error.go、go.mod、error_test.go、README.md、CONTRIBUTING.md、CHANGELOG.md、OPTIMIZATION_PLAN.md，以及按 AGENTS 追加/关闭的真实错误记录。阶段 1 测试、基准和用户改动保留；没有新增运行导出 API，未进入阶段 3～6、未执行提交/推送/tag/Release，也未首次升级 v2。

回滚定位：参照本轮 before 副本撤回 WithStack/GoDoc、go.mod 的最低版本、阶段 2 的空栈片子例与 TestWithStackChains、三份说明的阶段 2 差异，以及计划中本轮记录；不能用 reset/checkout 覆盖阶段 1 测试或用户改动。无外部状态、依赖或存储迁移要回滚。回滚仅为可审阅方法，未实际执行或验证旧 Go 工具链。

停止位置：阶段 2 的补栈与错误链兼容范围已完成；质量遗留项和耗时比较不算通过。仅在后续明确授权时继续对应阶段；本轮停在阶段 2，不自动进入快照实现。


## 15. 阶段 3：最小诊断快照执行记录

### 15.1 授权、当前事实与版本前置核对（代码改动前保存）

- 用户在阶段 2 完成汇报后要求“继续”；本轮仅执行既定步骤 3.1、3.2，不进入阶段 4～6。需求为主动读取内部诊断快照，并保留普通调用方函数作为可插拔出口。
- 已实际读取本轮适用全局/项目规则、错误记录、原计划、错误/兼容/结构/测试/文档/现代 Go Guide，以及 harness-implement、use-modern-go、go-error-handling、go-testing、go-enterprise-quality、go-concurrency 与 go-performance 相关正文/参考。读取不代表门禁已通过，不启动 OpenSpec。
- 本轮前状态为阶段 1/2 的未提交改动及用户 `.gitignore` 改动；新增 diagnostic 文件尚不存在。已将 21 份现有代码/测试/模块/文档及用户改动备份到 `/private/tmp/goerr-stage3.r58q5K/before`，指纹见 before-sha256.json；不以 HEAD 覆盖这些差异。
- 2026-10-21 重新执行 `/usr/local/bin/go version`：go1.27.2 darwin/amd64；go.mod 为 `go 1.27.0`、模块 `github.com/gtkit/goerr`、无外部依赖；Version 和最高本地 tag 均 v1.4.0。
- 远端核对：`git ls-remote --tags https://github.com/gtkit/goerr.git` 退出 0，最高正式 tag v1.4.0；GitHub `/repos/gtkit/goerr/releases?per_page=100` 经系统 curl 正常 TLS 校验读取，退出 0、列表为 `[]`。原始输出/时间/命令见本轮目录 remote-tags.txt/json、remote-releases.txt/json。由此核定本次新增能力的未发布 MINOR 元数据为 **v1.5.0**，保留 v1 module/import；不创建 tag/Release、不执行 Makefile 发布目标，CHANGELOG 继续用 Unreleased 并注明目标版本，不能冒充已发布。
- 前置查询中的网页 Cache miss 与 Python CA 失败已分别记录 ERR-20261021-009/010；系统 curl 与 git 的真实成功不抹去先前失败，不修改证书/用户配置、不禁用 TLS 校验。
- 既有缺口继续保留：Format 5 处 errcheck、resolve_test 1 处 ineffassign、旧 stackFrames 的 SplitSeq 建议、整体覆盖率未达 80%、耗时基线不稳定。它们不是本阶段修复项。

### 15.2 零决策实施步骤与验收（执行前已保存，状态按实际结果更新）

所有语义沿用第 13.1 节字段与契约，新增生产符号严格为 Inspect、Diagnostic、Cause；私有状态不成为选项或扩展框架。

| 步骤 | 前置条件 | 具体改动与文件 | 检查与可复核证据 | 停止条件 | 交付物 / 当前状态 |
| --- | --- | --- | --- | --- | --- |
| 3.1a 入口与有界读取 | 本节落盘并实际回读；远端版本已核定 | 新建根目录 diagnostic.go；实现三个既定导出符号、首位带键守卫与字段 GoDoc；私有读取状态、显式 DFS 工作栈，按子槽逐个推进；仅两个保护常量 256 与 16384 | diagnostic_test.go 覆盖 nil/非 nil 成功/零值 Item、普通错误、标准包装、Join/多 %w/索引参数顺序、状态覆盖、自定义 As 与第三方 Error；主状态与旧 Convert 比较、根全文相等 | 出现额外公共符号、改变旧路径/Item、隐式编码/日志、丢前缀、主状态漂移 | 生产实现与契约测试已完成；见 15.4 |
| 3.1b 边界与已有栈 | 3.1a 实现；不调用外部方法来判 typed-nil | 同文件实现当前路径循环检测：先 reflect.Value.Comparable 再加入局部 map；不可比较节点不作 interface 哈希。nil 子切片槽消耗预算；单值 Unwrap 返回 nil 表示没有子节点。未知 typed-nil 在调用方法前失败，已知 nil Item 按原契约。只解析真实 Item 的既有 PC；沿当前路径同一首元素地址及相同长度的栈只在首次实际携带节点保存，兄弟分支各自保存 | 直接/间接/不可比较循环、共享引用、嵌接口的不可比较值、宽 nil 分支、256/257 访问、16384/16385 rune（含无效 UTF-8、公开文案及帧文本）；有栈/无栈、路径共享/兄弟栈、调用方起源帧 | 预算绕过、比较 panic、循环后无界读取、诊断位置抓新栈、预算失败返回部分值 | 对应反例在 race 下通过；见 15.4 |
| 3.2a 原子性与所有权 | 3.1 完成 | 原因树有界遍历中按 errors.AsType 的前序规则记录首个 Item 匹配：直接类型优先，再调用节点自定义 As，命中后不再调用后续 As；保持虚拟 Item/typed-nil 的旧主状态选择。遍历后只读取主状态与根全文，不再次调用 Convert 或无界重遍历。新增数据由每次调用独占，不保存 error/Status/原 PC、无全局缓存 | 返回 Causes/Frames 篡改后原对象与另一快照不变；并发读取/各自修改返回值；错误输入后返回严格 Diagnostic{} 与非 nil error；计数夹具验证错误结构失败不再读取根 Error，不对虚拟 Item 建节点 | 部分成功、返回共享可写数据、引入缓存/锁/接口或修改业务错误 | 读取所有权与失败契约已验证；见 15.4 |
| 3.2b 对外说明与版本 | 3.2a 完成且新契约验证通过 | example_test.go 增加外部 ExampleInspect（实际检查诊断 error）及带键值用法；README 加入口/语义/限制/内部信息边界，CONTRIBUTING 加快照矩阵；CHANGELOG Unreleased 记 Added 与目标 v1.5.0，version.go 更新相同未发布 MINOR | 运行 Example，GoDoc/README/字段回读；导出清单与版本/module/依赖对比；只在新增 Example/测试中调用，不接入旧运行链 | 把快照当客户端 DTO、虚构 logger API/版本、写成已发布，或扩大到两个出口组合 | 文档、Example、版本元数据已核对；见 15.4 |
| 3.2c 本阶段验证与记录 | 上述代码/说明完成 | benchmark_test.go 只追加 Inspect 的 nil/普通/有栈/无栈、1/4/16/64 层、分支与访问上限基准；保持旧基准不改。按本节命令执行、保存环境/命令/退出码/原始结果，更新本计划 | 全量 race/build/vet；lint/fix 与已知缺口逐项比较；整体及 diagnostic 文件覆盖；基准 3 次样本、benchstat 与旧路径分配对比；检查范围和所有权 | 新增失败或旧状态/分配漂移；同一问题自纠 3 轮；不能复核证据 | 验证与取证已完成；既有门禁及耗时缺口保留；见 15.4 |

读取算法补充：DFS 每个单值子节点、每个多值切片子槽各访问一次；不复制宽子切片，不先排满工作队列。不可比较循环由访问预算终止；当前路径的可比较重复节点立即失败，共享兄弟对象不会误报。自定义 As 保留标准库首个匹配规则且与实际节点区分；外部方法必须有限、无阻塞、不 panic，读取期间输入保持稳定，本库不隔离方法内部的阻塞、分配或副作用。reflect 只用于主动诊断路径的 typed-nil/可比较性保护，不进入旧路径。

帧从已有 PC 经 runtime.CallersFrames 解析为独立 runtime.Frame 切片；空/无效 PC 无有效帧时保持 nil，Func 仅为运行时只读元数据。文本预算按 Go range 逐 rune 累减，累计 Error、Message、保存帧的 File/Function，超过时返回零快照，原字符串不重写、不截断。

验收命令采用本轮 run.py 的逐条退出码记录，GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off；Go/lint cache 与覆盖产物放本轮临时目录。先定向 `go test -json -race -count=1 -timeout=5m -run 'TestInspect|ExampleInspect' ./...`，再全量 race、`go build ./...`、`go vet ./...`、`golangci-lint run ./...`、`go fix -diff ./...`、`go test -count=1 -coverprofile=<本轮目录>/coverage.out ./...`；CPU 检查完成后单独 `go test -bench=. -benchmem -count=3 ./...` 并用既有 benchstat。固定环境只证明本机构建；不安装工具、不执行容器、数据库或发布流程。3 次样本不能消除阶段 1 耗时噪声，耗时门槛仍以可信度为前提，不为凑结论反复采样。

### 15.3 兼容、风险、失败与回滚

- 旧 Error 接口、Item/Status 结构、码表、New/Wrap/Convert/HasCode/WithStack 与弃用入口均不改；新增调用点只在新增测试/Example/基准中，外部消费者和其他平台未验证，不声称完整调用图覆盖。
- 风险：无限/阻塞或不稳定的外部方法无法由访问预算隔离；快照含内部全文/路径，不是客户端响应，调用方负责字段选择与脱敏；新增函数/字段形成 v1 的长期公共契约，因此不补无关入口、getter、Options 或自动日志。
- 任一读取失败为零快照和独立诊断 error，不替换原业务 error，不回退输出未经处理原文。真实新失败先记错误记录，定位后只修本阶段引入问题；历史质量缺口不自动修复。
- 回滚前先撤除新增 Inspect/Diagnostic/Cause 消费引用，再撤回 diagnostic.go、diagnostic_test.go、新增 Example/基准与本阶段文档及版本差异；参照本轮 before 副本逐项恢复，保留阶段 1/2 和用户改动。无外部服务、存储、依赖或后台资源迁移；不执行 reset/checkout 覆盖工作区。回滚方法待审阅，不宣称已真实回滚。
- 本轮停止位置固定为阶段 3：实现、验收、写回计划并回读后汇报，阶段 4/5/6 仍需用户继续授权。

### 15.4 实际结果

本阶段实现与功能验收已完成；整体质量门禁和耗时限制分别记录，不将未通过项勾选通过。


实际环境与产物：`/usr/local/bin/go` go1.27.2 darwin/amd64、GOAMD64=v1、CGO_ENABLED=1、GOWORK 为空；CPU 为基准实际头部的 Intel Core i7-9750H @ 2.60GHz，并发度 12。golangci-lint 2.14.0（go1.27.2 构建）、既有 benchstat 的 x/perf 版本 v0.0.0-20260908200009-22c9c6c9d4da（go1.27.1 构建）。environment.json 保存各命令真实输出及退出码，GOGC/GOMEMLIMIT/GOMAXPROCS/GOEXPERIMENT/GODEBUG 未设覆盖；供电与后台负载仍未控制，不能把本机样本扩展为其他平台或生产性能。

本轮完整证据目录为 `/private/tmp/goerr-stage3.r58q5K`，每项 run.py 命令都有独立 `.txt` / `.json`，verification-index.json 汇总最终命令、参数、环境和起止时间。网络查询仅核对版本，代码检查使用 `GOTOOLCHAIN=local`、`GOFLAGS=-mod=readonly`、`GOPROXY=off`；缓存/覆盖产物均在临时目录，未生成仓库内测试产物。

| 检查 | 真实命令（省略统一 run.py 前缀） | 退出码 | 结果 / 证据名 |
| --- | --- | --- | --- |
| targeted-race | `/usr/local/bin/go test -json -race -count=1 -timeout=5m -run 'TestInspect\|ExampleInspect' ./...` | 0 | 49 个 pass 事件，无 fail/skip/race；定向快照与 Example 验证；targeted-race.txt/json |
| final-race | `/usr/local/bin/go test -json -race -count=1 -timeout=5m ./...` | 0 | 两个包通过，125 个 pass 事件（含父/子测试与包事件），无 fail/skip/race；不是同数量的顶层测试；final-race.txt/json |
| final-build | `/usr/local/bin/go build ./...` | 0 | 通过；final-build.txt/json |
| final-vet | `/usr/local/bin/go vet ./...` | 0 | 通过；final-vet.txt/json |
| final-lint | `/Users/xiaozhaofu/go/bin/golangci-lint run ./...` | 1 | 未通过：只剩原有 Format 5 项 errcheck、resolve_test 1 项 ineffassign；新增文件/夹具无 lint 项；final-lint.txt/json |
| final-fix | `/usr/local/bin/go fix -diff ./...` | 1 | 未通过：仅旧 error_test.go:290 的 Split → SplitSeq 建议；未应用改写；final-fix.txt/json |
| final-coverage | `/usr/local/bin/go test -count=1 -coverprofile=/private/tmp/goerr-stage3.r58q5K/final-coverage.out ./...` | 0 | 命令成功，整体 66.2%，低于 80%；外部 test 子包无可统计语句；final-coverage.txt/json |
| final-bench | `/usr/local/bin/go test -bench=. -benchmem -count=3 ./...` | 0 | 两个包通过，27 组各 3 次、共 81 个样本；包含 17 组旧路径及 10 组 Inspect；final-bench.txt/json |
| final-benchstat | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-stage2.WY4TCO/bench.txt /private/tmp/goerr-stage3.r58q5K/final-bench.txt` | 0 | 统计输出已保留；每组 3 次不足以取得 95% 区间或判断显著差异；final-benchstat.txt/json |
| coverage-functions | `/usr/local/bin/go tool cover -func=/private/tmp/goerr-stage3.r58q5K/final-coverage.out` | 0 | Inspect/walk/frames/reserve 均 100.0%；整体 66.2%；coverage-functions.txt/json |
| inspect-godoc | `/usr/local/bin/go doc Inspect` | 0 | 新增入口 GoDoc 已回读；inspect-godoc.txt/json |
| diagnostic-godoc | `/usr/local/bin/go doc Diagnostic` | 0 | 字段/所有权/首位未导出守卫的文档已回读；diagnostic-godoc.txt/json |
| cause-godoc | `/usr/local/bin/go doc Cause` | 0 | 节点与已有栈字段文档已回读；cause-godoc.txt/json |

定向与完整回归包括：真 nil、非 nil 成功 Item、零值/typed-nil、根前缀、Join 与多个 %w（含索引参数顺序）、状态覆盖、自定义 As 的首个匹配/虚拟对象/false target 复用、第三方 Error 的实际节点业务码、不可比较值及接口内不可比较值、共享兄弟对象、直接/间接/不可比较循环、nil 子槽与极宽分支、256/257 访问边界、16384/16385 rune（含无效 UTF-8、公开文案和帧文本）、已有栈/空或无效栈、同路径去重/兄弟分支归属、返回值篡改和 32×16 次并发读取后各自修改快照。失败断言要求严格零 Diagnostic 与非 nil error；循环计数夹具确认不再读取根 Error 或重复调用 As/Unwrap。

新增 diagnostic.go 的 **86 条可统计语句中 86 条被覆盖（100.0%）**；这是本文件语句覆盖，不是所有外部错误实现、平台或方法行为已经验证。整体覆盖率 66.2% 仍未达标，不据新增文件覆盖情况勾选全仓 80% 通过。

首轮完整检查也通过 race/build/vet，但 lint 多出两处测试/基准夹具 ST1023；已按 ERR-20261021-011 先记录，再仅将已是 error 接口的 io.EOF 初始化改为类型推断。最终完整检查重跑，新增两项消失，原 6 项保留。009/010 的版本查询失败已通过替代只读途径完成核对，原工具失败仍保留；012 是后处理猜测阶段 2 文件名导致退出 1，实际基准退出 0，两者没有混同。已按清单改用真实 bench.txt/bench.json，后处理成功且回读；没有因此重采样或修改运行代码。009～012 均按真实处置关闭，003/004/008 的既有缺口未关闭。回读时发现的命令表转义/空格和概括措辞问题另记 013；已只修计划呈现，13 项命令表与原始 JSON 逐项比较一致，回读后关闭。

### 15.5 性能原始证据与解释

旧 17 组与阶段 2 的 B/op、allocs/op 样本集合逐组相同，`benchmark-analysis.json` 的 old_allocation_changes 为空；WithStackMissing 保持阶段 2 修复后 112 B/op、2 allocs/op，不拿阶段 1 的未补栈零分配作为本轮正常基线。既有错误处理、码表与脱敏实现亦与本轮 before 副本一致（Version 元数据按计划更新），新路径成本由调用 Inspect 的调用方承担。

以下是本轮每组 3 次样本的中位耗时和实测字节/分配集合，取值来自原始结果，用 Python statistics.median 计算。它们是本轮微基准观测，不是吞吐、p99、容量承诺或显著性结论：

| 组 | 样本数 | ns/op 中位数 | B/op 样本取值 | allocs/op 样本取值 |
| --- | --- | --- | --- | --- |
| BenchmarkNew | 3 | 1014 | 136 | 3 |
| BenchmarkWrap | 3 | 102.1 | 104 | 2 |
| BenchmarkWrapf | 3 | 240.9 | 120 | 3 |
| BenchmarkWrapStatus | 3 | 99.43 | 128 | 2 |
| BenchmarkConvertNil | 3 | 2.42 | 0 | 0 |
| BenchmarkConvertPlain | 3 | 46.83 | 80 | 1 |
| BenchmarkConvertItem | 3 | 6.809 | 0 | 0 |
| BenchmarkConvertDepth/1 | 3 | 13.43 | 0 | 0 |
| BenchmarkConvertDepth/4 | 3 | 43.53 | 0 | 0 |
| BenchmarkConvertDepth/16 | 3 | 139.7 | 0 | 0 |
| BenchmarkConvertDepth/64 | 3 | 500.8 | 0 | 0 |
| BenchmarkConvertParallel | 3 | 4.85 | 0 | 0 |
| BenchmarkHasCode | 3 | 15.83 | 0 | 0 |
| BenchmarkWithStackPlain | 3 | 523 | 112 | 2 |
| BenchmarkWithStackExisting | 3 | 13.31 | 0 | 0 |
| BenchmarkWithStackMissing | 3 | 506.4 | 112 | 2 |
| BenchmarkFormatVerbose | 3 | 2006 | 744 | 12 |
| BenchmarkInspect/nil | 3 | 9.558 | 0 | 0 |
| BenchmarkInspect/plain | 3 | 441.7 | 289 | 4 |
| BenchmarkInspect/item | 3 | 2393 | 1618 | 11 |
| BenchmarkInspect/no_stack | 3 | 853.1 | 754 | 7 |
| BenchmarkInspect/branches | 3 | 5709 | 4031 | 26 |
| BenchmarkInspect/depth_1 | 3 | 3071 | 1779 | 13 |
| BenchmarkInspect/depth_4 | 3 | 3997 | 2998 | 18 |
| BenchmarkInspect/depth_16 | 3 | 9874 | 12394 | 39 |
| BenchmarkInspect/depth_64 | 3 | 26262 | 29018 | 93 |
| BenchmarkInspect/visit_limit | 3 | 85614 | 116088, 116089 | 287 |

耗时门槛 **未验证**，不能报告 10% 预算通过。与阶段 2 中位数直接计算时，下列观测已超过名义阈值（计算为 本轮/阶段2 - 1）：

| 旧路径 | 阶段 2 ns/op 中位数 | 本轮 ns/op 中位数 | 观测差异 |
| --- | --- | --- | --- |
| BenchmarkConvertDepth/4 | 32.28 | 43.53 | +34.85% |
| BenchmarkConvertDepth/16 | 117.8 | 139.7 | +18.59% |
| BenchmarkConvertDepth/64 | 440.3 | 500.8 | +13.74% |

benchstat 明确给出 `need >= 6 samples for confidence interval at level 0.95` 与 `need >= 4 samples to detect a difference at alpha level 0.05`，区间为 ±∞。加上阶段 1 基线噪声与供电/后台负载未控制，不能将上述观测排除为无问题，也不能据此认定是新代码导致的性能退化；保留实际超阈值观测及未验证状态。不用含新增组的 geomean 代替旧路径验收，不在本阶段自动扩展采样、缓存或性能改造。

最终基准夹具 SHA-256（从保存的真实 source-review.json 引用）：`806d95fc70af21a448fee0ff3fc9ef44562c0403bd985c7be13bc9bafbae7226`。原始全部样本如下，final-bench.txt/json 和 final-benchstat.txt/json 同时保留：

```text
goos: darwin
goarch: amd64
pkg: github.com/gtkit/goerr
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkNew-12                  	 1000000	      1043 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1000000	      1014 ns/op	     136 B/op	       3 allocs/op
BenchmarkNew-12                  	 1248512	       971.2 ns/op	     136 B/op	       3 allocs/op
BenchmarkWrap-12                 	11941898	        98.89 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	12300891	       102.1 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrap-12                 	12290130	       102.9 ns/op	     104 B/op	       2 allocs/op
BenchmarkWrapf-12                	 5223624	       239.7 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 4858017	       250.6 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapf-12                	 5051485	       240.9 ns/op	     120 B/op	       3 allocs/op
BenchmarkWrapStatus-12           	11550325	       101.6 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	11688961	        99.43 ns/op	     128 B/op	       2 allocs/op
BenchmarkWrapStatus-12           	11154927	        99.08 ns/op	     128 B/op	       2 allocs/op
BenchmarkConvertNil-12           	486593035	         2.420 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	497484584	         2.415 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertNil-12           	494800058	         2.484 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertPlain-12         	26338915	        46.83 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	27102642	        47.43 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertPlain-12         	23882793	        46.48 ns/op	      80 B/op	       1 allocs/op
BenchmarkConvertItem-12          	182674196	         6.809 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	167730882	         6.896 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertItem-12          	182318007	         6.530 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	89580081	        13.76 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	95025154	        13.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/1-12       	83743045	        13.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	39465481	        39.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	37905969	        43.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/4-12       	32132380	        49.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	 8925120	       139.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	 7255671	       153.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/16-12      	 9320972	       124.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2362137	       569.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2100784	       500.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertDepth/64-12      	 2766693	       425.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	251357503	         4.669 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	509230066	         4.850 ns/op	       0 B/op	       0 allocs/op
BenchmarkConvertParallel-12      	424000342	         5.682 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	73821373	        15.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	71444271	        15.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkHasCode-12              	67245044	        18.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackPlain-12       	 2271877	       523.0 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 1819144	       607.3 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackPlain-12       	 2296662	       510.8 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackExisting-12    	83763364	        13.81 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	93857128	        13.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackExisting-12    	91606549	        13.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkWithStackMissing-12     	 2400548	       495.9 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackMissing-12     	 2296054	       507.7 ns/op	     112 B/op	       2 allocs/op
BenchmarkWithStackMissing-12     	 2292618	       506.4 ns/op	     112 B/op	       2 allocs/op
BenchmarkFormatVerbose-12        	  567955	      2037 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  608695	      2006 ns/op	     744 B/op	      12 allocs/op
BenchmarkFormatVerbose-12        	  598426	      1961 ns/op	     744 B/op	      12 allocs/op
BenchmarkInspect/nil-12          	100000000	        10.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkInspect/nil-12          	126348253	         9.558 ns/op	       0 B/op	       0 allocs/op
BenchmarkInspect/nil-12          	126141615	         9.511 ns/op	       0 B/op	       0 allocs/op
BenchmarkInspect/plain-12        	 2708616	       445.5 ns/op	     289 B/op	       4 allocs/op
BenchmarkInspect/plain-12        	 2728056	       441.7 ns/op	     289 B/op	       4 allocs/op
BenchmarkInspect/plain-12        	 2773762	       437.0 ns/op	     289 B/op	       4 allocs/op
BenchmarkInspect/item-12         	  502573	      2387 ns/op	    1618 B/op	      11 allocs/op
BenchmarkInspect/item-12         	  530073	      2393 ns/op	    1618 B/op	      11 allocs/op
BenchmarkInspect/item-12         	  457184	      2426 ns/op	    1618 B/op	      11 allocs/op
BenchmarkInspect/no_stack-12     	 1392027	       853.1 ns/op	     754 B/op	       7 allocs/op
BenchmarkInspect/no_stack-12     	 1388282	       855.1 ns/op	     754 B/op	       7 allocs/op
BenchmarkInspect/no_stack-12     	 1415230	       850.0 ns/op	     754 B/op	       7 allocs/op
BenchmarkInspect/branches-12     	  196423	      5709 ns/op	    4031 B/op	      26 allocs/op
BenchmarkInspect/branches-12     	  208754	      5787 ns/op	    4031 B/op	      26 allocs/op
BenchmarkInspect/branches-12     	  189387	      5639 ns/op	    4031 B/op	      26 allocs/op
BenchmarkInspect/depth_1-12      	  438422	      2782 ns/op	    1779 B/op	      13 allocs/op
BenchmarkInspect/depth_1-12      	  474236	      3071 ns/op	    1779 B/op	      13 allocs/op
BenchmarkInspect/depth_1-12      	  455428	      3552 ns/op	    1779 B/op	      13 allocs/op
BenchmarkInspect/depth_4-12      	  274756	      4362 ns/op	    2998 B/op	      18 allocs/op
BenchmarkInspect/depth_4-12      	  333284	      3657 ns/op	    2998 B/op	      18 allocs/op
BenchmarkInspect/depth_4-12      	  326593	      3997 ns/op	    2998 B/op	      18 allocs/op
BenchmarkInspect/depth_16-12     	  117068	      9874 ns/op	   12394 B/op	      39 allocs/op
BenchmarkInspect/depth_16-12     	  108675	     10832 ns/op	   12394 B/op	      39 allocs/op
BenchmarkInspect/depth_16-12     	  133016	      9602 ns/op	   12394 B/op	      39 allocs/op
BenchmarkInspect/depth_64-12     	   48714	     25149 ns/op	   29018 B/op	      93 allocs/op
BenchmarkInspect/depth_64-12     	   48273	     29101 ns/op	   29018 B/op	      93 allocs/op
BenchmarkInspect/depth_64-12     	   41150	     26262 ns/op	   29018 B/op	      93 allocs/op
BenchmarkInspect/visit_limit-12  	   14224	     86873 ns/op	  116089 B/op	     287 allocs/op
BenchmarkInspect/visit_limit-12  	   13819	     85614 ns/op	  116088 B/op	     287 allocs/op
BenchmarkInspect/visit_limit-12  	   14394	     83224 ns/op	  116089 B/op	     287 allocs/op
PASS
ok  	github.com/gtkit/goerr	103.238s
PASS
ok  	github.com/gtkit/goerr/test	0.655s
```

### 15.6 逐项验收、兼容回审与停止位置

| 验收项 | 实际状态 | 依据 / 限制 |
| --- | --- | --- |
| 3.1 最小快照 | 本阶段功能通过 | 只有 Inspect、Diagnostic、Cause；主状态与旧 Convert 一致，根全文、实际节点与分支栈保留；已验证预算、nil 与循环边界 |
| 3.2 所有权 / 失败 | 本阶段功能通过 | 篡改和并发读取在 race 下通过；零值失败不返回部分原因/帧，无旧路径调用 Inspect、无全局缓存或原对象存储 |
| 外部用法 / 文档 / 版本 | 已实施并回读 | 两个外部包 Example 在完整测试中通过；三个符号及所有字段有 GoDoc；README/CONTRIBUTING/Unreleased 与 Version=v1.5.0 对应；远端仍是查询时最高 v1.4.0、无 Release，不称已发布 |
| 旧行为 / 导出面 / 依赖 | 本轮范围已核对 | 原错误构造/包装/转换实现、error_test、模块和 go.sum 与阶段 3 前副本一致；Version 按 MINOR 更新，保留阶段 1/2 和用户 .gitignore；不改码表、Status、目录或 API/错误链；新增数据字段与第 13.1 节逐项一致 |
| build / vet / race | 本轮通过 | 最终命令分别退出 0；仅 darwin/amd64、当前模块默认构建条件，不代替跨平台或真实消费者验证 |
| lint / 现代化门禁 | 整体未通过，既有项保留 | lint 仍有原 6 项，fix 仍有旧 helper 建议；新增范围无告警，不自动处理旧缺口 |
| 覆盖率 | 整体未达门槛 | 本文件 100.0%（86/86），全仓 66.2% < 80%；不能以局部覆盖代替整体门槛 |
| 分配 / 耗时 | 分配样本已核对，耗时未验证 | 旧 17 组分配相同；27×3 样本如上，部分耗时观测超 10%，统计可信度不足且环境未控制；完整性能验收保留至获准阶段 |
| 回滚 / 发布 | 方法已写明，未执行 | 第 15.3 节要求先撤新增调用引用再按 before 副本撤阶段 3 差异；未真实回滚、提交、push、tag 或 Release；未首次升级 v2 |

调用方回审范围：当前工具未暴露 gopls/LSP 语义引用查询，使用 rg 与源码逐处核对当前单模块、默认构建条件（未发现 go:build 条件）。新生产入口仅在 diagnostic.go 声明，调用点位于 diagnostic_test.go、example_test.go、benchmark_test.go；测试、Example 与基准均检查诊断 error 后才使用快照。私有 walk/frames/reserve 的每条失败均返回到 Inspect 的零值失败路径。旧生产函数没有新调用；没有实际消费者仓库、动态/反射调用或其他平台证据，不能称完整调用图或外部接入已通过。

质量自审：diagnostic.go 为 211 行（含 GoDoc/字段说明），生产函数只有 Inspect 和 walk/frames/reserve 三个私有方法；遍历循环约 80 行，保留预算、类型保护、路径集合与入/出栈顺序在同一控制流，未为一次内部使用再拆注册器、策略、适配包或接口。已有 PC 读取与文本预算分别隔离，直接节点断言用于真实节点语义，不能改成会递归并提升虚拟 As 对象的旧查询。新增参数/字段均有契约依据，无配置和依赖版本要新增。

合规摘要：

- [x] 本轮 Go 1.27 编译/现代写法已核对；整体 fix 遗留另列。
- [x] 零外部运行依赖，目录、模块与既有错误体系不变。
- 不适用：Functional Options、UUID、JSON 编码器、服务/存储/容器迁移；本阶段没有这些入口或变更。
- [x] 新增导出 API/字段均有 GoDoc，核心入口与值类型有外部 Example。
- [ ] 整体覆盖率 ≥80%：实际 66.2%；新诊断文件为 100.0%。
- [x] Benchmark 使用 ReportAllocs/b.Loop，保留全部 81 个样本及警告。
- [x] 错误路径、失败原子性、数据所有权、并发读取与 API 兼容在本阶段范围已验证。
- [ ] 整体 lint/fix 全绿及 10% 耗时门槛：分别未通过/未验证，未勾选通过。
- [x] 新增代码职责与导出面积已回审，未知消费者/外部方法/平台限制明确，无编造结果。

交付文件：新增 diagnostic.go、diagnostic_test.go；本阶段修改 example_test.go、benchmark_test.go、version.go、README.md、CONTRIBUTING.md、CHANGELOG.md、OPTIMIZATION_PLAN.md，及按规则追加/关闭的 .harness/error-journal.md。本阶段运行代码未修改 error.go、resolve.go、status.go、errcode.go、sanitize.go 或模块依赖；它们与阶段 3 前的副本一致。

停止位置：阶段 3 的快照实现、功能验收与取证已完成，整体门禁遗留及耗时限制继续开放。阶段 4 的两个普通出口组合/日志接入说明、阶段 5、阶段 6 均未进入；不自动续跑，不提交或发布。

最终回读：13 项命令及退出码与原始 JSON 一致，计划中的 81 个原始样本与 final-bench.txt 逐行一致，基准夹具指纹一致；本轮代码及相关说明在最终验证后未再改变。回读记录保存在 delivery-readback.json。回读不代替前述未通过项，不构成下一阶段授权。


## 16. 阶段 4：可插拔出口与日志接入边界执行记录

### 16.1 执行前事实与授权

- 用户明确授权继续阶段 4，完成后写回真实结果并汇报，不自动进入阶段 5。本轮只执行既定 4.1、4.2；阶段 3 的既有 lint/fix、整体覆盖率和性能可信度缺口不因此自动修复。
- 已重新读取全局/项目 AGENTS、错误记录和原计划，以及结构、API 兼容、测试、文档、现代 Go Guide；按本次路径读取 harness-implement、use-modern-go、go-testing、go-observability 与 go-enterprise-quality 的相关正文/参考。不诊断 harness、不启动 OpenSpec、不重排阶段。
- 本轮 Go 实际入口 `/usr/local/bin/go` 为 go1.27.2 darwin/amd64；go.mod 为 1.27.0，模块仍 github.com/gtkit/goerr，无外部 require；Version 为阶段 3 准备的 v1.5.0 未发布元数据。本轮不增加生产/测试第三方依赖、不改版本或发布。
- 已按实际清单将 23 份既有源码/测试/模块/文档与用户 .gitignore 改动备份到 `/private/tmp/goerr-stage4.Ltg8SD/before`，指纹为 before-sha256.json；只在本轮授权文件写入，不用 HEAD 覆盖既有阶段或用户改动。
- logger API 参考使用已存在缓存 github.com/gtkit/logger@v1.10.0（module 为裸路径、最低 Go 1.26）与 zap@v1.28.0；logger go.mod 还锁定 logrotate v1.1.3、multierr v1.11.0，均已在本机缓存。缓存 logger.go 与本地干净 logger 仓库 v1.10.0 tag/HEAD 对应文件一致，真实命令、源码指纹见 logger-api-evidence.json。不据本地事实宣称远端最新版，不把该参考版本当调用方升级要求。
- 现有相邻项目的 logger 锁定版本并不统一；本轮不修改它们，也不从 skill 的 v2 示例推导 major 迁移。调用方应沿用自身已有路径/版本，当前 goerr 不选择 logger 运行依赖。

本轮来源耦合清单（取自已读源码，代码前固定）：

| 耦合点 | 真实契约 | 本轮对应 / 限制 |
| --- | --- | --- |
| Inspect / Diagnostic / Cause | 主状态、根全文和实际原因树；切片按每次 Inspect 独立，普通值复制仍共享切片 | 两个出口只读同一快照；不承诺任意出口写入切片时彼此隔离，不新增 Clone 或库内出口类型 |
| Go 业务错误 | Convert 选择、Error/Message/HTTPCode、errors.Is/HasCode 仍由原错误决定 | 出口错误单独处理，不回写原 error 或客户端状态；验证前后状态/文案与 cause 可达性 |
| writer | 使用标准 io.Writer，短写伴随非 nil error；可能已经写出部分字节 | 文本/JSON 返回实际写入 error，不自动重试，不宣称输出事务或部分字节回滚 |
| JSON | 测试只需零依赖公共库的标准数据互操作 | 使用 encoding/json 编码调用方私有字段结构，不编码整个 Diagnostic；README 记录标准库例外与非稳定协议边界 |
| logger.ErrorCtx | func(context.Context, string, ...zap.Field)，无返回值 | 调用方捕获自身 ctx，调用包级函数；返回 nil 只说明调用完成，不代表落盘/平台接收成功 |
| logger.WithContextFields | ContextFieldsFunc 为 func(context.Context) []zap.Field，ErrorCtx 合并已有 ctx 字段 | 沿用应用已有初始化和关联字段来源；不在 goerr 引入 context 或 provider，不伪造 request_id/trace_id API |
| zap 字段 | Int32(string,int32)、Int(string,int)、String(string,string) 返回 zap.Field | 仅临时调用方编译与 README 的已核实 API 参考；生产 goerr 不引用 zap |
| logger 生命周期 | New(opts...) error、Sync() 属调用方进程装配；ErrorCtx 为记录位置，未提供交付 ACK | 不在本库初始化/关闭日志、调整采样、栈或配置；不执行真实日志平台接入 |
| 信息边界 | 原文、公开文案与帧路径可能含内部信息；RedactForLog 只覆盖既有手机号/邮箱规则 | 示例显式选择字段/脱敏，未做通用密钥清洗承诺；日志自动 stacktrace 不冒充 Cause.Frames 起源栈 |

### 16.2 执行前实施步骤与验收（执行前已保存并回读，结果见 16.4）

| 步骤 | 前置条件 | 具体工作与文件 | 验证 / 可复核证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 4.1a 两个普通出口 | 本节保存并实际回读；不改生产接口 | 在 test/integration_test.go 的现有外部包添加两个私有输出函数，绑定 io.Writer 后类型均为 func(goerr.Diagnostic) error；文本和 JSON 各自写入内存，字段只选 operation/request_id/code/error/cause_count，Error 显式 RedactForLog | 同一 Inspect 结果切换输出变量、两个出口并存；核对选择时未写入另一个目标、独立格式内容和脱敏，快照与业务错误不变 | 新增生产接口/导出类型/注册器/配置/目录、隐式编码日志，或快照被误当客户端 DTO | 已完成；TestDiagnosticOutlets 专项及完整 race 通过 |
| 4.1b 失败边界 | 4.1a 完成 | 同一外部测试文件覆盖文本/JSON writer 在写出前失败与部分写出后失败；调用方读取失败/无错误输入时跳过出口；只在测试使用私有失败 writer，不增加生产抽象 | errors.Is 命中真实输出 error、每次只有一次写入调用、部分字节保留事实；原 error 全文、Code/HTTPCode/Message 及 cause 可达性不变；读取失败无输出且无原文回退 | 失败替换原业务错误、自动重试/二次日志、声明已回滚部分字节或缺检查诊断 error | 已完成；两种格式各两类 writer 失败、nil/超限跳过出口用例均通过 |
| 4.2a 项目日志边界 | 4.1 通过；本节 API 清单已实际核实 | README 在现有诊断与日志相邻处补普通函数出口和 gtkit/logger 包级接入说明；替换已有 log.Print/log.Printf 建议，原始 err 的日志注释改为 Inspect 同取全文/起源；只编辑本阶段直接相关日志说明 | 对照 cached logger v1.10.0/zap v1.28.0 的签名、字段和 void 返回限制回读；说明既有 ctx 字段/生命周期及敏感信息边界；不改支付/订单或贡献矩阵 | 引用未核实 API/major、返回 nil 被写成可靠交付、出现 Logger facade 或配置注入 | 已完成并回读；只给调用方普通函数参考，未加生产出口或日志依赖 |
| 4.2b 可编译参考 | 4.2a 示例已准备，引用已有缓存 | 在本轮临时目录创建隔离调用方编译工装，引用当前 goerr 工作区与既有 logger v1.10.0/zap v1.28.0 缓存，GOPROXY=off；编译 README 的同一消费者函数，不运行 logger 初始化或输出，不下载依赖、不写 goerr go.mod/go.sum | 真实 build 退出码、准确临时 go.mod/go.sum、源文件与 README 代码一致性；只能证明 Go 1.27.2 编译，不冒充日志平台接收或其他版本兼容 | 需要联网安装/升级或修改项目依赖时停止该项并报告未验证，不创建替身 SDK 假报通过 | 已完成；首次缺模块自身校验和退出 1，改用已核实真实缓存源码后 build 退出 0，代码块逐字一致 |
| 4.1/4.2 收尾 | 用例和说明完成 | 定向与完整 race、build、vet、lint、只读 fix、覆盖率各执行一次，保存实际命令/环境/退出码；检查仅本轮文件差异，写回本计划并回读 | 本轮检查结果、原 6 项 lint/旧 fix 的对比、API/依赖/版本/用户改动指纹、原始输出 | 新增失败或范围漂移；同一问题连续 3 轮不达标；不能复核证据 | 检查、写回与回读已完成；原门禁缺口未通过，结果见 16.4，阶段 5、6 未进入 |

检查使用本轮 run.py（来自已读阶段 3 记录器），缓存/产物在本轮临时目录，GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off。候选命令为 `go test -json -race -count=1 -timeout=5m -run '^TestDiagnostic' ./test`、完整 `go test -json -race -count=1 -timeout=5m ./...`、`go build ./...`、`go vet ./...`、`golangci-lint run ./...`、`go fix -diff ./...` 与 `go test -count=1 -coverprofile=<本轮目录>/coverage.out ./...`。本阶段没有生产/基准变化，不重新采集性能样本或将历史基准当本轮结果；阶段 6 的性能验收未获准。真实平台接入与外部业务工程运行不在当前授权范围。

### 16.3 兼容、风险、失败与回滚

- 文件白名单为 test/integration_test.go、README.md、OPTIMIZATION_PLAN.md，以及按规则记录本轮真实失败的 .harness/error-journal.md；临时工装只保存在 `/private/tmp/goerr-stage4.Ltg8SD`。所有生产 .go、模块/锁文件、Version、既有 Example/基准、CONTRIBUTING/CHANGELOG、规则与 OpenSpec 均不改；保留用户与前序阶段改动。
- 调用方按用途读取快照；共享同一个快照时，两个出口须只读。如出口需要修改切片，调用方自己安排独立数据，不给本库加克隆/保护开关。
- JSON 的示例字段属于调用方私有日志记录，不成为 goerr 的稳定 wire/schema；不把 PC、runtime.Frame.Func 或整个快照隐式编码。现有 encoding/json 仅用于零依赖公共库的测试互操作，不改变调用方已锁定 gtkit/json 约定。
- 记录操作失败可留下部分输出，不能撤销已写字节；不自动重试，避免重复事件。读取失败不退回未经处理根原文；业务响应仍由原 error 的 Convert/SanitizeForClient 决定。示例只说明边界，不实现实际 handler 或支付 ACK 策略。
- logger v1.10.0 是已核实本地参考，不是自动选型/升级；ErrorCtx 无返回 error，无法提供逐条交付 ACK。已有用户级日志配置、接收字段、真实平台连通性与其他版本兼容继续未验证，不用内存 writer 或编译替代这些证据。
- 回滚本阶段时只撤新增测试与 README 对应差异，参照 before 副本保留全部生产实现和用户改动；临时工装可停用，无服务/存储/依赖迁移、注册或后台资源要回收。本轮不实际回滚、不提交、推送、tag 或发布。
- 当前阶段固定为 4；完成后逐项记录真实结果并回读，不进入阶段 5 的支付说明、订单示例、贡献矩阵/未发布记录集中修订，也不进入阶段 6。

### 16.4 实际结果（2026-10-21）

#### 16.4.1 环境、命令与原始证据

本轮检查继续使用 Go 1.27.2 darwin/amd64。项目工作目录为 `/Users/xiaozhaofu/go/src/my-gtkit-package/goerr`；临时消费者工作目录为 `/private/tmp/goerr-stage4.Ltg8SD/logger-consumer`。所有下表中的实际子进程均经 run.py 单独记录命令、cwd、开始/结束 UTC、退出码和输出，文件为本轮证据目录内同名 .json / .txt；汇总见 results.json。

统一覆盖环境为 `GOTOOLCHAIN=local`、`GOCACHE=/private/tmp/goerr-stage4.Ltg8SD/gocache`、`GOLANGCI_LINT_CACHE=/private/tmp/goerr-stage4.Ltg8SD/lint-cache`、`GOFLAGS=-mod=readonly`、`GOPROXY=off`。未下载、安装或升级依赖，不更改用户级配置。

| 证据名 | cwd | 实际命令 | 退出码 | 真实结果 / 限制 |
| --- | --- | --- | --- | --- |
| targeted-race | 项目 | `/usr/local/bin/go test -json -race -count=1 -timeout=5m -run '^TestDiagnostic' ./test` | 0 | 出口专项竞态：3 个顶层测试，9 条含子测试 pass 事件，1 个包通过，无 fail/skip |
| logger-consumer-build | 临时消费者 | `/usr/local/bin/go build ./...` | 1 | 首次工装失败：缺 logger 模块自身 go.sum 项，不是 API 编译错误 |
| logger-consumer-build-corrected | 临时消费者 | `/usr/local/bin/go build ./...` | 0 | 改用真实缓存源码后编译通过，无输出；只编译、不运行 |
| full-race | 项目 | `/usr/local/bin/go test -json -race -count=1 -timeout=5m ./...` | 0 | 67 条顶层 pass 事件、132 条含子测试/Example 的 pass 事件，另 2 个包通过，无 fail/skip/race 报告 |
| build | 项目 | `/usr/local/bin/go build ./...` | 0 | 通过，无输出 |
| vet | 项目 | `/usr/local/bin/go vet ./...` | 0 | 通过，无输出 |
| lint-version | 项目 | `/Users/xiaozhaofu/go/bin/golangci-lint version` | 0 | 2.14.0，使用 go1.27.2 构建 |
| lint-corrected | 项目 | `/Users/xiaozhaofu/go/bin/golangci-lint run ./...` | 1 | 未通过：原有 errcheck 5 项、ineffassign 1 项；与阶段 3 输出逐字一致 |
| fix | 项目 | `/usr/local/bin/go fix -diff ./...` | 1 | 未通过：仅原 error_test.go stackFrames 的 Split→SplitSeq 建议；与阶段 3 输出逐字一致，未应用 |
| coverage | 项目 | `/usr/local/bin/go test -count=1 -coverprofile=/private/tmp/goerr-stage4.Ltg8SD/coverage.out ./...` | 0 | 测试通过；根包 66.2%，未达 80%，test 子包无可统计语句 |

另有一次 lint 启动错误：记录器尝试 `/usr/local/bin/golangci-lint run ./...`，报 `FileNotFoundError`，工具会话 42905 退出 1，**没有启动 lint 子进程，也没有该次 lint 的退出码或 run.py 元数据**。先记录 ERR-20261021-016，再用 `command -v golangci-lint` 核实实际入口 `/Users/xiaozhaofu/go/bin/golangci-lint`，退出 0；版本与实际 lint 结果列于上表。原始异常见当次工具输出，事实摘要保存为 lint-launch-failure.json，不能把启动失败或空 stdout 当作 lint 通过。

本阶段没有生产/基准修改，因此没有重新运行 benchmark、性能比较或 Fuzz；第 15 节的耗时可信度和预算限制继续未验证，不将历史样本充当本轮性能结果。未连接真实业务 HTTP 服务或日志平台，没有平台接收、持久化或上线证据。

#### 16.4.2 步骤 4.1：普通函数出口与失败边界

- `test/integration_test.go` 仍是从外部 import goerr 的现有包。新增两个私有普通函数，分别绑定文本和 JSON writer，返回同一种未命名函数类型 `func(goerr.Diagnostic) error`。没有新增生产类型、接口、注册器、开关或目录。
- `TestDiagnosticOutlets` 用一次 Inspect 的快照切换函数变量，再让两个出口并存；文本目标和 JSON 目标分别核对内容、事件数及选择隔离。两者显式选取 operation/request_id/code/error/cause_count，手机号/邮箱的合成夹具在出口脱敏；JSON 不包含整个快照或运行时元数据。输出后快照与另一份独立 Inspect 结果一致，根全文、主 Item 和原 cause 可达性保持。
- `TestDiagnosticOutletFailures` 覆盖文本/JSON 各自写出前失败、写出 7 字节后失败，共 4 个子用例。调用方得到 `io.ErrClosedPipe`；每次 writer 只调用 1 次，部分字节仍然存在，未做重试或假报回滚。原错误全文、Item 身份、Code/HTTPCode/Message、errors.Is 和 HasCode 均保持，输出错误没有进入业务错误链。
- `TestDiagnosticFailureSkipsOutputs` 覆盖 nil 与超过诊断文字预算的输入；均不调用出口。超限返回诊断错误且不输出原文，原业务错误仍可用于原有状态选择；不把诊断失败当作业务成功。
- 专项 race 的 3 个顶层测试及其 6 个子测试均通过；完整 race 的两个包通过。这是内存 writer 和响应字段投影的契约证据，不冒充真实 HTTP 响应或日志接收端验收。

#### 16.4.3 步骤 4.2：已核实日志参考与边界

- README 在既有诊断相邻位置说明普通函数出口、只读共享、调用方私有 JSON 字段、诊断与输出失败分离、部分写入和脱敏限制；替换原 log.Print/log.Printf 建议及原始外层错误 `%+v` 的日志假设。支付/订单、渠道重试、贡献矩阵与未发布记录集中修订均未进入。
- API 来自真实本机缓存 `github.com/gtkit/logger@v1.10.0`、`go.uber.org/zap@v1.28.0`，并按 16.1 清单读取源码。`logger.ErrorCtx(context.Context, string, ...zap.Field)` 无返回值；示例仅返回“调用完成”的 nil，没有可靠交付承诺。logger 的既有 ContextFields、初始化/关闭、采样与字段策略均由调用方掌握，不向 goerr 注入日志对象、context 或 provider。
- 临时消费者只把 README 中 `package app` / `loggerOutput` 代码块逐字提取为 logger_output.go，引用当前 goerr 工作区；编译中不执行 logger.New、logger.Sync、ErrorCtx 或应用初始化。首次 build 缺 logger 自身校验和退出 1 后，先记录 ERR-20261021-015，再将临时模块显式 replace 到已核实的真实 v1.10.0 缓存源码，修正后 build 退出 0。没有创建替身 SDK，也没有修改 goerr go.mod/go.sum。
- 临时 go.mod 锁定 goerr v1.5.0（replace 到当前未发布工作区）、logger v1.10.0（replace 到实际缓存源码）、zap v1.28.0，间接 logrotate v1.1.3、multierr v1.11.0；go 1.27.0。完整 go.mod/go.sum 与源码、首次失败和修正后输出均留在本轮目录，最终指纹见 verification-summary.json。只验证该真实源码组合在 Go 1.27.2 可编译，不声称完成远端来源、校验数据库或供应链审计。
- README 代码块与编译源逐字一致，实际编译源 SHA-256 为 `5921e28cd8459e26ab80fd849ed9e078a950e92643b270284754d99d433cb2a9`。这不是最新版声明、调用方依赖选择或 major 升级；其他锁定版本兼容与真实日志平台接收仍未验证。自动日志 stacktrace 与 Cause.Frames 起源栈明确区分；内部路径、业务秘密和日志关联字段仍按调用方策略处理。

#### 16.4.4 门禁、导出面与范围核对

| 项目 | 本轮实际结果 | 结论 |
| --- | --- | --- |
| Go 1.27 / 现代写法 | 当前测试使用 range-over-int、min；编译/竞态/vet 通过；fix 仍有原测试 helper 建议 | 本轮文件无新增 fix 建议，整体 fix 门禁未通过 |
| 外部依赖、API 与配置 | 生产 .go、go.mod/go.sum、Version 指纹均与阶段 4 开始时相同 | 导出面、运行依赖与配置未增加；无日志初始化或隐式输出 |
| 错误、并发与数据归属 | 输出失败与业务错误分离；独立并行夹具、只读消费；完整 race 无报告 | 本轮出口契约通过，不扩大为任意外部代码的线程安全保证 |
| lint | 只有 error.go:84/90/101/103/105 的 errcheck 与 resolve_test.go:24 的 ineffassign | 未通过；6 项均与阶段 3 原始输出逐字一致，无新增出口问题 |
| 语句覆盖率 | 从本轮 coverage.out 计数为 251/379，66.2%；其中 diagnostic.go 为 86/86，100% | 全包低于 80%，不通过覆盖门槛；外部 test 子包无可统计生产语句 |
| GoDoc / Example / Options / 基准 | 没有新增生产 API、配置或行为；既有 Example 在完整测试中通过 | 不额外扩导出面或添加 Options/基准；未采集本轮性能证据 |
| 源码与文档自审 | 逐行检查出口选择、显式投影/脱敏、单次写入及失败路径；README 对照真实 API 与编译源回读 | 未发现本轮未处置的契约问题；不借自审修无关旧问题 |
| 文件保护 | before-sha256.json 的 23 项仅 README、计划、test/integration_test.go 变化，其余 20 项不变 | 保留原用户 .gitignore、前序生产实现、既有测试/基准与文档；错误记录按规则追加 |
| 文档与日志边界 | 编译源码一致、无 std log/slog 接入、无 Logger facade；说明私有 JSON 投影和 std JSON 测试例外 | 必要说明已完成；没有用编译或内存 writer 替代真实平台验证 |

verification-summary.json 保存受保护文件、README/编译源码及与阶段 3 lint/fix 输出比较结果，相关比较均为 true。coverage 的百分比由本轮 profile 的已覆盖/总语句计算（251 ÷ 379 × 100，四舍五入至一位小数），并与真实命令输出 66.2% 对照；不是新增外部测试代码自身的覆盖率。

#### 16.4.5 本轮纠正、风险与回滚状态

- ERR-20261021-014：README 初稿提前写“已编译”，先记录并改回待验证；取得并读取修正后 build 退出 0 及代码一致性证据后才恢复实测结论，已关闭。
- ERR-20261021-015：临时消费者自身模块校验和缺失，首次编译退出 1；按真实缓存源码修正临时模块后退出 0，已关闭。首轮失败保留，未做网络安装或项目依赖变更。
- ERR-20261021-016：lint 工具入口猜错导致记录器启动失败；按核实入口实际执行，结果仅原有 6 项，已关闭。历史 ERR-20261021-003/004/008 的既有 lint/fix、覆盖率缺口继续开放。
- 回滚方法仍为 16.3：仅撤本阶段新增外部测试及 README 对应差异，参照 before 副本保留阶段 1～3 与用户改动；不使用 HEAD/整文件重置覆盖既有工作。没有运行依赖、服务、注册状态或存储迁移需要回退。**本轮未实际执行回滚**。
- 真实消费方锁定版本、初始化配置、关联字段、平台接收、内部路径/业务秘密策略和性能预算仍未验证；没有发布、提交、推送、tag 或修改 OpenSpec。

#### 16.4.6 交付与停止位置

阶段 4 的 4.1、4.2 功能与必要日志边界说明已完成；整体质量门禁仍按上文保留未通过项，不扩大为全项目验收已通过。

已实际回读本轮计划的授权/阶段状态、16.1～16.4、README 全部改动区段及 test/integration_test.go 全文。plan-readback.json 的 10 行命令与原始元数据比较均为 true；生产与其他受保护文件不变、阶段 1～3 历史记录不变、阶段 5/6 原计划不变、README 仅直接相关日志区段变化、代码块与编译源一致均为 true。错误记录的 014/015/016 已关闭，003/004/008 仍开放。文档回读未触发额外测试、代码修改或阶段推进。

停止在阶段 4，阶段 5、6 未授权、未进入，不自动继续；未提交或发布。
