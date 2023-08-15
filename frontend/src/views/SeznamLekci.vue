<script setup lang="ts">
import axios from "axios"
import BlokLekce from "../components/BlokLekce.vue";
import { onMounted, ref } from "vue"
import { checkTeapot, getToken, pridatOznameni } from "../utils";
import { useHead } from 'unhead'

useHead({
    title: "Lekce"
})

const lekce = ref([[]])
const dokoncene = ref([])

onMounted(() => {
    const header = getToken() ? { headers: { Authorization: `Bearer ${getToken()}` } } : {}
    axios.get("/lekce", header)
        .then(response => {
            lekce.value = response.data.lekce
            dokoncene.value = response.data.dokoncene
        }).catch(e => {
            if (!checkTeapot(e)) {
                pridatOznameni()
                console.log(e)
            }
        })
})
</script>

<template>
    <h1>Lekce</h1>
    <div id="seznam">
        <h2>Střední řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 4" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[0]" :sus="l" :pismena="l['pismena']"
            :jeDokoncena="dokoncene.includes(l['id'])" />
        <h2>Horní řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 5" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[1]" :sus="l" :pismena="l['pismena']"
            :jeDokoncena="dokoncene.includes(l['id'])" />
        <h2>Dolní řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 3" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[2]" :sus="l" :pismena="l['pismena']"
            :jeDokoncena="dokoncene.includes(l['id'])" />
        <h2>Diakritika</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 4" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[3]" :sus="l" :pismena="l['pismena']"
            :jeDokoncena="dokoncene.includes(l['id'])" />
        <h2>Interpunkce</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 3" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[4]" :sus="l" :pismena="l['pismena']"
            :jeDokoncena="dokoncene.includes(l['id'])" />
    </div>
</template>

<style scoped>
#seznam {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

@media screen and (max-width: 1000px) {
    #seznam {
        width: 70vw;
        align-items: center;
    }
    h2 {
        align-self: start;
    }
}
</style>