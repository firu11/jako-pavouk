<script setup lang="ts">
import { useHead } from 'unhead'
import { Oznacene, checkTeapot, getToken, napovedaKNavigaci, pridatOznameni } from '../utils';
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
const rychlosti = ref([] as number[])
const mobil = ref(document.body.clientWidth <= 1000)
const o = new Oznacene()
let randomCvic = 1

onMounted(() => {
    axios.get("/procvic", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        texty.value = response.data.texty
        rychlosti.value = response.data.rychlosti
        o.setMax(texty.value.length + 1)
        randomCvic = Math.floor(Math.random() * texty.value.length) + 1
    }).catch(e => {
        if (!checkTeapot(e)) {
            pridatOznameni()
        }
    })
    document.addEventListener('keydown', e1)
    document.addEventListener('keyup', e2)
    document.addEventListener('mousemove', zrusitVyber)
})

function e1(e: KeyboardEvent) {
    if (e.key == 'ArrowUp' || e.key == 'ArrowLeft') {
        e.preventDefault()
        o.mensi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 500 })
    } else if (e.key == 'ArrowDown' || e.key == 'ArrowRight') {
        e.preventDefault()
        o.vetsi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 200 })
    } else if (e.key == 'Enter') {
        e.preventDefault()
        let cvicE: HTMLElement | null = document.querySelector(`[i="true"]`)
        if (cvicE == null || o.bezOznaceni) {
            o.bezOznaceni = true
            o.index.value = randomCvic
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

function mobilKlik(e: MouseEvent) {
    e.preventDefault()
    pridatOznameni('Psaní na telefonech zatím neučíme...')
}

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
    document.removeEventListener('keyup', e2)
    document.removeEventListener('mousemove', zrusitVyber)
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
            <span v-if="rychlosti[i + 1] != -1"><b>{{ Math.round(rychlosti[i + 1] * 10) / 10 }}</b> CPM</span>
        </RouterLink>
        <div v-else v-for="t in texty" class="blok" @click="mobilKlik">
            <h2>{{ t }}</h2>
        </div>
        <h2>Na míru</h2>
        <RouterLink v-if="!mobil" :to="'/test-psani'" class="blok" :i="4 == o.index.value"
            :class="{ oznacene: 4 == o.index.value, nohover: o.index.value != 0 }">
            <h2>Test psaní</h2>
            <span v-if="texty.length != 0 && rychlosti[0] != -1"><b>{{ Math.round(rychlosti[0] * 10) / 10 }}</b>
                CPM</span>
        </RouterLink>
        <a v-else href="/test-psani" class="blok" :i="4 == o.index.value" style="user-select: none;" @click="mobilKlik">
            <h2>Test psaní</h2>
            <span v-if="texty.length != 0 && rychlosti[0] != -1"><b>{{ Math.round(rychlosti[0] * 10) / 10 }}</b> CPM</span>
        </a>
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
    justify-content: space-between;

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

.blok span {
    font-size: 1.2rem;
    display: flex;
    align-items: baseline;
    gap: 5px;
    justify-content: end;
    height: 34px;
    align-self: center;
}

.blok span b {
    font-size: 1.8rem;
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
        min-height: 48px;
        max-height: 100px;
        height: auto;
        transition-duration: 0.2s;

        /* kvuli tomu neprihlasenymu */
        cursor: pointer;
    }

    .blok span b {
        font-size: 1.2rem;
    }

    .blok span {
        font-size: 0.8rem;
        top: 13px;
        gap: 3px;
        height: 22px;
    }

    .blok h2 {
        font-size: 1.3rem;
    }
}
</style>