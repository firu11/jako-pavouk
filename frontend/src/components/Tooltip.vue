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
})
</script>

<template>
    <div style="display: flex; flex-direction: column; align-items: center; width: 100%;">
        <div id="obsah" ref="obsah">
            <slot />
        </div>
        <div id="tooltip" :style="{ top: `${y}px`, maxWidth: `${props.sirka == null ? obsah.getBoundingClientRect().width * 2.2 : props.sirka}px` }">
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
    z-index: 100;
}

#obsah:hover~#tooltip {
    visibility: visible;
}

#obsah {
    cursor: help;
}
</style>