import { createApp } from "vue"
import "./style.css"
import App from "./App.vue"
import router from "./router.ts"
import axios from "axios"
import { createHead } from "@unhead/vue"

const app = createApp(App)
const head = createHead()

axios.defaults.baseURL = "http://localhost:8080/api" // na production jen /api

app.use(router)
app.use(head)
app.mount("#app")
