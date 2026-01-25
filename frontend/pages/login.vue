<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white rounded-lg shadow p-8 w-full max-w-md">
      <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">用户登录</h1>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-gray-700 text-sm font-bold mb-2">
            用户名
          </label>
          <input
            v-model="formData.username"
            type="text"
            class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-bold mb-2">
            密码
          </label>
          <input
            v-model="formData.password"
            type="password"
            class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div v-if="error" class="text-red-500 text-sm">
          {{ error }}
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition disabled:bg-blue-300"
        >
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>

      <p class="mt-4 text-center text-gray-600 text-sm">
        还没有账号？
        <NuxtLink to="/register" class="text-blue-500 hover:underline">
          立即注册
        </NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'
import { useMessage } from '~/composables/useMessage'
import { reactive, ref } from 'vue'

const authStore = useAuthStore()
const api = useApi()
const messageStore = useMessage()

const formData = reactive({
  username: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await api.login(formData)

    if (response.data?.success) {
      authStore.setToken(response.data.data.token, response.data.data.user)
      messageStore.success('登录成功！欢迎回来，' + response.data.data.user.username)
      navigateTo('/dashboard')
    } else {
      error.value = response.data?.error || '登录失败'
    }
  } catch (e) {
    error.value = e.response?.data?.error || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}
</script>
