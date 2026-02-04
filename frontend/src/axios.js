import axios  from "axios";
import { ElMessage } from 'element-plus'
import { getToken, removeToken } from "~/composable/auth.js";
import {useUserStore} from '~/store/user'
import router from '~/router' // 👈 导入 router 实例

const service = axios.create( {
    baseURL: 'http://localhost:8080/api'
})

export default service;

// 添加请求拦截器
service.interceptors.request.use(function (config) {
    // 在发送请求之前做些什么
    //往headers添加token
    const token = getToken()

    if (token) {
        config.headers['Authorization'] = 'Bearer ' + token
    }
    return config;
}, function (error) {
    // 对请求错误做些什么
    return Promise.reject(error);
});

// 添加响应拦截器
service.interceptors.response.use(function (response) {
    // 2xx 范围内的状态码都会触发该函数。
    // 对响应数据做点什么
    if (response.data.success == false) {
        ElMessage({
            message: response.data.error,
            type: 'error',
            duration: 2000
        })
    }
    return response;
}, function (error) {
    // 超出 2xx 范围的状态码都会触发该函数。
    // 对响应错误做点什么
    if (error.status == 401) {
        ElMessage({
            message: '登录过期或token失效，请重新登录',
            type: 'error',
            duration: 2000
        })
        useUserStore().removeUserInfo() //清除pina的用户信息状态
        removeToken() //移除cookie的token
        router.push('/login')
    }
    return Promise.reject(error);
});