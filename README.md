# ToDoList 项目需求文档 (第一版)

## 1. 项目概述
ToDoList是一个简单的任务管理应用，用户可以创建、修改、删除和查看任务。项目采用前后端分离架构，后端使用Go语言的Gin框架，数据库使用MySQL，缓存使用Redis，前端使用Vue.js。

## 2. 系统需求

### 2.1 功能需求

#### 用户管理
- 用户注册：用户名、密码注册
- 用户登录：用户名和密码登录
- 用户登出：退出登录

#### 任务管理
- 创建任务：标题、描述、截止时间
- 更新任务：修改任务内容
- 删除任务：删除指定任务
- 查看任务：任务列表展示
- 任务状态：标记完成/未完成

### 2.2 技术需求
- 前端：Vue 3 + Element Plus + Pinia + Vite
- 后端：Go + Gin + Gorm
- 数据库：MySQL
- 缓存：Redis (用于JWT黑名单)
- 认证：JWT

## 3. 技术架构

### 3.1 前端架构
- 使用Vue 3作为前端框架
- 使用Element Plus作为UI组件库
- 使用Pinia进行状态管理
- 使用Vue Router进行路由管理
- 使用Vite作为构建工具

### 3.2 后端架构
- 使用Gin框架处理HTTP请求
- 使用Gorm进行数据库操作
- 使用JWT进行用户认证
- 使用Redis存储JWT黑名单

### 3.3 数据库设计

```sql
-- 用户表（users）
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 任务表（tasks）
CREATE TABLE tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    status TINYINT DEFAULT 0, -- 0: 未完成, 1: 已完成
    due_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## 4. API设计

### 4.1 用户模块
```
POST /api/v1/user/register    - 用户注册
请求体：
{
    "username": "string",
    "password": "string"
}

POST /api/v1/user/login       - 用户登录
请求体：
{
    "username": "string",
    "password": "string"
}

POST /api/v1/user/logout      - 用户登出
Header: Authorization: Bearer <token>
```

### 4.2 任务模块
```
POST /api/v1/tasks           - 创建任务
请求体：
{
    "title": "string",
    "description": "string",
    "due_date": "2024-01-01T15:04:05Z"
}

GET /api/v1/tasks           - 获取任务列表
查询参数：
page: int
page_size: int

GET /api/v1/tasks/:id      - 获取单个任务

PUT /api/v1/tasks/:id      - 更新任务
请求体：
{
    "title": "string",
    "description": "string",
    "due_date": "2024-01-01T15:04:05Z"
}

DELETE /api/v1/tasks/:id   - 删除任务

PATCH /api/v1/tasks/:id/status - 更新任务状态
请求体：
{
    "status": 0  // 0: 未完成, 1: 已完成
}
```

## 5. 项目结构

```
backend/
├── config/             # 配置文件
│   └── config.yaml
├── internal/
│   ├── api/           # API处理器
│   │   ├── user.go
│   │   └── task.go
│   ├── middleware/    # 中间件
│   │   └── auth.go
│   ├── model/        # 数据模型
│   │   ├── user.go
│   │   └── task.go
│   ├── repository/   # 数据访问层
│   │   ├── user.go
│   │   └── task.go
│   └── service/      # 业务逻辑层
│       ├── user.go
│       └── task.go
├── pkg/              # 公共包
│   ├── jwt/
│   └── utils/
└── main.go          # 入口文件

frontend/
├── src/
│   ├── api/         # API请求
│   ├── components/  # 组件
│   ├── views/       # 页面
│   │   ├── Login.vue
│   │   └── Tasks.vue
│   ├── router/      # 路由
│   └── store/       # 状态管理
└── package.json
```

## 6. 用户界面要求

### 6.1 登录/注册界面
- 简洁的登录表单
- 用户名、密码输入框
- 登录、注册按钮
- 表单验证

### 6.2 任务列表界面
- 任务列表展示
- 创建任务按钮
- 任务项操作按钮（编辑、删除、完成）
- 分页控件

### 6.3 任务编辑界面
- 任务标题输入
- 任务描述输入
- 截止日期选择
- 保存/取消按钮

## 7. 后期迭代计划

### 第二版
1. 添加任务分类功能
2. 添加任务优先级
3. 实现基础的任务筛选

### 第三版
1. 实现任务提醒功能
2. 添加任务标签系统
3. 优化性能，添加缓存

### 第四版
1. 添加管理员功能
2. 实现数据统计
3. 添加任务导入导出

## 8. 开发流程
1. 搭建基础项目结构
2. 实现用户认证系统
3. 实现基础的任务CRUD
4. 开发前端界面
5. 进行集成测试
6. 部署上线

## 9. 注意事项
1. 所有API请求需要进行参数验证
2. 密码需要加密存储
3. 所有API需要进行错误处理
4. 需要添加适当的日志记录
5. 代码需要编写单元测试