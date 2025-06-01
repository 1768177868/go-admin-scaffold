import request from '@/utils/request'

// 获取登录日志列表
export function getLoginLogs(params) {
  return request({
    url: '/admin/v1/logs/login',
    method: 'get',
    params
  })
}

// 获取操作日志列表
export function getOperationLogs(params) {
  return request({
    url: '/admin/v1/logs/operation',
    method: 'get',
    params
  })
}

// 获取用户日志
export function getUserLogs(userId) {
  return request({
    url: `/admin/v1/logs/user/${userId}`,
    method: 'get'
  })
}

// 清理日志
export function clearLogs(type, days) {
  return request({
    url: '/admin/v1/logs/clear',
    method: 'delete',
    data: { type, days }
  })
} 