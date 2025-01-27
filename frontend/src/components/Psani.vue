<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef, watch } from "vue";
import { onUnmounted } from "vue";
import Klavesnice from "../components/Klavesnice.vue";
import { Howl } from "howler";
import { MojeMapa, pridatOznameni } from "../utils";
import { useRoute } from "vue-router";
import { useHead } from "unhead";
import Tooltip from "../components/Tooltip.vue";

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
const textElem = useTemplateRef("textElem")
let indexPosunuti = -1
const mistaPosunuti = ref([0, 0] as number[])
const chybyPismenka = new MojeMapa()

const fullHideKlavesnice = ref(false)

const zvukyZaply = ref(true)
const zvuky: Howl[] = []

const capslock = ref(false)
let interval: number

const celyPsani = useTemplateRef("celyPsani")
const input = useTemplateRef("input")

let inputPredchoziDelka = 0
const unfocused = ref(false)

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
    }
})

onMounted(() => {
    document.addEventListener("keydown", specialniKlik)
    document.addEventListener("mousemove", enableKurzor)
    loadZvuk()
    loadHideKlavesnci()

    input.value?.focus()
    document.addEventListener("click", checkFocus)
})

onUnmounted(() => {
    document.removeEventListener("keydown", specialniKlik)
    document.removeEventListener("mousemove", enableKurzor)
    clearTimeout(timeoutID)
    document.removeEventListener("click", checkFocus)
})

function enableKurzor() {
    celyPsani.value?.classList.remove("bez-kurzoru")
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
    checkJestliPise()
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
    checkJestliPise()
    emit("pise")
}

function klik(e: Event) {
    e.preventDefault() // ať to nescrolluje a nehazí nějaký stupid zkratky
    startTimer()

    if (!(e instanceof InputEvent)) return // typescript je sus, nemůžu dát do parametru rovnou InputEvent https://github.com/microsoft/TypeScript/issues/39925
    if (e.data == "") return
    //                                kvůli macos ->                             kvůli linuxu ->
    if (e.inputType == "insertText" || e.inputType == "insertFromComposition" || (e.inputType == "insertCompositionText" && !e.isComposing)) {
        if (e.data === aktivniPismeno.value.znak) {
            if (zvukyZaply.value) zvuky[Math.floor(Math.random() * 2)].play()
            if (aktivniPismeno.value.spatne === 1) {
                aktivniPismeno.value.spatne = 2
            }
            counterSpatneSvislaCara = 0
        } else {
            if (zvukyZaply.value) zvuky[3].play()
            aktivniPismeno.value.spatne = 1
            chybyPismenka.put(aktivniPismeno.value.znak)

            if (aktivniPismeno.value.znak === "|") {
                counterSpatneSvislaCara++
                if (counterSpatneSvislaCara >= 2) pridatOznameni(`Znak "|" je lehce problematický a jeho poloha se může lišit. Pokud máte dvouřádkový Enter, je označená klávesa posunutá vlevo od něj. Pokud jen nefunguje zvýrazněná klávesa, pravděpodobně se znak schovává vpravo od levého Shiftu. S pozdravem, Firu`, 15_000)
            }
        }
        nextPismeno()
        posunoutRadek()
    } else if (e.inputType == "deleteContentBackward") {
        backPismeno()
        vratitRadek()
    } else if (e.inputType == "deleteWordBackward") { // tak dáme celé slovo pryč (Ctrl + Backspace zkratka)
        for (let i = 0; i < inputPredchoziDelka - input.value!.value.length; i++) {
            backPismeno()
            vratitRadek()
        }
    }

    if (aktivniPismeno.value.id === -1) { // konec
        console.log("skončeno předčasně velký špatný")

        clearInterval(interval)
        calcCas() // naposledy
        document.removeEventListener("keydown", specialniKlik)
        document.removeEventListener("click", checkFocus)
        emit("konec", opravene.value, preklepy.value, chybyPismenka)
        restart()
    }

    inputPredchoziDelka = input.value!.value.length
    celyPsani.value?.classList.add("bez-kurzoru")
    if (textViditelny.value[textViditelny.value.length - 1] == props.text[props.text.length - 1] && !props.nacitamNovej) emit("prodlouzit")
}

async function posunoutRadek() {
    let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y
    let lastY = document.getElementById("p" + (aktivniPismeno.value.id - 1))?.getBoundingClientRect().y
    if (lastY == undefined || aktualniY == undefined) return
    if (aktualniY - lastY > 30) {
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

async function vratitRadek() {
    let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y
    let lastY = document.getElementById("p" + (aktivniPismeno.value.id + 1))?.getBoundingClientRect().y
    if (lastY == undefined || aktualniY == undefined) return
    if (lastY - aktualniY > 30) {
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

function specialniKlik(e: KeyboardEvent) {
    capslockCheck(e)
    if (e.key.slice(0, 5) == "Arrow") { // vypnout sipky
        e.preventDefault()
        return
    } else if (e.key == "Delete") {
        e.preventDefault()

        if (route.fullPath.split("/")[1] == "prace" || route.fullPath == "/prvni-psani") return
        if (e.repeat) return

        resetTlacitko()
        animace()
    } else if (e.code == "Space" && unfocused.value) {
        e.preventDefault()
        input.value?.focus()
        unfocused.value = false
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
        document.removeEventListener("keydown", specialniKlik)
        document.removeEventListener("click", checkFocus)
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
    mistaPosunuti.value = [0, 0]
    chybyPismenka.clear()
    opravene.value = 0
    if (input.value != null) input.value.value = '' // nemusi byt asi
    clearTimeout(timeoutID)

    for (let i = 0; i < props.text.length; i++) {
        for (let j = 0; j < props.text[i].length; j++) {
            // eslint-disable-next-line vue/no-mutating-props
            props.text[i][j].spatne = 0;
        }
    }

    if (textElem.value?.hasAttribute("style")) textElem.value.style.top = "0rem" // reset posunuti
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

let timeoutID: number
const prestalPsat = ref(false)
async function checkJestliPise() {
    clearTimeout(timeoutID)
    prestalPsat.value = false
    timeoutID = setTimeout(() => {
        prestalPsat.value = true
        restart()
    }, 10000) // 10s
}

watch(fullHideKlavesnice, () => {
    localStorage.setItem("pavouk_hide_klavesnice", JSON.stringify(fullHideKlavesnice.value))
})

async function loadHideKlavesnci() {
    let x = localStorage.getItem("pavouk_hide_klavesnice")
    if (x == null) {
        localStorage.setItem("pavouk_hide_klavesnice", JSON.stringify(fullHideKlavesnice.value))
    } else {
        fullHideKlavesnice.value = JSON.parse(x)
    }
}

function checkFocus() {
    unfocused.value = document.activeElement !== input.value
}

function focusInput() {
    input.value?.focus()
}

defineExpose({ restart, aktivniPismeno, fullHideKlavesnice, focusInput })
</script>

<template>
    <div id="flex" ref="celyPsani">
        <div id="nabidka" :class="{ unfocused: unfocused }">
            <h2 id="cas">{{ casFormat }}</h2>
            <h2 :style="{ visibility: capslock ? 'visible' : 'hidden' }" id="capslock">CapsLock</h2>
            <h2 id="preklepy">Překlepy: {{ preklepy }}</h2>
        </div>

        <div id="ramecek" @click="input?.focus()" :class="{ unfocused: unfocused }">
            <div id="fade">
                <div id="text" ref="textElem" data-nosnippet>
                    <div class="slovo" v-for="s, i in textViditelny" :key="i">
                        <div v-for="p in s" :key="p.id" class="pismeno" :id="'p' + p.id" :class="{
                            'podtrzeni': p.id === aktivniPismeno.id,
                            'spatne-pismeno': p.spatne === 1 && aktivniPismeno.id > p.id,
                            'opravene-pismeno': p.spatne === 2 && aktivniPismeno.id > p.id,
                            'spravne-pismeno': (!p.spatne && aktivniPismeno.id > p.id) || !p.psat
                        }">

                            {{ (p.znak !== " " ? p.znak : p.spatne && p.id < aktivniPismeno.id ? "_" : "&nbsp;") }} </div>
                        </div>
                    </div>
                </div>
            </div>
            <span class="unfocused"
                :style="{ display: unfocused ? 'block' : 'none', top: (route.fullPath == '/prvni-psani' || route.fullPath.split('/')[1] == 'prace') ? '200px' : '235px' }"
                @click="input?.focus()">Klikni sem nebo zmáčkni <span class="klavesa-v-textu">Mezerník</span> !</span>

            <input type="text" autocomplete="off" ref="input" id="input" @input="klik">

            <Transition>
                <div v-show="klavesnice != ''">
                    <Klavesnice :typ="klavesnice" :aktivniPismeno="aktivniPismeno.znak" :rozmazat="hideKlavesnice || prestalPsat"
                        :cekame="(aktivniPismeno.id == 0 || aktivniPismeno.id == -1) && cass == 0" :full-hide="fullHideKlavesnice" />
                    <Tooltip v-if="props.resetBtn" zprava="Restart cvičení <span class='klavesa-v-textu-mensi'>Delete</span>" :sirka="120"
                        :vzdalenost="6" :xOffset="385" :yOffset="-154">
                        <div id="reset-btn" @click="resetTlacitko(); animace(); input?.focus();"
                            :class="{ schovat: route.fullPath == '/prvni-psani' }">
                            <img :style="{ transform: rotace }" src="../assets/icony/reset.svg" alt="Restart">
                        </div>
                    </Tooltip>
                    <div v-if="props.resetBtn" id="hide-btn" @click="fullHideKlavesnice = !fullHideKlavesnice; input?.focus()"
                        :class="{ schovat: route.fullPath == '/prvni-psani' }"
                        :style="{ top: route.fullPath.split('/')[1] == 'lekce' ? '-140px' : '-70px' }">
                        <img v-if="!fullHideKlavesnice" src="../assets/icony/oko.svg" alt="Schovat" width="34">
                        <img v-else src="../assets/icony/okoSkrtnuty.svg" alt="Schovat" width="34">
                    </div>

                    <div id="zvuk-btn" @click="toggleZvuk(); input?.focus()">
                        <img v-if="zvukyZaply" style="margin-top: 1px;" class="zvuk-icon" src="../assets/icony/zvukOn.svg" alt="Zvuky jsou zapnuté">
                        <img v-else style="margin-left: 1px;" class="zvuk-icon" src="../assets/icony/zvukOff.svg" alt="Zvuky jsou vypnuté">
                    </div>
                </div>
            </Transition>
            <Transition>
                <div id="nepise" v-if="prestalPsat"
                    :style="{ boxShadow: fullHideKlavesnice ? 'none' : '0px 0px 10px 2px rgba(0, 0, 0, 0.75)', top: (route.fullPath == '/prvni-psani' || route.fullPath.split('/')[1] == 'prace') ? '370px' : '403px' }">
                    <h3>Jsi tam ještě?</h3>
                    <p>
                        Přestal jsi psát a tak jsme museli cvičení přerušit.
                    </p>
                    <button class="tlacitko" @click="prestalPsat = false; input?.focus()">Jsem tu!</button>
                </div>
            </Transition>
        </div>
</template>

<style scoped>
.klavesa-v-textu {
    padding: 3px 5px 2px 5px;
    font-size: 20px;
}

span.unfocused {
    position: absolute;
    font-size: 1.6em;
    font-weight: 600;
    user-select: none;
}

#input {
    position: absolute;
    opacity: 0;
    cursor: default;
}

#hide-btn {
    top: 100px;
    position: relative;
    width: 55px;
    height: 55px;
    background-color: var(--tmave-fialova);
    border-radius: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
    left: 695px;
    cursor: pointer;
    transition: background-color 0.1s;
    user-select: none;
}

#nepise {
    background-color: var(--tmave-fialova);
    padding: 20px;
    border-radius: 10px;
    position: absolute;
    display: flex;
    gap: 6px;
    align-items: center;
    flex-direction: column;
}

#nepise h3 {
    font-size: 26px;
    font-weight: 500;
    margin-bottom: 12px;
}

#nepise p {
    font-size: 20px;
}

#nepise .tlacitko {
    font-size: 20px;
}

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
    cursor: pointer;
    transition: background-color 0.1s;
    user-select: none;
}

#reset-btn img {
    width: 30px;
    height: 30px;
    transition-duration: 0.2s;
    transition-timing-function: ease-out;
}

#reset-btn:hover,
#zvuk-btn:hover,
#hide-btn:hover {
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
    transition: filter 0.15s;
}

#nabidka.unfocused {
    filter: blur(2px);
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
    transition: filter 0.15s;
}

#ramecek.unfocused {
    filter: blur(5px) brightness(0.8) opacity(0.8);
}

#text {
    display: flex;
    flex-wrap: wrap;
    position: relative;
    top: 0;
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
    font-size: 25px;
    line-height: 35px;
    width: 17px;
    text-decoration: none;
    padding: 0 1px;
    margin-right: 1px;
    border-bottom: 3px solid rgba(255, 255, 255, 0);
    /* aby se nedojebala vyska lajny když jdu na dalsi radek*/
    color: var(--bila);
}

#nabidka h2 {
    font-size: 19px;
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