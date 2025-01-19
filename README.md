# CSS (Connect Service Scaffold)

一个现代化的全栈 Web 开发框架，基于 Connect RPC 构建，使用 Go 和 TypeScript 技术栈。专注于提供高性能、类型安全和开发者友好的开发体验。

## 核心特性

- 🚀 **高性能通信**

  - 基于 Connect RPC 的高效通信协议
  - 支持 HTTP/2 和多种 RPC 模式
  - 兼容 gRPC、gRPC-Web 和 Connect 协议

- 🛠 **完整工具链**

  - Protocol Buffers 代码自动生成
  - TypeScript 类型定义生成
  - Connect RPC 客户端和服务端代码生成
  - 内置开发工具和调试支持

- 🎯 **企业级架构**

  - 模块化的服务设计
  - 可扩展的中间件系统
  - 结构化的日志系统
  - 灵活的配置管理

- 🔒 **安全性和可靠性**

  - 内置 CORS 支持
  - 请求追踪和日志记录
  - 错误处理和恢复机制

- 🎨 **现代前端集成**
  - 智能的开发环境代理
  - 生产环境静态文件服务
  - SPA 路由自动支持
  - 开发服务状态检测

## 技术栈

### 后端技术

- Go 1.21+
- Connect RPC (通信协议)
- Protocol Buffers (数据序列化)
- Echo (HTTP 框架)
- Zap (日志系统)
- Viper (配置管理)

### 前端技术

- TypeScript
- React
- Connect-Web (RPC 客户端)
- Vite (构建工具)

## 快速开始

### 环境要求

- Go 1.21 或更高版本
- Node.js 20 或更高版本
- Protocol Buffers 编译器
- Buf CLI 工具

### 安装步骤

1. 克隆项目：

```bash
git clone <repository-url>
cd css
```

2. 安装依赖：

```bash
# 后端依赖
go mod download

# 前端依赖
cd frontend
pnpm install
```

3. 生成代码：

```bash
buf generate
```

### 开发流程

1. 启动后端服务：

```bash
go run backend serve
```

2. 启动前端开发服务器：

```bash
cd frontend
pnpm dev
```

## 配置说明

配置文件：`config.yaml`

```yaml
env: development # 运行环境：development/production
server:
  port: 21421 # 后端服务端口
  host: 0.0.0.0 # 服务监听地址
frontend:
  port: 21422 # 前端开发服务器端口
  dist: ./dist # 前端构建输出目录
  isSpa: true # 是否为单页应用
log:
  level: info # 日志级别
  format: text # 日志格式
```

## 项目结构

```
.
├── backend/           # 后端代码
│   ├── cmd/          # 命令行工具
│   ├── pkg/          # 公共包
│   │   ├── logger/   # 日志系统
│   │   └── util/     # 工具函数
│   ├── serve/        # HTTP 服务
│   └── service/      # 业务服务
├── frontend/         # 前端代码
├── proto/            # Protocol Buffers 定义
└── config.yaml      # 配置文件
```

## 开发模式

### 开发环境

- 自动检测前端开发服务
- 智能请求代理
- 实时热重载
- 开发者友好的错误提示

### 生产环境

- 高效的静态文件服务
- SPA 路由支持
- 优化的资源加载
- 完整的错误处理

## API 开发指南

1. 在 `proto/` 目录定义服务接口
2. 使用 `buf generate` 生成代码
3. 在 `backend/service/` 实现服务
4. 在 `backend/app.go` 注册服务

示例：

```protobuf
service ExampleService {
  rpc Method(Request) returns (Response) {}
}
```

## 构建部署

### 构建

```bash
# 后端构建
go build -o app backend/main.go

# 前端构建
cd frontend
pnpm build
```

### 部署步骤

1. 配置 `config.yaml`
2. 确保前端已构建
3. 运行服务：

```bash
./app serve
```

## 许可证

MIT License
