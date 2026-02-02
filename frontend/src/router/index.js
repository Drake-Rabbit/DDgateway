import { createWebHistory, createRouter } from 'vue-router'


import About from "~/pages/about.vue";
import NotFound from "~/pages/404.vue";
import Login from "~/pages/login/login.vue";
import Dashboard from "~/pages/dashboard.vue";
import Admin from "~/layout/admin.vue"

const routes = [
    { path: '/', component: Admin ,
        children:[
            { path: '/', component: Dashboard ,meta:{title:'仪表盘'}},
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