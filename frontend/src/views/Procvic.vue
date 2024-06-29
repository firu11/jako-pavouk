<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { checkTeapot, clone, getCisloProcvic, getToken, MojeMapa, pridatOznameni, saveNastaveni } from '../utils';
import SipkaZpet from '../components/SipkaZpet.vue';
import { computed, onMounted, ref, toRaw } from 'vue';
import axios from 'axios';
import Vysledek from '../components/Vysledek.vue';
import { useHead } from '@unhead/vue';
import Psani from '../components/Psani.vue';
import NastaveniBtn from "../components/NastaveniBtn.vue";
import PsaniMenu from "../components/PsaniMenu.vue";
import { mobil, nastaveniJmeno } from '../stores';

const router = useRouter()
const route = useRoute()
const id: string = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

useHead({
    title: "Procvičování" // po fetchi změnim
})

const text = ref([] as { id: number, znak: string, spatne: number, }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const cas = ref(0)
const nejcastejsiChyby = ref()

const jmeno = ref(". . .")

const psaniRef = ref()
const menuRef = ref()

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)
const nacitamNovej = ref(false)

const hideKlavecnice = ref(false)

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

function get() {
    nacitamNovej.value = true
    axios.get("/procvic/" + id + "/" + getCisloProcvic(id), {
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
        jmeno.value = response.data.jmeno

        loadAlternativy()
        toggleDiakritikaAVelkaPismena()

        if (response.data.klavesnice != undefined) menuRef.value.klavModel = response.data.klavesnice == "qwerty"

        useHead({
            title: jmeno.value
        })
    }).catch(e => {
        if (!checkTeapot(e)) {
            pridatOznameni()
            router.back()
        }

    }).finally(() => {
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
        menuRef.value.diakritika = obj.diakritika
        menuRef.value.velkaPismena = obj.velkaPismena
        menuRef.value.typ = obj.typ
        menuRef.value.delka = obj.delka
        menuRef.value.klavModel = obj.klavesnice
    }
    get()
})

function restart() {
    text.value = [] as { id: number, znak: string, spatne: number, }[][]
    delkaTextu.value = 0

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

const klavesnice = computed(() => {
    if (menuRef.value == undefined) return ""
    return menuRef.value.klavModel ? "qwerty" : "qwerty"
})

let puvodniText = [] as { id: number, znak: string, spatne: number }[][]
let textBezDiakritiky = [] as { id: number, znak: string, spatne: number }[][]
let textMalym = [] as { id: number, znak: string, spatne: number }[][]
let textOboje = [] as { id: number, znak: string, spatne: number }[][]

function toggleDiakritikaAVelkaPismena() {
    if (!menuRef.value.diakritika && !menuRef.value.velkaPismena) {
        text.value = clone(textOboje)
    } else if (!menuRef.value.diakritika) {
        text.value = clone(textBezDiakritiky)
    } else if (!menuRef.value.velkaPismena) {
        text.value = clone(textMalym)
    } else {
        text.value = clone(puvodniText)
    }
    psaniRef.value?.restart()
    saveNastaveni(menuRef.value.diakritika, menuRef.value.velkaPismena, menuRef.value.typ, menuRef.value.delka, menuRef.value.klavModel)
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
    nacitamNovej.value = true

    axios.get("/procvic/" + id + "/" + getCisloProcvic(id), {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        if (!menuRef.value.diakritika && !menuRef.value.velkaPismena) {
            for (let i = 0; i < response.data.text.length; i++) {
                response.data.text[i] = response.data.text[i].normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLocaleLowerCase()
            }
        } else if (!menuRef.value.diakritika && menuRef.value.velkaPismena) {
            for (let i = 0; i < response.data.text.length; i++) {
                response.data.text[i] = response.data.text[i].normalize("NFD").replace(/[\u0300-\u036f]/g, "")
            }
        } else if (menuRef.value.diakritika && !menuRef.value.velkaPismena) {
            for (let i = 0; i < response.data.text.length; i++) {
                response.data.text[i] = response.data.text[i].toLocaleLowerCase()
            }
        }

        let pocetSlov = text.value.length

        let lastSlovo = text.value[pocetSlov - 1]
        let lastPismeno = lastSlovo[lastSlovo.length - 1]
        if (lastPismeno.znak != " ") {
            delkaTextu.value++
            text.value[pocetSlov - 1].push({ id: delkaTextu.value, znak: " ", spatne: 0 })
        }

        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[pocetSlov + i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0 })
                delkaTextu.value++
            })
        })
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni()
        }
    }).finally(() => {
        nacitamNovej.value = false
    })
}

</script>

<template>
    <h1 class="nadpisSeSipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        {{ jmeno }}
    </h1>

    <Psani v-if="!konec" @konec="konecTextu" @restart="restart" @pise="hideKlavecnice = false" @prodlouzit="prodlouzit" :text="text"
        :klavesnice="klavesnice" :delkaTextu="delkaTextu" :hide-klavesnice="hideKlavecnice" :nacitam-novej="nacitamNovej"
        :cas="menuRef == undefined ? 15 : menuRef.delka" ref="psaniRef" />

    <Vysledek v-else @restart="restart" :preklepy="preklepy" :opravenych="opravenePocet" :delkaTextu="delkaNapsanehoTextu" :casF="casFormat" :cas="cas"
        :cislo="id" :posledni="true" :nejcastejsiChyby="nejcastejsiChyby" />

    <PsaniMenu :class="{ hide: konec || !hideKlavecnice }" @restart="restart(); psaniRef.restart()" @toggle="toggleDiakritikaAVelkaPismena"
        :vyberTextu="false" ref="menuRef" />

    <NastaveniBtn v-if="!konec" @klik="hideKlavecnice = !hideKlavecnice" />
</template>

<style scoped>
.hide {
    opacity: 0;
    z-index: -1000;
    user-select: none;
}
</style>