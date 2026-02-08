import axios from '~/axios.js'

// 获取服务列表
export function getServiceList(page_no, page_size) {
    return axios.get('/services/service_list', {
        params: {
            page_no,
            page_size
        }
    })
}


// 新增服务Http请求
export function addService(data) {
    return axios.post('/services/service_add_http', data)
}