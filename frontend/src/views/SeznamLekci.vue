<script setup lang="ts">
import axios from "axios";
import BlokLekce from "../components/BlokLekce.vue";
import Rada from "../components/Rada.vue"
import { onMounted, onUnmounted, ref } from "vue";
import { Oznacene, checkTeapot, getToken, pridatOznameni, napovedaKNavigaci } from "../utils";
import { useHead } from "unhead";
import { useRouter } from "vue-router";
import { mobil } from "../stores";

useHead({
    title: "Lekce",
    meta: [
        {
            name: "description",
            content: "Jako Pavouk nabízí spoustu lekcí s řadou typů cviční. Od jednotlivých písmen, přes slova až k celým větám. Naučte se psát všemi deseti!",
        },
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/lekce"
        }
    ]
})

const lekce = ref([[]] as { id: number, pismena: string, cislo: number }[][])
const dokoncene = ref([] as number[])
const o = new Oznacene()
const prvniNedokoncena = ref(1)

let dalsiCviceni: string

const router = useRouter()

const nacitam = ref(false)

onMounted(() => {
    nacitam.value = true

    const header = getToken() ? { headers: { Authorization: `Bearer ${getToken()}` } } : {}
    axios.get("/lekce", header)
        .then(response => {
            lekce.value = response.data.lekce
            dokoncene.value = response.data.dokoncene
            dalsiCviceni = response.data.dalsi_cviceni
            o.setMax(lekce.value.join(",").split(",").length) // pocet lekci

            let counter = 1
            let nebylaNedoko = true
            for (let i = 0; i < lekce.value.length; i++) {
                for (let j = 0; j < lekce.value[i].length; j++) {
                    lekce.value[i][j]["cislo"] = counter
                    counter += 1
                    if (dokoncene.value.includes(lekce.value[i][j]["id"]) && nebylaNedoko) prvniNedokoncena.value += 1
                    else nebylaNedoko = false
                }
            }
        }).catch(e => {
            if (!checkTeapot(e)) {
                pridatOznameni()
                console.log(e)
            }
        }).finally(() => {
            nacitam.value = false
        })


    document.addEventListener("keydown", e1)
    document.addEventListener("keyup", e2)
    document.addEventListener("mousemove", zrusitVyber)
})

onUnmounted(() => {
    document.removeEventListener("keydown", e1)
    document.removeEventListener("keyup", e2)
    document.removeEventListener("mousemove", zrusitVyber)
})

let jede = false
let ms = 120

function e1(e: KeyboardEvent) {
    if (e.key == "ArrowUp" || e.key == "ArrowLeft") {
        e.preventDefault()
        if (jede) return

        if (o.index.value == 0) o.index.value = prvniNedokoncena.value + 1
        o.mensi()
        let lekceE: HTMLElement | null = document.querySelector(`[i="${o.index.value}"]`)

        jede = true

        let scroll = document.body.getBoundingClientRect().top
        window.scrollTo({ top: lekceE!.offsetTop - 200, behavior: Math.abs(-scroll - lekceE!.offsetTop) > 600 ? "instant" : "smooth" })
        setTimeout(() => { jede = false }, ms)
    } else if (e.key == "ArrowDown" || e.key == "ArrowRight") {
        e.preventDefault()
        if (jede) return

        if (o.index.value == 0) o.index.value = prvniNedokoncena.value - 1
        o.vetsi()
        let lekceE: HTMLElement | null = document.querySelector(`[i="${o.index.value}"]`)

        jede = true

        let scroll = document.body.getBoundingClientRect().top
        window.scrollTo({ top: lekceE!.offsetTop - 200, behavior: Math.abs(-scroll - lekceE!.offsetTop) > 600 ? "instant" : "smooth" })
        setTimeout(() => { jede = false }, ms)
    } if (e.key == "Enter") {
        e.preventDefault()
        let lekceE: HTMLElement | null = document.querySelector(`.oznacene`)
        if (lekceE == null || o.bezOznaceni) {
            o.bezOznaceni = true
            o.index.value = prvniNedokoncena.value
            lekceE = document.querySelector(`[i="${o.index.value}"]`)
            window.scrollTo({ top: lekceE!.offsetTop - 200 })
        } else lekceE?.click()
    } else if (e.key == "Tab") {
        e.preventDefault()
        napovedaKNavigaci()
    }
}

function e2(e: KeyboardEvent) {
    if (e.key == "Enter") {
        e.preventDefault()
        let lekceE: HTMLElement | null = document.querySelector(`.oznacene`)
        router.push(lekceE!.getAttribute("href")!)  // chromium sus
    }
}

function zrusitVyber() {
    o.index.value = 0
}

function pokracovatOdPosledniho() {
    router.push("/lekce" + dalsiCviceni)
}

</script>

<template>
    <h1>Lekce</h1>
    <div id="seznam">
        <Rada :pocetDoko="nacitam ? -1 : dokoncene.length" />
        <button v-if="!nacitam && dalsiCviceni && !mobil" id="pokracovani" @click="pokracovatOdPosledniho">
            Pokračovat od posledního
            <img src="../assets/icony/start.svg" alt="Začít" width="35">
        </button>

        <h2>Střední řada</h2>
        <div v-if="lekce[0].length == 0" class="kategorie">
            <!-- jen aby tam něco bylo než se to načte -->
            <BlokLekce v-for="i in 4" pismena=". . ." :jeDokoncena="false" :cislo="i" :key="i" />
        </div>
        <div v-else class="kategorie">
            <BlokLekce v-for="l in lekce[0]" :pismena="l['pismena']" :key="l.id" :jeDokoncena="dokoncene.includes(l['id'])" :oznacena="o.is(l['id'])"
                :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" :cislo="l['cislo']" />
        </div>

        <h2>Horní řada</h2>
        <div v-if="lekce[0].length == 0" class="kategorie">
            <BlokLekce v-for="i in 5" pismena=". . ." :jeDokoncena="false" :cislo="i + 4" :key="i" />
        </div>
        <div v-else class="kategorie">
            <BlokLekce v-for="l in lekce[1]" :pismena="l['pismena']" :key="l.id" :jeDokoncena="dokoncene.includes(l['id'])" :oznacena="o.is(l['id'])"
                :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" :cislo="l['cislo']" />
        </div>

        <h2>Dolní řada</h2>
        <div v-if="lekce[0].length == 0" class="kategorie">
            <BlokLekce v-for="i in 3" pismena=". . ." :jeDokoncena="false" :cislo="i + 4 + 5" :key="i" />
        </div>
        <div v-else class="kategorie">
            <BlokLekce v-for="l in lekce[2]" :pismena="l['pismena']" :key="l.id" :jeDokoncena="dokoncene.includes(l['id'])" :oznacena="o.is(l['id'])"
                :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" :cislo="l['cislo']" />
        </div>

        <h2>Diakritika</h2>
        <div v-if="lekce[0].length == 0" class="kategorie">
            <BlokLekce v-for="i in 5" pismena=". . ." :jeDokoncena="false" :cislo="i + 4 + 5 + 3" :key="i" />
        </div>
        <div v-else class="kategorie">
            <BlokLekce v-for="l in lekce[3]" :pismena="l['pismena']" :key="l.id" :jeDokoncena="dokoncene.includes(l['id'])" :oznacena="o.is(l['id'])"
                :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" :cislo="l['cislo']" />
        </div>

        <h2>Závěr kurzu</h2>
        <div v-if="lekce[0].length == 0" class="kategorie">
            <BlokLekce v-for="i in 2" pismena=". . ." :jeDokoncena="false" :cislo="i + 4 + 5 + 3 + 5" :key="i" />
        </div>
        <div v-else class="kategorie">
            <BlokLekce v-for="l in lekce[4]" :pismena="l['pismena']" :key="l.id" :jeDokoncena="dokoncene.includes(l['id'])" :oznacena="o.is(l['id'])"
                :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" :cislo="l['cislo']" />
        </div>

        <h2>Pro programátory</h2>
        <div v-if="lekce[0].length == 0" class="kategorie">
            <BlokLekce v-for="i in 2" pismena=". . ." :jeDokoncena="false" :cislo="i + 4 + 5 + 3 + 5 + 2" :key="i" />
        </div>
        <div v-else class="kategorie">
            <BlokLekce v-for="l in lekce[5]" :pismena="l['pismena']" :key="l.id" :jeDokoncena="dokoncene.includes(l['id'])" :oznacena="o.is(l['id'])"
                :i="l['cislo']" :class="{ nohover: o.index.value != 0 }" :cislo="l['cislo']" />
        </div>

    </div>
</template>

<style scoped>
#pokracovani {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    border: none;
    width: 100%;
    height: 64px;
    align-self: center;
    padding: 12px 12px 12px 25px;
    color: white;
    font-size: 24px;
    font-weight: 500;
    transition-duration: 0.1s;
    margin-top: 5px;
}

#pokracovani:hover {
    background-color: var(--fialova);
    cursor: pointer;
}

#pokracovani>img {
    height: 32px;
    margin-right: 8px;
    margin-bottom: 1px;
}

#seznam {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}

#seznam>.kategorie {
    display: flex;
    flex-direction: column;
    gap: 20px;
    width: 100%;
}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

.nohover {
    pointer-events: none;
    transition-duration: 0.1s;
}

@media screen and (max-width: 1100px) {
    #seznam {
        width: 70vw;
        align-items: center;
    }

    h2 {
        align-self: start;
    }
    #pokracovani {
        height: 48px;
    }
    #pokracovani>img {
        width: 26px;
        margin-right: 11px;
    }
}
</style>