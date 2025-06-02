<template>
  <div class="app-container">
    <!-- 搜索栏 -->
    <div class="filter-container">
      <el-input
        v-model="listQuery.keyword"
        placeholder="请输入权限名称"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-button class="filter-item" type="primary" icon="Search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="Plus" @click="handleCreate">
        添加权限
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
      row-key="id"
    >
      <el-table-column label="权限名称" min-width="150">
        <template #default="{ row }">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      
      <el-table-column label="显示名称" min-width="150">
        <template #default="{ row }">
          <span>{{ row.display_name }}</span>
        </template>
      </el-table-column>
      
      <el-table-column label="所属模块" width="120">
        <template #default="{ row }">
          <el-tag type="info">{{ row.module }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="操作类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getActionTagType(row.action)">{{ row.action }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="资源类型" width="120">
        <template #default="{ row }">
          <span>{{ row.resource }}</span>
        </template>
      </el-table-column>
      
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="描述" min-width="200">
        <template #default="{ row }">
          <span>{{ row.description }}</span>
        </template>
      </el-table-column>
      
      <el-table-column align="center" label="操作" width="150">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogFormVisible"
      width="600px"
    >
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="100px"
        style="width: 500px; margin-left:20px;"
      >
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="temp.name" placeholder="如：dashboard:view" />
        </el-form-item>
        
        <el-form-item label="显示名称" prop="display_name">
          <el-input v-model="temp.display_name" placeholder="如：查看仪表盘" />
        </el-form-item>
        
        <el-form-item label="所属模块" prop="module">
          <el-input v-model="temp.module" placeholder="如：dashboard" />
        </el-form-item>
        
        <el-form-item label="操作类型" prop="action">
          <el-select v-model="temp.action" placeholder="请选择" style="width: 100%">
            <el-option label="查看" value="view" />
            <el-option label="创建" value="create" />
            <el-option label="编辑" value="edit" />
            <el-option label="删除" value="delete" />
            <el-option label="导入" value="import" />
            <el-option label="导出" value="export" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="资源类型" prop="resource">
          <el-input v-model="temp.resource" placeholder="如：dashboard" />
        </el-form-item>
        
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="temp.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input v-model="temp.description" type="textarea" rows="3" />
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
import { ref, reactive, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPermissionList, createPermission, updatePermission, deletePermission } from '@/api/permission'

export default {
  name: 'PermissionManagement',
  setup() {
    const dataForm = ref()
    
    const list = ref([])
    const listLoading = ref(true)
    const dialogFormVisible = ref(false)
    const dialogType = ref('')
    
    const listQuery = reactive({
      keyword: ''
    })
    
    const temp = reactive({
      id: undefined,
      name: '',
      display_name: '',
      description: '',
      module: '',
      action: 'view',
      resource: '',
      status: 1
    })

    const rules = {
      name: [{ required: true, message: '权限名称是必需的', trigger: 'blur' }],
      display_name: [{ required: true, message: '显示名称是必需的', trigger: 'blur' }],
      module: [{ required: true, message: '所属模块是必需的', trigger: 'blur' }],
      action: [{ required: true, message: '操作类型是必需的', trigger: 'change' }],
      resource: [{ required: true, message: '资源类型是必需的', trigger: 'blur' }]
    }

    const dialogTitle = computed(() => {
      return dialogType.value === 'create' ? '添加权限' : '编辑权限'
    })

    // 获取操作类型对应的标签类型
    const getActionTagType = (action) => {
      const types = {
        view: '',
        create: 'success',
        edit: 'warning',
        delete: 'danger',
        import: 'info',
        export: 'info'
      }
      return types[action] || ''
    }

    const getList = async () => {
      listLoading.value = true
      try {
        const { data } = await getPermissionList(listQuery)
        list.value = data.items
      } catch (error) {
        console.error('获取权限列表失败:', error)
      } finally {
        listLoading.value = false
      }
    }

    const resetTemp = () => {
      Object.assign(temp, {
        id: undefined,
        name: '',
        display_name: '',
        description: '',
        module: '',
        action: 'view',
        resource: '',
        status: 1
      })
    }

    const handleCreate = () => {
      resetTemp()
      dialogType.value = 'create'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value?.clearValidate()
      })
    }

    const createData = async () => {
      try {
        await dataForm.value?.validate()
        await createPermission(temp)
        dialogFormVisible.value = false
        ElMessage({
          message: '创建成功',
          type: 'success'
        })
        getList()
      } catch (error) {
        console.error('创建权限失败:', error)
      }
    }

    const handleUpdate = (row) => {
      Object.assign(temp, row)
      dialogType.value = 'update'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value?.clearValidate()
      })
    }

    const updateData = async () => {
      try {
        await dataForm.value?.validate()
        await updatePermission(temp.id, temp)
        dialogFormVisible.value = false
        ElMessage({
          message: '更新成功',
          type: 'success'
        })
        getList()
      } catch (error) {
        console.error('更新权限失败:', error)
      }
    }

    const handleDelete = (row) => {
      ElMessageBox.confirm(
        '确定要删除这个权限吗？',
        '警告',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }
      ).then(async () => {
        try {
          await deletePermission(row.id)
          ElMessage({
            message: '删除成功',
            type: 'success'
          })
          getList()
        } catch (error) {
          console.error('删除权限失败:', error)
        }
      })
    }

    const handleFilter = () => {
      getList()
    }

    // 初始化
    getList()

    return {
      dataForm,
      list,
      listLoading,
      listQuery,
      dialogFormVisible,
      dialogType,
      dialogTitle,
      temp,
      rules,
      handleFilter,
      handleCreate,
      handleUpdate,
      handleDelete,
      createData,
      updateData,
      getActionTagType
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