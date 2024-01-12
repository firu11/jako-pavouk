<script setup lang="ts">
import { useHead } from 'unhead'

useHead({
    title: "O nás",
    meta: [
        {
            name: "description",
            content: "Jako Pavouk je studentský projekt na výuku psaní všemi deseti. Obsahuje spoustu lekcí a cviční díky kterým si brzy osvojíte celou českou klávesnici.",
        }
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/o-nas"
        }
    ]
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
        se zděšením jsem zjistil, že co se týče české klávesnice, není moc možností, přičemž nějaké z nich vás budou stát
        skoro 1000 Kč.
    </p>
    <p>
        Jelikož je Jako Pavouk kompletně zdarma, je tu možnost tento projekt podpořit přes platformu Ko-fi níže. Moc si
        vážím vaší
        podpory! <br>
        <a href="https://ko-fi.com/jakopavouk" target=”_blank” rel="nofollow"><img id="kofi"
                src="../assets/kofi_button.webp" alt="Podpořit" width="231" height="36"></a>
    </p>
    <p>
        Objevil/a jsi nějakou chybu/bug? Máš nějaký zajímavý nápad jak vylepšit tuto aplikaci nebo snad chceš jen dát
        zpětnou
        vazbu? Využij náš Discord!<br>
        <a href="https://discord.gg/kMGRR6f3cb" target=”_blank” rel="nofollow"><img id="discord"
                src="../assets/discord.webp" alt="Discord" width="192" height="65"></a>
    </p>

    <p>
        Pro zvědavé tu mám i link na GitHub, kde je celý projekt volně k prohlédnutí.<br>
        <a href="https://github.com/gyarab/2023-4e-ruzicka-jako_pavouk" target=”_blank” rel="nofollow"><img id="github"
                src="../assets/github.webp" alt="GitHub" width="173" height="60"></a> <!-- random hash -->
    </p>
</template>

<style scoped>
p {
    margin-top: 25px;
    width: 110%;
    display: flex;
    flex-direction: column;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    align-items: center;
    padding: 30px 24px;
    gap: 20px;
}

#discord {
    width: 10em !important;
}

#github {
    width: 9em !important;
}

#kofi,
#discord,#github {
    width: 12em;
    height: auto;
    transition: 0.2s;
    cursor: pointer;
}

#kofi:hover,
#discord:hover,
#github:hover {
    transform: scale(1.05);
    transition: 0.2s;
}

@media screen and (max-width: 1100px) {
    p {
        width: 90%;
    }
}
</style>