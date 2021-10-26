<template>
  <div class="navbar">
    <div class="logo-container">
      <img :src="logo" class="logo" />
    </div>
    <div class="title">Goploy-Agent</div>
    <el-menu class="menu" mode="horizontal" :router="true">
      <el-menu-item
        v-for="item in navbarRoutes.children"
        :key="item.path"
        :index="item.path"
      >
        <i :class="item.meta.icon"></i>
        <span>{{ item.name }}</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script>
import logo from '@/assets/images/logo.png'
import { defineComponent } from 'vue'
import { navbarRoutes } from '@/router'

export default defineComponent({
  data() {
    return {
      logo: logo,
      navbarRoutes: navbarRoutes,
    }
  },
})
</script>

<style lang="scss" scoped>
.navbar {
  height: 61px;
  display: flex;
  .logo-container {
    padding: 12px 0 12px 15px;
    border-bottom: solid 1px #e6e6e6;
  }
  .logo {
    width: 40px;
    height: 40px;
  }
  .title {
    font-weight: 600;
    font-size: 16px;
    margin-top: 23px;
    padding: 0 15px;
    color: #606266;
    border-bottom: solid 1px #e6e6e6;
  }
  .menu {
    flex: 1;
  }
}
</style>
