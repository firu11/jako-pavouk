<script setup lang="ts">
import { onMounted, ref } from 'vue';

const props = defineProps({
    zprava: String,
    sirka: Number,
    vzdalenost: {
        type: Number,
        default: 15
    }
})

const obsah = ref({} as HTMLElement)
const y = ref(0)
const sirka = ref(props.sirka)

onMounted(() => {
    y.value = obsah.value.getBoundingClientRect().bottom + props.vzdalenost
    if (props.sirka == null) {
        sirka.value = obsah.value.getBoundingClientRect().width
    } else sirka.value = props.sirka
    console.log(sirka.value, y.value)
})
</script>

<template>
    <div style="display: flex; flex-direction: column; align-items: center;">
        <div id="obsah" ref="obsah">
            <slot />
        </div>
        <div id="tooltip" :style="{ top: `${y}px`, maxWidth: `calc(${sirka}px + 10%)` }">
            {{ zprava }}
        </div>
    </div>
</template>

<style scoped>
#tooltip {
    visibility: hidden;
    background-color: black;
    color: white;
    text-align: center;
    padding: 5px;
    border-radius: 6px;
    font-size: 0.9rem;

    position: absolute;
    z-index: 1;
}

#obsah:hover~#tooltip {
    visibility: visible;
}

#obsah {
    cursor: help;
}
</style>