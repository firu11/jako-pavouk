<script setup lang="ts">
import { useHead } from 'unhead'
import { ref } from 'vue';

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

const copy = ref()

function copyEmail() {
    navigator.clipboard.writeText("firu@jakopavouk.cz").then(function () {
        copy.value.style.transform = "rotate(-25deg)"
        setTimeout(() => {
            copy.value.style.transform = "rotate(15deg)"
            setTimeout(() => {
                copy.value.style.transform = "rotate(0deg)"
            }, 110)
        }, 110)
    }, function (err) {
        console.error('Async: Could not copy text: ', err);
    });
}

</script>

<template>
    <h1>O nás</h1>

    <p>
        Jako Pavouk je studentský projekt zaměřený na výuku psaní všemi deseti.
        <br>
        Nápad a motivace pro vývoj této online aplikace vzešly z mé osobní zkušenosti,
        kdy jsem sám psát všemi deseti neuměl a rozhodl jsem se to naučit.
        <br>
        Když jsem ale začal hledat vhodné programy, zjistil jsem, že jich není mnoho a některé z nich stojí téměř 1000 Kč.
    </p>
    <p>
        Jelikož je Jako Pavouk kompletně zdarma, je tu možnost tento projekt podpořit přes platformu Ko-fi níže. Moc si
        vážím vaší podpory! <br>
        <a href="https://ko-fi.com/jakopavouk" target=”_blank” rel="nofollow"><img id="kofi"
                src="../assets/kofi_button.webp" alt="Podpořit" width="231" height="36"></a>
    </p>
    <p>
        Objevil/a jsi nějakou chybu/bug?<br> Máš nějaký zajímavý nápad jak vylepšit tuto aplikaci nebo nám jen chceš dát
        zpětnou vazbu?
        Neboj se mi napsat na e-mail:
        <span id="btn" @click="copyEmail">
            <span>firu</span>
            <span style="display:none;">nebudete mi websrapovat email vy smradlavý smrdi</span>
            <span>@</span>
            <span style="display:none;">nebudete mi websrapovat email vy smradlavý smrdi</span>
            <span>jakopavouk.cz</span>
            <button @click="copyEmail" ref="copy"><img src="../assets/icony/copy.svg" alt="Kopírovat"></button>
        </span>
    </p>

    <p>
        Pro zvědavé tu mám i link na GitHub, kde je celý projekt volně k prohlédnutí.<br>
        Budu moc rád za hvězdičku!
        <a href="https://github.com/Firu115/jako_pavouk" target=”_blank” rel="nofollow"><img id="github"
                src="../assets/github.webp" alt="GitHub" width="173" height="60"></a> <!-- random hash -->
    </p>
</template>

<style scoped>
#btn {
    margin: -5px;
    background-color: var(--fialova);
    padding: 10px 14px 6px 14px;
    border-radius: 5px;
    cursor: pointer;
    transition: 0.2s;
}

#btn:hover {
    background-color: var(--svetle-fialova);
}

button {
    background-color: transparent;
    border: 0px;
    transition: 0.1s ease-in-out;
    cursor: pointer;
}

button img {
    width: 20px;
    margin-left: 10px;
    vertical-align: middle;
    padding-bottom: 5px;
}

span {
    font-weight: 500;
}

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
#discord,
#github {
    width: 12em;
    height: auto;
    transition: 0.2s;
    cursor: pointer;
    text-align: center;
    color: var(--svetle-fialova);
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