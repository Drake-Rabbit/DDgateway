<script setup>

</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-50">
    <img
      src="../assets/Gopher-404.png"
      alt="Go Gopher"
      class="w-64 h-64 mb-8"
    />
    <h1 class="text-6xl font-bold text-gray-800 mb-4">404</h1>
    <p class="text-xl text-gray-600 mb-8">页面走丢了，请检查URL或返回首页</p>
    <router-link
      to="/"
      class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
    >
      返回首页
    </router-link>
  </div>
</template>

<style scoped>

</style>