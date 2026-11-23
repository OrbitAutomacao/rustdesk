<template>
  <el-scrollbar class="scroll-sidebar" height="100vh">
    <menus></menus>
  </el-scrollbar>
</template>
<script>
  import Menus from '@/layout/components/menu/index.vue'
  import { defineComponent, ref, onMounted } from 'vue'

  export default defineComponent({
    name: 'GAside',
    components: { Menus },
  })
</script>

<style scoped lang="scss">
.scroll-sidebar {
  background: #0e243c;
  padding: 12px 8px;
  border-right: 1px solid #173b58;
}
</style>
