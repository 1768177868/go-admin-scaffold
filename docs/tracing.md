# 链路追踪使用文档

## 概述

链路追踪系统通过 `trace_id` 实现请求的全链路跟踪。每个请求从进入系统到响应的整个生命周期都可以通过同一个 `trace_id` 进行追踪。

## Trace ID 的生成和传递

### 1. 自动生成场景

1. **HTTP 请求入口**：
   - 中间件 `middleware.Trace()` 会自动为每个新的 HTTP 请求生成 trace_id
   - 如果请求头中已包含 `X-Trace-ID`，则会复用该值
   - 生成的 trace_id 会被注入到 gin.Context 中

```go
// middleware/trace.go
func Trace() gin.HandlerFunc {
    return func(c *gin.Context) {
        traceID := c.GetHeader("X-Trace-ID")
        if traceID == "" {
            traceID = generateTraceID() // 生成新的 trace_id
        }
        c.Set("trace_id", traceID)
        c.Header("X-Trace-ID", traceID)
        c.Next()
    }
}
```

2. **内部 RPC 调用**：
   - 服务间调用时自动传递 trace_id
   - 调用方将当前请求的 trace_id 放入请求头
   - 被调用方提取并继续使用该 trace_id

### 2. 日志记录

在日志中记录 trace_id，方便追踪请求链路：

```go
// 中间件中记录请求日志
func LogMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        reqLogger, _ := logger.NewBuilder().
            SetDriver("daily").
            SetPath("storage/logs/request.log").
            Build()
        defer reqLogger.Close()

        // 记录包含 trace_id 的请求日志
        reqLogger.Info("HTTP Request", map[string]interface{}{
            "trace_id": c.GetString("trace_id"),  // 获取 trace_id
            "route": c.FullPath(),
            "request": map[string]interface{}{
                "url":     c.Request.URL.String(),
                "method":  c.Request.Method,
                "headers": c.Request.Header,
                "ip":      c.ClientIP(),
            },
        })

        c.Next()
    }
}
```

### 3. 响应中的 Trace ID

所有 API 响应都应包含 trace_id：

```go
// response/response.go
type Response struct {
    Code     int         `json:"code"`
    Message  string      `json:"message"`
    Data     interface{} `json:"data,omitempty"`
    TraceID  string      `json:"trace_id"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
        TraceID: c.GetString("trace_id"),
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
        TraceID: c.GetString("trace_id"),
    })
}
```

## 使用示例

### 1. 控制器中使用

```go
func (ctrl *UserController) Login(c *gin.Context) {
    traceID := c.GetString("trace_id")
    
    // 记录业务日志
    ctrl.logger.Info("User login attempt", map[string]interface{}{
        "trace_id": traceID,
        "username": c.PostForm("username"),
        "ip": c.ClientIP(),
    })

    // 处理业务逻辑...

    // 返回响应
    response.Success(c, map[string]interface{}{
        "token": token,
    })
    // 响应会自动包含 trace_id
}
```

### 2. 服务层使用

```go
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
    traceID := ctx.Value("trace_id").(string)
    
    // 记录业务操作
    s.logger.Info("Creating new user", map[string]interface{}{
        "trace_id": traceID,
        "user_id": user.ID,
        "username": user.Username,
    })

    // 数据库操作...
    
    // 发送欢迎邮件
    s.emailService.SendWelcomeEmail(ctx, user.Email) // ctx 中包含 trace_id
    
    return nil
}
```

### 3. 外部服务调用

```go
func (s *ThirdPartyService) CallExternalAPI(ctx context.Context, data interface{}) error {
    traceID := ctx.Value("trace_id").(string)
    
    req, _ := http.NewRequest("POST", "http://api.example.com/v1/data", nil)
    req.Header.Set("X-Trace-ID", traceID) // 传递 trace_id 到外部服务
    
    // 记录外部调用
    s.logger.Info("Calling external API", map[string]interface{}{
        "trace_id": traceID,
        "api": "example.com",
        "data": data,
    })
    
    // 发送请求...
    return nil
}
```

## 最佳实践

1. **始终传递 Context**：
   - 在函数间传递 context.Context
   - context 中包含 trace_id
   - 使用 context 确保链路追踪的连续性

2. **统一日志格式**：
   - 所有日志都应包含 trace_id 字段
   - 使用结构化日志便于分析
   - 保持字段命名的一致性

3. **异步任务处理**：
   - 将 trace_id 传递给异步任务
   - 在任务执行时记录原始 trace_id
   - 保持异步操作的可追踪性

4. **错误处理**：
   - 错误日志必须包含 trace_id
   - 详细记录错误上下文
   - 便于问题定位和分析

5. **监控和分析**：
   - 基于 trace_id 聚合请求链路
   - 分析请求耗时和性能瓶颈
   - 跟踪错误传播路径 