# order-service

Go 电商交易链路服务 —— 求职项目（进行中）。

技术栈：Go · Gin · GORM/MySQL · Redis · Kafka(规划) · gRPC(规划) · Docker Compose

## 目录结构

```
order-service/
├── cmd/server/              # 程序入口
├── internal/                # 私有代码（Go 强制，外部模块无法 import）
│   ├── handler/             # HTTP 层：参数校验 + 响应，不含业务逻辑
│   ├── service/             # 业务逻辑
│   ├── repository/          # 数据访问：返回 (data, error)
│   ├── model/               # 数据结构
│   ├── config/              # 配置加载
│   └── pkg/                 # 包内共享工具
│       ├── response/        # 统一响应体
│       └── apperr/          # 业务错误码与 AppError
├── migrations/              # 建表 SQL
├── docs/                    # 设计文档、架构图、配置说明
├── docker-compose.yml       # 本地依赖：MySQL 8.0 + Redis 7
└── .env.example             # 环境变量模板（.env 含凭据，不提交）
```

## 设计约定

- **handler 薄层**：只做参数校验、调用 service、写响应
- **错误不返回 message**：repository / service 返回 error，文案由 handler 统一翻译
- **业务错误码与 HTTP status 解耦**：`code` 给业务判断，HTTP status 给网络层
- **接口由使用方定义**（9/23 决策）：`OrderRepository` 接口**定义在 `internal/service`**，
  `internal/repository` 只放**实现**且不 import service。
  好处：① service 只声明自己需要的方法（接口不会肥大）② service **不依赖 repository 包**（依赖倒置真正成立）
  ③ 单测时传一个假实现即可，不用连数据库。
  装配点在 `main.go` —— 两边互不认识，由 main 把它们接起来。
  > 备选做法是接口放实现方（repository 包），在 DDD/Java 背景的团队里很常见，
  > 代价是接口容易变胖、且 service 必须 import repository 包。本项目选「使用方定义」。

## 快速开始

```bash
# 1) 起依赖：MySQL + Redis
cp .env.example .env          # 填入你的本地密码
docker compose up -d
docker compose ps             # mysql 应显示 healthy

# 2) 跑服务
make tidy                     # 拉依赖
make run                      # 启动服务（默认 :8080）
make test                     # 跑测试
make race                     # 竞态检测（-race）
```

> 依赖配置的取舍说明见 **`docs/docker-compose-说明.md`**
> （为什么锁镜像版本 / 端口映射怎么读 / 为什么走 .env / 为什么 utf8mb4 /
> 为什么挂命名卷 / healthcheck 与 `depends_on` 的关系）

## 接口

| Method | Path | 说明 |
| :--- | :--- | :--- |
| GET | `/health` | 健康检查 |
| GET | `/api/v1/orders/:id` | 查询订单详情 |

### 统一响应格式

```json
{ "code": 0, "msg": "ok", "data": { ... } }
```

`code = 0` 表示成功；非 0 为业务错误码（如 `40401` 订单不存在、`40001` 参数不合法）。

## 开发进度

- [x] W1：项目骨架、分层、统一响应与错误码
- [ ] **W2 开工前（30 min）**：把 `cmd/server/main.go` 拆成 `main.go` / `wire.go` / `router.go`
      —— 装配与路由分开，避免 W2/W3 加依赖后 main 膨胀成流水账（**不引入 wire 代码生成**，手工 DI 足够）
- [ ] W2：下单接口、MySQL 事务与库存扣减
- [ ] W3：Redis+Lua 秒杀、Kafka 削峰
- [ ] W4：库存服务拆分（gRPC + etcd）、压测
