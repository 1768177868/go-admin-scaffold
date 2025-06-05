import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken } from '@/utils/auth'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

// 创建axios实例
const service = axios.create({
  baseURL: '/api', // api的base_url
  timeout: 10000 // 请求超时时间
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    // 在请求发送之前做一些处理
    const token = getToken()
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data

    // 如果自定义代码不是0，则判断为一个错误
    if (res.code !== 0) {
      ElMessage({
        message: res.message || 'Error',
        type: 'error',
        duration: 5 * 1000
      })

      // 401: 未授权，需要重新登录
      if (res.code === 401 || res.code === 10401 || res.code === 10005) {
        handleUnauthorized()
        return Promise.reject(new Error(res.message || 'Unauthorized'))
      }
      return Promise.reject(new Error(res.message || 'Error'))
    } else {
      return res
    }
  },
  error => {
    console.error('Response error:', error)
    
    // 处理 401 错误或自定义未授权错误码
    if (
      error.response && (
        error.response.status === 401 || 
        error.response.data?.code === 401 ||
        error.response.data?.code === 10005
      )
    ) {
      handleUnauthorized()
    } else {
      ElMessage({
        message: error.message || '请求失败',
        type: 'error',
        duration: 5 * 1000
      })
    }
    return Promise.reject(error)
  }
)

// 处理未授权的情况
function handleUnauthorized() {
  const authStore = useAuthStore()
  
  // 先重置状态
  authStore.resetState()
  
  // 强制重定向到登录页，不保留当前路径
  router.replace('/login')
  
  ElMessage({
    message: '登录已过期，请重新登录',
    type: 'warning'
  })
}

export default service 