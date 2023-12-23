<script setup lang="ts">
import { useHead } from 'unhead'

useHead({
    title: "Jak psát",
    meta: [
        {
            name: "description",
            content: "Teorie psaní všemi deseti na české klávesnici. Popsaný prstoklad i tipy na zlepšení přesnosti.",
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
                <b>písmena F, J</b>, kde často bývá fyzická značka, abyste se jich mohli držet.
                Ostatní prsty umístíme <b>do stran ve stejné řadě</b>.
                Tedy levý prostředníček bude na D, prsteníček na S a malíček na A. Zrcadlově na druhé ruce.
            </p>
        </div>
        <img src="../assets/vychoziKlavesy.webp" alt="Výchozí klávesy">
    </div>
    <div class="bloky">
        <img src="../assets/klavesniceSPavoukemMensi.webp" alt="Pavouk s klávesnicí">
        <div class="text">
            <h2>Co zbytek klávesnice?</h2>
            <p>
                Každý prst si hlídá svůj pomyslný <b>sloupeček</b> označený jednou barvou.
                Například zelené klávesy mačkáme v ideálním případě pouze pravým ukazováčem. <br>
                <b>Výjimkou jsou palce</b>, které se starají pouze o mezerník.
            </p>
        </div>
    </div>
    <div class="bloky">
        <div class="text" style="display: flex; flex-direction: column; text-align: center;">
            <h2>Státní zkouška</h2>
            <ol>
                <li>Základní (<b>200+</b> CPM*)</li>
                <li>Se zvýšenou rychlostí (<b>300+</b> CPM*)</li>
                <li>Mistrovská (<b>400+</b> CPM*)</li>
            </ol>
            <hr>
            <p>
                Přesnost min. <b>99.5 %</b>
                přičemž za chybu se odečítá <b>10 úhozů</b>.
            </p>
            <div id="spodni">
                <a href="https://www.npi.cz/statni-tesnopisny-ustav-statni-zkousky" target=”_blank” rel="nofollow">VÍCE
                    ZDE</a>
                <span>*(CPM = úhozů za minutu)</span>
            </div>
        </div>
        <a id="npi" href="https://www.npi.cz/statni-tesnopisny-ustav-statni-zkousky" target=”_blank” rel="nofollow">
            <img src="../assets/npi.webp" alt="Národní pedagogický institut" width="220" height="220">
        </a>
    </div>
    <div class="bloky">
        <img src="../assets/sipky.svg" alt="Pavouk s klávesnicí">
        <div class="text">
            <h2 style="margin-bottom: 10px;">Tipy a triky</h2>
            <ul>
                <li><b>Ctrl + Backspace</b> smaže celé slovo, abys to nemusel vyťukávat po jednom</li>
                <li><u>Pohyb po Pavoukovi:</u>
                    <ul class="mensiPadding">
                        <li><b>Šipkama</b> vybereš lekci / cviční</li>
                        <li><b>Enter</b> zvolí vybranou lekci / cvičení</li>
                        <li><b>Mezerník</b> opakuje cvičení</li>
                        <li><b>Šipka doprava / Enter</b> pokračuje na další cvičení</li>
                        <li><b>Escape</b> se vrátí zpět</li>
                    </ul>
                </li>
            </ul>
        </div>
    </div>
</template>

<style scoped>
li {
    font-size: 1.1rem;
}

.mensiPadding li {
    padding: 2px 0 !important;
}

h1 {
    margin-bottom: 0;
}

ol,
ul {
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

#npi {
    height: 100%;
    padding: 30px;
}

#npi img {
    border-radius: 1000rem;
    aspect-ratio: 1 / 1;
    width: 100%;
}

@media screen and (max-width: 1100px) {
    h1 {
        margin-bottom: 10px;
    }

    ol {
        margin-left: 0;
    }

    img {
        max-width: 20rem !important;
    }

    #npi img {
        max-width: 12rem !important;
        height: 100%;
    }

    #npi {
        padding: 0;
    }
}
</style>