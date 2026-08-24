# 事件总线系统（eventbus）

纯 Go 标准库实现的事件总线后端服务，零第三方依赖。用于管理事件主题、发布者、订阅者，并模拟事件投递、重试与死信处理。

## 运行

```bash
# 启动服务（默认监听 :8080）
go run ./cmd/server

# 自定义端口与配置
PORT=9090 MAX_PAGE_SIZE=100 MAX_ATTEMPTS=3 go run ./cmd/server
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8080 | 监听端口 |
| ADDR | 空（覆盖 PORT） | 完整监听地址 |
| MAX_PAGE_SIZE | 100 | 分页单页最大条数 |
| MAX_ATTEMPTS | 3 | 事件投递最大重试次数，超过则进入死信 |
| LOG_LEVEL | info | 日志级别（debug/info/warn/error） |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误时 `code` 非 0，HTTP 状态码对应 400/404/409/500。

### 主题 Topic

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/topics | 新增主题（名称唯一） |
| GET | /api/topics?keyword=&page=&size= | 分页列表 |
| GET | /api/topics/{id} | 查询主题 |
| PUT | /api/topics/{id} | 更新主题 |
| DELETE | /api/topics/{id} | 删除主题 |

### 发布者 Publisher

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/publishers | 新增发布者（校验 topic_id 存在） |
| GET | /api/publishers?topic_id=&status=&keyword=&page=&size= | 分页列表 |
| GET | /api/publishers/{id} | 查询发布者 |
| PUT | /api/publishers/{id} | 更新发布者 |
| DELETE | /api/publishers/{id} | 删除发布者 |
| POST | /api/publishers/{id}/activate | 启用发布者 |
| POST | /api/publishers/{id}/disable | 禁用发布者 |

### 订阅者 Subscriber

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/subscribers | 新增订阅者（校验 topic_id 存在） |
| GET | /api/subscribers?topic_id=&status=&keyword=&page=&size= | 分页列表 |
| GET | /api/subscribers/{id} | 查询订阅者 |
| PUT | /api/subscribers/{id} | 更新订阅者 |
| DELETE | /api/subscribers/{id} | 删除订阅者 |
| POST | /api/subscribers/{id}/activate | 启用订阅者 |
| POST | /api/subscribers/{id}/disable | 禁用订阅者 |
| POST | /api/subscribers/batch-disable | 批量禁用，body `{"ids":[...]}` |

状态枚举：`active`（启用）/ `inactive`（禁用）

### 事件 Event

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/events | 发布事件（校验 topic_id、publisher_id 存在） |
| GET | /api/events?topic_id=&publisher_id=&status=&keyword=&page=&size= | 分页列表 |
| GET | /api/events/{id} | 查询事件 |
| DELETE | /api/events/{id} | 删除事件 |
| POST | /api/events/{id}/deliver | 模拟投递成功（pending/failed→delivered） |
| POST | /api/events/{id}/fail | 模拟投递失败（pending→failed），body `{"reason":"..."}` |
| POST | /api/events/{id}/retry | 重试（attempts+1），body `{"success":bool,"reason":"..."}` |

事件状态枚举与流转：
- `pending`（待投递）→ `delivered`（成功）/ `failed`（失败）
- `failed`（失败）→ `delivered`（重试成功）/ `dead`（重试次数达到 MAX_ATTEMPTS 进入死信）

### 死信 DeadLetter

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/dead-letters?topic_id=&page=&size= | 分页列表 |
| GET | /api/dead-letters/{id} | 查询死信 |
| DELETE | /api/dead-letters/{id} | 删除死信 |

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/stats/overview | 整体概览统计 |
| GET | /healthz | 健康检查 |

## 项目结构

```
origin/
├── cmd/server/main.go       # 入口：配置加载、依赖装配、优雅关闭
├── internal/
│   ├── app/app.go           # 依赖装配 store -> service -> handler
│   ├── config/config.go     # 环境变量配置
│   ├── model/               # 领域模型 + 校验 + 状态机
│   ├── store/               # Store 接口 + 内存实现
│   ├── service/             # 业务逻辑（投递/重试/死信）
│   └── handler/             # HTTP 路由 + 处理器
└── pkg/
    ├── httpx/               # 统一响应、分页、JSON 解析
    ├── idgen/               # ID 生成
    └── logger/              # 分级日志
```

## 测试

```bash
go test ./...
```
