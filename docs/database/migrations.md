# 数据库迁移和填充

本框架提供了类似 Laravel 的强大数据库迁移和填充系统，支持版本控制、依赖管理和状态跟踪。

## 迁移功能

### 创建迁移

使用命令行工具创建迁移文件：

```bash
# 创建新的迁移文件
go run cmd/tools/main.go make:migration create_users_table

# 创建带外键的迁移文件
go run cmd/tools/main.go make:migration add_role_id_to_users
```

### 迁移文件结构

```go
func init() {
    migrations["20240310_create_users_table"] = &MigrationDefinition{
        Up: func(tx *gorm.DB) error {
            // 创建表结构
            return tx.AutoMigrate(&User{})
        },
        Down: func(tx *gorm.DB) error {
            // 回滚操作
            return tx.Migrator().DropTable("users")
        },
    }
}
```

### 运行迁移

```bash
# 查看迁移状态
go run cmd/tools/main.go migrate status

# 运行所有未执行的迁移
go run cmd/tools/main.go migrate run

# 回滚最后一批迁移
go run cmd/tools/main.go migrate rollback

# 回滚所有迁移
go run cmd/tools/main.go migrate reset

# 回滚并重新运行所有迁移
go run cmd/tools/main.go migrate refresh
```

## 数据填充

### 创建填充器

```go
func init() {
    Register("users", &Seeder{
        Name: "users",
        Description: "Create default users",
        Dependencies: []string{"roles"}, // 依赖于roles填充器
        Run: func(tx *gorm.DB) error {
            // 填充数据
            return tx.Create(&users).Error
        },
    })
}
```

### 运行填充

```bash
# 查看填充状态
go run cmd/tools/main.go seed status

# 运行所有填充器
go run cmd/tools/main.go seed run

# 运行指定填充器
go run cmd/tools/main.go seed run users roles

# 重置填充数据
go run cmd/tools/main.go seed reset
```

## 最佳实践

### 迁移文件命名

- 使用时间戳前缀：`20240310_create_users_table.go`
- 使用描述性名称：`create_users_table`, `add_role_id_to_users`
- 一个迁移文件只做一件事

### 迁移设计

1. 表结构设计
```go
type User struct {
    ID        uint      `gorm:"primarykey"`
    Username  string    `gorm:"size:50;not null;unique"`
    Email     string    `gorm:"size:100;not null;unique"`
    Status    int       `gorm:"default:1"`
    CreatedAt time.Time
}
```

2. 索引设计
```go
// 添加索引
tx.Exec(`CREATE INDEX idx_users_email ON users(email)`)

// 添加复合索引
tx.Exec(`CREATE INDEX idx_users_status_created ON users(status, created_at)`)
```

3. 外键关系
```go
// 添加外键
tx.Exec(`ALTER TABLE users ADD CONSTRAINT fk_users_role 
         FOREIGN KEY (role_id) REFERENCES roles(id)`)
```

### 填充器设计

1. 依赖管理
```go
Register("users", &Seeder{
    Dependencies: []string{"roles", "permissions"},
})
```

2. 批量插入
```go
tx.Create(&[]User{
    {Username: "user1", Email: "user1@example.com"},
    {Username: "user2", Email: "user2@example.com"},
})
```

3. 关联数据
```go
// 创建用户和角色关联
tx.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, 
        userId, roleId)
```

## 生产环境注意事项

1. 数据备份
```bash
# 执行迁移前备份数据库
mysqldump -u user -p database > backup.sql
```

2. 迁移验证
- 在测试环境完整运行迁移
- 验证回滚功能
- 检查数据完整性

3. 性能优化
- 大表添加索引时使用 `ALGORITHM=INPLACE`
- 批量插入时使用事务
- 避免长时间锁表操作

4. 监控和日志
- 记录迁移执行时间
- 监控数据库性能
- 保存详细的操作日志

## 常见问题

### 1. 迁移失败处理

如果迁移过程中出现错误：

1. 查看迁移状态
```bash
go run cmd/tools/main.go migrate status
```

2. 手动修复数据库

3. 标记迁移完成
```bash
go run cmd/tools/main.go migrate fix
```

### 2. 填充数据清理

重置特定填充数据：

```bash
# 重置用户数据
go run cmd/tools/main.go seed reset users
```

### 3. 开发建议

- 保持迁移文件的原子性
- 提供完整的回滚操作
- 记录修改原因和影响
- 遵循数据库设计最佳实践 