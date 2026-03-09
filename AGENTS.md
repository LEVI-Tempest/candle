# Candle (Go) Agent Guide

Scope: this file applies to the entire `candle/` repository.

## Communication

- 默认使用中文与用户沟通；若用户明确要求其他语言，再切换。
- 解释技术决策时，优先给出可验证的依据（代码、测试、数据），避免空泛结论。

## Engineering principles

- 遵循最小修改原则：优先做小而准的改动，不做无关重构。
- 架构优先：新增能力要放在清晰边界内（识别在 `identify`、展示在 `charting`、数据接入在 `datasource`），避免把决策逻辑塞进展示层。
- 优先复用成熟库：能用稳定、维护中的第三方库就不要重复造轮子；只有在缺能力、性能/许可受限或可控性必须时才自研。
- 若选择不使用现成库，必须在变更说明中写明取舍（为什么不用、风险是什么、后续迁移路径是什么）。
- 先修根因，不做表面补丁；若只能临时绕过，要明确标记后续修复点。
- 不顺手修与任务无关的问题；必要时在交付说明中单独提示。
- 保持模块边界清晰：
  - `pkg/identify`: 形态识别与趋势/评分逻辑
  - `pkg/charting`: 图表渲染与标注展示
  - `pkg/datasource`: 外部数据源访问
  - `pkg/risk`: 风险参数与管理逻辑
- 不要为了“复用”引入只被调用一次的小函数，除非能明显提升可读性。

## Generic codex essentials (通用工程精髓)

- 变更应保持“范围最小 + 行为清晰”，避免无关重命名与大规模风格漂移。
- 与现有代码风格保持一致；新增抽象前先证明有复用价值。
- 先做最贴近改动点的验证，再扩展到全量验证，节省定位成本。
- 改行为必须补文档/示例/测试中的至少一项对应证据，避免“代码与说明脱节”。
- 优先用可读、可维护的直接实现，不为“炫技”引入复杂度。
- 搜索代码和文件优先用 `rg` / `rg --files`。
- 不使用破坏性 git 操作（如 `git reset --hard`、`git checkout --`）除非用户明确要求。

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

## Done criteria (交付前硬性检查)

- 代码变更后必须至少完成：`gofmt -w` + 受影响包测试。
- 若涉及公共接口、跨包调用或评分规则变更，必须补充或更新对应测试。
- 若修改 `.proto`，必须完成生成并确认仓库内无生成物漂移。
- 若改变用户可见行为，必须更新 `README.md`、`PATTERN_MARKERS_IMPLEMENTATION.md` 或 `docs/` 对应文档之一。
- 交付说明必须包含：改了什么、为什么改、怎么验证、已知未覆盖风险。

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

## Signal output contract (面向 Agent 的输出契约)

- 面向上层 Agent 的信号输出应优先使用结构化 JSON，避免仅输出自由文本。
- 推荐字段（至少保持语义兼容）：`symbol`、`as_of`、`patterns`、`trend`、`score`、`evidence`、`counter_evidence`、`invalid_if`。
- 所有分数字段统一在 `[0, 1]` 区间；置信度统一使用离散等级（如 `strong/medium/weak`）。
- 证据项应包含“命中因子 + 阈值 + 当前值 + 结论”，保证可审计与可复盘。
- 不允许在输出中省略时间戳与数据来源信息。

## Domain guardrails (trading context)

- 不要把形态信号表述为确定性结论；应表述为“概率/条件性信号”。
- 调整阈值、评分或风险参数时：
  - 明确写出变更前后行为差异
  - 补充或更新测试，覆盖关键边界条件
- 涉及策略建议时，区分“数据事实”与“推断结论”。
- 对外输出必须带研究免责声明：仅供研究与决策支持，不构成投资建议，不直接触发交易。

## Delivery template (建议交付模板)

- `Summary`: 变更范围与目标。
- `Design`: 关键取舍（为何用库/为何不自研、模块落位、兼容性影响）。
- `Validation`: 执行过的构建/测试命令与结果。
- `Risk`: 未覆盖场景、潜在误判条件、回退方案。

## File hygiene

- 不提交本地构建产物和临时文件（例如根目录二进制、临时 HTML）。
- 未经用户明确要求，不修改与当前任务无关的文件。

## Git workflow

- 提交代码时，`git commit` 的 message 必须使用英文。
