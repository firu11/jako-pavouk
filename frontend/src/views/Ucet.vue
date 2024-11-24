<script setup lang="ts">
import axios from "axios";
import { prihlasen, role, tokenJmeno } from "../stores";
import { useRouter } from "vue-router";
import { onMounted, ref, useTemplateRef, watch } from "vue";
import { getToken, MojeMapa, pridatOznameni } from "../utils";
import { useHead } from "unhead";
import Tooltip from "../components/Tooltip.vue";
import AnimaceCisla from "../components/AnimaceCisla.vue";
import PrepinacTabu from "../components/PrepinacTabu.vue";

useHead({
    title: "Účet"
})

const router = useRouter()

type Stats = {
    rychlost: number
    cas: number
    nejcastejsiChyby: Map<string, number>
    uspesnost: number
}
const info = ref({ jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, 1: {} as Stats, 14: {} as Stats, "celkem": {} as Stats })
const uprava = ref(false)

const jmenoUprava = ref("")

const smazatPotvrzeni = ref(false)

const prepinacTabu = useTemplateRef("prepinac-tabu")

const rychlost = ref(-1)
const cas = ref(0)
const uspesnost = ref(-1)
const pismenaChyby = ref([] as { pismeno: string, pocet: number }[])

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
    axios.get("/ja", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(resp => {
        info.value = resp.data
        info.value[1].nejcastejsiChyby = new MojeMapa(Object.entries(info.value[1].nejcastejsiChyby)).top(6)
        info.value[14].nejcastejsiChyby = new MojeMapa(Object.entries(info.value[14].nejcastejsiChyby)).top(6)
        info.value["celkem"].nejcastejsiChyby = new MojeMapa(Object.entries(info.value["celkem"].nejcastejsiChyby)).top(6)
        jmenoUprava.value = resp.data.jmeno
        role.value = resp.data.role

        prepnoutStatistiky()
    })
    setTimeout(() => {
        let jmeno = document.getElementById("jmeno")
        let velikost = 2
        while (jmeno!.clientWidth > 300) { // hnus ale potřebuju to zmenšit natolik aby se to tam vešlo
            jmeno!.style.fontSize = `${velikost}em`
            velikost -= 0.2
        }
        let email = document.getElementById("email")
        velikost = 1.5
        while (email!.clientWidth > 300) { // hnus ale potřebuju to zmenšit natolik aby se to tam vešlo
            email!.style.fontSize = `${velikost}em`
            velikost -= 0.1
        }
    }, 1)
}

function postSmazat() {
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
        getInfo()
    }).catch(e => {
        if (e.response.data.error.search("uzivatel_jmeno_key")) {
            pridatOznameni("Takové jméno už někdo má")
        }
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

watch(() => prepinacTabu.value?.tab, prepnoutStatistiky)

function prepnoutStatistiky() {
    if (info.value.jmeno == "...") return
    if (prepinacTabu.value?.tab == "dnes") {
        rychlost.value = info.value[1].rychlost
        uspesnost.value = info.value[1].uspesnost
        cas.value = info.value[1].cas
        pismenaChyby.value = Array.from(info.value[1].nejcastejsiChyby, ([name, value]) => ({ pismeno: name as string, pocet: value as number }))
        pismenaChyby.value.sort((a, b) => b.pocet - a.pocet)

    } else if (prepinacTabu.value?.tab == "dva tydny") {
        rychlost.value = info.value[14].rychlost
        uspesnost.value = info.value[14].uspesnost
        cas.value = info.value[14].cas
        pismenaChyby.value = Array.from(info.value[14].nejcastejsiChyby, ([name, value]) => ({ pismeno: name as string, pocet: value as number }))
        pismenaChyby.value.sort((a, b) => b.pocet - a.pocet)

    } else if (prepinacTabu.value?.tab == "celkem") {
        rychlost.value = info.value["celkem"].rychlost
        uspesnost.value = info.value["celkem"].uspesnost
        cas.value = info.value["celkem"].cas
        pismenaChyby.value = Array.from(info.value["celkem"].nejcastejsiChyby, ([name, value]) => ({ pismeno: name as string, pocet: value as number }))
        pismenaChyby.value.sort((a, b) => b.pocet - a.pocet)
    }
}

</script>

<template>
    <div id="ucet">
        <img src="../assets/pavoucekBezPozadi.svg" alt="uzivatel">
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
                <div id="nacitani" :style="{ width: info.dokonceno + '%' }" />
            </div>
            <span class="popis" style="width: 100%;">
                <span style="margin-right: 4px;">Postup v kurzu: </span>
                <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.dokonceno)" /> %
            </span>
        </div>
        <div class="blok">
            <img src="../assets/icony/kalendar.svg" alt="Kalendář">
            <Tooltip
                :zprava="'Počítají se jak <b>cvičení</b> v kurzu, tak texty ze záložky <b>procvičování</b>.' + (role == 'student' ? ' (Práce ve třídě se nepočítají)' : '')"
                :sirka="200" :vzdalenost="3" style="width: 60%;">
                <span class="popis">
                    Počet dní v řadě:<br>
                    <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.daystreak)" :desetina-mista="0" />
                </span>
            </Tooltip>
        </div>
    </div>
    <div id="bloky">
        <div id="nadpis">
            <h2>Statistiky: </h2>
            <PrepinacTabu id="prepinac-tabu" :taby="[['celkem', 'Celkem'], ['dva tydny', 'Dva týdny'], ['dnes', 'Dnes']]" default-tab="celkem"
                ref="prepinac-tabu" />
        </div>

        <div class="blok">
            <img src="../assets/icony/rychlost.svg" alt="Rychlost" width="68">
            <Tooltip zprava="Za neopravené chyby je adekvátní penalizace. (<b>CPM</b> = úhozů za minutu)" :sirka="200" style="width: 60%;">
                <span v-if="rychlost == -1" class="popis">Rychlost:<br><span class="nic">Zatím nic</span></span>
                <span v-else class="popis">
                    Rychlost:<br>
                    <AnimaceCisla class="cislo" :cislo="zaokrouhlit(rychlost)" /> CPM
                </span>
            </Tooltip>
        </div>
        <div class="blok">
            <img src="../assets/icony/cas.svg" alt="Klavesnice" width="68">
            <span class="popis">
                Čas strávený psaním: <br>
                <span v-if="cas >= 3600">
                    <AnimaceCisla class="cislo" :cislo="(cas - (cas % 3600)) / 3600" :desetina-mista="0" /> h
                </span>
                <span v-if="(cas % 3600) / 60 >= 1">
                    <AnimaceCisla class="cislo" :cislo="((cas % 3600) - (cas % 60)) / 60" :desetina-mista="0" /> min
                </span>
                <span v-if="cas % 60 >= 1 && cas < 3660">
                    <AnimaceCisla class="cislo" :cislo="cas % 60" :desetina-mista="0" /> s
                </span>
                <span v-if="cas == 0" class="nic">Zatím nic</span>
            </span>
        </div>
        <div class="blok" id="chyby">
            <div id="presnost">
                <img src="../assets/icony/terc.svg" alt="Přesnost">
                <Tooltip zprava="Přesnost zahrunuje chyby <b>opravené</b> i <b>neopravené</b>. (Státní zkouška počítá pouze neopravené.)"
                    :sirka="210">
                    <span v-if="uspesnost == -1" class="popis">Přesnost:<br><span class="nic">Zatím nic</span></span>
                    <span v-else class="popis">
                        Přesnost:<br>
                        <AnimaceCisla class="cislo" :cislo="zaokrouhlit(uspesnost)" /> %
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
                        <span v-for="znak, i in pismenaChyby.slice(0, 2)" :key="znak.pismeno + znak.pocet"><span class="cisla">{{ i + 1 }}. </span>
                            <b :style="{ fontSize: znak.pismeno == ' ' ? '12px' : 'auto', fontWeight: znak.pismeno == ' ' ? '700' : '500' }">
                                {{ znak.pismeno == " " ? "┗━┛" : znak.pismeno }}
                            </b>
                        </span>
                    </div>
                    <div id="druhy">
                        <span v-for="znak, i in pismenaChyby.slice(2)" :key="znak.pismeno + znak.pocet"><span class="cisla">{{ i + 3 }}. </span>
                            <b :style="{ fontSize: znak.pismeno == ' ' ? '12px' : 'auto', fontWeight: znak.pismeno == ' ' ? '700' : '500' }">
                                {{ znak.pismeno == " " ? "┗━┛" : znak.pismeno }}
                            </b>
                        </span>
                    </div>
                </div>
                <div v-else style="margin: 27px 0;">
                    <span>Žádné!</span>
                </div>
            </div>
        </div>
    </div>

    <div id="tlacitka">
        <button @click="router.push('/klavesnice')" class="tlacitko">Klávesnice</button>
        <button @click="odhlasit" class="tlacitko">Odhlásit se</button>
        <button v-if="!smazatPotvrzeni && role != 'student'" @click="smazatPotvrzeni = true" class="cervene-tlacitko">Smazat účet</button>
        <button v-else-if="role != 'student'" @click="postSmazat" class="cervene-tlacitko">Opravdu?</button>
    </div>
</template>

<style scoped>
#nadpis {
    grid-row-start: 1;
    grid-column-start: 1;
    grid-column-end: 3;
    justify-self: center;
    margin: 30px 0 0 0;
    display: flex;
    gap: 8px;
    align-items: center;
    background-color: var(--tmave-fialova);
    border-radius: 30px;
    padding: 0 0 0 18px;
}

#prepinac-tabu {
    margin: 0;
}

.nic {
    font-weight: 600;
}

#pismena {
    display: flex;
    flex-direction: column;
    font-weight: 700;
    gap: 2px;
    width: 100%;
}

#prvni {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 32px;
    gap: 20%;
}

#druhy {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 22px;
    gap: 12%;
    opacity: 70%;
}

#prvni b,
#druhy b {
    font-family: "Red Hat Mono", monospace;
}

#chyby {
    height: auto;
    flex-direction: column;
    grid-row-start: 2;
    grid-row-end: 4;
    grid-column-start: 2;
    gap: 0;
    justify-content: space-between;
}

#chyby h2 {
    margin-bottom: 5px;
    margin-top: 12px;
    font-size: 19px;
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
    margin-top: 20px;
    gap: 20px;
    max-width: 80vw;
    flex-wrap: wrap;
    justify-content: center;
}

.tlacitko,
.cervene-tlacitko {
    margin-top: 0;
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
    font-size: 16px;
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
    font-family: "Red Hat Mono";
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
    height: 105px;
    transition-duration: 0.2s;
    padding: 10px;
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
    transition-duration: 0.2s;
    padding: 0 50px 0 20px;
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
    margin-bottom: 20px;
    padding: 10px 20px 10px 0px;
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
    font-size: 24px;
}

#ucet input:focus {
    outline: none !important;
}

#bloky div img {
    height: 58px;
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

@media screen and (max-width: 1100px) {
    #progres {
        order: -3;
    }

    .blok:nth-child(2) {
        order: -2;
    }

    .blok {
        width: 270px;
        height: 120px;
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
        font-size: 21px;
        text-align: start;
    }

    #ucet #nadpisy h1 {
        font-size: 27px;
        max-width: 100% !important;
    }

    #ucet input {
        width: 100%;
    }

    .tlacitko,
    .cervene-tlacitko {
        width: 120px;
    }

    #nadpisy form {
        max-width: 215px;
    }

    #chyby>div:first-child {
        gap: 20px;
    }

    #chyby #pismena>div:nth-child(2) {
        gap: 8%;
    }
}

@media screen and (max-width: 600px) {
    #bloky {
        grid-template-columns: auto;
        justify-items: center;
    }

    #nadpis {
        grid-row-start: 1;
        grid-column-start: 1;
        grid-column-end: 2;
        flex-direction: column;
        padding: 10px 0 0 0;
        border-radius: 10px;
        width: 270px;
        gap: 0;
        margin-top: 50px;
    }

    #nadpis h2 {
        font-size: 22px;
    }

    #prepinac-tabu {
        transform: scale(0.8);
        font-size: 17px;
    }

    #chyby {
        height: auto;
        grid-column-start: 1;
    }
}
</style>