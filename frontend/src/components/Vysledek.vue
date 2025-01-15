<script setup lang="ts">
import axios from "axios";
import { onMounted, onUnmounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { getCisloPochvaly, getToken, MojeMapa, napovedaKNavigaci } from "../utils";
import { levelyRychlosti } from "../stores";
import Tooltip from "./Tooltip.vue";
import AnimaceCisla from "../components/AnimaceCisla.vue";

const emit = defineEmits(["restart"])

const props = defineProps({
    preklepy: {
        type: Number,
        default: 0
    },
    opravenych: {
        type: Number,
        default: 0
    },
    delkaTextu: {
        type: Number,
        default: 1
    },
    cas: {
        type: Number,
        default: 1
    },
    pismena: {
        type: String,
        default: ""
    },
    nejcastejsiChyby: {
        type: MojeMapa,
        default: new MojeMapa()
    },
    cislo: String,
    posledni: {
        type: Boolean,
        default: true
    }
})

let rychlost = ((props.delkaTextu - (10 * props.preklepy)) / props.cas) * 60
const route = useRoute()
const router = useRouter()
const pochavly = ["Dobrá práce!", "Bravo!", "Pěkná práce!", "Skvělá práce!", "Výborně!", "Parádní!", "Skvělý výsledek!", "Paráda!", "Hezký!", "Super výkon!", "Parádní výkon!", "Skvělý výkon!"]
const vsechnyHodnoceni = [
    ["Pavouci jásají z tvé šikovnosti.", "Avšak i když už jsi profík, vždy je kam se posouvat.", "Píšeš krásně jako pavouk."], // parádní
    ["Ale můžeš ještě zapracovat na rychlosti.", "Leda rychlost jde ještě zlepšovat.", "Máš to za sebou ale rychlost můžeš ještě zlepšit."], // dobrý ale rychlost by šla zlepšit
    ["Ale můžeš ještě zapracovat na přesnosti.", "Leda přesnost jde ještě zlepšovat.", "Máš to za sebou ale přesnost můžeš ještě zlepšit."], // dobrý ale přesnost by šla zlepšit
    ["Ale můžeš se ještě zlepšit.", "Máš to za sebou ale ještě je kam růst."], // dobrý ale oboje jde zlepsit
    ["Dej tomu ale ještě chvíli. Jde psát i trochu rychleji.", "Zatím ale moc pomalé.", "Musíš ale ještě trochu zrychlit."], // rychlost není dostatečná
    ["Dej tomu ale ještě chvíli. Jde dělat i méně chyb.", "Zatím hodně chybuješ.", "Zaměř se i na přesnost, ještě to není ono."], // přesnost není dostatečná
    ["Dej tomu ale ještě chvíli. Zatím ti to moc nejde.", "Zkus to ale ještě jednou."]
]
const hodnoceni = ref("")

let presnost = (props.delkaTextu - props.preklepy - props.opravenych) / props.delkaTextu * 100
const nejcastejsiChybyTop3 = ref()

function reset() {
    emit("restart")
}

function dalsi() {
    if (props.cislo == undefined) return
    let r = route.path.split("/")
    r.pop()
    let c = r.join("/")
    if (props.posledni) router.push(c) // /lekce/pismena
    else router.push(c + "/" + (parseInt(props.cislo) + 1).toString()) // /lekce/pismena/cislo
}

function random(list: Array<string>) {
    return list[(Math.floor(Math.random() * list.length))]
}

onMounted(() => {
    hodnoceni.value = random(pochavly) + " " + random(vsechnyHodnoceni[getCisloPochvaly(rychlost, presnost)])

    document.addEventListener("keydown", e1)

    if (props.cislo == "prvni-psani") {
        hodnoceni.value = "Píšeš krásně, ale tohle byl jen začátek..."
        return
    }

    nejcastejsiChybyTop3.value = props.nejcastejsiChyby.top(3)

    if (props.pismena == "") { // je to procvicovani / test takze posilame jinam
        let cislo = props.cislo
        if (props.cislo == "test-psani") cislo = "0" // test psani

        axios.post("/dokonceno-procvic/" + cislo, {
            "neopravenePreklepy": props.preklepy,
            "cas": props.cas,
            "delkaTextu": props.delkaTextu,
            "nejcastejsiChyby": Object.fromEntries(props.nejcastejsiChyby)
        }, {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }).catch(function (e) {
            console.log(e)
        })
        return
    }

    if (props.pismena == "pracepraceprace") { // je to práce
        let id = props.cislo

        axios.post("/skola/dokoncit-praci/" + id, {
            "neopravenePreklepy": props.preklepy,
            "cas": props.cas,
            "delkaTextu": props.delkaTextu,
            "nejcastejsiChyby": Object.fromEntries(props.nejcastejsiChyby)
        }, {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }).catch(function (e) {
            console.log(e)
        })
        return
    }

    // jsme ve cviceni
    axios.post("/dokonceno/" + encodeURIComponent(props.pismena) + "/" + props.cislo, {
        "neopravenePreklepy": props.preklepy,
        "cas": props.cas,
        "delkaTextu": props.delkaTextu,
        "nejcastejsiChyby": Object.fromEntries(props.nejcastejsiChyby)
    }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).catch(function (e) {
        console.log(e)
    })
})

onUnmounted(() => {
    document.removeEventListener("keydown", e1)
})

function e1(e: KeyboardEvent) {
    if (e.key == "Delete") {
        e.preventDefault()
        if (props.cislo == 'prvni-psani' || props.pismena == 'pracepraceprace') return
        reset()
    } else if (e.key == "ArrowRight") {
        e.preventDefault()
        if (route.path.split('/')[1] != 'lekce') return
        dalsi()
    } else if (e.key == "Tab") {
        e.preventDefault()
        napovedaKNavigaci()
    }
}
</script>

<template>
    <div id="bloky" style="margin-top: 25px;">
        <div id="hodnoceni" class="blok" :style="{ width: cislo == 'prvni-psani' ? '400px' : '' }">
            <Tooltip :zprava="`Pro získání 3 hvězd je potřeba dosánout rychlosti <b>min. ${levelyRychlosti[2]} CPM</b>. Hodně štěstí!`"
                :vzdalenost="10" :sirka="250">
                <div id="hvezdy">
                    <img v-if="rychlost >= levelyRychlosti[0]" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
                    <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
                    <img v-if="rychlost >= levelyRychlosti[1]" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
                    <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
                    <img v-if="rychlost >= levelyRychlosti[2]" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
                    <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
                </div>
            </Tooltip>
            <div style="display: flex; align-items: center; height: 100%;">
                <h3 style="font-weight: 300; margin: 0">{{ hodnoceni }}</h3>
            </div>
        </div>
        <div v-if="cislo !== 'prvni-psani'" class="blok" id="chyby">
            <h2>Nejčastější chyby</h2>
            <hr>
            <div v-if="nejcastejsiChyby.size !== 0">
                <ol>
                    <li v-for="znak in nejcastejsiChybyTop3" :key="znak.znak">
                        <span :style="{ fontSize: znak.znak == ' ' ? '20px' : 'auto', paddingTop: znak.znak == ' ' ? '14px' : '0px'}">{{ znak.znak == " " ? "┗━┛" : znak.znak }}</span>
                    </li>
                </ol>
                <ul>
                    <li v-for="znak in nejcastejsiChybyTop3" :key="znak">
                        <span v-if="znak.pocet > 0">{{ znak.pocet }}</span>
                    </li>
                </ul>
            </div>
            <h3 v-else style="height: 100%; display: flex; justify-content: center; align-items: center; margin-top: 0;">
                <span>Žádné!</span>
            </h3>
        </div>
    </div>

    <div id="bloky">
        <div class="blok">
            <Tooltip
                zprava="Za <b>neopravené</b> chyby je adekvátní <b>penalizace</b>. Chybu opravíš pomocí klávesy <span class='klavesa-v-textu-mensi'>Backspace</span>"
                :sirka="180" :vzdalenost="6">
                <AnimaceCisla class="cislo" :cislo="rychlost > 0 ? Math.round(rychlost * 10) / 10 : 0" :desetina-mista="0" />
            </Tooltip>
            <hr>
            <p class="jednotka">CPM / úhozů</p>
            <p class="jednotka">&zwnj;</p>
            <h3>Rychlost</h3>
        </div>
        <div class="blok">
            <Tooltip zprava="Přesnost zahrnuje chyby <b>opravené</b> i <b>neopravené</b>" :sirka="200" :vzdalenost="6">
                <AnimaceCisla class="cislo" :cislo="Math.round(presnost * 10) / 10 <= 0 ? 0 : Math.round(presnost * 10) / 10" />
                <span class="procento">%</span>
            </Tooltip>
            <hr>
            <p v-if="preklepy == 1" class="jednotka">{{ preklepy }} neopravený</p>
            <p v-else-if="preklepy >= 2 && preklepy <= 4" class="jednotka">{{ preklepy }} neopravené</p>
            <p v-else-if="preklepy >= 5 || preklepy == 0" class="jednotka">{{ preklepy }} neopravených</p>
            <p v-if="opravenych == 1" class="jednotka">{{ opravenych }} opravený</p>
            <p v-else-if="opravenych >= 2 && opravenych <= 4" class="jednotka">{{ opravenych }} opravené</p>
            <p v-else-if="opravenych >= 5 || opravenych == 0" class="jednotka">{{ opravenych }} opravených</p>
            <h3>Přesnost</h3>
        </div>
        <div class="blok">
            <AnimaceCisla class="cislo" :cislo="cas < 60 ? Math.round(cas * 10) / 10 : Math.round(cas / 60)" :desetina-mista="0" />
            <hr>
            <p class="jednotka">{{ cas < 60 ? "Sekund" : (cas == 60) ? "Minuta" : (cas <= 240) ? "Minuty" : "Minut" }}</p>
                    <p class="jednotka">&zwnj;</p>
                    <h3>Čas</h3>
        </div>
    </div>

    <div id="tlacitka-kontainer">
        <span v-if="props.cislo == 'prvni-psani'">Líbí se ti aplikace?</span>
        <button v-if="props.cislo == 'prvni-psani'" class="tlacitko" @click="router.push('/registrace')">Vytvořit účet</button>
        <button v-if="props.pismena == 'pracepraceprace'" class="tlacitko" @click="router.push('/trida')">Zpět do třídy</button>

        <Tooltip v-if="props.cislo != 'prvni-psani' && props.pismena != 'pracepraceprace'"
            zprava="Také pomocí klávesy <span class='klavesa-v-textu-mensi'>Delete</span>" :sirka="130">
            <button class="tlacitko" @click="reset">Zkusit znovu</button>
        </Tooltip>

        <Tooltip v-if="route.path.split('/')[1] == 'lekce'"
            zprava="Také pomocí klávesy <span class='klavesa-v-textu-mensi'><img src='/src/assets/icony/sipkaL.svg' alt='Šipka' class='klav-sipka' style='transform: scaleX(-1) translateY(2px) translateX(-1px);'></span>"
            :sirka="130">
            <button class="tlacitko" @click="dalsi()">Pokračovat</button>
        </Tooltip>
    </div>

</template>

<style scoped>
li {
    font-size: 18px;
    opacity: 70%;
    height: 23.2px;
    display: flex;
    align-items: center;
}

li:first-child {
    font-size: 29px;
    margin-bottom: 4px;
    opacity: 100%;
    height: 38.4px;
}

ol span {
    font-weight: 500;
    font-family: "Red Hat Mono", monospace;
}

ol,
ul {
    display: flex;
    flex-direction: column;
    align-items: center;
    list-style-type: none;
}

#chyby {
    max-height: 155px;
    padding-bottom: 5px;
}

#chyby h2 {
    font-size: 26px;
    margin-bottom: 8px;
    font-size: 19px;
    font-family: 'Montserrat';
}

#chyby div {
    display: flex;
    align-items: center;
    gap: 44px;
    height: 65%;
}

.hvezda {
    width: 50px;
    height: 50px;
}

#hvezdy :nth-child(2) {
    position: relative;
    top: -5px;
}

#hvezdy {
    margin-top: 5px;
}

#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
    margin-top: 20px;
}

#hodnoceni {
    width: 380px;
    display: flex;
    gap: 10px;
    height: auto;
    padding: 20px;
}

.blok div {
    display: flex;
    align-items: baseline;
    justify-content: center;
    gap: 5px;
}

.procento {
    font-size: 24px;
    position: relative;
    left: 4px;
}

.blok {
    color: var(--bila);
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 220px;
    background-color: var(--tmave-fialova);
    transition-duration: 0.2s;
    padding: 15px 15px 20px 15px;
}

.blok .cislo {
    font-size: 40px;
    font-weight: 500;
    font-family: "Red Hat Mono";
}

.jednotka {
    font-size: 14px;
}

.blok hr {
    width: 180px;
    align-self: center;
    position: relative;
    top: -3px;
    margin-bottom: 4px;
}

.blok h3 {
    font-weight: 500;
    margin-top: 12px;
}

#tlacitka-kontainer {
    display: inline-flex;
    gap: 20px;
    margin-top: 20px;
    align-items: center;
}

#tlacitka-kontainer .tlacitko {
    margin: 0 !important;
}
</style>