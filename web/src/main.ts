import './style.css'
import { createApp } from 'vue'
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import MyTable from "@/components/table.vue";
import App from './App.vue'
import router from './router'



export const app = createApp(App)
app.component("my-table", MyTable);

app.use(router)
app.use(ElementPlus)

for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')