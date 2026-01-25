<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white rounded-lg shadow p-8 w-full max-w-md text-center">
      <h1 class="text-3xl font-bold text-gray-800 mb-4">网关服务管理系统</h1>
      <p class="text-gray-600 mb-8">多租户网关服务管理平台</p>

      <NuxtLink
        to="/login"
        class="inline-block bg-blue-500 text-white px-8 py-3 rounded-lg hover:bg-blue-600 transition"
      >
        立即登录
      </NuxtLink>
    </div>
  </div>
</template>

<script setup>
// 如果已登录，自动跳转到 dashboard
const authStore = useAuthStore()
onMounted(() => {
  if (authStore.isAuthenticated) {
    navigateTo('/dashboard')
  }
})
</script>
