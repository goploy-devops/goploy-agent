<template>
  <div class="app-wrapper">
    <div>
      <Navbar />
    </div>
    <div class="main-container">
      <AppMain />
    </div>
  </div>
</template>

<script>
import Navbar from './Navbar.vue'
import AppMain from './AppMain.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Layout',
  components: {
    Navbar,
    AppMain,
  },
})
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';

.app-wrapper {
  @include clearfix;
  position: relative;
  height: 100%;
  width: 100%;
  &.mobile.openSidebar {
    position: fixed;
    top: 0;
  }
}

.main-container {
  transition: margin-left 0.28s;
  position: relative;
  background-image: url('@/assets/images/background.jpg');
  background-color: #fff;
  background-position: left bottom;
  background-size: cover;
  overflow: hidden;
}
</style>
