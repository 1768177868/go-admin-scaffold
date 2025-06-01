# Go Admin 前端启动指南

## 快速开始

### 1. 安装 Node.js
确保安装了 Node.js 16.0.0 或更高版本：
```bash
node --version
npm --version
```

### 2. 进入前端目录
```bash
cd frontend
```

### 3. 安装依赖
```bash
npm install
```

### 4. 启动开发服务器
```bash
npm run dev
```

### 5. 访问应用
打开浏览器访问：http://localhost:3001

## 登录信息

### 管理员账号
- 用户名：`admin`
- 密码：`admin123`

### 普通用户账号
- 用户名：`user`
- 密码：`admin123`

## 项目结构说明

```
frontend/
├── src/
│   ├── api/           # API 接口定义
│   ├── assets/        # 静态资源
│   ├── components/    # 公共组件
│   ├── layout/        # 布局组件
│   ├── router/        # 路由配置
│   ├── stores/        # 状态管理 (Pinia)
│   ├── utils/         # 工具函数
│   └── views/         # 页面组件
│       ├── dashboard/ # 仪表盘
│       ├── system/    # 系统管理
│       │   ├── user/  # 用户管理
│       │   ├── role/  # 角色管理
│       │   └── permission/ # 权限管理
│       ├── log/       # 日志管理
│       │   ├── login/ # 登录日志
│       │   └── operation/ # 操作日志
│       ├── profile/   # 个人中心
│       └── login/     # 登录页面
```

## 功能模块

### ✅ 已完成功能
- [x] 用户认证（登录/登出）
- [x] JWT Token 管理
- [x] 权限控制和路由守卫
- [x] 用户管理（增删改查）
- [x] 角色管理
- [x] 权限管理
- [x] 登录日志查看
- [x] 操作日志查看
- [x] 个人中心
- [x] 响应式布局
- [x] 多标签页导航

### 🎨 UI 特性
- 基于 Element Plus 组件库
- 现代化设计风格
- 响应式布局支持
- 动画过渡效果
- 统一的样式规范

### 🔧 技术特性
- Vue 3 Composition API
- TypeScript 支持（可选）
- Vite 构建工具
- 热模块替换 (HMR)
- ESLint 代码检查
- Axios 请求拦截器
- 路由懒加载

## 开发命令

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# 代码检查
npm run lint
```

## 环境配置

### 开发环境
- 默认端口：3001
- API 代理：localhost:8080

### 生产环境
需要配置实际的 API 地址，可以通过环境变量设置：

```bash
# .env.production
VITE_API_BASE_URL=https://your-api-domain.com/api
```

## 常见问题解决

### 1. 依赖安装失败
```bash
# 清除缓存重新安装
rm -rf node_modules package-lock.json
npm install
```

### 2. 端口被占用
```bash
# 指定其他端口启动
npm run dev -- --port 3001
```

### 3. API 请求失败
检查后端服务是否在 localhost:8080 运行，或修改 vite.config.js 中的代理配置。

### 4. 登录后页面空白
通常是 Token 验证失败，检查：
- 后端是否正常运行
- API 响应格式是否正确
- 浏览器控制台是否有错误信息

## 部署说明

### 1. 构建项目
```bash
npm run build
```

### 2. 部署 dist 目录
将构建生成的 `dist` 目录部署到 Web 服务器。

### 3. Nginx 配置示例
```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://backend-server:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 下一步开发建议

1. **完善用户体验**
   - 添加更多的加载状态
   - 优化错误处理和提示
   - 添加操作确认对话框

2. **扩展功能模块**
   - 文件上传管理
   - 系统设置
   - 数据导入导出

3. **性能优化**
   - 组件懒加载
   - 图片懒加载
   - 缓存优化

4. **测试覆盖**
   - 单元测试
   - 集成测试
   - E2E 测试

## 技术支持

如遇到问题，请：
1. 查看浏览器控制台错误信息
2. 检查网络请求状态
3. 确认后端 API 正常运行
4. 参考项目文档和代码注释 