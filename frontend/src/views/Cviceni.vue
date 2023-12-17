<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { formatovanyPismena, getToken, pridatOznameni } from '../utils';
import SipkaZpet from '../components/SipkaZpet.vue';
import { computed, onMounted, ref } from 'vue';
import axios from 'axios';
import { onUnmounted } from 'vue';
import Vysledek from '../components/Vysledek.vue';
import Klavesnice from '../components/Klavesnice.vue';
import { useSound } from '@vueuse/sound';
import { useHead } from '@unhead/vue';

const router = useRouter()
const route = useRoute()
const pismena: string = Array.isArray(route.params.pismena) ? route.params.pismena[0] : route.params.pismena
const cislo: string = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

useHead({
    title: "Cvičení " + pismena
})

const text = ref([[]] as { id: number, znak: string, spatne: number, }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const counter = ref(0)
const counterSlov = ref(0)
const preklepy = ref(0)
const timerZacatek = ref(0)
const cas = ref(0)
const textElem = ref<HTMLInputElement>()
let indexPosunuti = -1

let predchoziZnak = ""

const posledni = ref(false)
const klavesnice = ref("")

const zvukyZaply = ref(true)
let tmp = localStorage.getItem("pavouk_zvuk")
if (tmp == null) {
    zvukyZaply.value = true
} else {
    zvukyZaply.value = JSON.parse(tmp) === true // nejde to dat na jednu lajnu TS sus
}
const zvuky = [useSound('/zvuky/klik1.ogg'), useSound('/zvuky/klik2.ogg'), useSound('/zvuky/klik3.ogg'), useSound('/zvuky/miss.ogg')]

const capslock = ref(false)
let interval: number
const konec = ref(false)

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

const progress = computed(() => {
    return delkaTextu.value !== 0 ? ((aktivniPismeno.value.id) / delkaTextu.value) * 100 : 0
})

const aktivniPismeno = computed(() => {
    if (counterSlov.value < text.value.length - 1) return text.value[counterSlov.value][counter.value]
    return { id: -1, znak: "", spatne: 0 }
})

function get() {
    axios.get("/cvic/" + encodeURIComponent(pismena) + "/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0 })
                delkaTextu.value++
            })
        })
        posledni.value = response.data.posledni
        klavesnice.value = response.data.klavesnice
    }).catch(_ => {
        pridatOznameni()
        router.back()
    });
}

onMounted(() => {
    get()
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
    if (text.value[counterSlov.value].length - 1 === counter.value) { // posledni pismeno ve slovu
        if (aktivniPismeno.value.spatne === 1) {
            preklepy.value++
        }
        counterSlov.value++
        counter.value = 0
    } else {
        if (aktivniPismeno.value.spatne === 1) preklepy.value++
        counter.value++
    }
}

function backPismeno() {
    if (counter.value === 0) { // prvni pismeno ve slovu
        counterSlov.value--
        counter.value = text.value[counterSlov.value].length - 1
        if (aktivniPismeno.value.spatne === 1) preklepy.value--
    } else {
        counter.value--
        if (aktivniPismeno.value.spatne === 1) preklepy.value--
    }
}

function jeSHackem(key: string) {
    let velkym = aktivniPismeno.value.znak.toLocaleUpperCase() === aktivniPismeno.value.znak
    if (predchoziZnak === "ˇ") {
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ď" && (!velkym && key === "d" || velkym && key === "D")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ň" && (!velkym && key === "n" || velkym && key === "N")) return true
        if (aktivniPismeno.value.znak.toLocaleLowerCase() === "ť" && (!velkym && key === "t" || velkym && key === "T")) return true
    } else if (predchoziZnak === "´") {
        if (aktivniPismeno.value.znak === "ó" && key === "o") return true
    }
}

function klik(this: any, e: KeyboardEvent) {
    e.preventDefault() // ať to nescrolluje a nehazí nějaký stupid zkratky
    startTimer()

    if (delkaTextu.value == 0) {
        console.log(e.key)
        return
    }

    if (e.key === aktivniPismeno.value.znak || jeSHackem(e.key)) {
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

    let aktualniY = document.getElementById("p" + aktivniPismeno.value.id)?.getBoundingClientRect().y!
    let lastY = document.getElementById("p" + (aktivniPismeno.value.id - 1))?.getBoundingClientRect().y!
    if (aktualniY - lastY > 30) {
        indexPosunuti++
        if (indexPosunuti > 0) textElem.value!.style.top = `${indexPosunuti * (-2.2 - 0.188)}rem` // posunuti dolu
    }

    if (aktivniPismeno.value.id === -1) { // konec
        clearInterval(interval)
        calcCas() // naposledy
        konec.value = true
        document.removeEventListener("keypress", klik)
        document.removeEventListener("keydown", specialniKlik)
    }
}

function specialniKlik(e: KeyboardEvent) {
    capslockCheck(e)
    if (e.key === "Backspace") {
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

function restart() {
    timerZacatek.value = 0
    cas.value = 0
    counter.value = 0
    counterSlov.value = 0
    preklepy.value = 0
    konec.value = false
    text.value = [[]] as { id: number, znak: string, spatne: number, }[][]
    delkaTextu.value = 0
    indexPosunuti = -1

    get()

    document.addEventListener("keypress", klik)
    document.addEventListener("keydown", specialniKlik)
}

function toggleZvuk() {
    zvukyZaply.value = !zvukyZaply.value
    localStorage.setItem("pavouk_zvuk", zvukyZaply.value.toString())
}

function format(p: string) {
    if (p === "zbylá diakritika") return "Zbylá diakritika"
    else if (p === "velká písmena (shift)") return "Velká písmena (Shift)"
    return formatovanyPismena(p)
}
</script>

<template>
    <h1 class="nadpisSeSipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Lekce: {{ format(pismena) }}
    </h1>
    <h2>Cviceni: {{ cislo }}</h2>

    <div id="flex" v-if="!konec">
        <div id="nabidka">
            <h3 id="cas">{{ casFormat }}s</h3>
            <h3 :style="{ visibility: capslock ? 'visible' : 'hidden' }" id="capslock">CapsLock</h3>
            <h3 id="preklepy">Překlepy: {{ preklepy }}</h3>
        </div>

        <div id="ramecek">
            <div id="fade">
                <div id="text" ref="textElem">
                    <div class="slovo" v-for="s in text">
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

            <Klavesnice v-if="klavesnice != ''" :typ="klavesnice" :aktivniPismeno="aktivniPismeno.znak"></Klavesnice>

            <div id="zvukBtn" @click="toggleZvuk">
                <img v-if="zvukyZaply" style="margin-top: 1px;" class="zvukIcon" src="../assets/icony/zvukOn.svg"
                    alt="Zvuky jsou zapnuté">
                <img v-else style="margin-left: 1px;" class="zvukIcon" src="../assets/icony/zvukOff.svg"
                    alt="Zvuky jsou vypnuté">
            </div>
        </div>

        <Vysledek v-else @restart="restart" :preklepy="preklepy" :delkaTextu="delkaTextu" :casF="casFormat" :cas="cas"
            :pismena="pismena" :cislo="cislo" :posledni="posledni"></Vysledek>
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
    bottom: 20px;
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