# API 文档

本文档详细说明了 Go Admin Scaffold 的 API 接口。

## API 概述

### 基础信息

- 基础URL: `http://your-domain/api/admin/v1`
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
    "code": 0,            // 状态码，0 表示成功
    "message": "success", // 状态信息
    "data": {            // 响应数据
        // 具体数据
    },
    "trace_id": "..."    // 请求追踪ID
}
```

3. 错误响应
```json
{
    "code": 400,           // 错误码
    "message": "error",    // 错误信息
    "data": null,         // 错误时通常为 null
    "trace_id": "..."     // 请求追踪ID
}
```

4. 分页响应
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [],      // 数据列表
        "total": 100,     // 总记录数
        "page": 1,        // 当前页码
        "per_page": 20    // 每页记录数
    },
    "trace_id": "..."
}
```

## 认证接口

### 登录

- 路径: `/api/admin/v1/auth/login`
- 方法: POST
- 描述: 用户登录获取 token
- 请求体:
```json
{
    "username": "admin",
    "password": "password123"
}
```
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "token_type": "Bearer",
        "expires_in": 3600
    },
    "trace_id": "..."
}
```

### 刷新 Token

- 路径: `/api/admin/v1/auth/refresh`
- 方法: POST
- 描述: 刷新访问令牌
- 请求头: `Authorization: Bearer {token}`
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "access_token": "eyJhbGciOiJIUzI1NiIs...",
        "token_type": "Bearer",
        "expires_in": 3600
    },
    "trace_id": "..."
}
```

## 用户管理

### 获取用户列表

- 路径: `/api/admin/v1/users`
- 方法: GET
- 描述: 获取用户列表
- 权限: `user:view`
- 请求头: `Authorization: Bearer {token}`
- 参数:
  - `page`: 页码 (默认: 1)
  - `per_page`: 每页记录数 (默认: 20)
  - `search`: 搜索关键词
  - `status`: 状态筛选 (0: 禁用, 1: 启用)
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "username": "admin",
                "email": "admin@example.com",
                "nickname": "Administrator",
                "avatar": "http://...",
                "status": 1,
                "last_login_at": "2024-03-10T10:00:00Z",
                "created_at": "2024-03-10T10:00:00Z",
                "updated_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 100,
        "page": 1,
        "per_page": 20
    },
    "trace_id": "..."
}
```

### 创建用户

- 路径: `/api/admin/v1/users`
- 方法: POST
- 描述: 创建新用户
- 权限: `user:create`
- 请求头: `Authorization: Bearer {token}`
- 请求体:
```json
{
    "username": "newuser",
    "password": "password123",
    "email": "new@example.com",
    "nickname": "New User",
    "avatar": "http://...",
    "status": 1,
    "role_ids": [1, 2]
}
```
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 2,
        "username": "newuser",
        "email": "new@example.com",
        "nickname": "New User",
        "avatar": "http://...",
        "status": 1,
        "created_at": "2024-03-10T10:00:00Z",
        "updated_at": "2024-03-10T10:00:00Z"
    },
    "trace_id": "..."
}
```

### 获取用户详情

- 路径: `/api/admin/v1/users/{id}`
- 方法: GET
- 描述: 获取用户详细信息
- 权限: `user:view`
- 请求头: `Authorization: Bearer {token}`
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "nickname": "Administrator",
        "avatar": "http://...",
        "status": 1,
        "roles": [
            {
                "id": 1,
                "name": "admin",
                "description": "Administrator"
            }
        ],
        "last_login_at": "2024-03-10T10:00:00Z",
        "created_at": "2024-03-10T10:00:00Z",
        "updated_at": "2024-03-10T10:00:00Z"
    },
    "trace_id": "..."
}
```

### 更新用户

- 路径: `/api/admin/v1/users/{id}`
- 方法: PUT
- 描述: 更新用户信息
- 权限: `user:edit`
- 请求头: `Authorization: Bearer {token}`
- 请求体:
```json
{
    "email": "updated@example.com",
    "nickname": "Updated User",
    "avatar": "http://...",
    "status": 1,
    "role_ids": [1, 2]
}
```
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "admin",
        "email": "updated@example.com",
        "nickname": "Updated User",
        "avatar": "http://...",
        "status": 1,
        "updated_at": "2024-03-10T10:00:00Z"
    },
    "trace_id": "..."
}
```

### 删除用户

- 路径: `/api/admin/v1/users/{id}`
- 方法: DELETE
- 描述: 删除用户
- 权限: `user:delete`
- 请求头: `Authorization: Bearer {token}`
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": null,
    "trace_id": "..."
}
```

## 角色管理

### 获取角色列表

- 路径: `/api/admin/v1/roles`
- 方法: GET
- 描述: 获取角色列表
- 权限: `role:view`
- 请求头: `Authorization: Bearer {token}`
- 参数:
  - `page`: 页码 (默认: 1)
  - `per_page`: 每页记录数 (默认: 20)
  - `search`: 搜索关键词
  - `status`: 状态筛选 (0: 禁用, 1: 启用)
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "name": "admin",
                "description": "Administrator",
                "status": 1,
                "permissions": [
                    {
                        "id": 1,
                        "name": "user:view",
                        "description": "View users"
                    }
                ],
                "created_at": "2024-03-10T10:00:00Z",
                "updated_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 10,
        "page": 1,
        "per_page": 20
    },
    "trace_id": "..."
}
```

### 创建角色

- 路径: `/api/admin/v1/roles`
- 方法: POST
- 描述: 创建新角色
- 权限: `role:create`
- 请求头: `Authorization: Bearer {token}`
- 请求体:
```json
{
    "name": "editor",
    "description": "Content Editor",
    "status": 1,
    "permission_ids": [1, 2, 3]
}
```
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 2,
        "name": "editor",
        "description": "Content Editor",
        "status": 1,
        "created_at": "2024-03-10T10:00:00Z",
        "updated_at": "2024-03-10T10:00:00Z"
    },
    "trace_id": "..."
}
```

## 权限管理

### 获取权限列表

- 路径: `/api/admin/v1/permissions`
- 方法: GET
- 描述: 获取所有权限
- 权限: `permission:view`
- 请求头: `Authorization: Bearer {token}`
- 参数:
  - `module`: 按模块筛选
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "name": "user:view",
                "description": "View users",
                "module": "user",
                "status": 1,
                "created_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 50
    },
    "trace_id": "..."
}
```

### 获取权限模块

- 路径: `/api/admin/v1/permissions/modules`
- 方法: GET
- 描述: 获取所有权限模块
- 权限: `permission:view`
- 请求头: `Authorization: Bearer {token}`
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "name": "user",
                "permissions": [
                    {
                        "id": 1,
                        "name": "user:view",
                        "description": "View users"
                    }
                ]
            }
        ]
    },
    "trace_id": "..."
}
```

## 日志管理

### 获取登录日志

- 路径: `/api/admin/v1/logs/login`
- 方法: GET
- 描述: 获取用户登录日志
- 权限: `log:view`
- 请求头: `Authorization: Bearer {token}`
- 参数:
  - `page`: 页码 (默认: 1)
  - `per_page`: 每页记录数 (默认: 20)
  - `username`: 用户名筛选
  - `status`: 状态筛选 (0: 失败, 1: 成功)
  - `start_time`: 开始时间
  - `end_time`: 结束时间
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "user_id": 1,
                "username": "admin",
                "ip": "127.0.0.1",
                "user_agent": "Mozilla/5.0...",
                "status": 1,
                "message": "login successful",
                "created_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 1000,
        "page": 1,
        "per_page": 20
    },
    "trace_id": "..."
}
```

### 获取操作日志

- 路径: `/api/admin/v1/logs/operation`
- 方法: GET
- 描述: 获取用户操作日志
- 权限: `log:view`
- 请求头: `Authorization: Bearer {token}`
- 参数:
  - `page`: 页码 (默认: 1)
  - `per_page`: 每页记录数 (默认: 20)
  - `username`: 用户名筛选
  - `module`: 模块筛选
  - `action`: 操作类型筛选
  - `status`: 状态筛选 (0: 失败, 1: 成功)
  - `start_time`: 开始时间
  - `end_time`: 结束时间
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "user_id": 1,
                "username": "admin",
                "ip": "127.0.0.1",
                "method": "POST",
                "path": "/api/admin/v1/users",
                "action": "user management",
                "module": "users",
                "request_params": "{\"username\":\"newuser\"}",
                "status": 1,
                "duration": 100,
                "user_agent": "Mozilla/5.0...",
                "created_at": "2024-03-10T10:00:00Z"
            }
        ],
        "total": 1000,
        "page": 1,
        "per_page": 20
    },
    "trace_id": "..."
}
```

## 文件上传

### 单文件上传

- 路径: `/api/admin/v1/upload/file`
- 方法: POST
- 描述: 上传单个文件
- 权限: `upload:create`
- 请求头: 
  - `Authorization: Bearer {token}`
  - `Content-Type: multipart/form-data`
- 参数:
  - `type`: 文件类型 (avatar/image/other)
  - `file`: 文件数据
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "url": "http://...",
        "path": "uploads/images/xxx.jpg",
        "name": "xxx.jpg",
        "size": 1024,
        "type": "image/jpeg"
    },
    "trace_id": "..."
}
```

### 多文件上传

- 路径: `/api/admin/v1/upload/files`
- 方法: POST
- 描述: 上传多个文件
- 权限: `upload:create`
- 请求头: 
  - `Authorization: Bearer {token}`
  - `Content-Type: multipart/form-data`
- 参数:
  - `type`: 文件类型 (avatar/image/other)
  - `files`: 文件数据（可多个）
- 响应:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 3,
        "success": 2,
        "failed": 1,
        "files": [
            {
                "url": "http://...",
                "path": "uploads/images/xxx.jpg",
                "name": "xxx.jpg",
                "size": 1024,
                "type": "image/jpeg"
            }
        ]
    },
    "trace_id": "..."
}
```

## 错误码说明

### 通用错误码

- 0: 成功
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

## 文件上传限制

1. 头像上传
   - 最大大小: 2MB
   - 支持格式: jpg/jpeg/png/gif/webp

2. 图片上传
   - 最大大小: 5MB
   - 支持格式: jpg/jpeg/png/gif/webp

3. 其他文件
   - 最大大小: 10MB
   - 支持格式: 所有类型

## 相关文档

- [认证系统](../features/authentication.md)
- [RBAC 权限系统](../features/rbac.md)
- [操作日志](../features/operation-log.md)
- [部署指南](../deployment/README.md) 