<script setup lang="ts">
import { useHead } from 'unhead'

useHead({
    title: "Jak psát všemi deseti",
    meta: [
        {
            name: "description",
            content: "Teorie jak psát všemi deseti na české klávesnici. Popsaný prstoklad i tipy na zlepšení přesnosti.",
        },
        {
            name: "keywords",
            content: "psani, psaní, všemi, vsemi, všema, vsema, deseti, teorie, jak, psát, klávesy, princip, teorie, prstoklad",
        }
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/jak-psat"
        }
    ]
})
</script>

<template>
    <h1>Jak psát všemi deseti</h1>

    <div class="bloky" id="blok">
        <img src="../assets/vychoziKlavesy.webp" alt="Výchozí klávesy" width="732" height="254" />
        <div class="text">
            <h2>Výchozí pozice</h2>
            <p>
                Základem je mít prsty ve správné výchozí pozici, která je velice prostá. Ukazováčky položíme na
                <b>písmena F, J</b>, kde často bývá fyzická značka, abychom se jich mohli držet.
                Ostatní prsty umístíme <b>do stran ve stejné řadě</b>.
                Tedy levý prostředníček bude na D, prsteníček na S a malíček na A. Zrcadlově na druhé ruce.
            </p>
        </div>
    </div>
    <div class="bloky">
        <img src="../assets/klavesniceSPavoukem.webp" alt="Pavouk s klávesnicí" width="300" height="233">
        <div class="text">
            <h2>Co zbytek klávesnice?</h2>
            <p>
                Každý prst si hlídá svůj pomyslný <b>sloupeček</b> označený jednou barvou.
                Například zelené klávesy mačkáme v ideálním případě pouze pravým ukazováčem.
                <b>Výjimkou jsou palce</b>, které se starají pouze o&nbsp;mezerník.
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
                Přesnost min. <b>99.5 %</b>, přičemž za neopravenou chybu se odečítá <b>10 úhozů</b>.
            </p>
            <div id="spodni">
                <a href="https://www.npi.cz/statni-tesnopisny-ustav-statni-zkousky" target=”_blank” rel="nofollow">VÍCE
                    ZDE</a>
                <span>*(CPM = úhozů za minutu)</span>
            </div>
        </div>
        <a id="npi" href="https://www.npi.cz/statni-tesnopisny-ustav-statni-zkousky" target=”_blank” rel="nofollow">
            <img src="../assets/npi.webp" alt="Národní pedagogický institut" width="225" height="225">
        </a>
    </div>
    <div class="bloky" id="pohyb">
        <img src="../assets/sipky.svg" alt="Pavouk s klávesnicí" width="345" height="244">
        <div class="text">
            <h2>Pohyb po Pavoukovi</h2>
            <ul class="mensiPadding" style="list-style: none;">
                <li><b>Výběr lekce / cvičení:</b></li>
                <li>
                    <ul style="margin: 0 0 1.1rem 0.5rem; list-style: none;">
                        <li><span class="klavesaVTextu">Enter</span> podržet - vybere to, kde jsem skončil</li>
                        <li><span class="klavesaVTextu" style="margin-right: 5px;"><img src="../assets/icony/sipkaL.svg" alt="Šipka"
                                    class="klavSipka"></span>
                            <span class="klavesaVTextu"><img src="../assets/icony/sipkaL.svg" alt="Šipka" class="klavSipka"
                                    style="transform: scaleX(-1);"></span> - přesný výběr,
                            <span class="klavesaVTextu" style="margin-left: 5px;">Enter</span>
                            - potvrzení
                        </li>
                    </ul>
                </li>

                <li><span class="klavesaVTextu"><img src="../assets/icony/sipkaL.svg" alt="Šipka" class="klavSipka"
                            style="transform: scaleX(-1);"></span> - pokračuje na další cvičení
                </li>
                <li><span class="klavesaVTextu">Enter</span> - opakuje cvičení</li>
                <li><span class="klavesaVTextu">Esc</span> - zpět</li>
            </ul>
        </div>
    </div>
</template>

<style scoped>
.bloky {
    height: 350px;
}

.bloky img {
    width: 300px;
}

#pohyb .text {
    width: 55%;
}

#pohyb p {
    text-align: start;
}

#pohyb>img {
    transform: scale(1.15);
}

#blok {
    flex-direction: column;
    height: auto;
    gap: 0px;
    padding-top: 32px;
}

#blok .text {
    width: auto;
}

#blok .text h2 {
    margin-bottom: 16px;
    margin-top: 24px;
}

#blok img {
    width: 100%;
    margin-bottom: 16px;
}

li {
    font-size: 1.1rem;
}

.mensiPadding li {
    padding: 3px 0 !important;
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
    width: 225px;
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

    #blok img {
        width: 100%;
        max-width: 500px !important;
    }

    #blok {
        padding: 16px 15px 25px 15px;
    }

    .bloky {
        height: auto;
    }

    .bloky img {
        width: 100%;
    }

    #pohyb .text {
        width: auto;
    }
}
</style>