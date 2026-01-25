<template>
  <div class="p-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800">服务管理</h1>
      <p class="text-gray-600 mt-1">管理系统中所有的服务配置</p>
    </div>

    <!-- 操作栏 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex justify-between items-center">
        <div class="flex items-center space-x-4">
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索服务名称..."
              class="pl-10 pr-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-64"
            />
            <svg class="w-5 h-5 absolute left-3 top-2.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <select v-model="loadTypeFilter" class="px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="">全部类型</option>
            <option value="0">HTTP</option>
            <option value="1">TCP</option>
            <option value="2">GRPC</option>
          </select>
        </div>
        <button
          @click="showCreateModal = true"
          class="bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-600 transition flex items-center"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          添加服务
        </button>
      </div>
    </div>

    <!-- 服务列表表格 -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">服务名称</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">负载类型</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">服务描述</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-if="loading">
            <td colspan="6" class="px-6 py-4 text-center text-gray-500">
              加载中...
            </td>
          </tr>
          <tr v-else-if="filteredServices.length === 0">
            <td colspan="6" class="px-6 py-4 text-center text-gray-500">
              暂无数据
            </td>
          </tr>
          <tr v-else v-for="service in paginatedServices" :key="service.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ service.id }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ service.service_name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <span
                class="px-2 py-1 text-xs rounded-full"
                :class="{
                  'bg-green-100 text-green-800': service.load_type === 0,
                  'bg-blue-100 text-blue-800': service.load_type === 1,
                  'bg-purple-100 text-purple-800': service.load_type === 2,
                }"
              >
                {{ getLoadTypeText(service.load_type) }}
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">{{ service.service_desc || '-' }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(service.created_at) }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button
                @click="editService(service)"
                class="text-blue-600 hover:text-blue-900 mr-4"
              >
                编辑
              </button>
              <button
                @click="confirmDelete(service)"
                class="text-red-600 hover:text-red-900"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- 分页 -->
      <div v-if="filteredServices.length > 0" class="bg-white px-6 py-4 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            显示第 {{ (currentPage - 1) * pageSize + 1 }} 到 {{ Math.min(currentPage * pageSize, filteredServices.length) }} 条，共 {{ filteredServices.length }} 条
          </div>
          <div class="flex items-center space-x-2">
            <button
              @click="currentPage--"
              :disabled="currentPage === 1"
              class="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              上一页
            </button>
            <span class="text-sm text-gray-700">
              第 {{ currentPage }} / {{ totalPages }} 页
            </span>
            <button
              @click="currentPage++"
              :disabled="currentPage === totalPages"
              class="px-3 py-1 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              下一页
            </button>
            <select
              v-model="pageSize"
              class="ml-4 px-2 py-1 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option :value="10">10条/页</option>
              <option :value="20">20条/页</option>
              <option :value="50">50条/页</option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑服务弹窗 -->
    <div v-if="showCreateModal || showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
        <h2 class="text-xl font-bold text-gray-800 mb-4">{{ showEditModal ? '编辑服务' : '添加服务' }}</h2>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-gray-700 text-sm font-bold mb-2">服务名称</label>
            <input
              v-model="formData.service_name"
              type="text"
              class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div>
            <label class="block text-gray-700 text-sm font-bold mb-2">负载类型</label>
            <select
              v-model="formData.load_type"
              class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            >
              <option :value="0">HTTP</option>
              <option :value="1">TCP</option>
              <option :value="2">GRPC</option>
            </select>
          </div>
          <div>
            <label class="block text-gray-700 text-sm font-bold mb-2">服务描述</label>
            <textarea
              v-model="formData.service_desc"
              class="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              rows="3"
            ></textarea>
          </div>
          <div class="flex justify-end space-x-4 pt-4">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 border rounded-lg hover:bg-gray-100 transition"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="submitting"
              class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition disabled:bg-blue-300"
            >
              {{ submitting ? '提交中...' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-sm p-6">
        <h2 class="text-xl font-bold text-gray-800 mb-4">确认删除</h2>
        <p class="text-gray-600 mb-6">确定要删除服务 "{{ serviceToDelete?.service_name }}" 吗？此操作不可恢复。</p>
        <div class="flex justify-end space-x-4">
          <button
            @click="showDeleteModal = false"
            class="px-4 py-2 border rounded-lg hover:bg-gray-100 transition"
          >
            取消
          </button>
          <button
            @click="handleDelete"
            :disabled="deleting"
            class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition disabled:bg-red-300"
          >
            {{ deleting ? '删除中...' : '删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useApi } from '~/composables/useApi'
import { useMessage } from '~/composables/useMessage'

const api = useApi()
const messageStore = useMessage()

const services = ref([])
const loading = ref(false)
const searchQuery = ref('')
const loadTypeFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const serviceToDelete = ref(null)
const editingServiceId = ref(null)

const formData = ref({
  service_name: '',
  load_type: 0,
  service_desc: ''
})

// 获取服务列表
const fetchServices = async () => {
  loading.value = true
  try {
    const response = await api.getServices()
    if (response.data?.success) {
      services.value = response.data.data || []
    } else {
      messageStore.error(response.data?.error || '获取服务列表失败')
    }
  } catch (error) {
    messageStore.error('获取服务列表失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 过滤后的服务列表
const filteredServices = computed(() => {
  let result = services.value

  if (searchQuery.value) {
    result = result.filter(s =>
      s.service_name?.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  if (loadTypeFilter.value !== '') {
    result = result.filter(s => s.load_type === parseInt(loadTypeFilter.value))
  }

  return result
})

// 分页后的服务列表
const paginatedServices = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredServices.value.slice(start, end)
})

// 总页数
const totalPages = computed(() => {
  return Math.ceil(filteredServices.value.length / pageSize.value)
})

// 监听分页大小变化
watch(pageSize, () => {
  currentPage.value = 1
})

// 监听过滤条件变化
watch([searchQuery, loadTypeFilter], () => {
  currentPage.value = 1
})

// 获取负载类型文本
const getLoadTypeText = (type) => {
  const types = {
    0: 'HTTP',
    1: 'TCP',
    2: 'GRPC'
  }
  return types[type] || '未知'
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 编辑服务
const editService = (service) => {
  editingServiceId.value = service.id
  formData.value = {
    service_name: service.service_name,
    load_type: service.load_type,
    service_desc: service.service_desc || ''
  }
  showEditModal.value = true
}

// 确认删除
const confirmDelete = (service) => {
  serviceToDelete.value = service
  showDeleteModal.value = true
}

// 处理删除
const handleDelete = async () => {
  deleting.value = true
  try {
    const response = await api.deleteService(serviceToDelete.value.id)
    if (response.data?.success) {
      messageStore.success('删除成功')
      showDeleteModal.value = false
      serviceToDelete.value = null
      await fetchServices()
    } else {
      messageStore.error(response.data?.error || '删除失败')
    }
  } catch (error) {
    messageStore.error('删除失败: ' + (error.message || '未知错误'))
  } finally {
    deleting.value = false
  }
}

// 关闭弹窗
const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  formData.value = {
    service_name: '',
    load_type: 0,
    service_desc: ''
  }
  editingServiceId.value = null
}

// 处理提交
const handleSubmit = async () => {
  submitting.value = true
  try {
    let response
    if (showEditModal.value) {
      response = await api.updateService(editingServiceId.value, formData.value)
    } else {
      response = await api.createService(formData.value)
    }

    if (response.data?.success) {
      messageStore.success(showEditModal.value ? '更新成功' : '创建成功')
      closeModal()
      await fetchServices()
    } else {
      messageStore.error(response.data?.error || '操作失败')
    }
  } catch (error) {
    messageStore.error('操作失败: ' + (error.message || '未知错误'))
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchServices()
})
</script>
