<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';

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

const zobrazeneCislo = ref("")
const jedenFrame = 1000 / 60
const dobaTrvani = 1400

const celkemFramu = Math.ceil(dobaTrvani / jedenFrame)

onMounted(() => {
    animace()
    setTimeout(() => {
        watch(props, () => {
            animace()
        })
    }, 300)
})

function animace() {
    let frame = 0

    const counter = setInterval(() => {
        frame++

        let t = frame / celkemFramu
        zobrazeneCislo.value = (Math.sqrt(1 - Math.pow(t - 1, 4)) * props.cislo).toFixed(props.desetineMista)

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