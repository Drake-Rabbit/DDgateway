<template>
  <div class="f-menu">
  <el-menu  @select="handleSelect" :default-active="defaultActive"
  unique-opened>
    <template v-for="item in AsideMenus" :key="item.name">
      <!-- 情况1：如果有子菜单（children）→ 用 el-sub-menu -->
      <el-sub-menu v-if="item.children && item.children.length" :index="item.frontpath">
        <template #title>
          <!-- 图标：动态使用 item.icon 对应的组件 -->
          <el-icon>
            <component :is="item.icon" />
          </el-icon>
          <!-- 标题文字 -->
          <span>{{ item.name }}</span>
        </template>

        <!-- 遍历子菜单 -->
        <el-menu-item v-for="child in item.children" :key="child.name" :index="child.frontpath">
          <el-icon>
            <component :is="child.icon" />
          </el-icon>
          <span>{{ child.name }}</span>
        </el-menu-item>
      </el-sub-menu>

      <!-- 情况2：如果没有子菜单 → 用 el-menu-item -->
      <el-menu-item v-else :index="item.frontpath">
        <el-icon>
          <component :is="item.icon" />
        </el-icon>
        <span>{{ item.name }}</span>
      </el-menu-item>

    </template>


  </el-menu>
  </div>
</template>

<script setup>
const AsideMenus = [
  { "name": "仪表盘", "icon": "Odometer", frontpath: "/" },
  { "name": "服务列表", "icon": "MagicStick", frontpath: "/service-list" },
  {
    "name": "AI功能", "icon": "place",index:"/ai",
    "children": [
      { "name": "AI分析", "icon": "pear", frontpath: "/ai-analysis" },
      { "name": "AI模型管理", "icon": "apple", frontpath: "/ai-modelList" },
    ]
  },
  { "name": "系统设置", "icon": "setting", frontpath: "/system-setting" },
]


import { ref } from 'vue'
import { useRouter ,useRoute} from 'vue-router'
const router = useRouter()
// 定义默认选中项
const route = useRoute()
const defaultActive =ref(route.path)


// 处理菜单点击事件
const handleSelect = (index) => {router.push(index)}
</script>

<style scoped>
.f-menu {
  width: 250px;
  top: 64px;
  bottom: 0;
  left: 0;
  overflow: auto;
  @apply shadow-md fixed;
  border: 0px;
}
</style>