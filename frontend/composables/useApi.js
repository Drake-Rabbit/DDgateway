import axios from 'axios'

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 10000
})

// 请求拦截器 - 添加 token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token 过期，清除登录状态
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// API 方法
export const useApi = () => {
  // 登录
  const login = (data) => {
    return apiClient.post('/auth/login', data)
  }

  // 注册
  const register = (data) => {
    return apiClient.post('/auth/register', data)
  }

  // 获取服务列表
  const getServices = (params = {}) => {
    return apiClient.get('/services/lists', { params })
  }

  // 获取服务详情
  const getServiceDetail = (id) => {
    return apiClient.get(`/services/detail/${id}`)
  }

  // 创建服务
  const createService = (data) => {
    return apiClient.post('/services/create', data)
  }

  // 更新服务
  const updateService = (id, data) => {
    return apiClient.post(`/services/update/${id}`, data)
  }

  // 删除服务
  const deleteService = (id) => {
    return apiClient.post(`/services/delete/${id}`)
  }

  // 获取租户列表
  const getTenants = () => {
    return apiClient.get('/tenants/list')
  }

  // 获取用户列表
  const getUsers = () => {
    return apiClient.get('/users/lists')
  }

  return {
    apiClient,
    login,
    register,
    getServices,
    getServiceDetail,
    createService,
    updateService,
    deleteService,
    getTenants,
    getUsers
  }
}
