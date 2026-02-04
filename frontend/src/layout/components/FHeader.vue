<script setup>
import { useUserStore } from '~/store/user'
import { ArrowDown } from '@element-plus/icons-vue'
const userstore = useUserStore()
import { useRouter } from 'vue-router'
import { logout, updatePassword } from '~/api/login.js'
import { ElMessage } from 'element-plus'
import { showModal } from '~/composable/util.js'
import { removeToken } from '~/composable/auth.js'

import { ref } from 'vue'
const upatepwd_drawer = ref(false)
const formRef = ref(null)

const router = useRouter()
const form = ref({
  oldPwd: '',
  newPwd: '',
  confirmPwd: ''
})
const rules = ref({
  oldPwd: [
    { required: true, message: '请输入旧密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  newPwd: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPwd: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
})

const handleCommand = (command) => {
  if (command === 'profile') {
    // 处理个人设置
    console.log('个人设置')
  } else if (command === 'logout') {
    // 处理退出登陆
    handlelogout()
  }else if (command === 'changepwd') {
    // 处理修改密码
    console.log('修改密码')
    upatepwd_drawer.value = true
  }

  
}

function handlelogout() {
  showModal('确定要退出登陆吗？', 'warning', '退出登陆').then(res => {
    logout()
      .finally(()=>{
        removeToken() //移除cookie的token
        userstore.removeUserInfo() //清除pina的用户信息状态
        router.push('/login')
        ElMessage.success('退出登陆成功!')
    })
  })
}
// 处理修改密码的提交
async function handleSubmit() {
 // 1. 先让 Element Plus 校验基础规则（必填、长度）
  formRef.value.validate((valid) => {
    if (!valid){
      return false
    }
    // 2. 再校验「两次新密码是否一致」
    if (form.value.newPwd !== form.value.confirmPwd) {
      ElMessage.error('两次输入的新密码不一致')
      return
    }

    // 3. 一致才提交到后端
    updatePassword(form.value)
      .then(res => {
        if(res.data.success == true){
          ElMessage.success('密码修改成功！,请重新登录')
          upatepwd_drawer.value = false
          // 清空表单
          form.value = { oldPwd: '', newPwd: '', confirmPwd: '' }
          formRef.value.clearValidate()
          //直接退出登陆
          removeToken() //移除cookie的token
          userstore.removeUserInfo() //清除pina的用户信息状态
          router.push('/login')
          }
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
            <el-dropdown-item command="changepwd">修改密码</el-dropdown-item>
            <el-dropdown-item command="logout">退出登陆</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
<!--  修改密码的右边弹出-->
  <el-drawer v-model="upatepwd_drawer" title="修改密码"
  size="45%" :close-on-click-modal="false"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" class="demo-ruleForm">
      <el-form-item label="旧密码" prop="oldPwd">
        <el-input v-model="form.oldPwd" type="password" autocomplete="off" />
      </el-form-item>
      <el-form-item label="新密码" prop="newPwd">
        <el-input v-model="form.newPwd" type="password" autocomplete="off" />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPwd">
        <el-input v-model="form.confirmPwd" type="password" autocomplete="off" />
      </el-form-item>
            <!-- 更新密码的按钮 -->
      <el-form-item class="text-center">
        <el-button type="primary" @click="handleSubmit">提交</el-button>
        <el-button type="default" @click="upatepwd_drawer = false">取消</el-button>
      </el-form-item>
    </el-form>


  </el-drawer>


</template>


<style scoped>
.f-header{
  @apply  flex bg-sky-500 text-white p-6  fixed top-0 left-0 right-0 items-center;
  height: 64px ;
}
</style>