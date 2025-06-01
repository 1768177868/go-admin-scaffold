import { ref, computed } from 'vue'

// 全局语言状态
const locale = ref(localStorage.getItem('locale') || 'zh')

// 翻译数据存储
const messages = ref({
  zh: {},
  en: {}
})

// 加载翻译文件
const loadMessages = async (lang) => {
  if (!messages.value[lang] || Object.keys(messages.value[lang]).length === 0) {
    try {
      const response = await fetch(`/api/admin/v1/i18n/translations?locale=${lang}`)
      const data = await response.json()
      if (data.code === 0) {
        messages.value[lang] = data.data
      }
    } catch (error) {
      console.error(`Failed to load ${lang} translations:`, error)
      // 如果加载失败，使用默认翻译
      if (lang === 'zh') {
        messages.value[lang] = {
          common: {
            welcome: '欢迎使用 Go Admin',
            save: '保存',
            cancel: '取消',
            delete: '删除',
            edit: '编辑',
            create: '创建',
            search: '搜索',
            actions: '操作'
          },
          user: {
            username: '用户名',
            password: '密码',
            email: '邮箱',
            status: '状态'
          }
        }
      } else {
        messages.value[lang] = {
          common: {
            welcome: 'Welcome to Go Admin',
            save: 'Save',
            cancel: 'Cancel',
            delete: 'Delete',
            edit: 'Edit',
            create: 'Create',
            search: 'Search',
            actions: 'Actions'
          },
          user: {
            username: 'Username',
            password: 'Password',
            email: 'Email',
            status: 'Status'
          }
        }
      }
    }
  }
}

// 获取翻译文本
const t = (key, params = {}) => {
  const keys = key.split('.')
  let value = messages.value[locale.value]
  
  for (const k of keys) {
    if (value && typeof value === 'object') {
      value = value[k]
    } else {
      value = undefined
      break
    }
  }
  
  if (typeof value === 'string') {
    // 简单的参数替换
    return value.replace(/\{(\w+)\}/g, (match, param) => {
      return params[param] !== undefined ? params[param] : match
    })
  }
  
  return key // 如果找不到翻译，返回原始key
}

// 设置语言
const setLocale = async (newLocale) => {
  if (newLocale !== locale.value) {
    locale.value = newLocale
    localStorage.setItem('locale', newLocale)
    await loadMessages(newLocale)
  }
}

// 初始化
const initI18n = async () => {
  await loadMessages(locale.value)
}

export function useI18n() {
  return {
    locale: computed(() => locale.value),
    t,
    setLocale,
    initI18n
  }
} 