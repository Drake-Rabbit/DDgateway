import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()
  
  // 不需要认证的页面
  const publicPages = ['/login', '/register']
  
  // 如果用户未登录且尝试访问受保护页面，重定向到登录页
  if (!authStore.token && !publicPages.includes(to.path)) {
    return navigateTo('/login')
  }
  
  // 如果用户已登录且尝试访问登录页或注册页，重定向到首页
  if (authStore.token && publicPages.includes(to.path)) {
    return navigateTo('/')
  }
})