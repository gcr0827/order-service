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
└── docs/                    # 设计文档、架构图
```

## 设计约定

- **handler 薄层**：只做参数校验、调用 service、写响应
- **错误不返回 message**：repository / service 返回 error，文案由 handler 统一翻译
- **业务错误码与 HTTP status 解耦**：`code` 给业务判断，HTTP status 给网络层
- **面向接口**：service 依赖 `OrderRepository` 接口，便于单测 mock

## 快速开始

```bash
make tidy    # 拉依赖
make run     # 启动服务（默认 :8080）
make test    # 跑测试
make race    # 竞态检测
```

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
- [ ] W2：下单接口、MySQL 事务与库存扣减
- [ ] W3：Redis+Lua 秒杀、Kafka 削峰
- [ ] W4：库存服务拆分（gRPC + etcd）、压测
