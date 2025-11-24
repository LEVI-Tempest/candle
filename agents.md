# Candle Project - AI Agent Configuration / 蜡烛图项目 - AI助手配置

## Project Overview / 项目概述

This is a Japanese Candlestick Charting project focused on pattern recognition and visualization.
这是一个专注于形态识别和可视化的日本蜡烛图项目。

## Build & Test / 构建与测试

- **Build**: `go build -o candle main.go`
- **Test**: `go test ./...`
- **Run examples**: `go run examples/precise_pattern_markers.go`
- **Generate charts**: Open generated HTML files in browser

## Core Principles / 核心原则

### 1. Language Policy / 语言政策
- **Always response in Chinese** / **始终使用中文回复**
  - All interactions should be conducted in Chinese
  - 所有交互都应使用中文进行

### 2. Code Documentation Standards / 代码文档标准
- **All code comments should be in Chinese and English** / **所有代码注释应使用中英文双语**
  - Provide bilingual comments for better understanding
  - 提供双语注释以便更好理解
  - Format: `// 中文说明 | English explanation`

### 3. Development Philosophy / 开发理念
- **Minimal Modification Principle** / **最小修改原则**
  - Make the smallest possible changes to achieve the desired outcome
  - 以最小的修改实现预期目标
  - Preserve existing architecture and functionality
  - 保持现有架构和功能

### 4. Architectural Thinking / 架构思维
- **Think like a computer architect, consider the architecture like a financial and trading expert** / **像计算机架构师一样思考，像金融交易专家一样考虑架构**
  - Apply systematic architectural principles
  - 应用系统化的架构原则
  - Consider performance, scalability, and reliability like financial systems
  - 像金融系统一样考虑性能、可扩展性和可靠性
  - Balance technical excellence with business requirements
  - 平衡技术卓越性与业务需求

## Project Structure / 项目结构

```
├── pkg/
│   ├── charting/     # 图表生成和可视化 | Chart generation and visualization
│   ├── identify/     # 蜡烛图形态识别 | Candlestick pattern recognition
│   ├── proto/        # Protocol buffer定义 | Protocol buffer definitions
│   ├── risk/         # 风险管理模块 | Risk management module
│   └── utils/        # 工具函数 | Utility functions
├── examples/         # 示例代码 | Example code
├── third_party/      # 第三方依赖 | Third-party dependencies
└── main.go          # 主程序入口 | Main program entry
```

## Architecture Overview / 架构概述

The system follows a modular design with clear separation of concerns:
系统采用模块化设计，职责分离清晰：

- **Pattern Recognition**: Identifies Japanese candlestick patterns
- **Charting Engine**: Generates interactive HTML charts with ECharts
- **Visualization**: Precise position marking of detected patterns
- **Risk Management**: Trading risk assessment and management

## Conventions & Patterns / 约定与模式

### File Naming / 文件命名
- Go files: `snake_case.go`
- Test files: `*_test.go`
- Example files: `examples/descriptive_name.go`

### Code Style / 代码风格
- Follow Go standard formatting (`gofmt`)
- Use meaningful variable names in English
- Add bilingual comments for complex logic
- Keep functions focused and small

### Pattern Detection / 形态检测
- All patterns implement the `Pattern` interface
- Strength scoring: 0.0 to 1.0 scale
- Only display patterns with strength ≥ 0.8
- Position marking must be precise to the exact candlestick

## Implementation Guidelines / 实施指南

### Code Review Checklist / 代码审查清单
- [ ] 中英文双语注释完整 | Bilingual comments complete
- [ ] 遵循最小修改原则 | Follows minimal modification principle
- [ ] 架构设计合理 | Architecture design is sound
- [ ] 性能考虑充分 | Performance considerations adequate
- [ ] 可扩展性良好 | Good scalability
- [ ] 错误处理完善 | Comprehensive error handling
- [ ] 形态检测准确性 | Pattern detection accuracy
- [ ] 图表可视化清晰 | Chart visualization clarity

### Communication Standards / 沟通标准
- Use Chinese for all discussions and documentation
- 所有讨论和文档使用中文
- Provide context and reasoning for architectural decisions
- 为架构决策提供上下文和推理
- Consider both technical and business implications
- 同时考虑技术和业务影响

### Quality Assurance / 质量保证
- Prioritize system stability and reliability
- 优先考虑系统稳定性和可靠性
- Implement robust testing strategies
- 实施强健的测试策略
- Validate pattern detection accuracy
- 验证形态检测准确性
- Test chart rendering across browsers
- 测试图表在不同浏览器中的渲染

## External Dependencies / 外部依赖

- **ECharts**: For interactive chart generation
- **go-echarts**: Go wrapper for ECharts
- **Protocol Buffers**: For data serialization
- **Third-party APIs**: Real-time market data (东方财富网)

## Security Considerations / 安全考虑

- Validate all input data before processing
- 处理前验证所有输入数据
- Sanitize HTML output to prevent XSS
- 清理HTML输出以防止XSS攻击
- Implement rate limiting for API calls
- 为API调用实施速率限制

---

*This configuration ensures consistent, high-quality development practices aligned with financial industry standards.*
*此配置确保符合金融行业标准的一致、高质量开发实践。*
