# 代码静态检查 (Linting)

本项目使用 [golangci-lint](https://golangci-lint.run/) 进行代码静态检查，以确保代码质量和一致性。

## 配置文件

项目包含以下 lint 配置文件：

1. [.golangci.yml](file:///Users/wangyan/Code/Self/gomicro/.golangci.yml) - 主配置文件，用于 CI 和团队标准
2. [.golangci.local.yml](file:///Users/wangyan/Code/Self/gomicro/.golangci.local.yml) - 本地开发配置文件（不会提交到版本控制）

## 启用的检查器 (Linters)

配置启用了以下检查器：

- `errcheck` - 检查未处理的错误
- `gosimple` - 简化代码建议
- `govet` - Go 官方工具集检查器
- `ineffassign` - 检测无效赋值
- `staticcheck` - 综合静态分析检查器
- `typecheck` - 类型断言检查
- `unused` - 检测未使用的常量、变量、函数和类型
- `gocyclo` - 检测函数圈复杂度
- `gosec` - 安全问题检查
- `gofmt` - 检查代码格式
- `goimports` - 检查导入格式和缺失的导入
- `misspell` - 拼写错误检查
- `lll` - 超长行检查
- `unconvert` - 移除不必要的类型转换
- `revive` - 可配置的 Go lint 工具
- `unparam` - 检测未使用的函数参数
- `dogsled` - 检测过多空白标识符的赋值
- `nakedret` - 检测裸返回
- `prealloc` - 建议预分配切片
- `dupl` - 检测重复代码

## 运行检查

### 使用脚本运行（推荐）

```bash
./scripts/lint.sh
```

该脚本会自动检测并安装 golangci-lint（如尚未安装），然后运行检查。

### 直接运行

```bash
# 使用默认配置
golangci-lint run ./...

# 使用本地配置
golangci-lint run -c .golangci.local.yml ./...
```

## Git 提交前自动检查

项目配置了 Git pre-commit 钩子，会在每次提交前自动运行代码检查。只有通过检查的代码才能被提交。

### 安装钩子

```bash
./scripts/install-hooks.sh
```

### 跳过检查

如果需要跳过检查（不推荐），可以使用以下命令：

```bash
git commit --no-verify
```

### 钩子工作机制

1. 在提交前运行 golangci-lint 检查所有 Go 文件
2. 检查 go.mod 和 go.sum 文件的一致性
3. 如果检查失败，提交会被阻止

## IDE 集成

### VS Code

安装 Go 插件后，在设置中添加：

```json
"go.lintTool": "golangci-lint",
"go.lintFlags": [
  "--fast"
]
```

### GoLand

在 Settings/Preferences | Go | Linting 中选择 golangci-lint 并配置路径。

## 忽略特定问题

### 使用注释忽略

在需要忽略警告的代码行前添加注释：

```go
//nolint:errcheck
_, err := someFunction()
```

忽略多种检查：

```go
//nolint:errcheck,govet
someCode()
```

### 在配置中排除

可以在 [.golangci.yml](file:///Users/wangyan/Code/Self/gomicro/.golangci.yml) 的 `issues.exclude-rules` 部分添加排除规则。

## 最佳实践

1. 在提交代码前运行 lint 检查
2. 修复所有 lint 错误而非简单忽略
3. 团队应统一 lint 规则，保持代码风格一致
4. 对于不可避免的警告，使用具体的 `//nolint` 注释并说明原因