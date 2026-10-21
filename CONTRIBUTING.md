# 参与贡献与发布检查清单

本地开发与检查使用 **Go 1.27.0 或更高版本**，与 go.mod 的最低版本一致；Go 1.26 工具链不满足当前未发布代码要求。

## 本地测试（建议合并前全部执行）

```bash
go test ./... -count=1
go test ./... -race -count=1
go vet ./...
go test -fuzz=FuzzClampMessage -fuzztime=3s ./...
```

## 测试矩阵（核心包 `github.com/gtkit/goerr`）

| 区域 | 内容 |
| --- | --- |
| 构造 | `New`、`NewFn`、`Newf`、`NewFrom`、`NewfFrom`、`NewStatus`、`MustNewStatus` |
| 包装 | `Wrap`、`Wrapf`、`WrapStatus`、`WrapStatusFrom`；`WithStack` 的缺栈补齐、有栈/typed-nil 身份、原 cause 与栈复用 |
| 解析 | `ResolveStatus`、`AsItem`、`Unwrap` / `errors.Is` |
| 展示 | `Message` vs `Error`、`Format`（`%+v`） |
| 诊断 | `Inspect` 的主状态/根全文、实际原因树、已有栈、typed-nil、不可比较节点、循环/预算、失败原子性与重复/并发读取所有权 |
| 安全 | `ClampMessage`、`SanitizeForClient`、`RedactForLog` |
| 错误码 | `ValidateCode`、`Code.Message` |
| 回归 | `test` 子包对外 API 冒烟 |

## 发版

1. 更新 `CHANGELOG.md`「未发布」小节为具体版本号与日期。  
2. 打 tag：`vMAJOR.MINOR.PATCH`。  
3. 若存在破坏性变更，递增 **MAJOR** 并在 README / CHANGELOG 中醒目标注迁移说明。
