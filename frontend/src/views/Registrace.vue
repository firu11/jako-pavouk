<script setup lang="ts">
import axios from 'axios';
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { prihlasen, tokenJmeno } from '../stores';
import { pridatOznameni } from '../utils';

const router = useRouter()

const heslo = ref("")
const jmeno = ref("")
const email = ref("")
const kod = ref("")
const spatnyHeslo = ref(false)
const spatnyJmeno = ref(false)
const spatnyEmail = ref(false)
const spatnyKod = ref(false)
const emailExistuje = ref(false)
const jmenoExistuje = ref(false)

const overeni = ref(false)

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
        .then(_ => {
            overeni.value = true
            heslo.value = "" // radsi uz smazem idk
        }).catch(e => {
            if (e.response.data.error.search("email") != -1) emailExistuje.value = true
            else if (e.response.data.error.search("jmen") != -1) jmenoExistuje.value = true
            else pridatOznameni()
        })
}

function overeniPost(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!spatnyKod.value && kod.value.length == 5) {
        axios
            .post('/overit-email', {
                "email": email.value,
                "kod": kod.value
            })
            .then(response => {
                localStorage.setItem(tokenJmeno, response.data.token)
                prihlasen.value = true
                router.push("/ucet")
            }).catch(e => {
                if (e.response.data.error.search("kod") != -1) spatnyKod.value = true
                else pridatOznameni()
            })
    }
}

function chekujUdaje(jaky: string) {
    if (jaky === 'email' && email.value) spatnyEmail.value = !/^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(email.value); //test jestli email
    else if (jaky === 'heslo' && heslo.value !== undefined) spatnyHeslo.value = !/^(?=.*[0-9])(?=.*[!@#$%^&*_])[a-zA-Z0-9!@#$%^&*_]{8,25}$/.test(heslo.value) //heslo 8-25 aspon jeden CAPS a *_!
    else if (jaky === 'jmeno' && jmeno.value !== undefined) spatnyJmeno.value = !/^[a-zA-Z0-9!@#$%^&*_ ]{3,12}$/.test(jmeno.value) //jmeno 3-12
    else if (jaky === 'kod' && kod.value !== undefined) spatnyKod.value = !/^\d{5}$/.test(kod.value) //kod
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
    <div v-if="!overeni">
        <form class="pruhledne">
            <h3 class="nadpis">Uživatelské jméno:</h3>
            <input :class="{ spatnej_input: spatnyJmeno || jmenoExistuje }" @:input="chekujUdaje('jmeno')" type="text"
                v-model="jmeno" placeholder="Např: Pepa z depa">
            <h4 :class="{ opacity0: !jmenoExistuje }" class="chybaExistuje">Uživatel s tímto jménem už existuje</h4>
            <h3 class="nadpis">Email:</h3>
            <input :class="{ spatnej_input: spatnyEmail || emailExistuje }" @:input="chekujUdaje('email')" type="text"
                v-model="email" placeholder="Např: pepa@zdepa.cz">
            <h4 :class="{ opacity0: !emailExistuje }" class="chybaExistuje">Uživatel s tímto emailem už existuje</h4>
            <h3 class="nadpis">Heslo: <img src="../assets/icony/info.svg" alt="info" @mouseover="open_info"
                    @mouseleave="close_info"></h3>
            <input :class="{ spatnej_input: spatnyHeslo }" @:input="chekujUdaje('heslo')" type="password" v-model="heslo"
                placeholder='Rozhodně ne "Pepa123"'>
            <button type="submit" class="tlacitko" @click="registr">Registrovat</button>
        </form>
        <div id="info_out" class="info">
            Heslo musí obsahovat:
            <ul>
                <li>Minimálně 8 znaků</li>
                <li>Alespoň jeden speciální znak (!@#$%^&*_)</li>
                <li>Alespoň jedna číslice</li>
            </ul>
        </div>
    </div>
    <div v-else>
        <form class="pruhledne">
            <h3 style="margin-bottom: 20px;">Na email vám byl zaslán <br>ověřovací kód</h3>
            <h3 class="nadpis">Zadejte kód z emailu:</h3>
            <input :class="{ spatnej_input: spatnyKod }" @:input="chekujUdaje('kod')" type="text" inputmode="numeric"
                v-model="kod" placeholder="Např: 12345">
            <button type="submit" class="tlacitko" @click="overeniPost">Potvrdit</button>
        </form>
    </div>
    <p>Máte už účet?
        <router-link to="/prihlaseni">Přihlášení</router-link>
    </p>
</template>

<style scoped>
@import "../loginRegisterForma.css";
</style>
