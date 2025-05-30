# 日志系统使用文档

## 基本用法

日志系统支持链式调用配置和结构化日志记录。基于 `zap` 和 `lumberjack` 实现，支持日志分割、压缩和自动清理等特性。

### 1. 创建日志实例

```go
logger, err := logger.NewBuilder().
    SetDriver("daily").                    // 设置按天分割
    SetPath("storage/logs/app.log").       // 设置日志文件路径
    SetLevel("info").                      // 设置日志级别
    SetMaxSize(100).                       // 设置单个文件最大尺寸(MB)
    SetMaxBackups(10).                     // 设置最大备份数
    SetMaxAge(30).                         // 设置日志保留天数
    SetCompress(true).                     // 设置是否压缩
    Build()

if err != nil {
    panic(err)
}
defer logger.Close()
```

### 2. 记录日志

```go
// 简单日志
logger.Info("Hello World")

// 带结构化数据的日志
logger.Info("User Login", map[string]interface{}{
    "user_id": 123,
    "ip": "192.168.1.1",
    "time": time.Now(),
})

// 错误日志
logger.Error("Database Error", map[string]interface{}{
    "error": err.Error(),
    "query": "SELECT * FROM users",
})
```

### 3. 在中间件中使用

```go
func LogMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 创建请求日志记录器
        reqLogger, err := logger.NewBuilder().
            SetDriver("daily").
            SetPath("storage/logs/request.log").
            Build()
        
        if err != nil {
            log.Printf("Failed to create logger: %v", err)
            c.Next()
            return
        }
        defer reqLogger.Close()

        // 记录请求信息
        reqLogger.Info("HTTP Request", map[string]interface{}{
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

### 4. 不同类型的日志实例

```go
// API请求日志
apiLogger, _ := logger.NewBuilder().
    SetDriver("daily").
    SetPath("storage/logs/api-requests.log").
    Build()

// 错误日志
errorLogger, _ := logger.NewBuilder().
    SetDriver("daily").
    SetPath("storage/logs/errors.log").
    SetLevel("error").
    Build()

// 调试日志
debugLogger, _ := logger.NewBuilder().
    SetDriver("daily").
    SetPath("storage/logs/debug.log").
    SetLevel("debug").
    Build()
```

## 特性说明

1. **按天分割**：
   - 启用 `SetDriver("daily")` 后，日志文件会自动按天分割
   - 文件名格式：`{原文件名}-{日期}.log`
   - 例如：`app-2024-03-14.log`

2. **日志级别**：
   - debug：调试信息
   - info：一般信息
   - warn：警告信息
   - error：错误信息

3. **结构化数据**：
   - 支持记录复杂的结构化数据
   - 自动转换为 JSON 格式
   - 支持任意嵌套的 map 和 slice

4. **自动清理**：
   - `SetMaxSize`：单个文件最大尺寸(MB)
   - `SetMaxBackups`：保留的旧文件数量
   - `SetMaxAge`：日志保留天数
   - `SetCompress`：是否压缩旧文件

5. **性能优化**：
   - 基于高性能的 zap 日志库
   - 异步写入
   - 自动批处理

## 最佳实践

1. **合理设置日志级别**：
   - 生产环境建议使用 info 或以上级别
   - 开发环境可以使用 debug 级别

2. **分类记录日志**：
   - 按功能分类：请求日志、错误日志、业务日志等
   - 按级别分类：debug、info、error 等

3. **结构化记录**：
   - 使用结构化数据而不是简单字符串
   - 包含必要的上下文信息
   - 使用统一的字段命名

4. **及时清理**：
   - 设置合理的 MaxAge 和 MaxBackups
   - 启用压缩以节省空间

5. **错误处理**：
   - 始终检查 logger 创建的错误
   - 使用 defer Close() 确保资源释放 