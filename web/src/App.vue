<!--
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Mon Mar 4 16:45:10 2024 +0800
-->
<template>
  <div class="nav">
    <div class="logo">
      <!-- <img src="./assets/logo.png" alt=""> -->
      <span class="logo_title">Syscare management system</span>
    </div>
    <div class="route">
      <el-text class="route_item" :class="{ route_item_active: routeName === 'hostList' }" tag="b"
        @click=" route2page('hostList')">环境管理</el-text>
      <el-text class="route_item" :class="{ route_item_active: routeName === 'patchList' }" tag="b"
        @click="route2page('patchList')">热补丁管理</el-text>
    </div>
    <div class="top_right">
      <div class="triangle"></div>
      <div class="question">
        <el-icon style="color:#fff;" size="20" title="help" @click="hintHelp">
          <QuestionFilled />
        </el-icon>
        <span class="question_helper">有任何问题，联系yangzhao1@kylinos.cn</span>
      </div>
    </div>
  </div>
  <el-config-provider :locale="zhCn">
    <router-view v-slot="{ Component }">
      <transition name="fade">
        <component :is="Component"></component>
      </transition>
    </router-view>
  </el-config-provider>
</template>
<script setup lang="ts">
import locale from "element-plus/es/locale/lang/zh-cn";
import { QuestionFilled } from '@element-plus/icons-vue';
import { directTo } from './router';
import { useRoute } from "vue-router";
import { ref, watchEffect } from "vue";

const zhCn = ref(locale);
const routeName = ref('hostList');
const route2page = (pathName: string) => {
  routeName.value = pathName;
  directTo(pathName);
}
const route = useRoute();
watchEffect(() => {
  let currentPath = route.fullPath;
  routeName.value = currentPath === '/' ? '' : currentPath.split('/')[1];
})

// 帮助提示
const hintHelp = () => {

}
</script>

<style scoped lang="less">
.nav {
  width: 100%;
  height: 44px;
  display: flex;
  justify-content: space-between;
  border-bottom: 2px solid var(--main-color);

  .logo {
    width: 240px;
    height: 44px;
    display: flex;
    align-items: center;
    justify-content: space-evenly;

    img {
      width: 36px;
    }

    &_title {
      width: 176px;
      color: var(--main-color);
      font-size: 16px;
      font-weight: 600;
      text-shadow: 2px 2px 2px var(--el-color-primary-light-5);
    }
  }

  .route {
    width: 240px;
    display: flex;
    align-items: center;
    justify-content: space-around;

    &_item {
      padding-left: 4px;
      border: 3px solid #fff;

      &:hover,
      &_active {
        cursor: pointer;
        color: var(--main-color);
        border-left: 3px solid var(--main-color);
      }
    }
  }

  .top_right {
    width: 60%;
    height: 44px;
    display: flex;
    justify-content: flex-end;
    align-items: center;

    .triangle {
      width: 0;
      height: 0;
      border-left: 88px solid transparent;
      border-bottom: 44px solid var(--main-color);
    }

    .question {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      background-color: var(--main-color);
      width: 40%;
      padding-right: 20px;
      height: 44px;

      &_helper {
        color: #fff;
        font-size: 14px;
        padding: 0 6px;
      }
    }

  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .3s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}
</style>
