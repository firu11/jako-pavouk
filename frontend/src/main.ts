import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import router from "./router.ts";
import axios from "axios";
import { createHead } from '@unhead/vue/client';
import vue3GoogleLogin from "vue3-google-login";

const app = createApp(App)
const head = createHead()

axios.defaults.baseURL = "http://localhost:1323/api" // http:||localhost:1323|api na production jen |api

app.use(vue3GoogleLogin, {
    clientId: "873874261202-aekjikhkkkfnmbdo68crs8l8e252b7rf.apps.googleusercontent.com"
})

app.use(router)
app.use(head)
app.mount("#app")

console.log("%cCo sem koukáš koloušku?", "color: white; font-size: x-large") // troulin
