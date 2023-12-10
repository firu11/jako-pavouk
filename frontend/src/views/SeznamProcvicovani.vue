<script setup lang="ts">
import { useHead } from 'unhead'
import { checkTeapot, pridatOznameni } from '../utils';
import axios from 'axios';
import { onMounted, ref } from 'vue';

useHead({
    title: "Procvičování"
})

const texty = ref([])
const mobil = ref(document.body.clientWidth <= 1000)

onMounted(() => {
    axios.get("/procvic")
        .then(response => {
            texty.value = response.data.texty
        }).catch(e => {
            if (!checkTeapot(e)) {
                pridatOznameni()
            }
        })
})

</script>

<template>
    <h1>Procvičování</h1>
    <div id="seznam">
        <h2>Texty</h2>
        <RouterLink v-if="!mobil" v-for="t, i in texty" :to="`/procvic/${i+1}`" class="blok">
            <h2>{{ t }}</h2>
        </RouterLink>
        <div v-else v-for="t in texty" class="blok" @click="pridatOznameni('Psaní na telefonech zatím neučíme...')">
            <h2>{{ t }}</h2>
        </div>
        <h2>Písmena na míru</h2>
        <div class="blok">
            <h2>Pavouci už na tom pilně pracují...</h2>
        </div>
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

.blok {
    display: flex;
    color: var(--bila);
    padding: 12px 20px 12px 25px;
    text-decoration: none;
    border-radius: 10px;
    width: 500px;
    background-color: var(--tmave-fialova);
    min-height: 64px;
    transition-duration: 0.2s;
    cursor: pointer;
    user-select: none;
    /* kvuli tomu neprihlasenymu */
}

.blok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.blok h2 {
    font-size: 24px;
    font-weight: 300;
    margin: 0;
    align-self: center;
}

@media screen and (max-width: 1100px) {
    #seznam {
        width: 70vw;
        align-items: center;
    }

    h2 {
        align-self: start;
    }
    .blok {
        min-width: 260px;
        width: 100%;
        background-color: var(--tmave-fialova);
        min-height: 64px;
        max-height: 100px;
        height: auto;
        transition-duration: 0.2s;

        /* kvuli tomu neprihlasenymu */
        cursor: pointer; 
    }
}
</style>