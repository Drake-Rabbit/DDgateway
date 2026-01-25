<template>
  <div class="min-h-[60vh] flex items-center justify-center">
    <div class="bg-white rounded-lg shadow p-8 w-full max-w-md">
      <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Register</h1>

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-gray-700 text-sm font-bold mb-2">
            Username
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
            Email
          </label>
          <input
            v-model="formData.email"
            type="email"
            class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-bold mb-2">
            Password
          </label>
          <input
            v-model="formData.password"
            type="password"
            class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
            minlength="6"
          />
        </div>

        <div>
          <label class="block text-gray-700 text-sm font-bold mb-2">
            Tenant Code
          </label>
          <input
            v-model="formData.tenant_code"
            type="text"
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
          {{ loading ? 'Registering...' : 'Register' }}
        </button>
      </form>

      <p class="mt-4 text-center text-gray-600 text-sm">
        Already have an account?
        <NuxtLink to="/login" class="text-blue-500 hover:underline">
          Login
        </NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { useAuthStore } from '~/stores/auth'
import { useApi } from '~/composables/useApi'
import { reactive, ref } from 'vue'

const authStore = useAuthStore()
const api = useApi()

const formData = reactive({
  username: '',
  email: '',
  password: '',
  tenant_code: '',
})

const loading = ref(false)
const error = ref('')

const handleRegister = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await api.register(formData)

    if (response.data?.success) {
      authStore.setToken(response.data.data.token, response.data.data.user)
      navigateTo('/')
    } else {
      error.value = response.data?.error || 'Registration failed'
    }
  } catch (e) {
    error.value = e.data?.data?.error || 'An error occurred'
  } finally {
    loading.value = false
  }
}
</script>
