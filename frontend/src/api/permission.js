import request from '@/utils/request'

// 获取权限列表
export function getPermissionList(params) {
  return request({
    url: '/admin/v1/permissions',
    method: 'get',
    params
  })
}

// 获取权限树
export function getPermissionTree() {
  return request({
    url: '/admin/v1/permissions/tree',
    method: 'get'
  })
}

// 获取权限详情
export function getPermissionDetail(id) {
  return request({
    url: `/admin/v1/permissions/${id}`,
    method: 'get'
  })
}

// 创建权限
export function createPermission(data) {
  return request({
    url: '/admin/v1/permissions',
    method: 'post',
    data
  })
}

// 更新权限
export function updatePermission(id, data) {
  return request({
    url: `/admin/v1/permissions/${id}`,
    method: 'put',
    data
  })
}

// 删除权限
export function deletePermission(id) {
  return request({
    url: `/admin/v1/permissions/${id}`,
    method: 'delete'
  })
}

// 更新角色权限
export function updateRolePermissions(roleId, data) {
  return request({
    url: `/admin/v1/roles/${roleId}/permissions`,
    method: 'put',
    data
  })
} 