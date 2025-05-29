# 认证系统

本框架提供了完整的用户认证和授权系统，支持多种认证方式和灵活的权限控制。

## 认证配置

在 `config/config.yaml` 中配置认证相关参数：

```yaml
auth:
  jwt:
    secret: "your-jwt-secret"
    ttl: 86400          # Token 有效期（秒）
    refresh_ttl: 604800 # 刷新 Token 有效期（秒）
  
  password:
    min_length: 8
    require_number: true
    require_uppercase: true
    require_special: true
  
  throttle:
    max_attempts: 5
    decay_minutes: 10
```

## 用户认证

### 1. 登录认证

```go
// internal/controllers/auth_controller.go
func (c *AuthController) Login(ctx *gin.Context) {
    var req LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    token, err := c.authService.Login(req.Email, req.Password)
    if err != nil {
        ctx.JSON(401, gin.H{"error": "Invalid credentials"})
        return
    }

    ctx.JSON(200, gin.H{
        "token": token,
        "type": "Bearer",
    })
}
```

### 2. 注册用户

```go
func (c *AuthController) Register(ctx *gin.Context) {
    var req RegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    user, err := c.authService.Register(req)
    if err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(201, user)
}
```

### 3. 认证中间件

```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "No token provided"})
            return
        }

        // 验证 token
        claims, err := jwt.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
            return
        }

        // 设置用户信息到上下文
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

## 授权系统

### 1. 角色和权限

```go
// internal/models/role.go
type Role struct {
    ID          uint   `gorm:"primarykey"`
    Name        string `gorm:"size:50;not null;unique"`
    Description string `gorm:"size:255"`
    Permissions []Permission `gorm:"many2many:role_permissions;"`
}

// internal/models/permission.go
type Permission struct {
    ID          uint   `gorm:"primarykey"`
    Name        string `gorm:"size:50;not null;unique"`
    Description string `gorm:"size:255"`
}
```

### 2. 权限检查

```go
// internal/middleware/permission.go
func RequirePermission(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetUint("user_id")
        
        hasPermission, err := auth.CheckPermission(userID, permission)
        if err != nil || !hasPermission {
            c.AbortWithStatusJSON(403, gin.H{"error": "Permission denied"})
            return
        }
        
        c.Next()
    }
}
```

### 3. 使用示例

```go
// 路由配置
router.POST("/users", 
    middleware.Auth(),
    middleware.RequirePermission("users.create"),
    userController.Create,
)
```

## 密码管理

### 1. 密码哈希

```go
// pkg/auth/password.go
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### 2. 密码重置

```go
func (s *AuthService) ResetPassword(token, newPassword string) error {
    // 验证重置令牌
    if !s.ValidateResetToken(token) {
        return errors.New("invalid reset token")
    }

    // 更新密码
    hashedPassword, err := auth.HashPassword(newPassword)
    if err != nil {
        return err
    }

    return s.userRepo.UpdatePassword(userID, hashedPassword)
}
```

## 社交登录

### 1. OAuth2 配置

```yaml
auth:
  oauth:
    google:
      client_id: "your-client-id"
      client_secret: "your-client-secret"
      redirect_url: "http://localhost:8080/auth/google/callback"
    github:
      client_id: "your-client-id"
      client_secret: "your-client-secret"
      redirect_url: "http://localhost:8080/auth/github/callback"
```

### 2. 实现社交登录

```go
func (c *AuthController) SocialLogin(ctx *gin.Context) {
    provider := ctx.Param("provider")
    redirectURL := c.authService.GetSocialLoginURL(provider)
    ctx.Redirect(302, redirectURL)
}

func (c *AuthController) SocialCallback(ctx *gin.Context) {
    provider := ctx.Param("provider")
    code := ctx.Query("code")

    userInfo, err := c.authService.ProcessSocialLogin(provider, code)
    if err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 生成 JWT token
    token, err := c.authService.GenerateToken(userInfo.ID)
    if err != nil {
        ctx.JSON(500, gin.H{"error": "Could not generate token"})
        return
    }

    ctx.JSON(200, gin.H{
        "token": token,
        "user": userInfo,
    })
}
```

## 安全最佳实践

1. 密码安全：
   - 使用强密码策略
   - 密码必须加密存储
   - 定期提醒用户更新密码

2. Token 安全：
   - 使用安全的 JWT 密钥
   - 合理设置 Token 过期时间
   - 实现 Token 刷新机制

3. 访问控制：
   - 实施最小权限原则
   - 定期审查用户权限
   - 记录重要操作日志

4. 防护措施：
   - 实现登录尝试限制
   - 启用 CSRF 保护
   - 使用 HTTPS
   - 实现 IP 黑名单

## 常见问题

### 1. Token 过期处理

```go
func RefreshToken(ctx *gin.Context) {
    oldToken := ctx.GetHeader("Authorization")
    newToken, err := auth.RefreshToken(oldToken)
    if err != nil {
        ctx.JSON(401, gin.H{"error": "Could not refresh token"})
        return
    }
    
    ctx.JSON(200, gin.H{"token": newToken})
}
```

### 2. 会话管理

```go
func (s *AuthService) InvalidateAllSessions(userID uint) error {
    // 将用户的所有 token 加入黑名单
    return s.tokenBlacklist.Add(userID)
}

func (s *AuthService) CheckTokenBlacklist(token string) bool {
    return s.tokenBlacklist.Contains(token)
}
```

### 3. 权限缓存

```go
func (s *AuthService) GetUserPermissions(userID uint) ([]string, error) {
    // 先从缓存获取
    if permissions, exists := s.cache.Get(fmt.Sprintf("permissions:%d", userID)); exists {
        return permissions.([]string), nil
    }

    // 从数据库获取并缓存
    permissions, err := s.permissionRepo.GetUserPermissions(userID)
    if err != nil {
        return nil, err
    }

    s.cache.Set(fmt.Sprintf("permissions:%d", userID), permissions, time.Hour)
    return permissions, nil
}
``` 