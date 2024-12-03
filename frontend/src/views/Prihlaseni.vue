<script setup lang="ts">
import axios from "axios";
import { ref } from "vue";
import { useRouter } from "vue-router";
import { prihlasen, tokenJmeno, uziv, role } from "../stores";
import { pridatOznameni } from "../utils";
import { useHead } from "unhead";

useHead({
    title: "Přihlášení",
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/prihlaseni"
        }
    ]
})

const router = useRouter()

const heslo = ref("")
const email = ref("")
const spatnyHeslo = ref(false)
const spatnyEmail = ref(false)

function login(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!heslo.value) { //pokud uzivatel nic nenapsal
        spatnyHeslo.value = true
    }
    if (!email.value) {
        spatnyEmail.value = true
    }
    if (spatnyEmail.value || spatnyHeslo.value) return //nezkoušet ani

    axios.post("/prihlaseni", {
        "email": email.value,
        "heslo": heslo.value
    }).then(response => {
        localStorage.setItem(tokenJmeno, response.data.token)
        prihlasen.value = true

        uziv.value.email = response.data.email
        uziv.value.jmeno = response.data.jmeno
        role.value = response.data.role
        
        router.push("/statistiky")
    }).catch(e => {
        if (e.response.status == 400 || e.response.status == 401) {
            if (e.response.data.error.search("Email") !== -1) {
                spatnyEmail.value = true
                pridatOznameni("Špatný email / jméno")
            }
            else if (e.response.data.error.search("Heslo") !== -1) {
                spatnyHeslo.value = true
                pridatOznameni("Špatné heslo")
            }
            else if (e.response.data.error.search("google") !== -1) {
                pridatOznameni(e.response.data.error)
            }
            else pridatOznameni()
        } else {
            pridatOznameni()
        }
    })
}

function zmena() { // pokud zacnu znova psat tak zrusim znaceni spatnyho inputu
    spatnyEmail.value = false
    spatnyHeslo.value = false
}

const handleLoginSuccess = (response: { credential: string }) => {
    axios.post("/google", {
        "access_token": response.credential,
    }).then(response => {
        localStorage.setItem(tokenJmeno, response.data.token)
        prihlasen.value = true

        uziv.value.email = response.data.email
        uziv.value.jmeno = response.data.jmeno
        role.value = response.data.role

        if (response.data.novy) router.push("/klavesnice?kam=statistiky")
        else router.push("/statistiky")
    }).catch(() => {
        pridatOznameni()
    })
}
</script>

<template>
    <h2>Přihlášení</h2>
    <form>
        <h3 class="nadpis">Email nebo jméno:</h3>
        <input :class="{ 'spatnej-input': spatnyEmail }" :oninput="zmena" type="text" v-model="email" placeholder="Např: pan@pavouk.cz"
            inputmode="email">
        <h3 class="nadpis">Heslo:</h3>
        <input :class="{ 'spatnej-input': spatnyHeslo }" :oninput="zmena" type="password" v-model="heslo" placeholder='Rozhodně ne "Pavouk123"'>
        <button type="submit" class="tlacitko" @click="login">Přihlásit</button>

        <hr id="predel">

        <GoogleLogin id="google" :callback="handleLoginSuccess" :error="pridatOznameni" :buttonConfig="{ text: 'continue_with' }" />
    </form>
    <p>Nemáš ještě účet?
        <router-link to="/registrace">Registrace</router-link>
    </p>

    <router-link to="/zapomenute-heslo">Zapomenuté heslo?</router-link>
</template>

<style scoped>
@import "../loginRegisterForma.css";
</style>
