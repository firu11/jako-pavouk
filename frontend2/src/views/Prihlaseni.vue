<script setup lang="ts">
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { prihlasen, token_jmeno } from '../stores';

const router = useRouter()

const heslo = ref("")
const email = ref("")
const spatny_heslo = ref(false)
const spatny_email = ref(false)

function login(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!heslo.value) { //pokud uzivatel nic nenapsal
        spatny_heslo.value = true
    }
    if (!email.value) {
        spatny_email.value = true
    }
    if (spatny_email.value || spatny_heslo.value) return //nezkoušet ani

    axios.post('/prihlaseni', {
        "email": email.value,
        "heslo": heslo.value
    }).then(response => {
        localStorage.setItem(token_jmeno, response.data.token)
        prihlasen.value = true
        router.push("/ucet")
    }).catch(e => {
        try {
            if (e.response.status == 400 || e.response.status == 401) {
                spatny_heslo.value = true
                spatny_email.value = true
            } else {
                alert("Něco se pokazilo na naší straně...")
            }
        } catch {
            console.log(e)
        }
    })
}

function zmena() { // pokud zacnu znova psat tak zrusim znaceni spatnyho inputu
    spatny_email.value = false
    spatny_heslo.value = false
}

</script>

<template>
    <h2>Přihlášení</h2>
    <form class="pruhledne">
        <h3 class="nadpis">Email:</h3>
        <input style="margin-bottom: 20px" :class="{ spatnej_input: spatny_email }" :oninput="zmena" type="text"
            v-model="email" placeholder="Např: pepa@zdepa.cz">
        <h3 class="nadpis">Heslo:</h3>
        <input class="margin" :class="{ spatnej_input: spatny_heslo }" :oninput="zmena" type="password" v-model="heslo"
            placeholder='Rozhodně ne "Pepa123"'>
        <button class="tlacitko" @click="login">Přihlásit</button>

    </form>
    <p>Nemáte ještě účet?
        <router-link to="/register">Registrace</router-link>
    </p>
</template>

<style scoped>
@import "../loginRegisterForma.css";
</style>
