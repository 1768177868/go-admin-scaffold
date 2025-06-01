<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <div class="lang-select">
      <el-icon><Switch /></el-icon>
      <span class="lang-text">{{ currentLang.label }}</span>
      <el-icon class="arrow"><ArrowDown /></el-icon>
    </div>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item
          v-for="lang in languages"
          :key="lang.value"
          :command="lang.value"
          :class="{ active: currentLang.value === lang.value }"
        >
          {{ lang.label }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script>
import { Switch, ArrowDown } from '@element-plus/icons-vue'
import { useI18n } from '@/composables/useI18n'

export default {
  name: 'LangSelect',
  components: {
    Switch,
    ArrowDown
  },
  setup() {
    const { locale, setLocale } = useI18n()
    
    const languages = [
      { label: '中文', value: 'zh' },
      { label: 'English', value: 'en' }
    ]
    
    const currentLang = computed(() => {
      return languages.find(lang => lang.value === locale.value) || languages[0]
    })
    
    const handleCommand = (command) => {
      setLocale(command)
      ElMessage.success(`语言已切换为${languages.find(l => l.value === command)?.label}`)
    }
    
    return {
      languages,
      currentLang,
      handleCommand
    }
  }
}
</script>

<style scoped>
.lang-select {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.lang-select:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.lang-text {
  margin: 0 4px;
  font-size: 14px;
}

.arrow {
  font-size: 12px;
  transition: transform 0.3s;
}

.el-dropdown-menu__item.active {
  color: #409eff;
  background-color: #ecf5ff;
}
</style> 