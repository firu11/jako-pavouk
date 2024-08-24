<script setup lang="ts">
import { useHead } from "unhead";
import { Oznacene, checkTeapot, getToken, napovedaKNavigaci, pridatOznameni, naJednoDesetiny } from "../utils";
import axios from "axios";
import { onMounted, onUnmounted, ref } from "vue";
import { mobil } from "../stores";
import ObtiznostBar from "../components/ObtiznostBar.vue";
import Tooltip from "../components/Tooltip.vue";

useHead({
    title: "Procvičování",
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/procvic"
        }
    ]
})

const texty = ref(new Map<string, { id: number, jmeno: string, cpm: number, cislo: number, obtiznost: number }[]>())
const testPsaniCPM = ref(-1)
const o = new Oznacene()
let randomCvic = 1

onMounted(() => {
    axios.get("/procvic", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        var i = 2
        response.data.texty = response.data.texty.sort((a: { id: number, jmeno: string, kategorie: string, obtiznost: number }, b: { id: number, jmeno: string, kategorie: string, obtiznost: number }) => a.kategorie.localeCompare(b.kategorie))
        response.data.texty.forEach((txt: { id: number, jmeno: string, cpm: number, kategorie: string, obtiznost: number }) => {
            let a = texty.value.get(txt.kategorie)

            if (a) {
                a.push({ id: txt.id, jmeno: txt.jmeno, cpm: txt.cpm, cislo: i, obtiznost: txt.obtiznost })
            } else {
                texty.value.set(txt.kategorie, [{ id: txt.id, jmeno: txt.jmeno, cpm: txt.cpm, cislo: i, obtiznost: txt.obtiznost }])
            }
            i++
        })
        testPsaniCPM.value = response.data.testPsaniCPM
        o.setMax(response.data.texty.length + 1)
        randomCvic = Math.floor(Math.random() * response.data.texty.length) + 1
    }).catch(e => {
        if (!checkTeapot(e)) {
            pridatOznameni()
        }
    })
    document.addEventListener("keydown", e1)
    document.addEventListener("keyup", e2)
    document.addEventListener("mousemove", zrusitVyber)
})

function e1(e: KeyboardEvent) {
    if (e.key == "ArrowUp" || e.key == "ArrowLeft") {
        e.preventDefault()
        o.mensi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 500 })
    } else if (e.key == "ArrowDown" || e.key == "ArrowRight") {
        e.preventDefault()
        o.vetsi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 200 })
    } else if (e.key == "Enter") {
        e.preventDefault()
        let cvicE: HTMLElement | null = document.querySelector(`[i="true"]`)
        if (cvicE == null || o.bezOznaceni) {
            o.bezOznaceni = true
            o.index.value = randomCvic
        } else cvicE?.click()
    } else if (e.key == "Tab") {
        e.preventDefault()
        napovedaKNavigaci()
    }
}

function e2(e: KeyboardEvent) {
    if (e.key == "Enter") {
        e.preventDefault()
        let cvicE: HTMLElement | null = document.querySelector(`[i="true"]`)
        cvicE?.click()
    }
}

function zrusitVyber() {
    o.index.value = 0
}

function mobilKlik(e: MouseEvent) {
    e.preventDefault()
    pridatOznameni("Psaní na telefonech zatím neučíme...")
}

onUnmounted(() => {
    document.removeEventListener("keydown", e1)
    document.removeEventListener("keyup", e2)
    document.removeEventListener("mousemove", zrusitVyber)
})

</script>

<template>
    <h1>Procvičování</h1>

    <div id="seznam">
        <RouterLink v-if="!mobil" :to="'/test-psani'" class="blok" :i="1 == o.index.value"
            :class="{ oznacene: 1 == o.index.value, nohover: o.index.value != 0 }" style="margin-top: 5px;">
            <h3>Test psaní</h3>
            <span v-if="texty.size != 0 && testPsaniCPM != -1"><b>{{ naJednoDesetiny(testPsaniCPM) }}</b> CPM</span>
        </RouterLink>
        <a v-else href="/test-psani" class="blok" :i="1 == o.index.value" style="user-select: none; margin-top: 5px;" @click="mobilKlik">
            <h3>Test psaní</h3>
            <span v-if="texty.size != 0 && testPsaniCPM != -1"><b>{{ naJednoDesetiny(testPsaniCPM) }}</b> CPM</span>
        </a>

        <div v-if="texty.size != 0" v-for="k in texty.keys()" style="width: 100%;">
            <h2>{{ k }}</h2>

            <RouterLink v-if="!mobil" v-for="t in texty.get(k)" :to="`/procvic/${t.id}`" class="blok" :i="t.cislo == o.index.value"
                :class="{ oznacene: t.cislo == o.index.value, nohover: o.index.value != 0 }">
                <h3>
                    <Tooltip :sirka="100" :zprava="`${t.obtiznost == 1 ? 'Jednoduchá' : (t.obtiznost == 2 ? 'Střední' : 'Těžká')} obtížnost`" :xOffset="-38" :vzdalenost="5">
                        <ObtiznostBar :o="t.obtiznost" />
                    </Tooltip>
                    {{ t.jmeno }}
                </h3>
                <span v-if="t.cpm != -1"><b>{{ naJednoDesetiny(t.cpm) }}</b> CPM</span>
            </RouterLink>
            <div v-else v-for="t in texty.get(k)" class="blok" @click="mobilKlik">
                <h3>
                    <ObtiznostBar :o="t.obtiznost" />
                    {{ t.jmeno }}
                </h3>
            </div>
        </div>
        <div v-else>
            <h2>Knihy</h2>
            <div v-if="texty.size == 0" v-for="_ in 3" class="blok">
                <h3 style="margin-left: 8px;">. . .</h3>
            </div>

            <h2>Naučné</h2>
            <div v-if="texty.size == 0" v-for="_ in 2" class="blok">
                <h3 style="margin-left: 8px;">. . .</h3>
            </div>

            <h2>Zábavné</h2>
            <div v-if="texty.size == 0" v-for="_ in 2" class="blok">
                <h3 style="margin-left: 8px;">. . .</h3>
            </div>
        </div>
    </div>
</template>

<style scoped>
#seznam {
    display: flex;
    gap: 20px;
    text-align: left;
    flex-direction: column;
}

#seznam>div {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

.blok {
    display: flex;
    color: var(--bila);
    padding: 12px 20px 12px calc(25px - 8px);
    text-decoration: none;
    border-radius: 10px;
    width: 500px;
    background-color: var(--tmave-fialova);
    min-height: 64px;
    transition-duration: 0.2s;
    justify-content: space-between;

    cursor: pointer;
    user-select: none;
    /* kvuli tomu neprihlasenymu */
}

.blok:first-child {
    padding: 12px 20px 12px 25px;
}

.blok:hover,
.oznacene {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.blok h3 {
    font-size: 24px;
    font-weight: 300;
    margin: 0;
    align-self: center;
    display: flex;
    align-items: center;
    line-height: 40px;
}

.blok span {
    font-size: 1.2rem;
    display: flex;
    align-items: baseline;
    gap: 5px;
    justify-content: end;
    height: 34px;
    align-self: center;
    position: relative;
    top: -2px;
}

.blok span b {
    font-family: 'Red Hat Mono';
    font-size: 1.8rem;
}

@media screen and (max-width: 1100px) {
    #seznam {
        width: 70vw;
        align-items: center;
    }

    #seznam div {
        align-items: center;
    }

    h2 {
        align-self: start;
    }

    .blok {
        min-width: 254px;
        width: 100%;
        height: 48px;
        min-height: auto;
        padding: 12px 10px 12px 10px;

        /* kvuli tomu neprihlasenymu */
        cursor: pointer;
    }

    .blok span b {
        font-size: 1.35rem;
    }

    .blok span {
        font-size: 0.8rem;
        top: -2.5px;
        gap: 3px;
        height: 22px;
    }

    .blok h3 {
        font-size: 1.3rem;
    }


}
</style>