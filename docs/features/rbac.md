# 角色权限控制 (RBAC)

## 概述

本项目使用基于角色的访问控制（Role-Based Access Control，RBAC）来管理用户权限。

## 核心概念

- **用户（User）**：系统的使用者
- **角色（Role）**：用户的分组，如管理员、普通用户等
- **权限（Permission）**：具体的操作权限，如创建用户、删除文章等
- **资源（Resource）**：系统中的实体，如用户、文章等

## 数据结构

```go
type User struct {
    ID       uint     `json:"id"`
    Username string   `json:"username"`
    Roles    []Role   `json:"roles"`
}

type Role struct {
    ID          uint         `json:"id"`
    Name        string       `json:"name"`
    Permissions []Permission `json:"permissions"`
}

type Permission struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}
```

## API 接口

### 角色管理

```http
# 创建角色
POST /api/v1/roles
Content-Type: application/json

{
    "name": "editor",
    "description": "Content editor role",
    "permissions": [1, 2, 3]
}

# 获取角色列表
GET /api/v1/roles

# 更新角色
PUT /api/v1/roles/:id

# 删除角色
DELETE /api/v1/roles/:id
```

### 权限管理

```http
# 获取所有权限
GET /api/v1/permissions

# 为角色分配权限
POST /api/v1/roles/:id/permissions
```

## 使用示例

### 权限检查中间件

```go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := auth.GetCurrentUser(c)
        if !user.HasPermission(permission) {
            c.AbortWithStatus(403)
            return
        }
        c.Next()
    }
}
```

### 在路由中使用

```go
router.POST("/articles", 
    middleware.RequirePermission("articles.create"),
    articleController.Create,
)
```

## 最佳实践

1. 遵循最小权限原则
2. 定期审查角色权限
3. 记录权限变更日志
4. 实现权限缓存机制

## 常见问题

1. 如何处理角色继承
2. 动态权限更新
3. 权限缓存更新

## 相关文档

- [用户认证](authentication.md)
- [错误处理](../advanced/error-handling.md)
- [API 文档](../api/README.md) 