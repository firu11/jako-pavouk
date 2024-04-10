<script setup lang="ts">
import PsaniHomePage from '../components/PsaniHomePage.vue';
import { useHead } from 'unhead'
import { useRouter } from 'vue-router';
import { prihlasen } from '../stores';

const router = useRouter()

useHead({
    title: "",
    meta: [
        {
            name: "description",
            content: "Webová aplikace na výuku psaní všemi deseti zdarma. Chceš se naučit psát rychle Jako Pavouk? Buď pánem své klávesnice!",
        },
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz"
        }
    ]
})

function zacit() {
    if (prihlasen.value) {
        router.push('/lekce')
    } else {
        router.push('/prvni-psani')
    }
}
</script>

<template>
    <PsaniHomePage></PsaniHomePage>
    <h1 id="schovat">Psaní všemi deseti Jako Pavouk</h1> <!-- Pro SEO (nevim jestli funguje) -->
    <div class="bloky">
        <div class="text">
            <p style="font-size: 1.3em; line-height: 1.5em;">Webová apllikace na výuku <br><b>psaní všemi deseti</b>.<br>Chceš se naučit psát jako pavouk?</p>
            <button class="tlacitko" @click="zacit()" style="transform: scale(1.2); margin-top: 2em;">Začít psát</button>
        </div>
        <img src="../assets/pavoukStudent.svg" alt="Pavouk s klávesnicí" width="300" height="300">
    </div>
    <div class="bloky">
        <img src="../assets/klavesniceSPavoukem.webp" alt="Pavouk s klávesnicí" width="285" height="220">
        <div class="text">
            <h2>Proč psát všemi deseti?</h2>
            <ul>
                <li>Ušetříš spoustu svého drahocenného času</li>
                <li>Snížíš chybovost při psaní</li>
                <li>Předejdeš bolesti krku a hlavy z ustavičného koukání na klávesnici</li>
            </ul>
        </div>
    </div>
    <div class="bloky">
        <div class="text">
            <h2>Co tě to bude stát?</h2>
            <p>
                <span class="tlusty"><b>Vůbec nic.</b></span> <br>
                Celý kurz psaní všemi deseti je zcela <b>zdarma</b>! <br>
                Avšak než si osvojíš všechny klávesy, bude to chvíli trvat.
            </p>
        </div>
        <img src="../assets/pavoukNaZlate.svg" width="287" height="180" alt="Bohatý pavouk">
    </div>
</template>

<style scoped>
@media screen and (max-width: 1100px) {
    .tlacitko:nth-of-type(1) {
        transform: scale(1) !important;
        margin: 1em 0 !important;
    }
    .text:nth-of-type(1) p {
        font-size: 1.2em !important;
    }
    .bloky img {
        max-height: 230px;
        width: auto;
    }
}
</style>