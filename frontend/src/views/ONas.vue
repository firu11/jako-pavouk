<script setup lang="ts">
import { useHead } from 'unhead'

useHead({
    title: "O nás"
})
</script>

<template>
    <h1>O nás</h1>

    <p>
        Jako Pavouk je studentský projekt na výuku psaní všemi deseti.
        <br>
        Nápad a motivace pro vývoj této webové aplikace byl fakt,
        že já sám jsem neuměl psát všemi deseti, a když jsem se to rozhodl naučit,
        <br>
        se zděšením jsem zjistil, že co se týče české klávesnice, není moc možností, přičemž nějaké z nich vás budou stát skoro 1000 Kč.
    </p>
    <p>
        Jelikož je Jako Pavouk kompletně zdarma, je tu možnost tento projekt podpořit přes platformu Ko-fi níže. Moc si vážím vaší
        podpory! <br>
        <a href="https://ko-fi.com/jakopavouk" target=”_blank”><img id="kofi" src="../assets/kofi_button.webp" alt="Podpořit"></a>
    </p>
    <p>
        Objevili jste nějakou chybu/bug? Máte nějaký zajímavý nápad jak vylepšit tuto aplikaci nebo snad chcete jen dát zpětnou
        vazbu? Využijte náš Discord! <br>
        <a href="https://discord.gg/kMGRR6f3cb" target=”_blank”><img id="discord" src="../assets/discord.png" alt="Discord"></a>
    </p>
</template>

<style scoped>
p {
    margin-bottom: 25px;
    width: 110%;
    display: flex;
    flex-direction: column;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    align-items: center;
    padding: 30px 20px;
}

#discord {
    width: 10em !important;
}

#kofi, #discord {
    width: 12em;
    transition: 0.2s;
    cursor: pointer;
    margin-top: 20px;
}

#kofi:hover, #discord:hover {
    filter: contrast(65%);
    transition: 0.2s;
}

@media screen and (max-width: 1000px) {
    p {
        width: 90%;
    }
}
</style>