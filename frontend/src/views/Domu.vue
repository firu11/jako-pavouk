<script setup lang="ts">
import PsaniHomePage from '../components/PsaniHomePage.vue';
import { useHead } from 'unhead'

useHead({
    title: "",
    meta: [
        {
            name: "keywords",
            content: "jako, pavouk, psani, psaní, všemi, vsemi, všema, vsema, deseti, prsty, klavesnici, klavesnice, prst, klávesnice, klávesnici, rychle, rychlé, výuka, vyuka, zdarma, free, pavoucek, pociatci, pocitac, notebooku, notebook, stroj, stroji, na, online, program, aplikace, webovka, atf, procvicovani, procvičování, test",
        }
    ]
})
</script>

<template>
    <PsaniHomePage></PsaniHomePage>
    <h1 style="opacity: 0; user-select: none; margin: 0;">Jako Pavouk</h1> <!-- Pro SEO -->
    <div class="bloky">
        <img src="/klavesniceSPavoukem.png" alt="Pavouk s klávesnicí">
        <div class="text">
            <h2>Proč psát všemi deseti?</h2>
            <ul>
                <li>Ušetříte spoustu vašeho drahocenného času</li>
                <li>Snížíte chybovost při psaní</li>
                <li>Předejdete bolesti krku a hlavy z ustavičného koukání na klávesnici</li>
            </ul>
        </div>
    </div>
    <div class="bloky">
        <div class="text">
            <h2>Co vás to bude stát?</h2>
            <p>
                <span class="tlusty">Vůbec nic.</span> <br>
                Celý kurz Jako Pavouk je zcela zdarma!
                Avšak než si osvojíte všechny klávesy, bude to chvíli trvat.
                Ale od toho, aby to netrvalo věčnost, jsme tu my!
            </p>
        </div>
        <img src="/pavoukNaZlate.svg" alt="Bohatý pavouk">
    </div>
</template>

<style scoped></style>