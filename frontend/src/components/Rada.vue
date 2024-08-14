<script setup lang="ts">
import { useRouter } from 'vue-router';
import { mobil, prihlasen, role } from "../stores"

const router = useRouter()

defineProps({
    pocetDoko: {
        type: Number,
        default: -1
    }
})

function hop() {
    router.push("/jak-psat")
}
</script>

<template>
    <div v-if="!mobil && (!prihlasen || (pocetDoko < 2 && pocetDoko != -1))" id="box">
        <div>
            <h2>Nevíš jak na to?</h2>
            <p>Mrkni na trochu teorie! Základem je výchozí pozice.</p>
            <button class="tlacitko" :onclick="hop">Jak psát</button>
        </div>
        <img src="../assets/theThinker.svg" alt="Pavouk the thinker" width="180" height="157">
    </div>
    <div v-else-if="!mobil && pocetDoko >= 12 - 2 && role == 'basic'" id="box">
        <div>
            <h2>Líbí se ti projekt?</h2>
            <p>Pokud už píšeš Jako Pavouk, oceníme finanční podporu!</p>
            <a href="https://ko-fi.com/jakopavouk" target="_blank" rel="nofollow">
                <img id="kofi" src="../assets/kofi_button.webp" alt="Podpořit" width="230" height="36">
            </a>
        </div>
        <img src="../assets/theLover.svg" alt="Pavouk the thinker" width="180" height="157" style="transform: scale(1.15);"> <!-- idk proc je mensi -->
    </div>
</template>

<style scoped>
a {
    margin-top: 10px;
    height: 36px;
}

#kofi {
    width: auto !important;
    transition: 0.2s;
}

#kofi:hover {
    transform: scale(1.05);
}

#box {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 12px 15px 12px 25px;
    display: flex;
    width: 500px;
    align-items: center;
    justify-content: space-between;
    margin-top: 5px;
}

#box div {
    display: flex;
    flex-direction: column;
}

h2 {
    font-size: 1.6rem;
    margin-bottom: 12px;
    font-weight: 400;
    width: 100%;
}

p {
    font-size: 1.15rem;
    margin-bottom: 5px;
}

img {
    width: 180px;
}
</style>