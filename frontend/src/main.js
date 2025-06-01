import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import 'nprogress/nprogress.css'

import App from './App.vue'
import router from './router'
import './assets/styles/index.css'
import { useI18n } from '@/composables/useI18n'

const app = createApp(App)

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus, {
  locale: {
    name: 'zh-cn'
  }
})

// 挂载应用
app.mount('#app')

// 异步初始化国际化（不阻塞应用启动）
const { initI18n } = useI18n()
initI18n().catch(console.error) 