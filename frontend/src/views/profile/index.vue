<template>
  <div class="app-container">
    <el-row :gutter="20">
      <!-- 左侧用户信息 -->
      <el-col :span="8">
        <el-card class="user-info-card">
          <div class="user-info">
            <div class="avatar-container">
              <img :src="userInfo.avatar || defaultAvatar" class="user-avatar" alt="头像">
              <div class="avatar-actions">
                <el-button type="primary" size="small" @click="showUploadDialog">
                  <el-icon><Upload /></el-icon>
                  更换头像
                </el-button>
              </div>
            </div>
            
            <div class="user-details">
              <h2>{{ userInfo.username }}</h2>
              <p class="user-email">{{ userInfo.email }}</p>
              <div class="user-roles">
                <el-tag v-for="role in userInfo.roles" :key="role.id" type="primary" size="small">
                  {{ role.name }}
                </el-tag>
              </div>
            </div>
          </div>
          
          <el-divider />
          
          <div class="user-stats">
            <div class="stat-item">
              <div class="stat-number">{{ userStats.loginCount }}</div>
              <div class="stat-label">登录次数</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ userStats.operationCount }}</div>
              <div class="stat-label">操作次数</div>
            </div>
            <div class="stat-item">
              <div class="stat-number">{{ formatDate(userInfo.last_login_at) }}</div>
              <div class="stat-label">最后登录</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 右侧信息编辑 -->
      <el-col :span="16">
        <el-card>
          <template #header>
            <el-tabs v-model="activeTab">
              <el-tab-pane label="基本信息" name="basic" />
              <el-tab-pane label="修改密码" name="password" />
              <el-tab-pane label="登录记录" name="logs" />
            </el-tabs>
          </template>
          
          <!-- 基本信息 -->
          <div v-show="activeTab === 'basic'" class="tab-content">
            <el-form
              ref="profileForm"
              :model="profileData"
              :rules="profileRules"
              label-width="100px"
              style="max-width: 600px;"
            >
              <el-form-item label="用户名" prop="username">
                <el-input v-model="profileData.username" disabled />
              </el-form-item>
              
              <el-form-item label="邮箱" prop="email">
                <el-input v-model="profileData.email" />
              </el-form-item>
              
              <el-form-item label="昵称" prop="nickname">
                <el-input v-model="profileData.nickname" />
              </el-form-item>
              
              <el-form-item label="头像" prop="avatar">
                <el-upload
                  class="avatar-uploader"
                  :show-file-list="false"
                  :on-success="handleAvatarSuccess"
                  :before-upload="beforeAvatarUpload"
                >
                  <img v-if="profileData.avatar" :src="profileData.avatar" class="avatar" />
                  <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
                </el-upload>
              </el-form-item>
              
              <el-form-item>
                <el-button type="primary" @click="updateProfile">保存修改</el-button>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 修改密码 -->
          <div v-show="activeTab === 'password'" class="tab-content">
            <el-form
              ref="passwordForm"
              :model="passwordData"
              :rules="passwordRules"
              label-width="100px"
              style="max-width: 600px;"
            >
              <el-form-item label="当前密码" prop="old_password">
                <el-input v-model="passwordData.old_password" type="password" show-password />
              </el-form-item>
              
              <el-form-item label="新密码" prop="new_password">
                <el-input v-model="passwordData.new_password" type="password" show-password />
              </el-form-item>
              
              <el-form-item label="确认密码" prop="confirm_password">
                <el-input v-model="passwordData.confirm_password" type="password" show-password />
              </el-form-item>
              
              <el-form-item>
                <el-button type="primary" @click="updatePassword">修改密码</el-button>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 登录记录 -->
          <div v-show="activeTab === 'logs'" class="tab-content">
            <el-table :data="loginLogs" border>
              <el-table-column label="登录IP" prop="ip" width="140" />
              <el-table-column label="登录地址" prop="location" />
              <el-table-column label="浏览器" prop="user_agent" />
              <el-table-column label="登录时间" prop="created_at" width="160">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
            
            <div style="margin-top: 20px; text-align: center;">
              <el-button @click="loadMoreLogs">加载更多</el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 头像上传对话框 -->
    <el-dialog title="上传头像" v-model="uploadDialogVisible" width="500px">
      <el-upload
        class="avatar-uploader"
        action="#"
        :show-file-list="false"
        :before-upload="beforeAvatarUpload"
        :http-request="handleAvatarUpload"
      >
        <img v-if="newAvatar" :src="newAvatar" class="avatar-preview" />
        <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
      </el-upload>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmUploadAvatar" :disabled="!newAvatar">
            确认上传
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { Upload, Plus } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { getUserLogs } from '@/api/log'
import dayjs from 'dayjs'

export default {
  name: 'Profile',
  components: {
    Upload,
    Plus
  },
  setup() {
    const authStore = useAuthStore()
    const profileForm = ref()
    const passwordForm = ref()
    
    const activeTab = ref('basic')
    const uploadDialogVisible = ref(false)
    const newAvatar = ref('')
    const loginLogs = ref([])
    const userStats = ref({
      loginCount: 0,
      operationCount: 0
    })
    
    const defaultAvatar = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'
    
    const userInfo = computed(() => authStore.userInfo || {})
    
    const profileData = reactive({
      username: '',
      email: '',
      nickname: '',
      avatar: ''
    })
    
    const passwordData = reactive({
      old_password: '',
      new_password: '',
      confirm_password: ''
    })
    
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== passwordData.new_password) {
        callback(new Error('两次输入密码不一致'))
      } else {
        callback()
      }
    }
    
    const profileRules = {
      email: [
        { required: true, message: '请输入邮箱地址', trigger: 'blur' },
        { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
      ]
    }
    
    const passwordRules = {
      old_password: [
        { required: true, message: '请输入当前密码', trigger: 'blur' }
      ],
      new_password: [
        { required: true, message: '请输入新密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
      ],
      confirm_password: [
        { required: true, message: '请确认新密码', trigger: 'blur' },
        { validator: validateConfirmPassword, trigger: 'blur' }
      ]
    }
    
    const initProfileData = () => {
      Object.assign(profileData, {
        username: userInfo.value.username || '',
        email: userInfo.value.email || '',
        nickname: userInfo.value.nickname || '',
        avatar: userInfo.value.avatar || ''
      })
    }
    
    const updateProfile = async () => {
      await profileForm.value.validate(async (valid) => {
        if (valid) {
          try {
            // 这里应该调用更新用户信息的API
            ElMessage.success('个人信息更新成功')
            await authStore.getUserInfo() // 重新获取用户信息
          } catch (error) {
            console.error('Failed to update profile:', error)
          }
        }
      })
    }
    
    const updatePassword = async () => {
      await passwordForm.value.validate(async (valid) => {
        if (valid) {
          try {
            // 这里应该调用修改密码的API
            ElMessage.success('密码修改成功')
            // 重置表单
            Object.assign(passwordData, {
              old_password: '',
              new_password: '',
              confirm_password: ''
            })
          } catch (error) {
            console.error('Failed to update password:', error)
          }
        }
      })
    }
    
    const showUploadDialog = () => {
      newAvatar.value = ''
      uploadDialogVisible.value = true
    }
    
    const beforeAvatarUpload = (file) => {
      const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
      const isLt2M = file.size / 1024 / 1024 < 2
      
      if (!isJPG) {
        ElMessage.error('上传头像图片只能是 JPG/PNG 格式!')
        return false
      }
      if (!isLt2M) {
        ElMessage.error('上传头像图片大小不能超过 2MB!')
        return false
      }
      return true
    }
    
    const handleAvatarUpload = (options) => {
      const file = options.file
      const reader = new FileReader()
      reader.onload = (e) => {
        newAvatar.value = e.target.result
      }
      reader.readAsDataURL(file)
    }
    
    const confirmUploadAvatar = async () => {
      try {
        // 这里应该调用上传头像的API
        ElMessage.success('头像上传成功')
        uploadDialogVisible.value = false
        await authStore.getUserInfo() // 重新获取用户信息
      } catch (error) {
        console.error('Failed to upload avatar:', error)
      }
    }
    
    const loadLoginLogs = async () => {
      try {
        const { data } = await getUserLogs(userInfo.value.id)
        loginLogs.value = data.items || []
        userStats.value = {
          loginCount: data.login_count || 0,
          operationCount: data.operation_count || 0
        }
      } catch (error) {
        console.error('Failed to load login logs:', error)
      }
    }
    
    const loadMoreLogs = () => {
      // 实现加载更多日志的逻辑
      ElMessage.info('加载更多功能待实现')
    }
    
    const formatDate = (date) => {
      if (!date) return '-'
      return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    }
    
    watch(activeTab, (newTab) => {
      if (newTab === 'logs') {
        loadLoginLogs()
      }
    })
    
    onMounted(() => {
      initProfileData()
    })
    
    watch(userInfo, () => {
      initProfileData()
    })
    
    return {
      activeTab,
      uploadDialogVisible,
      newAvatar,
      loginLogs,
      userStats,
      defaultAvatar,
      userInfo,
      profileData,
      passwordData,
      profileRules,
      passwordRules,
      profileForm,
      passwordForm,
      updateProfile,
      updatePassword,
      showUploadDialog,
      beforeAvatarUpload,
      handleAvatarUpload,
      confirmUploadAvatar,
      loadMoreLogs,
      formatDate
    }
  }
}
</script>

<style lang="scss" scoped>
.app-container {
  padding: 20px;
}

.user-info-card {
  .user-info {
    text-align: center;
    
    .avatar-container {
      margin-bottom: 20px;
      
      .user-avatar {
        width: 120px;
        height: 120px;
        border-radius: 50%;
        object-fit: cover;
        border: 4px solid #fff;
        box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
      }
      
      .avatar-actions {
        margin-top: 15px;
      }
    }
    
    .user-details {
      h2 {
        margin-bottom: 10px;
        color: #303133;
      }
      
      .user-email {
        color: #909399;
        margin-bottom: 15px;
      }
      
      .user-roles {
        .el-tag {
          margin-right: 5px;
        }
      }
    }
  }
  
  .user-stats {
    display: flex;
    justify-content: space-around;
    
    .stat-item {
      text-align: center;
      
      .stat-number {
        font-size: 24px;
        font-weight: 600;
        color: #409eff;
        margin-bottom: 5px;
      }
      
      .stat-label {
        font-size: 14px;
        color: #909399;
      }
    }
  }
}

.tab-content {
  padding: 20px 0;
}

.avatar-uploader {
  :deep(.el-upload) {
    border: 1px dashed #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: 0.2s;
    
    &:hover {
      border-color: #409eff;
    }
  }
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.avatar-preview {
  width: 178px;
  height: 178px;
  display: block;
  object-fit: cover;
}
</style> 