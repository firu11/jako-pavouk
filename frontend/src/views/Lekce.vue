<script setup lang="ts">
import axios from 'axios'
import { onMounted, onUnmounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router'
import { Oznacene, format, getToken, napovedaKNavigaci } from '../utils';
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
const dokoncene = ref([] as { id: number, cpm: number, presnost: number }[])
const fetchProbehl = ref(false)
const o = new Oznacene()
const prvniNedokoncene = ref(1)

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

        if (dokoncene.value.length != cviceni.value.length && dokoncene.value.length != 0) {
            let dokoIds = dokoncene.value.map(a => a.id);
            for (let i = 0; i < cviceni.value.length; i++) {
                if (dokoIds.includes(cviceni.value[i].id)) prvniNedokoncene.value += 1
                else break
            }
        } else if (dokoncene.value.length == cviceni.value.length) {
            prvniNedokoncene.value = Math.ceil(cviceni.value.length/2)
        }
        
    }).catch(_ => {
        router.push('/404')
    })

    document.addEventListener('keydown', e1)
    document.addEventListener('keyup', e2)
    document.addEventListener('mousemove', zrusitVyber)
})

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
    document.removeEventListener('keyup', e2)
    document.removeEventListener('mousemove', zrusitVyber)
})

function e1(e: KeyboardEvent) {
    if (e.key == 'ArrowLeft' || e.key == 'ArrowUp') {
        e.preventDefault()
        if (o.index.value == 0) o.index.value = prvniNedokoncene.value + 1
        o.mensi()
    } else if (e.key == 'ArrowRight' || e.key == 'ArrowDown') {
        e.preventDefault()
        if (o.index.value == 0) o.index.value = prvniNedokoncene.value - 1
        o.vetsi()
    } else if (e.key == 'Enter') {
        e.preventDefault()
        let cvicE: HTMLElement | null = document.querySelector(`[i="true"]`)
        if (cvicE == null || o.bezOznaceni) {
            o.bezOznaceni = true
            o.index.value = prvniNedokoncene.value
        } else cvicE?.click()
    } else if (e.key == 'Tab') {
        e.preventDefault()
        napovedaKNavigaci()
    }
}

function e2(e: KeyboardEvent) {
    if (e.key == 'Enter') {
        e.preventDefault()
        let cvicE: HTMLElement | null = document.querySelector(`[i="true"]`)
        cvicE?.click()
    }
}

function zrusitVyber() {
    o.index.value = 0
}

function jeDokoncene(id: number) {
    for (const cvic of dokoncene.value) {
        if (cvic.id == id) return true
    }
    return false
}

function cvicID(id: number) {
    for (const cvic of dokoncene.value) {
        if (cvic.id == id) return cvic
    }
    return { id: 0, cpm: 0, presnost: 0 }
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
                :rychlost="cvicID(id).cpm" :presnost="cvicID(id).presnost" :i="index + 1 == o.index.value"
                :class="{ nohover: o.index.value != 0 }" :oznacena="index + 1 == o.index.value" />
        </div>
        <p v-else-if="cviceni.length == 0 && fetchProbehl">Tato lekce zatím nemá žádná cvičení</p>
    </div>
</template>

<style scoped>
h1 {
    direction: ltr;
}

.kontejnr {
    width: 101%;
    display: flex;
    gap: 15px;
    flex-wrap: wrap;
    justify-content: center;
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
