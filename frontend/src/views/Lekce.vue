<script setup lang="ts">
import axios from 'axios'
import { onMounted, onUnmounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router'
import { Oznacene, formatovanyPismena, getToken } from '../utils';
import BlokCviceni from '../components/BlokCviceni.vue';
import SipkaZpet from '../components/SipkaZpet.vue';
import { useHead } from 'unhead'

const route = useRoute().params.pismena
const pismena = Array.isArray(route) ? route[0] : route // sus proste kdyby nahodou to byl array jakoze nebude tak to indexnem
const router = useRouter()

useHead({
    title: "Lekce " + pismena
})

const cviceni = ref([] as { id: number, typ: string }[])
const dokoncene = ref([] as any[])
const fetchProbehl = ref(false)
const o = new Oznacene()

onMounted(() => {
    axios.get('/lekce/' + encodeURIComponent(pismena), {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        if (response.data.cviceni === null) {
            router.push('/404')
        }
        cviceni.value = response.data.cviceni
        dokoncene.value = response.data.dokoncene
        fetchProbehl.value = true
        o.setMax(cviceni.value.length)
    }).catch(_ => {
        router.push('/404')
    })
    document.addEventListener('keydown', e1)
    document.addEventListener('mousemove', e2)
})

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
    document.removeEventListener('mousemove', e2)
})

function e1(e: KeyboardEvent) {
    if (e.key == 'ArrowLeft') {
        e.preventDefault()
        o.mensi()
    } else if (e.key == 'ArrowRight') {
        e.preventDefault()
        o.vetsi()
    } else if (e.key == 'Enter') {
        e.preventDefault()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        lekce?.click()
    }
}

function e2() {
    o.index.value = 0
}

function jeDokoncene(id: number) {
    for (const cvic of dokoncene.value) {
        if (cvic.Id == id) return true
    }
    return false
}

function cvicID(id: number) {
    for (const cvic of dokoncene.value) {
        if (cvic.Id == id) return cvic
    }
    return { Id: 0, Cpm: 0, Presnost: 0 }
}

function format(p: string) {
    if (p === "zbylá diakritika") return "Zbylá diakritika"
    else if (p === "velká písmena (shift)") return "Velká písmena (Shift)"
    return formatovanyPismena(p)
}

</script>

<template>
    <h1 class="nadpisSeSipkou">
        <SipkaZpet />
        Lekce: {{ format(pismena) }}
    </h1>
    <div class="kontejnr">
        <div v-if="cviceni.length !== 0 && fetchProbehl" v-for="({ id, typ }, index) in cviceni">
            <BlokCviceni :dokonceno="jeDokoncene(id)" :typ="typ" :index="index + 1" :pismena="pismena"
                :rychlost="cvicID(id).Cpm" :presnost="cvicID(id).Presnost" :fetchProbehl="fetchProbehl"
                :i="index + 1 == o.index.value" :class="{ nohover: o.index.value != 0 }"
                :oznacena="index + 1 == o.index.value" />
        </div>
        <div v-else-if="cviceni.length === 0 && !fetchProbehl" v-for="index in 5">
            <BlokCviceni :dokonceno="false" typ="..." :index="index" :pismena="pismena" :fetchProbehl="fetchProbehl" />
        </div>
        <p v-else>Tato lekce zatím nemá žádná cvičení</p>
    </div>
</template>

<style scoped>
h1 {
    direction: ltr;
}

.kontejnr {
    display: flex;
    gap: 15px;
    max-width: 700px;
    flex-wrap: wrap;
    justify-content: center;
    margin-top: 10px;
    margin: 10px 20px 10px 20px;
}

@media screen and (max-width: 1100px) {
    h1 {
        font-size: 1.7em;
        vertical-align: sub;
        max-width: 70%;
    }
}
</style>
