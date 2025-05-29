# Go Admin Framework Documentation

## 目录

### 入门指南
- [快速开始](getting-started/quick-start.md)
- [项目结构](getting-started/structure.md)
- [配置说明](getting-started/configuration.md)

### 功能特性
- [用户认证](features/authentication.md)
- [角色权限](features/rbac.md)
- [缓存系统](features/cache.md)
- [任务调度](features/scheduling.md)

### API 文档
- [API 概述](api/overview.md)
- [API 参考](api/README.md)

### 测试指南
- [测试指南](testing.md)

## 核心功能

### 认证与授权
- [认证系统](features/authentication.md)
  - JWT Token 认证
  - 刷新令牌机制
  - 多端登录控制
  - 登录日志记录

- [权限控制](features/rbac.md)
  - 基于角色的访问控制
  - 权限管理
  - 动态权限分配
  - 权限缓存机制

### 系统功能
- [缓存与队列](features/cache.md)
  - Redis 缓存管理
  - 消息队列系统
  - 分布式锁
  - 数据同步

- [任务调度](features/scheduling.md)
  - 定时任务管理
  - 异步任务处理
  - 任务监控
  - 失败重试机制

### API 开发
- [API 设计](api/overview.md)
  - RESTful API 规范
  - 请求/响应格式
  - 错误处理
  - 参数验证

### 项目开发
- [项目配置](getting-started/configuration.md)
  - 环境配置
  - 数据库配置
  - 缓存配置
  - 日志配置

- [测试规范](testing.md)
  - 单元测试
  - 集成测试
  - 性能测试
  - 测试覆盖率

## 快速链接

- [项目主页](../README.md)
- [API 文档](api/README.md)
- [快速开始](getting-started/quick-start.md)
- [测试指南](testing.md)

## 文档更新

本文档持续更新中。如果发现任何问题或有改进建议，请提交 Issue 或 Pull Request。

## 贡献指南

如果您想为项目文档做出贡献，请遵循以下步骤：

1. Fork 项目仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-doc`)
3. 提交您的改动 (`git commit -m 'Add some amazing doc'`)
4. 推送到分支 (`git push origin feature/amazing-doc`)
5. 创建一个 Pull Request 