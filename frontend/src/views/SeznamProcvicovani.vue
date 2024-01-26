<script setup lang="ts">
import { useHead } from 'unhead'
import { Oznacene, checkTeapot, pridatOznameni } from '../utils';
import axios from 'axios';
import { onMounted, onUnmounted, ref } from 'vue';

useHead({
    title: "Procvičování",
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/procvic"
        }
    ]
})

const texty = ref([])
const mobil = ref(document.body.clientWidth <= 1000)
const o = new Oznacene()

onMounted(() => {
    axios.get("/procvic")
        .then(response => {
            texty.value = response.data.texty
            o.setMax(texty.value.length)
        }).catch(e => {
            if (!checkTeapot(e)) {
                pridatOznameni()
            }
        })
    document.addEventListener('keydown', e1)
    document.addEventListener('mousemove', e2)
})

function e1(e: KeyboardEvent) {
    if (e.key == 'ArrowUp') {
        e.preventDefault()
        o.mensi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 500 })
    } else if (e.key == 'ArrowDown') {
        e.preventDefault()
        o.vetsi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 200 })
    } else if (e.key == 'Enter') {
        e.preventDefault()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        lekce?.click()
    }
}

function e2() {
    o.index.value = 0
}

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
    document.removeEventListener('mousemove', e2)
})

</script>

<template>
    <h1>Procvičování</h1>
    <div id="seznam">
        <h2>Texty</h2>
        <div v-if="texty.length == 0" v-for="_ in 3" class="blok">
            <h2>. . .</h2>
        </div>
        <RouterLink v-else-if="!mobil" v-for="t, i in texty" :to="`/procvic/${i + 1}`" class="blok"
            :i="i + 1 == o.index.value" :class="{ oznacene: i + 1 == o.index.value, nohover: o.index.value != 0 }">
            <h2>{{ t }}</h2>
        </RouterLink>
        <div v-else v-for="t in texty" class="blok" @click="pridatOznameni('Psaní na telefonech zatím neučíme...')">
            <h2>{{ t }}</h2>
        </div>
        <h2>Texty na míru</h2>
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

.blok:hover,
.oznacene {
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