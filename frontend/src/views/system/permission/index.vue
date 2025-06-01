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
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <el-table-column label="权限名称" min-width="200">
        <template #default="{ row }">
          <span>{{ row.name }}</span>
        </template>
      </el-table-column>
      
      <el-table-column label="权限代码" width="200">
        <template #default="{ row }">
          <el-tag>{{ row.code }}</el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getPermissionTypeTag(row.type)">
            {{ getPermissionTypeName(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="路径" width="200">
        <template #default="{ row }">
          {{ row.path }}
        </template>
      </el-table-column>
      
      <el-table-column label="排序" width="80" align="center">
        <template #default="{ row }">
          {{ row.sort }}
        </template>
      </el-table-column>
      
      <el-table-column label="描述">
        <template #default="{ row }">
          {{ row.description }}
        </template>
      </el-table-column>
      
      <el-table-column align="center" label="操作" width="200">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button v-if="row.type === 'menu'" type="success" size="small" @click="handleAddChild(row)">
            添加子项
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
        <el-form-item label="父级权限" prop="parent_id">
          <el-tree-select
            v-model="temp.parent_id"
            :data="permissionTreeOptions"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            placeholder="选择父级权限（可为空）"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="temp.name" />
        </el-form-item>
        
        <el-form-item label="权限代码" prop="code">
          <el-input v-model="temp.code" />
        </el-form-item>
        
        <el-form-item label="权限类型" prop="type">
          <el-select v-model="temp.type" placeholder="请选择" style="width: 100%">
            <el-option label="菜单" value="menu" />
            <el-option label="按钮" value="button" />
            <el-option label="接口" value="api" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="路径/接口" prop="path">
          <el-input v-model="temp.path" placeholder="菜单路径或API接口" />
        </el-form-item>
        
        <el-form-item label="图标" prop="icon">
          <el-input v-model="temp.icon" placeholder="图标名称（仅菜单需要）" />
        </el-form-item>
        
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="temp.sort" :min="0" style="width: 100%" />
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
import { getPermissionList, createPermission, updatePermission, deletePermission } from '@/api/permission'

export default {
  name: 'PermissionManagement',
  setup() {
    const dataForm = ref()
    
    const list = ref([])
    const listLoading = ref(true)
    const dialogFormVisible = ref(false)
    const dialogType = ref('')
    const permissionTreeOptions = ref([])
    
    const listQuery = reactive({
      keyword: ''
    })
    
    const temp = reactive({
      id: undefined,
      parent_id: null,
      name: '',
      code: '',
      type: 'menu',
      path: '',
      icon: '',
      sort: 0,
      description: ''
    })

    const rules = {
      name: [{ required: true, message: '权限名称是必需的', trigger: 'blur' }],
      code: [{ required: true, message: '权限代码是必需的', trigger: 'blur' }],
      type: [{ required: true, message: '权限类型是必需的', trigger: 'change' }]
    }

    const dialogTitle = computed(() => {
      return dialogType.value === 'create' ? '添加权限' : '编辑权限'
    })

    const getList = async () => {
      listLoading.value = true
      try {
        const { data } = await getPermissionList(listQuery)
        list.value = data.items || []
        
        // 构建树形选择器选项
        permissionTreeOptions.value = buildTreeOptions(data.items || [])
      } catch (error) {
        console.error('Failed to get permission list:', error)
      } finally {
        listLoading.value = false
      }
    }

    const buildTreeOptions = (permissions) => {
      const result = [{ id: null, name: '顶级权限', children: [] }]
      
      const buildTree = (items, parentId = null) => {
        return items
          .filter(item => item.parent_id === parentId)
          .map(item => ({
            id: item.id,
            name: item.name,
            children: buildTree(items, item.id)
          }))
      }
      
      result[0].children = buildTree(permissions)
      return result
    }

    const handleFilter = () => {
      getList()
    }

    const resetTemp = () => {
      temp.id = undefined
      temp.parent_id = null
      temp.name = ''
      temp.code = ''
      temp.type = 'menu'
      temp.path = ''
      temp.icon = ''
      temp.sort = 0
      temp.description = ''
    }

    const handleCreate = () => {
      resetTemp()
      dialogType.value = 'create'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value.clearValidate()
      })
    }

    const handleAddChild = (row) => {
      resetTemp()
      temp.parent_id = row.id
      dialogType.value = 'create'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value.clearValidate()
      })
    }

    const createData = async () => {
      await dataForm.value.validate(async (valid) => {
        if (valid) {
          try {
            await createPermission(temp)
            dialogFormVisible.value = false
            ElMessage.success('创建成功')
            getList()
          } catch (error) {
            console.error('Failed to create permission:', error)
          }
        }
      })
    }

    const handleUpdate = (row) => {
      Object.assign(temp, row)
      dialogType.value = 'update'
      dialogFormVisible.value = true
      nextTick(() => {
        dataForm.value.clearValidate()
      })
    }

    const updateData = async () => {
      await dataForm.value.validate(async (valid) => {
        if (valid) {
          try {
            await updatePermission(temp.id, temp)
            dialogFormVisible.value = false
            ElMessage.success('更新成功')
            getList()
          } catch (error) {
            console.error('Failed to update permission:', error)
          }
        }
      })
    }

    const handleDelete = async (row) => {
      await ElMessageBox.confirm('此操作将永久删除该权限及其子权限, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      
      try {
        await deletePermission(row.id)
        ElMessage.success('删除成功')
        getList()
      } catch (error) {
        console.error('Failed to delete permission:', error)
      }
    }

    const getPermissionTypeName = (type) => {
      const types = {
        menu: '菜单',
        button: '按钮',
        api: '接口'
      }
      return types[type] || type
    }

    const getPermissionTypeTag = (type) => {
      const tags = {
        menu: 'primary',
        button: 'success',
        api: 'warning'
      }
      return tags[type] || ''
    }

    onMounted(() => {
      getList()
    })

    return {
      dataForm,
      list,
      listLoading,
      dialogFormVisible,
      dialogType,
      dialogTitle,
      permissionTreeOptions,
      listQuery,
      temp,
      rules,
      getList,
      handleFilter,
      handleCreate,
      handleAddChild,
      createData,
      handleUpdate,
      updateData,
      handleDelete,
      getPermissionTypeName,
      getPermissionTypeTag
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