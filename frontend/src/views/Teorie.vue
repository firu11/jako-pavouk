<script setup lang="ts">
import { useHead } from 'unhead'

useHead({
    title: "Jak psát",
    meta: [
        {
            name: "description",
            content: "Teorie psaní všemi deseti na české klávesnici",
        },
        {
            name: "keywords",
            content: "psani, psaní, všemi, vsemi, všema, vsema, deseti, teorie, jak, psát, klávesy, princip",
        }
    ]
})
</script>

<template>
    <h1>Jak psát všemi deseti</h1>

    <div class="bloky">
        <div class="text">
            <h2>Výchozí pozice</h2>
            <p>
                Základním principem je mít prsty ve správné výchozí pozici, která je velice prostá. <br> Ukazováčky položíme
                na
                <span class="tlusty">písmena F, J</span>, kde často bývá fyzická značka, abyste se jich mohli držet.
                Ostatní prsty umístíme <span class="tlusty">do stran ve stejné řadě</span>.
                Tedy levý prostředníček bude na D, prsteníček na S a malíček na A. Zrcadlově na druhé ruce.
            </p>
        </div>
        <img src="/vychoziKlavesy.png" alt="Výchozí klávesy">
    </div>
    <div class="bloky">
        <img src="/klavesniceSPavoukem.png" alt="Pavouk s klávesnicí">
        <div class="text">
            <h2>Co zbytek klávesnice?</h2>
            <p>
                Každý prst si hlídá svůj pomyslný <span class="tlusty">sloupeček</span> označený jednou barvou.
                Například zelené klávesy mačkáme v ideálním případě pouze pravým ukazováčem. <br>
                <span class="tlusty">Výjimkou jsou palce</span>, které se starají pouze o mezerník.
            </p>
        </div>
    </div>
    <div class="bloky">
        <div class="text" style="display: flex; flex-direction: column; text-align: center;">
            <h2>Státní zkouška</h2>
            <ol>
                <li>Základní (<span class="tlusty">200+</span> CPM*)</li>
                <li>Se zvýšenou rychlostí (<span class="tlusty">300+</span> CPM*)</li>
                <li>Mistrovská (<span class="tlusty">400+</span> CPM*)</li>
            </ol>
            <hr>
            <p>
                Přesnost min. <span class="tlusty">99.5 %</span>
                přičemž za chybu se odečítá <span class="tlusty">10 úhozů</span>.
            </p>
            <div id="spodni">
                <a href="https://www.npi.cz/statni-tesnopisny-ustav-statni-zkousky" target=”_blank”>VÍCE ZDE</a>
                <span>*(CPM = úhozů za minutu)</span>
            </div>
        </div>
        <a href="https://www.npi.cz/statni-tesnopisny-ustav-statni-zkousky" target=”_blank”
            style="height: 100%; padding: 30px;">
            <img src="/npi.jpg" alt="Národní pedagogický institut"
                style="border-radius: 1000px; width: 100%; height: 100%;">
        </a>
    </div>
</template>

<style scoped>
h1 {
    margin-bottom: 0;
}

ol {
    text-align: center;
    margin-left: 20px;
}

#spodni {
    align-self: flex-end;
    margin-top: 20px;
    display: flex;
    justify-content: space-between;
    width: 100%;
}

#spodni a {
    color: var(--bila);
    width: auto;
}

@media screen and (max-width: 1000px) {
    h1 {
        margin-bottom: 10px;
    }

    ol {
        margin-left: 0;
    }
}
</style>