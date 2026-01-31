import { createWebHistory, createRouter } from 'vue-router'


import Index from "~/pages/index.vue";
import About from "~/pages/about.vue";
import NotFound from "~/pages/404.vue";
import Login from "~/pages/login/index.vue";

const routes = [
    { path: '/', component: Index },
    { path: '/login', component: Login },
    { path: '/about', component: About },
    { path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFound}
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router