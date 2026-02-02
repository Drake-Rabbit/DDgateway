import axios  from "axios";
import { ElMessage } from 'element-plus'
import { getToken } from "~/composable/auth.js";

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
    return Promise.reject(error);
});