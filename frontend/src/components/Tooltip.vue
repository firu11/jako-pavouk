<script setup lang="ts">
import { onMounted, ref } from 'vue';

const props = defineProps({
    zprava: String,
    sirka: Number,
    xOffset: {
        type: Number,
        default: 0
    },
    vzdalenost: {
        type: Number,
        default: 15
    }
})

const obsah = ref({} as HTMLElement)
const tip = ref({} as HTMLElement)
const y = ref(0)

onMounted(() => {
    y.value = obsah.value.getBoundingClientRect().bottom + props.vzdalenost

    if (props.xOffset != 0) {
        let rect = tip.value.getBoundingClientRect()
        tip.value.style.left = `${props.xOffset + rect.left}px`
    }
})
</script>

<template>
    <div style="display: flex; flex-direction: column; align-items: center;">
        <div id="obsah" ref="obsah">
            <slot />
        </div>
        <div id="tooltip" :style="{ top: `${y}px`, maxWidth: `${props.sirka == null ? obsah.getBoundingClientRect().width * 2.2 : props.sirka}px` }"
            v-html="zprava" ref="tip" />
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
    line-height: 1rem;
    pointer-events: none;
}

#obsah:hover~#tooltip {
    visibility: visible;
}

#obsah {
    cursor: help;
}
</style>