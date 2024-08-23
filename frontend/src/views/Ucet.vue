<script setup lang="ts">
import axios from "axios"
import { prihlasen, role, tokenJmeno } from "../stores";
import { useRouter } from "vue-router";
import { onMounted, ref } from "vue";
import { checkTeapot, getToken, MojeMapa, pridatOznameni } from "../utils";
import { useHead } from "unhead"
import Tooltip from "../components/Tooltip.vue";
import AnimaceCisla from "../components/AnimaceCisla.vue";

useHead({
    title: "Účet"
})

const router = useRouter()

const info = ref({ jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, medianRychlosti: -1, uspesnost: -1, klavesnice: "QWERTZ", celkovyCas: 0, nejcastejsiChyby: new Map })
const uprava = ref(false)

const klavesniceUprava = ref(false)
const jmenoUprava = ref("")

const pismenaChyby = ref([] as any[][])

const smazatPotvrzeni = ref(false)

function odhlasit() {
    localStorage.removeItem(tokenJmeno)
    role.value = "basic"
    prihlasen.value = false
    router.push("/prihlaseni")
}

function zaokrouhlit(cislo: number | null) {
    if (cislo == null) {
        return -1
    }
    return Math.round(cislo * 10) / 10
}

onMounted(() => {
    getInfo()
})

async function getInfo() {
    try {
        let resp = await axios.get("/ja", {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        })
        info.value = resp.data
        info.value.nejcastejsiChyby = new MojeMapa(Object.entries(info.value.nejcastejsiChyby)).top(6)
        pismenaChyby.value = Array.from(info.value.nejcastejsiChyby, ([name, value]) => ([name, value]))
        pismenaChyby.value.sort((a, b) => b[1] - a[1])
        jmenoUprava.value = resp.data.jmeno
        klavesniceUprava.value = resp.data.klavesnice.toLocaleLowerCase() == "qwerty"
        role.value = resp.data.role
    }
    catch (e: any) {
        if (!checkTeapot(e)) {
            router.push("/prihlaseni")
            prihlasen.value = false
        }
    }
    setTimeout(() => {
        let jmeno = document.getElementById("jmeno")
        let velikost = 2
        while (jmeno?.clientWidth! > 300) { // hnus ale potřebuju to zmenšit natolik aby se to tam vešlo
            jmeno!.style.fontSize = `${velikost}em`
            velikost -= 0.2
        }
        let email = document.getElementById("email")
        velikost = 1.5
        while (email?.clientWidth! > 300) { // hnus ale potřebuju to zmenšit natolik aby se to tam vešlo
            email!.style.fontSize = `${velikost}em`
            velikost -= 0.1
        }
    }, 1)
}

function postSmazat() {
    axios.post("/ucet-zmena", { "zmena": "smazat" }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(_ => {
        prihlasen.value = false
        localStorage.removeItem("pavouk_token")
        router.push("/prihlaseni")
    }).catch(e => {
        console.log(e)
        pridatOznameni()
    })
}

function postJmeno() {
    axios.post("/ucet-zmena", { "zmena": "jmeno", "hodnota": jmenoUprava.value }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(_ => {
        getInfo()
    }).catch(e => {
        if (e.response.data.error.search("uzivatel_jmeno_key")) {
            pridatOznameni("Takové jméno už někdo má")
        }
    })
}

function postKlavesnice() {
    let klavesnice = klavesniceUprava.value ? "qwerty" : "qwertz"
    axios.post("/ucet-zmena", { "zmena": "klavesnice", "hodnota": klavesnice }, { headers: { Authorization: `Bearer ${getToken()}` } }).then(_ => {
        getInfo()
    }).catch(e => {
        checkTeapot(e)
    })
}

function zmenaJmena(e: Event) {
    e.preventDefault()
    if (jmenoUprava.value == info.value.jmeno) {
        uprava.value = false
        return
    }
    if (/^[a-zA-Z0-9ěščřžýáíéůúťňďóĚŠČŘŽÝÁÍÉŮÚŤŇĎÓ_\-+*! ]{3,12}$/.test(jmenoUprava.value)) {
        postJmeno()
        uprava.value = false
    } else {
        if (jmenoUprava.value.length < 3) pridatOznameni("Jméno je moc krátké.<br>(3-12 znaků)")
        else if (jmenoUprava.value.length > 12) pridatOznameni("Jméno je moc dlouhé.<br>(3-12 znaků)")
        else pridatOznameni("Jméno může obsahovat jen velká a malá písmena, čísla a znaky _-+*!?")
    }
}

</script>

<template>
    <div id="ucet">
        <img src="../assets//pavoucekBezPozadi.svg" alt="uzivatel">
        <div id="nadpisy">
            <h1 v-if="!uprava" id="jmeno">{{ info.jmeno }}
                <img @click="uprava = true" id="upravit" src="../assets/icony/upravit.svg" alt="Upravit">
            </h1>
            <h2 v-if="!uprava" id="email">{{ info.email }}</h2>
            <form v-if="uprava">
                <input v-model="jmenoUprava" type="text">
                <button type="submit" @click="zmenaJmena" id="tlacitko">Uložit</button>
            </form>
        </div>
    </div>

    <div id="bloky">
        <div class="blok" id="progres">
            <div id="nacitani-pozadi">
                <div id="nacitani" :style="{ width: info.dokonceno + '%' }"></div>
            </div>
            <span class="popis" style="width: 100%;">
                <span style="margin-right: 4px;">Dokončeno: </span>
                <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.dokonceno)" /> %
            </span>

        </div>
        <div class="blok">
            <img src="../assets/icony/rychlost.svg" alt="Rychlost" width="75">
            <Tooltip zprava="Za neopravené chyby je adekvátní penalizace. (CPM = úhozů za minutu)" :sirka="200" style="width: 60%;">
                <span v-if="info.medianRychlosti == -1" class="popis">Rychlost:<br><span class="nic">Zatím nic</span></span>
                <span v-else class="popis">
                    Rychlost:<br>
                    <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.medianRychlosti)" /> CPM
                </span>
            </Tooltip>
        </div>
        <div class="blok">
            <img src="../assets/icony/kalendar.svg" alt="Kalendář">
            <span class="popis" style="width: 60%;">
                Počet dní v řadě:<br>
                <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.daystreak)" :desetine-mista="0" />
            </span>
        </div>
        <div class="blok">
            <img src="../assets/icony/iconaKlavesnice.svg" alt="Klavesnice" width="75">
            <span class="popis">
                Klávesnice:
                <input type="checkbox" id="toggle1" class="toggleCheckbox" v-model="klavesniceUprava" @change="postKlavesnice" />
                <label for="toggle1" class="toggleContainer">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </span>
        </div>
        <div class="blok" id="chyby">
            <div id="presnost">
                <img src="../assets/icony/terc.svg" alt="Přesnost">
                <Tooltip zprava="Přesnost zahrunuje chyby opravené i neopravené." :sirka="210">
                    <span v-if="info.uspesnost == -1" class="popis">Přesnost:<br><span class="nic">Zatím nic</span></span>
                    <span v-else class="popis">
                        Přesnost:<br>
                        <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.uspesnost)" /> %
                    </span>
                </Tooltip>
            </div>

            <div>
                <h2>Nejčastější chyby:</h2>
                <hr>
            </div>
            <div style="width: 100%;">
                <div v-if="pismenaChyby.length !== 0" id="pismena">
                    <div id="prvni">
                        <span v-for="znak, i in pismenaChyby.slice(0, 2)"><span class="cisla">{{ i + 1 }}. </span>
                            <b :style="{ fontSize: znak[0] == ' ' ? '12px' : 'auto', fontWeight: znak[0] == ' ' ? '700' : '500' }">
                                {{ znak[0] == " " ? "┗━┛" : znak[0] }}
                            </b>
                        </span>
                    </div>
                    <div id="druhy">
                        <span v-for="znak, i in pismenaChyby.slice(2)"><span class="cisla">{{ i + 3 }}. </span>
                            <b :style="{ fontSize: znak[0] == ' ' ? '12px' : 'auto', fontWeight: znak[0] == ' ' ? '700' : '500' }">
                                {{ znak[0] == " " ? "┗━┛" : znak[0] }}
                            </b>
                        </span>
                    </div>
                </div>
                <div v-else style="margin: 30px 0;">
                    <span>Žádné!</span>
                </div>
            </div>
        </div>
    </div>

    <div id="tlacitka">
        <button @click="odhlasit" class="tlacitko">Odhlásit</button>
        <button v-if="!smazatPotvrzeni" @click="smazatPotvrzeni = true" class="cerveneTlacitko">Smazat účet</button>
        <button v-else @click="postSmazat" class="cerveneTlacitko">Opravdu?</button>
    </div>
</template>

<style scoped>
.nic {
    font-weight: 600;
}

.toggleContainer {
    font-size: 16px;
    width: calc(100% - 20px);
    height: 34px;
    margin: 10px 10px 0 10px;
}

.toggleContainer div {
    padding: 7px;
}

#pismena {
    display: flex;
    flex-direction: column;
    font-weight: 700;
    gap: 10px;
    width: 100%;
}

#prvni {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2em;
    gap: 20%;
}

#druhy {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.4em;
    gap: 12%;
    opacity: 70%;
}

#prvni b,
#druhy b {
    font-family: 'Red Hat Mono', monospace;
}

#chyby {
    height: auto;
    flex-direction: column;
    grid-row-start: 1;
    grid-row-end: 3;
    grid-column-start: 2;
    gap: 0;
    justify-content: space-between;
}

#chyby h2 {
    font-size: 1.6em;
    margin-bottom: 6px;
    margin-top: 15px;
    font-size: 1.2em;
}

#chyby>div:first-child {
    height: 65%;
}

#chyby hr {
    width: 180px;
    align-self: center;
    position: relative;
    top: -3px;
    margin-bottom: 4px;
}

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
}

.cislo {
    font-size: 28pt;
    font-weight: 480;
    font-family: 'Red Hat Mono';
}

#bloky {
    display: grid;
    width: 100%;
    gap: 20px;
    grid-template-columns: auto auto;
    justify-content: center;
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
    display: flex;
    flex-direction: column;
    text-decoration: none;
    align-items: normal;
}

#presnost {
    display: flex;
    text-decoration: none;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    background-color: var(--tmave-fialova);
    height: 120px;
    transition-duration: 0.2s;
    padding: 0 40px 0 25px;
    gap: 10px;
    max-height: 100px;
}

#nadpisy {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: start;
}

#nadpisy h1 {
    margin-bottom: 0;
    align-self: flex-start;
    direction: ltr;
    display: flex;
    flex-wrap: nowrap;
    align-items: center;
    max-width: 100%;
    font-weight: 500;
}

#nadpisy h2 {
    overflow: hidden;
    text-overflow: ellipsis !important;
    width: 100%;
    text-align: left;
}

#nadpisy h2:hover {
    overflow: visible;
    background-color: var(--tmave-fialova);
    width: auto;
    border-radius: 5px;
}

#nadpisy form {
    margin: 0;
    max-width: 250px;
}

#ucet img {
    height: 100px;
}

#upravit {
    width: 30px;
    height: 25px !important;
    cursor: pointer;
    margin: 3px;
}

#ucet {
    display: flex;
    background-color: var(--tmave-fialova);
    margin-bottom: 30px;
    padding: 15px 25px 15px 5px;
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
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 1.5em;
}

#ucet input:focus {
    outline: none !important;
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
    transition: 1000ms;
}

#druhKlavesnice {
    display: flex;
}

@media screen and (max-width: 1100px) {
    .toggleContainer {
        font-size: 16px;
        width: 100%;
        margin: 10px 0 0 0;
    }

    #chyby {
        order: -1;
    }

    #progres {
        order: -3;
    }

    .blok:nth-child(2) {
        order: -2;
    }

    .blok {
        width: 270px;
    }

    #ucet {
        flex-direction: column;
        padding: 25px 25px 25px 25px;
    }

    #ucet #nadpisy h2 {
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis !important;
        width: 100%;
        font-size: 1.3em;
        text-align: start;
    }

    #ucet #nadpisy h1 {
        font-size: 1.7em;
        max-width: 100% !important;
    }

    #ucet input {
        width: 100%;
    }

    .tlacitko,
    .cerveneTlacitko {
        width: 120px;
    }

    #nadpisy form {
        max-width: 215px;
    }

    #chyby>div:first-child {
        gap: 20px;
    }
}

@media screen and (max-width: 600px) {
    #bloky {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        max-width: 720px;
    }
}
</style>