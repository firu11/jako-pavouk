<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { format, getToken, MojeMapa, pridatOznameni } from "../utils";
import SipkaZpet from "../components/SipkaZpet.vue";
import { onMounted, ref } from "vue";
import axios from "axios";
import Vysledek from "../components/Vysledek.vue";
import { useHead } from "@unhead/vue";
import Psani from "../components/Psani.vue";
import { getCas } from "../stores";

const router = useRouter()
const route = useRoute()
const pismena: string = Array.isArray(route.params.pismena) ? route.params.pismena[0] : route.params.pismena
const cislo: string = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

useHead({
    title: "Cvičení " + pismena
})

const text = ref([] as { id: number, znak: string, spatne: number, psat: boolean }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const nejcastejsiChyby = ref()

const posledni = ref(false)
const klavesnice = ref("")
const typTextu = ref("")

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)
const nacitamNovej = ref(false)

function get() {
    nacitamNovej.value = true
    axios.get("/cvic/" + encodeURIComponent(pismena) + "/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: true })
                delkaTextu.value++
            })
        })

        posledni.value = response.data.posledni
        klavesnice.value = response.data.klavesnice
        typTextu.value = response.data.typ
    }).catch(e => {
        console.log(e)
        pridatOznameni()
        router.back()
    }).finally(() => {
        nacitamNovej.value = false
    })
}

onMounted(() => {
    get()
})

function restart() {
    text.value = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]
    delkaTextu.value = 0

    get()
    konec.value = false
}

function konecTextu(o: number, p: number, n: MojeMapa, d: number) {
    opravenePocet.value = o
    preklepy.value = p
    nejcastejsiChyby.value = new MojeMapa(n)
    konec.value = true
    delkaNapsanehoTextu.value = d
}

async function prodlouzit() {
    nacitamNovej.value = true
    axios.get("/cvic/" + encodeURIComponent(pismena) + "/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        let pocetSlov = text.value.length
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[pocetSlov + i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: true })
                delkaTextu.value++
            })
        })
    }).catch(e => {
        console.log(e)
        pridatOznameni()
    }).finally(() => {
        nacitamNovej.value = false
    })

}
</script>

<template>
    <h1 class="nadpis-se-sipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Lekce: {{ format(pismena) }}
    </h1>
    <h2>Cvičení: {{ cislo }}</h2>

    <Psani v-if="!konec" @konec="konecTextu" @restart="restart" @prodlouzit="prodlouzit" :text :klavesnice :hide-klavesnice="false" :nacitamNovej
        :cas="getCas(typTextu)" :delkaTextu />

    <Vysledek v-else @restart="restart" :preklepy :opravenych="opravenePocet" :delkaTextu="delkaNapsanehoTextu" :cas="getCas(typTextu)" :cislo
        :posledni :pismena :nejcastejsiChyby />
</template>

<style scoped>

h2 {
    margin-top: 5px;
    color: rgb(194, 194, 194);
}

</style>