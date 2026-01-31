<template>
  <div class="login-container">
    <div class="w-full max-w-sm">
      <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
        <!-- Logo Header -->
        <div class="bg-gradient-to-r from-blue-600 to-indigo-600 p-5 text-center">
          <div class="flex justify-center mb-3">
            <div class="w-14 h-14 bg-white rounded-full flex items-center justify-center">
              <svg class="w-8 h-8 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                <path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838L7.667 9.088l1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zM9.3 16.573A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z"/>
              </svg>
            </div>
          </div>
          <h2 class="text-xl font-bold text-white">Gateway Service</h2>
          <p class="text-blue-100 text-sm">登录您的账户</p>
        </div>

        <!-- Login Form -->
        <div class="p-6">
          <el-form @submit.prevent="handleLogin" :model="loginForm" class="w-full">
            <!-- Username Field -->
            <el-form-item class="sm text-center">
              <el-input
                v-model="loginForm.username"
                placeholder="请输入用户名"
                prefix-icon="User"
                clearable
                class="w-64 mx-auto"
              />
            </el-form-item>

            <!-- Password Field -->
            <el-form-item class="text-center">
              <el-input
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码"
                prefix-icon="Lock"
                clearable
                show-password
                class="w-64 mx-auto"
              />
            </el-form-item>

            <!-- Forgot Password -->
            <div class="text-center mb-4">
              <el-link type="primary" @click="handleForgotPassword">忘记密码？</el-link>
            </div>

            <!-- Login Button -->
            <el-form-item class="text-center ">
              <el-button
                type="primary"
                class="w-full center"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>


          </el-form>

          <!-- Divider -->
          <div class="my-4">
            <div class="flex items-center">
              <div class="flex-1 border-t border-gray-300"></div>
              <span class="flex-shrink mx-3 text-gray-400 text-sm">或</span>
              <div class="flex-1 border-t border-gray-300"></div>
            </div>
          </div>

          <!-- Social Login -->
          <div class="flex justify-center space-x-4">
            <el-button
              type="text"
              class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-gray-200 transition-colors"
            >
              <svg class="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 24 24">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
            </el-button>
            <el-button
              type="text"
              class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-gray-200 transition-colors"
            >
              <svg class="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.09.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.577.688.482A10.019 10.019 0 0022 12c0-5.523-4.477-10-10-10z"/>
              </svg>
            </el-button>
            <el-button
              type="text"
              class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-gray-200 transition-colors"
            >
              <svg class="w-5 h-5 text-gray-600" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zM5.838 12a6.162 6.162 0 1112.324 0 6.162 6.162 0 01-12.324 0zM12 16a4 4 0 110-8 4 4 0 010 8zm4.965-10.405a1.44 1.44 0 112.881.001 1.44 1.44 0 01-2.881-.001z"/>
              </svg>
            </el-button>
          </div>
        </div>

        <!-- Footer -->
        <div class="bg-gray-50 px-6 py-4 text-center">
          <p class="text-sm text-gray-600">
            还没有账户？ <el-link type="primary" href="/register">立即注册</el-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const loginForm = ref({
  username: '',
  password: ''
})
const showPassword = ref(false)

const handleLogin = () => {
  // 这里可以添加实际的登录逻辑
  console.log('Login attempt:', loginForm.value)

  // 模拟登录成功
  if (loginForm.value.username && loginForm.value.password) {
    ElMessage.success('登录成功！')
    router.push('/')
  } else {
    ElMessage.error('请填写用户名和密码')
  }
}

const handleForgotPassword = () => {
  ElMessage.info('忘记密码功能开发中...')
}
</script>

<style scoped>
.login-container{
  @apply min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-0;
}
</style>