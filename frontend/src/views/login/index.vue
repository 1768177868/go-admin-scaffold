<template>
  <div class="login-container">
    <div class="login-form">
      <div class="login-header">
        <h2>Go Admin 后台管理系统</h2>
        <p>欢迎登录</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        size="large"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
            clearable
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            show-password
            clearable
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        
        <el-form-item prop="captcha_code">
          <div class="captcha-container">
            <el-input
              v-model="loginForm.captcha_code"
              placeholder="请输入验证码"
              clearable
              style="flex: 1; margin-right: 10px;"
            />
            <div class="captcha-image" @click="refreshCaptcha">
              <img v-if="captchaImage" :src="captchaImage" alt="验证码" />
              <div v-else class="captcha-loading">加载中...</div>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            style="width: 100%;"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-tips">
        <div class="demo-account">
          <p><strong>演示账号：</strong></p>
          <p>管理员：admin / admin123</p>
          <p>普通用户：user / admin123</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { User, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getCaptcha } from '@/api/auth'

export default {
  name: 'Login',
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()
    const route = useRoute()

    const loginFormRef = ref()
    const loading = ref(false)
    const captchaImage = ref('')
    const captchaId = ref('')
    
    const loginForm = reactive({
      username: 'admin',
      password: 'admin123',
      captcha_code: '',
      captcha_id: ''
    })
    
    const loginRules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ],
      captcha_code: [
        { required: true, message: '请输入验证码', trigger: 'blur' }
      ]
    }

    // 获取验证码
    const refreshCaptcha = async () => {
      try {
        const { data } = await getCaptcha()
        captchaImage.value = data.captcha_image
        captchaId.value = data.captcha_id
        loginForm.captcha_id = data.captcha_id
        loginForm.captcha_code = ''
      } catch (error) {
        console.error('获取验证码失败:', error)
        ElMessage.error('获取验证码失败')
      }
    }

    const handleLogin = async () => {
      if (!loginFormRef.value) return
      
      await loginFormRef.value.validate(async (valid) => {
        if (valid) {
          loading.value = true
          try {
            await authStore.login(loginForm)
            ElMessage.success('登录成功')
            
            // 跳转到目标页面
            const redirect = route.query.redirect || '/'
            router.push(redirect)
          } catch (error) {
            console.error('Login failed:', error)
            // 登录失败后刷新验证码
            refreshCaptcha()
          } finally {
            loading.value = false
          }
        }
      })
    }

    // 组件挂载时获取验证码
    onMounted(() => {
      refreshCaptcha()
    })

    return {
      User,
      Lock,
      loginFormRef,
      loginForm,
      loginRules,
      loading,
      captchaImage,
      handleLogin,
      refreshCaptcha
    }
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
}

.login-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
}

.login-form {
  width: 400px;
  background: white;
  border-radius: 8px;
  padding: 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  color: #333;
  margin-bottom: 10px;
  font-weight: 600;
}

.login-header p {
  color: #666;
  font-size: 14px;
}

.captcha-container {
  display: flex;
  align-items: center;
  width: 100%;
}

.captcha-image {
  width: 120px;
  height: 40px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  transition: border-color 0.3s;
}

.captcha-image:hover {
  border-color: #409eff;
}

.captcha-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 3px;
}

.captcha-loading {
  font-size: 12px;
  color: #909399;
}

.login-tips {
  margin-top: 20px;
  text-align: center;
}

.demo-account {
  background: #f8f9fa;
  border-radius: 4px;
  padding: 15px;
  font-size: 13px;
  color: #666;
  text-align: left;
}

.demo-account p {
  margin: 5px 0;
  line-height: 1.5;
}

:deep(.el-input__inner) {
  height: 50px;
}
</style> 