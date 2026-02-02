<script setup>
import { useUserStore } from '~/store/user'
import { ArrowDown } from '@element-plus/icons-vue'
const userstore = useUserStore()
import { useRouter } from 'vue-router'
import { logout } from '~/api/login.js'
import { ElMessage } from 'element-plus'
import { showModal } from '~/composable/util.js'
import { removeToken } from '~/composable/auth.js'

const router = useRouter()

const handleCommand = (command) => {
  if (command === 'profile') {
    // 处理个人设置
    console.log('个人设置')
  } else if (command === 'logout') {
    // 处理退出登陆
    handlelogout()
  }
}

function handlelogout() {
  showModal('确定要退出登陆吗？', 'warning', '退出登陆').then(res => {
    logout()
      .finally(()=>{
        removeToken() //移除cookie的token
        userstore.setUserInfo({}) //清除pina的用户信息状态
        router.push('/login')
        ElMessage.success('退出登陆成功！')
    })
  })
}


</script>

<template>
  <header class="f-header">
    <!-- 左侧 Logo -->
    <div class="text-2xl font-bold">GateWay Service</div>

    <!-- 中间占位（可选） -->
    <div class="flex-1"></div>

    <!-- 右侧用户菜单 -->
    <div class="flex items-center">
      <el-dropdown @command="handleCommand">
        <span class="el-dropdown-link flex items-center cursor-pointer">
          <el-avatar class="mr-2"> {{ userstore.userinfo.username?.[0] || 'U' }} </el-avatar>
          <span class="text-white mr-2">{{userstore.userinfo.username  }}</span>
          <el-icon class="el-icon--right ml-1">
            <ArrowDown />
          </el-icon>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人设置</el-dropdown-item>
            <el-dropdown-item command="logout">退出登陆</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>


<style scoped>
.f-header{
  @apply  flex bg-sky-500 text-white p-6  fixed top-0 left-0 right-0 items-center;
  height: 64px ;
}
</style>