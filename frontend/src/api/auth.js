import request from '@/utils/request'

// 获取验证码
export function getCaptcha() {
  return request({
    url: '/admin/v1/auth/captcha',
    method: 'get'
  })
}

// 登录
export function login(data) {
  return request({
    url: '/admin/v1/auth/login',
    method: 'post',
    data
  })
}

// 刷新token
export function refreshToken() {
  return request({
    url: '/admin/v1/auth/refresh',
    method: 'post'
  })
}

// 获取用户信息
export function getUserInfo() {
  return request({
    url: '/admin/v1/profile',
    method: 'get'
  })
}

// 退出登录
export function logout() {
  return request({
    url: '/admin/v1/auth/logout',
    method: 'post'
  })
} 