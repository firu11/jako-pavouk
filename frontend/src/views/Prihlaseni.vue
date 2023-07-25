<script setup lang="ts">
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { prihlasen, tokenJmeno } from '../stores';
import { pridatOznameni } from '../utils';

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

    axios.post('/prihlaseni', {
        "email": email.value,
        "heslo": heslo.value
    }).then(response => {
        localStorage.setItem(tokenJmeno, response.data.token)
        prihlasen.value = true
        router.push("/ucet")
    }).catch(e => {
        if (e.response.status == 400 || e.response.status == 401) {
            if (e.response.data.error.search("Email") !== -1) spatnyEmail.value = true
            else if (e.response.data.error.search("Heslo") !== -1) spatnyHeslo.value = true
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

</script>

<template>
    <h2>Přihlášení</h2>
    <form class="pruhledne">
        <h3 class="nadpis">Email nebo jméno:</h3>
        <input style="margin-bottom: 20px" :class="{ spatnej_input: spatnyEmail }" :oninput="zmena" type="text"
            v-model="email" placeholder="Např: pepa@zdepa.cz">
        <h3 class="nadpis">Heslo:</h3>
        <input class="margin" :class="{ spatnej_input: spatnyHeslo }" :oninput="zmena" type="password" v-model="heslo"
            placeholder='Rozhodně ne "Pepa123"'>
        <button class="tlacitko" @click="login">Přihlásit</button>

    </form>
    <p>Nemáte ještě účet?
        <router-link to="/registrace">Registrace</router-link>
    </p>
</template>

<style scoped>
@import "../loginRegisterForma.css";
</style>
