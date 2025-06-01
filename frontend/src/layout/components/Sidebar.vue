<template>
  <div class="sidebar-container">
    <!-- Logo -->
    <div class="sidebar-logo">
      <router-link to="/" class="sidebar-logo-link">
        <h1 class="sidebar-title">{{ collapse ? 'GA' : 'Go Admin' }}</h1>
      </router-link>
    </div>

    <!-- 菜单 -->
    <el-scrollbar class="sidebar-menu-container">
      <el-menu
        :default-active="activeMenu"
        :collapse="collapse"
        :unique-opened="false"
        :collapse-transition="false"
        mode="vertical"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <sidebar-item
          v-for="route in routes"
          :key="route.path"
          :item="route"
          :base-path="route.path"
        />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script>
import { useRoute } from 'vue-router'
import { constantRoutes, asyncRoutes } from '@/router'
import SidebarItem from './SidebarItem.vue'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'Sidebar',
  components: {
    SidebarItem
  },
  setup() {
    const route = useRoute()
    const authStore = useAuthStore()

    const collapse = ref(false)

    // 计算当前激活的菜单
    const activeMenu = computed(() => {
      const { meta, path } = route
      if (meta.activeMenu) {
        return meta.activeMenu
      }
      return path
    })

    // 计算可见的路由
    const routes = computed(() => {
      const permissions = authStore.userPermissions
      const allRoutes = [...constantRoutes, ...asyncRoutes]
      return filterMenuRoutes(allRoutes, permissions)
    })

    // 过滤菜单路由
    function filterMenuRoutes(routes, permissions) {
      const res = []
      routes.forEach(route => {
        const tmp = { ...route }
        
        // 隐藏的路由不显示在菜单中
        if (tmp.meta?.hidden) {
          return
        }

        // 检查权限
        if (hasPermission(permissions, tmp)) {
          if (tmp.children) {
            tmp.children = filterMenuRoutes(tmp.children, permissions)
          }
          res.push(tmp)
        }
      })
      return res
    }

    // 权限检查
    function hasPermission(permissions, route) {
      if (route.meta && route.meta.permission) {
        return permissions.some(permission => permission === route.meta.permission)
      } else {
        return true
      }
    }

    return {
      collapse,
      activeMenu,
      routes
    }
  }
}
</script>

<style lang="scss" scoped>
.sidebar-container {
  height: 100%;
  background-color: #304156;
  width: 100% !important;
}

.sidebar-logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  border-bottom: 1px solid #2c3e50;
  
  .sidebar-logo-link {
    display: inline-block;
    text-decoration: none;
    
    .sidebar-title {
      display: inline-block;
      margin: 0;
      color: #fff;
      font-weight: 600;
      font-size: 18px;
      vertical-align: middle;
      transition: all 0.3s;
    }
  }
}

.sidebar-menu-container {
  height: calc(100% - 60px);
  
  :deep(.el-menu) {
    border-right: none;
    height: 100%;
    width: 100% !important;
  }
}
</style> 