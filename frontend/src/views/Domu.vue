<script setup lang="ts">
import PsaniHomePage from "../components/PsaniHomePage.vue";
import { useHead } from "unhead"
import { useRouter } from "vue-router";
import { mobil, prihlasen } from "../stores";

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
        router.push("/lekce")
    } else {
        router.push("/prvni-psani")
    }
}

</script>

<template>
    <PsaniHomePage />
    <h1 id="schovat">Psaní všemi deseti zdarma Jako Pavouk</h1> <!-- Pro SEO (nevim jestli funguje) -->
    <div class="bloky">
        <div class="text">
            <p style="font-size: 1.3em; line-height: 1.5em;">Webová aplikace na výuku <br><b
                    style="font-weight: 700;">psaní všemi deseti</b> zdarma.<br>Chceš se naučit psát jako pavouk?</p>
            <button v-if="!mobil" class="tlacitko" @click="zacit()"
                style="transform: scale(1.25); margin-top: 2em;">Začít psát</button>
        </div>
        <img src="../assets/pavoukStudent.svg" alt="Pavouk student" width="300" height="253" style="max-height: 230px;">
    </div>
    <div class="bloky">
        <img src="../assets/klavesniceSPavoukem.webp" alt="Pavouk s klávesnicí" width="300" height="223">
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
                Celý kurz psaní všemi deseti je zcela <b>ZDARMA</b>! <br>
                Avšak než si osvojíš všechny klávesy, bude to chvíli trvat.
            </p>
        </div>
        <img src="../assets/pavoukNaZlate.svg" width="290" height="182" alt="Bohatý pavouk" style="width: 290px;">
    </div>
    <div class="bloky">
        <img src="../assets/pavoukSkola.svg" width="290" height="182" alt="Bohatý pavouk" style="width: 290px;">
        <div class="text">
            <h2>Systém pro školy</h2>
            <p>
                Lorem ipsum dolor sit, amet consectetur adipisicing elit. Exercitationem assumenda ratione dolore voluptate in hic.
            </p>
            <button class="tlacitko" @click="router.push('/info-skola')">Číst dál</button>
        </div>
    </div>
</template>

<style scoped>
.bloky img {
    width: 300px;
}

@media screen and (max-width: 1100px) {
    .tlacitko:nth-of-type(1) {
        transform: scale(1) !important;
        margin: 1em 0 !important;
    }

    .text:nth-of-type(1) p {
        font-size: 1.2em !important;
    }

    .bloky img {
        width: 90% !important;
    }
}
</style>