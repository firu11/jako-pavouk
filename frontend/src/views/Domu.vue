<script setup lang="ts">
import PsaniHomePage from '../components/PsaniHomePage.vue';
</script>

<template>
    <PsaniHomePage></PsaniHomePage>
    <h1 style="opacity: 0; user-select: none; margin: 0;">Jako Pavouk</h1> <!-- Pro SEO -->
    <div class="bloky">
        <img src="/klavesniceSPavoukem.png" alt="Pavouk s klávesnicí">
        <div class="text">
            <h2>Proč psát všema deseti?</h2>
            <ul>
                <li>Ušetříte spoustu vašeho drahocenného času</li>
                <li>Snižuje chybovost při psaní</li>
                <li>Předchází bolesti krku a hlavy z ustavičného koukaní na klávesnici</li>
            </ul>
        </div>
    </div>
    <div class="bloky">
        <div class="text">
            <h2>Co vás to bude stát?</h2>
            <p>
                <span class="tlusty">Vůbec nic.</span> <br> 
                Celý kurz Jako Pavouk je zcela zdarma! 
                Avšak než si osvojíte všechny klávesy, budete to chvíli trvat. 
                Od toho, aby to ale netrvalo věčnost, jsme tu my!
            </p>
        </div>
        <img src="/pavoukNaZlate.svg" alt=""> <!-- TODO -->
    </div>
</template>

<style scoped>
.text {
    width: 60%;
    list-style-position: inside;
}

.text li {
    text-align: left;
    text-indent: -1em;
    margin-left: 10px;
}

.text p, .text li {
    font-size: 1.1em;
    padding: 8px;
    line-height: 1.3em;
}

.text h2 {
    margin-bottom: 20px;
    font-size: 1.6em;
}

img {
    width: 40%;
    padding: 0 25px;
}

.bloky {
    margin-top: 80px;
    width: 110%;
    height: 20em;
    display: flex;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    align-items: center;
    padding: 20px;
}

@media screen and (max-width: 1000px) {
    .bloky {
        margin-top: 25px;
        flex-direction: column-reverse;
        width: 80%;
        height: auto;
        padding: 25px 15px;
        gap: 20px;
    }

    div:nth-of-type(3).bloky {
        flex-direction: column;
    }

    .text {
        width: 100%;
        height: 60%;
    }

    img {
        width: 100%;
        height: 40%;
    }
}
</style>