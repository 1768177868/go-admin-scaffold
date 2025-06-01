<template>
  <div class="tags-view-container">
    <el-scrollbar ref="scrollPane" class="tags-view-wrapper" @wheel.prevent="handleScroll">
      <router-link
        v-for="tag in visitedViews"
        ref="tag"
        :key="tag.path"
        :class="isActive(tag) ? 'active' : ''"
        :to="{ path: tag.path, query: tag.query, fullPath: tag.fullPath }"
        tag="span"
        class="tags-view-item"
        @click.middle="!isAffix(tag) && closeSelectedTag(tag)"
        @contextmenu.prevent="openMenu(tag, $event)"
      >
        {{ tag.title }}
        <span v-if="!isAffix(tag)" class="el-icon-close" @click.prevent.stop="closeSelectedTag(tag)">
          <el-icon><Close /></el-icon>
        </span>
      </router-link>
    </el-scrollbar>
    
    <!-- 右键菜单 -->
    <ul v-show="visible" :style="{ left: left + 'px', top: top + 'px' }" class="contextmenu">
      <li @click="refreshSelectedTag(selectedTag)">刷新</li>
      <li v-if="!isAffix(selectedTag)" @click="closeSelectedTag(selectedTag)">关闭</li>
      <li @click="closeOthersTags">关闭其他</li>
      <li @click="closeAllTags(selectedTag)">关闭所有</li>
    </ul>
  </div>
</template>

<script>
import { Close } from '@element-plus/icons-vue'

export default {
  name: 'TagsView',
  components: {
    Close
  },
  setup() {
    const route = useRoute()
    const router = useRouter()

    const visible = ref(false)
    const top = ref(0)
    const left = ref(0)
    const selectedTag = ref({})
    const affixTags = ref([])
    const scrollPane = ref()
    const tag = ref()

    // 模拟visitedViews状态
    const visitedViews = ref([
      {
        name: 'Dashboard',
        path: '/dashboard',
        title: '仪表盘',
        affix: true,
        fullPath: '/dashboard'
      }
    ])

    const isActive = (route) => {
      return route.path === useRoute().path
    }

    const isAffix = (tag) => {
      return tag.affix
    }

    const addTags = () => {
      const { name } = route
      if (name) {
        const view = {
          name: route.name,
          path: route.path,
          title: route.meta?.title || 'no-title',
          fullPath: route.fullPath,
          query: route.query,
          affix: route.meta?.affix
        }
        
        // 检查是否已存在
        const exist = visitedViews.value.find(v => v.path === view.path)
        if (!exist) {
          visitedViews.value.push(view)
        }
      }
    }

    const closeSelectedTag = (view) => {
      const index = visitedViews.value.findIndex(v => v.path === view.path)
      if (index > -1) {
        visitedViews.value.splice(index, 1)
        if (isActive(view)) {
          toLastView(visitedViews.value, view)
        }
      }
    }

    const closeOthersTags = () => {
      router.push(selectedTag.value)
      visitedViews.value = visitedViews.value.filter(view => {
        return view.affix || view.path === selectedTag.value.path
      })
    }

    const closeAllTags = (view) => {
      visitedViews.value = visitedViews.value.filter(tag => tag.affix)
      if (affixTags.value.some(tag => tag.path === view.path)) {
        return
      }
      toLastView(visitedViews.value, view)
    }

    const toLastView = (visitedViews, view) => {
      const latestView = visitedViews.slice(-1)[0]
      if (latestView) {
        router.push(latestView.fullPath)
      } else {
        if (view.name === 'Dashboard') {
          router.replace({ path: '/redirect' + view.fullPath })
        } else {
          router.push('/')
        }
      }
    }

    const refreshSelectedTag = (view) => {
      router.replace({
        path: '/redirect' + view.fullPath
      })
    }

    const openMenu = (tag, e) => {
      const menuMinWidth = 105
      const offsetLeft = scrollPane.value.$el.getBoundingClientRect().left
      const offsetWidth = scrollPane.value.$el.offsetWidth
      const maxLeft = offsetWidth - menuMinWidth
      const left = e.clientX - offsetLeft + 15

      if (left > maxLeft) {
        left.value = maxLeft
      } else {
        left.value = left
      }

      top.value = e.clientY
      visible.value = true
      selectedTag.value = tag
    }

    const closeMenu = () => {
      visible.value = false
    }

    const handleScroll = (e) => {
      const eventDelta = e.wheelDelta || -e.deltaY * 40
      const $scrollWrapper = scrollPane.value.wrap
      $scrollWrapper.scrollLeft = $scrollWrapper.scrollLeft + eventDelta / 4
    }

    // 监听路由变化
    watch(route, addTags, { immediate: true })

    // 监听点击事件关闭菜单
    watch(visible, (value) => {
      if (value) {
        document.body.addEventListener('click', closeMenu)
      } else {
        document.body.removeEventListener('click', closeMenu)
      }
    })

    onMounted(() => {
      addTags()
    })

    return {
      visible,
      top,
      left,
      selectedTag,
      affixTags,
      scrollPane,
      tag,
      visitedViews,
      isActive,
      isAffix,
      closeSelectedTag,
      closeOthersTags,
      closeAllTags,
      refreshSelectedTag,
      openMenu,
      handleScroll
    }
  }
}
</script>

<style lang="scss" scoped>
.tags-view-container {
  height: 34px;
  width: 100%;
  background: #fff;
  border-bottom: 1px solid #d8dce5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, .12), 0 0 3px 0 rgba(0, 0, 0, .04);

  .tags-view-wrapper {
    .tags-view-item {
      display: inline-block;
      position: relative;
      cursor: pointer;
      height: 26px;
      line-height: 26px;
      border: 1px solid #d8dce5;
      color: #495057;
      background: #fff;
      padding: 0 8px;
      font-size: 12px;
      margin-left: 5px;
      margin-top: 4px;
      text-decoration: none;

      &:first-of-type {
        margin-left: 15px;
      }

      &:last-of-type {
        margin-right: 15px;
      }

      &.active {
        background-color: #42b983;
        color: #fff;
        border-color: #42b983;

        &::before {
          content: '';
          background: #fff;
          display: inline-block;
          width: 8px;
          height: 8px;
          border-radius: 50%;
          position: relative;
          margin-right: 2px;
        }
      }
    }
  }

  .contextmenu {
    margin: 0;
    background: #fff;
    z-index: 3000;
    position: absolute;
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 400;
    color: #333;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, .3);

    li {
      margin: 0;
      padding: 7px 16px;
      cursor: pointer;

      &:hover {
        background: #eee;
      }
    }
  }
}

.tags-view-wrapper {
  .el-icon-close {
    width: 16px;
    height: 16px;
    vertical-align: 2px;
    border-radius: 50%;
    text-align: center;
    transition: all .3s cubic-bezier(.645, .045, .355, 1);
    transform-origin: 100% 50%;

    &:before {
      transform: scale(.6);
      display: inline-block;
      vertical-align: -3px;
    }

    &:hover {
      background-color: #b4bccc;
      color: #fff;
    }
  }
}
</style> 