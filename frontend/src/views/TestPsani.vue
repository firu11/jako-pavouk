<script setup lang="ts">
import { checkTeapot, getToken, pridatOznameni } from '../utils';
import { computed, onMounted, ref, toRaw } from 'vue';
import axios from 'axios';
import Vysledek from '../components/Vysledek.vue';
import { useHead } from '@unhead/vue';
import Psani from '../components/Psani.vue';
import { nastaveniJmeno, prihlasen } from '../stores';
import { useRouter } from 'vue-router';

useHead({
    title: "Test psaní",
    meta: [
        {
            name: "description",
            content: "Test psaní všemi deseti. Zjisti jak rychle píšeš a jak jsi přesný!",
        }
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/test-psani"
        }
    ]
})

const router = useRouter()

const text = ref([[]] as { id: number, znak: string, spatne: number, }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const cas = ref(0)
const nejcastejsiChyby = ref()

const typ = ref(true) // false = slova, true = vety
const delka = ref(1)

const klavesnice = ref("")
const diakritika = ref(true)
const velkaPismena = ref(false)

const psaniRef = ref()

const konec = ref(false)

const hideKlavecnice = ref(false)

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

function get() {
    axios.post("/test-psani",
        {
            typ: typ.value ? "vety" : "slova",
            delka: delka.value,
        },
        {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }
    ).then(response => {
        let text2 = [[]] as { id: number, znak: string, spatne: number, }[][]
        response.data.text.forEach((slovo: string, i: number) => {
            text2.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text2[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0 })
                delkaTextu.value++
            })
        })
        text.value = text2
        loadAlternativy()
        toggleDiakritikaAVelkaPismena()
        klavesnice.value = response.data.klavesnice
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni()
        }
    })
}

onMounted(() => {
    const mobil = document.body.clientWidth <= 1000
    if (mobil) {
        router.back()
        pridatOznameni('Psaní na telefonech zatím neučíme...')
        return
    }
    let nastaveni = localStorage.getItem(nastaveniJmeno)
    if (nastaveni !== null) {
        let obj = JSON.parse(nastaveni)
        diakritika.value = obj.diakritika
        velkaPismena.value = obj.velkaPismena
        typ.value = obj.typ
        if (!typ.value) {
            delka.value = 10
        }
    }
    get()
})

function restart() {
    delkaTextu.value = 0

    get()
    konec.value = false
}

function konecTextu(c: number, o: number, p: number, n: any[]) {
    cas.value = c
    opravenePocet.value = o
    preklepy.value = p
    nejcastejsiChyby.value = n
    konec.value = true
}

function d(x: number) {
    delka.value = x
    restart()
    psaniRef.value.restart()
}

function disabledBtn(e: KeyboardEvent) {
    e.preventDefault()
}

const rotaceStupne = ref(0)
function animace() {
    rotaceStupne.value += 45
}

const rotace = computed(() => {
    return `rotate(${rotaceStupne.value}deg)`
})

const klavModel = ref(false)
function switchKlavesnice() {
    if (klavesnice.value == "qwertz") klavesnice.value = "qwerty"
    else klavesnice.value = "qwertz"
}

let puvodniText = [[]] as { id: number, znak: string, spatne: number }[][]
let textBezDiakritiky = [[]] as { id: number, znak: string, spatne: number }[][]
let textMalym = [[]] as { id: number, znak: string, spatne: number }[][]
let textOboje = [[]] as { id: number, znak: string, spatne: number }[][]

function toggleDiakritikaAVelkaPismena() {
    if (!diakritika.value && !velkaPismena.value) {
        text.value = structuredClone(textOboje)
    } else if (!diakritika.value) {
        text.value = structuredClone(textBezDiakritiky)
    } else if (!velkaPismena.value) {
        text.value = structuredClone(textMalym)
    } else {
        text.value = structuredClone(puvodniText)
    }
    psaniRef.value.restart()
    localStorage.setItem(nastaveniJmeno, JSON.stringify({ "diakritika": diakritika.value, "velkaPismena": velkaPismena.value, "typ": typ.value }))
}

async function loadAlternativy() {
    puvodniText = structuredClone(toRaw(text.value))
    textBezDiakritiky = structuredClone(toRaw(text.value))
    textBezDiakritiky.forEach(slovo => {
        slovo.forEach(pismeno => {
            pismeno.znak = pismeno.znak.normalize("NFD").replace(/[\u0300-\u036f]/g, "")
        })
    })
    textMalym = structuredClone(toRaw(text.value))
    textMalym.forEach(slovo => {
        slovo.forEach(pismeno => {
            pismeno.znak = pismeno.znak.toLocaleLowerCase()
        })
    })
    textOboje = structuredClone(toRaw(text.value))
    textOboje.forEach(slovo => {
        slovo.forEach(pismeno => {
            pismeno.znak = pismeno.znak.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLocaleLowerCase()
        })
    })
}
</script>

<template>
    <h1 style="margin: 0">Test psaní</h1>

    <Psani v-if="!konec" @konec="konecTextu" @pise="hideKlavecnice = false" :text="text" :delkaTextu="delkaTextu"
        :klavesnice="klavesnice" :hide-klavesnice="hideKlavecnice" ref="psaniRef" />

    <Vysledek v-else @restart="restart" :preklepy="preklepy" :opravenych="opravenePocet" :delkaTextu="delkaTextu"
        :casF="casFormat" :cas="cas" :cislo="'test-psani'" :posledni="true" :nejcastejsiChyby="nejcastejsiChyby" />

    <Transition>
        <div v-if="!konec && hideKlavecnice" id="psani-menu">

            <div class="kontejner">
                <input v-model="typ" v-on:change="d(typ ? 1 : 10)" type="checkbox" id="toggle" class="toggleCheckbox" />
                <label for="toggle" class="toggleContainer">
                    <div>Slova</div>
                    <div>Věty</div>
                </label>

                <div v-if="typ" id="delka" :class="{ odsunout: prihlasen }">
                    <button @keyup="disabledBtn" :class="{ aktivni: 1 == delka }" @click="d(1)">1</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 3 == delka }" @click="d(3)">3</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 5 == delka }" @click="d(5)">5</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 10 == delka }" @click="d(10)">10</button>
                </div>
                <div v-else id="delka" :class="{ odsunout: prihlasen }">
                    <button @keyup="disabledBtn" :class="{ aktivni: 10 == delka }" @click="d(10)">10</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 25 == delka }" @click="d(25)">25</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 50 == delka }" @click="d(50)">50</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 100 == delka }" @click="d(100)">100</button>
                </div>

                <input v-if="!prihlasen" @change="switchKlavesnice" v-model="klavModel" type="checkbox" id="toggle1"
                    class="toggleCheckbox" />
                <label v-if="!prihlasen" for="toggle1" class="toggleContainer">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </div>

            <hr id="predel">

            <div class="kontejner">
                <label for="toggle2" class="kontejner">
                    <input v-model="velkaPismena" @change="toggleDiakritikaAVelkaPismena" type="checkbox"
                        id="toggle2" class="radio" />
                    Velká písmena
                </label>

                <label for="toggle3" class="kontejner">
                    <input v-model="diakritika" @change="toggleDiakritikaAVelkaPismena" type="checkbox" id="toggle3"
                        class="radio" />
                    Diakritika
                </label>
            </div>
        </div>
    </Transition>

    <div v-if="!konec && klavesnice != ''" id="nastaveniBtn" @click="hideKlavecnice = !hideKlavecnice; animace()">
        <img :style="{ transform: rotace }" src="../assets/icony/nastaveni.svg" alt="Nastavení">
    </div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
    transition: opacity 0.2s;
}

.v-enter-from,
.v-leave-to {
    opacity: 0;
}

#predel {
    margin: 12px 0 15px 0;
    width: 92%;
    border: 1px solid var(--fialova);
}

.kontejner {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
    margin: 0 10px;
    cursor: pointer;
    transition: filter 0.2s;
}

label.kontejner:hover {
    filter: brightness(120%);
}

.odsunout {
    margin-left: 18px;
}

#nastaveniBtn {
    position: relative;
    width: 55px;
    height: 55px;
    background-color: var(--tmave-fialova);
    border-radius: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
    left: 385px;
    bottom: 236px;
    cursor: pointer;
    transition: background-color 0.1s;
}

#nastaveniBtn img {
    width: 30px;
    transition-duration: 0.4s;
}

#nastaveniBtn:hover {
    background-color: var(--fialova);
}

.radio {
    appearance: none;
    -webkit-appearance: none;
    border: 0.15rem solid var(--fialova);
    border-radius: 10rem;
    transition: filter 0.1s;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.radio::before {
    content: "";
    width: 14px;
    height: 14px;
    transform: scale(0);
    background-color: var(--fialova);
    border-radius: 10rem;
    transition: 0.1s;
    display: block; 
}

.radio:checked:before {
    transform: scale(1);
}

.aktivni {
    color: var(--svetle-fialova) !important;
}

#delka {
    display: flex;
    gap: 6px;
    justify-content: center;
    width: 120px;
}

#psani-menu {
    background-color: var(--tmave-fialova);
    padding: 10px;
    border-radius: 8px;
    min-height: 50px;
    margin-bottom: 186px;
    margin-top: 40px;
    display: flex;
    flex-direction: column;
    gap: 0 10px;
    position: absolute;
    top: 400px;
    max-width: 420px;
    flex-wrap: wrap;
    align-items: center;
}

#psani-menu button {
    background-color: transparent;
    border: none;
    color: var(--bila);
    transition: 0.1s;
    font-size: 1em;
    border-radius: 5px;
    padding: 0 2px;
}

#psani-menu button:hover {
    color: var(--svetle-fialova);
    font-weight: 500;
    color: white;
    cursor: pointer;
}

.toggleCheckbox {
    display: none;
}

.toggleContainer {
    position: relative;
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    width: fit-content;
    font-weight: bold;
    color: var(--tmave-fialova);
    cursor: pointer;
    background: transparent;
    font-size: 1em;
    border-radius: 8px;
    border: 1px var(--fialova) solid;
    justify-self: start;
    height: 30px;
}

.toggleContainer::before {
    content: '';
    position: absolute;
    width: 50%;
    height: 100%;
    left: 0%;
    border-radius: 6px;
    background: var(--fialova);
    transition: all 0.3s;
}

.toggleCheckbox:checked+.toggleContainer::before {
    left: 50%;
}

.toggleContainer div {
    padding: 6px;
    text-align: center;
    z-index: 1;
    user-select: none;
    position: relative;
    top: -1px;
}

.toggleCheckbox:checked+.toggleContainer div:first-child {
    color: transparent;
    transition: color 0.3s;
}

.toggleCheckbox:checked+.toggleContainer div:last-child {
    color: white;
    transition: color 0.3s;
}

.toggleCheckbox+.toggleContainer div:first-child {
    color: white;
    transition: color 0.3s;
}

.toggleCheckbox+.toggleContainer div:last-child {
    color: transparent;
    transition: color 0.3s;
}
</style>