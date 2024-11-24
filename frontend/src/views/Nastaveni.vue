<script setup lang="ts">
import axios from "axios";
import { useHead } from "unhead";
import { onMounted, ref, useTemplateRef } from "vue";
import { getToken, pridatOznameni, postKlavesnice } from "../utils";
import { prihlasen, uziv } from "../stores";
import { useRouter } from "vue-router";
import { role } from "../stores";

useHead({
    title: "Nastavení"
})

const router = useRouter()

const smazatPotvrzeni = ref(false)
const info = ref({ id: -1, role: "basic", klavesnice: "", jmeno: "...", email: "...@..." })
const klavesniceUprava = ref(false)
const jmenoUprava = ref("")

onMounted(() => {
    if (!getToken()) {
        pridatOznameni("Nejsi přihlášený")
        router.push("/prihlaseni")
    }
    get()
})

function get() {
    axios.get("/nastaveni", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(resp => {
        info.value = resp.data
        klavesniceUprava.value = info.value.klavesnice == "qwerty"
        jmenoUprava.value = info.value.jmeno

        uziv.value.jmeno = info.value.jmeno
    }).catch(e => {
        console.log(e)
    })
}

function postSmazat(e: Event) {
    e.preventDefault()
    axios.post("/ucet-zmena", { "zmena": "smazat" }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(() => {
        prihlasen.value = false
        localStorage.removeItem("pavouk_token")
        router.push("/prihlaseni")
    }).catch(e => {
        console.log(e)
        pridatOznameni()
    })
}

function postJmeno() {
    axios.post("/ucet-zmena", { "zmena": "jmeno", "hodnota": jmenoUprava.value }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(() => {
        get()
        jmenoInput.value?.blur() // lose focus
    }).catch(e => {
        if (e.response.data.error.search("uzivatel_jmeno_key")) {
            pridatOznameni("Takové jméno už někdo má")
        }
    })
}

function zmenaJmena(e: Event) {
    e.preventDefault()
    if (jmenoUprava.value == info.value.jmeno) {
        return
    }
    if (/^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$/.test(jmenoUprava.value)) {
        postJmeno()
    } else {
        if (jmenoUprava.value.length < 3) pridatOznameni("Jméno je moc krátké.<br>(3-12 znaků)")
        else if (jmenoUprava.value.length > 12) pridatOznameni("Jméno je moc dlouhé.<br>(3-12 znaků)")
        else pridatOznameni("Jméno může obsahovat jen velká a malá písmena, čísla a znaky _-+*!?")
    }
}

let timeoutID = 0
function zrusitPotvrzeni() {
    timeoutID = setTimeout(() => { smazatPotvrzeni.value = false }, 800)
}

function zrusitTimeout() {
    clearTimeout(timeoutID)
}

const jmenoInput = useTemplateRef("jmenoInput")

</script>
<template>
    <h1>Nastavení</h1>
    <div id="nastaveni">
        <form @submit="zmenaJmena">
            <div>
                <h3>Uživatelské jméno:</h3>
                <input ref="jmenoInput" v-model="jmenoUprava" type="text">
            </div>
            <button type="submit" class="tlacitko">Uložit</button>
        </form>

        <form>
            <div>
                <h3>Preference klávesnice:</h3>

                <input type="checkbox" id="toggle1" class="toggle-checkbox" v-model="klavesniceUprava" @change="postKlavesnice(klavesniceUprava)" />
                <label for="toggle1" class="toggle-contejner">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </div>
            <RouterLink to="/klavesnice?kam=nastaveni">Nevím, jakou mám klávesnici.</RouterLink>

            <span>
                Klávesnice ovlivňuje nejen <b>grafickou klácesnici</b> pod textem, ale i kurz samotný. Od první lekce horní řady se totiž lyší množina
                již naučených písmen podle toho, jestli už umíme <b>z</b> nebo <b>y</b>.
            </span>
        </form>

        <form>
            <div>
                <h3>Změna hesla:</h3>
                <button type="button" @click="router.push('/zapomenute-heslo?email=' + uziv.email)" class="tlacitko">Změnit heslo</button>
            </div>
        </form>

        <form>
            <div>
                <h3>Školní systém:</h3>
                <button type="button" @click="router.push('/zapis')" class="tlacitko">Zapsat se</button>
            </div>
        </form>

        <form>
            <div>
                <h3>Smazání účtu:</h3>
                <button v-if="!smazatPotvrzeni && role != 'student'" @click="smazatPotvrzeni = true" class="cervene-tlacitko">Smazat účet</button>
                <button v-else-if="role != 'student'" @click="postSmazat" class="cervene-tlacitko" @mouseleave="zrusitPotvrzeni"
                    @mouseenter="zrusitTimeout">Opravdu?</button>
            </div>
        </form>
    </div>
</template>
<style scoped>
.cervene-tlacitko {
    min-width: 150px;
}

form input[type=text] {
    width: 43%;
    height: 36px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    transition: all 0.2s cubic-bezier(0.5, 0, 0.5, 1) 0s;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 20px;
}

form input[type=text]:focus {
    outline: none !important;
    width: calc(43% + 10px);
}

form input::placeholder {
    color: rgba(255, 255, 255, 0.42);
}

form>:nth-child(2) {
    margin-top: 6px !important;
}

/* eslint-disable-next-line vue-scoped-css/no-unused-selector */
a {
    color: var(--bila);
    text-decoration: underline;
    text-decoration-color: #f0f0f07c;
    transition: 0.2s;
}

/* eslint-disable-next-line vue-scoped-css/no-unused-selector */
a:hover {
    text-decoration-color: #f0f0f0c5;
}

@media screen and (max-width: 460px) {
    form input[type=text] {
        width: 90%;
    }

    form input[type=text]:focus {
        width: calc(90% + 10px);
    }
}
</style>