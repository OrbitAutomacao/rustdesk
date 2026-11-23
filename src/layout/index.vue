<template>
  <el-config-provider :locale="appStore.setting.locale.value">
    <el-container class="app-layout" :class="{ 'mobile-menu-open': appStore.setting.mobileMenuOpen }" :style="{'--sideBarWidth': sideBarWidth}">
      <button v-if="appStore.setting.mobileMenuOpen" class="mobile-menu-backdrop" type="button" aria-label="Fechar menu" @click="appStore.closeMobileMenu()"></button>
      <el-aside :width="leftWidth" class="app-left" @click="onSidebarClick">
        <g-aside></g-aside>
      </el-aside>
      <el-container class="app-container ">
        <el-header class="app-header">
          <g-header></g-header>
        </el-header>
        <div class="header-tags">
          <tags></tags>
        </div>

        <el-main class="app-main">
          <router-view v-slot="{ Component }">
            <transition mode="out-in" name="el-fade-in-linear">
              <keep-alive :include="cachedTags">
                <component :is="Component"/>
              </keep-alive>
            </transition>
          </router-view>
        </el-main>
      </el-container>
    </el-container>
  </el-config-provider>
</template>

<script setup>
  import { useAppStore } from '@/store/app'
  import { useTagsStore } from '@/store/tags'
  import { ref, computed, watch } from 'vue'
  import { useRoute } from 'vue-router'
  import Tags from '@/layout/components/tags/index.vue'
  import GAside from '@/layout/components/aside.vue'
  import GHeader from '@/layout/components/header.vue'

  const appStore = useAppStore()
  const tagStore = useTagsStore()
  const route = useRoute()
  watch(() => route.fullPath, () => appStore.closeMobileMenu())
  const onSidebarClick = (event) => {
    if (event.target.closest('.el-menu-item')) appStore.closeMobileMenu()
  }
  const sideBarWidth = computed(() => appStore.setting.locale.sideBarWidth)
  const leftWidth = computed(() => appStore.setting.sideIsCollapse ? '64px' : 'var(--sideBarWidth)')

  const cachedTags = ref([])

  cachedTags.value = tagStore.cached
</script>

<style lang="scss" scoped>
.app-header {
  height: 68px;
  padding: 0 26px;
  display: flex;
  align-items: center;
  color: #17334b;
  background: #fff;
  border-bottom: 1px solid #e3edf3;
}

.header-tags {
  min-height: 48px;
  padding: 8px 24px;
  display: flex;
  align-items: center;
  gap: 8px;
  background: #fff;
  border-bottom: 1px solid #e3edf3;
}

.app-left {
  transition: width .25s ease;
  background: #0e243c;
  border-right: 1px solid #14344f;
}

.app-container {
  min-width: 0;
  min-height: 100vh;
  background: var(--orbit-canvas);
}

.app-main { padding: 24px; }

@media (max-width: 700px) {
  .app-layout { width: 100%; min-width: 0; }
  .app-left {
    position: fixed;
    z-index: 101;
    top: 0;
    bottom: 0;
    left: 0;
    width: min(280px, calc(100vw - 48px)) !important;
    transform: translateX(-100%);
    transition: transform .25s ease;
    box-shadow: 12px 0 36px rgba(5, 20, 35, .18);
  }
  .mobile-menu-open .app-left { transform: translateX(0); }
  .mobile-menu-backdrop {
    position: fixed;
    z-index: 100;
    inset: 0;
    border: 0;
    background: rgba(3, 18, 34, .58);
  }
  .app-container { width: 100%; max-width: 100%; }
  .app-main { min-width: 0; }
  .app-header { height: 60px; }
  .header-tags { min-width: 0; }
}

@media (max-width: 900px) {
  .app-header { padding: 0 16px; }
  .header-tags { padding: 8px 14px; overflow-x: auto; }
  .app-main { padding: 14px; }
}
</style>


