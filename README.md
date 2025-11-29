# Simple CLI

一个基于 Go 语言的轻量级 RESTful API 服务脚手架，采用清晰的分层架构和依赖注入设计。

## ✨ 特性

- 🚀 **现代化技术栈**: Gin + GORM + Cobra + Viper
- 🔧 **依赖注入**: 支持手动 DI 和 Wire 自动 DI 两种方式
- 📦 **模块化设计**: 清晰的分层架构，易于扩展
- 🗃️ **SQLite 存储**: 开箱即用，无需额外配置数据库
- ⚙️ **灵活配置**: 支持配置文件、命令行参数和环境变量

## 🏗️ 项目结构

```
simple-cli/
├── cmd/
│   └── main.go              # 程序入口
├── internal/
│   ├── app/
│   │   └── container.go     # 依赖容器（手动 DI）
│   ├── cmd/
│   │   ├── cmd.go           # CLI 根命令
│   │   ├── serve/           # serve 子命令
│   │   └── version/         # version 子命令
│   ├── config/
│   │   └── db.go            # 数据库配置
│   ├── handler/             # HTTP 处理层
│   │   ├── health/          # 健康检查
│   │   ├── order/           # 订单模块 (Wire DI)
│   │   ├── product/         # 产品模块
│   │   └── user/            # 用户模块
│   ├── interfaces/          # 接口定义
│   ├── model/               # 数据模型
│   ├── repository/          # 数据访问层
│   ├── server/              # HTTP 服务器
│   ├── service/             # 业务逻辑层
│   └── types/               # 通用类型
├── build/                   # 构建输出目录
├── config.yaml              # 配置文件
├── Makefile                 # 构建脚本
├── go.mod
└── go.sum
```

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Make (可选)
- Wire (用于生成依赖注入代码)

```bash
# 安装 Wire
go install github.com/google/wire/cmd/wire@latest
```

### 构建项目

```bash
# 使用 Make（推荐）
make all

# 或手动构建
go mod tidy
go build -o build/simple-cli cmd/main.go
```

### 运行服务

```bash
# 使用默认配置（端口 9001）
./build/simple-cli serve

# 指定端口
./build/simple-cli serve --port 8080

# 指定配置文件
./build/simple-cli serve --config ./config.yaml
```

### 查看版本

```bash
./build/simple-cli version
```

## 📚 API 接口

服务启动后，默认监听 `http://localhost:9001`

### 健康检查

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |

### 用户管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/users` | 创建用户 |
| GET | `/users/:id` | 获取用户详情 |
| PUT | `/users/:id` | 更新用户 |
| DELETE | `/users/:id` | 删除用户 |

### 产品管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/products` | 创建产品 |
| GET | `/products` | 获取产品列表 |
| GET | `/products/:id` | 获取产品详情 |
| PUT | `/products/:id` | 更新产品 |
| DELETE | `/products/:id` | 删除产品 |

### 订单管理

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/orders` | 创建订单 |
| GET | `/orders` | 获取订单列表 |
| GET | `/orders/:id` | 获取订单详情 |
| PUT | `/orders/:id/status` | 更新订单状态 |
| POST | `/orders/:id/cancel` | 取消订单 |

## 🔧 配置说明

### 配置文件 (config.yaml)

```yaml
port: 9001
```

### 环境变量

所有配置项都可以通过环境变量覆盖，前缀为 `SIMPLE_CLI_`：

```bash
export SIMPLE_CLI_PORT=8080
```

### 命令行参数

```bash
./build/simple-cli serve --port 8080 --config ./config.yaml
```

**优先级**: 命令行参数 > 环境变量 > 配置文件

## 🛠️ 开发指南

### Make 命令

```bash
make all       # 执行完整构建流程 (tidy + wire + imports + build)
make tidy      # 整理依赖
make wire      # 生成 Wire 依赖注入代码
make imports   # 格式化 imports
make build     # 编译项目
make run       # 运行项目
make test      # 运行测试
make clean     # 清理构建产物
make help      # 显示帮助信息
```

### 添加新模块

1. 在 `internal/model/` 创建数据模型
2. 在 `internal/interfaces/` 定义服务接口
3. 在 `internal/repository/` 实现数据访问层
4. 在 `internal/service/` 实现业务逻辑
5. 在 `internal/handler/` 实现 HTTP 处理器
6. 在 `internal/server/server.go` 注册模块

### 依赖注入方式

项目支持两种依赖注入方式：

**1. 手动 DI（通过 Container）**

适用于简单场景，参见 `user` 和 `product` 模块：

```go
// internal/app/container.go
c.UserService, err = userSrv.NewUserService(userSrv.WithUserRepository(c.UserRepo))
```

**2. Wire 自动 DI**

适用于复杂依赖关系，参见 `order` 模块：

```go
// internal/handler/order/wire.go
func InitOrderHandler(db *gorm.DB) (*OrderHandler, error) {
    wire.Build(ProviderSet)
    return nil, nil
}
```

## 📦 技术栈

| 组件 | 库 | 用途 |
|------|------|------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) | HTTP 路由和中间件 |
| ORM | [GORM](https://gorm.io/) | 数据库操作 |
| CLI | [Cobra](https://github.com/spf13/cobra) | 命令行接口 |
| 配置 | [Viper](https://github.com/spf13/viper) | 配置管理 |
| DI | [Wire](https://github.com/google/wire) | 编译时依赖注入 |
| 数据库 | SQLite | 嵌入式数据库 |

## 📄 License

MIT License

