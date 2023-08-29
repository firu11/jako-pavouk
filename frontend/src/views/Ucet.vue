<script setup lang="ts">
import axios from 'axios'
import { prihlasen, tokenJmeno } from '../stores';
import { useRouter } from 'vue-router';
import { onMounted, ref } from 'vue';
import { checkTeapot, getToken, pridatOznameni } from '../utils';
import { useHead } from 'unhead'

useHead({
    title: "Účet"
})

const router = useRouter()

const info = ref({ jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, prumerRychlosti: -1, uspesnost: -1, klavesnice: "QWERTZ" })
const uprava = ref(false)

const klavesniceUprava = ref("")
const jmenoUprava = ref("")

const smazatPotvrzeni = ref(false)

function odhlasit() {
    localStorage.removeItem(tokenJmeno)
    prihlasen.value = false
    router.push("/prihlaseni")
}

function zaokrouhlit(cislo: number | null) {
    if (cislo == null) {
        return "-"
    }
    return Math.round(cislo * 10) / 10
}

onMounted(() => {
    getInfo()
})

async function getInfo() {
    try {
        let resp = await axios.get('/ja', {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        })
        info.value = resp.data
        jmenoUprava.value = resp.data.jmeno
        klavesniceUprava.value = resp.data.klavesnice
    }
    catch (e: any) {
        if (!checkTeapot(e)) {
            router.push("/prihlaseni")
            prihlasen.value = false
        }
    }
}

function postSmazat() {
    axios.post('/ucet-zmena', { "zmena": "smazat" }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(_ => {
        router.push("/prihlaseni")
    }).catch(e => {
        console.log(e)
        pridatOznameni()
    })
}

function postJmeno() {
    axios.post('/ucet-zmena', { "zmena": "jmeno", "hodnota": jmenoUprava.value }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(_ => {
        getInfo()
    }).catch(e => {
        if (e.response.data.error.search("uzivatel_jmeno_key")) {
            pridatOznameni("Takové jméno už někdo má")
        }
    })
}

function postKlavesnice() {
    klavesniceUprava.value = klavesniceUprava.value == "QWERTZ" ? "QWERTY" : "QWERTZ" // otocim
    axios.post('/ucet-zmena', { "zmena": "klavesnice", "hodnota": klavesniceUprava.value }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(_ => {
        getInfo()
    }).catch(e => {
        checkTeapot(e)
    })
}

function zmenaJmena() {
    if (jmenoUprava.value != info.value.jmeno) {
        if (/^[a-zA-Z0-9!@#$%^&*_ ]{3,12}$/.test(jmenoUprava.value)) {
            postJmeno()
        } else {
            pridatOznameni("Jméno musí obsahovat jen znaky !@#$%^&*_ a může být 3-12 znaků dlouhé")
        }
    }
    uprava.value = false
}

</script>

<template>
    <div id="ucet">
        <img src="/pavoucekBezPozadi.svg" alt="uzivatel">
        <div id="nadpisy">
            <h1 v-if="!uprava">{{ info.jmeno }} <img v-if="!uprava" @click="uprava = true" id="upravit"
                    src="../assets/icony/upravit.svg" alt="Upravit"></h1>
            <h2 v-if="!uprava">{{ info.email }}</h2>
            <input v-if="uprava" v-model="jmenoUprava" type="text">
            <button v-if="uprava" type="submit" @click="zmenaJmena" id="tlacitko">Uložit</button>
        </div>
    </div>
    <div id="progres">
        <div id="nacitani-pozadi">
            <div id="nacitani" :style="{ width: info.dokonceno + '%' }"></div>
        </div>
        <span class="popis" style="width: 100%;">Dokončeno: <span class="cislo">{{ zaokrouhlit(info.dokonceno)
        }}%</span></span>
    </div>
    <div id="bloky">
        <div class="blok">
            <img src="../assets/icony/rychlost.svg" alt="Rychlost" width="75">
            <span v-if="info.prumerRychlosti == -1">Zatím nic</span>
            <span v-else class="popis">Rychlost: <br><span class="cislo">{{ zaokrouhlit(info.prumerRychlosti) }}</span>
                CPM</span>
        </div>
        <div class="blok">
            <img src="../assets/icony/iconaKlavesnice.svg" alt="Rychlost" width="75">
            <span class="popis">Klávesnice: <br>
                <button id="tlacitko" @click="postKlavesnice">{{ klavesniceUprava }}</button>
            </span>
        </div>
        <div class="blok">
            <img src="../assets/icony/terc.svg" alt="Přesnost">
            <span v-if="info.uspesnost == -1">Zatím nic</span>
            <span v-else class="popis">Přesnost: <br><span class="cislo">{{ zaokrouhlit(info.uspesnost) }}</span> %</span>
        </div>
        <div class="blok">
            <img src="../assets/icony/kalendar.svg" alt="Kalendář">
            <span class="popis">Počet dní v řadě: <br><span class="cislo">{{ zaokrouhlit(info.daystreak) }}</span></span>
        </div>
    </div>

    <div id="tlacitka">
        <button @click="odhlasit" class="tlacitko">Odhlásit</button>
        <button v-if="!smazatPotvrzeni" @click="smazatPotvrzeni = true" class="cerveneTlacitko">Smazat účet</button>
        <button v-else @click="postSmazat" class="cerveneTlacitko">Opravdu?</button>
    </div>
</template>

<style scoped>
#tlacitka {
    display: inline-flex;
    margin-top: 10px;
    gap: 20px;
}

.tlacitko {
    background-color: var(--tmave-fialova);
}

.tlacitko:hover {
    background-color: var(--fialova);
}

#tlacitko {
    width: 120px;
    height: 40px;
    border: none;
    border-radius: 5px;
    color: var(--bila);
    font-size: 1em;
    margin: 10px 0 0 0;
    background-color: var(--fialova);
    transition: 0.2s;
    cursor: pointer;
    align-self: center;
}

#tlacitko:hover {
    background-color: var(--svetle-fialova);
    transition: 0.2s;
}

.popis {
    font-size: 15pt;
    width: 60%;
}

.cislo {
    font-size: 28pt;
    font-weight: 500;
}

#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
    flex-wrap: wrap-reverse;
    justify-content: center;
    align-content: flex-end;
}

.blok {
    display: flex;
    text-decoration: none;
    border-radius: 10px;
    justify-content: space-evenly;
    align-items: center;
    width: 320px;
    background-color: var(--tmave-fialova);
    height: 120px;
    transition-duration: 0.2s;
    padding: 15px;
    gap: 10px;
}

#progres {
    margin-bottom: 20px;
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 420px;
    background-color: var(--tmave-fialova);
    height: 110px;
    transition-duration: 0.2s;
    padding: 15px;
    gap: 10px;
}

#nadpisy {
    display: flex;
    flex-direction: column;
    justify-content: center;
}

#nadpisy h1 {
    margin-bottom: 0;
    align-self: flex-start;
}

#ucet img {
    height: 100px;
}

#upravit {
    width: 30px;
    height: 25px !important;
}

#ucet {
    display: flex;
    background-color: var(--tmave-fialova);
    margin-bottom: 40px;
    padding: 15px 30px 15px 5px;
    border-radius: 10px;
    gap: 5px;
    justify-content: space-around;
}

#ucet input {
    max-width: 250px;
    height: 39px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    transition: all 0.15s cubic-bezier(0.5, 0, 0.5, 1) 0s;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 1.5em;
}

#ucet input:focus {
    outline: none !important;
    transition: all 0.15s cubic-bezier(0.5, 0, 0.5, 1) 0s;
}

#bloky div img {
    height: 65px;
}

#bloky div h2 {
    font-weight: 500;
}

#nacitani-pozadi {
    height: 20px;
    background-color: var(--fialova);
    border-radius: 10px;
    padding: 0;
    overflow: hidden;
}

#nacitani {
    background-color: var(--bila);
    height: 20px;
    position: relative;
    left: -1px
}

#druhKlavesnice {
    display: flex;
}

@media screen and (max-width: 1000px) {
    #bloky {
        align-items: center;
        justify-content: center;
    }

    #progres {
        width: 320px;
    }

    #ucet {
        flex-direction: column;
        padding: 25px 25px 25px 25px
    }

}
</style>