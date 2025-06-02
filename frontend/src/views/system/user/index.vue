<template>
  <div class="app-container">
    <!-- 搜索栏 -->
    <div class="filter-container">
      <el-input
        v-model="listQuery.username"
        placeholder="请输入用户名"
        style="width: 150px;"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-input
        v-model="listQuery.email"
        placeholder="请输入邮箱"
        style="width: 150px;"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-select v-model="listQuery.status" placeholder="状态" clearable style="width: 120px" class="filter-item">
        <el-option label="启用" :value="1" />
        <el-option label="禁用" :value="0" />
      </el-select>
      <el-select v-model="listQuery.role_id" placeholder="角色" clearable style="width: 120px" class="filter-item">
        <el-option
          v-for="role in roleList"
          :key="role.id"
          :label="role.name"
          :value="role.id"
        />
      </el-select>
      <el-button class="filter-item" type="primary" icon="Search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" type="default" icon="Refresh" @click="resetFilter">
        重置
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="Plus" @click="handleCreate">
        添加用户
      </el-button>
    </div>

    <!-- 表格 -->
    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column align="center" label="ID" width="95">
        <template #default="{ row }">
          {{ row.id }}
        </template>
      </el-table-column>
      
      <el-table-column label="用户名" prop="username" />
      
      <el-table-column label="邮箱" prop="email" />
      
      <el-table-column label="昵称" prop="nickname" />
      
      <el-table-column label="手机号" width="120">
        <template #default="{ row }">
          {{ row.phone }}
        </template>
      </el-table-column>
      
      <el-table-column label="角色" width="120">
        <template #default="{ row }">
          <el-tag v-for="role in row.roles" :key="role.id" size="small" style="margin-right: 5px;">
            {{ role.name }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="状态" prop="status">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column align="center" prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">
          <span>{{ formatDate(row.created_at) }}</span>
        </template>
      </el-table-column>
      
      <el-table-column align="center" label="操作" width="250">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button
            v-if="row.status === 1"
            size="small"
            type="warning"
            @click="handleModifyStatus(row, 0)"
          >
            禁用
          </el-button>
          <el-button
            v-else
            size="small"
            type="success"
            @click="handleModifyStatus(row, 1)"
          >
            启用
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <pagination
      v-show="total > 0"
      :total="total"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.limit"
      @pagination="getList"
    />

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogType === 'create' ? '添加用户' : '编辑用户'"
      v-model="dialogFormVisible"
      width="600px"
    >
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="100px"
        style="width: 400px; margin-left:50px;"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="temp.username" :disabled="dialogType === 'update'" />
        </el-form-item>
        
        <el-form-item v-if="dialogType === 'create'" label="密码" prop="password">
          <el-input v-model="temp.password" type="password" show-password />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="temp.email" />
        </el-form-item>
        
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="temp.nickname" />
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="temp.phone" />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-select v-model="temp.status" placeholder="请选择">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="角色" prop="role_ids">
          <el-select
            v-model="temp.role_ids"
            multiple
            placeholder="请选择角色"
            style="width: 100%"
          >
            <el-option
              v-for="role in roleList"
              :key="role.id"
              :label="role.name"
              :value="role.id"
              :disabled="role.code === 'admin' && dialogType === 'create'"
            />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogFormVisible = false">
            取消
          </el-button>
          <el-button type="primary" @click="dialogType === 'create' ? createData() : updateData()">
            确认
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { getUserList, createUser, updateUser, deleteUser, updateUserStatus, updateUserRoles } from '@/api/user'
import { getRoleList } from '@/api/role'
import Pagination from '@/components/Pagination/index.vue'
import dayjs from 'dayjs'

export default {
  name: 'UserManagement',
  components: {
    Pagination
  },
  setup() {
    const dataForm = ref()
    
    const list = ref([])
    const total = ref(0)
    const listLoading = ref(true)
    const dialogFormVisible = ref(false)
    const dialogType = ref('')
    const roleList = ref([])
    
    const listQuery = reactive({
      page: 1,
      limit: 20,
      username: '',
      email: '',
      status: '',
      role_id: ''
    })
    
    const temp = reactive({
      id: undefined,
      username: '',
      email: '',
      nickname: '',
      password: '',
      phone: '',
      status: 1,
      role_ids: []
    })
    
    const rules = {
      username: [{ required: true, message: '用户名不能为空', trigger: 'blur' }],
      password: [{ required: true, message: '密码不能为空', trigger: 'blur' }],
      email: [
        { required: true, message: '邮箱不能为空', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
      ],
      role_ids: [{ required: true, message: '请至少选择一个角色', trigger: 'change' }]
    }

    const getList = async () => {
      listLoading.value = true
      try {
        const params = {
          page: listQuery.page,
          page_size: listQuery.limit
        }
        
        // 添加搜索参数
        if (listQuery.username) params.username = listQuery.username
        if (listQuery.email) params.email = listQuery.email
        if (listQuery.status !== '') params.status = listQuery.status
        if (listQuery.role_id) params.role_id = listQuery.role_id
        
        const { data } = await getUserList(params)
        list.value = data.items || []
        total.value = data.pagination?.total || 0
      } catch (error) {
        console.error('获取用户列表失败:', error)
        ElMessage.error('获取用户列表失败')
      } finally {
        listLoading.value = false
      }
    }

    const getRoles = async () => {
      try {
        const { data } = await getRoleList({ page: 1, page_size: 100 })
        roleList.value = data.items || []
      } catch (error) {
        console.error('获取角色列表失败:', error)
      }
    }

    const handleFilter = () => {
      listQuery.page = 1
      getList()
    }

    const resetFilter = () => {
      listQuery.username = ''
      listQuery.email = ''
      listQuery.status = ''
      listQuery.role_id = ''
      listQuery.page = 1
      getList()
    }

    const resetTemp = () => {
      temp.id = undefined
      temp.username = ''
      temp.password = ''
      temp.email = ''
      temp.nickname = ''
      temp.phone = ''
      temp.status = 1
      temp.role_ids = []
    }

    const handleCreate = () => {
      resetTemp()
      dialogType.value = 'create'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value.clearValidate()
      })
    }

    const createData = () => {
      dataForm.value.validate(async (valid) => {
        if (valid) {
          try {
            // 创建用户
            const { data } = await createUser(temp)
            // 分配角色
            if (temp.role_ids.length > 0) {
              await updateUserRoles(data.id, temp.role_ids)
            }
            dialogFormVisible.value = false
            ElMessage.success('创建成功')
            getList()
          } catch (error) {
            console.error('创建用户失败:', error)
            ElMessage.error('创建用户失败')
          }
        }
      })
    }

    const handleUpdate = (row) => {
      Object.assign(temp, row)
      temp.role_ids = row.roles?.map(role => role.id) || []
      dialogType.value = 'update'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value.clearValidate()
      })
    }

    const updateData = () => {
      dataForm.value.validate(async (valid) => {
        if (valid) {
          try {
            const updateData = { ...temp }
            delete updateData.password // 移除密码字段
            // 更新用户基本信息
            await updateUser(temp.id, updateData)
            // 更新用户角色
            await updateUserRoles(temp.id, temp.role_ids)
            dialogFormVisible.value = false
            ElMessage.success('更新成功')
            getList()
          } catch (error) {
            console.error('更新用户失败:', error)
            ElMessage.error('更新用户失败')
          }
        }
      })
    }

    const handleDelete = (row) => {
      ElMessageBox.confirm('此操作将永久删除该用户, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await deleteUser(row.id)
          const index = list.value.findIndex(v => v.id === row.id)
          list.value.splice(index, 1)
          ElMessage.success('删除成功')
          getList()
        } catch (error) {
          console.error('删除用户失败:', error)
          ElMessage.error('删除用户失败')
        }
      })
    }

    const handleModifyStatus = async (row, status) => {
      try {
        await updateUserStatus(row.id, status)
        row.status = status
        ElMessage.success('状态更新成功')
        getList()
      } catch (error) {
        console.error('更新状态失败:', error)
        ElMessage.error('更新状态失败')
      }
    }

    const formatDate = (date) => {
      if (!date) return ''
      return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    }

    onMounted(() => {
      getList()
      getRoles()
    })

    return {
      list,
      total,
      listLoading,
      listQuery,
      dialogFormVisible,
      dialogType,
      temp,
      rules,
      dataForm,
      roleList,
      getList,
      handleFilter,
      handleCreate,
      createData,
      handleUpdate,
      updateData,
      handleDelete,
      handleModifyStatus,
      formatDate,
      resetFilter
    }
  }
}
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.filter-container {
  padding-bottom: 10px;
  
  .filter-item {
    display: inline-block;
    vertical-align: middle;
    margin-bottom: 10px;
    margin-right: 10px;
  }
}
</style> 