# Todo 应用示例

本教程将指导您使用我们的管理框架创建一个简单的待办事项(Todo)应用，展示基本的增删改查(CRUD)操作。

## 项目结构

```
.
├── internal
│   ├── core
│   │   ├── models
│   │   │   └── todo.go           # Todo 数据模型
│   │   ├── repositories
│   │   │   └── todo_repository.go # Todo 数据访问层
│   │   └── services
│   │       └── todo_service.go   # Todo 业务逻辑
│   ├── api
│   │   └── admin
│   │       └── v1
│   │           └── todo.go       # Todo 处理器和路由
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

## 3. 创建数据访问层

创建 Todo 仓库 (`internal/core/repositories/todo_repository.go`):

```go
package repositories

import (
    "context"
    "app/internal/core/models"
    "gorm.io/gorm"
)

type TodoRepository struct {
    db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
    return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
    return r.db.Create(todo).Error
}

func (r *TodoRepository) List(ctx context.Context, pagination *models.Pagination) ([]models.Todo, error) {
    var todos []models.Todo
    query := r.db.Model(&models.Todo{})
    
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

func (r *TodoRepository) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
    var todo models.Todo
    if err := r.db.First(&todo, id).Error; err != nil {
        return nil, err
    }
    return &todo, nil
}

func (r *TodoRepository) Update(ctx context.Context, todo *models.Todo) error {
    return r.db.Save(todo).Error
}

func (r *TodoRepository) Delete(ctx context.Context, id uint) error {
    return r.db.Delete(&models.Todo{}, id).Error
}
```

## 4. 创建服务层

创建 Todo 服务 (`internal/core/services/todo_service.go`):

```go
package services

import (
    "context"
    "app/internal/core/models"
    "app/internal/core/repositories"
)

type TodoService struct {
    repo *repositories.TodoRepository
}

func NewTodoService(repo *repositories.TodoRepository) *TodoService {
    return &TodoService{repo: repo}
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
    if err := s.repo.Create(ctx, todo); err != nil {
        return nil, err
    }
    return todo, nil
}

func (s *TodoService) List(ctx context.Context, pagination *models.Pagination) ([]models.Todo, error) {
    return s.repo.List(ctx, pagination)
}

func (s *TodoService) GetByID(ctx context.Context, id uint) (*models.Todo, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *TodoService) Update(ctx context.Context, id uint, req *UpdateTodoRequest) (*models.Todo, error) {
    todo, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    todo.Title = req.Title
    todo.Description = req.Description
    todo.Completed = req.Completed
    
    if err := s.repo.Update(ctx, todo); err != nil {
        return nil, err
    }
    return todo, nil
}

func (s *TodoService) Delete(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
}
```

## 5. 创建处理器

创建 Todo 处理器 (`internal/api/admin/v1/todo.go`):

```go
package v1

import (
    "app/internal/core/services"
    "app/pkg/response"
    "github.com/gin-gonic/gin"
)

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
func CreateTodo(c *gin.Context) {
    var req services.CreateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    todo, err := todoService.Create(c.Request.Context(), &req)
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
func ListTodos(c *gin.Context) {
    pagination := &models.Pagination{
        Page:     c.GetInt("page"),
        PageSize: c.GetInt("page_size"),
    }

    todos, err := todoService.List(c.Request.Context(), pagination)
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
func GetTodo(c *gin.Context) {
    id := c.GetUint("id")
    todo, err := todoService.GetByID(c.Request.Context(), id)
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
func UpdateTodo(c *gin.Context) {
    id := c.GetUint("id")
    var req services.UpdateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err.Error())
        return
    }

    todo, err := todoService.Update(c.Request.Context(), id, &req)
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
func DeleteTodo(c *gin.Context) {
    id := c.GetUint("id")
    if err := todoService.Delete(c.Request.Context(), id); err != nil {
        response.ServerError(c)
        return
    }

    response.Success(c, nil)
}
```

## 6. 注册路由

在 `internal/api/admin/v1/todo.go` 中添加路由注册：

```go
// RegisterTodoRoutes 注册 Todo 相关路由
func RegisterTodoRoutes(r *gin.RouterGroup) {
    todos := r.Group("/todos")
    todos.Use(middleware.RBAC("todo:manage"))
    {
        todos.POST("", wrapHandler(CreateTodo))
        todos.GET("", wrapHandler(ListTodos))
        todos.GET("/:id", wrapHandler(GetTodo))
        todos.PUT("/:id", wrapHandler(UpdateTodo))
        todos.DELETE("/:id", wrapHandler(DeleteTodo))
    }
}
```

然后在 `internal/routes/router.go` 中调用：

```go
// 在 SetupRoutes 函数中的 adminV1Protected 路由组中添加：
adminv1.RegisterTodoRoutes(adminV1Protected)
```

## 7. API 使用示例

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

## 8. 运行项目

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

1. 分层架构：
   - 模型层 (Models)
   - 数据访问层 (Repositories)
   - 服务层 (Services)
   - 处理器层 (Handlers)
   - API 路由层 (Routes)
2. 数据库迁移
3. 权限控制
4. API 文档
5. 统一的响应格式
6. 错误处理

通过这个示例，您可以了解到框架的主要特性和最佳实践。每一层都有其特定的职责：
- Models: 定义数据结构和验证规则
- Repositories: 处理数据持久化，封装数据库操作
- Services: 实现业务逻辑，协调多个仓库
- Handlers: 处理 HTTP 请求，参数验证，调用服务层
- Routes: 定义 API 路由，配置中间件

这个示例展示了如何使用我们的框架构建一个完整的 Todo CRUD 应用，包括：

1. 分层架构：
   - 模型层 (Models)
   - 数据访问层 (Repositories)
   - 服务层 (Services)
   - 处理器层 (Handlers)
   - API 路由层 (Routes)
2. 数据库迁移
3. 权限控制
4. API 文档
5. 统一的响应格式
6. 错误处理

通过这个示例，您可以了解到框架的主要特性和最佳实践。每一层都有其特定的职责：
- Models: 定义数据结构和验证规则
- Repositories: 处理数据持久化，封装数据库操作
- Services: 实现业务逻辑，协调多个仓库
- Handlers: 处理 HTTP 请求，参数验证，调用服务层
- Routes: 定义 API 路由，配置中间件 