import request from '@/utils/request'

// 获取用户列表
export function getUserList(params) {
  return request({
    url: '/admin/v1/users',
    method: 'get',
    params
  })
}

// 获取用户列表（别名）
export function listUsers(params) {
  return request({
    url: '/admin/v1/users',
    method: 'get',
    params
  })
}

// 获取用户详情
export function getUserDetail(id) {
  return request({
    url: `/admin/v1/users/${id}`,
    method: 'get'
  })
}

// 创建用户
export function createUser(data) {
  return request({
    url: '/admin/v1/users',
    method: 'post',
    data
  })
}

// 更新用户
export function updateUser(id, data) {
  return request({
    url: `/admin/v1/users/${id}`,
    method: 'put',
    data
  })
}

// 删除用户
export function deleteUser(id) {
  return request({
    url: `/admin/v1/users/${id}`,
    method: 'delete'
  })
}

// 重置用户密码
export function resetPassword(id, data) {
  return request({
    url: `/admin/v1/users/${id}/password`,
    method: 'put',
    data
  })
}

// 更新用户状态
export function updateUserStatus(id, status) {
  console.log('updateUserStatus called with:', { id, status, statusType: typeof status })
  const data = { status }
  console.log('Sending data:', data)
  
  return request({
    url: `/admin/v1/users/${id}/status`,
    method: 'put',
    data
  })
}

// 更新用户角色
export function updateUserRoles(id, roleIds) {
  return request({
    url: `/admin/v1/users/${id}/roles`,
    method: 'put',
    data: { role_ids: roleIds }
  })
} 