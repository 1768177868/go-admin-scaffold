import { defineStore } from 'pinia'
import { login, logout, getUserInfo } from '@/api/auth'
import menuApi from '@/api/menu'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { asyncRoutes } from '@/router'
import { ElMessage } from 'element-plus'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getToken(),
    userInfo: null,
    roles: [],
    permissions: [],
    menus: [], // 存储用户菜单
    menuError: null // 存储菜单加载错误
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    userName: (state) => state.userInfo?.username || '',
    userAvatar: (state) => state.userInfo?.avatar || '',
    userRoles: (state) => state.roles || [],
    userPermissions: (state) => state.permissions || [],
    userMenus: (state) => state.menus || [] // 用户菜单的getter
  },

  actions: {
    // 登录
    async login(loginForm) {
      try {
        const { data } = await login(loginForm)
        this.token = data.access_token
        setToken(data.access_token)
        return data
      } catch (error) {
        this.resetState()
        throw error
      }
    },

    // 登出
    async logout() {
      try {
        await logout()
      } finally {
        this.resetState()
      }
    },

    // 获取用户信息
    async getUserInfo() {
      try {
        const { data } = await getUserInfo()
        this.userInfo = data.user || data  // Support both nested and flat structure
        this.roles = data.user?.roles || data.roles || []
        this.permissions = data.permissions || []
        
        // 尝试获取用户菜单，但不要因为菜单失败而导致整个getUserInfo失败
        try {
          await this.getUserMenus()
        } catch (menuError) {
          console.warn('获取菜单失败，但用户信息获取成功:', menuError)
          // 菜单获取失败不应该影响用户登录状态
        }
        
        return data
      } catch (error) {
        // 只有在获取用户信息本身失败时才清除token
        console.error('获取用户信息失败:', error)
        this.resetState()
        throw error
      }
    },

    // 获取用户菜单
    async getUserMenus() {
      try {
        const { data } = await menuApi.getUserMenus()
        this.menus = data || []
        this.menuError = null
        console.log('菜单获取成功:', data)
        return data
      } catch (error) {
        console.error('获取用户菜单失败:', error)
        this.menuError = error
        
        // 只有在严重的授权错误时才重置状态
        if (error.response?.status === 401 && error.response?.data?.code === 10005) {
          console.warn('Token已过期，需要重新登录')
          this.resetState()
          throw error
        }
        
        // 对于其他错误，显示警告但继续运行
        console.warn('菜单加载失败，使用空菜单列表')
        this.menus = []
        
        // 不抛出错误，让应用继续运行
        return []
      }
    },

    // 生成路由
    async generateRoutes() {
      const accessedRoutes = this.filterAsyncRoutes(asyncRoutes, this.permissions)
      
      // 添加404路由
      accessedRoutes.push({ path: '/:pathMatch(.*)*', redirect: '/404', hidden: true })
      
      return accessedRoutes
    },

    // 过滤异步路由
    filterAsyncRoutes(routes, permissions) {
      const res = []
      
      routes.forEach(route => {
        const tmp = { ...route }
        
        if (this.hasPermission(permissions, tmp)) {
          if (tmp.children) {
            tmp.children = this.filterAsyncRoutes(tmp.children, permissions)
          }
          res.push(tmp)
        }
      })
      
      return res
    },

    // 检查权限
    hasPermission(permissions, route) {
      if (route.meta && route.meta.permission) {
        return permissions.some(permission => permission === route.meta.permission)
      } else {
        return true
      }
    },

    // 重置状态
    resetState() {
      this.token = ''
      this.userInfo = null
      this.roles = []
      this.permissions = []
      this.menus = []
      this.menuError = null
      removeToken()
    }
  }
}) 