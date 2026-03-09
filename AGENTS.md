# Candle (Go) Agent Guide

Scope: this file applies to the entire `candle/` repository.

## Communication

- 默认使用中文与用户沟通；若用户明确要求其他语言，再切换。
- 解释技术决策时，优先给出可验证的依据（代码、测试、数据），避免空泛结论。

## Engineering principles

- 遵循最小修改原则：优先做小而准的改动，不做无关重构。
- 保持模块边界清晰：
  - `pkg/identify`: 形态识别与趋势/评分逻辑
  - `pkg/charting`: 图表渲染与标注展示
  - `pkg/datasource`: 外部数据源访问
  - `pkg/risk`: 风险参数与管理逻辑
- 不要为了“复用”引入只被调用一次的小函数，除非能明显提升可读性。

## Go code style

- 所有 Go 代码改动后运行 `gofmt -w`（至少覆盖变更文件）。
- 变量命名使用清晰英文；注释可中英混合，但要以可维护性为准，避免机械双语。
- 错误处理要显式，不要静默吞错。
- 避免过度抽象，优先直接、可读、可测试的实现。

## Build, test, and validation

- 本地构建：`go build ./...`
- 测试策略（按范围递进）：
  1. 先跑受影响包（例如 `go test ./pkg/identify`）
  2. 再跑全量测试 `go test ./...`
- 进行较大改动（算法、公共类型、跨包接口）时，额外运行：`go test -race ./...`
- 若改动涉及数据源/网络行为，说明哪些测试在离线环境无法覆盖。

## Proto and generated files

- 不要手改 `*.pb.go`（这些是生成文件）。
- 若修改了 `pkg/proto/*.proto` 或 `third_party/**/*.proto`：
  - 运行 `make config` 重新生成代码
  - 将生成结果与 `.proto` 变更一起提交，避免漂移
- 若本机缺少插件，先安装生成依赖（如 `protoc-gen-go`、`protoc-gen-go-grpc`、`protoc-gen-go-http`）再执行生成。

## Docs sync requirements

- 改了用户可见行为（CLI 输出、图表标记、示例流程），同步更新 `README.md` 或相关文档。
- 改了形态标注/显示逻辑（位置、强度阈值、筛选规则），同步更新 `PATTERN_MARKERS_IMPLEMENTATION.md`。
- 改了量价关系或决策逻辑，优先同步 `docs/` 下对应设计文档，保证“实现与文档”一致。

## Domain guardrails (trading context)

- 不要把形态信号表述为确定性结论；应表述为“概率/条件性信号”。
- 调整阈值、评分或风险参数时：
  - 明确写出变更前后行为差异
  - 补充或更新测试，覆盖关键边界条件
- 涉及策略建议时，区分“数据事实”与“推断结论”。

## File hygiene

- 不提交本地构建产物和临时文件（例如根目录二进制、临时 HTML）。
- 未经用户明确要求，不修改与当前任务无关的文件。

## Git workflow

- 提交代码时，`git commit` 的 message 必须使用英文。
