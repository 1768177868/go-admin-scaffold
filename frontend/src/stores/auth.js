import { defineStore } from 'pinia'
import { login, logout, getUserInfo } from '@/api/auth'
import { getToken, setToken, removeToken } from '@/utils/auth'
import { asyncRoutes } from '@/router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getToken(),
    userInfo: null,
    roles: [],
    permissions: []
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    userName: (state) => state.userInfo?.username || '',
    userAvatar: (state) => state.userInfo?.avatar || '',
    userRoles: (state) => state.roles || [],
    userPermissions: (state) => state.permissions || []
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
        throw error
      }
    },

    // 获取用户信息
    async getUserInfo() {
      try {
        const { data } = await getUserInfo()
        this.userInfo = data
        this.roles = data.roles || []
        this.permissions = data.permissions || []
        return data
      } catch (error) {
        throw error
      }
    },

    // 登出
    async logout() {
      try {
        await logout()
      } catch (error) {
        console.error('Logout error:', error)
      } finally {
        this.token = ''
        this.userInfo = null
        this.roles = []
        this.permissions = []
        removeToken()
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
      removeToken()
    }
  }
}) 