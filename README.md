# mq Project

基于 DDD (Domain-Driven Design) 架构的 Go 项目。

## 技术栈

### 核心框架
- Go 1.22+
- go-zero v1.6.0 (微服务框架)
- JWT 认证
- Redis 缓存

### 数据库
- MySQL 
- Redis 

### CI/CD
- GitHub Actions

## 项目结构

```
.
├── cmd/                    # 应用程序入口
│   └── main.go        # 主程序入口
├── domain/                # 领域层
│   ├── model/            # 领域模型
│   ├── repository/       # 仓储接口定义
│   └── service/          # 领域服务
├── application/          # 应用层
│   ├── dto/             # 数据传输对象
│   ├── service/         # 应用服务
│   └── usecase/         # 用例实现
├── infrastructure/       # 基础设施层
│   ├── persistence/     # 持久化实现
│   ├── config/         # 配置管理
│   ├── intergration/   # 集成服务
│   ├── provider/       # 提供者
│   ├── sql/           # 数据表DDL
│   ├── startup/        # 服务启动初始化
│   └── svc/            # 服务上下文
├── interfaces/          # 接口层
│   ├── api/            # API 接口
│   │   └── handler/    # HTTP 处理器
│   └── middleware/     # 中间件
├── common/             # 公共组件
│   ├── errors/        # 错误定义
│   ├── utils/         # 工具函数
│   └── constants/     # 常量定义
├── etc/               # 配置文件
├── pkg/               # 可重用的包
├── scripts/           # 脚本文件
├── test/              # 测试文件
├── vendor/            # 依赖包
├── go.mod             # Go 模块文件
├── go.sum             # Go 依赖版本锁定
├── Makefile           # 构建脚本
├── Dockerfile         # Docker 构建文件
└── docker-compose.yml # Docker 编排文件
```

## 目录说明

### cmd/
- 包含应用程序的入口点
- 负责配置和启动应用程序

### domain/
- 包含核心业务逻辑和规则
- 定义领域模型、实体和值对象
- 包含领域服务接口

### application/
- 实现用例和业务流程
- 协调领域对象和基础设施
- 处理事务和事件

### infrastructure/
- 提供技术实现细节
- 实现持久化和外部服务集成
- 处理配置和日志

### interfaces/
- 处理外部请求和响应
- 实现 API 和 gRPC 接口
- 包含中间件和路由配置

### common/
- 提供共享的工具和组件
- 定义通用错误类型
- 包含常量定义

## 开发指南

1. 遵循 DDD 原则进行开发
2. 使用 Go 1.22 或更高版本
3. 遵循 RESTful API 设计规范
4. 实现适当的错误处理和日志记录
5. 编写单元测试和集成测试

## 构建和运行

```bash
# 安装依赖
go mod download

# 运行测试
make test

# 构建项目
make build

# 运行服务
make run
```