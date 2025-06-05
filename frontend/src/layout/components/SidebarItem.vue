<template>
  <div v-if="!item.meta?.hidden">
    <!-- 只有一个子菜单且不是总是显示根菜单 -->
    <template v-if="hasOneShowingChild(item.children, item) && (!onlyOneChild.children || onlyOneChild.noShowingChildren) && !item.alwaysShow">
      <app-link v-if="onlyOneChild.meta" :to="resolvePath(onlyOneChild.path)">
        <el-menu-item :index="resolvePath(onlyOneChild.path)" :class="{ 'submenu-title-noDropdown': !isNest }">
          <el-icon v-if="onlyOneChild.meta.icon">
            <component :is="onlyOneChild.meta.icon" />
          </el-icon>
          <template #title>
            <span>{{ onlyOneChild.meta.title }}</span>
          </template>
        </el-menu-item>
      </app-link>
    </template>

    <!-- 多个子菜单 -->
    <el-sub-menu v-else :index="resolvePath(item.path)" popper-append-to-body>
      <template #title>
        <el-icon v-if="item.meta?.icon">
          <component :is="item.meta.icon" />
        </el-icon>
        <span>{{ item.meta?.title }}</span>
      </template>
      
      <sidebar-item
        v-for="child in item.children"
        :key="child.path"
        :item="child"
        :is-nest="true"
        :base-path="resolvePath(child.path)"
        class="nest-menu"
      />
    </el-sub-menu>
  </div>
</template>

<script>
import path from 'path-browserify'
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import AppLink from './AppLink.vue'

export default {
  name: 'SidebarItem',
  components: {
    AppLink
  },
  props: {
    item: {
      type: Object,
      required: true
    },
    isNest: {
      type: Boolean,
      default: false
    },
    basePath: {
      type: String,
      default: ''
    }
  },
  setup(props) {
    const route = useRoute()
    const onlyOneChild = ref({})

    function hasOneShowingChild(children = [], parent) {
      if (!children) {
        return false
      }

      const showingChildren = children.filter(item => {
        if (item.meta?.hidden) {
          return false
        } else {
          // 临时设置，当只有一个显示的子路由时使用
          onlyOneChild.value = item
          return true
        }
      })

      // 当只有一个子路由时，默认显示子路由
      if (showingChildren.length === 1) {
        return true
      }

      // 如果没有子路由，显示父路由
      if (showingChildren.length === 0) {
        onlyOneChild.value = { ...parent, path: '', noShowingChildren: true }
        return true
      }

      return false
    }

    function resolvePath(routePath) {
      if (isExternal(routePath)) {
        return routePath
      }
      if (isExternal(props.basePath)) {
        return props.basePath
      }
      return path.resolve(props.basePath, routePath)
    }

    function isExternal(path) {
      return /^(https?:|mailto:|tel:)/.test(path)
    }

    return {
      onlyOneChild,
      hasOneShowingChild,
      resolvePath
    }
  }
}
</script>

<style lang="scss" scoped>
.nest-menu {
  :deep(.el-sub-menu__title) {
    padding-left: 40px !important;
  }
}
</style> 