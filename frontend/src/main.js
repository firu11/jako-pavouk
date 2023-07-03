import { createApp } from 'vue'
import App from '@/App.vue'
import router from '@/router.js'
import axios from 'axios'
import '@/main.css'
import { createWatcher } from 'next-vue-storage-watcher';
import { formatovany_pismena } from '@/utils'
import axiosRetry from 'axios-retry';


export const lsWatcher = createWatcher({
    prefix:"pavouk_" //https://github.com/dreambo8563/next-vue-storage-watcher
})

axiosRetry(axios, { retries: 2, retryDelay: axiosRetry.exponentialDelay });
axios.defaults.baseURL = "http://localhost:8080/";

const app = createApp(App)

app.config.globalProperties.$format = formatovany_pismena;

app.config.globalProperties.$ls = lsWatcher
app.provide('ls', lsWatcher)

app.use(router)
app.use(lsWatcher)
app.mount('#app')
