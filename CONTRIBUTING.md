# 参与贡献与发布检查清单

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
| 构造 | `New`、`NewFn`、`Newf`、`NewFrom`、`NewfFrom` |
| 包装 | `Wrap`、`Wrapf`、`WrapStatus`、`WrapStatusFrom` |
| 解析 | `ResolveStatus`、`AsItem`、`Unwrap` / `errors.Is` |
| 展示 | `Message` vs `Error`、`Format`（`%+v`） |
| 安全 | `ClampMessage`、`SanitizeForClient`、`RedactForLog` |
| 错误码 | `ValidateCode`、`Code.Message` |
| 回归 | `test` 子包对外 API 冒烟 |

## 发版

1. 更新 `CHANGELOG.md`「未发布」小节为具体版本号与日期。  
2. 打 tag：`vMAJOR.MINOR.PATCH`。  
3. 若存在破坏性变更，递增 **MAJOR** 并在 README / CHANGELOG 中醒目标注迁移说明。
