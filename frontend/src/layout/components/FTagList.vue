<template>
  <div class="f-tag-list">
    <el-tabs v-model="activeTab" type="card" class="demo-tabs" closable @tab-remove="removeTab" @tab-change="changeTab"
      style="min-width: 100px;">
      <el-tab-pane class="bg-white-100" v-for="item in tabList" :key="item.path" :label="item.title" :name="item.path"
      :closable="item.path != '/'">
      </el-tab-pane>
    </el-tabs>
  </div>
  <div style="height: 8px;"></div>
</template>


<script setup>
import { ref } from 'vue'
import { onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router'
const route = useRoute()
const router = useRouter()
import { useCookies } from "@vueuse/integrations/useCookies";
const cookies = useCookies()
import { onMounted } from 'vue'

let tabIndex = 2
const activeTab = ref(route.path)
const tabList = ref([
  {
    title: '仪表盘',
    path: '/',
  },
])

// 添加标签
function addTab(tab) {
  let noTab = tabList.value.findIndex(item => item.path === tab.path) == -1
  if (noTab) {
    tabList.value.push(tab)
  }
  cookies.set('tabList', tabList.value)
}

// 监听路由变化
onBeforeRouteUpdate((to, from) => {
  activeTab.value = to.path
  // console.log(to.path)
  addTab({
    title: to.meta.title,
    path: to.path,
  })
})

// 切换标签
function changeTab(tab) {
  activeTab.value = tab //更新高亮的path
  router.push(tab) // url切换路由
}
//关闭标签
function removeTab(tabPath) {
  let tabs = tabList.value
  let a = activeTab.value
  if (a == tabPath) {
    tabs.forEach((item, index) => {
      if (item.path === tabPath) {
        const nextTab = tabs[index + 1] || tabs[index - 1]
        if (nextTab) {
          a = nextTab.path
        }else{
          a = '/'
        }
      }
    })
  }
  activeTab.value = a
  router.push(a);
  tabList.value = tabList.value.filter(item => item.path !== tabPath)
  cookies.set('tabList', tabList.value)
}

//初始化标签导航列表
function initTabList() {
  let tbs = cookies.get('tabList')
  if (tbs) {
    tabList.value = tbs
  }
}
onMounted(() => {
  initTabList()
})

</script>


<style scoped>
.f-tag-list {
  margin-bottom: 20px;
  right: 0px;
  height: 42px;
  @apply bg-light-400 items-center;
}

:deep(.el-tabs__header) {
  border: 0 !important;
  @apply mb-0;
}

:deep(.el-tabs__nav) {
  border: 0 !important;
}

:deep(.el-tabs__item) {
  border: 0 !important;
  height: 32px;
  line-height: 32px;
  @apply bg-white mx-1 rounded;
}

:deep(.el-tabs__next),
:deep(.el-tabs__prev) {
  @apply bg-white;
  height: 32px;
  line-height: 32px;
}

:deep(.is-disabled) {
  cursor: not-allowed;
  @apply text-gray-400;
}

:deep(.el-tabs__content) {
  display: none !important;
  border: none !important;
  padding: 0 !important;
  margin: 0 !important;
  background: transparent !important;
  height: 0 !important;
  overflow: hidden;
}
</style>