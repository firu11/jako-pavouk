<script setup lang="ts">
import { clone, getToken, MojeMapa, pridatOznameni, saveNastaveni } from "../utils";
import { computed, onMounted, ref, toRaw } from "vue";
import axios from "axios";
import Vysledek from "../components/Vysledek.vue";
import { useHead } from "@unhead/vue";
import Psani from "../components/Psani.vue";
import { mobil, okZnaky } from "../stores";
import { useRouter } from "vue-router";
import NastaveniBtn from "../components/NastaveniBtn.vue";
import PsaniMenu from "../components/PsaniMenu.vue";
import SipkaZpet from '../components/SipkaZpet.vue';

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

const text = ref([] as { id: number, znak: string, spatne: number, psat: boolean }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const nejcastejsiChyby = ref()

const psaniRef = ref()
const menuRef = ref()

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)
const nacitamNovej = ref(false)

const hideKlavecnice = ref(false)

function get() {
    nacitamNovej.value = true
    axios.post("/test-psani",
        {
            typ: menuRef.value.typ ? "vety" : "slova",
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
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: !okZnaky.test(pismeno) })
                delkaTextu.value++
            })
        })

        if (menuRef.value == null) return

        loadAlternativy()
        toggleDiakritikaAVelkaPismena()

        if (response.data.klavesnice != undefined) menuRef.value.klavModel = response.data.klavesnice == "qwerty"

    }).catch(e => {
        console.log(e)
        pridatOznameni()
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
    get()
})

function restart() {
    delkaTextu.value = 0
    text.value = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]

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

const klavesnice = computed(() => {
    if (menuRef.value == undefined || (delkaTextu.value == 0 && puvodniText.length == 0)) return ""
    return menuRef.value.klavModel ? "qwerty" : "qwertz"
})

let puvodniText = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]
let textBezDiakritiky = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]
let textMalym = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]
let textOboje = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]

async function toggleDiakritikaAVelkaPismena() {
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

    axios.post("/test-psani",
        {
            typ: menuRef.value.typ ? "vety" : "slova",
        },
        {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }
    ).then(response => {
        if (menuRef.value == null) return

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
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[pocetSlov + i].push({ id: delkaTextu.value, znak: pismeno, spatne: 0, psat: !okZnaky.test(pismeno) })
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

function refocus() {
    psaniRef.value.focusInput()
}

</script>

<template>
    <h1 class="nadpis-se-sipkou" style="margin: 0; direction: ltr;">
        <SipkaZpet />
        Test psaní
    </h1>
    <h2>Typ: {{ menuRef != undefined ? (menuRef.typ ? "Věty" : "Slova") : " " }}</h2>

    <Psani v-if="!konec" @konec="konecTextu" @restart="restart" @pise="hideKlavecnice = false" @prodlouzit="prodlouzit"
        :text="text" :klavesnice="klavesnice" :delkaTextu="delkaTextu" :hide-klavesnice="hideKlavecnice"
        :nacitam-novej="nacitamNovej" :cas="menuRef == undefined ? 15 : menuRef.delka" ref="psaniRef" />

    <Vysledek v-else @restart="restart" :preklepy="preklepy" :opravenych="opravenePocet"
        :delkaTextu="delkaNapsanehoTextu" :cas="menuRef == undefined ? 15 : menuRef.delka" :cislo="'test-psani'"
        :posledni="true" :nejcastejsiChyby="nejcastejsiChyby" />

    <PsaniMenu class="psaniMenu" :class="{ hide: konec || !hideKlavecnice }" @restart="restart(); psaniRef.restart()"
        @toggle="toggleDiakritikaAVelkaPismena" @click="refocus" :bez-stinu="psaniRef?.fullHideKlavesnice"
        ref="menuRef" />

    <NastaveniBtn v-if="!konec && klavesnice != ''"
        @klik="hideKlavecnice = !hideKlavecnice; refocus(); psaniRef?.restart()" />
</template>

<style scoped>
.psaniMenu {
    transition: opacity 0.2s;
}

.hide {
    opacity: 0;
    user-select: none;
    z-index: -1000;
}

h2 {
    margin-top: 5px;
    color: rgb(194, 194, 194);
}
</style>
