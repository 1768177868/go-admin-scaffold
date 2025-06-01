<template>
  <div class="app-container">
    <!-- 搜索栏 -->
    <div class="filter-container">
      <el-input
        v-model="listQuery.username"
        placeholder="请输入用户名"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-input
        v-model="listQuery.ip"
        placeholder="请输入IP地址"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter="handleFilter"
      />
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        class="filter-item"
        @change="handleDateChange"
      />
      <el-button class="filter-item" type="primary" icon="Search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" type="warning" icon="Delete" @click="handleClearLogs">
        清理日志
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
      <el-table-column align="center" label="ID" width="80">
        <template #default="{ row }">
          {{ row.id }}
        </template>
      </el-table-column>
      
      <el-table-column label="用户名" width="120">
        <template #default="{ row }">
          {{ row.username }}
        </template>
      </el-table-column>
      
      <el-table-column label="登录IP" width="140">
        <template #default="{ row }">
          {{ row.ip }}
        </template>
      </el-table-column>
      
      <el-table-column label="登录地址" width="200">
        <template #default="{ row }">
          {{ row.location || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column label="浏览器">
        <template #default="{ row }">
          {{ row.user_agent || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column label="登录状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
            {{ row.status === 'success' ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="失败原因" width="150">
        <template #default="{ row }">
          {{ row.message || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column align="center" prop="created_at" label="登录时间" width="160">
        <template #default="{ row }">
          <span>{{ formatDate(row.created_at) }}</span>
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

    <!-- 清理日志对话框 -->
    <el-dialog title="清理登录日志" v-model="clearDialogVisible" width="400px">
      <el-form label-width="120px">
        <el-form-item label="保留天数：">
          <el-input-number v-model="clearDays" :min="1" :max="365" />
          <div style="margin-top: 5px; color: #666; font-size: 12px;">
            将删除 {{ clearDays }} 天前的所有登录日志
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="clearDialogVisible = false">取消</el-button>
          <el-button type="danger" @click="confirmClearLogs">确认清理</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { getLoginLogs, clearLogs } from '@/api/log'
import Pagination from '@/components/Pagination/index.vue'
import dayjs from 'dayjs'

export default {
  name: 'LoginLog',
  components: {
    Pagination
  },
  setup() {
    const list = ref([])
    const total = ref(0)
    const listLoading = ref(true)
    const clearDialogVisible = ref(false)
    const clearDays = ref(30)
    const dateRange = ref([])
    
    const listQuery = reactive({
      page: 1,
      limit: 20,
      username: '',
      ip: '',
      start_date: '',
      end_date: ''
    })

    const getList = async () => {
      listLoading.value = true
      try {
        const { data } = await getLoginLogs(listQuery)
        list.value = data.items || []
        total.value = data.total || 0
      } catch (error) {
        console.error('Failed to get login logs:', error)
      } finally {
        listLoading.value = false
      }
    }

    const handleFilter = () => {
      listQuery.page = 1
      getList()
    }

    const handleDateChange = (dates) => {
      if (dates && dates.length === 2) {
        listQuery.start_date = dayjs(dates[0]).format('YYYY-MM-DD')
        listQuery.end_date = dayjs(dates[1]).format('YYYY-MM-DD')
      } else {
        listQuery.start_date = ''
        listQuery.end_date = ''
      }
    }

    const handleClearLogs = () => {
      clearDialogVisible.value = true
    }

    const confirmClearLogs = async () => {
      try {
        await ElMessageBox.confirm(
          `确定要清理 ${clearDays.value} 天前的登录日志吗？此操作不可恢复！`,
          '确认清理',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await clearLogs('login', clearDays.value)
        clearDialogVisible.value = false
        ElMessage.success('日志清理成功')
        getList()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('Failed to clear logs:', error)
        }
      }
    }

    const formatDate = (date) => {
      return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    }

    onMounted(() => {
      getList()
    })

    return {
      list,
      total,
      listLoading,
      clearDialogVisible,
      clearDays,
      dateRange,
      listQuery,
      getList,
      handleFilter,
      handleDateChange,
      handleClearLogs,
      confirmClearLogs,
      formatDate
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