<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useHead } from "@unhead/vue";
import Psani from "../components/Psani.vue";
import { pridatOznameni } from "../utils";
import Vysledek from "../components/Vysledek.vue";
import { mobil } from "../stores"

useHead({
    title: "První krůčky"
})

const text = ref([[]] as { id: number, znak: string, spatne: number, }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)
const cas = ref(0)

const psaniRef = ref()

const konec = ref(false)

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

function konecTextu(c: number, o: number, p: number) {
    cas.value = c
    opravenePocet.value = o
    preklepy.value = p
    konec.value = true
}

function setup() {
    if (mobil.value) {
        pridatOznameni("Psaní na telefonech zatím neučíme...")
        return
    }

    text.value = [[]] as { id: number, znak: string, spatne: number, }[][]
    let textRaw = "ffff jjjj ffjj jjff fjfj jfjf fjjj jfff jfjj fjff jjfj ffjf fjjf jffj"
    let slovoCounter = -1
    for (let i = 0; i < textRaw.length; i++) {
        if (i == 0 || textRaw[i - 1] == " ") {
            text.value.push([])
            slovoCounter++
        }
        text.value[slovoCounter].push({ id: delkaTextu.value, znak: textRaw[i], spatne: 0 })
        delkaTextu.value++
    }
}

onMounted(() => {
    setup()
})

function restart() {
    konec.value = false
    delkaTextu.value = 0

    setup()
}

const ok = ref(false)

</script>

<template>
    <h1 style="margin: 0">První krůčky</h1>

    <Psani v-if="!konec" @konec="konecTextu" @pise="ok = true" :text="text" :delkaTextu="delkaTextu" :klavesnice="'qwertz'" :hideKlavesnice="!ok"
        :nacitam-novej="false" ref="psaniRef" />

    <Vysledek v-else :preklepy="preklepy" :opravenych="opravenePocet" :delkaTextu="delkaTextu" :casF="casFormat" :cas="cas" :cislo="'prvni-psani'"
        :posledni="true" @restart="restart" />
    <Transition>
        <div id="napoveda" v-if="!ok">
            <h3>Nápověda</h3>
            <ul>
                <li>Podtržení ukazuje na písmeno, které máš napsat.</li>
                <li>Pokud uděláš chybu, můžeš se vrátit zpět pomocí klávesy <span class="klavesaVTextu">Backspace</span> a překlep opravit.</li>
                <li>Kdyby jsi chtěl začít odznova, zmáčkni klávesu <span class="klavesaVTextu">Enter</span>.</li>
            </ul>
            <button class="tlacitko" @click="ok = true">Jdu na to!</button>
        </div>
    </Transition>
</template>

<style scoped>
#napoveda {
    background-color: var(--tmave-fialova);
    padding: 20px;
    border-radius: 10px;
    position: relative;
    top: -295px;
    max-width: 470px;
    display: flex;
    gap: 10px;
    align-items: center;
    flex-direction: column;
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
}

#napoveda h3 {
    font-size: 1.4rem;
    font-weight: 500;
}

#napoveda ul {
    padding-left: 10px;
    list-style-type: none;
    display: flex;
    flex-direction: column;
    gap: 7px;
    line-height: 1.4rem;
}


.v-leave-active {
    transition: 0.3s !important;
}

.v-leave-to {
    opacity: 0;
}
</style>