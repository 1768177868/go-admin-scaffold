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
      <el-select v-model="listQuery.method" placeholder="请求方法" clearable style="width: 120px" class="filter-item">
        <el-option label="GET" value="GET" />
        <el-option label="POST" value="POST" />
        <el-option label="PUT" value="PUT" />
        <el-option label="DELETE" value="DELETE" />
      </el-select>
      <el-input
        v-model="listQuery.path"
        placeholder="请求路径"
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
      
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          {{ row.description || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column label="请求方法" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getMethodTag(row.method)">
            {{ row.method }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="请求路径" width="250">
        <template #default="{ row }">
          <el-tooltip :content="row.path" placement="top">
            <span>{{ row.path }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      
      <el-table-column label="IP地址" width="140">
        <template #default="{ row }">
          {{ row.ip }}
        </template>
      </el-table-column>
      
      <el-table-column label="响应状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTag(row.status_code)">
            {{ row.status_code }}
          </el-tag>
        </template>
      </el-table-column>
      
      <el-table-column label="耗时(ms)" width="100" align="center">
        <template #default="{ row }">
          {{ row.duration || '-' }}
        </template>
      </el-table-column>
      
      <el-table-column align="center" prop="created_at" label="操作时间" width="160">
        <template #default="{ row }">
          <span>{{ formatDate(row.created_at) }}</span>
        </template>
      </el-table-column>
      
      <el-table-column align="center" label="操作" width="100">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="showDetail(row)">
            详情
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

    <!-- 详情对话框 -->
    <el-dialog title="操作详情" v-model="detailDialogVisible" width="800px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户名">{{ currentLog.username }}</el-descriptions-item>
        <el-descriptions-item label="操作描述">{{ currentLog.description }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">
          <el-tag :type="getMethodTag(currentLog.method)">{{ currentLog.method }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="请求路径">{{ currentLog.path }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ currentLog.ip }}</el-descriptions-item>
        <el-descriptions-item label="响应状态">
          <el-tag :type="getStatusTag(currentLog.status_code)">{{ currentLog.status_code }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="响应时间">{{ currentLog.duration }}ms</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ formatDate(currentLog.created_at) }}</el-descriptions-item>
      </el-descriptions>
      
      <div style="margin-top: 20px;">
        <h4>请求参数：</h4>
        <el-input
          v-model="currentLog.request_body"
          type="textarea"
          :rows="6"
          readonly
          style="margin-top: 10px;"
        />
      </div>
      
      <div style="margin-top: 20px;">
        <h4>响应数据：</h4>
        <el-input
          v-model="currentLog.response_body"
          type="textarea"
          :rows="6"
          readonly
          style="margin-top: 10px;"
        />
      </div>
    </el-dialog>

    <!-- 清理日志对话框 -->
    <el-dialog title="清理操作日志" v-model="clearDialogVisible" width="400px">
      <el-form label-width="120px">
        <el-form-item label="保留天数：">
          <el-input-number v-model="clearDays" :min="1" :max="365" />
          <div style="margin-top: 5px; color: #666; font-size: 12px;">
            将删除 {{ clearDays }} 天前的所有操作日志
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
import { getOperationLogs, clearLogs } from '@/api/log'
import Pagination from '@/components/Pagination/index.vue'
import dayjs from 'dayjs'

export default {
  name: 'OperationLog',
  components: {
    Pagination
  },
  setup() {
    const list = ref([])
    const total = ref(0)
    const listLoading = ref(true)
    const detailDialogVisible = ref(false)
    const clearDialogVisible = ref(false)
    const clearDays = ref(30)
    const dateRange = ref([])
    const currentLog = ref({})
    
    const listQuery = reactive({
      page: 1,
      limit: 20,
      username: '',
      method: '',
      path: '',
      start_date: '',
      end_date: ''
    })

    const getList = async () => {
      listLoading.value = true
      try {
        const { data } = await getOperationLogs(listQuery)
        list.value = data.items || []
        total.value = data.total || 0
      } catch (error) {
        console.error('Failed to get operation logs:', error)
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

    const showDetail = (row) => {
      currentLog.value = { ...row }
      detailDialogVisible.value = true
    }

    const handleClearLogs = () => {
      clearDialogVisible.value = true
    }

    const confirmClearLogs = async () => {
      try {
        await ElMessageBox.confirm(
          `确定要清理 ${clearDays.value} 天前的操作日志吗？此操作不可恢复！`,
          '确认清理',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await clearLogs('operation', clearDays.value)
        clearDialogVisible.value = false
        ElMessage.success('日志清理成功')
        getList()
      } catch (error) {
        if (error !== 'cancel') {
          console.error('Failed to clear logs:', error)
        }
      }
    }

    const getMethodTag = (method) => {
      const tags = {
        GET: 'info',
        POST: 'success',
        PUT: 'warning',
        DELETE: 'danger'
      }
      return tags[method] || ''
    }

    const getStatusTag = (status) => {
      if (status >= 200 && status < 300) return 'success'
      if (status >= 300 && status < 400) return 'warning'
      if (status >= 400) return 'danger'
      return 'info'
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
      detailDialogVisible,
      clearDialogVisible,
      clearDays,
      dateRange,
      currentLog,
      listQuery,
      getList,
      handleFilter,
      handleDateChange,
      showDetail,
      handleClearLogs,
      confirmClearLogs,
      getMethodTag,
      getStatusTag,
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