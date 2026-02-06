<script setup>
import FHeader from "~/layout/components/FHeader.vue";
import FMenu from "~/layout/components/FMenu.vue";
import FTagList from "~/layout/components/FTagList.vue";

</script>

<template>
  <el-container>
    <!--1.总容器的上部分作为头部-->
    <el-header style="z-index: 999;">
      <FHeader />
    </el-header>

    <!--总容器的下面分成两容器,左右分割-->
    <el-container>
      <!--侧边栏-->
      <el-aside>
        <FMenu />
      </el-aside>
      <el-main>
        <!--导航栏-->
        <FTagList />

        <router-view v-slot="{ Component }">
          <keep-alive :max="5">
            <component :is="Component" />
          </keep-alive>
        </router-view>

      </el-main>

    </el-container>

  </el-container>

</template>

<style scoped></style>