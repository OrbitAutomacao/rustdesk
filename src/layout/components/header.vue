<template>
  <el-icon class="ex-icon" @click="expandOrFoldSlider">
    <el-icon-expand v-if="isMobile ? !setting.mobileMenuOpen : setting.sideIsCollapse"></el-icon-expand>
    <el-icon-fold v-else></el-icon-fold>
  </el-icon>
  <div class="header-logo">
    <img :src="setting.logo" alt="" class="logo">
    <div class="title">{{setting.title}}</div>
  </div>
  <Setting></Setting>
</template>

<script>
  import { defineComponent, computed, onMounted, onUnmounted, ref } from 'vue'
  import HeaderMenu from '@/layout/components/menu/index.vue'
  import Setting from '@/layout/components/setting/index.vue'
  import { useAppStore } from '@/store/app'
  import GTags from '@/layout/components/tags/index.vue'

  export default defineComponent({
    name: 'LayerHeader',
    created () {
    },
    components: { HeaderMenu, Setting, GTags },
    watch: {},
    setup (props) {
      const appStore = useAppStore()
      const setting = computed(() => appStore.setting)
      const isMobile = ref(false)
      const updateMobile = () => { isMobile.value = window.matchMedia('(max-width: 700px)').matches }
      onMounted(() => { updateMobile(); window.addEventListener('resize', updateMobile) })
      onUnmounted(() => window.removeEventListener('resize', updateMobile))
      const expandOrFoldSlider = () => {
        if (isMobile.value) appStore.toggleMobileMenu()
        else appStore.sideCollapse()
      }
      return {
        setting,
        isMobile,
        expandOrFoldSlider,
      }
    },

  })
</script>

<style scoped lang="scss">
.ex-icon {
  display: flex;
  align-items: center;
  height: 100%;
  margin-right: 18px;
  font-size: 20px;
  cursor: pointer;
  color: #5a7489;
}

.header-logo {
  display: flex;
  align-items: center;
  min-width: 0;

  .logo {
    display: block;
    width: 108px;
    max-height: 40px;
    object-fit: contain;
  }

  .title {
    margin-left: 16px;
    padding-left: 16px;
    border-left: 1px solid #dce7ef;
    color: #1c3e59;
    font-size: 16px;
    font-weight: 700;
    white-space: nowrap;
  }
}

@media (max-width: 700px) {
  .ex-icon { flex: 0 0 auto; margin-right: 10px; min-width: 36px; justify-content: center; }
  .header-logo { overflow: hidden; }
  .header-logo .title { overflow: hidden; text-overflow: ellipsis; }
  .header-logo .logo { flex: 0 0 auto; }
  .header-logo .logo { width: 78px; }
  .header-logo .title { font-size: 13px; margin-left: 10px; padding-left: 10px; }
}
@media (max-width: 480px) {
  .header-logo .title { display: none; }
}
</style>
<style lang="scss">

</style>
