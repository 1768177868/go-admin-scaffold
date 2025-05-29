# Todo 应用示例

本教程将指导您使用我们的管理框架创建一个简单的待办事项(Todo)应用，展示基本的增删改查(CRUD)操作。

## 项目结构

```
.
├── internal
│   ├── models
│   │   └── todo.go           # Todo 数据模型
│   ├── handlers
│   │   └── todo_handler.go   # Todo 处理器
│   └── routes
│       └── todo_routes.go    # 路由配置
└── database
    └── migrations
        └── create_todos_table.go  # 数据库迁移
```

## 1. 创建数据模型

首先，创建 Todo 模型 (`internal/models/todo.go`):

```go
package models

import (
    "time"
    "app/pkg/database"
)

type Todo struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

func (Todo) TableName() string {
    return "todos"
}
```

## 2. 创建数据库迁移

创建迁移文件 (`database/migrations/create_todos_table.go`):

```go
package migrations

import (
    "app/internal/models"
    "app/pkg/database"
    "gorm.io/gorm"
)

func init() {
    database.RegisterMigration("create_todos_table", func(db *gorm.DB) error {
        return db.AutoMigrate(&models.Todo{})
    })
}
```

## 3. 请求/响应数据结构

### 请求数据结构

```go
// 创建待办事项请求
type CreateTodoRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
}

// 更新待办事项请求
type UpdateTodoRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

// 分页查询参数
type TodoListQuery struct {
    Page     int    `form:"page" binding:"min=1"`
    PageSize int    `form:"page_size" binding:"min=1,max=100"`
    Status   string `form:"status" binding:"omitempty,oneof=all completed uncompleted"`
}
```

### 响应数据结构

```go
// 通用响应结构
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 分页数据响应
type PaginatedResponse struct {
    Total       int64       `json:"total"`
    CurrentPage int         `json:"current_page"`
    PageSize    int         `json:"page_size"`
    Data        interface{} `json:"data"`
}
```

## 4. 实现处理器

修改处理器 (`internal/handlers/todo_handler.go`) 以支持分页和更详细的请求/响应处理：

```go
package handlers

import (
    "net/http"
    "app/internal/models"
    "app/pkg/response"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type TodoHandler struct{}

// 创建待办事项
func (h *TodoHandler) Create(c *gin.Context) {
    var req CreateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
        return
    }

    todo := models.Todo{
        Title:       req.Title,
        Description: req.Description,
    }

    if err := database.DB.Create(&todo).Error; err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to create todo: "+err.Error())
        return
    }

    response.Success(c, todo)
}

// 获取待办事项列表（支持分页）
func (h *TodoHandler) List(c *gin.Context) {
    var query TodoListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        response.Error(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
        return
    }

    // 设置默认值
    if query.Page == 0 {
        query.Page = 1
    }
    if query.PageSize == 0 {
        query.PageSize = 10
    }

    // 构建查询
    db := database.DB.Model(&models.Todo{})

    // 根据状态筛选
    switch query.Status {
    case "completed":
        db = db.Where("completed = ?", true)
    case "uncompleted":
        db = db.Where("completed = ?", false)
    }

    // 计算总数
    var total int64
    if err := db.Count(&total).Error; err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to count todos: "+err.Error())
        return
    }

    // 获取分页数据
    var todos []models.Todo
    offset := (query.Page - 1) * query.PageSize
    if err := db.Offset(offset).Limit(query.PageSize).
        Order("created_at DESC").
        Find(&todos).Error; err != nil {
        response.Error(c, http.StatusInternalServerError, "Failed to fetch todos: "+err.Error())
        return
    }

    // 返回分页响应
    response.Success(c, PaginatedResponse{
        Total:       total,
        CurrentPage: query.Page,
        PageSize:    query.PageSize,
        Data:        todos,
    })
}

// 获取单个待办事项
func (h *TodoHandler) Get(c *gin.Context) {
    id := c.Param("id")
    var todo models.Todo

    if err := database.DB.First(&todo, id).Error; err != nil {
        response.Error(c, http.StatusNotFound, "Todo not found")
        return
    }

    response.Success(c, todo)
}

// 更新待办事项
func (h *TodoHandler) Update(c *gin.Context) {
    id := c.Param("id")
    var todo models.Todo

    if err := database.DB.First(&todo, id).Error; err != nil {
        response.Error(c, http.StatusNotFound, "Todo not found")
        return
    }

    if err := c.ShouldBindJSON(&todo); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    if err := database.DB.Save(&todo).Error; err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, todo)
}

// 删除待办事项
func (h *TodoHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    if err := database.DB.Delete(&models.Todo{}, id).Error; err != nil {
        response.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(c, gin.H{"message": "Todo deleted successfully"})
}
```

## 5. API 使用示例

### 创建待办事项
```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "完成项目文档",
    "description": "编写项目的技术文档和使用说明"
  }'
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "title": "完成项目文档",
        "description": "编写项目的技术文档和使用说明",
        "completed": false,
        "created_at": "2024-01-20T10:30:00Z",
        "updated_at": "2024-01-20T10:30:00Z"
    }
}
```

### 获取待办事项列表（带分页）
```bash
# 基本分页
curl "http://localhost:8080/api/todos?page=1&page_size=10"

# 筛选已完成的待办事项
curl "http://localhost:8080/api/todos?page=1&page_size=10&status=completed"

# 筛选未完成的待办事项
curl "http://localhost:8080/api/todos?page=1&page_size=10&status=uncompleted"
```

响应示例：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 25,
        "current_page": 1,
        "page_size": 10,
        "data": [
            {
                "id": 1,
                "title": "完成项目文档",
                "description": "编写项目的技术文档和使用说明",
                "completed": false,
                "created_at": "2024-01-20T10:30:00Z",
                "updated_at": "2024-01-20T10:30:00Z"
            },
            // ... 更多待办事项
        ]
    }
}
```

### 获取单个待办事项
```bash
curl http://localhost:8080/api/todos/1
```

### 更新待办事项
```bash
curl -X PUT http://localhost:8080/api/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "完成项目文档",
    "description": "编写项目的技术文档和使用说明",
    "completed": true
  }'
```

### 删除待办事项
```bash
curl -X DELETE http://localhost:8080/api/todos/1
```

## 6. 运行项目

1. 首先运行数据库迁移：
```bash
go run cmd/artisan/main.go migrate
```

2. 启动服务器：
```bash
go run cmd/server/main.go
```

现在，您可以通过访问 `http://localhost:8080/api/todos` 来测试 Todo 应用的各项功能。

## 总结

这个示例展示了如何使用我们的框架快速构建一个具有完整 CRUD 功能的 RESTful API。通过这个简单的 Todo 应用，您可以了解到：

1. 数据模型的定义和迁移
2. 处理器的实现
3. 路由的配置
4. API 的使用方法

这些基础概念可以帮助您开发更复杂的应用程序。 