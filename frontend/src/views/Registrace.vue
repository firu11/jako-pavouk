<script setup lang="ts">
import axios from "axios";
import { ref } from "vue";
import { onBeforeRouteLeave, useRouter } from "vue-router";
import { prihlasen, tokenJmeno } from "../stores";
import { pridatOznameni } from "../utils";
import { useHead } from "unhead"

useHead({
    title: "Registrace",
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/registrace"
        }
    ]
})

const router = useRouter()

const heslo = ref("")
const jmeno = ref("")
const email = ref("")
const kod = ref("")
const spatnyHeslo = ref(false)
const spatnyJmeno = ref(false)
const spatnyEmail = ref(false)
const spatnyKod = ref(false)

const overeni = ref(false)

const posilame = ref(false)

function registr(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!heslo.value) spatnyHeslo.value = true
    if (!email.value) spatnyEmail.value = true
    if (!jmeno.value) spatnyJmeno.value = true

    if (spatnyEmail.value || spatnyHeslo.value || spatnyJmeno.value) {
        if (spatnyJmeno.value && jmeno.value.length > 12) pridatOznameni("Jméno je moc krátké.<br>(3-12 znaků)")
        else if (spatnyJmeno.value && jmeno.value.length < 3) pridatOznameni("Jméno je moc dlouhé.<br>(3-12 znaků)")
        else if (spatnyJmeno.value) pridatOznameni("Jméno může obsahovat jen velká a malá písmena, čísla a znaky _-+*!?")
        else if (spatnyEmail.value) pridatOznameni("Email není validní.")
        else if (spatnyHeslo.value) pridatOznameni("Heslo musí obsahovat alespoň 5 znaků. Toť vše.")
        return
    }

    posilame.value = true

    axios
        .post("/registrace", {
            "jmeno": jmeno.value,
            "email": email.value,
            "heslo": heslo.value
        })
        .then(_ => {
            overeni.value = true
            heslo.value = "" // radsi uz smazem idk
            posilame.value = false
        }).catch(e => {
            if (e.response.data.error.search("email") != -1) {
                spatnyEmail.value = true
                pridatOznameni("Uživatel s tímto emailem už existuje.")
            }
            else if (e.response.data.error.search("docasne") != -1) {
                spatnyJmeno.value = true
                pridatOznameni("Jméno je dočasně rezervováno. Za max. 10min možná bude volné.")
            }
            else if (e.response.data.error.search("jmenem") != -1) {
                spatnyJmeno.value = true
                pridatOznameni("Uživatel s tímto jménem už existuje.")
            }
            else pridatOznameni()
            posilame.value = false
        })
}

function overeniPost(e: Event) {
    e.preventDefault(); //aby se nerefreshla stranka

    if (!spatnyKod.value && kod.value.length == 5) {
        axios
            .post("/overit-email", {
                "email": email.value,
                "kod": kod.value
            })
            .then(response => {
                localStorage.setItem(tokenJmeno, response.data.token)
                prihlasen.value = true
                router.push("/ucet")
            }).catch(e => {
                if (e.response.data.error.search("kod") != -1) {
                    spatnyKod.value = true
                    pridatOznameni("Špatný kód")
                }
                else if (e.response.data.error.search("Cas") != -1) {
                    kod.value = ""
                    pridatOznameni("Čas pro ověření vypršel. Zkus to prosím znovu.")
                    overeni.value = false
                }
                else pridatOznameni()
            })
    } else {
        pridatOznameni("Kód je 5ti místné číslo")
    }
}

function chekujUdaje(jaky: string) {
    if (jaky === "email" && email.value) spatnyEmail.value = !/^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/g.test(email.value)
    else if (jaky === "heslo" && heslo.value !== undefined) spatnyHeslo.value = !/^(?=.*[a-zA-Z]).{5,128}$/.test(heslo.value) //heslo min 5 znaku
    else if (jaky === "jmeno" && jmeno.value !== undefined) spatnyJmeno.value = !/^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$/.test(jmeno.value) //jmeno 3-12
    else if (jaky === "kod" && kod.value !== undefined) spatnyKod.value = !/^\d{5}$/.test(kod.value) //kod 5 dlouhy
    if (jaky === "email" && email.value.length === 0) spatnyEmail.value = false
    else if (jaky === "jmeno" && jmeno.value.length === 0) spatnyJmeno.value = false
    else if (jaky === "heslo" && heslo.value.length === 0) spatnyHeslo.value = false
}

function openInfo() {
    document.getElementsByClassName("info")[0].id = "infoShow";
}

function closeInfo() {
    document.getElementsByClassName("info")[0].id = "infoHide";
}

onBeforeRouteLeave(() => {
    if (!overeni.value || prihlasen.value) return
    const answer = window.confirm("Hej! Email už je na cestě! \nOpravdu chceš odejít a zahodit tento pokus?")
    if (!answer) return false
})

const handleLoginSuccess = (response: any) => {
    axios.post("/google", {
        "access_token": response.credential,
    }).then(response => {
        localStorage.setItem(tokenJmeno, response.data.token)
        prihlasen.value = true
        router.push("/ucet")
    }).catch(_ => {
        pridatOznameni()
    })
}

</script>

<template>
    <h2>Registrace</h2>
    <div v-if="!overeni">
        <form>
            <h3 class="nadpis">Uživatelské jméno:</h3>
            <input :class="{ spatnej_input: spatnyJmeno }" @:input="chekujUdaje('jmeno')" type="text" v-model="jmeno"
                placeholder="Např: Pan Pavouk">
            <h3 class="nadpis">Email:</h3>
            <input :class="{ spatnej_input: spatnyEmail }" @:input="chekujUdaje('email')" type="email" v-model="email"
                placeholder="Např: pan@pavouk.cz" inputmode="email">
            <h3 class="nadpis infoNadpis">Heslo: <img src="../assets/icony/info.svg" alt="info" @mouseover="openInfo"
                    @mouseleave="closeInfo"></h3>
            <input :class="{ spatnej_input: spatnyHeslo }" @:input="chekujUdaje('heslo')" type="password"
                v-model="heslo" placeholder="Rozhodně ne 'Pavouk123'">
            <button type="submit" class="tlacitko" @click="registr" :disabled="posilame">{{ posilame ? ". . ." :
        "Registrovat" }}</button>

            <hr id="predel">
            <GoogleLogin id="google" :callback="handleLoginSuccess" :error="pridatOznameni"
                :buttonConfig="{ text: 'continue_with' }" />
        </form>
        <div id="infoHide" class="info">
            Doporučujeme:
            <ul>
                <li>Minimálně 5 znaků</li>
                <li>Alespoň jedna číslice</li>
                <li>Alespoň jedno velké písmeno</li>
            </ul>
        </div>
    </div>
    <div v-else>
        <form id="overeni">
            <h3 style="margin-bottom: 20px;">Na email: "<i>{{ email }}</i>" ti byl zaslán ověřovací kód</h3>
            <h3 class="nadpis">Zadej kód z emailu:</h3>
            <input :class="{ spatnej_input: spatnyKod }" @:input="chekujUdaje('kod')" type="text" inputmode="numeric"
                v-model.trim="kod" placeholder="Např: 12345">
            <button type="submit" class="tlacitko" @click="overeniPost">Potvrdit</button>
        </form>
    </div>
    <p>Máš už účet?
        <router-link to="/prihlaseni">Přihlášení</router-link>
    </p>
</template>

<style scoped>
@import "../loginRegisterForma.css";
</style>
