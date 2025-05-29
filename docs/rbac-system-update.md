# RBAC系统升级指南

## 概述

本次更新将简单的JSON权限列表升级为完整的RBAC（基于角色的访问控制）系统，提供更灵活和标准化的权限管理。

## 主要变更

### 1. 数据库结构变更

#### 新增表
- `permissions` - 权限表
  - `id` - 主键
  - `name` - 权限名称（如 `user:create`）
  - `display_name` - 显示名称
  - `description` - 权限描述
  - `module` - 所属模块
  - `action` - 操作类型
  - `resource` - 资源类型
  - `status` - 状态

- `role_permissions` - 角色权限关联表
  - `id` - 主键
  - `role_id` - 角色ID
  - `permission_id` - 权限ID

#### 修改表
- `roles` 表新增字段：
  - `code` - 角色编码
  - `status` - 状态
  - 移除 `perm_list` JSON字段

### 2. 模型更新

#### Permission模型
```go
type Permission struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`        // user:create
    DisplayName string `json:"display_name"` // 创建用户
    Description string `json:"description"`  // 创建新用户
    Module      string `json:"module"`       // user
    Action      string `json:"action"`       // create
    Resource    string `json:"resource"`     // user
    Status      int    `json:"status"`       // 1-启用 0-禁用
    // ...
}
```

#### Role模型
```go
type Role struct {
    ID          uint           `json:"id"`
    Name        string         `json:"name"`
    Code        string         `json:"code"`        // 新增
    Description string         `json:"description"`
    Status      int            `json:"status"`      // 新增
    Permissions []Permission   `json:"permissions"` // 关联权限
    // ...
}
```

### 3. 服务层更新

#### RBACService
- 实现真正的权限检查逻辑
- 支持复杂权限查询
- 提供用户权限获取方法

```go
func (s *RBACService) CheckPermission(ctx context.Context, user interface{}, permission string) (bool, error)
func (s *RBACService) GetUserPermissions(ctx context.Context, userID uint) ([]string, error)
func (s *RBACService) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error)
```

#### PermissionService
全新的权限管理服务，提供权限的CRUD操作。

#### RoleService
更新角色服务，支持权限关联管理：
- 使用 `permission_ids` 替代 `permissions` 字符串数组
- 支持事务性权限分配
- 提供权限查询方法

### 4. API更新

#### 新增权限管理API
```
GET    /api/admin/v1/permissions          # 获取权限列表
POST   /api/admin/v1/permissions          # 创建权限
GET    /api/admin/v1/permissions/:id      # 获取权限详情
PUT    /api/admin/v1/permissions/:id      # 更新权限
DELETE /api/admin/v1/permissions/:id      # 删除权限
GET    /api/admin/v1/permissions/modules  # 按模块获取权限
```

#### 更新路由权限检查
权限名称从旧格式 `users.manage` 更新为新格式 `user:view`、`user:create` 等。

#### 新增用户资料API
```
GET    /api/admin/v1/profile              # 获取当前用户资料
PUT    /api/admin/v1/profile              # 更新当前用户资料
```

### 5. 权限命名规范

采用 `resource:action` 格式：

#### 用户管理
- `user:view` - 查看用户
- `user:create` - 创建用户
- `user:edit` - 编辑用户
- `user:delete` - 删除用户

#### 角色管理
- `role:view` - 查看角色
- `role:create` - 创建角色
- `role:edit` - 编辑角色
- `role:delete` - 删除角色

#### 权限管理
- `permission:view` - 查看权限
- `permission:create` - 创建权限
- `permission:edit` - 编辑权限
- `permission:delete` - 删除权限

#### 系统功能
- `dashboard:view` - 查看仪表盘
- `log:view` - 查看日志
- `profile:view` - 查看个人资料
- `profile:edit` - 编辑个人资料

## 默认权限配置

### 管理员角色
拥有所有权限（16个权限）

### 管理者角色
- dashboard:view
- user:view, user:edit
- role:view
- log:view
- profile:view, profile:edit

### 普通用户角色
- dashboard:view
- profile:view, profile:edit

## 迁移指南

### 1. 运行迁移
```bash
go run cmd/tools/main.go migrate run
```

### 2. 运行数据填充
```bash
go run cmd/tools/main.go seed run
```

### 3. 验证权限
登录系统后，不同角色的用户应该能看到相应的菜单和功能。

## 开发注意事项

### 1. 权限检查
使用中间件进行权限检查：
```go
users.GET("", middleware.RBAC("user:view"), adminv1.ListUsers)
```

### 2. 编程式权限检查
在业务逻辑中检查权限：
```go
rbacSvc := c.MustGet("rbacService").(*services.RBACService)
hasPermission, err := rbacSvc.CheckPermission(ctx, user, "user:create")
```

### 3. 角色管理
创建角色时使用权限ID数组：
```json
{
  "name": "Editor",
  "code": "editor", 
  "description": "Content Editor",
  "permission_ids": [1, 2, 3, 15, 16]
}
```

## 兼容性说明

此次更新是破坏性的，旧的权限检查方式将不再工作。请确保：

1. 更新所有权限检查的代码
2. 重新配置角色权限
3. 测试所有权限相关功能

## 故障排除

### 1. 权限检查失败
- 确认用户已分配正确角色
- 检查角色是否已分配所需权限
- 验证权限名称格式是否正确

### 2. API访问被拒绝
- 检查路由上的RBAC中间件权限名称
- 确认用户拥有对应权限
- 验证JWT token有效性

### 3. 数据不一致
- 运行 `seed reset` 然后 `seed run` 重新初始化数据
- 检查迁移是否全部执行成功 