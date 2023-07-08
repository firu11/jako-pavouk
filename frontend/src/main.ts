import { createApp } from "vue"
import "./style.css"
import App from "./App.vue"
import router from "./router.ts"
import axios from "axios"

const app = createApp(App)

axios.defaults.baseURL = "http://localhost:8080/api" // na production jen /api

app.use(router)
app.mount("#app")
