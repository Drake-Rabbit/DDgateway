import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import 'virtual:windi.css'
import router from '~/router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import  pinia  from '~/store/index'



const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component)
}
import 'virtual:windi.css'

app.use(pinia)
app.use(ElementPlus)
app.use(router)


import '~/permission.js' // 引入路由守卫
app.mount('#app')