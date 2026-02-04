import axios from '~/axios.js'


//登陆login
export function login(username, password) {
    return axios.post('/auth/login', {
        username,
        password
    })
}

//登出logout
export function logout() {
    return axios.post('/auth/logout')
}

//获取用户信息getUserInfo-只返回username和userid
export function getUserInfo() {
    return axios.post('/auth/userinfo')
}

//更新密码
export function updatePassword(data) {
    return axios.post('/auth/updatepassword', data)
}