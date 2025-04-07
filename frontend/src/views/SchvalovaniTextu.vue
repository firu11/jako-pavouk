<script setup lang="ts">
import { onMounted, ref, useTemplateRef } from 'vue';
import TextZadani from '../components/ucitel/TextZadani.vue';
import axios from 'axios';
import { getToken, pridatOznameni } from '../utils';

const textovePole = useTemplateRef("textove-pole")

const nazev = ref("")

onMounted(() => {
    get()
})

function get() {
    axios.get("/sus", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(resp => {
        textovePole.value!.text = resp.data.text
    }).catch(() => {
        pridatOznameni()
    })
}
</script>
<template>
    <div id="kontejner">
        <TextZadani ref="textove-pole" />

        <textarea ref="textarea" placeholder="Název" v-model="nazev" id="nazev" />
    </div>

    <button class="tlacitko">Potvrdit</button>
</template>
<style scoped>
#kontejner {
    position: relative;
    top: -60px;
    padding-top: 55px;
    width: 410px;
    height: 86vh;
}

#nazev {
    background-color: var(--tmave-fialova);
    border: 0;
    width: 100%;
    font-size: 16px !important;
    text-decoration: none;
    color: var(--bila);
    padding: 10px;
    resize: none;
    z-index: 2;
    line-height: 21px;
    font-family: "Red Hat Mono", monospace;
    border-radius: 5px;
    margin-top: 10px;
}

h1 {
    z-index: 1;
}
</style>