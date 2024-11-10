<script setup lang="ts">
import { onMounted, ref, watch } from "vue";

const props = defineProps({
    cislo: {
        type: Number,
        required: true
    },
    desetinaMista: {
        type: Number,
        default: 1
    },
})

const zobrazeneCislo = ref("0")
const jedenFrame = 1000 / 60 // fps
const dobaTrvani = 1400
let puvodniCislo = 0

let counter = 0

const celkemFramu = Math.ceil(dobaTrvani / jedenFrame)

onMounted(() => {
    animace()
    setTimeout(() => {
        watch(props, () => {
            clearInterval(counter)
            animace()
        })
    }, 300)
})

function animace() {
    let frame = 0

    if (zobrazeneCislo.value != "") puvodniCislo = parseInt(zobrazeneCislo.value)

    counter = setInterval(() => {
        frame++

        let t = frame / celkemFramu
        zobrazeneCislo.value = transform(Math.sqrt(1 - Math.pow(t - 1, 6)), puvodniCislo, props.cislo).toFixed(props.desetinaMista)

        if (frame === celkemFramu) {
            clearInterval(counter)
            puvodniCislo = props.cislo
        }
    }, jedenFrame)
}

function transform(x: number, a: number, b: number) {
    if (a == b) return a
    else if (a == 0) return b * x
    else if (b == 0) return 0
    else return a + (b - a) * x
}

</script>
<template>
    <span>{{ zobrazeneCislo }}</span>
</template>
<style scoped></style>