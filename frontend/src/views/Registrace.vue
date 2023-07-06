<script setup lang="ts">
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { prihlasen, tokenJmeno } from '../stores';

const router = useRouter()

const heslo = ref("")
const jmeno = ref("")
const email = ref("")
const spatnyHeslo = ref(false)
const spatnyJmeno = ref(false)
const spatnyEmail = ref(false)
const emailExistuje = ref(false)
const jmenoExistuje = ref(false)

function registr(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!heslo.value) spatnyHeslo.value = true
    if (!email.value) spatnyEmail.value = true
    if (!jmeno.value) spatnyJmeno.value = true

    if (spatnyEmail.value || spatnyHeslo.value || spatnyJmeno.value) return;

    axios
        .post('/registrace', {
            "jmeno": jmeno.value,
            "email": email.value,
            "heslo": heslo.value
        })
        .then(response => {
            if (response.data['error'] === "email") emailExistuje.value = true
            else if (response.data['error'] === "jmeno") jmenoExistuje.value = true
            else {
                localStorage.setItem(tokenJmeno, response.data.token)
                prihlasen.value = true
                router.push("/ucet")
            }
        })
}

function chekuj_udaje(jaky: string) {
    if (jaky === 'email' && email.value) spatnyEmail.value = !/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(email.value); //test jestli email
    else if (jaky === 'heslo' && heslo.value !== undefined) spatnyHeslo.value = !/^(?=.*[0-9])(?=.*[!@#$%^&*_])[a-zA-Z0-9!@#$%^&*_]{8,25}$/.test(heslo.value) //heslo 8-25 aspon jeden CAPS a *_!
    else if (jaky === 'jmeno' && jmeno.value !== undefined) spatnyJmeno.value = !/^[a-zA-Z0-9!@#$%^&*_ ]{3,25}$/.test(jmeno.value) //jmeno 3-25
    if (jaky === 'email') emailExistuje.value = false
    else if (jaky === 'jmeno') jmenoExistuje.value = false
}

function open_info() {
    document.getElementsByClassName('info')[0].id = 'info_show';
}

function close_info() {
    document.getElementsByClassName('info')[0].id = 'info_out';
}

</script>

<template>
    <h2>Registrace</h2>
    <form class="pruhledne">
        <h3 class="nadpis">Uživatelské jméno:</h3>
        <input :class="{ spatnej_input: spatnyJmeno || jmenoExistuje }" @:input="chekuj_udaje('jmeno')" type="text"
            v-model="jmeno" placeholder="Např: Pepa z depa">
        <h4 :class="{ opacity0: !jmenoExistuje }" class="chybaExistujee">Uživatel s tímto jménem už existuje</h4>
        <h3 class="nadpis">Email:</h3>
        <input :class="{ spatnej_input: spatnyEmail || emailExistuje }" @:input="chekuj_udaje('email')" type="text"
            v-model="email" placeholder="Např: pepa@zdepa.cz">
        <h4 :class="{ opacity0: !emailExistuje }" class="chybaExistujee">Uživatel s tímto emailem už existuje</h4>
        <h3 class="nadpis">Heslo: <img src="../assets/icony/info.svg" alt="info" @mouseover="open_info"
                @mouseleave="close_info"></h3>
        <input :class="{ spatnej_input: spatnyHeslo }" @:input="chekuj_udaje('heslo')" type="text" v-model="heslo"
            placeholder='Rozhodně ne "Pepa123"'>
        <button class="tlacitko" @click="registr">Registrovat</button>
    </form>
    <div id="info_out" class="info">
        Heslo musí obsahovat:
        <ul>
            <li>Minimálně 8 znaků</li>
            <li>Alespoň jeden speciální znak (!@#$%^&*_)</li>
            <li>Alespoň jedna číslice</li>
        </ul>
    </div>

    <p>Máte už účet?
        <router-link to="/prihlaseni">Přihlášení</router-link>
    </p>
</template>

<style scoped>
@import "../loginRegisterForma.css";
</style>
