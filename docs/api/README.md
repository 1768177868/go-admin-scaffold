# API 文档

本文档详细说明了 Go Admin Scaffold 的 API 接口。

## API 概述

### 基础信息

- 基础URL: `http://your-domain/api`
- 版本: v1
- 格式: JSON
- 认证: JWT Bearer Token

### 通用规范

1. 请求格式
   - GET 请求参数使用 Query String
   - POST/PUT 请求参数使用 JSON Body
   - 文件上传使用 multipart/form-data

2. 响应格式
```json
{
    "code": 200,           // 状态码
    "message": "success",  // 状态信息
    "data": {             // 响应数据
        // 具体数据
    }
}
```

3. 错误响应
```json
{
    "code": 400,           // 错误码
    "message": "error",    // 错误信息
    "errors": {           // 详细错误信息
        "field": ["error message"]
    }
}
```

4. 分页响应
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "items": [],      // 数据列表
        "total": 100,     // 总记录数
        "page": 1,        // 当前页码
        "per_page": 20    // 每页记录数
    }
}
```

## 认证接口

### 登录

- 路径: `/api/auth/login`
- 方法: POST
- 描述: 用户登录获取 token
- 请求体:
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "user": {
            "id": 1,
            "username": "admin",
            "email": "admin@example.com",
            "role": "admin"
        }
    }
}
```

### 登出

- 路径: `/api/auth/logout`
- 方法: POST
- 描述: 用户登出
- 请求头: `Authorization: Bearer {token}`
- 响应:
```json
{
    "code": 200,
    "message": "success"
}
```

### 刷新 Token

- 路径: `/api/auth/refresh`
- 方法: POST
- 描述: 刷新访问令牌
- 请求头: `Authorization: Bearer {token}`
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs..."
    }
}
```

## 用户管理

### 获取用户列表

- 路径: `/api/users`
- 方法: GET
- 描述: 获取用户列表
- 权限: `users.list`
- 参数:
  - `page`: 页码 (默认: 1)
  - `per_page`: 每页记录数 (默认: 20)
  - `search`: 搜索关键词
  - `role`: 角色筛选
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "username": "admin",
                "email": "admin@example.com",
                "role": "admin",
                "created_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "per_page": 20
    }
}
```

### 创建用户

- 路径: `/api/users`
- 方法: POST
- 描述: 创建新用户
- 权限: `users.create`
- 请求体:
```json
{
    "username": "newuser",
    "email": "new@example.com",
    "password": "password123",
    "role": "user"
}
```
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 2,
        "username": "newuser",
        "email": "new@example.com",
        "role": "user",
        "created_at": "2024-03-10T10:00:00Z"
    }
}
```

### 更新用户

- 路径: `/api/users/{id}`
- 方法: PUT
- 描述: 更新用户信息
- 权限: `users.update`
- 请求体:
```json
{
    "username": "updated",
    "email": "updated@example.com",
    "role": "admin"
}
```
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 1,
        "username": "updated",
        "email": "updated@example.com",
        "role": "admin",
        "updated_at": "2024-03-10T10:00:00Z"
    }
}
```

### 删除用户

- 路径: `/api/users/{id}`
- 方法: DELETE
- 描述: 删除用户
- 权限: `users.delete`
- 响应:
```json
{
    "code": 200,
    "message": "success"
}
```

## 角色管理

### 获取角色列表

- 路径: `/api/roles`
- 方法: GET
- 描述: 获取角色列表
- 权限: `roles.list`
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "name": "admin",
                "description": "Administrator",
                "permissions": ["users.list", "users.create"]
            }
        ],
        "total": 10,
        "page": 1,
        "per_page": 20
    }
}
```

### 创建角色

- 路径: `/api/roles`
- 方法: POST
- 描述: 创建新角色
- 权限: `roles.create`
- 请求体:
```json
{
    "name": "editor",
    "description": "Content Editor",
    "permissions": ["posts.list", "posts.create"]
}
```
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": 2,
        "name": "editor",
        "description": "Content Editor",
        "permissions": ["posts.list", "posts.create"],
        "created_at": "2024-03-10T10:00:00Z"
    }
}
```

## 权限管理

### 获取权限列表

- 路径: `/api/permissions`
- 方法: GET
- 描述: 获取所有权限
- 权限: `permissions.list`
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "name": "users.list",
                "description": "List users"
            }
        ]
    }
}
```

## 系统管理

### 获取系统信息

- 路径: `/api/system/info`
- 方法: GET
- 描述: 获取系统信息
- 权限: `system.info`
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "version": "1.0.0",
        "environment": "production",
        "database": {
            "driver": "mysql",
            "version": "5.7"
        },
        "redis": {
            "version": "6.0"
        },
        "queue": {
            "driver": "redis",
            "status": "running"
        }
    }
}
```

### 获取系统日志

- 路径: `/api/system/logs`
- 方法: GET
- 描述: 获取系统日志
- 权限: `system.logs`
- 参数:
  - `page`: 页码
  - `per_page`: 每页记录数
  - `level`: 日志级别
  - `start_date`: 开始日期
  - `end_date`: 结束日期
- 响应:
```json
{
    "code": 200,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "level": "info",
                "message": "User logged in",
                "context": {
                    "user_id": 1,
                    "ip": "127.0.0.1"
                },
                "created_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 1000,
        "page": 1,
        "per_page": 20
    }
}
```

## 错误码说明

### 通用错误码

- 200: 成功
- 400: 请求参数错误
- 401: 未认证
- 403: 无权限
- 404: 资源不存在
- 422: 数据验证错误
- 500: 服务器错误

### 业务错误码

- 1001: 用户不存在
- 1002: 密码错误
- 1003: 账号已禁用
- 1004: Token 已过期
- 1005: Token 无效
- 2001: 角色不存在
- 2002: 角色名称已存在
- 3001: 权限不存在
- 3002: 权限已分配

## 相关文档

- [认证系统](../features/authentication.md)
- [RBAC 权限系统](../features/rbac.md)
- [操作日志](../features/operation-log.md)
- [部署指南](../deployment/README.md) 