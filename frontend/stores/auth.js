import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    user: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
  },

  actions: {
    setToken(token, user) {
      this.token = token
      this.user = user
      localStorage.setItem('auth_token', token)
      localStorage.setItem('auth_user', JSON.stringify(user))
    },

    loadFromStorage() {
      try {
        const token = localStorage.getItem('auth_token')
        const userStr = localStorage.getItem('auth_user')
        if (token && userStr) {
          this.token = token
          this.user = JSON.parse(userStr)
        }
      } catch (error) {
        console.warn('本地存储解析失败', error)
        this.logout()
      }
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
    },
  },

  hydrate(state) {
    const token = localStorage.getItem('auth_token')
    const userStr = localStorage.getItem('auth_user')
    if (token && userStr) {
      try {
        state.token = token
        state.user = JSON.parse(userStr)
      } catch (error) {
        console.warn('Hydrate解析失败', error)
      }
    }
  },
})