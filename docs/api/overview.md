# API 文档

本框架提供了 RESTful API 接口，支持 JSON 格式的请求和响应。

## API 版本控制

API 使用 URL 前缀进行版本控制：

```
https://api.example.com/v1/users
https://api.example.com/v2/users
```

## 认证

大多数 API 端点需要认证。使用 Bearer Token 进行认证：

```http
Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...
```

## 请求格式

### Content-Type

所有请求应该设置：

```http
Content-Type: application/json
```

### 分页参数

列表接口支持分页：

```
GET /api/v1/users?page=1&per_page=20
```

### 排序参数

支持字段排序：

```
GET /api/v1/users?sort=created_at:desc
```

### 过滤参数

支持字段过滤：

```
GET /api/v1/users?status=active&role=admin
```

## 响应格式

### 成功响应

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
    }
}
```

### 列表响应

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [
            {
                "id": 1,
                "name": "John Doe"
            },
            {
                "id": 2,
                "name": "Jane Doe"
            }
        ],
        "pagination": {
            "current_page": 1,
            "per_page": 20,
            "total": 50,
            "total_pages": 3
        }
    }
}
```

### 错误响应

```json
{
    "code": 400,
    "message": "Validation failed",
    "errors": {
        "email": ["The email field is required"],
        "password": ["The password must be at least 8 characters"]
    }
}
```

## 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 422 | 验证失败 |
| 429 | 请求过于频繁 |
| 500 | 服务器错误 |

## API 端点

### 认证相关

#### 登录
```http
POST /api/v1/auth/login

Request:
{
    "email": "user@example.com",
    "password": "password"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
        "expires_in": 3600
    }
}
```

#### 注册
```http
POST /api/v1/auth/register

Request:
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password",
    "password_confirmation": "password"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com"
    }
}
```

### 用户相关

#### 获取用户列表
```http
GET /api/v1/users

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [...],
        "pagination": {...}
    }
}
```

#### 创建用户
```http
POST /api/v1/users

Request:
{
    "name": "John Doe",
    "email": "john@example.com",
    "role": "admin"
}

Response:
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "role": "admin"
    }
}
```

## 开发指南

### 1. 控制器结构

```go
// internal/controllers/user_controller.go
type UserController struct {
    userService *services.UserService
}

func (c *UserController) List(ctx *gin.Context) {
    page := ctx.DefaultQuery("page", "1")
    perPage := ctx.DefaultQuery("per_page", "20")
    
    users, pagination, err := c.userService.List(page, perPage)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    ctx.JSON(200, gin.H{
        "code": 0,
        "message": "success",
        "data": gin.H{
            "items": users,
            "pagination": pagination,
        },
    })
}
```

### 2. 请求验证

```go
// internal/requests/user_request.go
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

func (r *CreateUserRequest) Validate() error {
    if strings.TrimSpace(r.Name) == "" {
        return errors.New("name is required")
    }
    // 更多验证规则...
    return nil
}
```

### 3. 响应封装

```go
// pkg/response/response.go
func Success(data interface{}) gin.H {
    return gin.H{
        "code": 0,
        "message": "success",
        "data": data,
    }
}

func Error(code int, message string) gin.H {
    return gin.H{
        "code": code,
        "message": message,
    }
}

func ValidationError(errors map[string][]string) gin.H {
    return gin.H{
        "code": 422,
        "message": "Validation failed",
        "errors": errors,
    }
}
```

## API 测试

### 1. 使用 curl 测试

```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'

# 获取用户列表
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 2. 使用 Postman

1. 导入 Postman 集合（在 `docs/postman` 目录）
2. 设置环境变量
3. 运行测试集合

## 最佳实践

1. API 设计：
   - 使用合适的 HTTP 方法
   - 保持 URL 结构一致
   - 提供详细的错误信息

2. 安全性：
   - 始终验证输入数据
   - 实现速率限制
   - 使用 HTTPS
   - 验证文件上传

3. 性能：
   - 合理使用缓存
   - 优化数据库查询
   - 实现数据分页
   - 压缩响应数据

4. 文档：
   - 保持文档更新
   - 提供示例代码
   - 说明所有参数
   - 列出可能的错误 