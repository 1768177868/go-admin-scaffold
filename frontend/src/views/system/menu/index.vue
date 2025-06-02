<template>
  <div class="menu-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>菜单管理</span>
          <el-button type="primary" @click="handleAdd" v-permission="'menu:create'">
            <el-icon><Plus /></el-icon>
            新增菜单
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="menuData"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        border
      >
        <el-table-column prop="title" label="菜单名称" min-width="200">
          <template #default="{ row }">
            <el-icon v-if="row.icon" style="margin-right: 8px">
              <component :is="row.icon" />
            </el-icon>
            <span>{{ row.title }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="name" label="路由名称" width="150" />
        
        <el-table-column prop="path" label="路由路径" width="200" />
        
        <el-table-column prop="component" label="组件路径" width="200" show-overflow-tooltip />
        
        <el-table-column prop="permission" label="权限标识" width="150" />
        
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        
        <el-table-column prop="type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.type === 1 ? 'primary' : 'info'">
              {{ row.type === 1 ? '菜单' : '按钮' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="visible" label="可见" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.visible === 1 ? 'success' : 'danger'">
              {{ row.visible === 1 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" align="center">
          <template #default="{ row }">
            <el-button 
              link 
              type="primary" 
              size="small" 
              @click="handleEdit(row)"
              v-permission="'menu:edit'"
            >
              编辑
            </el-button>
            <el-button 
              link 
              type="primary" 
              size="small" 
              @click="handleAddChild(row)"
              v-permission="'menu:create'"
            >
              新增子菜单
            </el-button>
            <el-button 
              link 
              type="danger" 
              size="small" 
              @click="handleDelete(row)"
              v-permission="'menu:delete'"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 菜单表单对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="isEdit ? '编辑菜单' : '新增菜单'" 
      width="700px"
    >
      <el-form 
        ref="formRef" 
        :model="form" 
        :rules="rules" 
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="菜单名称" prop="title">
              <el-input v-model="form.title" placeholder="请输入菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="路由名称" prop="name">
              <el-input v-model="form.name" placeholder="请输入路由名称" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="上级菜单" prop="parent_id">
              <el-tree-select
                v-model="form.parent_id"
                :data="menuTreeOptions"
                :props="{ value: 'id', label: 'title', children: 'children' }"
                :render-after-expand="false"
                placeholder="请选择上级菜单"
                clearable
                check-strictly
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="菜单图标" prop="icon">
              <el-input v-model="form.icon" placeholder="请输入图标名称">
                <template #prefix>
                  <el-icon v-if="form.icon">
                    <component :is="form.icon" />
                  </el-icon>
                </template>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="路由路径" prop="path">
              <el-input v-model="form.path" placeholder="请输入路由路径" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="组件路径" prop="component">
              <el-input v-model="form.component" placeholder="请输入组件路径" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="权限标识" prop="permission">
              <el-input v-model="form.permission" placeholder="请输入权限标识" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="排序" prop="sort">
              <el-input-number v-model="form.sort" :min="0" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="菜单类型" prop="type">
              <el-radio-group v-model="form.type">
                <el-radio :label="1">菜单</el-radio>
                <el-radio :label="2">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="是否可见" prop="visible">
              <el-radio-group v-model="form.visible">
                <el-radio :label="1">显示</el-radio>
                <el-radio :label="0">隐藏</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="菜单状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :label="1">启用</el-radio>
                <el-radio :label="0">禁用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="是否缓存" prop="keep_alive">
              <el-switch v-model="form.keep_alive" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="是否外链" prop="external">
              <el-switch v-model="form.external" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import menuApi from '@/api/menu'

// 响应式数据
const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const menuData = ref([])
const menuTreeOptions = ref([])
const formRef = ref(null)

// 表单数据
const initForm = {
  id: null,
  name: '',
  title: '',
  icon: '',
  path: '',
  component: '',
  parent_id: null,
  sort: 0,
  type: 1,
  visible: 1,
  status: 1,
  keep_alive: false,
  external: false,
  permission: '',
  meta: {
    title: '',
    icon: '',
    hidden: false,
    alwaysShow: false,
    noCache: false,
    affix: false,
    breadcrumb: true,
    activeMenu: ''
  }
}

const form = reactive({ ...initForm })

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入菜单名称', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入路由名称', trigger: 'blur' }
  ],
  path: [
    { required: true, message: '请输入路由路径', trigger: 'blur' }
  ]
}

// 获取菜单列表
const getMenuList = async () => {
  try {
    loading.value = true
    const { data } = await menuApi.getMenuTree()
    menuData.value = data || []
    
    // 构建菜单树选项（用于父菜单选择）
    menuTreeOptions.value = buildTreeOptions(data || [])
  } catch (error) {
    console.error('获取菜单列表失败:', error)
    ElMessage.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

// 构建树形选项数据
const buildTreeOptions = (data) => {
  return data.map(item => ({
    id: item.id,
    title: item.title,
    children: item.children ? buildTreeOptions(item.children) : []
  }))
}

// 重置表单
const resetForm = () => {
  Object.assign(form, initForm)
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

// 新增菜单
const handleAdd = () => {
  resetForm()
  isEdit.value = false
  dialogVisible.value = true
}

// 新增子菜单
const handleAddChild = (row) => {
  resetForm()
  form.parent_id = row.id
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑菜单
const handleEdit = async (row) => {
  try {
    const { data } = await menuApi.getMenu(row.id)
    Object.assign(form, data)
    isEdit.value = true
    dialogVisible.value = true
  } catch (error) {
    console.error('获取菜单详情失败:', error)
    ElMessage.error('获取菜单详情失败')
  }
}

// 删除菜单
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜单"${row.title}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await menuApi.deleteMenu(row.id)
    ElMessage.success('删除成功')
    await getMenuList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除菜单失败:', error)
      ElMessage.error('删除菜单失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    submitLoading.value = true
    
    // 更新 meta 信息
    form.meta.title = form.title
    form.meta.icon = form.icon
    form.meta.hidden = form.visible === 0
    
    if (isEdit.value) {
      await menuApi.updateMenu(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await menuApi.createMenu(form)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    await getMenuList()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  } finally {
    submitLoading.value = false
  }
}

// 初始化
onMounted(() => {
  getMenuList()
})
</script>

<style scoped>
.menu-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 