<script setup lang="ts">
import { onMounted, ref } from 'vue';

const props = defineProps({
    cislo: {
        type: Number,
        required: true
    },
    desetineMista: {
        type: Number,
        default: 1
    }
})

const zobrazeneCislo = ref(0)
const jedenFrame = 1000 / 60
const dobaTrvani = 800

const celkemFramu = Math.ceil(dobaTrvani / jedenFrame)

onMounted(() => {
    animace()
})

function animace() {
    let frame = 0

    const counter = setInterval(() => {
        frame++

        zobrazeneCislo.value = Math.round(props.cislo * ((frame / celkemFramu) === 1 ? 1 : 1 - Math.pow(2, -10 * (frame / celkemFramu))) * 10 ** props.desetineMista) / 10 ** props.desetineMista

        if (frame === celkemFramu) {
            clearInterval(counter)
        }
    }, jedenFrame)
}
</script>
<template>
    <span>{{ zobrazeneCislo }}</span>
</template>
<style scoped></style>