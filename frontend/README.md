# Go Admin Frontend

基于 Vue3 + Element Plus 的后台管理系统前端

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Element Plus** - Vue 3 组件库
- **Vue Router** - Vue.js 官方路由管理器
- **Pinia** - Vue 状态管理库
- **Axios** - HTTP 客户端
- **Vite** - 下一代前端构建工具
- **SCSS** - CSS 预处理器

## 功能特性

### 🔐 用户认证
- 用户登录/登出
- JWT Token 认证
- 自动 Token 刷新
- 权限控制

### 👥 用户管理
- 用户 CRUD 操作
- 用户状态管理
- 角色分配
- 密码重置

### 🛡️ 角色权限
- 角色管理
- 权限管理
- 基于角色的访问控制 (RBAC)
- 动态菜单生成

### 📊 日志监控
- 登录日志
- 操作日志
- 日志搜索和过滤
- 日志清理

### 💼 个人中心
- 个人信息管理
- 头像上传
- 密码修改
- 登录记录

### 🎨 界面特色
- 响应式设计
- 现代化 UI
- 暗色主题支持
- 多标签页导航
- 面包屑导航

## 项目结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API 接口
│   │   ├── auth.js        # 认证相关
│   │   ├── user.js        # 用户管理
│   │   ├── role.js        # 角色管理
│   │   ├── permission.js  # 权限管理
│   │   └── log.js         # 日志管理
│   ├── assets/            # 静态资源
│   │   └── styles/        # 样式文件
│   ├── components/        # 公共组件
│   │   └── Pagination/    # 分页组件
│   ├── layout/           # 布局组件
│   │   ├── components/   # 布局子组件
│   │   └── index.vue     # 主布局
│   ├── router/           # 路由配置
│   ├── stores/           # 状态管理
│   ├── utils/            # 工具函数
│   │   ├── auth.js       # 认证工具
│   │   ├── request.js    # HTTP 请求封装
│   │   └── scroll-to.js  # 滚动工具
│   ├── views/            # 页面组件
│   │   ├── dashboard/    # 仪表盘
│   │   ├── login/        # 登录页
│   │   ├── system/       # 系统管理
│   │   ├── log/          # 日志管理
│   │   ├── profile/      # 个人中心
│   │   └── error/        # 错误页面
│   ├── App.vue           # 根组件
│   └── main.js           # 入口文件
├── index.html            # HTML 模板
├── package.json          # 项目配置
├── vite.config.js        # Vite 配置
└── README.md            # 项目说明
```

## 开始使用

### 环境要求

- Node.js >= 16.0.0
- npm 或 yarn

### 安装依赖

```bash
# 使用 npm
npm install

# 或使用 yarn
yarn install
```

### 开发环境

```bash
# 启动开发服务器
npm run dev

# 或
yarn dev
```

访问 http://localhost:3001

### 生产构建

```bash
# 构建生产版本
npm run build

# 或
yarn build
```

### 代码检查

```bash
# 运行 ESLint
npm run lint

# 或
yarn lint
```

## 配置说明

### API 代理配置

在 `vite.config.js` 中配置了 API 代理：

```javascript
server: {
  port: 3001,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      secure: false
    }
  }
}
```

### 环境变量

创建 `.env.local` 文件来配置环境变量：

```bash
# API 基础 URL
VITE_API_BASE_URL=http://localhost:8080/api

# 应用标题
VITE_APP_TITLE=Go Admin 后台管理系统
```

## 登录凭据

### 演示账号

- **管理员**
  - 用户名: `admin`
  - 密码: `admin123`

- **普通用户**
  - 用户名: `user`
  - 密码: `admin123`

## 部署说明

### 生产环境部署

1. 构建项目：
```bash
npm run build
```

2. 将 `dist` 目录下的文件部署到 Web 服务器

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        root /path/to/dist;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 开发指南

### 添加新页面

1. 在 `src/views` 下创建页面组件
2. 在 `src/router/index.js` 中添加路由配置
3. 在 `src/api` 中添加相关 API 接口
4. 如需权限控制，在路由 meta 中添加 permission 字段

### 添加新组件

1. 在 `src/components` 下创建组件
2. 使用 PascalCase 命名
3. 导出为默认模块

### 状态管理

使用 Pinia 进行状态管理，stores 位于 `src/stores` 目录：

```javascript
import { defineStore } from 'pinia'

export const useExampleStore = defineStore('example', {
  state: () => ({
    // 状态定义
  }),
  getters: {
    // 计算属性
  },
  actions: {
    // 方法定义
  }
})
```

## 常见问题

### Q: 登录后页面空白？
A: 检查后端 API 是否正常运行，确认 JWT Token 格式正确。

### Q: 权限控制不生效？
A: 确认用户角色和权限配置正确，检查路由 meta.permission 字段。

### Q: 构建失败？
A: 检查 Node.js 版本是否 >= 16.0.0，清除 node_modules 重新安装依赖。

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88

## 开源协议

MIT License

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 完整的用户管理功能
- 角色权限系统
- 日志监控
- 个人中心

## 联系方式

如有问题或建议，请创建 Issue 或联系开发团队。 