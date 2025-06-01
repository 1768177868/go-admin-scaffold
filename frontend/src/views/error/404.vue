<template>
  <div class="error-404">
    <div class="error-content">
      <div class="error-image">
        <img src="https://via.placeholder.com/300x200/f0f0f0/666666?text=404" alt="404">
      </div>
      <div class="error-text">
        <h1>404</h1>
        <h2>页面不存在</h2>
        <p>抱歉，您访问的页面不存在或已被删除</p>
        <div class="error-actions">
          <el-button type="primary" @click="goHome">
            <el-icon><House /></el-icon>
            返回首页
          </el-button>
          <el-button @click="goBack">
            <el-icon><Back /></el-icon>
            返回上页
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { House, Back } from '@element-plus/icons-vue'

export default {
  name: 'Error404',
  components: {
    House,
    Back
  },
  setup() {
    const router = useRouter()

    const goHome = () => {
      router.push('/')
    }

    const goBack = () => {
      router.back()
    }

    return {
      goHome,
      goBack
    }
  }
}
</script>

<style lang="scss" scoped>
.error-404 {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.error-content {
  text-align: center;
  background: #fff;
  border-radius: 20px;
  padding: 60px 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  max-width: 500px;
}

.error-image {
  margin-bottom: 30px;
  
  img {
    max-width: 100%;
    height: auto;
  }
}

.error-text {
  h1 {
    font-size: 72px;
    font-weight: 700;
    color: #409eff;
    margin-bottom: 10px;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  h2 {
    font-size: 24px;
    color: #303133;
    margin-bottom: 15px;
  }
  
  p {
    font-size: 16px;
    color: #606266;
    margin-bottom: 30px;
    line-height: 1.6;
  }
}

.error-actions {
  .el-button {
    margin: 0 10px;
  }
}
</style> 