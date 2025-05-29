# 队列系统文档

## 概述

Go Admin Scaffold 提供了一个功能强大的队列系统，支持异步任务处理。系统支持多种驱动（Redis、Database），提供任务重试、延迟执行、优先级队列等功能。

## 特性

- 🚀 **多驱动支持**: Redis、Database（MySQL）
- ⏰ **延迟任务**: 支持任务延迟执行
- 🔄 **自动重试**: 任务失败自动重试，支持退避策略
- 📊 **优先级队列**: 支持不同优先级的队列
- 🛠️ **命令行工具**: 提供队列管理和状态查询工具
- 📈 **监控**: 实时查询队列状态和任务数量
- 🔧 **灵活配置**: 支持多种配置选项

## 配置

### 基础配置

在 `configs/config.yaml` 中配置队列：

```yaml
queue:
  # 默认驱动: redis, database
  driver: "redis"
  # 默认队列名称
  queue: "default"
  # 工作进程配置
  worker:
    # 无任务时休眠时间(秒)
    sleep: 3
    # 最大处理任务数(0表示无限制)
    max_jobs: 0
    # 最大运行时间(0表示无限制)
    max_time: 0
    # 处理完一个任务后休息时间(秒)
    rest: 0
    # 内存限制(MB)
    memory: 128
    # 任务最大重试次数
    tries: 3
    # 任务超时时间(秒)
    timeout: 60
  # 队列配置
  queues:
    # 默认队列
    default:
      priority: 1
      processes: 1
      timeout: 60
      tries: 3
      retry_after: 60
      backoff: [60, 300, 900]
    # 高优先级队列
    high:
      priority: 2
      processes: 2
      timeout: 30
      tries: 5
      retry_after: 30
      backoff: [30, 60, 180]
    # 低优先级队列
    low:
      priority: 0
      processes: 1
      timeout: 120
      tries: 2
      retry_after: 120
      backoff: [120, 300, 600]
```

### Redis 配置

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### MySQL 配置

```yaml
mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  database: "go_admin"
```

## 使用方法

### 1. 创建任务

#### 基础任务结构

```go
package jobs

import (
    "app/pkg/queue"
    "time"
)

type ExampleJob struct {
    queue.BaseJob
    Message string `json:"message"`
}

func (j *ExampleJob) Handle() error {
    // 处理任务逻辑
    fmt.Printf("Processing message: %s\n", j.Message)
    return nil
}
```

#### 邮件任务示例

```go
type SendWelcomeEmailJob struct {
    queue.BaseJob
    Email   string `json:"email"`
    Name    string `json:"name"`
    Subject string `json:"subject"`
}

func (j *SendWelcomeEmailJob) Handle() error {
    // 发送欢迎邮件
    return sendEmail(j.Email, j.Name, j.Subject)
}
```

### 2. 推送任务

#### 基本推送

```go
package main

import (
    "context"
    "app/pkg/queue"
    "app/internal/core/jobs"
)

func main() {
    // 创建队列管理器
    config := queue.Config{
        Driver: "redis",
        Options: map[string]interface{}{
            "connection": "redis://localhost:6379/0",
            "queue":      "default",
        },
    }
    
    manager, err := queue.NewManager(config)
    if err != nil {
        panic(err)
    }
    defer manager.Close()
    
    // 创建任务
    job := &jobs.ExampleJob{
        BaseJob: queue.BaseJob{
            Queue:       "default",
            MaxAttempts: 3,
            Timeout:     60 * time.Second,
        },
        Message: "Hello, Queue!",
    }
    
    // 推送任务
    ctx := context.Background()
    err = manager.Push(ctx, job)
    if err != nil {
        panic(err)
    }
}
```

#### 延迟任务

```go
// 延迟5分钟执行
delay := 5 * time.Minute
err = manager.Later(ctx, job, delay)
```

#### 原始数据推送

```go
payload := map[string]interface{}{
    "type": "email",
    "to":   "user@example.com",
    "body": "Hello World",
}

rawData, _ := json.Marshal(payload)
err = manager.PushRaw(ctx, "emails", rawData, map[string]interface{}{
    "delay":        2 * time.Second,
    "max_attempts": 3,
    "timeout":      30 * time.Second,
})
```

### 3. 启动工作进程

#### 使用命令行工具

```bash
# 启动队列服务
./queue-cmd.exe -start

# 停止所有队列
./queue-cmd.exe -stop

# 清空指定队列
./queue-cmd.exe -clear -queue=default
```

#### 编程方式启动

```go
package main

import (
    "app/internal/core/services"
)

func main() {
    // 创建队列服务
    queueService, err := services.NewQueueService()
    if err != nil {
        panic(err)
    }
    
    // 启动服务
    err = queueService.Start()
    if err != nil {
        panic(err)
    }
    
    // 等待信号...
    
    // 停止服务
    queueService.Stop()
}
```

### 4. 自定义工作进程

```go
package main

import (
    "context"
    "app/pkg/queue"
)

func main() {
    // 创建管理器
    manager, _ := queue.NewManager(config)
    
    // 创建工作进程选项
    options := queue.WorkerOptions{
        Sleep:   3 * time.Second,
        MaxJobs: 100,
        Memory:  256,
        Tries:   3,
        Timeout: 60 * time.Second,
    }
    
    // 创建工作进程
    worker := queue.NewWorker(manager, []string{"default", "high"}, options)
    
    // 启动工作进程
    worker.Start()
}
```

## 命令行工具

### 队列管理工具 (queue-cmd.exe)

```bash
# 显示帮助
./queue-cmd.exe

# 启动队列服务
./queue-cmd.exe -start

# 停止队列服务
./queue-cmd.exe -stop

# 列出所有队列
./queue-cmd.exe -list

# 清空指定队列
./queue-cmd.exe -clear -queue=default

# 使用自定义配置文件
./queue-cmd.exe -config=configs/production.yaml -start
```

### 队列状态查询工具 (queue-status.exe)

```bash
# 查看所有队列状态
./queue-status.exe -all

# 查看指定队列状态
./queue-status.exe -queue=default

# 查看指定驱动的队列状态
./queue-status.exe -queue=high -driver=database

# 显示帮助
./queue-status.exe
```

## API 参考

### Queue Manager

#### 创建管理器

```go
func NewManager(config Config) (*Manager, error)
```

#### 推送任务

```go
func (m *Manager) Push(ctx context.Context, job JobInterface) error
func (m *Manager) PushRaw(ctx context.Context, queue string, payload []byte, options map[string]interface{}) error
func (m *Manager) Later(ctx context.Context, job JobInterface, delay time.Duration) error
```

#### 获取任务

```go
func (m *Manager) Pop(ctx context.Context, queue string) (JobInterface, error)
```

#### 队列管理

```go
func (m *Manager) Size(ctx context.Context, queue string) (int64, error)
func (m *Manager) Clear(ctx context.Context, queue string) error
func (m *Manager) Delete(ctx context.Context, queue string, job JobInterface) error
func (m *Manager) Release(ctx context.Context, queue string, job JobInterface, delay time.Duration) error
```

### Job Interface

```go
type JobInterface interface {
    Handle() error
    GetQueue() string
    GetAttempts() int
    GetMaxAttempts() int
    GetDelay() time.Duration
    GetTimeout() time.Duration
    GetRetryAfter() time.Duration
    GetBackoff() []time.Duration
    GetPayload() []byte
    SetPayload(payload []byte)
    SetAttempts(attempts int)
    GetID() string
    SetID(id string)
    SetReservedAt(t *time.Time)
}
```

### Worker Options

```go
type WorkerOptions struct {
    Sleep   time.Duration // 无任务时休眠时间
    MaxJobs int64         // 最大处理任务数
    MaxTime time.Duration // 最大运行时间
    Rest    time.Duration // 处理完任务后休息时间
    Memory  int64         // 内存限制(MB)
    Tries   int           // 最大重试次数
    Timeout time.Duration // 任务超时时间
}
```

## 驱动说明

### Redis 驱动

**优点:**
- 高性能，低延迟
- 支持延迟任务（使用 Sorted Set）
- 内存存储，速度快

**缺点:**
- 数据可能丢失（重启时）
- 内存限制

**适用场景:**
- 高并发场景
- 对性能要求高的任务
- 临时性任务

### Database 驱动

**优点:**
- 数据持久化
- 事务支持
- 数据不会丢失

**缺点:**
- 相对较慢
- 数据库负载

**适用场景:**
- 重要任务
- 需要数据持久化
- 对可靠性要求高

## 最佳实践

### 1. 任务设计

```go
// ✅ 好的做法
type ProcessOrderJob struct {
    queue.BaseJob
    OrderID int64 `json:"order_id"`
}

func (j *ProcessOrderJob) Handle() error {
    // 幂等性处理
    if j.isProcessed() {
        return nil
    }
    
    // 业务逻辑
    return j.processOrder()
}

// ❌ 避免的做法
type BadJob struct {
    queue.BaseJob
    LargeData []byte `json:"large_data"` // 避免大数据
}
```

### 2. 错误处理

```go
func (j *MyJob) Handle() error {
    // 区分可重试和不可重试的错误
    if err := j.doSomething(); err != nil {
        if isRetryableError(err) {
            return err // 会重试
        }
        // 记录日志，不重试
        log.Printf("Non-retryable error: %v", err)
        return nil
    }
    return nil
}
```

### 3. 监控和日志

```go
func (j *MyJob) Handle() error {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("Job %s completed in %v", j.GetID(), duration)
    }()
    
    // 任务逻辑
    return nil
}
```

### 4. 队列选择

- **high**: 重要且紧急的任务（支付、通知）
- **default**: 一般任务（邮件发送、数据处理）
- **low**: 不紧急的任务（日志清理、报表生成）

## 故障排除

### 常见问题

1. **任务不执行**
   - 检查工作进程是否启动
   - 检查队列配置是否正确
   - 查看日志错误信息

2. **任务重复执行**
   - 确保任务具有幂等性
   - 检查任务是否正确删除

3. **内存使用过高**
   - 调整 worker.memory 配置
   - 减少并发工作进程数量
   - 优化任务处理逻辑

4. **Redis 连接失败**
   - 检查 Redis 服务状态
   - 验证连接配置
   - 检查网络连接

5. **数据库连接失败**
   - 检查数据库服务状态
   - 验证数据库配置
   - 检查数据库权限

### 调试技巧

```bash
# 查看队列状态
./queue-status.exe -all

# 查看特定队列
./queue-status.exe -queue=default

# 清空问题队列
./queue-cmd.exe -clear -queue=problematic_queue

# 启动调试模式
./queue-cmd.exe -start -config=configs/debug.yaml
```

## 性能优化

### 1. Redis 优化

```yaml
redis:
  # 使用连接池
  max_idle_conns: 10
  max_open_conns: 100
  # 设置合适的超时
  read_timeout: 3s
  write_timeout: 3s
```

### 2. 工作进程优化

```yaml
queue:
  worker:
    # 根据 CPU 核心数调整
    processes: 4
    # 适当的休眠时间
    sleep: 1
    # 内存限制
    memory: 512
```

### 3. 任务优化

- 保持任务轻量级
- 避免长时间运行的任务
- 使用批处理减少队列操作
- 实现任务幂等性

## 示例项目

查看 `examples/queue/` 目录获取完整示例：

- `main.go`: 基本使用示例
- `worker.go`: 自定义工作进程
- `jobs/`: 各种任务示例

## 更新日志

### v1.0.0
- 初始版本发布
- 支持 Redis 和 Database 驱动
- 基本队列功能
- 命令行工具

## 生产环境部署

### 运行队列服务

在生产环境中，队列系统提供了两个可执行文件用于运行和管理队列：

1. **queue-cmd.exe**: 队列管理工具
2. **queue-status.exe**: 队列状态查看工具
3. **worker.exe**: 队列工作进程

#### 启动队列

方式一：使用队列管理工具
```bash
# 启动队列服务
./queue-cmd.exe -start

# 使用指定配置文件启动
./queue-cmd.exe -config=configs/production.yaml -start

# 停止所有队列
./queue-cmd.exe -stop

# 列出所有队列
./queue-cmd.exe -list

# 清空指定队列
./queue-cmd.exe -clear -queue=default
```

方式二：直接运行工作进程
```bash
# 直接运行工作进程
./worker.exe
```

#### 查看队列状态
```bash
# 查看所有队列状态
./queue-status.exe -all

# 查看指定队列状态
./queue-status.exe -queue=default

# 查看指定驱动的队列状态
./queue-status.exe -queue=high -driver=database
```

### 进程管理

#### Linux (systemd)

创建服务配置文件 `/etc/systemd/system/go-admin-queue.service`:

```ini
[Unit]
Description=Go Admin Queue Worker
After=network.target redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/your/app
ExecStart=/path/to/your/app/worker.exe
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

管理服务：
```bash
# 启动服务
sudo systemctl start go-admin-queue

# 设置开机自启
sudo systemctl enable go-admin-queue

# 查看状态
sudo systemctl status go-admin-queue

# 查看日志
sudo journalctl -u go-admin-queue -f
```

#### Windows (NSSM)

使用 NSSM 安装为 Windows 服务：

```powershell
# 安装服务
nssm install GoAdminQueue "C:\path\to\your\app\worker.exe"
nssm set GoAdminQueue AppDirectory "C:\path\to\your\app"
nssm set GoAdminQueue DisplayName "Go Admin Queue Worker"
nssm set GoAdminQueue Description "Go Admin Queue Worker Service"

# 启动服务
nssm start GoAdminQueue

# 停止服务
nssm stop GoAdminQueue

# 重启服务
nssm restart GoAdminQueue
```

### 生产环境最佳实践

1. **配置文件管理**
   - 使用单独的配置文件（如 `configs/production.yaml`）
   - 确保敏感信息（如数据库密码）通过环境变量注入
   - 根据服务器资源调整配置参数

2. **日志管理**
   - 配置日志输出到文件
   - 使用日志轮转工具（如 logrotate）管理日志文件
   - 定期检查日志，及时发现问题

3. **监控告警**
   - 定期使用 `queue-status.exe` 检查队列状态
   - 设置队列积压告警阈值
   - 监控工作进程状态
   - 配置任务处理失败告警

4. **启动顺序**
   - 确保依赖服务（Redis/MySQL）已启动
   - 先启动主应用服务
   - 然后启动队列工作进程

5. **高可用部署**
   - 考虑使用多个工作进程
   - 配置适当的进程数量（建议 CPU 核心数的 1-2 倍）
   - 使用负载均衡（如需要）

6. **备份策略**
   - 定期备份队列数据（特别是使用 Database 驱动时）
   - 保存重要的队列配置
   - 制定故障恢复方案

7. **性能优化**
   - 根据服务器资源调整工作进程数量
   - 配置合适的任务超时时间
   - 优化任务处理逻辑
   - 使用批处理减少队列操作

8. **安全建议**
   - 限制队列管理工具的访问权限
   - 使用非 root 用户运行服务
   - 定期更新系统和依赖包
   - 配置防火墙规则

### 故障排查

1. **队列不处理任务**
   - 检查工作进程是否运行
   - 验证配置文件是否正确
   - 检查数据库/Redis 连接
   - 查看错误日志

2. **任务积压**
   - 检查工作进程数量是否足够
   - 查看任务处理时间是否过长
   - 检查服务器资源使用情况
   - 考虑增加工作进程或优化任务

3. **内存使用过高**
   - 检查 `worker.memory` 配置
   - 减少并发工作进程数量
   - 优化任务处理逻辑
   - 检查是否存在内存泄漏

4. **连接问题**
   - 检查网络连接
   - 验证数据库/Redis 服务状态
   - 检查连接配置
   - 查看连接池设置

### 常用维护命令

```bash
# 查看队列状态
./queue-status.exe -all

# 清空问题队列
./queue-cmd.exe -clear -queue=problematic_queue

# 重启队列服务
./queue-cmd.exe -stop
./queue-cmd.exe -start

# 查看工作进程日志
tail -f /path/to/queue.log

# 检查系统资源
top
htop
free -m
```

## 本地开发环境运行

### Windows 环境

1. **直接运行**
```powershell
# 启动队列服务
.\queue-cmd.exe -start

# 查看队列状态
.\queue-status.exe -all

# 停止队列服务
.\queue-cmd.exe -stop
```

2. **使用 PowerShell 后台运行**
```powershell
# 启动队列服务到后台
Start-Process -NoNewWindow .\queue-cmd.exe -ArgumentList "-start"

# 查看后台进程
Get-Process queue-cmd

# 停止队列服务
Stop-Process -Name "queue-cmd"
```

3. **使用 Windows 服务（开发环境）**
```powershell
# 安装为 Windows 服务（需要管理员权限）
nssm install GoAdminQueueDev "C:\path\to\your\app\worker.exe"
nssm set GoAdminQueueDev AppDirectory "C:\path\to\your\app"
nssm set GoAdminQueueDev DisplayName "Go Admin Queue Worker (Dev)"
nssm set GoAdminQueueDev Description "Go Admin Queue Worker Service for Development"

# 启动服务
nssm start GoAdminQueueDev

# 停止服务
nssm stop GoAdminQueueDev
```

### Mac 环境

1. **直接运行**
```bash
# 启动队列服务
./queue-cmd -start

# 查看队列状态
./queue-status -all

# 停止队列服务
./queue-cmd -stop
```

2. **使用后台运行**
```bash
# 启动队列服务到后台
nohup ./queue-cmd -start > queue.log 2>&1 &

# 查看进程
ps aux | grep queue-cmd

# 停止队列服务
pkill -f queue-cmd
```

3. **使用 launchd（开发环境）**

创建服务配置文件 `~/Library/LaunchAgents/com.go-admin.queue.dev.plist`:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.go-admin.queue.dev</string>
    <key>ProgramArguments</key>
    <array>
        <string>/path/to/your/app/worker</string>
    </array>
    <key>WorkingDirectory</key>
    <string>/path/to/your/app</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>/path/to/your/app/logs/queue-error.log</string>
    <key>StandardOutPath</key>
    <string>/path/to/your/app/logs/queue.log</string>
</dict>
</plist>
```

管理服务：
```bash
# 加载服务
launchctl load ~/Library/LaunchAgents/com.go-admin.queue.dev.plist

# 启动服务
launchctl start com.go-admin.queue.dev

# 停止服务
launchctl stop com.go-admin.queue.dev

# 卸载服务
launchctl unload ~/Library/LaunchAgents/com.go-admin.queue.dev.plist
```

### 开发环境配置

1. **配置文件**
创建开发环境配置文件 `configs/config.dev.yaml`:
```yaml
queue:
  driver: "redis"  # 或 "database"
  queue: "default"
  worker:
    sleep: 1
    max_jobs: 0
    max_time: 0
    rest: 0
    memory: 256
    tries: 3
    timeout: 30
  queues:
    default:
      priority: 1
      processes: 2
      timeout: 30
      tries: 3
      retry_after: 60
      backoff: [30, 60, 120]

# Redis 开发环境配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

# 或 MySQL 开发环境配置
mysql:
  host: "localhost"
  port: 3306
  username: "root"
  password: "dev_password"
  database: "go_admin_dev"
```

2. **启动命令**
```bash
# Windows
.\queue-cmd.exe -start -config=configs/config.dev.yaml

# Mac
./queue-cmd -start -config=configs/config.dev.yaml
```

3. **开发环境最佳实践**

- **日志输出**
  - 开发环境建议将日志输出到控制台
  - 可以设置更详细的日志级别
  - 使用 `-v` 参数查看详细日志

- **调试模式**
  - 使用 `-debug` 参数启动调试模式
  - 可以查看更详细的处理信息
  - 方便排查问题

- **快速重启**
  - 开发时可以使用 `-watch` 参数
  - 配置文件修改后自动重启
  - 提高开发效率

- **测试数据**
  - 使用 `examples/queue/` 目录下的示例
  - 可以快速测试队列功能
  - 验证配置是否正确

4. **常用开发命令**

Windows:
```powershell
# 启动带调试信息的队列
.\queue-cmd.exe -start -config=configs/config.dev.yaml -debug

# 查看详细日志
.\queue-cmd.exe -start -config=configs/config.dev.yaml -v

# 启动并监视配置变化
.\queue-cmd.exe -start -config=configs/config.dev.yaml -watch
```

Mac:
```bash
# 启动带调试信息的队列
./queue-cmd -start -config=configs/config.dev.yaml -debug

# 查看详细日志
./queue-cmd -start -config=configs/config.dev.yaml -v

# 启动并监视配置变化
./queue-cmd -start -config=configs/config.dev.yaml -watch
```

5. **开发工具集成**

- **VSCode 任务配置**
创建 `.vscode/tasks.json`:
```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Start Queue",
            "type": "shell",
            "command": "${workspaceFolder}/queue-cmd",
            "args": ["-start", "-config=configs/config.dev.yaml", "-debug"],
            "windows": {
                "command": "${workspaceFolder}/queue-cmd.exe"
            },
            "group": {
                "kind": "build",
                "isDefault": true
            }
        },
        {
            "label": "Stop Queue",
            "type": "shell",
            "command": "${workspaceFolder}/queue-cmd",
            "args": ["-stop"],
            "windows": {
                "command": "${workspaceFolder}/queue-cmd.exe"
            }
        }
    ]
}
```

- **JetBrains IDE 运行配置**
  - 创建新的运行配置
  - 设置程序路径为 queue-cmd
  - 添加参数：-start -config=configs/config.dev.yaml -debug
  - 设置工作目录为项目根目录

---

如有问题或建议，请提交 Issue 或 Pull Request。 