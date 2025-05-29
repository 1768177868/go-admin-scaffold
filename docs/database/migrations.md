# 数据库迁移

数据库迁移用于管理数据库结构的版本控制。本框架提供了简单而强大的迁移系统。

## 生成迁移文件

使用命令行工具创建迁移文件：

```bash
# 创建新的迁移文件
go run cmd/tools/main.go make:migration create_users_table

# 创建带外键的迁移文件
go run cmd/tools/main.go make:migration add_role_id_to_users
```

生成的文件位于 `internal/database/migrations` 目录。

## 迁移文件结构

```go
package migrations

import (
    "app/internal/models"
    "gorm.io/gorm"
)

func init() {
    migrations["create_users_table"] = &Migration{
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&models.User{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&models.User{})
        },
    }
}
```

## 运行迁移

### 开发环境

```bash
# 运行所有未执行的迁移
go run cmd/tools/main.go migrate

# 回滚最后一次迁移
go run cmd/tools/main.go migrate:rollback

# 回滚所有迁移
go run cmd/tools/main.go migrate:reset

# 回滚并重新运行所有迁移
go run cmd/tools/main.go migrate:refresh
```

### 生产环境

```bash
# 使用编译后的工具
./bin/dbtools migrate
```

## 编写迁移

### 创建表

```go
func init() {
    migrations["create_posts_table"] = &Migration{
        Up: func(db *gorm.DB) error {
            type Post struct {
                ID        uint      `gorm:"primarykey"`
                Title     string    `gorm:"size:255;not null"`
                Content   string    `gorm:"type:text"`
                UserID    uint      `gorm:"not null"`
                Status    int       `gorm:"default:1"`
                CreatedAt time.Time
                UpdatedAt time.Time
                DeletedAt gorm.DeletedAt `gorm:"index"`
            }
            return db.AutoMigrate(&Post{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable("posts")
        },
    }
}
```

### 修改表

```go
func init() {
    migrations["add_category_to_posts"] = &Migration{
        Up: func(db *gorm.DB) error {
            // 添加字段
            return db.Migrator().AddColumn("posts", "category")
        },
        Down: func(db *gorm.DB) error {
            // 删除字段
            return db.Migrator().DropColumn("posts", "category")
        },
    }
}
```

### 添加索引

```go
func init() {
    migrations["add_posts_indexes"] = &Migration{
        Up: func(db *gorm.DB) error {
            // 添加索引
            return db.Exec("CREATE INDEX idx_posts_title ON posts(title)").Error
        },
        Down: func(db *gorm.DB) error {
            // 删除索引
            return db.Exec("DROP INDEX idx_posts_title ON posts").Error
        },
    }
}
```

## 最佳实践

1. 迁移文件命名：
   - 使用时间戳前缀
   - 使用描述性名称
   - 例如：`20240301120000_create_users_table.go`

2. 迁移设计：
   - 每个迁移文件只做一件事
   - 确保 Up 和 Down 方法是对应的
   - 添加必要的注释说明变更原因

3. 生产环境注意事项：
   - 总是在测试环境验证迁移
   - 备份生产数据库
   - 选择合适的时间执行迁移
   - 准备回滚方案

4. 性能考虑：
   - 大表添加索引时注意性能影响
   - 考虑是否需要分批处理
   - 避免长时间锁表操作

## 常见问题

### 1. 迁移失败处理

如果迁移过程中出现错误：

```bash
# 查看迁移状态
go run cmd/tools/main.go migrate:status

# 手动修复数据库后，标记迁移为已完成
go run cmd/tools/main.go migrate:fix
```

### 2. 数据库版本控制

- 迁移历史记录在 `migrations` 表中
- 每个迁移文件都有唯一的ID
- 系统会追踪已执行的迁移

### 3. 开发建议

- 不要修改已提交的迁移文件
- 如需修改，创建新的迁移文件
- 保持迁移文件的向后兼容性
- 在代码审查中特别关注迁移文件 