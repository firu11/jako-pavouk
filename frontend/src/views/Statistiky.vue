<script setup lang="ts">
import { useHead } from "unhead";
import { onMounted, ref } from "vue";
import { role } from "../stores";
import { getToken, MojeMapa } from "../utils";
import axios from "axios";
import PrepinacTabu from "../components/PrepinacTabu.vue";
import Tooltip from "../components/Tooltip.vue";
import AnimaceCisla from "../components/AnimaceCisla.vue";
import GrafStatistiky from "../components/GrafStatistiky.vue";

useHead({
    title: "Statistiky"
})

const info = ref({ rychlost: -1, cas: 0, uspesnost: -1, postupVKurzu: 0, daystreak: 0, nejcastejsiChyby: new Map<string, number>() })

const nejcastejsiChyby = ref([] as { znak: string, pocet: number }[])

async function getInfo() {
    axios.get("/statistiky", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(resp => {
        info.value = resp.data
        nejcastejsiChyby.value = new MojeMapa(Object.entries(info.value.nejcastejsiChyby)).top(6)
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

function zaokrouhlit(cislo: number | null) {
    if (cislo == null) {
        return -1
    }
    return Math.round(cislo * 10) / 10
}
/* 
function prepnoutStatistiky() {
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
*/

onMounted(() => {
    getInfo()
})
</script>
<template>
    <h1>Statistiky</h1>

    <div id="bloky">
        <div id="mini">
            <div class="blok" id="progres">
                <div id="nacitani-pozadi" :style="{ bottom: `${(info.postupVKurzu - 16) / 1.15}%` }" />
                <div id="nacitani-pozadi2" :style="{ height: `${info.postupVKurzu >= 100 ? 100 : info.postupVKurzu - 10}%` }" />
                <span class="popis">
                    Kurz: <br>
                    <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.postupVKurzu)" /> %
                </span>
            </div>
            <div class="blok">
                <Tooltip
                    :zprava="'Počítají se jak <b>cvičení</b> v kurzu, tak texty ze záložky <b>procvičování</b>.' + (role == 'student' ? ' (Práce ve třídě se nepočítají)' : '')"
                    :sirka="200" :vzdalenost="3">
                    <span class="popis">
                        Dní v řadě:<br>
                        <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.daystreak)" :desetina-mista="0" />
                    </span>
                </Tooltip>
            </div>
        </div>
        <div class="blok">
            <img src="../assets/icony/rychlost.svg" alt="Rychlost" width="68">
            <Tooltip
                :zprava="`Za <b>neopravené</b> chyby je adekvátní penalizace. <br>${zaokrouhlit(info.rychlost)} CPM = ${info.rychlost == -1 ? '?' : zaokrouhlit(info.rychlost / 5)} WPM <br>(<b>CPM</b> = úhozů za minutu, <b>WPM</b> = slov za minutu)`"
                :sirka="200" style="width: 60%;">
                <span v-if="info.rychlost == -1" class="popis">Rychlost:<br><span class="nic">Zatím nic</span></span>
                <span v-else class="popis">
                    Rychlost:<br>
                    <AnimaceCisla class="cislo" :cislo="zaokrouhlit(info.rychlost)" /> CPM
                </span>
            </Tooltip>
        </div>
        <div class="blok" id="chyby">
            <div id="presnost">
                <img src="../assets/icony/terc.svg" alt="Přesnost">
                <Tooltip zprava="Přesnost zahrunuje chyby <b>opravené</b> i <b>neopravené</b>. (Státní zkouška počítá pouze neopravené.)"
                    :sirka="210">
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
                <div v-if="nejcastejsiChyby.length !== 0" id="pismena">
                    <div id="prvni">
                        <span v-for="znak, i in nejcastejsiChyby.slice(0, 2)" :key="znak.znak + znak.pocet"><span class="cisla">{{ i + 1 }}. </span>
                            <b :style="{ fontSize: znak.znak == ' ' ? '12px' : 'auto', fontWeight: znak.znak == ' ' ? '700' : '500' }">
                                {{ znak.znak == " " ? "┗━┛" : znak.znak }}
                            </b>
                        </span>
                    </div>
                    <div id="druhy">
                        <span v-for="znak, i in nejcastejsiChyby.slice(2)" :key="znak.znak + znak.pocet"><span class="cisla">{{ i + 3 }}. </span>
                            <b :style="{ fontSize: znak.znak == ' ' ? '12px' : 'auto', fontWeight: znak.znak == ' ' ? '700' : '500' }">
                                {{ znak.znak == " " ? "┗━┛" : znak.znak }}
                            </b>
                        </span>
                    </div>
                </div>
                <div v-else style="margin: 27px 0;">
                    <span>Žádné!</span>
                </div>
            </div>
        </div>
        <GrafStatistiky id="graf" />
    </div>

    <div id="bloky">
        <PrepinacTabu id="prepinac-tabu" :taby="[['celkem', 'Celkem'], ['dva tydny', 'Dva týdny'], ['dnes', 'Dnes']]" default-tab="celkem"
            ref="prepinac-tabu" />
        <div class="blok">
            <img src="../assets/icony/cas.svg" alt="Čas" width="68">
            <span class="popis">
                Čas strávený psaním: <br>
                <span v-if="info.cas >= 3600">
                    <AnimaceCisla class="cislo" :cislo="(info.cas - (info.cas % 3600)) / 3600" :desetina-mista="0" /> h
                </span>
                <span v-if="(info.cas % 3600) / 60 >= 1">
                    <AnimaceCisla class="cislo" :cislo="((info.cas % 3600) - (info.cas % 60)) / 60" :desetina-mista="0" /> min
                </span>
                <span v-if="info.cas % 60 >= 1 && info.cas < 3660">
                    <AnimaceCisla class="cislo" :cislo="info.cas % 60" :desetina-mista="0" /> s
                </span>
                <span v-if="info.cas == 0" class="nic">Zatím nic</span>
            </span>
        </div>
        <div class="blok">
            <img src="../assets/icony/cas.svg" alt="Čas" width="68">
            <span class="popis">
                Napsaných písmen: <br>
                <span v-if="info.cas >= 3600">
                    <AnimaceCisla class="cislo" :cislo="(info.cas - (info.cas % 3600)) / 3600" :desetina-mista="0" /> h
                </span>
                <span v-if="(info.cas % 3600) / 60 >= 1">
                    <AnimaceCisla class="cislo" :cislo="((info.cas % 3600) - (info.cas % 60)) / 60" :desetina-mista="0" /> min
                </span>
                <span v-if="info.cas % 60 >= 1 && info.cas < 3660">
                    <AnimaceCisla class="cislo" :cislo="info.cas % 60" :desetina-mista="0" /> s
                </span>
                <span v-if="info.cas == 0" class="nic">Zatím nic</span>
            </span>
        </div>
    </div>
</template>
<style scoped>
#graf {
    grid-column-start: 1;
    grid-column-end: 3;
}

#mini>.blok {
    width: 150px !important;
    padding: 8px;
}

#mini {
    display: flex;
    gap: 20px;
}

#prepinac-tabu {
    grid-row-start: 1;
    grid-column-start: 1;
    grid-column-end: 3;
    justify-self: center;
    margin: 30px 0 0 0;
    background-color: var(--tmave-fialova);
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
    grid-row-start: 1;
    grid-row-end: 3;
    grid-column-start: 2;
    gap: 0;
    justify-content: space-between;
    min-height: 230px;
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
    position: relative;
    overflow: hidden;
}

#progres>.popis {
    z-index: 1;
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

#bloky div img {
    height: 58px;
}

#bloky div h2 {
    font-weight: 500;
}

#graf {
    height: 190px;
    width: 660px;
}

#nacitani-pozadi {
    background-color: var(--fialova);
    padding: 0;
    position: absolute;
    left: 0;
    width: 100%;
    height: 50%;
    z-index: 1;
    animation: animate 3.5s ease-in-out infinite;

    transition: bottom 1.2s;
}

#nacitani-pozadi2 {
    background-color: var(--fialova);
    padding: 0;
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    z-index: 1;
    height: 0;

    transition: height 1.2s;
}

@keyframes animate {

    0%,
    100% {
        clip-path: polygon(0% 47.5%, 9% 47.5%, 18.25% 48.25%, 28.25% 50.25%, 36.83% 52.75%, 45.25% 56.25%, 54% 60%, 62% 62.75%, 70.28% 64.95%, 78.5% 66.5%, 86.46% 67%, 93.5% 67.25%, 100% 67%, 100% 100%, 0% 100%);
    }

    50% {
        clip-path: polygon(0% 67%, 9% 67.25%, 18.25% 67%, 28.25% 66.5%, 36.83% 64.95%, 45.25% 62.75%, 54% 60%, 62%56.25%, 70.28% 52.75%, 78.5% 50.25%, 86.46% 48.25%, 93.5% 47.5%, 100% 47.5%, 100% 100%, 0% 100%);
    }
}

@media screen and (max-width: 730px) {
    .blok {
        width: 320px;
        height: 105px;
    }

    #mini>.blok {
        width: 150px !important;
    }

    #bloky:first-of-type>.blok:nth-of-type(2) {
        order: -1;
    }

    #bloky {
        grid-template-columns: auto;
        justify-items: center;
    }

    #prepinac-tabu {
        grid-row-start: 1;
        grid-column-start: 1;
        grid-column-end: 2;
        margin-top: 40px;
    }

    #nadpis h2 {
        font-size: 22px;
    }

    #chyby {
        height: auto;
        grid-column-start: 1;
    }

    #graf {
        display: none;
    }
}
</style>