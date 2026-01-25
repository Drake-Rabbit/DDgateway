<template>
  <NuxtLayout>
    <NuxtPage />

    <!-- 全局消息提示组件 -->
    <div class="fixed top-4 right-4 z-50 space-y-2">
      <div
        v-for="msg in messageStore.messages"
        :key="msg.id"
        :class="[
          'px-4 py-3 rounded-lg shadow-lg flex items-center justify-between transition-all transform animate-slide-in',
          {
            'bg-green-100 text-green-800 border border-green-200': msg.type === 'success',
            'bg-red-100 text-red-800 border border-red-200': msg.type === 'error',
            'bg-blue-100 text-blue-800 border border-blue-200': msg.type === 'info',
            'bg-yellow-100 text-yellow-800 border border-yellow-200': msg.type === 'warning',
          }
        ]"
      >
        <span class="flex-1">{{ msg.text }}</span>
        <button
          @click="messageStore.removeMessage(msg.id)"
          class="ml-2 text-current hover:opacity-70 transition-opacity"
        >
          &times;
        </button>
      </div>
    </div>



  </NuxtLayout>
</template>

<script setup>
import { onMounted, provide } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useMessage } from '~/composables/useMessage'

// 全局初始化：从localStorage加载登录状态
onMounted(() => {
  useAuthStore().loadFromStorage()
})

// 提供消息服务
const messageStore = useMessage()
provide('messageStore', messageStore)
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background-color: #f5f5f5;
}

/* 动画效果 */
.animate-slide-in {
  animation: slideInRight 0.3s ease-out;
}

@keyframes slideInRight {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
</style>
