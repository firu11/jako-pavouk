<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { onUnmounted } from 'vue';
import Klavesnice from '../components/Klavesnice.vue';
import { useSound } from '@vueuse/sound';

const emit = defineEmits(["konec", "pise"])

const props = defineProps<{
    text: { id: number, znak: string, spatne: number }[][]
    delkaTextu: number,
    klavesnice: string,
    hideKlavesnice: boolean
}>()

const counter = ref(0)
const counterSlov = ref(0)
const preklepy = ref(0)
const opravene = new Map<String, Boolean>()
const timerZacatek = ref(0)
const cas = ref(0)
const textElem = ref<HTMLInputElement>()
let indexPosunuti = -1

let predchoziZnak = ""

const zvukyZaply = ref(true)
let tmp = localStorage.getItem("pavouk_zvuk")
if (tmp == null) {
    zvukyZaply.value = true
} else {
    zvukyZaply.value = JSON.parse(tmp) === true // nejde to dat na jednu lajnu TS sus
}
const zvuky = [useSound(new URL('../assets/zvuky/klik1.ogg', import.meta.url).href), useSound(new URL('../assets/zvuky/klik2.ogg', import.meta.url).href), useSound(new URL('../assets/zvuky/klik3.ogg', import.meta.url).href), useSound(new URL('../assets/zvuky/miss.ogg', import.meta.url).href)]

const capslock = ref(false)
let interval: number

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

const progress = computed(() => {
    return props.delkaTextu !== 0 ? ((aktivniPismeno.value.id) / props.delkaTextu) * 100 : 0
})

const aktivniPismeno = computed(() => {
    if (counterSlov.value < props.text.length - 1) return props.text[counterSlov.value][counter.value]
    return { id: -1, znak: "", spatne: 0 }
})

onMounted(() => {
    document.addEventListener("keypress", klik) // je depracated ale je O TOLIK LEPSI ZE HO BUDU POUZIVAT PROSTE https://stackoverflow.com/questions/52882144/replacement-for-deprecated-keypress-dom-event
    document.addEventListener("keydown", specialniKlik)
})

onUnmounted(() => {
    document.removeEventListener("keypress", klik)
    document.removeEventListener("keydown", specialniKlik)
})

function capslockCheck(e: KeyboardEvent) { // TODO chtelo by to checknout hned po nacteni stranky ale nevim jestli to jde (spíš ne)
    capslock.value = e.getModifierState("CapsLock")
}

function nextPismeno() {
    if (props.text[counterSlov.value].length - 1 === counter.value) { // posledni pismeno ve slovu
        if (aktivniPismeno.value.spatne === 1) {
            preklepy.value++
        }
        counterSlov.value++
        counter.value = 0
    } else {
        if (aktivniPismeno.value.spatne === 1) preklepy.value++
        counter.value++
    }
    emit("pise")
}

function backPismeno() {
    if (counter.value === 0) { // prvni pismeno ve slovu
        counterSlov.value--
        counter.value = props.text[counterSlov.value].length - 1
        if (aktivniPismeno.value.spatne === 1) {
            preklepy.value--
            opravene.set(`${counterSlov.value}${counter.value}`, true)
        }
    } else {
        counter.value--
        if (aktivniPismeno.value.spatne === 1) {
            preklepy.value--
            opravene.set(`${counterSlov.value}${counter.value}`, true)
        }
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

function klik(this: any, e: KeyboardEvent) {
    e.preventDefault() // ať to nescrolluje a nehazí nějaký stupid zkratky
    startTimer()

    if (props.delkaTextu == 0) {
        console.log(props.text)
        return
    }

    let hacek = jeSHackem(e.key)
    if (hacek) predchoziZnak = ""

    if (e.key === aktivniPismeno.value.znak || hacek) {
        if (zvukyZaply.value) zvuky[Math.floor(Math.random() * 2)].play()
        if (aktivniPismeno.value.spatne === 1) {
            aktivniPismeno.value.spatne = 2
        }
        nextPismeno()
    } else {
        if (zvukyZaply.value) zvuky[3].play()
        aktivniPismeno.value.spatne = 1
        nextPismeno()
    }

    posunoutRadek()

    if (aktivniPismeno.value.id === -1) { // konec
        clearInterval(interval)
        calcCas() // naposledy
        document.removeEventListener("keypress", klik)
        document.removeEventListener("keydown", specialniKlik)
        emit("konec", cas.value, opravene.size, preklepy.value)
        restart()
    }

    if (predchoziZnak != "") predchoziZnak = ""
}

function posunoutRadek() {
    let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y!
    let lastY = document.getElementById("p" + (aktivniPismeno.value.id - 1))?.getBoundingClientRect().y!
    if (aktualniY - lastY > 30) {
        indexPosunuti++
        if (indexPosunuti > 0) textElem.value!.style.top = `${indexPosunuti * (-2.2 - 0.188)}rem` // posunuti dolu
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
        if (aktivniPismeno.value.id !== 0) {
            if (e.ctrlKey) { // tak dáme celé slovo pryč (Ctrl + Backspace zkratka)
                let lastY = document.getElementById("p" + (aktivniPismeno.value.id))?.getBoundingClientRect().y!
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
                let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y!
                if (lastY - aktualniY > 30) {
                    indexPosunuti--
                    if (indexPosunuti > -1) textElem.value!.style.top = `${indexPosunuti * (-2.2 - 0.188)}rem`
                }
            }
            else {
                backPismeno()
                let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y!
                let lastY = document.getElementById("p" + (aktivniPismeno.value.id + 1))?.getBoundingClientRect().y!
                if (lastY - aktualniY > 30) {
                    indexPosunuti--
                    if (indexPosunuti > -1) textElem.value!.style.top = `${indexPosunuti * (-2.2 - 0.188)}rem`
                }
            }
            if (zvukyZaply.value) zvuky[Math.floor(Math.random() * 2)].play()
        }
    }
}

function startTimer() {
    if (timerZacatek.value === 0) {
        timerZacatek.value = Date.now()
        interval = setInterval(calcCas, 100)
    }
}

function calcCas() {
    cas.value = (Date.now() - timerZacatek.value) / 1000
}

function toggleZvuk() {
    zvukyZaply.value = !zvukyZaply.value
    localStorage.setItem("pavouk_zvuk", zvukyZaply.value.toString())
}

function restart() {
    clearInterval(interval)
    timerZacatek.value = 0
    cas.value = 0
    counter.value = 0
    counterSlov.value = 0
    preklepy.value = 0
    indexPosunuti = -1
    textElem.value!.style.top = "0rem" // reset posunuti
}

defineExpose({ restart })
</script>

<template>
    <div id="flex">
        <div id="nabidka">
            <h3 id="cas">{{ casFormat }}s</h3>
            <h3 :style="{ visibility: capslock ? 'visible' : 'hidden' }" id="capslock">CapsLock</h3>
            <h3 id="preklepy">Překlepy: {{ preklepy }}</h3>
        </div>

        <div id="ramecek">
            <div id="fade">
                <div id="text" ref="textElem">
                    <div class="slovo" v-for="s in text.slice(0, 60)">
                        <div v-for="p in s" class="pismeno" :id="'p' + p.id"
                            :class="{ podtrzeni: p.id === aktivniPismeno.id, spatnePismeno: p.spatne === 1 && aktivniPismeno.id > p.id, opravenePismeno: p.spatne === 2 && aktivniPismeno.id > p.id, spravnePismeno: !p.spatne && aktivniPismeno.id > p.id }">
                            {{ (p.znak !== " " ? p.znak : p.spatne && p.id < aktivniPismeno.id ? "_" : "&nbsp") }} </div>
                        </div>
                    </div>
                </div>
            </div>
            <div id="bar">
                <div :style="'width: ' + progress + '%'" id="progress">&nbsp{{ Math.floor(progress) }}%&nbsp
                </div>
            </div>

            <Klavesnice v-if="klavesnice != '' && !hideKlavesnice" :typ="klavesnice" :aktivniPismeno="aktivniPismeno.znak"></Klavesnice>

            <div id="zvukBtn" @click="toggleZvuk">
                <img v-if="zvukyZaply" style="margin-top: 1px;" class="zvukIcon" src="../assets/icony/zvukOn.svg"
                    alt="Zvuky jsou zapnuté">
                <img v-else style="margin-left: 1px;" class="zvukIcon" src="../assets/icony/zvukOff.svg"
                    alt="Zvuky jsou vypnuté">
            </div>
        </div>
</template>

<style scoped>

.zvukIcon {
    width: 45px;
    height: 35px;
    margin-top: 1px;
}

#zvukBtn {
    position: absolute;
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
}

#zvukBtn:hover {
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
    font-weight: bold;
}

#ramecek {
    padding: 10px;
    height: 200px;
    border-radius: 10px 10px 0 0;
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
    overflow: hidden;
}

#text {
    display: flex;
    flex-wrap: wrap;
    position: relative;
    transition: ease 0.2s;
    top: 0em;
}

#fade {
    mask-image: linear-gradient(180deg, var(--tmave-fialova) 75%, transparent 97%);
    -webkit-mask-image: linear-gradient(180deg, var(--tmave-fialova) 75%, transparent 97%);
    height: 190px;
}

.slovo {
    display: flex;
    flex-wrap: nowrap;
}

.pismeno {
    border-radius: 3px;
    display: inline-flex;
    font-family: 'Red Hat Mono', monospace;
    font-weight: 400;
    font-size: 1.56rem;
    line-height: 2.2rem;
    text-decoration: none;
    padding: 0 1px;
    margin-right: 1px;
    border-bottom: 3px solid rgba(255, 255, 255, 0);
    /* aby se nedojebala vyska lajny když jdu na dalsi radek*/
    color: var(--bila);
    transition: 60ms;
}

#progress {
    height: 20px;
    background-color: var(--fialova);
    width: 0;
    border-bottom-left-radius: 10px;
    transition: ease 0.22s;
    text-align: right;
}

#bar {
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
    border-radius: 0 0 10px 10px;
    overflow: hidden;
}

.spravnePismeno {
    color: #9c9c9c;
}

.podtrzeni {
    border-bottom: 3px solid var(--bila);
    border-radius: 0;
    transition: 60ms;
}

.spatnePismeno {
    color: #ff0000;
}

.opravenePismeno {
    color: #b1529c;
}
</style>