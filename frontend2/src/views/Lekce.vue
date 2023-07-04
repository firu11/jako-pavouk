<script setup lang="ts">
import axios from 'axios'
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router'
import { formatovany_pismena, get_token } from '../utils';
import BlokCviceni from '../components/BlokCviceni.vue';

const route = useRoute().params.pismena
const pismena = Array.isArray(route) ? route[0] : route // sus proste kdyby nahodou to byl array jakoze nebude tak to indexnem
const router = useRouter()

const cviceni = ref([] as {id: number, typ: string}[])
const dokoncene = ref([] as number[])
const fetch_probehl = ref(false)

onMounted(() => {
    axios.get('/lekce/' + pismena, {
        headers: {
            Authorization: `Bearer ${get_token()}`
        }
    }).then(response => {
        if (response.data.cviceni === null) {
            router.push('/404')
        }
        cviceni.value = response.data.cviceni
        dokoncene.value = response.data.cviceni
        fetch_probehl.value = true

    }).catch(_ => {
        router.push('/404')
    })
})

</script>

<template>
    <h1>
        <router-link class="tlacZpet" :to="'/lekce'">
            <img src="../assets/icony/sipkaL.svg" alt="Zpět">
        </router-link>
        Lekce: {{ formatovany_pismena(pismena) }}
    </h1>
    <div class="kontejnr">
        <div v-if="cviceni.length !== 0 && fetch_probehl" v-for="({id, typ}, index) in cviceni">
            <BlokCviceni :dokonceno="dokoncene.includes(id)" :typ="typ" :index="index + 1" :pismena="pismena"/>
        </div>
        <div v-else-if="cviceni.length === 0 && !fetch_probehl" v-for="index in 5">
            <BlokCviceni :dokonceno="false" typ="..." :index="index" :pismena="pismena"/>
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
}

h1 {
    display: inline-flex;
    position: relative;
    right: 25px;
    /* posunuti o pulku sipky */
}
</style>
