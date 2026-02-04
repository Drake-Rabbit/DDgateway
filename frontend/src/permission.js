import router from '~/router'
import {getToken} from '~/composable/auth.js'
import { ElMessage } from 'element-plus'

const  whitelist = ['/login', '/register','/404']

//全局路由守卫--监测路由的
//前置
 router.beforeEach((to, from, next) => {
    // 检查用户是否已登录
    const isAuthenticated = getToken() // 假设 getToken 是检查登录状态的函数

    if (!isAuthenticated && !whitelist.includes(to.path)) {
        //如果用户未登录且尝试访问非登录页，则重定向到登录页
        ElMessage({
            message: '请先登录',
            type: 'error',
            duration: 2000
        })
         next('/login')
    } else if (isAuthenticated && to.path == '/login') {
        // 如果用户已登录且尝试访问登录页，则重定向到首页
        next('/')
    } else {
        // 其他情况，继续导航
        //设置页面标题
        let title = (to.meta.title? to.meta.title : '')+"-Gateway Service"
        document.title = title
        next()
    }
})

