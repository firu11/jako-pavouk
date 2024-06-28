<script setup lang="ts">
import { checkTeapot, clone, getToken, MojeMapa, pridatOznameni } from "../utils";
import { computed, onMounted, ref, toRaw } from "vue";
import axios from "axios";
import Vysledek from "../components/Vysledek.vue";
import { useHead } from "@unhead/vue";
import Psani from "../components/Psani.vue";
import { nastaveniJmeno, prihlasen, mobil } from "../stores";
import { useRouter } from "vue-router";

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

const text = ref([] as { id: number, znak: string, spatne: number, }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const cas = ref(0)
const nejcastejsiChyby = ref()

const typ = ref(true) // false = slova, true = vety
const delka = ref(60)

const klavesnice = ref("")
const diakritika = ref(true)
const velkaPismena = ref(false)

const psaniRef = ref()

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)
const nacitamNovej = ref(false)

const hideKlavecnice = ref(false)

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

function get() {
    nacitamNovej.value = true
    axios.post("/test-psani",
        {
            typ: typ.value ? "vety" : "slova",
            cas: 1,
        },
        {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }
    ).then(response => {
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0 })
                delkaTextu.value++
            })
        })

        loadAlternativy()
        toggleDiakritikaAVelkaPismena()
        klavesnice.value = response.data.klavesnice
        nacitamNovej.value = false
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni()
        }
        nacitamNovej.value = false
    })
}

onMounted(() => {
    if (mobil.value) {
        router.back()
        pridatOznameni("Psaní na telefonech zatím neučíme...")
        return
    }
    let nastaveni = localStorage.getItem(nastaveniJmeno)
    if (nastaveni !== null) {
        let obj = JSON.parse(nastaveni)
        diakritika.value = obj.diakritika
        velkaPismena.value = obj.velkaPismena
        typ.value = obj.typ
        delka.value = obj.delka
    }
    get()
})

function restart() {
    delkaTextu.value = 0
    text.value = [] as { id: number, znak: string, spatne: number, }[][]

    get()
    konec.value = false
}

function konecTextu(c: number, o: number, p: number, n: MojeMapa, d: number) {
    cas.value = Math.round(c * 100) / 100
    opravenePocet.value = o
    preklepy.value = p
    nejcastejsiChyby.value = new MojeMapa(n)
    konec.value = true
    delkaNapsanehoTextu.value = d
}

function d(x: number) {
    delka.value = x
    restart()
    psaniRef.value?.restart()
}

function disabledBtn(e: KeyboardEvent) {
    e.preventDefault()
    saveNastaveni()
}

const rotacePocet = ref(0)
function animace() {
    rotacePocet.value += 1
}

const rotace = computed(() => {
    return `rotate(${rotacePocet.value * 45}deg)`
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
        text.value = clone(textOboje)
    } else if (!diakritika.value) {
        text.value = clone(textBezDiakritiky)
    } else if (!velkaPismena.value) {
        text.value = clone(textMalym)
    } else {
        text.value = clone(puvodniText)
    }
    psaniRef.value?.restart()
    saveNastaveni()
}

async function loadAlternativy() {
    puvodniText = clone(toRaw(text.value))
    textBezDiakritiky = clone(toRaw(text.value))
    textBezDiakritiky.forEach(slovo => {
        slovo.forEach(pismeno => {
            pismeno.znak = pismeno.znak.normalize("NFD").replace(/[\u0300-\u036f]/g, "")
        })
    })
    textMalym = clone(toRaw(text.value))
    textMalym.forEach(slovo => {
        slovo.forEach(pismeno => {
            pismeno.znak = pismeno.znak.toLocaleLowerCase()
        })
    })
    textOboje = clone(toRaw(text.value))
    textOboje.forEach(slovo => {
        slovo.forEach(pismeno => {
            pismeno.znak = pismeno.znak.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLocaleLowerCase()
        })
    })
}

async function prodlouzit() {
    axios.post("/test-psani",
        {
            typ: typ.value ? "vety" : "slova",
        },
        {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }
    ).then(response => {
        if (!diakritika.value && !velkaPismena.value) {
            for (let i = 0; i < response.data.text.length; i++) {
                response.data.text[i] = response.data.text[i].normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLocaleLowerCase()
            }
        } else if (!diakritika.value && velkaPismena.value) {
            for (let i = 0; i < response.data.text.length; i++) {
                response.data.text[i] = response.data.text[i].normalize("NFD").replace(/[\u0300-\u036f]/g, "")
            }
        } else if (diakritika.value && !velkaPismena.value) {
            for (let i = 0; i < response.data.text.length; i++) {
                response.data.text[i] = response.data.text[i].toLocaleLowerCase()
            }
        }

        let pocetSlov = text.value.length
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[pocetSlov - 1 + i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0 })
                delkaTextu.value++
            })
        })

    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni()
        }
    })
}

function saveNastaveni() {
    localStorage.setItem(nastaveniJmeno, JSON.stringify({ "diakritika": diakritika.value, "velkaPismena": velkaPismena.value, "typ": typ.value, "delka": delka.value }))
}
</script>

<template>
    <h1 style="margin: 0">Test psaní</h1>

    <Psani v-if="!konec" @konec="konecTextu" @restart="restart" @pise="hideKlavecnice = false" @prodlouzit="prodlouzit" :text="text"
        :klavesnice="klavesnice" :delkaTextu="delkaTextu" :hide-klavesnice="hideKlavecnice" :nacitam-novej="nacitamNovej" :cas="delka"
        ref="psaniRef" />

    <Vysledek v-else @restart="restart" :preklepy="preklepy" :opravenych="opravenePocet" :delkaTextu="delkaNapsanehoTextu" :casF="casFormat"
        :cas="cas" :cislo="'test-psani'" :posledni="true" :nejcastejsiChyby="nejcastejsiChyby" />

    <Transition>
        <div v-if="!konec && hideKlavecnice" id="psani-menu">

            <div class="kontejner" style="gap: 20px;">
                <input v-model="typ" type="checkbox" id="toggle" class="toggleCheckbox" @change="saveNastaveni"/>
                <label for="toggle" class="toggleContainer">
                    <div>Slova</div>
                    <div>Věty</div>
                </label>
                <input v-if="!prihlasen" @change="switchKlavesnice" v-model="klavModel" type="checkbox" id="toggle1" class="toggleCheckbox" />
                <label v-if="!prihlasen" for="toggle1" class="toggleContainer">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </div>

            <div class="kontejner">
                <div id="delka" :class="{ odsunout: prihlasen }">
                    <button @keyup="disabledBtn" :class="{ aktivni: 15 == delka }" @click="d(15)">15s</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 30 == delka }" @click="d(30)">30s</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 60 == delka }" @click="d(60)">1min</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 120 == delka }" @click="d(120)">2min</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 300 == delka }" @click="d(300)">5min</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 600 == delka }" @click="d(600)">10min</button>
                </div>
            </div>

            <hr id="predel">

            <div class="kontejner">
                <label for="toggle2" class="kontejner">
                    <input v-model="velkaPismena" @change="toggleDiakritikaAVelkaPismena" type="checkbox" id="toggle2" class="radio" />
                    Velká písmena
                </label>

                <label for="toggle3" class="kontejner">
                    <input v-model="diakritika" @change="toggleDiakritikaAVelkaPismena" type="checkbox" id="toggle3" class="radio" />
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
    margin: 4px 0;
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
    bottom: 220px;
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

#delka {
    display: flex;
    gap: 6px;
    justify-content: center;
}
</style>