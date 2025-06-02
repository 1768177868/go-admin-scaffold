import request from '@/utils/request'

const menuApi = {
  // 获取菜单列表
  getMenuList() {
    return request({
      url: '/admin/v1/menus',
      method: 'get'
    })
  },

  // 获取菜单树
  getMenuTree() {
    return request({
      url: '/admin/v1/menus/tree',
      method: 'get'
    })
  },

  // 获取用户菜单
  getUserMenus() {
    return request({
      url: '/admin/v1/menus/user',
      method: 'get'
    })
  },

  // 获取菜单详情
  getMenu(id) {
    return request({
      url: `/admin/v1/menus/${id}`,
      method: 'get'
    })
  },

  // 创建菜单
  createMenu(data) {
    return request({
      url: '/admin/v1/menus',
      method: 'post',
      data
    })
  },

  // 更新菜单
  updateMenu(id, data) {
    return request({
      url: `/admin/v1/menus/${id}`,
      method: 'put',
      data
    })
  },

  // 删除菜单
  deleteMenu(id) {
    return request({
      url: `/admin/v1/menus/${id}`,
      method: 'delete'
    })
  },

  // 更新菜单角色
  updateMenuRoles(id, roleIds) {
    return request({
      url: `/admin/v1/menus/${id}/roles`,
      method: 'put',
      data: { role_ids: roleIds }
    })
  }
}

export default menuApi 