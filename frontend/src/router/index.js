import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import NProgress from 'nprogress'

// 导入布局组件
import Layout from '@/layout/index.vue'

// 公共路由
export const constantRoutes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', hidden: true }
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '404', hidden: true }
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'Odometer', affix: true }
      }
    ]
  }
]

// 动态路由
export const asyncRoutes = [
  {
    path: '/system',
    component: Layout,
    name: 'System',
    meta: { title: '系统管理', icon: 'Setting' },
    children: [
      {
        path: 'user',
        name: 'User',
        component: () => import('@/views/system/user/index.vue'),
        meta: { title: '用户管理', icon: 'User', permission: 'user:view' }
      },
      {
        path: 'role',
        name: 'Role',
        component: () => import('@/views/system/role/index.vue'),
        meta: { title: '角色管理', icon: 'UserFilled', permission: 'role:view' }
      },
      {
        path: 'permission',
        name: 'Permission',
        component: () => import('@/views/system/permission/index.vue'),
        meta: { title: '权限管理', icon: 'Key', permission: 'permission:view' }
      },
      {
        path: 'menu',
        name: 'Menu',
        component: () => import('@/views/system/menu/index.vue'),
        meta: { title: '菜单管理', icon: 'Menu', permission: 'menu:view' }
      }
    ]
  },
  {
    path: '/log',
    component: Layout,
    name: 'Log',
    meta: { title: '日志管理', icon: 'Document' },
    children: [
      {
        path: 'login',
        name: 'LoginLog',
        component: () => import('@/views/log/login/index.vue'),
        meta: { title: '登录日志', icon: 'Key', permission: 'log:view' }
      },
      {
        path: 'operation',
        name: 'OperationLog',
        component: () => import('@/views/log/operation/index.vue'),
        meta: { title: '操作日志', icon: 'Document', permission: 'log:view' }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    meta: { hidden: true },
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心', icon: 'User' }
      }
    ]
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes: [...constantRoutes, ...asyncRoutes] // 暂时将所有路由都添加为静态路由
})

// 简化的路由守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()

  const authStore = useAuthStore()
  const token = authStore.token

  // 白名单路由
  const whiteList = ['/login', '/404']

  if (token) {
    if (to.path === '/login') {
      next({ path: '/' })
      NProgress.done()
    } else {
      // 如果没有用户信息，获取用户信息
      if (!authStore.userInfo?.id) {
        try {
          await authStore.getUserInfo()
          next()
        } catch (error) {
          console.error('获取用户信息失败:', error)
          authStore.resetState()
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      } else {
        next()
      }
    }
  } else {
    // 没有token
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})

export default router 