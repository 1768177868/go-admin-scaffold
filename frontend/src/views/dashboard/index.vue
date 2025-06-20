<template>
  <div class="dashboard">
    <div class="dashboard-header">
      <h1>Go Admin 后台管理系统</h1>
      <p>欢迎使用后台管理系统</p>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stats-card">
          <div class="stats-icon user">
            <el-icon><User /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-number">{{ stats.users }}</div>
            <div class="stats-label">用户总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stats-card">
          <div class="stats-icon role">
            <el-icon><UserFilled /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-number">{{ stats.roles }}</div>
            <div class="stats-label">角色总数</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stats-card">
          <div class="stats-icon login">
            <el-icon><Key /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-number">{{ stats.logins }}</div>
            <div class="stats-label">今日登录</div>
          </div>
        </div>
      </el-col>
      
      <el-col :xs="24" :sm="12" :lg="6">
        <div class="stats-card">
          <div class="stats-icon operation">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stats-content">
            <div class="stats-number">{{ stats.operations }}</div>
            <div class="stats-label">今日操作</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 快捷入口 -->
    <div class="quick-actions">
      <h2>快捷操作</h2>
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="action-card" @click="$router.push('/system/user')">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="action-card" @click="$router.push('/system/role')">
            <el-icon><UserFilled /></el-icon>
            <span>角色管理</span>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="action-card" @click="$router.push('/log/login')">
            <el-icon><Key /></el-icon>
            <span>登录日志</span>
          </div>
        </el-col>
        
        <el-col :xs="24" :sm="12" :lg="6">
          <div class="action-card" @click="$router.push('/log/operation')">
            <el-icon><Document /></el-icon>
            <span>操作日志</span>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 系统信息 -->
    <div class="system-info">
      <h2>系统信息</h2>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-card>
            <template #header>
              <span>服务器信息</span>
            </template>
            <div class="info-item">
              <span>操作系统:</span>
              <span>{{ systemInfo.os }}</span>
            </div>
            <div class="info-item">
              <span>Go版本:</span>
              <span>{{ systemInfo.goVersion }}</span>
            </div>
            <div class="info-item">
              <span>运行时间:</span>
              <span>{{ systemInfo.uptime }}</span>
            </div>
            <div class="info-item">
              <span>内存使用:</span>
              <span>{{ systemInfo.memory }}</span>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="12">
          <el-card>
            <template #header>
              <span>项目信息</span>
            </template>
            <div class="info-item">
              <span>项目名称:</span>
              <span>Go Admin</span>
            </div>
            <div class="info-item">
              <span>项目版本:</span>
              <span>v1.0.0</span>
            </div>
            <div class="info-item">
              <span>技术栈:</span>
              <span>Go + Gin + Vue3 + Element Plus</span>
            </div>
            <div class="info-item">
              <span>更新时间:</span>
              <span>{{ new Date().toLocaleDateString() }}</span>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { User, UserFilled, Key, Document } from '@element-plus/icons-vue'

export default {
  name: 'Dashboard',
  components: {
    User,
    UserFilled,
    Key,
    Document
  },
  setup() {
    const stats = reactive({
      users: 0,
      roles: 0,
      logins: 0,
      operations: 0
    })

    const systemInfo = reactive({
      os: 'Linux',
      goVersion: 'Go 1.21',
      uptime: '15天 6小时',
      memory: '256MB / 2GB'
    })

    // 模拟加载统计数据
    const loadStats = () => {
      // 这里应该调用实际的API
      setTimeout(() => {
        stats.users = 128
        stats.roles = 8
        stats.logins = 45
        stats.operations = 156
      }, 1000)
    }

    onMounted(() => {
      loadStats()
    })

    return {
      stats,
      systemInfo
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  padding: 20px;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 30px;
  
  h1 {
    font-size: 28px;
    color: #303133;
    margin-bottom: 10px;
  }
  
  p {
    font-size: 16px;
    color: #909399;
  }
}

.stats-row {
  margin-bottom: 30px;
}

.stats-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
  
  &:hover {
    transform: translateY(-5px);
  }
}

.stats-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20px;
  
  .el-icon {
    font-size: 24px;
    color: #fff;
  }
  
  &.user {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
  
  &.role {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  }
  
  &.login {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  }
  
  &.operation {
    background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
  }
}

.stats-content {
  flex: 1;
}

.stats-number {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 5px;
}

.stats-label {
  font-size: 14px;
  color: #909399;
}

.quick-actions, .system-info {
  margin-bottom: 30px;
  
  h2 {
    font-size: 20px;
    color: #303133;
    margin-bottom: 20px;
  }
}

.action-card {
  background: #fff;
  border-radius: 8px;
  padding: 30px 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 25px 0 rgba(0, 0, 0, 0.15);
  }
  
  .el-icon {
    font-size: 32px;
    color: #409eff;
    margin-bottom: 10px;
  }
  
  span {
    display: block;
    font-size: 16px;
    color: #303133;
  }
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  
  &:last-child {
    margin-bottom: 0;
  }
}

:deep(.el-card__header) {
  background: #f8f9fa;
  border-bottom: 1px solid #ebeef5;
}
</style> 