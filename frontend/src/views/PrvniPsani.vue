<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useHead } from "@unhead/vue";
import Psani from "../components/Psani.vue";
import { useRouter } from "vue-router";
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

const konec = ref(false)
const router = useRouter()

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

function konecTextu(c: number, o: number, p: number) {
    cas.value = c
    opravenePocet.value = o
    preklepy.value = p
    konec.value = true
}

onMounted(() => {
    if (mobil.value) {
        router.push("/registrace")
        pridatOznameni("Psaní na telefonech zatím neučíme. Registrovat se ale můžeš.")
        return
    }

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
})

</script>

<template>
    <h1 style="margin: 0">První krůčky</h1>

    <Psani v-if="!konec" @konec="konecTextu" :text="text" :delkaTextu="delkaTextu" :klavesnice="'qwertz'"
        :hideKlavesnice="false" />

    <Vysledek v-else :preklepy="preklepy" :opravenych="opravenePocet" :delkaTextu="delkaTextu"
        :casF="casFormat" :cas="cas" :cislo="'prvni-psani'" :posledni="true" />
</template>

<style scoped></style>