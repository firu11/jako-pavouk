<script setup lang="ts">
import axios from 'axios'
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router'
import { formatovanyPismena, getToken } from '../utils';
import BlokCviceni from '../components/BlokCviceni.vue';
import SipkaZpet from '../components/SipkaZpet.vue';
import { useHead } from 'unhead'

const route = useRoute().params.pismena
const pismena = Array.isArray(route) ? route[0] : route // sus proste kdyby nahodou to byl array jakoze nebude tak to indexnem
const router = useRouter()

useHead({
    title: "Lekce " + route
})

const cviceni = ref([] as {id: number, typ: string}[])
const dokoncene = ref([] as number[])
const fetchProbehl = ref(false)

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

    }).catch(_ => {
        router.push('/404')
    })
})

</script>

<template>
    <h1 class="nadpisSeSipkou">
        <SipkaZpet/>
        Lekce: {{ formatovanyPismena(pismena) }}
    </h1>
    <div class="kontejnr">
        <div v-if="cviceni.length !== 0 && fetchProbehl" v-for="({id, typ}, index) in cviceni">
            <BlokCviceni :dokonceno="dokoncene.includes(id)" :typ="typ" :index="index + 1" :pismena="pismena" :fetchProbehl="fetchProbehl"/>
        </div>
        <div v-else-if="cviceni.length === 0 && !fetchProbehl" v-for="index in 5">
            <BlokCviceni :dokonceno="false" typ="..." :index="index" :pismena="pismena" :fetchProbehl="fetchProbehl"/>
        </div>
        <p v-else>Tato lekce zatím nemá žádná cvičení</p>
    </div>
</template>

<style scoped>
.kontejnr {
    display: flex;
    gap: 15px;
    max-width: 700px;
    flex-wrap: wrap;
    justify-content: center;
    margin-top: 10px;
}
</style>
