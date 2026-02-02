// stores/user.js
import { tr } from 'element-plus/es/locale/index.mjs'
import { defineStore } from 'pinia'
import { ref } from 'vue'



export const useUserStore = defineStore('user', () => {
  // 响应式状态
  const userinfo = ref({})

  // 同步 action（直接修改）
  function setUserInfo(newUserinfo) {
    userinfo.value = newUserinfo
  }

  // 返回要暴露的内容（必须显式 return）
  return {
    userinfo,
    setUserInfo,
  }

},
    {
    persist:{
      enabled: true,
       key: 'user-key',
    }
    }
)