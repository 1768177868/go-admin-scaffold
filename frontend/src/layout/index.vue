<template>
  <div class="layout-container">
    <el-container direction="horizontal">
      <!-- 侧边栏 -->
      <el-aside width="240px" class="layout-aside">
        <Sidebar />
      </el-aside>
      
      <!-- 主体区域 -->
      <el-container direction="vertical" class="layout-right">
        <!-- 顶部导航 -->
        <el-header height="60px" class="layout-header">
          <Navbar />
        </el-header>
        
        <!-- 标签页导航 -->
        <div class="tags-view-wrapper">
          <TagsView />
        </div>
        
        <!-- 主内容区 -->
        <el-main class="layout-main">
          <router-view v-slot="{ Component }">
            <transition name="fade-transform" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
import Sidebar from './components/Sidebar.vue'
import Navbar from './components/Navbar.vue'
import TagsView from './components/TagsView.vue'

export default {
  name: 'Layout',
  components: {
    Sidebar,
    Navbar,
    TagsView
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.layout-container .el-container {
  height: 100vh;
}

.layout-aside {
  background: #304156;
  overflow: hidden;
  height: 100vh;
}

.layout-right {
  height: 100vh;
  overflow: hidden;
  flex: 1;
}

.layout-header {
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  display: flex;
  align-items: center;
  padding: 0;
  z-index: 1000;
}

.tags-view-wrapper {
  height: 34px;
  flex-shrink: 0;
}

.layout-main {
  background: #f5f5f5;
  overflow-y: auto;
  flex: 1;
  padding: 20px;
  margin: 0;
}

/* 页面切换动画 */
.fade-transform-leave-active,
.fade-transform-enter-active {
  transition: all 0.3s;
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style> 