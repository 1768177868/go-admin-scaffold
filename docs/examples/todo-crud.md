# Todo 应用示例

本教程将指导您使用我们的管理框架创建一个简单的待办事项(Todo)应用，展示基本的增删改查(CRUD)操作。

## 项目结构

```
.
├── internal
│   ├── core
│   │   ├── models
│   │   │   └── todo.go           # Todo 数据模型
│   │   ├── handlers
│   │   │   └── todo_handler.go   # Todo 处理器
│   │   └── services
│   │       └── todo_service.go   # Todo 业务逻辑
│   ├── api
│   │   └── admin
│   │       └── v1
│   │           └── todo.go       # Todo API 路由处理
│   └── database
│       └── migrations
│           └── create_todos_table.go  # 数据库迁移
```

## 1. 创建数据模型

创建 Todo 模型 (`internal/core/models/todo.go`):

```go
package models

import (
    "app/internal/core/models"
)

type Todo struct {
    models.BaseModel
    Title       string `json:"title" gorm:"not null"`
    Description string `json:"description"`
    Completed   bool   `json:"completed" gorm:"default:false"`
}

func (Todo) TableName() string {
    return "todos"
}
```

## 2. 创建数据库迁移

创建迁移文件 (`internal/database/migrations/20240315_create_todos_table.go`):

```go
package migrations

import (
    "app/internal/core/models"
    "gorm.io/gorm"
)

type CreateTodosTable struct{}

func (m *CreateTodosTable) Up(db *gorm.DB) error {
    return db.AutoMigrate(&models.Todo{})
}

func (m *CreateTodosTable) Down(db *gorm.DB) error {
    return db.Migrator().DropTable(&models.Todo{})
}
```

## 3. 创建服务层

创建 Todo 服务 (`internal/core/services/todo_service.go`):

```go
package services

import (
    "context"
    "app/internal/core/models"
)

type TodoService struct {
    db *gorm.DB
}

func NewTodoService(db *gorm.DB) *TodoService {
    return &TodoService{db: db}
}

type CreateTodoRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
}

type UpdateTodoRequest struct {
    Title       string `json:"title" binding:"required"`
    Description string `json:"description"`
    Completed   bool   `json:"completed"`
}

func (s *TodoService) Create(ctx context.Context, req *CreateTodoRequest) (*models.Todo, error) {
    todo := &models.Todo{
        Title:       req.Title,
        Description: req.Description,
    }
    if err := s.db.Create(todo).Error; err != nil {
        return nil, err
    }
    return todo, nil
}

func (s *TodoService) List(ctx context.Context, pagination *models.Pagination) ([]models.Todo, error) {
    var todos []models.Todo
    query := s.db.Model(&models.Todo{})
    
    if err := query.Count(&pagination.Total).Error; err != nil {
        return nil, err
    }
    
    if err := query.Offset(pagination.GetOffset()).
        Limit(pagination.GetLimit()).
        Order("created_at DESC").
        Find(&todos).Error; err != nil {
        return nil, err
    }
    
    return todos, nil
}

func (s *TodoService) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
    var todo models.Todo
    if err := s.db.First(&todo, id).Error; err != nil {
        return nil, err
    }
    return &todo, nil
}

func (s *TodoService) Update(ctx context.Context, id uint, req *UpdateTodoRequest) (*models.Todo, error) {
    todo, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    todo.Title = req.Title
    todo.Description = req.Description
    todo.Completed = req.Completed
    
    if err := s.db.Save(todo).Error; err != nil {
        return nil, err
    }
    return todo, nil
}

func (s *TodoService) Delete(ctx context.Context, id uint) error {
    return s.db.Delete(&models.Todo{}, id).Error
}
```

## 4. 创建处理器

创建 Todo 处理器 (`internal/core/handlers/todo_handler.go`):

```go
package handlers

import (
    "app/internal/core/services"
    "app/pkg/response"
    "github.com/gin-gonic/gin"
)

type TodoHandler struct {
    todoService *services.TodoService
}

func NewTodoHandler(todoService *services.TodoService) *TodoHandler {
    return &TodoHandler{todoService: todoService}
}

// @Summary Create todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body services.CreateTodoRequest true "Todo info"
// @Success 200 {object} response.Response{data=models.Todo}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /admin/v1/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
    var req services.CreateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    todo, err := h.todoService.Create(c.Request.Context(), &req)
    if err != nil {
        response.ServerError(c)
        return
    }

    response.Success(c, todo)
}

// @Summary List todos
// @Description Get a paginated list of todos
// @Tags todos
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} response.Response{data=response.PageData{list=[]models.Todo}}
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /admin/v1/todos [get]
func (h *TodoHandler) List(c *gin.Context) {
    pagination := &models.Pagination{
        Page:     c.GetInt("page"),
        PageSize: c.GetInt("page_size"),
    }

    todos, err := h.todoService.List(c.Request.Context(), pagination)
    if err != nil {
        response.ServerError(c)
        return
    }

    response.PageSuccess(c, todos, pagination.Total, pagination.Page, pagination.PageSize)
}

// @Summary Get todo
// @Description Get todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} response.Response{data=models.Todo}
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /admin/v1/todos/{id} [get]
func (h *TodoHandler) Get(c *gin.Context) {
    id := c.GetUint("id")
    todo, err := h.todoService.GetByID(c.Request.Context(), id)
    if err != nil {
        response.NotFoundError(c)
        return
    }

    response.Success(c, todo)
}

// @Summary Update todo
// @Description Update todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body services.UpdateTodoRequest true "Todo info"
// @Success 200 {object} response.Response{data=models.Todo}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /admin/v1/todos/{id} [put]
func (h *TodoHandler) Update(c *gin.Context) {
    id := c.GetUint("id")
    var req services.UpdateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    todo, err := h.todoService.Update(c.Request.Context(), id, &req)
    if err != nil {
        response.ServerError(c)
        return
    }

    response.Success(c, todo)
}

// @Summary Delete todo
// @Description Delete todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Security Bearer
// @Router /admin/v1/todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
    id := c.GetUint("id")
    if err := h.todoService.Delete(c.Request.Context(), id); err != nil {
        response.ServerError(c)
        return
    }

    response.Success(c, nil)
}
```

## 5. 注册路由

在 `internal/routes/router.go` 中添加 Todo 路由：

```go
// 在 SetupRoutes 函数中的 adminV1Protected 路由组中添加：
todos := adminV1Protected.Group("/todos")
todos.Use(middleware.RBAC("todo:manage"))
{
    todoHandler := handlers.NewTodoHandler(services.NewTodoService(database.GetDB()))
    todos.POST("", wrapHandler(todoHandler.Create))
    todos.GET("", wrapHandler(todoHandler.List))
    todos.GET("/:id", wrapHandler(todoHandler.Get))
    todos.PUT("/:id", wrapHandler(todoHandler.Update))
    todos.DELETE("/:id", wrapHandler(todoHandler.Delete))
}
```

## 6. API 使用示例

### 创建待办事项
```bash
curl -X POST http://localhost:8080/api/admin/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "title": "完成项目文档",
    "description": "编写项目的技术文档和使用说明"
  }'
```

### 获取待办事项列表
```bash
curl "http://localhost:8080/api/admin/v1/todos?page=1&page_size=10" \
  -H "Authorization: Bearer your-token"
```

### 获取单个待办事项
```bash
curl "http://localhost:8080/api/admin/v1/todos/1" \
  -H "Authorization: Bearer your-token"
```

### 更新待办事项
```bash
curl -X PUT http://localhost:8080/api/admin/v1/todos/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -d '{
    "title": "完成项目文档",
    "description": "编写项目的技术文档和使用说明",
    "completed": true
  }'
```

### 删除待办事项
```bash
curl -X DELETE http://localhost:8080/api/admin/v1/todos/1 \
  -H "Authorization: Bearer your-token"
```

## 7. 运行项目

1. 运行数据库迁移：
```bash
go run cmd/tools/main.go migrate
```

2. 启动服务器：
```bash
go run cmd/server/main.go
```

3. 访问 Swagger 文档：
```
http://localhost:8080/swagger/index.html
```

## 总结

这个示例展示了如何使用我们的框架构建一个完整的 Todo CRUD 应用，包括：

1. 分层架构：模型层、服务层、处理器层
2. 数据库迁移
3. 权限控制
4. API 文档
5. 统一的响应格式
6. 错误处理

通过这个示例，您可以了解到框架的主要特性和最佳实践。 