<template>
  <div class="navbar">
    <!-- 左侧 -->
    <div class="navbar-left">
      <!-- 面包屑 -->
      <el-breadcrumb class="app-breadcrumb" separator="/">
        <el-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
          <span v-if="item.redirect === 'noRedirect' || !item.path" class="no-redirect">
            {{ item.meta.title }}
          </span>
          <router-link v-else :to="item.path">
            {{ item.meta.title }}
          </router-link>
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 右侧 -->
    <div class="navbar-right">
      <!-- 语言选择器 -->
      <LangSelect class="lang-select-item" />
      
      <!-- 头像下拉菜单 -->
      <el-dropdown class="avatar-container" trigger="click">
        <div class="avatar-wrapper">
          <img :src="avatar" class="user-avatar">
          <el-icon class="el-icon-caret-bottom">
            <CaretBottom />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <router-link to="/profile">
              <el-dropdown-item>个人中心</el-dropdown-item>
            </router-link>
            <el-dropdown-item divided @click="logout">
              <span style="display:block;">退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { CaretBottom } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import LangSelect from '@/components/LangSelect/index.vue'

export default {
  name: 'Navbar',
  components: {
    CaretBottom,
    LangSelect
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const authStore = useAuthStore()

    // 面包屑
    const breadcrumbs = ref([])

    // 用户头像
    const avatar = computed(() => {
      return authStore.userAvatar || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
    })

    // 生成面包屑
    function getBreadcrumb() {
      let matched = route.matched.filter(item => item.meta && item.meta.title)
      const first = matched[0]

      if (!isDashboard(first)) {
        matched = [{ path: '/dashboard', meta: { title: '仪表盘' } }].concat(matched)
      }

      breadcrumbs.value = matched.filter(item => item.meta && item.meta.title && item.meta.breadcrumb !== false)
    }

    function isDashboard(route) {
      const name = route && route.name
      if (!name) {
        return false
      }
      return name.trim().toLocaleLowerCase() === 'Dashboard'.toLocaleLowerCase()
    }

    // 退出登录
    async function logout() {
      await ElMessageBox.confirm('确定注销并退出系统吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })

      await authStore.logout()
      router.push(`/login?redirect=${route.fullPath}`)
    }

    // 监听路由变化
    watch(route, getBreadcrumb, { immediate: true })

    return {
      breadcrumbs,
      avatar,
      logout
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
}

.navbar-left {
  flex: 1;
}

.navbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.app-breadcrumb {
  display: inline-block;
  font-size: 14px;
  line-height: 50px;
  margin-left: 8px;

  .no-redirect {
    color: #97a8be;
    cursor: text;
  }
}

.avatar-container {
  margin-right: 30px;

  .avatar-wrapper {
    margin-top: 5px;
    position: relative;
    display: flex;
    align-items: center;
    cursor: pointer;

    .user-avatar {
      cursor: pointer;
      width: 40px;
      height: 40px;
      border-radius: 50%;
    }

    .el-icon-caret-bottom {
      cursor: pointer;
      position: absolute;
      right: -20px;
      top: 25px;
      font-size: 12px;
    }
  }
}

.lang-select-item {
  margin-right: 8px;
}
</style> 