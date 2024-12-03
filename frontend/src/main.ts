import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import router from "./router.ts";
import axios from "axios";
import { createHead, setHeadInjectionHandler } from "@unhead/vue";
import vue3GoogleLogin from "vue3-google-login";

const app = createApp(App)
const head = createHead()

axios.defaults.baseURL = "http://localhost:1323/api" // http:||localhost:1323|api na production jen |api

app.use(vue3GoogleLogin, {
    clientId: "873874261202-aekjikhkkkfnmbdo68crs8l8e252b7rf.apps.googleusercontent.com"
})

app.use(router)
setHeadInjectionHandler(() => head) // zmizí warning: "inject() can only be used inside setup() or functional components." https://github.com/unjs/unhead/discussions/375
app.mount("#app")

console.log("%cCo sem koukáš koloušku?", "color: white; font-size: x-large") // troulin