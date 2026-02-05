import { createWebHistory, createRouter } from 'vue-router'


import About from "~/pages/about.vue";
import NotFound from "~/pages/404.vue";
import Login from "~/pages/login/login.vue";
import Dashboard from "~/pages/dashboard.vue";
import Admin from "~/layout/admin.vue"
import ServiceList from "~/pages/services/list.vue"
import Analysis from "~/pages/AI/analysis.vue"
import ModelList from "~/pages/AI/modelList.vue"

const routes = [
    { path: '/', component: Admin ,
        children:[
            { path: '/', component: Dashboard ,meta:{title:'仪表盘'}},
            { path: '/service-list', component: ServiceList ,meta:{title:'服务列表'}},
            { path: '/ai-analysis', component: Analysis ,meta:{title:'AI分析'}},
            { path: '/ai-modelList', component: ModelList ,meta:{title:'AI模型管理'}},
        ]
    },

    { path: '/login', component: Login ,meta:{title:'登陆'} },
    { path: '/about', component: About ,meta:{title:'关于'}},

    { path: '/:pathMatch(.*)*',//404
      name: 'not-found',
      component: NotFound}

]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router