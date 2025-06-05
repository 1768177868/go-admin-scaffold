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
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
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

    // 使用后端返回的菜单数据
    const routes = computed(() => {
      return authStore.userMenus
    })

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