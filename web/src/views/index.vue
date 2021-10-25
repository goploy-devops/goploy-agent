<template>
  <div class="login-container">1</div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'Index',
  data() {
    return {}
  },
  methods: {},
})
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #fff;
$light_gray: #fff;
$cursor: #2f2f2f;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 42px;
    width: 85%;

    input {
      background: $bg;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: #2f2f2f;
      height: 40px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
    }
  }
  .el-form-item__error {
    padding-top: 4px;
  }
  .el-form-item {
    background: #fff;
  }
}
</style>

<style lang="scss" scoped>
$bg: #fff;
$dark_gray: #889aa4;
$light_gray: #eee;

.login-container {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  background-image: url('@/assets/images/background.jpg');
  background-position: left bottom;
  background-size: cover;
  width: 100%;
  overflow: hidden;
  .login-mark {
    margin-top: 58px;
  }

  .login-form {
    flex: 1;
    align-self: center;
    bottom: 0;
    width: 520px;
    padding: 0 35px;
    overflow: hidden;
    &-input {
      border: 1px solid rgba(229, 230, 231, 1);
      border-radius: 35px;
    }
    &-btn {
      margin-top: 15px;
      height: 40px;
      border-radius: 39px;
      font-size: 16px;
      background-color: #2580fa;
    }
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 20px;
    color: #c1c9d7;
    font-size: 16px;
    vertical-align: middle;
    width: 40px;
    display: inline-block;
  }

  .title-container {
    position: relative;
    padding-left: 5px;
    text-align: center;
    .title {
      margin: 20px 0;
      display: inline-block;
      font-size: 20px;
      font-family: PingFang SC;
      color: #586069;
    }
  }

  .show-pwd {
    position: absolute;
    right: 20px;
    top: 10px;
    font-size: 16px;
    color: #c1c9d7;
    cursor: pointer;
    user-select: none;
  }
}
</style>
