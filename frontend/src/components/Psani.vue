<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { onUnmounted } from "vue";
import Klavesnice from "../components/Klavesnice.vue";
import { Howl } from "howler";
import { MojeMapa, pridatOznameni } from "../utils";
import { useRoute } from "vue-router";
import { useHead } from "unhead";

useHead({
    meta: [
        { name: "robots", content: "noindex" }
    ]
})

const emit = defineEmits(["konec", "pise", "restart", "prodlouzit"])

interface TextInterface {
    id: number;
    znak: string;
    spatne: number;
    psat: boolean;
}
type TextTyp = TextInterface[][];

const props = defineProps({
    text: {
        type: Array as () => TextTyp, //whataheeeeeel
        required: true,
    },
    cas: {
        type: Number,
        required: true
    },
    klavesnice: {
        type: String,
        required: true
    },
    hideKlavesnice: {
        type: Boolean,
        required: true
    },
    nacitamNovej: {
        type: Boolean,
        required: true
    },
    delkaTextu: {
        type: Number,
        required: true
    },
    resetBtn: {
        type: Boolean,
        default: true
    }
})

const route = useRoute()

const counter = ref(0)
const counterSlov = ref(0)
const preklepy = ref(0)
const opravene = ref(0)
const timerZacatek = ref(0)
const cass = ref(0)
const textElem = ref<HTMLInputElement>()
let indexPosunuti = -1
const mistaPosunuti = ref([0, 0] as number[])
const chybyPismenka = new MojeMapa()

let predchoziZnak = ""

const zvukyZaply = ref(true)
const zvuky: Howl[] = []

const capslock = ref(false)
let interval: number

const celyPsani = ref()

let counterSpatneSvislaCara = 0

const casFormat = computed(() => {
    let zobrazeny = props.cas - cass.value
    return zobrazeny < 60 ? Math.floor(zobrazeny).toString() : `${Math.floor(zobrazeny / 60)}:${zobrazeny % 60 < 10 ? "0" + Math.floor(zobrazeny % 60).toString() : Math.floor(zobrazeny % 60)}`
})

const aktivniPismeno = computed(() => {
    if (counterSlov.value < props.text!.length) return props.text[counterSlov.value][counter.value]
    return { id: -1, znak: "", spatne: 0, psat: false }
})

watch(props.text, () => {
    if (!aktivniPismeno.value.psat && aktivniPismeno.value.id == 0) {
        counter.value++
        console.log("sus")
    }
})

onMounted(() => {
    document.addEventListener("keypress", klik) // je depracated ale je O TOLIK LEPSI ZE HO BUDU POUZIVAT PROSTE https://stackoverflow.com/questions/52882144/replacement-for-deprecated-keypress-dom-event
    document.addEventListener("keydown", specialniKlik)
    document.addEventListener("mousemove", enableKurzor)
})

onUnmounted(() => {
    document.removeEventListener("keypress", klik)
    document.removeEventListener("keydown", specialniKlik)
    document.removeEventListener("mousemove", enableKurzor)
})

function enableKurzor() {
    celyPsani.value.classList.remove("bez-kurzoru")
}

function capslockCheck(e: KeyboardEvent) { // TODO chtelo by to checknout hned po nacteni stranky ale nevim jestli to jde (spíš ne)
    capslock.value = e.getModifierState("CapsLock")
}

function nextPismeno() {
    if (aktivniPismeno.value.spatne === 1) preklepy.value++

    if (props.text[counterSlov.value].length - 1 === counter.value) { // posledni pismeno ve slovu
        counterSlov.value++
        counter.value = 0
    } else {
        counter.value++
    }

    if (!aktivniPismeno.value.psat) {
        nextPismeno()
    }
    emit("pise")
}

function backPismeno() {
    if (counter.value === 0) { // prvni pismeno ve slovu
        counterSlov.value--
        counter.value = props.text[counterSlov.value].length - 1
    } else {
        counter.value--
    }

    if (aktivniPismeno.value.spatne === 1) {
        preklepy.value--
        opravene.value++
    }

    if (!aktivniPismeno.value.psat) {
        backPismeno()
    }
    emit("pise")
}

function jeSHackem(key: string) {
    let velkym = aktivniPismeno.value.znak.toLocaleUpperCase() === aktivniPismeno.value.znak
    if (predchoziZnak === "ˇ") {
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ď" && (!velkym && key === "d" || velkym && key === "D")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ň" && (!velkym && key === "n" || velkym && key === "N")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ť" && (!velkym && key === "t" || velkym && key === "T")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ž" && (!velkym && key === "z" || velkym && key === "Z")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ř" && (!velkym && key === "r" || velkym && key === "R")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "č" && (!velkym && key === "c" || velkym && key === "C")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "š" && (!velkym && key === "s" || velkym && key === "S")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ě" && (!velkym && key === "e" || velkym && key === "E")) return true
    } else if (predchoziZnak === "´") {
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ó" && (!velkym && key === "o" || velkym && key === "O")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "é" && (!velkym && key === "e" || velkym && key === "E")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "í" && (!velkym && key === "i" || velkym && key === "I")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "á" && (!velkym && key === "a" || velkym && key === "A")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ý" && (!velkym && key === "y" || velkym && key === "Y")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ú" && (!velkym && key === "u" || velkym && key === "U")) return true
    } else if (predchoziZnak === "°") {
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ů" && (!velkym && key === "u" || velkym && key === "U")) return true
    } else {
        return false
    }
}

function klik(this: unknown, e: KeyboardEvent) {
    e.preventDefault() // ať to nescrolluje a nehazí nějaký stupid zkratky
    startTimer()
    loadZvuk()


    let hacek = jeSHackem(e.key)
    if (hacek) predchoziZnak = ""

    if (e.key === aktivniPismeno.value.znak || hacek) {
        if (zvukyZaply.value) zvuky[Math.floor(Math.random() * 2)].play()
        if (aktivniPismeno.value.spatne === 1) {
            aktivniPismeno.value.spatne = 2
        }
        nextPismeno()
        counterSpatneSvislaCara = 0
    } else {
        if (zvukyZaply.value) zvuky[3].play()
        aktivniPismeno.value.spatne = 1
        chybyPismenka.put(aktivniPismeno.value.znak)

        if (aktivniPismeno.value.znak === "|") {
            counterSpatneSvislaCara++
            if (counterSpatneSvislaCara >= 2) pridatOznameni(`Znak "|" je lehce problematický a jeho poloha se může lišit. Pokud máte dvouřádkový Enter, je označená klávesa posunutá vlevo od něj. Pokud jen nefunguje zvýrazněná klávesa, pravděpodobně se znak schovává vpravo od levého Shiftu. S pozdravem, Firu`, 15_000)
        }
        nextPismeno()
    }

    posunoutRadek()

    if (aktivniPismeno.value.id === -1) { // konec
        console.log("skončeno předčasně velký špatný")

        clearInterval(interval)
        calcCas() // naposledy
        document.removeEventListener("keypress", klik)
        document.removeEventListener("keydown", specialniKlik)
        emit("konec", opravene.value, preklepy.value, chybyPismenka)
        restart()
    }

    predchoziZnak = ""
    celyPsani.value.classList.add("bez-kurzoru")

    if (textViditelny.value[textViditelny.value.length - 1] == props.text[props.text.length - 1] && !props.nacitamNovej) emit("prodlouzit")
}

function posunoutRadek() {
    let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y
    let lastY = document.getElementById("p" + (aktivniPismeno.value.id - 1))?.getBoundingClientRect().y
    if (aktualniY! - lastY! > 30) {
        textElem.value!.classList.add("animace")
        indexPosunuti++
        if (indexPosunuti == 1) textElem.value!.style.top = "-2.35rem" // posunuti dolu
        else if (indexPosunuti > 1) textElem.value!.style.top = "-4.7rem" // posunuti dolu

        let slovoID = counterSlov.value
        setTimeout(() => {
            textElem.value!.classList.remove("animace")
            mistaPosunuti.value.push(slovoID)
            if (indexPosunuti > 0) textElem.value!.style.top = "-2.35rem" // posunuti dolu
        }, 200)
    }
}

function specialniKlik(e: KeyboardEvent) {
    capslockCheck(e)
    if (e.key === "Dead" && e.code === "Equal") { // kvůli macos :)
        e.preventDefault()
        if (e.shiftKey) predchoziZnak = "ˇ"
        else predchoziZnak = "´"
    } else if (e.key === "Dead" && e.code === "Backquote") {
        e.preventDefault()
        if (e.shiftKey) predchoziZnak = "°"
    } else if (e.key === "Backspace" || e.code === "Backspace" || e.keyCode == 8) {
        e.preventDefault()
        if (aktivniPismeno.value.id == 0 || props.nacitamNovej) return
        if (e.ctrlKey) { // tak dáme celé slovo pryč (Ctrl + Backspace zkratka)
            let lastY = document.getElementById("p" + (aktivniPismeno.value.id))?.getBoundingClientRect().y
            if (aktivniPismeno.value.znak == " ") backPismeno()
            if (counter.value == 0) backPismeno(); backPismeno()
            while (aktivniPismeno.value.znak != " ") {
                if (aktivniPismeno.value.id !== 0) {
                    backPismeno()
                } else {
                    break
                }
            }
            if (aktivniPismeno.value.id !== 0) nextPismeno()
            let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y
            if (lastY! - aktualniY! > 30) {
                indexPosunuti--
                textElem.value!.classList.add("animace")
                textElem.value!.style.top = "0rem"
                setTimeout(() => {
                    textElem.value!.classList.remove("animace")
                    mistaPosunuti.value.pop()
                    if (indexPosunuti > 0) textElem.value!.style.top = "-2.35rem" // posunuti dolu
                }, 200)
            }
        }
        else {
            backPismeno()
            let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y
            let lastY = document.getElementById("p" + (aktivniPismeno.value.id + 1))?.getBoundingClientRect().y
            if (lastY! - aktualniY! > 30) {
                indexPosunuti--
                textElem.value!.classList.add("animace")
                textElem.value!.style.top = "0rem"
                setTimeout(() => {
                    textElem.value!.classList.remove("animace")
                    mistaPosunuti.value.pop()
                    if (indexPosunuti > 0) textElem.value!.style.top = "-2.35rem" // posunuti dolu
                }, 200)
            }
        }
        if (zvukyZaply.value) zvuky[Math.floor(Math.random() * 2)].play()
    } else if (e.key == "Enter") {
        e.preventDefault()
        if (e.repeat) return
        resetTlacitko()
        animace()
    }
}

function startTimer() {
    if (timerZacatek.value === 0) {
        timerZacatek.value = Date.now()
        calcCas()
        interval = setInterval(calcCas, 200)
    }
}

function calcCas() {
    cass.value = (Date.now() - timerZacatek.value) / 1000

    if (props.cas - cass.value <= 0) {
        clearInterval(interval)
        document.removeEventListener("keypress", klik)
        document.removeEventListener("keydown", specialniKlik)
        emit("konec", opravene.value, preklepy.value, chybyPismenka, aktivniPismeno.value.id)
        restart()
    }
}

function toggleZvuk() {
    zvukyZaply.value = !zvukyZaply.value
    localStorage.setItem("pavouk_zvuk", zvukyZaply.value.toString())
}

function restart() {
    clearInterval(interval)
    timerZacatek.value = 0
    cass.value = 0
    counter.value = 0
    counterSlov.value = 0
    preklepy.value = 0
    indexPosunuti = -1
    textElem.value!.style.top = "0rem" // reset posunuti
    mistaPosunuti.value = [0, 0]
    chybyPismenka.clear()
    opravene.value = 0
}

async function loadZvuk() {
    let tmp = localStorage.getItem("pavouk_zvuk")
    if (tmp == null) {
        zvukyZaply.value = true
    } else {
        zvukyZaply.value = JSON.parse(tmp) === true // nejde to dat na jednu lajnu TS sus
    }

    zvuky.push(
        new Howl({
            src: "/zvuky/klik1.ogg",
            pool: 10,
            preload: true,
        }), new Howl({
            src: "/zvuky/klik2.ogg",
            pool: 10,
            preload: true,
        }),
        new Howl({
            src: "/zvuky/klik3.ogg",
            pool: 10,
            preload: true,
        }),
        new Howl({
            src: "/zvuky/miss.ogg",
            pool: 10,
            preload: true,
        })
    )
}

const textViditelny = computed(() => {
    return props.text.slice(mistaPosunuti.value[mistaPosunuti.value.length - 3], mistaPosunuti.value[mistaPosunuti.value.length - 2] + 42)
})

const rotaceStupne = ref(0)
function animace() {
    rotaceStupne.value -= 60
    setTimeout(() => { rotaceStupne.value = 0 }, 180)
}

const rotace = computed(() => {
    return `rotate(${rotaceStupne.value}deg)`
})

function resetTlacitko() {
    if (props.nacitamNovej) return
    emit("restart")
    restart()
}

defineExpose({ restart, aktivniPismeno })
</script>

<template>
    <div id="flex" ref="celyPsani">
        <div id="nabidka">
            <h2 id="cas">{{ casFormat }}</h2>
            <h2 :style="{ visibility: capslock ? 'visible' : 'hidden' }" id="capslock">CapsLock</h2>
            <h2 id="preklepy">Překlepy: {{ preklepy }}</h2>
        </div>

        <div id="ramecek">
            <div id="fade">
                <div id="text" ref="textElem" data-nosnippet>
                    <div class="slovo" v-for="s, i in textViditelny" :key="i">
                        <div v-for="p in s" :key="p.id" class="pismeno" :id="'p' + p.id"
                            :class="{ podtrzeni: p.id === aktivniPismeno.id, 'spatne-pismeno': p.spatne === 1 && aktivniPismeno.id > p.id, 'opravene-pismeno': p.spatne === 2 && aktivniPismeno.id > p.id, 'spravne-pismeno': (!p.spatne && aktivniPismeno.id > p.id) || !p.psat }">
                            {{ (p.znak !== " " ? p.znak : p.spatne && p.id < aktivniPismeno.id ? "_" : "&nbsp;") }} </div>
                        </div>
                    </div>
                </div>
            </div>

            <Transition>
                <Klavesnice v-if="klavesnice != ''" :typ="klavesnice" :aktivniPismeno="aktivniPismeno.znak" :rozmazat="hideKlavesnice"
                    :cekame="(aktivniPismeno.id == 0 || aktivniPismeno.id == -1) && cass == 0" />
            </Transition>
            <Transition>
                <div v-if="klavesnice != '' && props.resetBtn" id="reset-btn" @click="resetTlacitko(); animace()"
                    :class="{ schovat: route.fullPath == '/prvni-psani' }">
                    <img :style="{ transform: rotace }" src="../assets/icony/reset.svg" alt="Nastavení">
                </div>
            </Transition>

            <div id="zvuk-btn" @click="toggleZvuk">
                <img v-if="zvukyZaply" style="margin-top: 1px;" class="zvuk-icon" src="../assets/icony/zvukOn.svg" alt="Zvuky jsou zapnuté">
                <img v-else style="margin-left: 1px;" class="zvuk-icon" src="../assets/icony/zvukOff.svg" alt="Zvuky jsou vypnuté">
            </div>
        </div>

        <audio> <!-- načteme soubory do cache, aby je měl howler rychle -->
            <source src="/zvuky/klik1.ogg" type="audio/ogg">
            <source src="/zvuky/klik2.ogg" type="audio/ogg">
            <source src="/zvuky/klik3.ogg" type="audio/ogg">
            <source src="/zvuky/miss.ogg" type="audio/ogg">
        </audio>
</template>

<style scoped>
.schovat {
    opacity: 0;
    cursor: auto !important;
}

/*  eslint-disable-next-line vue-scoped-css/no-unused-selector */
.bez-kurzoru {
    cursor: none;
}

.v-enter-active {
    transition: 0.3s !important;
}

.v-enter-from {
    opacity: 0;
}

.zvuk-icon {
    width: 45px;
    height: 35px;
    margin-top: 1px;
}

#zvuk-btn {
    position: fixed;
    right: 30px;
    bottom: 25px;
    background-color: var(--tmave-fialova);
    border-radius: 100px;
    width: 55px;
    height: 55px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background-color 0.1s;
    user-select: none;
}

#reset-btn {
    position: relative;
    width: 55px;
    height: 55px;
    background-color: var(--tmave-fialova);
    border-radius: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
    left: 385px;
    top: -154px;
    cursor: pointer;
    transition: background-color 0.1s;
    user-select: none;
}

#reset-btn img {
    width: 30px;
    padding-bottom: 1px;
    transition-duration: 0.2s;
    transition-timing-function: ease-out;
}

#reset-btn:hover,
#zvuk-btn:hover {
    background-color: var(--fialova);
}

#flex {
    display: flex;
    flex-direction: column;
    align-items: center;
}

#nabidka {
    margin: 20px 0 6px 0;
    width: var(--sirka-textoveho-pole);
}

#cas {
    float: left;
    width: 150px;
    display: block;
    text-align: left;
}

#preklepy {
    float: right;
    display: block;
    width: 150px;
    text-align: right;
}

#capslock {
    display: inline-block;
    color: red;
    font-weight: bold !important;
}

#ramecek {
    padding: 10px;
    height: 190px;
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
    overflow: hidden;
    user-select: none;
}

#text {
    display: flex;
    flex-wrap: wrap;
    position: relative;
    top: 0em;
}

/*  eslint-disable-next-line vue-scoped-css/no-unused-selector */
.animace {
    transition: ease 0.2s;
}

#fade {
    mask-image: linear-gradient(180deg, var(--tmave-fialova) 75%, transparent 97%);
    -webkit-mask-image: linear-gradient(180deg, var(--tmave-fialova) 75%, transparent 97%);
    height: 107%;
}

.slovo {
    display: flex;
    flex-wrap: nowrap;
}

.pismeno {
    border-radius: 3px;
    display: inline-flex;
    font-family: "Red Hat Mono", monospace;
    font-weight: 400;
    font-size: 1.56rem;
    line-height: 2.2rem;
    width: 17px;
    text-decoration: none;
    padding: 0 1px;
    margin-right: 1px;
    border-bottom: 3px solid rgba(255, 255, 255, 0);
    /* aby se nedojebala vyska lajny když jdu na dalsi radek*/
    color: var(--bila);
}

#nabidka h2 {
    font-size: 1.2em;
    font-weight: 300;
}

.spravne-pismeno {
    color: #8c8c8c;
}

.podtrzeni {
    border-bottom: 3px solid var(--bila);
    border-radius: 0;
}

.spatne-pismeno {
    color: #ff0000;
}

.opravene-pismeno {
    color: #b1529c;
}
</style>