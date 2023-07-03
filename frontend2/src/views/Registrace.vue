<script setup lang="ts">
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { prihlasen, token_jmeno } from '../stores';

const router = useRouter()

const heslo = ref("")
const jmeno = ref("")
const email = ref("")
const spatny_heslo = ref(false)
const spatny_jmeno = ref(false)
const spatny_email = ref(false)
const email_existuje = ref(false)
const jmeno_existuje = ref(false)

function registr(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!heslo.value) spatny_heslo.value = true
    if (!email.value) spatny_email.value = true
    if (!jmeno.value) spatny_jmeno.value = true

    if (spatny_email.value || spatny_heslo.value || spatny_jmeno.value) return;

    axios
        .post('/registrace', {
            "jmeno": jmeno.value,
            "email": email.value,
            "heslo": heslo.value
        })
        .then(response => {
            if (response.data['error'] === "email") email_existuje.value = true
            else if (response.data['error'] === "jmeno") jmeno_existuje.value = true
            else {
                localStorage.setItem(token_jmeno, response.data.token)
                prihlasen.value = true
                router.push("/ucet")
            }
        })
}

function chekuj_udaje(jaky: string) {
    if (jaky === 'email' && email.value) spatny_email.value = !/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(email.value); //test jestli email
    else if (jaky === 'heslo' && heslo.value !== undefined) spatny_heslo.value = !/^(?=.*[0-9])(?=.*[!@#$%^&*_])[a-zA-Z0-9!@#$%^&*_]{8,25}$/.test(heslo.value) //heslo 8-25 aspon jeden CAPS a *_!
    else if (jaky === 'jmeno' && jmeno.value !== undefined) spatny_jmeno.value = !/^[a-zA-Z0-9!@#$%^&*_ ]{3,25}$/.test(jmeno.value) //jmeno 3-25
    if (jaky === 'email') email_existuje.value = false
    else if (jaky === 'jmeno') jmeno_existuje.value = false
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
        <input :class="{ spatnej_input: spatny_jmeno || jmeno_existuje }" @:input="chekuj_udaje('jmeno')" type="text"
            v-model="jmeno" placeholder="Např: Pepa z depa">
        <h4 :class="{ opacity0: !jmeno_existuje }" class="chybaExistujee">Uživatel s tímto jménem už existuje</h4>
        <h3 class="nadpis">Email:</h3>
        <input :class="{ spatnej_input: spatny_email || email_existuje }" @:input="chekuj_udaje('email')" type="text"
            v-model="email" placeholder="Např: pepa@zdepa.cz">
        <h4 :class="{ opacity0: !email_existuje }" class="chybaExistujee">Uživatel s tímto emailem už existuje</h4>
        <h3 class="nadpis">Heslo: <img src="../assets/icony/info.svg" alt="info" @mouseover="open_info"
                @mouseleave="close_info"></h3>
        <input :class="{ spatnej_input: spatny_heslo }" @:input="chekuj_udaje('heslo')" type="text" v-model="heslo"
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
