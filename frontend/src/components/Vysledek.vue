<script setup lang="ts">
import axios from 'axios';
import { onMounted, onUnmounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getToken, napovedaKNavigaci } from '../utils';
import { levelyRychlosti, levelyPresnosti } from '../stores';

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
    casF: String,
    pismena: {
        type: String,
        default: ""
    },
    nejcastejsiChyby: {
        type: Array<any>,
        default: ["prvni-psani"]
    },
    cislo: String,
    posledni: Boolean
})

let rychlost = Math.round((props.delkaTextu / props.cas) * 60 * 10) / 10
const route = useRoute()
const router = useRouter()
const pochavly = ["Dobrá práce!", "Bravo!", "Pěkná práce!", "Skvělá práce!", "Výborně!", "Parádní!", "Skvělý výsledek!", "Paráda!", "Hezký!", "Super výkon!", "Parádní výkon!", "Skvělý výkon!"]
const vsechnyHodnoceni = [
    ["Pavouci jásají z tvé šikovnosti.", "Avšak i když už jsi profík, vždy je kam se posouvat.", "Píšeš krásně jako pavouk."], // parádní
    ["Ale můžeš ještě zapracovat na rychlosti.", "Leda rychlost jde ještě zlepšovat.", "Cvičení máš hotové ale rychlost můžeš ještě zlepšit."], // dobrý ale rychlost by šla zlepšit
    ["Ale můžeš ještě zapracovat na přesnosti.", "Leda přesnost jde ještě zlepšovat.", "Cvičení máš hotové ale přesnost můžeš ještě zlepšit."], // dobrý ale přesnost by šla zlepšit
    ["Ale můžeš se ještě zlepšit.", "Cvičení máš hotové ale ještě je kam růst."], // dobrý ale oboje jde zlepsit
    ["Dej tomu jěště chvíli. Jde psát i trochu rychleji.", "Zatím ale moc pomalé.", "Musíš ale ještě trochu zrychlit."], // rychlost není dostatečná
    ["Dej tomu jěště chvíli. Jde dělat i méně chyb.", "Zatím hodně chybuješ.", "Zaměř se i na přesnost, ještě to není ono."], // přesnost není dostatečná
    ["Dej tomu jěště chvíli. Zatím ti to moc nejde.", "Zkus to ale ještě jednou."]
]
const hodnoceni = ref("")
const hvezdy = ref(0)

let presnost = (props.delkaTextu - props.preklepy) / props.delkaTextu * 100

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
    hodnoceni.value += random(pochavly) + " "
    if (rychlost >= levelyRychlosti[1] && presnost >= levelyPresnosti[1]) { // paradni
        hodnoceni.value += random(vsechnyHodnoceni[0])
        hvezdy.value = 3
    } else if (rychlost >= levelyRychlosti[0] && rychlost < levelyRychlosti[1] && presnost >= levelyPresnosti[1]) { // rychlost muze byt lepsi
        hodnoceni.value += random(vsechnyHodnoceni[1])
        hvezdy.value = 2
    } else if (presnost >= levelyPresnosti[0] && presnost < levelyPresnosti[1] && rychlost >= levelyRychlosti[1]) { // presnost muze byt lepsi
        hodnoceni.value += random(vsechnyHodnoceni[2])
        hvezdy.value = 2
    } else if (presnost >= levelyPresnosti[0] && presnost < levelyPresnosti[1] && rychlost >= levelyRychlosti[0] && rychlost <= levelyRychlosti[1]) { // oboje muze byt lepsi
        hodnoceni.value += random(vsechnyHodnoceni[3])
        hvezdy.value = 1
    } else if (rychlost < levelyRychlosti[0] && presnost < levelyPresnosti[0]) { // oboje bad
        hodnoceni.value += random(vsechnyHodnoceni[6])
        hvezdy.value = 0
    } else if (rychlost < levelyRychlosti[0]) { // rychlost bad
        hodnoceni.value += random(vsechnyHodnoceni[4])
        hvezdy.value = 0
    } else if (presnost < levelyPresnosti[0]) { // presnost bad
        hodnoceni.value += random(vsechnyHodnoceni[5])
        hvezdy.value = 0
    }

    document.addEventListener('keydown', e1)

    if (props.cislo == "" || props.cislo == 'test-psani') return // je to procvicovani / test takze neposilame

    if (props.cislo == "prvni-psani") {
        hodnoceni.value = "Píšeš krásně, ale tohle byl jen začátek..."
        return
    }

    axios.post('/dokonceno/' + encodeURIComponent(props.pismena) + '/' + props.cislo, {
        "cpm": rychlost,
        "preklepy": props.preklepy,
        "cas": props.cas,
        "delkaTextu": props.delkaTextu
    }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).catch(function (e) {
        console.log(e)
    })
})

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
})

function e1(e: KeyboardEvent) {
    if (e.key == " ") {
        e.preventDefault()
        reset()
    } else if (e.key == "ArrowRight" || e.key === "Enter") {
        e.preventDefault()
        dalsi()
    } else if (e.key == 'Tab') {
        e.preventDefault()
        napovedaKNavigaci()
    }
}
</script>

<template>
    <div id="bloky" style="margin-top: 25px;">
        <div id="hodnoceni" class="blok" :style="{width: cislo == 'prvni-psani' ? '400px' : ''}">
            <div id="hvezdy">
                <img v-if="hvezdy >= 1" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
                <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
                <img v-if="hvezdy >= 2" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
                <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
                <img v-if="hvezdy == 3" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
                <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
            </div>
            <div style="display: flex; align-items: center; height: 100%;">
                <h3 style="font-weight: 300; margin: 0">{{ hodnoceni }}</h3>
            </div>
        </div>
        <div v-if="cislo !== 'prvni-psani'" class="blok" id="chyby">
            <h2>Nejčastější chyby</h2>
            <hr>
            <div v-if="nejcastejsiChyby.length !== 0">
                <ol>
                    <li v-for="znak in nejcastejsiChyby"><span>{{ znak[0] == " " ? "_" : znak[0] }}</span></li>
                </ol>
                <ul>
                    <li v-for="znak in nejcastejsiChyby"><span v-if="znak[1] > 0">{{ znak[1] }}</span></li>
                </ul>
            </div>
            <h3 v-else style="margin-top: 32px;">Žádné!</h3>
        </div>
    </div>

    <div id="bloky">
        <div class="blok">
            <h2>{{ rychlost }}</h2>
            <hr>
            <p class="jednotka">CPM / úhozů</p>
            <p class="jednotka">&zwnj;</p>
            <h3>Rychlost</h3>
        </div>
        <div class="blok">
            <h2>{{ Math.round(presnost * 10) / 10 }}<span class="procento">%</span></h2>
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
            <h2>{{ cas < 60 ? Math.round(cas * 10) / 10 : `${Math.floor(cas / 60)}:${Math.floor(cas % 60 * 10) / 10 < 10
                    ? "0" + Math.floor(cas % 60 * 10) / 10 : Math.floor(cas % 60 * 10) / 10}` }}</h2>
                    <hr>
                    <p class="jednotka">{{ cas < 60 ? "Sekund" : "MM:SS" }}</p>
                            <p class="jednotka">&zwnj;</p>
                            <h3>Čas</h3>
        </div>
    </div>

    <div v-if="props.cislo != 'prvni-psani' && props.cislo != 'test-psani'" id="tlacitka_kontainer">
        <button class="tlacitko" @click="reset">Zkusit znovu</button>
        <button class="tlacitko" @click="dalsi()">Pokračovat</button>
    </div>
    <div v-else-if="props.cislo == 'test-psani'" id="tlacitka_kontainer">
        <button class="tlacitko" @click="reset">Zkusit znovu</button>
    </div>
    <div v-else id="tlacitka_kontainer" style="align-items: center;">
        <span>Líbí se ti aplikace?</span>
        <button class="tlacitko" @click="router.push('/registrace')">Vytvořit účet</button>
    </div>
</template>

<style scoped>
li {
    font-size: 1.1em;
    opacity: 70%;
}

li:first-child {
    font-size: 1.8em;
    margin-bottom: 4px;
    opacity: 100%;
}

ol li:first-child span {
    font-weight: 700;
}

ol li span {
    font-weight: 500;
}

ol, ul {
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
    font-size: 1.6em;
    margin-bottom: 8px;
    font-size: 1.2em;
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
    font-size: 0.7em;
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

.blok h2 {
    font-size: 40px;
    font-weight: 500;
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

#tlacitka_kontainer {
    display: inline-flex;
    gap: 20px;
    margin-top: 20px;
}

#tlacitka_kontainer .tlacitko {
    margin: 0 !important;
}
</style>