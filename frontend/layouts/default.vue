<template>
  <div class="min-h-screen bg-gray-100">
    <template v-if="isAuthenticated">
      <!-- 主布局：左侧导航 + 右侧内容 -->
      <div class="flex min-h-screen">
        <!-- 左侧导航栏 -->
        <aside class="w-64 bg-gray-900 text-white flex flex-col">
          <!-- Logo -->
          <div class="p-6 border-b border-gray-700">
            <h1 class="text-xl font-bold text-blue-400">Gateway Service</h1>
          </div>

          <!-- 导航菜单 -->
          <nav class="flex-1 p-4">
            <ul class="space-y-2">
              <li>
                <NuxtLink
                  to="/dashboard"
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-800 transition"
                  :class="{ 'bg-blue-600': isActive('/dashboard') }"
                >
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
                  </svg>
                  首页
                </NuxtLink>
              </li>
              <li>
                <NuxtLink
                  to="/services"
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-800 transition"
                  :class="{ 'bg-blue-600': isActive('/services') }"
                >
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
                  </svg>
                  服务管理
                </NuxtLink>
              </li>
              <li>
                <NuxtLink
                  to="/tenants"
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-800 transition"
                  :class="{ 'bg-blue-600': isActive('/tenants') }"
                >
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                  租户管理
                </NuxtLink>
              </li>
              <li>
                <NuxtLink
                  to="/users"
                  class="flex items-center px-4 py-3 rounded-lg hover:bg-gray-800 transition"
                  :class="{ 'bg-blue-600': isActive('/users') }"
                >
                  <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                  </svg>
                  用户管理
                </NuxtLink>
              </li>
            </ul>
          </nav>

          <!-- 用户信息 + 退出 -->
          <div class="p-4 border-t border-gray-700">
            <div class="flex items-center mb-3">
              <div class="w-10 h-10 bg-blue-500 rounded-full flex items-center justify-center">
                <span class="text-white font-bold">{{ userInitial }}</span>
              </div>
              <div class="ml-3 flex-1">
                <p class="text-sm font-medium">{{ authStore.user?.username }}</p>
                <p class="text-xs text-gray-400">{{ authStore.user?.role }}</p>
              </div>
            </div>
            <button
              @click="logout"
              class="w-full flex items-center justify-center px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              退出登录
            </button>
          </div>
        </aside>

        <!-- 右侧内容区 -->
        <main class="flex-1 overflow-auto">
          <slot />
        </main>
      </div>
    </template>

    <!-- 未登录时显示顶部导航 -->
    <template v-else>
      <nav class="bg-white shadow-lg">
        <div class="max-w-7xl mx-auto px-4">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <NuxtLink to="/" class="text-xl font-bold text-blue-600">
                Gateway Service
              </NuxtLink>
            </div>
            <div class="flex items-center space-x-4">
              <NuxtLink
                to="/login"
                class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
              >
                Login
              </NuxtLink>
            </div>
          </div>
        </div>
      </nav>
      <main class="max-w-7xl mx-auto px-4 py-6">
        <slot />
      </main>
    </template>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useMessage } from '~/composables/useMessage'
import { computed } from 'vue'

const authStore = useAuthStore()
const messageStore = useMessage()
const route = useRoute()

const isAuthenticated = computed(() => !!authStore.token)

const userInitial = computed(() => {
  return authStore.user?.username?.charAt(0).toUpperCase() || 'U'
})

const isActive = (path) => {
  return route.path === path
}

const logout = () => {
  authStore.logout()
  messageStore.info('退出登录成功')
  navigateTo('/login')
}
</script>
