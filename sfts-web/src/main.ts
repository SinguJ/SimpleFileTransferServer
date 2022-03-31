import { createApp } from 'vue'
import App from './App.vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import ApiPlugin from './api'

// English: Create Vue application instance
// 汉语：创建 Vue 应用实例
const app = createApp(App)

// English: Register `axios` components
// 汉语：注册 `axios` 组件
app.use(VueAxios, axios)

// English: Register `api` components
// 汉语：注册 `api` 组件
app.use(ApiPlugin, axios)

// English: Mount this Vue application instance on the `#app` node
// 汉语：在 `#app` 节点上挂载这个 Vue 应用实例
app.mount('#app')
