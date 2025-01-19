# CSS (Connect Service Scaffold)

一个基于 Connect RPC 的全栈 Web 开发框架，使用 Go 和 TypeScript 构建。

## 特性

- 🚀 基于 Connect RPC 的高性能通信
- 🔄 支持多种 RPC 模式：Unary、Server Stream、Client Stream、Bidirectional Stream
- 🌐 内置 HTTP/2 支持
- 🛠 完整的开发工具链
  - Protocol Buffers 代码生成
  - TypeScript 类型生成
  - Connect RPC 客户端和服务端代码生成
- 🎯 模块化的服务架构
- 🔌 可扩展的插件系统
- ⚡️ 支持 gRPC、gRPC-Web 和 Connect 协议
- 🔒 内置 CORS 支持
- 📝 结构化的日志系统

## 技术栈

### 后端

- Go
- Connect RPC
- Protocol Buffers
- Echo (HTTP 框架)
- Viper (配置管理)

### 前端

- TypeScript
- React
- Connect RPC Web Client

## 快速开始

### 前置条件

- Go 1.21+
- Node.js 20+
- Protocol Buffers 编译器
- Buf CLI 工具

### 安装

1. 克隆仓库：

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

3. 生成 Protocol Buffers 代码：

```bash
buf generate
```

### 开发

1. 启动后端服务：

```bash
# 或者你可以使用 go run backend/... serve
go run backend serve
```

2. 启动前端开发服务器：

```bash
cd frontend
pnpm dev
```

### 配置

项目使用 `config.yaml` 进行配置，支持以下配置项：

```yaml
env: development
server:
  port: 21421
  host: 0.0.0.0
log:
  level: info
  format: text
db:
  host: 127.0.0.1
  port: 5432
  user: postgres
  password: postgres
  name: app
  sslmode: disable
  timezone: Asia/Shanghai
  pool_max_conns: 10
  pool_max_idle_conns: 5
  pool_max_lifetime: 10m
  pool_max_idle_time: 5m
```

## 项目结构

```
.
├── backend/             # 后端代码
│   ├── cmd/            # 命令行工具
│   ├── core/           # 核心功能
│   ├── pkg/            # 生成的代码
│   └── service/        # 业务服务
├── frontend/           # 前端代码
├── proto/              # Protocol Buffers 定义
└── config.yaml         # 配置文件
```

## API 开发

1. 在 `proto/` 目录下定义服务接口
2. 使用 `buf generate` 生成代码
3. 在 `backend/service/` 实现服务接口
4. 在 `backend/app.go` 注册服务

示例服务定义：

```protobuf
service TestService {
  rpc Test1(Test1Request) returns (Test1Response) {}
  rpc Test2(Test2Request) returns (Test2Response) {}
  rpc Test3(Test3Request) returns (stream Test3Response) {}
}
```

## 构建和部署

### 构建

```bash
# 构建后端
go build -o app backend/main.go

# 构建前端
cd frontend
pnpm build
```

### 部署

1. 配置 `config.yaml`
2. 运行编译后的二进制文件：

```bash
./app serve
```

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交变更
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License
