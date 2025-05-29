# API 文档

## 概述

本文档详细说明了 Go Admin 后台管理系统的 API 接口规范和使用方法。

## 基础信息

- 基础路径：`/api/v1`
- 内容类型：`application/json`
- 认证方式：JWT Bearer Token

## 认证

所有需要认证的 API 都需要在请求头中携带 Token：

```http
Authorization: Bearer <your-token>
```

获取 Token：

```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "admin",
    "password": "password"
}
```

## 响应格式

### 成功响应

```json
{
    "code": 0,
    "data": {},
    "message": "success"
}
```

### 错误响应

```json
{
    "code": 40001,
    "message": "未授权的访问",
    "data": null
}
```

## API 目录

### 认证相关

- [认证概述](overview.md#authentication)
- [登录认证](overview.md#login)
- [刷新令牌](overview.md#refresh-token)

### 用户管理

- [用户管理概述](overview.md#user-management)
- [用户接口](overview.md#user-api)

### 角色权限

- [权限管理](overview.md#permissions)
- [角色管理](overview.md#roles)

## 错误码

| 错误码 | 说明 |
|--------|------|
| 40001  | 未授权 |
| 40002  | 参数错误 |
| 40003  | 资源不存在 |
| 40004  | 权限不足 |
| 50001  | 服务器错误 |

## 使用示例

### cURL

```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}'

# 获取用户列表
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer <your-token>"
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

func main() {
    // 登录请求
    loginData := map[string]string{
        "username": "admin",
        "password": "password",
    }
    jsonData, _ := json.Marshal(loginData)
    
    resp, _ := http.Post(
        "http://localhost:8080/api/v1/auth/login",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    // ... 处理响应
}
```

## 相关文档

- [认证系统](../features/authentication.md)
- [项目结构](../getting-started/structure.md)
- [配置说明](../getting-started/configuration.md) 