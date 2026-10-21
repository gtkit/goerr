# goerr 优化方案与分阶段实施计划

日期：2026-10-21（Asia/Taipei）  
状态：功能与既有门禁续项已完成；用户现已明确授权提交、发版并推送 GitHub。当前执行第 21 节发布方案，先落盘回读，再核对发布元数据、完成必要检查、提交/推送和创建 Release；结果待实际执行。
当前授权与执行边界以第 21 节为准；第 13～20 节保留当时范围、状态和真实证据。
需求来源：前序 oops 对照分析、阶段 0～6 及门禁修复；本轮明确要求先补排查与复核方案并回读，再区分真实性能回退与环境波动，验证并写回真实结果。

## 1. 本次授权、规则与修订记录

- 用户本轮授权排查与复核剩余 7 组性能未确认项，第 20 节是补充续项，不重排原阶段 0～6，不追写前轮授权或结果。
- 先保存并回读完整排查方案，再执行固定的源码/构建核对、A/A 对照、交错基准和统计分析；不修改运行代码、仓库测试、依赖、配置、导出面或性能门槛。
- 本轮只写本计划及真实错误的处置记录；临时模块、二进制、采样和分析工具均放证据目录。已有十份工作区差异全部保留。
- 不筛选样本、不择时等低负载、不关闭其他进程、不修改电源/CPU/GC 全局配置，不启动 OpenSpec，不安装工具，不提交、推送、tag、发布或首次 v2。
- 第 13～19 节的历史证据保留；未执行、失败、未确认和通过分别记录。测量结束后只写回并回读，不顺带重跑非本次质量检查。

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
| 本轮计划 | OPTIMIZATION_PLAN.md | 原阶段 0～6 与各续项的授权、契约、真实结果和遗留项；本轮性能状态见第 20 节 |
| 阶段 1 测试 | error_test.go、拟新增 benchmark_test.go | 原实现基准、旧行为契约与预期失败反例；不引用尚不存在的 API |
| 阶段 1 错误记录 | .harness/error-journal.md | 仅追加本轮真实失败及处置，遵守 AGENTS.md，不改规则或示例历史 |
| 阶段 2 执行范围 | error.go、error_test.go、go.mod | WithStack 补栈修复、错误链回归和 Go 1.27.0 声明；仅在第 14 节授权范围执行 |
| 阶段 3 已新增 | diagnostic.go、diagnostic_test.go | 读取入口、必要值类型、私有有界遍历及契约测试；未拆目录或扩展框架 |
| 阶段 3/4/6 测试 | example_test.go、benchmark_test.go、test/integration_test.go、sanitize_test.go | 阶段 3/4 已完成 Example/基准与普通出口；阶段 6 只修改后两个文件，补齐截断性质与共享快照并发出口，既有 Example/基准不改 |
| 阶段 6 结果文档 | CONTRIBUTING.md、OPTIMIZATION_PLAN.md、.harness/error-journal.md | 同步实际测试状态、结果和真实失败；其他生产源码、依赖、README/CHANGELOG、规则保持阶段 6 开始时的内容 |
| 门禁修复续项 | error.go、error_test.go、resolve_test.go、status_test.go 及必要文档/错误记录 | 仅 Format 错误处理和对应 lint/fix/公开契约覆盖；完整范围与实际结果见第 19 节 |
| 阶段 5 指定文档 | README.md、CHANGELOG.md、CONTRIBUTING.md | 集中同步既有能力、渠道/订单示例、贡献矩阵与 Unreleased；只编辑并回读，不执行文档中的命令 |
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

保留阶段 0～6；所有阶段的获准执行与取证已结束，整体门禁未通过。各阶段真实结果、未通过与未验证项分别保留，不自动扩展为发布或额外修复。

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

### 阶段 5：说明、示例和契约同步（已完成，仅文档同步与回读）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 5.1 | 相关能力已验收且阶段 5 获准 | 集中核对 Go 1.27、最小 API、主状态/原因区别、内部快照与客户端响应、logger 与 JSON 边界 | 与实际实现和规则逐项回读 | 候选能力被写成已实现或需修改代码 | 已完成指定说明同步和回读；见第 17 节，非本轮运行验证 |
| 5.2 | 5.1 完成且渠道原文已核对 | 纠正支付重试说明与订单过期示例，不改码表和 HTTP 映射 | 原文与说明回读 | 未核实渠道假设或擅自改变业务状态 | 已按具体官方协议和既有码表修正文档并回读；未改实现 |
| 5.3 | 5.1～5.2 完成 | 同步贡献矩阵和未发布记录，核对已实施阶段的 MINOR 元数据，不改历史版本号/日期 | 文档与版本记录回读 | 文档任务扩展为测试、发布或规则维护 | 已同步 CONTRIBUTING 与 Unreleased 并回读；v1.5.0 未发布，未运行检查 |

阶段验收：只修改指定说明并回读，不运行任何测试、构建或 OpenSpec。整体运行验收归阶段 6。

### 阶段 6：集成验收、性能评估与交付（执行取证结束，整体门禁未通过）

| 步骤 | 前置条件 | 具体工作 | 检查与证据 | 停止条件 | 交付物 / 状态 |
| --- | --- | --- | --- | --- | --- |
| 6.1 | 阶段 1～5 已验收且阶段 6 获准 | 补齐外部包、并发和既有截断性质用例，按最新代码规则执行必要检查 | 第 10 节矩阵、实际 Go 1.27 环境与真实结果 | 兼容/数据归属/导出面/日志规则未达标 | 已补齐性质/并发出口测试；完整 race/build/vet/Fuzz 通过，旧门禁未通过；见第 18 节 |
| 6.2 | 功能结果完整且基准环境可比较 | 测原路径与快照读取成本，不自动实施缓存或扩展配置 | 原始样本、allocs/op、B/op、中位耗时和差异原因 | 超预算未解决或环境不具可比性 | 440 个真实样本与 benchstat 已保存；正常路径分配/样本中位数达标，稳定耗时门槛未确认 |
| 6.3 | 前述结果可复核 | 核对文件、导出清单、版本、限制、接入与撤回步骤，更新本计划状态 | 每项实际证据与回滚定位 | 存在阻塞问题或拟自动提交/发布 | 导出/依赖/范围核对及结果写回已完成；已回读核对 |

阶段验收：必要条件均有真实证据，未验证范围明确；不包含消费者仓库修改、部署、提交、推送、tag、Release 或首次 v2。

## 10. 验收矩阵与未来检查

本矩阵是完整方案验收目标；阶段 1～5 的历史结果见第 13～17 节，原阶段 6 的逐项结果见 18.4，本轮门禁修复与剩余性能未确认项见 19.5，不以执行结束替代验收通过。

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
| 2026-10-21 | 阶段 5 | 只编辑 README、CHANGELOG、CONTRIBUTING 与本计划；对照源码和官方具体渠道文档同步说明、示例、契约及 Unreleased；真实文档结果见第 17 节 | 指定四份文档同步与回读已完成；未改代码、未执行运行检查；阶段 6 未授权、未进入 |
| 2026-10-21 | 阶段 6 | 仅补截断性质与共享快照并发出口；完整 race/build/vet/Fuzz 通过；440 个基准样本、AST/引用/依赖核对与真实命令见第 18 节 | 执行取证结束，结果已写回，已回读核对；整体门禁未通过，稳定耗时门槛及外部消费者/平台未验证；未改生产源码、依赖、README/CHANGELOG，未提交或发布 |
| 2026-10-21 | 门禁修复续项 | 先补方案并回读；修 Format 五处 error、测试无效赋值与 SplitSeq；补 136 个状态及 14 个 Is 契约；完整质量检查、100.0% 覆盖率和 780 个样本见第 19 节 | 代码质量门禁通过；7 组稳定耗时门槛未确认。结果已写回并回读核对；未改 API/依赖/Version，未提交或发布 |

阶段 0 未写错误记录或其他规则。本轮 `.gitignore` 的既有改动保持原样；阶段 1 仅按 AGENTS.md 追加真实失败记录，不改规则。各阶段记录授权范围、实际基线、文件与导出差异、命令/环境/退出码、逐项结果、失败与未验证项、回滚定位和停止位置。

## 12. 依据与未确认事项

- 当前约束依据：根目录 AGENTS.md、全局 AGENTS.md，以及已读取的 pkg-structure、pkg-api-compat、pkg-docs Guide；文档任务不触发代码门禁或工作流。
- 前序 oops 核心源码分析仅作设计参考；其 main 未固定 commit，不据此承诺发布版本能力，也不继承其全部接口。
- [Go errors 文档](https://pkg.go.dev/errors)用于错误树背景；前序已读的 Go 1.26 fmt 源码只解释当前现象，不作为 Go 1.27 已验证证据。后续代码阶段以实际 Go 1.27 源码/工具链重新核对相关行为。
- [微信支付回调通知注意事项](https://pay.wechatpay.cn/doc/v3/merchant/4012075420)为前序已核对的渠道依据；阶段 5 未能直接抓取该页面，本轮实际依据为[微信支付 APP 支付成功回调通知](https://pay.wechatpay.cn/doc/v3/merchant/4013070368)的应答与重复通知规则，来源和读取限制见 17.1，不外推其他渠道。
- 未确认：实际消费者与旧工具链需求、消费方 gtkit/logger/JSON 的锁定路径版本与平台字段、接入人员和时间、性能绝对预算，以及第 13 节单列的环境限制。阶段 4 的 logger v1.10.0 / zap v1.28.0 仅为已核实并编译的本地 API 参考，不能据此宣称其他消费者已兼容、平台已接入或供应链验证已完成。
- 阶段 1 未修改运行代码；阶段 2 仅执行补栈、最低 Go 声明和对应回归/说明，不升级已装工具链或增加运行依赖；阶段 3 已实现与验证最小快照，真实结果与未通过项见第 15 节；阶段 4 已完成普通出口和日志边界说明，未改生产源码/API/依赖；阶段 5 只同步指定文档；阶段 6 本轮补测试并完成取证，生产源码/依赖未改，真实结果及整体门禁缺口见第 18 节。

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

## 17. 阶段 5：说明、示例和契约同步

### 17.1 本轮授权、依据与范围

- 用户仅授权既定 5.1～5.3，白名单为 README.md、CHANGELOG.md、CONTRIBUTING.md、OPTIMIZATION_PLAN.md；只编辑和回读文档，不修改源码、测试、模块/配置、规则、错误记录或 OpenSpec，不运行测试、构建、vet、lint、fix、覆盖率、Fuzz、基准或日志示例编译，不自动进入阶段 6。
- 已实际读取最新全局/项目 AGENTS、完整错误记录、原计划、三份目标文档及 pkg-docs/pkg-api-compat Guide。对照读取当前 go.mod、version.go、诊断实现及相关旧 API、错误码/预构建状态和测试/Example/基准声明；这些是文档依据，不是本轮运行验证。
- 当前静态元数据仍为 `github.com/gtkit/goerr`、Go 1.27.0、Version v1.5.0，尚未发布；本轮不重新查询 tag/Release，不调整版本或历史版本日期。阶段 4 的 SDK 编译和代码检查仅作为历史证据引用。
- 渠道核对：原链接 `https://pay.wechatpay.cn/doc/v3/merchant/4012075420` 两次直接抓取均超时，`4012084431` 页面也未能直接读取，不能宣称这些页面本轮已核实。已从限定官方域名的检索结果读取[微信支付 APP 支付成功回调通知](https://pay.wechatpay.cn/doc/v3/merchant/4013070368)的 2.2 应答与 2.3 重复通知机制；本轮修订只采用这一具体协议依据，不推导支付宝或所有渠道的重试规则。网页直接读取失败不影响已有官方原文提供的这一项文档依据。

### 17.2 执行步骤、兼容与回读标准

| 步骤 | 前置条件与指定改动 | 回读验收与停止条件 | 当前状态 |
| --- | --- | --- | --- |
| 5.1 | 对照实际源码和阶段 1～4 记录，集中核对 README 的 Go 1.27、Inspect/Diagnostic/Cause、主状态与原因树、数据归属、客户端响应及 logger/JSON 边界 | 说明与既有实现一致；不新增 API、配置或协议承诺，需要修改源码时停止对应项 | 已同步并回读；原有准确说明保留，仅补最小导出面及首个 Item 查询边界 |
| 5.2 | 官方具体渠道原文已取得；纠正支付重试概括，订单过期使用现有 StatusOrderExpired/ErrOrderExpired | 不把业务 HTTP 映射当渠道 ACK，不外推其他渠道；不改码表与 HTTP 映射，不虚构接入代码或已验证运行结果 | 已修正并回读；依据为实际源码与官方 APP 支付通知协议，旧链接本轮未读到的限制保留 |
| 5.3 | 5.1～5.2 完成后同步 CONTRIBUTING 矩阵/后续检查说明与 CHANGELOG Unreleased，核对已有 v1.5.0 MINOR 元数据 | 历史结果、当前文档回读及待验收目标分开；不改历史版本区段或源码版本，不执行文档中的命令，不发版 | 已同步并回读；未发布元数据一致，未运行任何文档内代码检查命令 |

兼容约束：仅校正文档，保持所有现有签名、字段、状态码、HTTP 映射、依赖和错误行为；保留调用方普通函数出口，不增加库内 Logger 或插件框架。风险是把协议举例外推、把历史检查写成本轮结果、把诊断快照写成响应协议；回读逐项排除这些表述。撤回时仅逆转本轮四份文档中的对应差异，保留用户与阶段 1～4 改动，不重置整个工作区；本轮不实际执行回滚。

### 17.3 真实结果与停止位置

| 交付项 | 本轮真实修改 / 回读结果 | 证据与限制 |
| --- | --- | --- |
| README 契约 | 保留 Go 1.27.0、已有读取预算/失败/所有权及 logger/JSON 边界；补充诊断仅三个新增导出符号、跨包带键初始化和 HasCode 首个 Item 语义 | 已对照 diagnostic.go、error.go、go.mod、version.go 及前序记录回读；没有新增库内接口、选项或运行依赖 |
| 订单示例 | 改用 StatusOrderExpired()，明确 ErrOrderExpired=10010707、10010701=ErrOrderNotFound | 对照 errcode.go 与 status.go 的真实声明；只改 README 片段及 Unreleased 文档更正，未修改源码或执行示例 |
| 支付说明 | 用具体渠道协议取代“4xx 停止重试”概括，移除直接将业务 HTTP/码拼作通用 ACK 的示意 handler；保留业务码分类与显式脱敏出口 | 官方 APP 支付通知 2.2/2.3 节提供依据，README 与 Unreleased 均引用同一页面；不声称原超时页面已读取，不外推支付宝或其他通知接口，不改业务 HTTP 映射 |
| CONTRIBUTING | 同步旧 API/主状态、原因树、出口成功/失败、Example、基准矩阵；当前 Fuzz 只有调用、没有性质断言的缺口单列待补充 | 对照真实测试/Example/基准声明回读；命令仅供后续授权代码阶段使用，纯文档明确只回读；矩阵不表示本轮执行或全部门禁通过 |
| CHANGELOG / MINOR | 只补 [Unreleased] 的既有诊断契约、出口说明、文档更正及贡献说明；继续标注 v1.5.0 尚未发布 | 与 README、CONTRIBUTING 及 version.go/go.mod 静态内容一致；v1.4.0 及更早区段、历史版本号/日期均未改；旧支付描述由 Unreleased 明确更正 |
| 读取与编辑 | 本轮必要文件读取及修改区段回读命令均退出 0；apply_patch 只写白名单四份文档 | 渠道直接抓取的失败已在 17.1 区分记录；没有新增临时模块、检查产物或错误日志条目，不把网页失败当作项目运行失败 |
| 未执行 / 未验证 | 测试、构建、vet、lint、fix、覆盖率、Fuzz、基准、SDK 编译、真实支付通知/日志平台、回滚、提交和发布均未执行 | 不将阶段 4 的历史 race/build/vet、66.2% 覆盖率或 SDK 编译冒充本轮执行结果；既有 6 项 lint、旧 fix 建议、覆盖率与性能可信度缺口未解决 |

已实际回读 README 改动区段、CONTRIBUTING 全文、CHANGELOG 的 Unreleased 以及本计划的当前授权、阶段 5 步骤、记录与第 17 节，相关读取命令退出 0。对本轮写入前后文档内容的回读比较确认：CHANGELOG 从 v1.4.0 开始的历史区段逐字不变，README 的既有可插拔出口/logger 参考代码与边界说明逐字不变；两个比较结果均为 true，不是运行测试。版本继续为 v1.5.0 未发布，Go 最低声明继续为 1.27.0。

阶段 5 的 5.1～5.3 已完成；本轮只修改指定四份文档并回读，没有代码/测试/构建或其他运行检查，没有提交、发布或实际回滚。停止在阶段 5，阶段 6 未授权、未进入。

## 18. 阶段 6：集成验收、性能评估与交付

### 18.1 执行前事实、授权与兼容约束

- 用户要求继续后续阶段，按原编号执行 6.1～6.3。阶段 1～5 的功能/文档范围已完成，不等于 lint、fix、覆盖率和性能门槛已全部通过；本轮重新取证并如实记录。
- 已实际重读全局/项目规则、完整错误记录、原计划和相关源码/测试；读取 harness-implement、use-modern-go、go-testing、go-performance、go-concurrency、go-error-handling、go-enterprise-quality 及对应参考/项目 Guide。仅按本阶段路径使用，不启动 OpenSpec 或工具安装诊断。
- 当前 HEAD 为 `f6a30726b372b3dc07feda42e161d27de6cb5d77`，前序实现已成为已提交基线；进入本轮前工作区只有 CHANGELOG、CONTRIBUTING、OPTIMIZATION_PLAN、README 的现有修改。原实现参考 commit 为 `d1f31a90c7080fffce42fc591ca4003a744a0eef`。两者由真实 git log/rev-parse 核实，不重置工作区。
- 当前实际 `/usr/local/bin/go` 为 go1.27.2 darwin/amd64，GOAMD64=v1、CGO_ENABLED=1、GOWORK 为空；模块仍 `github.com/gtkit/goerr`、go 1.27.0、Version v1.5.0 开发元数据，无外部 require。golangci-lint 实际入口为 `/Users/xiaozhaofu/go/bin/golangci-lint`，2.14.0/go1.27.2；benchstat 实际入口为 `/Users/xiaozhaofu/go/bin/benchstat`，x/perf v0.0.0-20260908200009-22c9c6c9d4da（工具由 go1.27.1 构建）。不升级工具或依赖。
- 25 份当前源码/测试/模块/文档/规则和错误记录已备份到 `/private/tmp/goerr-stage6.iourzyn1/before`，指纹见 before-sha256.json；preflight.json 保存真实 revision 命令与退出码。先保存并回读本节，再修改测试。
- 项目写入白名单：sanitize_test.go、test/integration_test.go、CONTRIBUTING.md、OPTIMIZATION_PLAN.md，以及按规则追加/处置本轮真实失败的 .harness/error-journal.md。所有生产 .go、既有基准、go.mod/go.sum、Version、README/CHANGELOG 及其现有改动、AGENTS、规则、目录架构保持不变。测试仅用标准库，不新增 API、Options、日志接口或编码依赖。
- 只补原 M5/M6/M8/M9 所需的集成与性质验证；不借总体验收修改原 Format 实现、无关 ResolveStatus 测试、旧 stackFrames helper 或批量补全状态 getter 覆盖率。若这些既有缺口仍失败，保留为最终验收未通过项；不得为凑门槛扩大实施范围。临时分析工具/基准模块只写本轮证据目录。

### 18.2 实施步骤、证据与停止条件（已执行，真实结果见 18.4）

| 步骤 | 前置条件 | 具体改动或执行 | 回读/验收证据 | 停止条件与交付状态 |
| --- | --- | --- | --- | --- |
| 6.1a 截断契约 | 本节已保存并回读；sanitize.go 语义已读取 | sanitize_test.go 的 ClampMessage 表格加命名、空串/中文/emoji/组合字符/省略号、零负阈值、极大极小 int、合法/非法 UTF-8 边界；Fuzz 取消把阈值强制变非负的归一化，保留真实输入，增加种子与性质断言 | 不截断时原字节保持；截断时原字节前缀、rune 数和省略号语义；合法 UTF-8 仍合法、非法字节不被隐式转换；重复截断幂等；不复制生产 range 算法当 oracle | 若反例指向需修改运行行为，先记录并停在对应失败处，不削弱断言；已完成，15 个边界用例、9 个 Fuzz 种子及性质断言通过 |
| 6.1b 外部并发出口 | 6.1a 已保存，沿用阶段 4 两个私有普通输出函数 | test/integration_test.go 增加有界 16 个 worker 同时只读同一 Join/标准包装快照，各自 writer，使用 WaitGroup.Go 并等待全部退出；核对主状态、原因结构、显式字段/脱敏和原错误不变 | 两种格式真实内存结果、无共享 writer、快照/根全文/原 cause/Convert/HasCode 保持；完整 race 覆盖现有独立数据用例与新集成用例 | 不向生产加同步/框架，不在 worker 调用 Fatal、不使用 sleep；已完成，16 个 worker 只读共享快照的集成 race 通过 |
| 6.1c 功能与质量验证 | 上述测试完成；工具路径已核实 | 专项 race 后执行完整 go test -json -race -shuffle=on -count=1 -timeout=5m ./...、build、vet、lint、只读 fix、覆盖率及 cover -func；最后单独跑根包 FuzzClampMessage 30s，最多 4 个 worker | 原始输出、命令、环境、退出码及 test/pass/fail/skip 事件；对照 M1～M9；覆盖率与工具退出分别记录，保留旧失败 | 新失败先记录再定位；同一问题连续 3 轮不达标停止，不自动修既有实现；已执行，race/build/vet/Fuzz 通过，lint/fix/覆盖率未达标 |
| 6.2a 同条件基准准备 | 功能验证完成，真实 Git 源码和保存夹具均可读 | 临时 old 模块只放原 commit 的生产 Go/模块文件，临时 new 模块放当前生产文件；共同使用已读阶段 1 final-benchmark_test.go，函数体不改。临时 old 的 go 声明仅对齐 1.27.0，当前项目不改；两者都用 -trimpath 去掉临时路径影响 | fixture 与当前共同部分一致性、源码来源/指纹、完整临时 go.mod、同一 Go/环境/参数；不引用历史计时充作本轮基线 | 源码/夹具不一致、需要依赖或无法核实环境时标未验证，不生成替身实现；已完成，共同夹具及 old/new 源码/模块指纹已核对 |
| 6.2b 旧路径与快照成本 | 临时源码和夹具已核对，其余重任务结束 | old/new 各 10 轮交替采样，奇偶轮交换先后；共同 17 组基准使用 -run '^$' -bench . -benchmem -benchtime=100ms -count=1 -trimpath。当前工作区另采 BenchmarkInspect 的 10 组、每组 10 次；保存 benchstat 比较、原始样本、median/B/op/allocs 与不确定性 | 正常旧路径分配不增加及同条件中位耗时回退 ≤10% 分别判断；WithStack 缺栈新增捕获成本单列；深度 1/4/16/64、分支、节点预算、有栈/无栈/重复读取全有样本 | 超预算或波动不能解释即保留不通过/未验证，不自动做缓存/性能改造；最多按统计警告补 1 轮必要样本，不无限重测；已完成 440 个样本，正常路径分配/样本中位数达标，稳定耗时门槛未确认 |
| 6.3 交付核对 | 所有可做结果已保存，失败与未验证已分开 | 临时标准库 AST 工具比较原 commit 与当前导出声明，允许只有既定三个诊断符号和 Version 元数据差异；核对本轮生产文件不变、包依赖、既有 Example、项目日志/JSON 边界及版本/回滚限制。必要的 gopls/rg 调用方查询使用已核实入口/文档，不把仓库内查询当外部消费者证明 | 导出与文件指纹、实际依赖、M1～M11、结果/遗留矩阵；CONTRIBUTING 仅同步实际性质用例状态，本计划写回并回读 | 不能复核证据或需发布/扩展实现时停在交付记录，不宣称全部门禁通过；导出/依赖/范围核对与结果写回已完成，已回读核对 |

所有命令分别记录 stdout/stderr、cwd、UTC 时间、环境覆盖和退出码；记录器沿用已读阶段 4 实现并仅在临时工装补充启动异常记录。统一 GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off，Go/lint 缓存及覆盖率/原始输出位于本轮目录。Fuzz 若产生真实失败语料，先记录事实并保留原始语料，不创造成功样本。本库无服务、存储、网络或容器依赖，不运行实际支付通知、日志平台或生产链路。

### 18.3 风险、失败处理与回滚

- 用例只固定已有字节/rune、状态/原因与只读出口契约；测试共享 snapshot 时不修改切片，各 worker 独占 writer。输入错误对象保持稳定，不承诺能隔离任意外部阻塞或 panic。
- 不把 Fuzz/竞态一次通过外推到所有实现、平台和消费者；当前构建为 darwin/amd64、无自定义标签。后台负载、供电或热状态未取得证据时保持未知，不能据样本数保证性能门槛。
- 已核实 gopls v0.23.0 及 references 命令文档；如查询失败按真实结果记录，保留文本/源码范围限制，不安装或臆造语义工具。
- 出现本轮新增失败先如实追加错误记录，区分夹具、运行缺陷和既有门禁；同一问题最多纠正 3 轮，不能删除断言或改语义凑绿。旧 lint/fix、覆盖率与性能缺口只做本轮取证，保留最终验收未通过的可能。
- 回滚仅逆转本轮两个测试文件和必要文档差异，参照 before 保留已提交生产实现与进入本轮前的文档修改；保留失败证据及错误记录，不用 reset/clean。临时基准模块可停用，没有服务、存储、依赖或注册状态需迁移；本轮不主动执行回滚或发布。

### 18.4 实际结果（2026-10-21）

#### 18.4.1 本轮改动、环境与证据

- 新增/修订仅两个测试文件：`sanitize_test.go` 的 15 个命名边界用例、9 个 Fuzz 种子和私有性质断言；`test/integration_test.go` 的 16 worker 共享只读快照、独占 writer 的 Join/标准包装文本与 JSON 集成用例。没有修改运行行为或新增导出/配置/扩展点。
- 性质断言使用原字节前缀、独立 rune 计数、后缀与 UTF-8/幂等性约束，不把生产遍历算法复制为 oracle。真实有符号阈值、int 极值、组合字符和非法 UTF-8 均保留。并发用例等待所有 worker 退出，输出后检查快照、根全文、主状态及两个原始 cause。
- 验证使用 `/usr/local/bin/go` go1.27.2、darwin/amd64、GOAMD64=v1、CGO_ENABLED=1、无 GOWORK。实际 benchmark 头部为 Intel Core i7-9750H @ 2.60GHz、并发度后缀 12。golangci-lint 2.14.0/go1.27.2、gopls v0.23.0、benchstat x/perf v0.0.0-20260908200009-22c9c6c9d4da（工具自身 go1.27.1）均按本机入口使用，未安装或升级。
- 下表共对应 42 次有完成元数据的命令。证据根目录为 `/private/tmp/goerr-stage6.iourzyn1`；每个记录名对应同名 `.json`（完整 argv、cwd、环境覆盖、UTC 起止及退出码）和 `.txt`（合并输出）。`verification-summary.json` 保存事件汇总，`benchmark-analysis.json` 保存全部 440 个原始数值与基于实测输入的中位数计算，`benchmark-order.json` 保存真实采样顺序。临时目录可能被系统清理；本计划同时保留关键数值与判定，原始文件不是永久归档。
- 所有记录器命令覆盖 `GOTOOLCHAIN=local`、`GOFLAGS=-mod=readonly`、`GOPROXY=off`，Go/lint 缓存分别为证据目录下 gocache、lint-cache。下表 cwd 的 root 为项目绝对路径 `/Users/xiaozhaofu/go/src/my-gtkit-package/goerr`，old/new 分别为证据目录下同名临时模块。
- 原 `targeted-race` 在用户中断后未留下完成元数据，会话 27584 已不可查询，原退出码/结果未知；没有计入上述 42 次。另名重新执行的 `targeted-race-resumed` 才是专项通过依据，ERR-20261021-017 已据此关闭。

| 记录名（.json/.txt） | cwd | 实际命令 | 退出码 |
| --- | --- | --- | --- |
| tool-go-version | root | `/usr/local/bin/go version` | 0 |
| tool-go-env | root | `/usr/local/bin/go env -json GOVERSION GOOS GOARCH GOAMD64 CGO_ENABLED GOWORK GOMOD GOTOOLCHAIN` | 0 |
| tool-lint-version | root | `/Users/xiaozhaofu/go/bin/golangci-lint version` | 0 |
| tool-benchstat-version | root | `/usr/local/bin/go version -m /Users/xiaozhaofu/go/bin/benchstat` | 0 |
| tool-gopls-version | root | `/Users/xiaozhaofu/go/bin/gopls version` | 0 |
| targeted-race-resumed | root | `/usr/local/bin/go test -json -race -count=1 -timeout=5m -run 'TestClampMessage\|TestDiagnosticConcurrentOutlets\|FuzzClampMessage' ./...` | 0 |
| full-race | root | `/usr/local/bin/go test -json -race -shuffle=on -count=1 -timeout=5m ./...` | 0 |
| build | root | `/usr/local/bin/go build ./...` | 0 |
| vet | root | `/usr/local/bin/go vet ./...` | 0 |
| lint | root | `/Users/xiaozhaofu/go/bin/golangci-lint run ./...` | 1 |
| fix | root | `/usr/local/bin/go fix -diff ./...` | 1 |
| coverage | root | `/usr/local/bin/go test -count=1 -coverprofile=/private/tmp/goerr-stage6.iourzyn1/coverage.out ./...` | 0 |
| coverage-functions | root | `/usr/local/bin/go tool cover -func=/private/tmp/goerr-stage6.iourzyn1/coverage.out` | 0 |
| fuzz | root | `/usr/local/bin/go test . -run '^$' -fuzz '^FuzzClampMessage$' -fuzztime=30s -parallel=4` | 0 |
| bench-old-01～10（10 条相同命令） | old | `/usr/local/bin/go test -run '^$' -bench . -benchmem -benchtime=100ms -count=1 -trimpath .` | 0 |
| bench-new-01～10（10 条相同命令） | new | `/usr/local/bin/go test -run '^$' -bench . -benchmem -benchtime=100ms -count=1 -trimpath .` | 0 |
| benchmark-inspect | root | `/usr/local/bin/go test -run '^$' -bench '^BenchmarkInspect$' -benchmem -benchtime=100ms -count=10 -trimpath .` | 0 |
| benchstat-compare | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-stage6.iourzyn1/benchmark-old.txt /private/tmp/goerr-stage6.iourzyn1/benchmark-new.txt` | 0 |
| benchstat-inspect | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-stage6.iourzyn1/benchmark-inspect.txt` | 0 |
| api-audit | root | `/usr/local/bin/go run /private/tmp/goerr-stage6.iourzyn1/audit_api.go /private/tmp/goerr-stage6.iourzyn1/old /private/tmp/goerr-stage6.iourzyn1/new /private/tmp/goerr-stage6.iourzyn1/api-comparison.json` | 0 |
| references-inspect | root | `/Users/xiaozhaofu/go/bin/gopls references -d /Users/xiaozhaofu/go/src/my-gtkit-package/goerr/diagnostic.go:61:6` | 0 |
| references-clamp | root | `/Users/xiaozhaofu/go/bin/gopls references -d /Users/xiaozhaofu/go/src/my-gtkit-package/goerr/sanitize.go:46:6` | 0 |
| module-list | root | `/usr/local/bin/go list -m all` | 0 |
| dependency-list | root | `/usr/local/bin/go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' ./...` | 0 |

#### 18.4.2 功能与质量结果

| 项目 | 本轮真实结果 | 判定与范围 |
| --- | --- | --- |
| 专项竞态重跑 | 退出 0；27 条带 Test 字段的 pass 事件，两个包通过，无 fail/skip | 新增截断/并发出口及 Fuzz 种子通过；原中断进程仍是未知 |
| 完整竞态/随机顺序 | 退出 0；156 条带 Test 字段的 pass 事件，两个包通过，无 fail/skip | 包含子测试、Fuzz 种子和 6 个 Example，不把事件数当顶层测试函数数；shuffle 输出种子为 1792577695916812000、1792577698628725000 |
| build / vet | 各退出 0，无输出 | 当前 darwin/amd64、当前模块和无自定义标签下通过 |
| golangci-lint | 退出 1，6 项 | **未通过**：error.go:84/90/101/103/105 的 fmt 写入 errcheck；resolve_test.go:24 的 ineffassign。输出与阶段 4 lint-corrected.txt 逐字一致，无本轮新增项 |
| go fix -diff | 退出 1，只输出差异，未应用 | **未通过**：error_test.go 的 stackFrames 仍建议 Split→SplitSeq；输出与阶段 4 fix.txt 逐字一致 |
| 覆盖率 | 命令和 cover -func 各退出 0；根包 66.2%，外部 test 包无运行语句 | **未达 ≥80%**；Inspect/walk/frames/reserve、ClampMessage 和其余脱敏函数均 100% 语句覆盖。弃用 Is 为 0%，多项 Status getter 为 0%，不能以局部 100% 替代整体门槛或完整路径证明 |
| Fuzz | 30s 配置、4 worker、9/9 初始种子完成；444,896 次执行，35 个 new interesting（共 44）；PASS，测试耗时 31.774s，退出 0 | 性质断言通过；有限时间随机探索不是全输入证明，新增 interesting 不是缺陷数 |

完整 race 包含 ExampleNewStatus、ExampleMustNewStatus、ExampleClampMessage、ExampleRedactForLog、ExampleInspect、ExampleDiagnostic。输出切换/并存、写入前及部分写入后失败、诊断失败跳过输出的既有用例也在本轮实际执行。这里仅验证内存 writer 与调用方投影，不是日志平台 ACK、网络故障或真实客户端服务验证。

ERR-20261021-018 已保留 lint/fix/覆盖率失败，状态 open；原 003/004/008 也未关闭。未修改 Format、旧 ResolveStatus 测试、stackFrames 或为凑覆盖率批量补 getter 用例。没有本轮新增的运行缺陷反例，也没有据此扩大生产修复范围。

#### 18.4.3 性能方法与旧路径对比

- old 取自真实 commit `d1f31a90c7080fffce42fc591ca4003a744a0eef` 的六个生产 Go 文件及 go.mod/go.sum；new 为当前七个生产 Go 文件及同一模块文件。两边仅放同一份已保存的阶段 1 final-benchmark_test.go；其内容与当前 benchmark_test.go 的共同部分核对一致（当前只额外含 Inspect 基准及所需 errors import）。夹具 SHA-256 为 `63a9170c87471187908ec9e95fe7f6a6372253c3d22bb19d11ad5bb832eb3319`，来源、逐文件指纹、git 命令/退出和模块全文见 benchmark-preparation.json。
- 原 old go.mod 为 `module github.com/gtkit/goerr`、`go 1.26`；**仅临时副本**将声明改为 go 1.27.0。new 也是 `module github.com/gtkit/goerr`、`go 1.27.0`，两份 go.sum 均为空，无 require。比较的是两种实现同在 Go 1.27 下的结果，不是 Go 1.26 与 Go 1.27 的跨工具链比较，也未改项目模块文件。
- 两边均 `-trimpath`；前 10 轮每轮各跑 old/new 一次，奇数轮 old→new，偶数轮 new→old，20 次命令均退出 0；随后单独采当前工作区 Inspect 的 10 组、每组 10 次。共同基准各 17×10，Inspect 10×10，共 440 个样本。每组 benchtime=100ms，没有同时进行其他本任务重验证。后台负载、供电与热状态仍未知。
- 表中 ns/op 为每组 10 个原始值的 `statistics.median`（偶数样本取中间两数平均）；Δ=(new/old−1)×100%。B/op 与 allocs/op 也取中位数；正常旧路径的各轮原始分配值均相同。此为真实采样后的计算，不是估算性能。

| 共同基准 | old ns/op | new ns/op | 中位耗时 Δ | B/op old→new | allocs/op old→new |
| --- | --- | --- | --- | --- | --- |
| New | 1,029.5 | 1,044 | +1.41% | 136→136 | 3→3 |
| Wrap | 97.055 | 92.81 | -4.37% | 104→104 | 2→2 |
| Wrapf | 223.5 | 232.5 | +4.03% | 120→120 | 3→3 |
| WrapStatus | 90.865 | 95.935 | +5.58% | 128→128 | 2→2 |
| ConvertNil | 2.3715 | 2.3695 | -0.08% | 0→0 | 0→0 |
| ConvertPlain | 42.055 | 43.395 | +3.19% | 80→80 | 1→1 |
| ConvertItem | 6.6755 | 6.547 | -1.92% | 0→0 | 0→0 |
| ConvertDepth/1 | 13.2 | 12.425 | -5.87% | 0→0 | 0→0 |
| ConvertDepth/4 | 34.59 | 32.46 | -6.16% | 0→0 | 0→0 |
| ConvertDepth/16 | 133.7 | 127.05 | -4.97% | 0→0 | 0→0 |
| ConvertDepth/64 | 487 | 443.4 | -8.95% | 0→0 | 0→0 |
| ConvertParallel | 2.2885 | 2.5025 | +9.35% | 0→0 | 0→0 |
| HasCode | 15.03 | 12.93 | -13.97% | 0→0 | 0→0 |
| WithStackPlain | 503.75 | 488.6 | -3.01% | 112→112 | 2→2 |
| WithStackExisting | 13.76 | 12.98 | -5.67% | 0→0 | 0→0 |
| WithStackMissing（新增捕栈，单列） | 6.951 | 465.85 | +6601.91% | 0→112 | 0→2 |
| FormatVerbose | 1,865 | 1,816.5 | -2.60% | 648→648 | 12→12 |

16 条正常旧路径的 B/op、allocs/op 未增加，采样中位耗时最大回退为 ConvertParallel 的 **9.35%**，本轮样本均在 10% 预算内。但 benchstat 的时间比较对这 16 条均为 `~`（p≥0.052），未检出显著差异不等于证明等价或稳定满足 10% 上界；例如 old ConvertParallel 为 ±121%，new 为 ±29%。**样本中位数条件满足，稳定耗时门槛仍未确认**，未据此宣称优化提速或总性能验收通过。按既定 10 轮结束，未补测。

WithStackMissing 从 6.951 ns、0 B、0 次分配变为 465.85 ns、112 B、2 次分配；这是由“返回无栈 Item”改为实际捕获缺失起源栈后的成本，单列为预期修复代价。benchstat 包含该组的 17 组时间 geomean +25.30% 不能用作正常旧路径整体回退结论；分配 geomean 因含 0 无法计算，也不是命令失败。没有实施缓存、池化或生产性能重构。

#### 18.4.4 Inspect 主动读取成本

以下为当前工作区独立测量，每组 10 次；包括有栈、无栈、Join 分支、不同深度与预算边界，同一输入被反复读取。表内仅报告成本，没有虚构绝对性能 SLO。

| Inspect 场景 | 中位 ns/op | 中位 B/op | 中位 allocs/op |
| --- | --- | --- | --- |
| nil | 9.2535 | 0 | 0 |
| plain | 394.05 | 289 | 4 |
| item | 2,214.5 | 1,618 | 11 |
| no_stack | 767.85 | 754 | 7 |
| branches | 4,793 | 4,031 | 26 |
| depth_1 | 2,332 | 1,779 | 13 |
| depth_4 | 3,018 | 2,998 | 18 |
| depth_16 | 8,374 | 12,394 | 39 |
| depth_64 | 18,617.5 | 29,018 | 93 |
| visit_limit | 66,386.5 | 116,088 | 287 |

benchstat 报告 item 耗时 ±52%、depth_16 ±25%、depth_64 ±24%、visit_limit ±15% 等波动，不能从这些数值推导其他平台容量或真实请求延迟。深度 64 中位约 18.62µs、29,018 B、93 次分配，visit_limit 约 66.39µs、116,088 B、287 次分配；这些量化主动诊断的预算边界成本，不代表有全局缓存或零分配保证。

#### 18.4.5 导出、依赖、契约与逐项验收

- 标准库 AST 工具按生产文件比较导出函数/方法签名、完整类型声明（含字段/接口/守卫）及 const/var 声明；const 组保留顺序以核对 iota，组作为一个条目。原 189 个声明条目、当前 192 个；**仅新增 func Inspect、type Cause、type Diagnostic，无删除，唯一既有声明差异为 Version v1.4.0→v1.5.0**，退出 0、matches_expected=true。记录为 api-comparison.json；不是正式 apidiff 或所有行为兼容证明，未安装第三方审查工具。
- gopls references 实际退出均为 0；Inspect 得到 29 个位置、ClampMessage 得到 6 个位置，均包含声明。坐标由实际源码计算，记录见 references-locations.json。配合本轮逐行 diff 回审，两个测试文件的改动仅对应性质/并发出口，未发现需改调用方的本轮生产行为变化。查询限于当前工作区，外部消费者、反射、动态使用及其他标签/平台仍未验证。
- go list -m all 仅本模块；非标准包依赖列表只有本模块与 test 子包。go.mod/go.sum 未变，库无外部运行依赖；测试 JSON 沿用标准库，是零依赖公共库互操作验证的既定取舍。未增加 logger、JSON、接口层、Options、注册器、后台任务或新生产文件。
- README/CHANGELOG 与 Example 对照保持阶段 5 语义：三个诊断符号、内部原文/路径与客户端 DTO 分离、普通函数出口、首个 Item 主状态、真实 cause、切片所有权及调用方显式投影/脱敏。日志使用沿用调用项目锁定的 gtkit/logger/JSON；本轮没有重跑 SDK 参考编译或连接日志平台，历史真实缓存 API 编译证据仅见第 16 节。渠道原文依据仍是阶段 5 的具体协议，不声称本轮重新抓取或真实支付验证。

| 矩阵 | 本轮结论 | 可复核依据与限制 |
| --- | --- | --- |
| M1 兼容 | 已测场景通过；弃用 Is 行为未验证 | 旧签名/类型/常量结构比对、完整 race 和构造/转换契约；弃用 Is 语句覆盖 0%，不声称所有弃用 API 已验证 |
| M2 补栈 | 当前平台用例通过 | 缺栈/有效栈/标准包装、身份、首帧、原 cause 的既有回归均在 full-race 执行 |
| M3 分支 | 仓库内用例通过 | Join、多 %w、第三方 Error/虚拟 As、首个 Item 规则；新增并发 Join 出口检查两分支数据 |
| M4 边界 | 用例通过，外部方法限制保留 | nil/typed-nil、不可比较/重复/循环、256 节点/子槽及 16384 rune、失败零快照；不能终止外部阻塞/panic |
| M5 所有权 | 用例与 race 通过 | 独立返回值修改、重复/并发 Inspect、多栈分支及共享快照只读出口；并发修改同份切片仍须调用方同步 |
| M6 出口 | 内存集成通过 | 两个普通函数切换/并存、独占 writer、写入前/部分失败与诊断失败跳过；未验证真实平台接收 |
| M7 规则 | 本轮边界核对通过 | 依赖列表、导出/文件差异及 README 既有日志约定；具体调用方版本和平台仍未知 |
| M8 信息边界 | 既有契约及模拟出口通过 | 显式投影/脱敏，业务错误/响应字段不被日志错误替换；不承诺自动检测全部敏感信息或真实公网服务已验收 |
| M9 截断 | 表格、种子与有界 Fuzz 通过 | 15 个用例、9 个种子、444,896 次真实执行；保持既有字节/rune/省略号语义 |
| M10 性能 | 分配及样本中位数条件满足，稳定耗时门槛未确认 | 440 个原始样本、同源夹具/Go/参数、benchstat 与中位数；正常路径最大 +9.35%，缺栈捕获单列，波动和未知负载保留 |
| M11 版本与文档 | 本地核对通过，未发版 | go 1.27.0、实际 go1.27.2、三个新增符号、Version v1.5.0 与 Unreleased；6 个 Example 通过。远端 tag/Release 仅沿用阶段 3 的历史核实，本轮未重新查询 |

合规结果分别记录：本轮测试使用 range-over-int、WaitGroup.Go、并行子测试等现代写法；新增 Options/泛型/UUID 不适用。库最小导出面、标准库依赖、已有 GoDoc/Example、错误/并发与本轮职责范围已核对。**全仓现代化 fix 目标、lint、≥80% 覆盖率和稳定性能结论均未达成或未确认，不能给出全部勾选的合规摘要。**

#### 18.4.6 范围、回读、回滚与停止位置

已逐段回读本轮方案变更、完整阶段 6 结果和 CONTRIBUTING，并将表内命令与 42 份完成元数据逐项核对，均一致。25 份阶段 6 开始时的文件指纹中，仅两个测试文件、CONTRIBUTING、本计划和错误记录这 5 份白名单文件变化，其他 20 份受保护内容全部一致，无白名单外改动；结果见证据目录 result-docs-readback.json。CONTRIBUTING 只同步并发、性质与性能的实际测试状态，没有因最后文档写回重新运行测试或构建。

阶段 1～5 执行记录（第 13～17 节）与阶段 6 开始时的保存内容一致。生产源码、既有基准/Example、模块/版本、README/CHANGELOG 及进入本轮前的改动、AGENTS/.gitignore 等受保护内容保留。

回滚仍按 18.3，仅撤本轮两个测试文件和必要文档差异；本计划/CONTRIBUTING 撤回须保护进入阶段 6 前的四份文档改动。真实失败记录和证据保留，临时 old/new 模块停用即可，无服务/存储/注册状态需迁移。未实际执行回滚，未把备份存在当作回滚验证。

**停止位置：阶段 6 执行与取证结束，结果已写回，已回读核对；整体门禁未通过。** 没有修改生产源码、依赖或公开契约，没有自动修旧门禁、进入额外实施阶段、修改 OpenSpec、提交、推送、tag、发布或首次 v2。外部消费者适配、其他平台/工具链、真实日志/支付链路、供应链及生产回滚未在本轮验证。


## 19. 剩余门禁修复续项（本轮获准，先方案后实施）

### 19.1 需求、事实与目标

- 来源为用户本轮明确要求：先补充修复方案并回读，再分步解决 6 项 lint、1 处 fix、覆盖率不足，并重新核对性能门槛；不提交或发布。本节不冒充原阶段 7，不重排既定 0～6。
- 当前 HEAD 经真实 rev-parse 核实仍为 `f6a30726b372b3dc07feda42e161d27de6cb5d77`，原实现参考为 `d1f31a90c7080fffce42fc591ca4003a744a0eef`。进入本轮已有 CHANGELOG、CONTRIBUTING、OPTIMIZATION_PLAN、README、sanitize_test.go、test/integration_test.go 六份工作区差异，全部保留。
- 本轮已备份 25 份实际文件至 `/private/tmp/goerr-gates-p3pz_g87/before`，指纹与 revision/status 命令见 before-sha256.json、preflight.json。证据目录可能被系统清理，计划需同时保留关键结果。
- 当前已读源码对应五处 Format fmt 写入未检查 error、ResolveStatus 测试一次无效 got 赋值，以及 stackFrames 的 Split 循环。阶段 6 的真实 lint/fix 输出和 66.2% 覆盖率仍为历史基线，不当作本轮新执行。
- 当前模块 `github.com/gtkit/goerr`、go 1.27.0、Version v1.5.0（Unreleased），无外部 require。本轮实际 go env 为 go1.27.2、darwin/amd64、GOAMD64=v1、CGO_ENABLED=1、GOEXPERIMENT/GOWORK 空；沿用已核对入口，不升级工具或依赖。
- 已逐行读取 Format、相关测试、Status 及 codeTable；真实源码有 136 个 Status getter、136 个预构建映射、136 个数字码与元信息，无映射缺失。新测试固定这些公开业务契约和单例身份，不改码表、不靠删语句或降低门槛提高覆盖率。
- 已读取当前 Go 标准库 fmt.State/Formatter 文档及 fmt 实现：Format 无 error 返回；标准 fmt 内部 pp.Write 只写缓冲并返回 nil，外层 Fprintf/Fprint 返回最终 writer 的 n/error。自定义直接 State 调用可出现写入错误。修复只让 Format 遇到写错误停止，不能增加返回值或谎称能够经 Formatter 回传该错误。
- 目标：lint 退出 0；go fix -diff 退出 0 且输出为空；完整 build/vet/race 通过；根包语句覆盖率 ≥80%，导出 Status/弃用 Is 与 Format 错误路径有具体断言；同条件复核正常旧路径分配不增加及中位耗时 ≤10% 回退，缺栈新增捕获成本单列。未取得证据前均待验证。

### 19.2 范围、设计与兼容约束

| 文件 | 本轮具体范围 |
| --- | --- |
| error.go | 仅 Format 五处 fmt 调用检查 error；检测到错误即 return，含正文与栈帧失败。补充 GoDoc 说明接口限制，不改其他构造/解析/补栈逻辑 |
| error_test.go | 私有 fmt.State 故障夹具；正常格式/直接状态写失败/外层 Fprintf 部分失败反例、原错误不变；弃用 Is 的 nil/哨兵/包装/Join/typed-nil/不相等契约；stackFrames 仅 Split→SplitSeq |
| resolve_test.go | 原 nil func 错误分支不再赋给未使用的 got，保留原 error 断言与其余用例 |
| status_test.go | 136 个 getter 的静态表格契约：固定数字业务码、HTTP 码、默认文案、Code.Message 一致性、重复调用同一只读实例；不修改包级单例 |
| README.md | 仅补 Format 写失败终止、接口不返回 error、调用方处理外层 fmt 写入错误的必要说明 |
| CHANGELOG.md | 仅在现有 Unreleased 的 Fixed 中同步本次 Format 修正和门禁测试；保留原版本与历史内容 |
| CONTRIBUTING.md | 仅更新实际门禁/覆盖率/性能结果和新增契约测试状态 |
| OPTIMIZATION_PLAN.md、.harness/error-journal.md | 方案/真实结果/回读，以及真实失败追加与有证据的关闭；不改规则和示例历史 |

其他生产源码、benchmark_test.go、example_test.go、sanitize_test.go、test/integration_test.go、go.mod/go.sum/version.go、AGENTS/.gitignore、Makefile、目录与既有架构保持本轮开始时内容。临时分析工具/基准模块只放证据目录。

兼容约束：

1. 无新增/删除导出符号、字段、接口、配置或依赖，module/import、最低 Go 与 Version 均不改；可插拔出口继续是调用方普通函数，不引入 logger/log/slog 或接口层。
2. Format 的 `%v/%s/%q/%+v`、未知 verb、nil/零值、有栈/无栈的成功输出保持；不增加捕栈、递归打印 cause 或重复栈，不改宽度/精度现有语义。唯一失败行为差异为自定义 State 写入报错后停止后续写入，已写部分无法撤销，不重试、补日志或改业务错误。
3. 在五个 fmt 调用点就地检查 error；这是标准错误处理语句，保留各格式分支的语义，不为去除相似 if 新增生产 helper/抽象或增加格式化分配。无返回值的 Formatter 不能传播 State 错误，由 State 自身保存或外层 fmt writer 调用方按接口处理。
4. 测试期望从已读取的现有公开契约冻结为字面量，运行时不读取 codeTable/Status 实现生成 oracle。缺陷若指向现有码表语义，先停下记录，不擅自变更业务映射。
5. 不因门禁修改启用忽略项、nolint、排除文件、降低覆盖阈值、替换真实库、删除断言或升级工具；不扩大到未授权的性能改造、缓存或池化。

### 19.3 分步实施、验收证据与停止条件（已执行并回读）

| 步骤 | 前置条件 | 具体工作 | 检查与交付证据 | 停止条件 / 当前状态 |
| --- | --- | --- | --- | --- |
| G0 方案回读 | 已读最新指令、规则、错误记录、原计划与相关源码/参考 | 保存并回读本节；记录本轮基线、白名单及性能判定方法 | 实际文档内容、快照/指纹和来源；代码未动 | 方案/规则无法读取或存在关键缺口即停止；已先保存并实际回读；G0 完成当时 Go 文件未修改 |
| G1 格式化与旧 lint/fix | G0 回读完成 | 先添加正常格式与可失败 State/外层 writer 回归，执行原 Format 的反例；真实预期失败先记错误再修五处 error 检查；修 got 赋值与 SplitSeq | 原反例和修后专项 race/完整受影响包测试与 lint、fix；核对停止后不再写、部分字节、nil/格式/业务错误/cause；读完整函数和调用链 | 非预期行为差异或需改签名/额外实现即停；同一问题最多自纠 3 轮；已完成，原 4 个反例转绿，根包 race/lint/fix 通过 |
| G2 公开契约覆盖 | G1 相关检查可复核 | status_test.go 固定 136 个 getter 的数字码/HTTP/文案/实例；error_test.go 固定弃用 Is 与标准 errors.Is 的已知输入结果 | 专项测试、完整 race、coverprofile 与 cover -func；统计导出函数/未覆盖语句，目标根包 ≥80%，不只报告命令退出 0 | 测试揭示现有业务映射异常先记录，不顺改码表；已完成，152 条专项 Test pass，根包覆盖率 100.0% |
| G3 完整质量检查 | G1/G2 源码与测试已保存、gofmt 完成 | build、vet、golangci-lint run ./...、只读 go fix -diff、go test -json -race -shuffle=on -count=1 -timeout=5m ./...、覆盖率及 cover -func | 独立 stdout/stderr、cwd、环境、UTC 起止、退出码；fix 无输出、lint 无告警、race 无 fail/skip、根包 ≥80% | 新失败先记录再修本轮问题；无依赖变化，不跑供应链/容器/外部平台；已完成，build/vet/lint/fix/race/覆盖率全部通过 |
| G4 性能复核 | 完整功能结果可复核，重验证结束 | 临时 old/new 模块使用原 commit 与修后生产文件及同一保存夹具；仅临时 old go 声明对齐 1.27.0；两边 trimpath/PGO off、cpu=12、benchtime=500ms、count=1，各 10 轮交替 AB/BA；当前 Inspect 10 组各 10 次、100ms 单独测 | 原始样本、源码/fixture/模块指纹、工具/环境、顺序与 benchstat；17 共同组区分 16 正常路径与 WithStackMissing，报告中位数、分配和不确定性；方法见下 | 依据结果最多补 10 轮有必要的共同基准，仍不确定就保留未验证，不反复筛样本或自动调实现；已完成 20 轮及 780 样本，7 组稳定耗时仍未确认 |
| G5 交付核对与写回 | 所有可做结果已读取、失败与未验证分开 | 对比本轮前/后导出声明须完全相同；核对生产只改 Format、受保护指纹、依赖/Version/未发布元数据；结合 gopls/rg 回审 Format 调用方；同步指定文档与真实结果并回读 | API/文件/命令与结果、调用方错误分支及边界、覆盖率/性能结论、旧错误的真实关闭证据 | 自动提交/发布、改外部消费者或扩大范围即停；导出/调用方/范围核对与写回已完成并回读 |

性能方法与判定提前固定：

- 使用现有 `/usr/local/bin/go`、`/Users/xiaozhaofu/go/bin/benchstat`，版本执行后记录；共同夹具为真实阶段 1 final-benchmark_test.go，先核对与当前共同部分一致，再复制，不重写原始实现。标准库、编译参数、Go 声明、CPU 取值、输入与夹具保持一致。
- 共同基准每侧从 10 个 500ms 样本开始，奇偶轮交换先后，不与 race/lint 并行。明确 `-trimpath -pgo=off -cpu=12`，不修改生产配置。Go 相关 GC/PGO 环境、系统负载/供电只用本机实际可读证据记录，未取得者保持未知，不读取完整环境或进程参数。
- 先分别检查 B/op、allocs/op 不增加和采样中位数 Δ≤10%；缺栈实际捕获成本单列。benchstat 的显著性检验不证明等价；保留其原始区间和 p 值。
- 为减少只看点估计的误判，另计算每组正常路径的保守中位数比值上界：n 个样本排序，取最大 k 使 `P(Bin(n,0.5)≤k−1)≤0.025`，old 第 k 个值作分母下界，new 第 n−k+1 个值作分子上界；两个单侧界按并集界给出至少 95% 的逐组比值上界。在独立/稳定采样假设下，上界≤1.10 才记录本轮测量范围内的耗时门槛通过；跨界保留未确认，不能把统计假设写成已控制所有机器负载。
- 若首批有跨界/不足样本/明显波动，仅追加一批最多 10 轮共同基准并合并全部样本，方法不变；仍超界/不确定则分别记录失败/未确认，不无限重测、删异常值或改预算。该区间是基于真实样本的计算，不冒充工具原生输出，不代表跨平台或真实请求 SLO。

### 19.4 风险、失败处理与回滚

- 提前 return 只针对 State 写入失败，风险为改变直接调用 Format 的异常写出次数；用正文前失败、正文部分失败、栈帧失败和外层 Fprintf writer 失败固定边界。成功文本、原 error/状态/cause/起源栈须保持，不许吞掉业务错误。
- 静态目录测试增加契约维护量，但能识别 getter 指向错误状态、码号/HTTP/文案漂移和实例重建；不将码表改动跟随生成测试来掩盖变化。136 项映射逐行回读，运行时无源码解析/动态 oracle。
- 代码位置变化可能影响帧行号与字符串长度，正常输出比较不锁死文件绝对路径；共同基准使用相同文件名/fixture 行号与 trimpath。性能受背景负载与热状态影响，不能保证复核一定达标，未确认如实保留。
- 新错误先记真实失败再处理，同一问题最多 3 轮自纠；遇到范围、外部状态或依赖缺口停止对应步骤并完成不依赖它的授权工作。原错误只在有对应修复和真实检查证据后关闭，不能用 plan 状态代替事实。
- 回滚只逆转本轮 error.go Format、三个测试文件及必要文档差异，参照 before 中本轮开始时内容；保护已有六份工作区差异及其此前实现，不 reset/clean。新失败记录保留，临时基准停用即可。无存储、服务、配置、依赖或导出 API 迁移；本轮不主动回滚、提交或发布。

### 19.5 实际结果（2026-10-21）

#### 19.5.1 修复与兼容结果

- **6 项 lint 和 1 处 fix 已解决。** 生产仅改 `Item.Format` 的五处 fmt error 检查及必要 GoDoc，写入失败后立即 return；方法保持不足 50 行，没有生产 helper、配置或扩展点。正常格式、nil/零值、原错误/状态/文案/cause/栈保持。`fmt.Formatter` 无 error 返回，未把 State 故障错误写入业务链；外层 Fprintf 的部分 n/error 由调用方处理。
- `resolve_test.go` 仅把 nil func 分支的无效 got 赋值改为 `_`；`error_test.go` 的栈帧 helper 仅 Split→SplitSeq，保留原过滤逻辑。未改转换、包装、起源栈捕获、状态码表或其他生产实现。
- 新回归含 9 个正常格式场景、8 个直接 State 故障场景、2 个外层 Fprintf 故障场景、14 个弃用 Is 场景。原 Format 的 4 个 verbose 故障子用例确实失败：正文/帧写错误后仍调用 writer 共 4 次；正常格式与外层 writer 用例通过。原命令退出 1、17 条 Test pass，保留为修前反例，不混入最终失败计数。修后原断言完整保留并通过。
- Status 测试固定 136 个 getter 的数字业务码、HTTP 码、默认文案、Code.Message 和共享实例身份。已逐行回读所有静态条目，来源/指纹与冻结数值见 catalog-contract.json；测试运行时不读取源码或 codeTable 生成期望。函数较长来自 136 行独立契约数据，只有一个短断言循环；没有移动文件、引入生成流程或抽象成公共能力。
- 本轮没有新增/删除导出符号、接口、字段、Options、依赖或版本变更。可插拔出口仍为调用方普通函数；JSON/日志边界保持，未引入日志库或接口层。

#### 19.5.2 真实命令与环境

证据目录为 `/private/tmp/goerr-gates-p3pz_g87`，包含各记录名的 `.json`（完整 argv/cwd、环境覆盖、UTC 起止与退出码）和 `.txt`（合并 stdout/stderr）。下表从已读元数据生成；old/new 分别为该目录下临时同名模块，root 为项目 `/Users/xiaozhaofu/go/src/my-gtkit-package/goerr`。记录器统一 GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off，GOCACHE/GOLANGCI_LINT_CACHE 分别为该目录 gocache/lint-cache。实际 go1.27.2、darwin/amd64、GOAMD64=v1、CGO_ENABLED=1，GOEXPERIMENT/GOWORK 空；lint 2.14.0/go1.27.2、gopls v0.23.0、benchstat x/perf v0.0.0-20260908200009-22c9c6c9d4da（工具由 go1.27.1 构建），均未升级。

共 78 份完成元数据；40 次共同基准合并成 old/new 两行展示，其余逐条保留。首批分析发生的数据/记录文件重名已单独保留原执行证据为 analysis-initial-metadata，并由另名 analyze-first 从原始输入重建，详见 19.5.6。临时目录不是永久归档，本计划同时保存关键数值与结论。

| 记录名（.json/.txt） | cwd | 实际命令 | 退出码 |
| --- | --- | --- | --- |
| analysis-initial-metadata | root | `python3 /private/tmp/goerr-gates-p3pz_g87/analyze_benchmarks.py` | 0 |
| analyze-final | root | `python3 /private/tmp/goerr-gates-p3pz_g87/analyze_benchmarks.py` | 0 |
| analyze-first | root | `python3 /private/tmp/goerr-gates-p3pz_g87/analyze_benchmarks.py` | 0 |
| api-audit | root | `/usr/local/bin/go run /private/tmp/goerr-gates-p3pz_g87/audit_api.go /private/tmp/goerr-gates-p3pz_g87/before /private/tmp/goerr-gates-p3pz_g87/new /private/tmp/goerr-gates-p3pz_g87/api-comparison.json` | 0 |
| bench-new-01～20（20 条） | new | `/usr/local/bin/go test -run '^$' -bench . -benchmem -benchtime=500ms -count=1 -trimpath -pgo=off -cpu=12 .` | 0 |
| bench-old-01～20（20 条） | old | `/usr/local/bin/go test -run '^$' -bench . -benchmem -benchtime=500ms -count=1 -trimpath -pgo=off -cpu=12 .` | 0 |
| benchmark-inspect | root | `/usr/local/bin/go test -run '^$' -bench '^BenchmarkInspect$' -benchmem -benchtime=100ms -count=10 -trimpath -pgo=off -cpu=12 .` | 0 |
| benchstat-final | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-gates-p3pz_g87/benchmark-old.txt /private/tmp/goerr-gates-p3pz_g87/benchmark-new.txt` | 0 |
| benchstat-first | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-gates-p3pz_g87/benchmark-old.txt /private/tmp/goerr-gates-p3pz_g87/benchmark-new.txt` | 0 |
| benchstat-inspect | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-gates-p3pz_g87/benchmark-inspect.txt` | 0 |
| benchstat-version | root | `/usr/local/bin/go version -m /Users/xiaozhaofu/go/bin/benchstat` | 0 |
| build | root | `/usr/local/bin/go build ./...` | 0 |
| coverage-functions | root | `/usr/local/bin/go tool cover -func=/private/tmp/goerr-gates-p3pz_g87/coverage.out` | 0 |
| coverage | root | `/usr/local/bin/go test -count=1 -coverprofile=/private/tmp/goerr-gates-p3pz_g87/coverage.out ./...` | 0 |
| dependency-list | root | `/usr/local/bin/go list -deps -f '{{if not .Standard}}{{.ImportPath}}{{end}}' ./...` | 0 |
| error-journal-read | root | `bash .harness/scripts/read-error-journal.sh .` | 0 |
| fix | root | `/usr/local/bin/go fix -diff ./...` | 0 |
| fmt-formatter-doc | root | `/usr/local/bin/go doc fmt.Formatter` | 0 |
| fmt-state-doc | root | `/usr/local/bin/go doc fmt.State` | 0 |
| format-text-references | root | `rg -n -e 'Format\(' -e '%\+v' -g '*.go' .` | 0 |
| full-race | root | `/usr/local/bin/go test -json -race -shuffle=on -count=1 -timeout=5m ./...` | 0 |
| g1-fix | root | `/usr/local/bin/go fix -diff ./...` | 0 |
| g1-lint | root | `/Users/xiaozhaofu/go/bin/golangci-lint run ./...` | 0 |
| g1-race | root | `/usr/local/bin/go test -json -race -count=1 -timeout=5m .` | 0 |
| g2-contracts | root | `/usr/local/bin/go test -json -count=1 -run '^(TestStatusCatalog\|TestIsCompatibility)$' .` | 0 |
| go-env | root | `/usr/local/bin/go env -json GOVERSION GOOS GOARCH GOAMD64 GOROOT CGO_ENABLED GOMOD GOWORK GOEXPERIMENT GOTOOLCHAIN` | 0 |
| gofmt-g1 | root | `/usr/local/go/bin/gofmt -w error.go error_test.go resolve_test.go` | 0 |
| gofmt-g2 | root | `/usr/local/go/bin/gofmt -w error_test.go status_test.go` | 0 |
| gofmt-regression | root | `/usr/local/go/bin/gofmt -w error_test.go` | 0 |
| gopls-version | root | `/Users/xiaozhaofu/go/bin/gopls version` | 0 |
| help-build | root | `/usr/local/bin/go help build` | 0 |
| help-testflag | root | `/usr/local/bin/go help testflag` | 0 |
| lint-version | root | `/Users/xiaozhaofu/go/bin/golangci-lint version` | 0 |
| lint | root | `/Users/xiaozhaofu/go/bin/golangci-lint run ./...` | 0 |
| module-list | root | `/usr/local/bin/go list -m all` | 0 |
| original-format-counterexample | root | `/usr/local/bin/go test -json -count=1 -run 'TestItem_FormatTextContract\|TestItem_FormatStopsAfterWriteError\|TestItem_FprintfWriterFailure' .` | 1 |
| power-source-after | root | `/usr/bin/pmset -g batt` | 0 |
| power-source | root | `/usr/bin/pmset -g batt` | 0 |
| references-format | root | `/Users/xiaozhaofu/go/bin/gopls references -d /Users/xiaozhaofu/go/src/my-gtkit-package/goerr/error.go:81:16` | 0 |
| vet | root | `/usr/local/bin/go vet ./...` | 0 |

#### 19.5.3 功能与质量验收

| 项目 | 实际结果 | 判定 |
| --- | --- | --- |
| G1 受影响根包完整 race | 退出 0，167 条 Test pass，根包通过，无 fail | 原 4 个 verbose 反例转绿，正常输出及外层部分失败契约保持 |
| G2 公开契约专项 | 退出 0，152 条 Test pass，无 fail/skip | 136 个状态及 14 个 Is 子用例，加两个父测试；不是 152 个顶层函数 |
| build / vet | 各退出 0，无输出 | 通过 |
| golangci-lint | 退出 0，`0 issues.` | 原 6 项及新增测试无告警，无忽略/排除项 |
| go fix -diff | 退出 0，输出为空，未自动应用全仓改写 | 通过 |
| 完整 race / shuffle | 退出 0，330 条带 Test 字段的 pass、两个包通过，无 fail/skip | 含子测试、种子及 6 个 Example；不是 330 个顶层函数 |
| 覆盖率 / cover -func | 各退出 0；根包 100.0%，384/384 条语句，183 个函数条目均 100% | 超过 ≥80% 门槛；外部 test 子包无运行语句。mode:set 的语句覆盖不等于所有分支/输入证明 |
| AST 导出比较 | 退出 0，192→192 条声明，added/removed/changed 均空，matches_expected=true | 相对本轮开始时，无导出/类型/字段/接口/常量差异；const 组按一个声明条目统计 |
| 依赖 | 两项 go list 均退出 0；模块只有 goerr，非标准包只有 goerr 与其 test 子包 | 无外部运行依赖，go.mod/go.sum/Version 未变 |

完整 race 包含 ExampleNewStatus、ExampleMustNewStatus、ExampleClampMessage、ExampleRedactForLog、ExampleInspect、ExampleDiagnostic；shuffle 输出种子为 1792579556789432000、1792579557491172000。既有独立快照、Join、16 worker 只读出口、输出失败和截断种子均由本轮完整测试执行。ClampMessage 未改，本轮未重跑 30 秒持续 Fuzz；阶段 6 的真实 Fuzz 记录仍见 18.4，不混作本轮执行。

#### 19.5.4 性能复核

- old 为真实原 commit `d1f31a90c7080fffce42fc591ca4003a744a0eef` 的六个生产 Go 文件及模块文件，new 为修后的七个生产 Go 文件及模块文件。共同使用已读阶段 1 final-benchmark_test.go，SHA-256 `63a9170c87471187908ec9e95fe7f6a6372253c3d22bb19d11ad5bb832eb3319`；与当前 benchmark_test.go 的共同部分一致，未改函数体、输入或行号。源码、git 退出码与指纹见 benchmark-preparation.json。
- 两份临时模块均为 `module github.com/gtkit/goerr`、`go 1.27.0`、空 go.sum、无 require；只有 old 临时副本把原 `go 1.26` 声明对齐。当前项目 go.mod/go.sum 不改，不把该实验写成 Go 1.26 兼容证明。
- 共同基准采用同一 Go、`-trimpath -pgo=off -cpu=12 -benchtime=500ms -count=1`；奇数轮 old→new，偶数轮 new→old。首批 10 轮后，14 个正常路径的保守区间未确认，依计划追加一次 10 轮并合并全部样本。40 次共同基准命令均退出 0，各侧 17×20=340 个样本；Inspect 单独 10×10=100 个样本，总计 **780**。没有筛选异常值、修改预算或继续追加采样。
- 实际基准头部为 Intel Core i7-9750H @2.60GHz、darwin/amd64、后缀 12；选定 GOGC/GOMEMLIMIT/GOMAXPROCS/GOEXPERIMENT 环境均未设置，CPU 与 PGO 由上述测试参数明确控制。测前/测后 pmset 均报告 AC Power；未连续监测供电、温度或 CPU 频率。每个共同命令前后记录的 1 分钟系统负载范围为 3.85～32.67；没有读取完整环境或其他进程参数，也没有据负载认定具体性能根因。
- 表中中位数由真实样本 `statistics.median` 计算，Δ=(new/old−1)×100%。按 19.3 预先固定方法，n=20、k=6；old 第 6 个值与 new 第 15 个值构成保守比值上界。并集界覆盖率下限约 95.86%，只针对逐组中位数比值、依赖独立/平稳采样假设；不声明所有组同时有该置信度，不保证跨平台或真实请求 SLO。样本、点估计、区间及命令成功分别判定。

| 正常共同基准 | old 中位 ns/op | new 中位 ns/op | 中位 Δ | B/op old→new | allocs/op old→new | 保守上界 Δ | 区间判定 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| New | 954.45 | 987.9 | +3.50% | 136→136 | 3→3 | +29.81% | 未确认 |
| Wrap | 90.065 | 93.485 | +3.80% | 104→104 | 2→2 | +26.84% | 未确认 |
| Wrapf | 212.55 | 215.1 | +1.20% | 120→120 | 3→3 | +13.65% | 未确认 |
| WrapStatus | 92.525 | 93.02 | +0.53% | 128→128 | 2→2 | +12.58% | 未确认 |
| ConvertNil | 2.166 | 2.1 | -3.05% | 0→0 | 0→0 | +18.14% | 未确认 |
| ConvertPlain | 43.525 | 42.16 | -3.14% | 80→80 | 1→1 | +15.35% | 未确认 |
| ConvertItem | 6.0265 | 5.8705 | -2.59% | 0→0 | 0→0 | +13.35% | 未确认 |
| ConvertDepth/1 | 12.26 | 10.935 | -10.81% | 0→0 | 0→0 | +5.07% | 满足预算（条件区间） |
| ConvertDepth/4 | 31.85 | 28.49 | -10.55% | 0→0 | 0→0 | +3.27% | 满足预算（条件区间） |
| ConvertDepth/16 | 124.75 | 105 | -15.83% | 0→0 | 0→0 | -2.90% | 满足预算（条件区间） |
| ConvertDepth/64 | 445.35 | 387.95 | -12.89% | 0→0 | 0→0 | -2.91% | 满足预算（条件区间） |
| ConvertParallel | 2.1305 | 2.0235 | -5.02% | 0→0 | 0→0 | +8.29% | 满足预算（条件区间） |
| HasCode | 13.975 | 11.925 | -14.67% | 0→0 | 0→0 | +1.16% | 满足预算（条件区间） |
| WithStackPlain | 491.6 | 464.65 | -5.48% | 112→112 | 2→2 | +4.87% | 满足预算（条件区间） |
| WithStackExisting | 12.515 | 11.64 | -6.99% | 0→0 | 0→0 | +4.33% | 满足预算（条件区间） |
| FormatVerbose | 1,772 | 1,672 | -5.64% | 648→648 | 12→12 | +2.41% | 满足预算（条件区间） |

16 条正常路径的各轮原始 B/op、allocs/op 均未增加；采样中位耗时最大回退为 Wrap **+3.80%**，均在 10% 点估计预算内。保守条件区间有 9 组满足预算，**7 组未确认**：New、Wrap、Wrapf、WrapStatus、ConvertNil、ConvertPlain、ConvertItem；最大上界为 New +29.81%。因此**分配及样本中位数门槛通过，稳定耗时门槛仍未确认，整体性能验收不勾选通过**。

benchstat-final 对部分组检出负向差异，对其余正常组给出 `~`；无显著差异不是等价证明。其 New 的 new 时间区间为 ±22%、ConvertParallel 的 old ±83%、FormatVerbose 的 old ±40%，与系统负载读数一起保留为波动证据，不将下降值全部归因于本轮修复。FormatVerbose 当前中位 1,672 ns，B/op=648、allocs/op=12，未增加分配。

WithStackMissing 单列：old 6.5105 ns、0 B、0 次分配；new 440.95 ns、112 B、2 次分配，反映前序补栈修复的实际捕获代价。该组不是正常旧路径预算项；统计 JSON 对它保留原始相对比较和 expected_new_capture 标记，不能把其原始 fail/分配差异判为本轮正常路径失败。benchstat 的 17 组时间 geomean +21.30% 包含该组，不能用于正常路径总体回退结论；零分配组的分配 geomean 不适用。

首批 440 样本的最大正常中位回退为 +1.76%、保守上界最大 +29.40%，14 组未确认；首批结果和原输入另存 benchmark-analysis-first-result.json、benchmark-old-first.txt、benchmark-new-first.txt、benchstat-first.txt。最终输入为合并的 20 轮，不能把最终聚合内容当作首批命令执行时的输入。没有修改性能实现、增加缓存或池化。

Inspect 的主动成本仅报告本轮实测，未设绝对 SLO。各组 10 个样本、100ms，重复读取同一稳定输入；原始值和 benchstat-inspect 保留：

| Inspect 场景 | 中位 ns/op | 中位 B/op | 中位 allocs/op |
| --- | --- | --- | --- |
| nil | 12.29 | 0 | 0 |
| plain | 544.3 | 289 | 4 |
| item | 2,467 | 1,618 | 11 |
| no_stack | 939.95 | 754 | 7 |
| branches | 6,484.5 | 4,031 | 26 |
| depth_1 | 2,354.5 | 1,779 | 13 |
| depth_4 | 3,343.5 | 2,998 | 18 |
| depth_16 | 8,797 | 12,394 | 39 |
| depth_64 | 24,385 | 29,018 | 93 |
| visit_limit | 80,216.5 | 116,088.5 | 287 |

#### 19.5.5 调用方、范围与合规核对

- gopls v0.23.0 references 实际退出 0，准确位置为 error.go:81:16，返回声明和 error_test.go:267 的直接 State 调用。rg 补查到普通/未知格式、nil、栈 helper、回归和 benchmark 等 fmt 分派候选，已逐项回读；搜索行数不当作调用方数量。实际 Go fmt 源码通过 Formatter 接口分派到 Format，内部 State 仅写缓冲，标准 Sprintf/Fprintf 的正常输出无须适配；外层 writer 返回值仍由调用方处理，并有 0/7 字节部分失败回归。直接 State 故障调用已按终止写出契约覆盖。
- 未新增 error 返回或接口参数。行为差异仅为直接 State 写入失败后的终止；没有业务码、cause、返回快照所有权、日志出口或客户端响应适配需求。本轮当前平台/无自定义标签范围内核对，不把语义/文本搜索外推为仓库外消费者、反射、其他标签或平台已验证；真实日志/支付平台、供应链和生产回滚未运行。
- AST 工具比较本轮 before 与修后 new，192 个导出声明条目完全一致，包含类型字段/接口/初始化守卫和 Version。go list 仅本模块及 test 包，零外部运行依赖；真实代码和文档的 API、logger/JSON 边界与 v1.5.0 Unreleased 相符。没有重新查询远端发版状态。
- 文件指纹比较显示 25 份本轮开始时文件中仅白名单 9 份变化：error.go、error_test.go、resolve_test.go、status_test.go、README、CHANGELOG、CONTRIBUTING、本计划、错误记录；其余 16 份受保护内容一致。生产唯一差异为 Format；既有 sanitize_test.go、外部集成测试、benchmark/Example、模块/Version、规则/.gitignore 等均与本轮开始时一致，进入本轮前的改动保留。
- 第 13～18 节与本轮 before 的历史内容一致；当前状态另由本节记录。README 只补格式化错误契约，CHANGELOG 只补 Unreleased，CONTRIBUTING 只同步本轮契约与真实质量/性能结果。已逐段回读本轮方案/结果变更和三份说明文档的修改段落，并核对 78 份完成元数据、780 个样本与 25 份文件指纹，均与实际证据一致；回读记录见本轮目录 result-docs-readback.json。最后文档写回后未再运行测试或构建。

合规结果：

- [x] 实际 Go 1.27.2 与 go 1.27.0 声明一致；go fix -diff 无输出。
- [x] build/vet/lint、完整 race、6 个 Example 通过；根包语句覆盖率 100.0% ≥80%。
- [x] 零外部运行依赖、导出面/版本/可插拔边界保持；Format GoDoc/README/Unreleased 对应。
- [x] 真实写错误反例先失败后转绿，原业务错误与输出部分成功边界核对；代码职责单一，未新增生产抽象。
- [x] 基准 ReportAllocs、相同夹具/Go/参数、原始样本、benchstat 与不确定性分别保存。
- [ ] 稳定耗时门槛全部通过：7 组仍未确认，保留为整体未通过/未确认项。

新增 Options、UUID、泛型、配置、依赖升级及存储/服务迁移不适用；不勾选不存在的能力。没有用局部测试/语句覆盖冒充全输入、全平台或真实平台交付证明。

#### 19.5.6 错误记录、回滚与停止位置

- 原 003/004/008/018 的 lint/fix/覆盖率条目已据本轮真实通过结果关闭，历史失败正文保留。019（加载批次）、020（组合查询状态）按实际过程记录并处置；021 为修前 Format 故障反例，按修后证据关闭。
- 022 为临时分析结果与命令记录同名覆盖：初次退出 0 但产物 schema 为元数据，未当作有效数据文件。原始基准输出始终保留；另名 analyze-first 重建并核对 17/17/10 组、440 样本，首批副本与独立元数据均存在，已关闭。最终 analyze-final 与数据 benchmark-analysis.json 分名，回读为 780 样本；未因此重跑已成功的基准或修改统计方法。
- **023 保持 open**：20 轮之后仍有 7 组稳定耗时未确认；没有降预算、筛样本或自动改性能实现。该条不是新增功能回归的判定，也不表示已确认某个运行性能根因。
- 回滚按 19.4，只逆转本轮 Format、三个测试文件和必要文档差异，参照本轮 before；进入本轮前六份工作区改动保留。错误证据保留，无存储/配置/依赖/导出 API 迁移；没有实际执行回滚，不以备份存在冒充回滚验证。

**停止位置：G0～G5 的获准执行、记录与回读已结束。6 项 lint、1 处 fix、覆盖率不足已解决；稳定耗时仍有 7 组未确认，整体不宣称全部门禁通过。** 本轮未提交、推送、打 tag、发布、首次 v2、改变目录或规则、启动 OpenSpec、修改外部消费者或连接生产。


## 20. 七组性能未确认项的排查与复核（本轮获准，先方案后验证）

### 20.1 来源、事实、目标与非目标

- 来源：用户要求处理 New、Wrap、Wrapf、WrapStatus、ConvertNil、ConvertPlain、ConvertItem 七组，先补排查与复核方案并回读，区分真实回退与环境波动；不降低门槛、不筛选样本、不提交或发布。
- 本轮已实际读取最新全局/项目 AGENTS、错误记录、原计划、go-performance 及 Benchmark 参考、项目 skills-index/pkg-testing、七组夹具和原始分析。read-error-journal.sh 仅显示前 240 行，不能代替本轮已另行读取的完整错误记录。ERR-20261021-023 仍开放。
- 本轮 HEAD 实测仍为 f6a30726b372b3dc07feda42e161d27de6cb5d77；Go 为 /usr/local/bin/go 的 go1.27.2 darwin/amd64。已将 25 份实际文件及状态保存至 /private/tmp/goerr-perf-review-tnpiuhvy/before，指纹和三个独立命令退出 0 见 before-sha256.json、preflight.json。进入本轮已有十份工作区差异，均须保留。
- 第 19 节证据 /private/tmp/goerr-gates-p3pz_g87 有 20 轮、780 个样本，正常路径分配不增、最大中位回退 +3.80%；七组的保守上界 +12.58%～+29.81%，不是已经确认的回退。系统负载 3.85～32.67，频率/温度和具体后台来源未确认。
- 已读 diff 显示，七组所用 New/Wrap/Wrapf/WrapStatus/Convert 及 newItem/wrapItem/wrapStatusItem/stackFor/callers 等函数体相对原 commit d1f31a90c7080fffce42fc591ca4003a744a0eef 未改；生产差异在 Format、WithStack、diagnostic.go 和版本元数据。源码相同不能证明运行耗时等价，尚须核对编译结果、运行时与测量波动。
- 目标：保存可复核的对照、所有采样和统计证据，按原预算逐组给出通过、超预算或未确认，并明确哪些证据支持环境/测量波动、哪些支持版本相关差异。绝不把“无显著差异”“源码未改”当通过证明。
- 非目标：不优化运行代码，不添加缓存/池化/抽象/配置/出口接口，不改 Benchmark 函数体、输入、b.Loop、ReportAllocs、module/import、Go 声明、Version、日志/JSON 边界、目录或依赖；不重跑与本次无关的 lint、race、覆盖率或完整构建门禁，不声称前轮质量结果是本轮新执行。

### 20.2 范围、兼容、环境与证据

| 位置 | 具体工作与兼容约束 |
| --- | --- |
| OPTIMIZATION_PLAN.md | 当前授权/方案、步骤状态及真实结果；只修改总览与新增本节，第 13～19 节原始结果完整保留 |
| .harness/error-journal.md | 新的实际失败先记录；023 依据证据更新处置，只有缺口真实关闭才关闭，不覆盖历史失败 |
| /private/tmp/goerr-perf-review-tnpiuhvy | 本轮 before、原始历史输入副本/指纹、old/new 临时模块、固定二进制、执行记录、原始样本、统计及核对；临时产物不是永久归档，关键结果同时写入本节 |
| 其余仓库文件 | 所有 Go/测试/示例、README/CHANGELOG/CONTRIBUTING、模块、版本、规则/工作流、目录保持本轮开始时字节一致；不操作 Git 索引或外部系统 |

- old 取真实原 commit 的六个生产文件及模块；new 取当前七个生产文件及模块。共同夹具沿用阶段 1 final-benchmark_test.go，SHA-256 为 63a9170c87471187908ec9e95fe7f6a6372253c3d22bb19d11ad5bb832eb3319，先核对与当前共同部分一致；不重写夹具。
- 两个临时模块的 module 均为 github.com/gtkit/goerr、go 1.27.0、空 go.sum、无 require；仅 old 临时副本将原 go 1.26 声明对齐，不修改仓库，也不当作旧 Go 兼容测试。
- 沿用固定 Go/benchstat，重新记录版本、所选 Go 环境与 CPU/供电；GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off、证据目录 GOCACHE。不读取完整环境、凭据或进程参数，不安装工具。参数能力先由本机 Go 文档核实。
- 两边先按相同 -trimpath -pgo=off 编译测试二进制；仅包含原共同夹具，不带仓库其他测试/覆盖/race 插桩。构建完成后才采样，测量期间无并行构建或其他本任务检查。
- 二进制直接调用使用已核对的 -test.run=^$、七组精确 -test.bench、-test.benchmem、-test.benchtime=500ms、-test.count=1、-test.cpu=12；与第 19 节计时/CPU/夹具条件一致。每次只测一个已指定组，以缩短同组 old/new 间隔。只选这七组是本轮事先确定的目标范围，不按样本结果剔除基准。
- 每个命令独立保存 argv/cwd/选定环境/UTC/退出码/输出；采样另保存顺序、前后系统 load、墙钟与子进程 user/system CPU 时间。每个大块间记录 pmset 供电。未取得的温度/频率保持未知，不以 load 或 CPU/墙钟比推定某个进程或硬件原因。
- 数据文件与命令记录分别命名（result-* 与 exec-*），已有路径拒绝覆盖。原 780 样本、首批/最终结果与 benchstat 均保留并核对指纹。

### 20.3 预先固定的测量与判定（已按此执行，方法保留）

1. **固定一次复核，不按中间结果追加。** 24 个 A/B 主块，每块每组四次：奇数块 old→new→new→old，偶数块 new→old→old→new；七组顺序按块号循环移位。每三主块插入一个 A/A 块，共 8 个；其顺序也交替 ABBA/BAAB，但两个标签都执行同一个 old 二进制，文件 SHA 必须一致。主实验每组每侧 48 样本，A/A 每标签 16 样本，合计 896 个新样本。预先保存完整顺序清单，全部完成后才作最终统计；不会因为中途已达标而提前结束，不因高负载丢弃/等待/重跑某一轮。
2. A/A 能直接观察同一程序在本次协议下的时间波动与标签/顺序效应，不能证明所有版本差异都来自环境。主实验和 A/A 分开，不将 A/A 混作 old 或 new 扩充样本；不从 A/A 推导一个补偿因子来修正主实验。
3. 原第 6.4/19.3 的预算与方法保持：正常路径所有新样本 B/op、allocs/op 相对基线均不得增加；样本中位数 new/old≤1.10。对 n 个等量样本取最大 k，使 Bin(n,0.5) 的 P(X≤k−1)≤0.025，区间下界=new[k]/old[n−k+1]、上界=new[n−k+1]/old[k]（一基索引）。逐组在独立/平稳假设下，两个单侧界的并集界覆盖率至少 95%；上界≤1.10 才通过，下界>1.10 才超预算，跨界未确认；不改为只看 p 值或点估计。
4. 并列报告原 20 轮、此次 48/侧及原 20+新 48=68/侧的全部样本结果。相同源码/夹具/工具/计时条件经核对后才合并；合并是描述性复核，不假称跨时段平稳性已证明。此次与合并区间均满足预算、分配不增且无明显顺序/阶段冲突时，才将对应缺口记录为本机测量范围内通过。旧 100ms 等不同协议数据另保留，不偷偷混入或删除；如果来源不一致则停在未确认并说明。
5. 另算每主块四次中的 new 中位/old 中位比值，保存全部 24 个块比值及非参数区间；A/A 同算 8 个。这是抵抗短期共同漂移的辅助估计，估计对象不同，不代替第 3 条的门槛。首/后半、先/后标签与负载相关性只作预定诊断，不作为样本筛选规则或独立显著性主张。全部原始值、benchstat 输出、警告和方法限制保留。
6. **区分原因**：若同一二进制 A/A 也有大幅度/阶段性漂移，或编译路径相同而标签效应随次序改变，可确认存在环境/测量噪声，但不据此排除较小版本效应；若新实验与累计数据的下界均>1.10、块比值也持续超界，则记可复现的预算回退，再据符号/反汇编定位版本差异。无法区分则保持未确认。需要温度/频率、独占机器或运行代码优化才能解释的，不猜测或擅自调系统。
7. 若出现上述可复现回退，仅对受影响组再做一次 old/new 的 5 秒 CPU profile 与本地 pprof/top；profile 样本不混入计时验收，不开启服务或网络监听。无该证据不采 profile，不以诊断授权顺带改代码；任何修复另需先补具体方案和授权。

### 20.4 步骤、验收、失败与停止条件

| 步骤 | 前置条件 | 具体工作 | 检查/可复核证据 | 状态与停止条件 |
| --- | --- | --- | --- | --- |
| P0 方案落盘回读 | 最新规则、原计划、错误、七组原始输入已读 | 保存总览与本节，回读目标/门槛/固定顺序与样本数/白名单 | before/preflight、实际文档回读；运行代码和夹具未改 | 已保存并完整回读；plan-readback.json 核对 25 文件仅计划变化，历史 13～19 节一致；P0 完成，未先采样 |
| P1 来源与环境诊断 | P0 完成 | 归档原数据/指纹，核对来源/函数调用与工具；构建固定 old/new，记录符号/关键反汇编并作保留地址/重定位区别的解释 | source/tool/env、构建 argv/退出码、二进制 SHA、nm/objdump、历史原始分布；源码相同不冒充机器码/性能证明 | 已完成；归档 105 份历史文件，11 个源函数体一致；两个构建和工具命令退出 0，14 对关键符号归一化指令一致；不据此当性能通过 |
| P2 固定对照采样 | P1 完成、计划顺序文件保存，未并行跑检查 | 完整 24 主块+8 A/A，每组四次，load/CPU 时间及供电伴随记录 | 32 块×7 组×4 次，合计 896 个单样本命令；主实验 672、对照 224；全部退出码、时间与顺序 | 已完成；24 主块+8 A/A、896 个样本均按固定顺序完成且退出 0，主实验 672/对照 224，无筛选或补采 |
| P3 全样本分析与归因 | P2 样本完整且输入指纹一致 | 按 20.3 分析历史/新/累计和 A/A、块比值/预定分层，运行现有 benchstat；必要时才按条件作 profile | result-* JSON 与 exec-* 元数据分名；各组原始值、区间/分配/中位数、顺序/阶段/A/A 及限制 | 已完成；历史/本轮/累计/A/A、28 分层及三份 benchstat 全读回；七组原保守上界通过，未触发 profile 条件 |
| P4 写回与回读 | 可执行的诊断完成，结果和失败已实际读取 | 更新本节各步骤真实状态、命令/数值、原因与未知；核对 25 文件、历史段落/数据和 HEAD，处置 023，回读最终方案 | 指纹/命令表/样本计数和审计、文档实际回读；其余 24 文件保持（错误记录有真实处置时除外） | 已完成；独立审计、方案/结果回读及指纹核对通过，023/024/025 据证关闭；最终证据见 result-docs-readback.json，不再跑测试或构建 |

### 20.5 风险、失败处理、回滚与最终验收

- 风险为微基准时间短、进程初始化/GC/代码布局、负载/温度/频率漂移与时间相邻相关性。ABBA/次序轮换/固定二进制减少部分混杂但不保证独立平稳；A/A、分层及原始值提供实测限制。新数据量增加不能替代对假设的说明。
- 全样本方案固定后不剔除慢样本、异常值或不利轮次，不降低 10%/零新增分配要求，也不将新算法、标准库/编译器升级或不同 CPU/GC 条件混作同条件优化收益。
- 命令失败先追加真实错误并停止依赖步骤，区分编译/测量/分析工装失败与库性能回退；工装修复不删原始输出，最多三轮且要有正进展。基准机器仍不稳定时如实交付未确认及下一步所需条件，不无限采样。
- 源码及 API/错误/栈/状态/日志/JSON/可插拔出口、模块/依赖/Go/Version 均保持；无存储迁移和发布回滚需求。本轮可逆内容仅本计划、错误处置与临时证据：如需撤回只逆转本轮文档差异，依据 before，保留进入本轮前的十份改动及真实错误证据；临时文件留存，不批量删除，不执行 git reset/clean。备份存在不是回滚已实测。
- 验收包含：方案先保存回读；输入与工具真实一致；896 新样本固定计划全保留；所有七组同时展示此次和累计的分配/点估计/原保守区间及 A/A 解释；每项通过/超预算/未确认有依据；真实结果与未知写回并回读；仅白名单文件变化，无提交或发布。没有满足统计证据的组继续保持未确认，不能以“完成了命令”代替通过。

### 20.6 实际结果（2026-10-21）

#### 20.6.1 实际执行、来源与命令

- P0 先保存并完整回读第 20 节，plan-readback.json 证明当时仅计划变化；第 13～19 节原文保留。P1 已归档 105 份前轮文件，其原件和副本 SHA 均未变；旧源码来自真实 commit，新源码与本轮开始及第 19 节候选逐文件一致，夹具 SHA/行号/输入/函数体未改。11 个相关源函数体一致，14 对已保留的关键符号在仅归一化局部跳转目标后指令与其他操作数一致；New/Wrap/WrapStatus 在基准中的内联路径也包含于核对。不是整个二进制的字节或性能等价证明。
- 两个临时模块均 github.com/gtkit/goerr、go 1.27.0、空 go.sum；实际 Go 1.27.2、darwin/amd64、GOAMD64=v1、CGO_ENABLED=1，GOROOT=/usr/local/go、GOWORK/GOEXPERIMENT 空。benchstat 仍为 x/perf v0.0.0-20260908200009-22c9c6c9d4da（由 Go 1.27.1 构建）；没有下载或升级。old/new 固定二进制均 7,076,688 字节，SHA 分别为 `1ebab032b69335360509c8c15cb5bcb431202345a8d8402ea68c8d114fc5267c`、`1980c1a391f47a7a1044fd70d1af680487fcad52bdca97514adaacabd8ca9b94`；A/A 两个标签始终执行同一 old.test。
- P2 于 2026-10-21T11:50:14.985068+00:00～12:00:20.556732+00:00 实际完成全部 24 主块与 8 对照块。主实验 7×48×2=672、A/A 7×16×2=224，共 **896** 个单样本命令，全部退出 0，顺序与预先保存的 schedule.json 逐项相同。无中途停止、追加、删值、择时或重跑；原 780+新 896=**1676** 个样本保留。
- CPU 实际头部为 Intel Core i7-9750H @2.60GHz，cpu=12、500ms、同一夹具/构建参数。896 个元数据中 GOGC/GOMEMLIMIT/GOMAXPROCS/GOEXPERIMENT/GOAMD64 环境变量均未设置；GOAMD64 实际默认值来自 go env。采样前后 1 分钟 load 为 3.66～21.12，33 次离散供电观察均为 AC；未连续监测供电/温度/频率，也未读取其他进程参数，不据此认定具体根因。
- 证据目录 `/private/tmp/goerr-perf-review-tnpiuhvy`。result-performance.json 保留全部新样本、历史/本轮/累计/A/A 与辅助统计；result-verification.json 核对原始拼接、顺序、历史/源码和计数。以下 **948** 份完成命令元数据均退出 0，命令逐项 argv/cwd、选定环境、UTC、CPU/墙钟与输出见各 exec-*.json/.txt；批量样本和供电合并展示。独立审计运行前为 946 份，之后含审计及 023 关闭为 948 份。另有初始只读/来源命令和一次解析失败的工具记录，不能把 948 份成功元数据表述为本轮从未失败。

| 记录名 | cwd（root=项目，evidence=本轮目录） | 实际命令 | 退出码 |
| --- | --- | --- | --- |
| exec-analyze | root | `python3 /private/tmp/goerr-perf-review-tnpiuhvy/analyze.py` | 0 |
| exec-benchstat-aa | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-perf-review-tnpiuhvy/result-aa-A.txt /private/tmp/goerr-perf-review-tnpiuhvy/result-aa-B.txt` | 0 |
| exec-benchstat-main | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-perf-review-tnpiuhvy/result-main-old.txt /private/tmp/goerr-perf-review-tnpiuhvy/result-main-new.txt` | 0 |
| exec-benchstat-pooled | root | `/Users/xiaozhaofu/go/bin/benchstat /private/tmp/goerr-perf-review-tnpiuhvy/result-pooled-old.txt /private/tmp/goerr-perf-review-tnpiuhvy/result-pooled-new.txt` | 0 |
| exec-benchstat-version | root | `/usr/local/bin/go version -m /Users/xiaozhaofu/go/bin/benchstat` | 0 |
| exec-build-new | evidence/new | `/usr/local/bin/go test -c -o /private/tmp/goerr-perf-review-tnpiuhvy/new.test -trimpath -pgo=off .` | 0 |
| exec-build-old | evidence/old | `/usr/local/bin/go test -c -o /private/tmp/goerr-perf-review-tnpiuhvy/old.test -trimpath -pgo=off .` | 0 |
| exec-close-023 | root | `bash .harness/scripts/close-error-journal.sh . ERR-20261021-023 '第20节固定896样本与原780全保留；7组本轮48/侧和累计68/侧均按原保守上界≤1.10通过，最大+3.24%/+3.25%，分配不增；28个分层也通过。分析、benchstat、独立审计均退出0，仅关闭本机条件范围的统计缺口，具体环境根因仍未确认；未改代码/门槛或发布。'` | 0 |
| exec-close-024 | root | `bash .harness/scripts/close-error-journal.sh . ERR-20261021-024 '修正仅针对真实 :-1 来源字段；exec-compiled-paths 退出 0，完整解析 14 对符号，保留填充/操作数和重定位差异，原始输出未改。896 项采样顺序已固定，未产生或丢弃任何样本。'` | 0 |
| exec-compiled-paths | root | `python3 /private/tmp/goerr-perf-review-tnpiuhvy/compiled_paths.py` | 0 |
| exec-doc-nm | root | `/usr/local/bin/go doc cmd/nm` | 0 |
| exec-doc-objdump | root | `/usr/local/bin/go doc cmd/objdump` | 0 |
| exec-env | root | `/usr/local/bin/go env -json GOVERSION GOOS GOARCH GOAMD64 GOROOT CGO_ENABLED GOMOD GOWORK GOEXPERIMENT GOTOOLCHAIN` | 0 |
| exec-nm-new | root | `/usr/local/bin/go tool nm -size /private/tmp/goerr-perf-review-tnpiuhvy/new.test` | 0 |
| exec-nm-old | root | `/usr/local/bin/go tool nm -size /private/tmp/goerr-perf-review-tnpiuhvy/old.test` | 0 |
| exec-objdump-new | root | `/usr/local/bin/go tool objdump -s '^github.com/gtkit/goerr\.(Benchmark(New\|Wrap\|Wrapf\|WrapStatus\|ConvertNil\|ConvertPlain\|ConvertItem)\|New\|Wrap\|Wrapf\|WrapStatus\|Convert\|newItem\|wrapItem\|wrapStatusItem\|stackFor\|callers)$' /private/tmp/goerr-perf-review-tnpiuhvy/new.test` | 0 |
| exec-objdump-old | root | `/usr/local/bin/go tool objdump -s '^github.com/gtkit/goerr\.(Benchmark(New\|Wrap\|Wrapf\|WrapStatus\|ConvertNil\|ConvertPlain\|ConvertItem)\|New\|Wrap\|Wrapf\|WrapStatus\|Convert\|newItem\|wrapItem\|wrapStatusItem\|stackFor\|callers)$' /private/tmp/goerr-perf-review-tnpiuhvy/old.test` | 0 |
| exec-sampling | root | `python3 /private/tmp/goerr-perf-review-tnpiuhvy/sample.py` | 0 |
| exec-verify | root | `python3 /private/tmp/goerr-perf-review-tnpiuhvy/verify.py` | 0 |
| exec-sample-main-*（672） | evidence/old 或 new | `<old/new.test> -test.run=^$ -test.bench=^Benchmark<组名>$ -test.benchmem -test.benchtime=500ms -test.count=1 -test.cpu=12` | 全 0 |
| exec-sample-aa-*（224） | evidence/old | `old.test` 及相同参数；A/B 标签均为同一文件 | 全 0 |
| exec-power-*（33） | root | `/usr/bin/pmset -g batt` | 全 0 |

初始 preflight 三个独立命令（git rev-parse、git status、go version）及八次 git show 均退出 0，实际 argv/输出/指纹分别保留于 preflight/source-provenance。首次内联 Python 反汇编解析确实退出 1，详情见 20.6.5；没有伪造其独立 UTC 元数据。


#### 20.6.2 原门槛的七组验收

以下均由真实原始值计算，Δ=(new/old−1)×100%。本轮每侧 n=48、k=17，累计保留原 20+本轮 48，即 n=68、k=26；对应每个**单侧比值界**的覆盖率下限约 97.07%/96.15%，前提仍是独立/平稳。算法、0.025 单尾与 **≤10% 耗时回退、B/op/allocs/op 不增加**预算均与第 19 节一致。完整双侧区间、各组同时和跨数据集同时覆盖不声称有同一个 95% 保证；没有用 `p>0.05` 证明等价。

| 组 | 本轮 old→new ns/op | 本轮中位 Δ | 本轮上界 Δ | 累计 old→new ns/op | 累计中位 Δ | 累计上界 Δ | 两边全部 B/op / allocs/op | 判定 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| New | 998.3→1004.5 | +0.62% | +2.61% | 994.5→1003 | +0.85% | +3.15% | 136 / 3 | 通过 |
| Wrap | 94.875→95.33 | +0.48% | +2.27% | 94.585→95.305 | +0.76% | +2.45% | 104 / 2 | 通过 |
| Wrapf | 219.95→220.65 | +0.32% | +1.28% | 219.35→220.4 | +0.48% | +1.65% | 120 / 3 | 通过 |
| WrapStatus | 94.965→94.76 | -0.22% | +1.58% | 94.595→94.615 | +0.02% | +1.78% | 128 / 2 | 通过 |
| ConvertNil | 2.2935→2.299 | +0.24% | +3.24% | 2.29→2.2875 | -0.11% | +3.25% | 0 / 0 | 通过 |
| ConvertPlain | 43.94→43.72 | -0.50% | +2.59% | 43.835→43.44 | -0.90% | +2.13% | 80 / 1 | 通过 |
| ConvertItem | 6.21→6.206 | -0.06% | +2.48% | 6.1825→6.172 | -0.17% | +2.36% | 0 / 0 | 通过 |

七组本轮最大中位回退 **+0.62%**、最大上界 **+3.24%**；累计最大中位回退 **+0.85%**、最大上界 **+3.25%**。所有样本的 B/op、allocs/op 数值均相同，非仅中位数相同。因此这 **7 组原未确认项在本机固定协议范围内全部通过**；不将其外推为所有平台、真实服务延迟/吞吐或发布证明。

benchstat-main 的七组均显示 `~`、n=48，old/new 时间区间约 ±1%～3%；pooled 的七组 n=68，约 ±1%～3%。这仅与主计算交叉核对，验收仍依据上述保守上界。pooled 输入完整保留原 17 组：其他 10 组仍只有历史 20/侧，并非本轮复测；其中 WithStackMissing 为前序新增捕栈成本，0→112 B、0→2 alloc，时间 geomean +21.36% 含此项，不能当七组/正常路径回退。100ms 和 Inspect 原始记录保留但不混入 500ms 门槛样本，A/A 不混入主实验。

#### 20.6.3 A/A、顺序与环境诊断

A/A 每标签 16 个样本，n=16/k=4；以下 A、B 都是同一个 old.test，不是不同版本。区间只是逐个方向的原方法比值界；原值范围包含所有慢样本。

| 同一二进制组 | A→B 中位 ns/op | 标签中位 Δ | 比值下/上界 Δ | A 原始范围 ns/op | B 原始范围 ns/op |
| --- | --- | --- | --- | --- | --- |
| New | 1007.5→1003 | -0.45% | -6.23%～+5.75% | 826.3～1097 | 887.8～1070 |
| Wrap | 94.155→93.805 | -0.37% | -4.19%～+2.66% | 88.11～100.8 | 89.63～103.5 |
| Wrapf | 219.4→221.35 | +0.89% | -2.74%～+3.81% | 202.4～228.7 | 202.3～225.4 |
| WrapStatus | 94.565→94.365 | -0.21% | -4.13%～+5.38% | 86.67～98.48 | 91.07～110.8 |
| ConvertNil | 2.3045→2.2895 | -0.65% | -4.05%～+3.89% | 2.204～4.352 | 2.22～2.75 |
| ConvertPlain | 43.495→43.93 | +1.00% | -2.92%～+4.43% | 37.77～59.17 | 41.12～50.98 |
| ConvertItem | 6.1505→6.3 | +2.43% | -2.64%～+5.69% | 5.558～6.965 | 5.402～6.757 |

- **可确认存在测量波动。** 同一二进制 ConvertNil 的 A 标签也有 2.204～4.352 ns/op（接近两倍），New 为 826.3～1097 ns/op；未删这些值。七组 A/B 中位差为 -0.65%～+2.43%，各组界均包含 1.0，没有稳定的单向标签效应证据。这不能证明上一轮所有变化均来自环境，也不排除小于预算的版本效应。
- **未发现可复现的超预算版本回退。** 11 个源函数体/14 对关键编译路径无新增指令的证据、主实验与累计上界达标、A/A 实测噪声及下列预定分层相互一致。上一轮“未确认”属于统计精度未收敛，不能追写成“已经发现并修复了性能缺陷”。当前资料仍不能定位温度、频率或某个后台进程，也不能单独证明换协议带来的改善有多少来自环境变化。
- 主块比值的辅助估计使用全部 24 个块，n=24/k=7；其自身中位数双侧界下限约 97.73%，估计对象是“块内比值的中位数”，没有替代原门槛。预定的前/后半与 old/new 先运行四类各 n=24/侧，共 28 组，其上界也全部≤1.10，最大 +6.66%，无与主结论冲突的阶段或顺序超预算结果。

| 组 | 24 主块比值中位 Δ | 主块比值中位界 Δ | 前半 / 后半 / old先 / new先 中位 Δ | 四类最大上界 Δ |
| --- | --- | --- | --- | --- |
| New | +1.14% | -1.47%～+1.58% | +0.15% / +0.79% / -0.17% / +0.75% | +4.22% |
| Wrap | -0.10% | -1.44%～+1.66% | +0.60% / +0.32% / +1.03% / +0.27% | +3.82% |
| Wrapf | +0.23% | -1.21%～+1.76% | +0.80% / +0.20% / +0.45% / +0.25% | +2.47% |
| WrapStatus | -0.20% | -1.35%～+1.18% | -0.27% / +0.05% / +0.19% / -0.33% | +3.59% |
| ConvertNil | +0.03% | -0.88%～+1.77% | +0.22% / +0.32% / +0.61% / +0.24% | +4.35% |
| ConvertPlain | -0.32% | -2.53%～+1.45% | -1.89% / +1.75% / +0.21% / -0.66% | +5.42% |
| ConvertItem | -0.62% | -2.44%～+0.97% | -1.36% / -1.24% / +1.56% / -1.98% | +6.66% |

负载关联与 CPU 时间仅作为描述性伴随记录；Spearman 是本轮样本时间与采样前 1 分钟 load 的秩相关，未作显著性/因果推断。CPU/墙钟为整个基准子进程（含校准与 GC 等）时间比，不是单次函数耗时或 CPU 频率；大于 1 可包含多个运行时线程，不能解释成“超过 100% 的单核效率”。

| 组 | 时间/load Spearman old / new | 子进程 CPU/墙钟中位 old / new |
| --- | --- | --- |
| New | +0.171 / +0.158 | 1.021 / 1.020 |
| Wrap | +0.183 / -0.051 | 1.240 / 1.239 |
| Wrapf | -0.012 / +0.159 | 1.111 / 1.113 |
| WrapStatus | -0.071 / +0.020 | 1.292 / 1.292 |
| ConvertNil | -0.181 / -0.135 | 0.991 / 0.992 |
| ConvertPlain | -0.038 / -0.237 | 1.423 / 1.420 |
| ConvertItem | -0.211 / -0.098 | 0.992 / 0.992 |

#### 20.6.4 完整性、范围与兼容核对

- exec-verify 真实退出 0：896 项命令与固定顺序完全一致、原始聚合为精确拼接、A/A 始终同一二进制、全部样本分配不增；逐项核对 20/48/68/16 的计数和二项尾界，没有筛样本或改算法。
- 105 份历史原件与归档 SHA 未变，第 13～19 节原文逐字节保留。25 份本轮快照中仅计划和错误记录变化，其余 23 份（含全部 Go/测试/示例、README/CHANGELOG/CONTRIBUTING、模块/版本、规则/.gitignore）保持。本轮未新增/删除导出 API、Options、配置、依赖或目录，可插拔出口和日志/JSON 约定沿用现状。
- 审计内 git rev-parse/status 各退出 0；HEAD 仍为 f6a30726b372b3dc07feda42e161d27de6cb5d77，进入本轮的十份 Git 工作区差异保留。没有提交、推送、tag、发布、首次 v2、容器/外部服务或 OpenSpec 操作。
- 本轮只编译两个临时基准二进制并按方案测量，没有重跑仓库功能测试、完整 build/vet/lint/race、覆盖率或 Fuzz；第 19 节通过记录仍是前轮真实证据，不写成此次执行。没有触发 20.3 的可复现回退条件，因此未采 CPU profile、未改运行代码或开展优化。

#### 20.6.5 错误处置、回滚与停止位置

- ERR-20261021-024：首次内联 Python 解析退出 1，真实 objdump 的无文件名 `:-1` 填充行不满足原正则；已先记录，再仅放宽来源 token 解析。保留所有填充指令、PC/字节/操作数及原始输出；另名 exec-compiled-paths 退出 0，14 对符号全解析，024 据真实证据关闭。无性能样本因该问题重跑或删除。
- ERR-20261021-023：保留原 20 轮未确认正文，追加本节实际证据后，关闭脚本真实退出 0；关闭的是七组性能门槛未确认缺口，具体环境根因仍未定位。当前真实 ERR 条目无 open；没有把历史失败抹成成功。
- ERR-20261021-025：最终文档回读发现命令表中的正则竖线未转义、argv 拼接未保留参数边界；已先记录，再仅改成保留参数边界的引号并转义表格竖线。真实 argv/输出、统计与采样均未改；逐行回读命令表后关闭。
- 回滚范围仍仅本轮文档差异与临时证据，参照 before 保留进入本轮前改动；运行代码无需回滚。没有实际执行回滚、删除临时文件或改变系统配置，不把快照存在当作回滚已验证。
- **停止位置：P0～P4 已完成，七组按原门槛全部通过；方案/真实结果/错误处置已写回并实际回读核对。** 固定采样上限已到，不再追加、优化代码、重跑其他门禁或进入其他阶段；最终文档写回后只回读核对，不运行测试或构建。


## 21. 提交、v1.5.0 发版与 GitHub 推送（本轮获准，先方案后执行）

### 21.1 来源、现状、目标与范围

- 用户最新指令为“和功能无关的这些先不做，那就提交，发版，并 push 到 GitHub”。这明确替代此前不提交/发布的边界，仅授权本次既有功能的提交与发布，不授权工具链升级、CI 建设、Makefile 改造、汇总文档清理、目录调整或新增功能。
- 已实际读取最新全局/项目规则、全部错误记录、原计划、发布/兼容 Guide，以及工程治理和质量 skill。当前 main 的 HEAD 为 f6a30726b372b3dc07feda42e161d27de6cb5d77，已有一个尚未推送的诊断功能提交；进入本轮十份工作区差异为 CHANGELOG、CONTRIBUTING、本计划、README、error.go、error_test.go、resolve_test.go、sanitize_test.go、status_test.go、test/integration_test.go，索引为空，全部为本次功能与验收链的待交付内容。
- 本轮已将 26 份文件及现有差异备份至 /private/tmp/goerr-release-z_ka1_80/before；指纹、HEAD/状态/工具链的独立命令与 patch 见 before-sha256.json、preflight.json、before-working-tree.patch。
- 实际 GitHub 目标为 gtkit 远端的 gtkit/goerr；另一 ydsd 远端不在本次 GitHub 授权范围，不推送。git ls-remote --heads gtkit main 退出 0，远端 main 为 d1f31a90c7080fffce42fc591ca4003a744a0eef；远端 tags 查询退出 0，最高 v1.4.0，无 v1.5.0；本地最高 tag 也为 v1.4.0。浏览器真实页面显示当前没有 GitHub Release，已登录且可创建 Release。
- module 保持 github.com/gtkit/goerr、go 1.27.0，无外部运行依赖；Version 已为 v1.5.0。该版本承接已实现的诊断快照、补栈、Format 修正与契约测试，使用现有 v1 MINOR，不首次升级 v2，不盲目再增为 v1.5.1，不运行会自动 bump/提交/推送的 make tag。
- 当前 /usr/local/bin/go 已实际为 go1.27.3 darwin/amd64，与前轮的 go1.27.2 不同；本轮未安装或升级它。18 份 Go/测试/模块文件与前轮质量及性能证据指纹完全相同。前轮 build/vet/lint/fix/race/coverage 均真实退出 0、根包 100.0%，严格性能 7 组已按原门槛通过；这些保留为历史证据，不冒充当前工具链新执行。
- 目标：将当前功能及必要发布元数据提交到 main，发布唯一 v1.5.0 附注标签，推送至 GitHub 并创建正式 GitHub Release，核对远端对象与页面；在当前 Go 下复核必要发布检查。不修改运行代码、接口、数据/存储契约、日志/JSON/可插拔边界、依赖或最低 Go 声明。

### 21.2 文件、工具、兼容及发布设计

| 范围 | 本轮具体动作 |
| --- | --- |
| 当前十份已有差异 | 逐份审阅并显式暂存提交，只提交既有功能/测试/必要说明，不使用 git add . 扩大范围 |
| CHANGELOG.md | 仅将已有 Unreleased 内容整理为 [v1.5.0] - 2026-10-21；保留空 Unreleased 入口、原历史版本及变更正文 |
| OPTIMIZATION_PLAN.md | 仅更新当前授权/状态和追加本节；保留第 13～20 节历史；写回真实结果与回读 |
| version.go / go.mod / go.sum / 全部 Go 文件 | 保持本轮开始时字节不变；Version 已满足目标，不额外 bump，不新增代码或导出 API |
| .harness/error-journal.md | 仅真实失败先记录、据证处置；不强制加入被忽略的本地错误记录，不改规则 |
| 临时证据目录 | 原始检查/命令/退出码、版本与差异核对、Release/标签/提交说明和最终回执；不提交构建、覆盖率或日志产物 |
| Git / GitHub | 仅本仓库 main 与 v1.5.0 标签及对应 Release；不动其他远端/分支/已有标签，不强推或改写历史 |

- 现有 Makefile 只有 tool/check/tag，Guide 描述的 release-* 入口在本项目并不存在；本轮不猜命令或修脚本，使用显式、逐步检查并停止的 Git 命令。tool 会改写格式，tag 会自动 bump 与推送，均不调用。
- Go/已有 lint 工具先核对版本，GOTOOLCHAIN=local、GOFLAGS=-mod=readonly、GOPROXY=off，缓存和覆盖产物放证据目录；不下载工具、模块或镜像，不跑新的安全扫描/供应链检查，不把用户暂缓项标成通过。
- 当前工具链已变化，必要复核为 build、vet、lint、只读 go fix -diff、完整 race/count=1、覆盖率及 cover -func，以及现有 benchmark/count=3。Benchmark 使用 -run=^$ 排除重复功能用例，原始输出和分配完整保留；三样本只作本工具链基准执行检查，不把它或历史 Go1.27.2 的严格时间界写成 Go1.27.3 下新完成的 10% 比较。不开新一轮性能排查、加采样或优化。
- 兼容以既有契约与冻结指纹为据：最小 Go 1.27.0，保留业务码/HTTP/文案/cause、首个 Item、已有栈/typed-nil、快照所有权和普通函数出口；不承诺 Go1.26、目标生产平台或真实日志接收已验证。既有 AST/契约证据保留，不安装 apidiff 或扩展接口来完成发版。
- 发布说明来自当前版本 CHANGELOG，并注明最低 Go、当前真实检查与历史性能证据范围；不把本机验证写成生产验收，也不发布临时路径、错误记录或内部执行日志作为附件。
- 本机没有 gh 和可用的 GitHub 发布专用工具；已有登录 Chrome 的创建 Release 页面可用，因此 Git 用 CLI，GitHub Release 用已授权浏览器 UI，不读取/回显凭据，不更改认证/权限或用户全局配置。
- 提交顺序预先固定：先提交功能/发布元数据（C1），在干净工作区创建附注 tag 并向 gtkit 原子推送 main+v1.5.0；若远端不支持原子推送，先记录原失败且确认远端未变化，再按 main→tag 的正常快进顺序推送。随后创建指向既有远端 tag 的正式 Release。最后把实际发布结果写回本节、回读，单独提交必要发布回执文档（C2）并快进推送 main；已发布 tag 始终指向 C1，不移动到 C2。C2 自身哈希/推送的最终结果在外部回执和最终汇报记录，避免文档自引用或谎报未执行动作。

### 21.3 分步计划、验收与停止条件

| 步骤 | 前置条件 | 具体动作 | 可复核证据/验收 | 当前状态与停止条件 |
| --- | --- | --- | --- | --- |
| R0 发布方案 | 最新指令/规则/原计划/错误记录与远端已读 | 保存本节并完整回读，冻结 26 文件、18 源文件及十份差异清单 | 实际方案回读、before/preflight、版本/分支/目标来源 | 已完成：plan-readback.json 记录全文回读与冻结核对；事实/授权/目标不符即停，不先提交再补计划 |
| R1 发布元数据与差异 | R0 完成，v1.5.0 确认不存在 | 只整理 CHANGELOG 版本/日期，审阅现有全部 diff，核对显式文件清单；git diff --check，必要时仅清本轮已改文档行的非语义空白 | CHANGELOG/Version/module 一致，无泄密/额外文件；运行/测试/规则/Makefile 指纹不变 | 已完成：版本/日期及十份差异已核对，空白修正后 diff --check 退出 0；出现新运行差异、版本冲突或未知文件即停 |
| R2 当前工具链复核 | 元数据与代码范围已冻结、Go/lint 入口核实 | 按 21.2 运行必要发布检查，保存独立 argv/cwd/所选环境/UTC/退出码/完整输出，回读覆盖率与基准样本 | build/vet/lint/race/coverage 通过，fix 输出空，根包≥80%；基准 count=3 实际完成且分配核对，历史严格区间明确分开 | 已完成：八项命令退出 0；330 用例/6 Example/两包通过，384/384 覆盖，27 组各三次基准完成；跨 trimpath 观察已核对；失败先记录并停依赖发布步骤，不降门槛或改代码 |
| R3 提交与远端 tag | R1/R2 通过，远端仍无冲突，当前文件与验证指纹一致 | 显式暂存十份差异、回读 staged diff、创建 C1；干净后创建 v1.5.0 附注 tag，正常原子推送 main+该 tag 到 gtkit | 提交对象/文件清单、tag 类型/剥离目标、push 退出 0；ls-remote 精确验证 main/tag 与 C1 一致 | 已完成：C1=9a6f745，附注标签及 main 原子 push 退出 0，远端精确一致；不强推/删重打/擅自合并 |
| R4 GitHub Release | main 与 tag 已真实推送且一致 | 在已登录仓库页面选择现有 v1.5.0 标签，填写已核对中文说明，发布正式 Release，回读页面并保存可核对截图 | 正式 Release URL/版本/日期/标签目标及页面可见结果，无草稿或预发布误标 | 已完成：正式 Release v1.5.0 页面显示 Latest、标签与 C1，一并保存截图和页面证据 |
| R5 实际结果与交付 | 已完成的发布动作有真实证据 | 只写本节的实际检查与 C1/tag/push/Release 状态，回读；提交发布回执文档 C2 并仅推送 main，最终核对工作区/远端/tag/Release | 文档回读与受保护指纹、C2/push 实际结果、干净状态、固定 tag、链接与保留的未知项 | 实际发布结果已落盘并回读；按本步骤提交/推送 C2，C2 自身结果由最终外部回执确认，避免自引用；不补做暂缓项，不移动 tag 或部署生产 |

### 21.4 风险、失败处理与回滚

- 远端可能在查询后变化；只做正常快进/原子 push，拒绝就停并读取真实状态，不 force、reset/clean、stash 用户工作、删标签或重打既有版本。
- 版本标签发布后可能被模块代理永久缓存；既有 tag 不删除/覆盖。发布后若发现问题，调用方锁回上一已验证版本/制品；仓库修复需后续明确范围并发新 PATCH，撤回 Release/标签或改历史不在本次自动操作范围。本轮不是生产部署，也不宣称执行了生产回滚。
- C1 前的本轮新增内容只涉及 CHANGELOG 版本区段与计划，撤回时仅按 before 逆转这部分，保留进入本轮前十份差异；提交失败可保留索引和工作区，推送失败保留本地提交/tag 便于续做，不用破坏性清理。
- 浏览器权限/网络失败与 Git 推送分别报告；Release 未发布时不能仅凭标签说正式 GitHub Release 已完成。所有命令失败先记真实错误、保留输出；同一问题最多三次且必须有正进展，无法推进时说明具体缺口，不自行安装 gh 或读取凭据。
- 用户暂缓的 Go 升级、CI、脚本改造、汇总文档和目标环境/安全补充保持原状；当前 Go1.27.3 是已观察到的外部环境变化，不冒充本轮完成了升级。已有文档中的历史未验证状态保留，由本节提供当前发布事实，不扩大为文档清理。

### 21.5 实际结果

- R0/R1 已完成：方案实际全文回读，版本、远端、文件清单与差异已核对；CHANGELOG 整理为 v1.5.0 / 2026-10-21，Version 和 module 未改。
- 初次 git diff --check 退出 2，仅本轮计划头部两行 Markdown 行尾双空格；已先记录 ERR-20261021-026，仅去除这两处空白后复核退出 0、输出空，条目已关闭。原失败与成功输出均保留于证据目录，不写成首次全过。
- 当前 Go 环境及 lint 版本核对均退出 0：go1.27.3 darwin/amd64；golangci-lint 2.14.0（构建 Go1.27.2）。R2 发布检查已完成；R3/R4 提交、标签推送和正式 Release 已完成，详见下方真实回执。

- R2 本轮真实命令（均在项目根目录、21.2 的离线/只读模块环境执行）：`go build ./...`、`go vet ./...`、`golangci-lint run ./...`、`go fix -diff ./...`、`go test -json -race -count=1 -timeout=5m ./...`、`go test -count=1 -coverprofile=/private/tmp/goerr-release-z_ka1_80/coverage.out ./...`、`go tool cover -func=/private/tmp/goerr-release-z_ka1_80/coverage.out`、`go test -run '^$' -bench . -benchmem -count=3 ./...`，八项均退出 0。lint 为 0 issues，fix 原始输出为空；race 330 个用例通过（含 6 个 Example），两包通过；根包 384/384 语句覆盖，100.0%。逐项 argv/cwd/环境/UTC/退出码/原始输出见 exec-build/vet/lint/fix/race/coverage/cover-func/benchmark 的 json/txt。
- Benchmark 共 27 组、81 个样本完整保留，实际耗时 113.740 秒；分配次数均未增加。FormatVerbose 为 744 B/op、12 allocs/op，历史隔离对比为 648/12，但真实历史命令使用 -trimpath，本轮未使用；该路径会格式化栈文件名，不能把跨编译参数字节差异判成同条件回退。已读取原工作区十个历史样本，均为 744/12，与本轮三样本一致。该观察、两套原始值与命令均保留，不筛样本或改门槛；本轮未重新裁定 Go1.27.3 下严格 10% 耗时比较。详见 result-release-checks.json，原 Go1.27.2 严格结论仍限定第 20 节条件。
- 辅助汇总曾读取未确认的历史文件名而退出 1（ERR-20261021-027），已先记录并改用实际发现的证据；当前 Go 检查均为真实退出 0。跨 trimpath 观察记录 ERR-20261021-028，按命令和工作区原始样本核对后关闭，不修改源码。23 份受保护文件、历史第 13～20 节均未变化，当前十份待提交清单不变。

- R3 已完成：显式暂存上述十份文件，逐文件索引内容与已审阅工作区一致，原始 staged-release.patch 与 result-staged-review.json 已保存。`git commit -F commit-message.txt` 退出 0，C1 为 `9a6f7459bd5c74c6b83056ffb484784a6efae99b`，提交后工作区干净。`git tag -a v1.5.0 -F tag-message.txt` 退出 0，附注对象 `94f74e795fd9123a5f8015c354b56e5cd5840bf9`，剥离目标为 C1。`git push --atomic gtkit refs/heads/main:refs/heads/main refs/tags/v1.5.0:refs/tags/v1.5.0` 退出 0；后续 ls-remote 退出 0，main、tag 剥离目标均精确等于 C1，附注对象也相等。未推送 ydsd，未强推或重打标签。
- R4 已完成：真实 GitHub UI 选择已存在 v1.5.0 标签，填写本版 CHANGELOG 和实际验证说明，点击 Publish release 后页面为 [v1.5.0 正式 Release](https://github.com/gtkit/goerr/releases/tag/v1.5.0)，显示 Latest、released this、2026-10-21 21:04 GMT+8、Tag v1.5.0 和 Commit 9a6f745；不是草稿或预发布。已回读发布正文，保存 github-release-state.txt、github-release.json 与 github-release.jpg，页面保留为交付入口。
- R5 文档回执：本次仅写回本节和当前发布状态；第 13～20 节及源码/测试/依赖/Version/Makefile/规则继续保持冻结。回读后按既定顺序创建仅含本计划的 C2 并仅快进推送 main；此处记录的是提交前已核实事实，C2 自身哈希/推送退出码与最终工作区/远端核对由证据目录的最终回执及用户汇报确认，已发布 v1.5.0 继续固定 C1。本轮未修改运行代码、安装/升级工具、建设 CI、改造 Makefile、补做安全扫描或连接生产。
