<template>
  <div class="p-6">
    <!-- 欢迎横幅 -->
    <div class="bg-gradient-to-r from-blue-500 to-blue-600 rounded-lg shadow-lg p-6 mb-6 text-white">
      <h1 class="text-3xl font-bold mb-2">欢迎使用网关服务管理系统</h1>
      <p class="text-blue-100">{{ authStore.user?.username }}，您好！</p>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">服务总数</p>
            <p class="text-3xl font-bold text-gray-800 mt-2">{{ stats.serviceCount }}</p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">租户总数</p>
            <p class="text-3xl font-bold text-gray-800 mt-2">{{ stats.tenantCount }}</p>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">用户总数</p>
            <p class="text-3xl font-bold text-gray-800 mt-2">{{ stats.userCount }}</p>
          </div>
          <div class="w-12 h-12 bg-purple-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-gray-500 text-sm">HTTP 服务</p>
            <p class="text-3xl font-bold text-gray-800 mt-2">{{ stats.httpServiceCount }}</p>
          </div>
          <div class="w-12 h-12 bg-orange-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- 快捷操作 -->
    <div class="bg-white rounded-lg shadow p-6">
      <h2 class="text-xl font-bold text-gray-800 mb-4">快捷操作</h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <NuxtLink
          to="/services"
          class="flex items-center p-4 border rounded-lg hover:bg-blue-50 hover:border-blue-300 transition"
        >
          <div class="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center mr-4">
            <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div>
            <p class="font-medium text-gray-800">添加服务</p>
            <p class="text-sm text-gray-500">创建新的服务配置</p>
          </div>
        </NuxtLink>

        <NuxtLink
          to="/tenants"
          class="flex items-center p-4 border rounded-lg hover:bg-green-50 hover:border-green-300 transition"
        >
          <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center mr-4">
            <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div>
            <p class="font-medium text-gray-800">添加租户</p>
            <p class="text-sm text-gray-500">创建新的租户</p>
          </div>
        </NuxtLink>

        <NuxtLink
          to="/users"
          class="flex items-center p-4 border rounded-lg hover:bg-purple-50 hover:border-purple-300 transition"
        >
          <div class="w-10 h-10 bg-purple-100 rounded-full flex items-center justify-center mr-4">
            <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </div>
          <div>
            <p class="font-medium text-gray-800">添加用户</p>
            <p class="text-sm text-gray-500">创建新的用户</p>
          </div>
        </NuxtLink>
      </div>
    </div>

    <!-- 最近活动 -->
    <div class="bg-white rounded-lg shadow p-6 mt-6">
      <h2 class="text-xl font-bold text-gray-800 mb-4">系统信息</h2>
      <div class="space-y-3">
        <div class="flex items-center justify-between py-2 border-b">
          <span class="text-gray-600">当前用户</span>
          <span class="font-medium text-gray-800">{{ authStore.user?.username }}</span>
        </div>
        <div class="flex items-center justify-between py-2 border-b">
          <span class="text-gray-600">用户角色</span>
          <span class="px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800">{{ authStore.user?.role }}</span>
        </div>
        <div class="flex items-center justify-between py-2">
          <span class="text-gray-600">最后登录</span>
          <span class="font-medium text-gray-800">{{ lastLoginTime }}</span>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'
import { useMessage } from '~/composables/useMessage'

const authStore = useAuthStore()
const api = useApi()
const messageStore = useMessage()

const stats = ref({
  serviceCount: 0,
  tenantCount: 0,
  userCount: 0,
  httpServiceCount: 0
})

const lastLoginTime = computed(() => {
  if (!authStore.user?.last_login_at) return '首次登录'
  const date = new Date(authStore.user.last_login_at)
  return date.toLocaleString('zh-CN')
})

// 获取统计数据
const fetchStats = async () => {
  try {
    const [servicesRes, tenantsRes, usersRes] = await Promise.all([
      api.getServices(),
      api.getTenants(),
      api.getUsers()
    ])

    if (servicesRes.data?.success) {
      const services = servicesRes.data.data || []
      stats.value.serviceCount = services.length
      stats.value.httpServiceCount = services.filter(s => s.load_type === 0).length
    }

    if (tenantsRes.data?.success) {
      stats.value.tenantCount = (tenantsRes.data.data || []).length
    }

    if (usersRes.data?.success) {
      stats.value.userCount = (usersRes.data.data || []).length
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

onMounted(() => {
  fetchStats()
})
</script>
