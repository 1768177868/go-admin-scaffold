# 测试指南

本文档介绍了如何运行和编写项目的测试用例。

## 目录

- [运行测试](#运行测试)
- [测试覆盖率](#测试覆盖率)
- [测试结构](#测试结构)
- [Mock 对象](#mock-对象)
- [测试示例](#测试示例)

## 运行测试

### 运行所有测试

```bash
# 运行所有测试
go test ./...

# 运行所有测试并显示详细信息
go test -v ./...

# 运行所有测试并生成覆盖率报告
go test -cover ./...
```

### 运行特定包的测试

```bash
# 运行服务层测试
go test ./internal/core/services

# 运行 API 处理器测试
go test ./internal/api/admin/v1

# 运行特定测试函数
go test -run TestUserService_ExportUserList ./internal/core/services
```

## 测试覆盖率

生成测试覆盖率报告：

```bash
# 生成覆盖率数据
go test -coverprofile=coverage.out ./...

# 在浏览器中查看覆盖率报告
go tool cover -html=coverage.out
```

## 测试结构

项目的测试文件组织结构如下：

```
internal/
  ├── core/
  │   ├── services/
  │   │   ├── user_svc.go
  │   │   └── user_svc_test.go      # 服务层测试
  │   └── repositories/
  │       ├── user_repo.go
  │       └── user_repo_test.go     # 仓储层测试
  └── api/
      └── admin/
          └── v1/
              ├── user.go
              └── user_test.go       # API 处理器测试
```

## Mock 对象

项目使用 `testify/mock` 包来创建 mock 对象。以下是主要的 mock 对象：

### 1. MockUserRepository

用于测试服务层的仓储层 mock：

```go
type MockUserRepository struct {
    mock.Mock
    *repositories.BaseRepository
}

// 示例：模拟查找用户
mockRepo := new(MockUserRepository)
mockRepo.On("FindByUsername", ctx, "testuser").Return(&models.User{}, nil)
```

### 2. MockUserService

用于测试 API 处理器的服务层 mock：

```go
type MockUserService struct {
    mock.Mock
}

// 示例：模拟创建用户
mockSvc := new(MockUserService)
mockSvc.On("Create", mock.Anything, mock.AnythingOfType("*services.CreateUserRequest")).
    Return(createdUser, nil)
```

## 测试示例

### 1. 服务层测试示例

```go
func TestUserService_ExportUserList(t *testing.T) {
    mockRepo := new(MockUserRepository)
    mockLogSvc := new(MockLogService)
    userSvc := NewUserService(mockRepo, mockLogSvc)

    ctx := context.Background()
    req := &ExportUserListRequest{
        Username:  "user",
        Email:     "@example.com",
        Status:    &[]int{1}[0],
        StartTime: time.Now().Add(-24 * time.Hour),
        EndTime:   time.Now(),
    }

    testUsers := []models.User{
        {
            BaseModel: models.BaseModel{ID: 1},
            Username:  "user1",
            Email:     "user1@example.com",
            Status:    1,
            Roles:     []models.Role{{Name: "admin"}},
        },
    }

    mockRepo.On("GetDB").Return(&gorm.DB{})
    mockRepo.On("Find", mock.Anything).Return(testUsers, nil)
    mockLogSvc.On("RecordOperationLog", mock.Anything, mock.AnythingOfType("*models.OperationLog")).
        Return(nil)

    users, err := userSvc.ExportUserList(ctx, req)
    assert.NoError(t, err)
    assert.Len(t, users, 1)
    assert.Equal(t, "user1", users[0].Username)
}
```

### 2. API 处理器测试示例

```go
func TestExportUsers(t *testing.T) {
    r, mockSvc := setupTestRouter()
    r.POST("/users/export", ExportUsers)

    req := services.ExportUserListRequest{
        Username: "user",
        Email:    "@example.com",
        Status:   &[]int{1}[0],
    }

    users := []models.User{
        {
            BaseModel: models.BaseModel{ID: 1},
            Username:  "user1",
            Email:     "user1@example.com",
            Status:    1,
            Roles:     []models.Role{{Name: "admin"}},
        },
    }

    mockSvc.On("ExportUserList", mock.Anything, mock.AnythingOfType("*services.ExportUserListRequest")).
        Return(users, nil)

    body, _ := json.Marshal(req)
    w := httptest.NewRecorder()
    httpReq, _ := http.NewRequest("POST", "/users/export", bytes.NewBuffer(body))
    httpReq.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, httpReq)

    assert.Equal(t, http.StatusOK, w.Code)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, float64(0), response["code"])
    assert.Equal(t, "success", response["msg"])
}
```

## 最佳实践

1. 每个测试函数应该专注于测试一个特定的功能或场景
2. 使用有意义的测试函数名称，如 `TestUserService_ExportUserList`
3. 在测试中包含成功和失败的场景
4. 使用 mock 对象来隔离依赖
5. 确保测试覆盖所有关键的业务逻辑
6. 定期运行测试并保持测试用例的更新

## 常见问题

### 1. 测试无法运行

确保：
- 已安装所有依赖 `go mod tidy`
- 测试文件名以 `_test.go` 结尾
- 测试函数名以 `Test` 开头

### 2. Mock 对象方法未被调用

检查：
- mock 对象的设置是否正确
- 参数类型是否匹配
- 是否使用了正确的 mock 方法名

### 3. 测试覆盖率低

考虑：
- 添加更多的测试场景
- 测试边界条件
- 测试错误处理路径 