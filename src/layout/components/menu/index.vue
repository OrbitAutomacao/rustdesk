<template>
  <el-menu
          class="menus"
          :collapse="isCollapse"
          :default-active="activeIndex"
          background-color="#0e243c"
          text-color="#bfd0df"
          active-text-color="#55d2f4"
          router
  >
    <menu-item v-for="(route,index) in routes" :key="route.name" :route="route"></menu-item>
  </el-menu>
</template>

<script>
  import { defineComponent, ref, onMounted, watch, computed } from 'vue'
  import { useRouteStore } from '@/store/router'
  import MenuItem from '@/layout/components/menu/item.vue'
  import { useRoute } from 'vue-router'
  import { useAppStore } from '@/store/app'
  import { useMediaQuery } from '@vueuse/core'

  export default defineComponent({
    name: 'Menu',
    created () {
    },
    components: { MenuItem },
    setup () {
      const routes = ref([])
      const route = useRoute()
      const app = useAppStore()
      const isMobile = useMediaQuery('(max-width: 700px)')
      const isCollapse = computed(() => !isMobile.value && app.setting.sideIsCollapse)
      const activeIndex = computed(() => route.name)

      routes.value = useRouteStore().routes
      return {
        routes,
        activeIndex,
        isCollapse,
      }
    },

  })
</script>

<style lang="scss" scoped>
.menus {
  min-height: calc(100vh - 24px);
  border-right: none;
  &:not(.el-menu--collapse) { width: 100%; }

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    min-height: 44px;
    height: auto;
    margin: 3px 0;
    border-radius: 9px;
    font-weight: 500;
    line-height: 1.25;
  }

  :deep(.el-menu-item:hover),
  :deep(.el-sub-menu__title:hover) { background: #193b55; }

  :deep(.el-menu-item.is-active) {
    background: #183e58;
    color: #55d2f4;
    font-weight: 700;
    box-shadow: inset 3px 0 #22bce8;
  }
}
</style>
<style>
</style>
