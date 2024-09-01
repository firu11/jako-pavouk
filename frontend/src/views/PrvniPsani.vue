<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useHead } from "@unhead/vue";
import Psani from "../components/Psani.vue";
import { pridatOznameni } from "../utils";
import Vysledek from "../components/Vysledek.vue";
import { mobil } from "../stores";

useHead({
    title: "První krůčky"
})

const text = ref([] as { id: number, znak: string, spatne: number, psat: boolean }[][]) // spatne: 0 ok, 1 spatne, 2 opraveno
const delkaTextu = ref(0)
const preklepy = ref(0)
const opravenePocet = ref(0)

const konec = ref(false)
const delkaNapsanehoTextu = ref(0)

function konecTextu(o: number, p: number, _: any, d: number) {
    opravenePocet.value = o
    preklepy.value = p
    konec.value = true
    delkaNapsanehoTextu.value = d
}

function setup() {
    if (mobil.value) {
        pridatOznameni("Psaní na telefonech zatím neučíme...")
        return
    }

    text.value = [] as { id: number, znak: string, spatne: number, psat: boolean }[][]
    let textRaw = "ffff jjjj ffjj jjff fjfj jfjf fjjj jfff jfjj fjff jjfj ffjf fjjf jffj "
    let slovoCounter = -1
    for (let j = 0; j < 35; j++) {
        for (let i = 0; i < textRaw.length; i++) {
            if (i == 0 || textRaw[i - 1] == " ") {
                text.value.push([])
                slovoCounter++
            }
            text.value[slovoCounter].push({ id: delkaTextu.value, znak: textRaw[i], spatne: 0, psat: true })
            delkaTextu.value++
        }
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

    <Psani v-if="!konec" @konec="konecTextu" @restart="restart" @pise="ok = true" :text :klavesnice="'qwertz'" :hide-klavesnice="!ok"
        :nacitamNovej="false" :cas="30" :delkaTextu />

    <Vysledek v-else :preklepy="preklepy" :opravenych="opravenePocet" :delkaTextu="delkaNapsanehoTextu" :cas="15" :cislo="'prvni-psani'" :posledni="true"
        @restart="restart" />

    <Transition>
        <div id="napoveda" v-if="!ok">
            <h3>Nápověda</h3>
            <ul>
                <li>Podtržení ukazuje na písmeno, které máš napsat.</li>
                <li>Pokud uděláš chybu, můžeš se vrátit zpět pomocí klávesy <span class="klavesa-v-textu">Backspace</span> a překlep opravit.</li>
                <li>Kdyby jsi chtěl začít odznova, zmáčkni klávesu <span class="klavesa-v-textu">Enter</span>.</li>
            </ul>
        </div>
    </Transition>
</template>

<style scoped>
#napoveda {
    background-color: var(--tmave-fialova);
    padding: 20px;
    border-radius: 10px;
    position: relative;
    top: -270px;
    max-width: 470px;
    display: flex;
    gap: 12px;
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
    gap: 10px;
    line-height: 1.5rem;
}


.v-leave-active {
    transition: 0.3s !important;
}

.v-leave-to {
    opacity: 0;
}
</style>