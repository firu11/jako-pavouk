<script setup lang="ts">
import { useHead } from "unhead";
import { ref } from "vue";
import { pridatOznameni } from "../utils";

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
        console.error("Async: Could not copy text: ", err);
    })
    pridatOznameni("Zkopírováno!", undefined, "copy")
}

</script>

<template>
    <h1>O nás</h1>
    <p>
        Jako Pavouk je studentský projekt zaměřený na výuku psaní všemi deseti.
        <span style="margin: -14px;" /> <!-- trochu menší mezera než 2x <br> -->
        Když jsem se rozhodl, že bych se chtěl naučit psát rychleji, začal jsem hledat vhodné programy, které by mi s tím pomohly.
        <br>
        Zjistil jsem ale, že zejména pro českou klávesnici jich není mnoho a některé z nich dokonce stojí téměř 1000 Kč.
        A tak jsem se rozhodl vytvořit si vlastní.
        <span style="margin: -14px;" /> <!-- trochu menší mezera než 2x <br> -->
        A ano, čteš správně!
        Text je psán v první osobě. Možná by se tato stránka měla jmenovat spíše "O mně", jelikož tento web spravuje pouze jeden člověk.
    </p>
    <p>
        Chtěl bych poděkovat těmto webům, které mi poskytly svůj obsah, ať už ve formě vzdělávacích textů, článků, nebo pohádek.
        Většina textů je upravena.
        <span id="pytel-partneru">
            <a href="https://www.dejepis.com/" target="_blank">dejepis.com</a>
            <!-- <a href="https://dobryzpravy.cz/" target="_blank">dobryzpravy.cz</a> -->
            <a href="https://www.pohadkozem.cz/" target="_blank">pohadkozem.cz</a>
            <a href="https://www.pohadky.org/" target="_blank">pohadky.org</a>
            <a href="https://www.zemepis24.cz/" target="_blank">zemepis24.cz</a>
        </span>
    </p>
    <p>
        Jelikož je Jako Pavouk kompletně zdarma, je tu možnost tento projekt podpořit přes platformu Ko-fi níže.
        Moc si vážím vaší podpory!
        <a href="https://ko-fi.com/jakopavouk" target="_blank" rel="nofollow">
            <img id="kofi" src="../assets/kofi_button.webp" alt="Podpořit" width="230" height="36">
        </a>
    </p>
    <p>
        Objevil/a jsi nějakou chybu/bug?
        <br>
        Máš nějaký zajímavý nápad jak vylepšit tuto aplikaci nebo jen chceš dát zpětnou vazbu? Neboj se mi napsat na e-mail:
        <span id="btn" @click="copyEmail">
            <span style="unicode-bidi:bidi-override; direction: rtl;">zc</span>
            <span style="display: none;">nechci dostavat spam</span>
            <span>&#46;</span>
            <span style="display: none;">proste to nedelej</span>
            <span style="unicode-bidi:bidi-override; direction: rtl;">kuovapokaj</span>
            <span style="display: none;">nebudete mi websrapovat email vy smradlavý smrdi</span>
            <span>&#64;</span>
            <span style="display: none;">ani to nezkoušej</span>
            <span style="unicode-bidi:bidi-override; direction: rtl;">urif</span>

            <button ref="copy"><img src="../assets/icony/copy.svg" alt="Kopírovat"></button>
        </span>
    </p>
    <p>
        Pro zvědavé tu mám i link na GitHub, kde je celý projekt volně k prohlédnutí.
        <br>
        Budu moc rád za hvězdičku!
        <a href="https://github.com/Firu115/jako_pavouk" target=”_blank” rel="nofollow">
            <img id="github" src="../assets/github.webp" alt="GitHub" width="173" height="60">
        </a>
    </p>
</template>

<style scoped>
#pytel-partneru {
    display: flex;
    gap: 20px;
    flex-wrap: wrap;
    justify-content: center;
}

#pytel-partneru a {
    padding: 10px 16px;
    background-color: var(--fialova);
    border-radius: 8px;
    text-decoration: none;
    color: var(--bila);
    font-weight: 500;
    transition: 0.2s;
}

#pytel-partneru a:hover {
    background-color: var(--svetle-fialova);
}

#btn {
    margin: -5px;
    background-color: var(--fialova);
    padding: 10px 14px 6px 14px;
    border-radius: 8px;
    cursor: grab;
    transition: 0.2s;
}

#btn:hover {
    background-color: var(--svetle-fialova);
}

button {
    background-color: transparent;
    border: 0px;
    transition: 0.1s ease-in-out;
    cursor: grab;
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
    line-height: 1.5rem;
}

#github {
    width: 9em !important;
}

#kofi,
#github {
    width: 12em;
    height: auto;
    transition: 0.2s;
    cursor: pointer;
    text-align: center;
    color: var(--svetle-fialova);
}

#kofi:hover,
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