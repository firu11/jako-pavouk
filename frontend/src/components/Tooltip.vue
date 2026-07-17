<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue';

interface Props {
    zprava?: string;
    sirka: number;
    xOffset?: number;
    yOffset?: number;
    vzdalenost?: number;
    vzdalenostX?: number;
    hoverDelay?: string;
}
const props = withDefaults(defineProps<Props>(), {
    zprava: '',
    xOffset: 0,
    yOffset: 0,
    vzdalenost: 15,
    vzdalenostX: 0,
    hoverDelay: '0.4s',
});

const wrapper = useTemplateRef<HTMLElement>('wrapper');
const obsah = useTemplateRef<HTMLElement>('obsah');
const x = ref(0);
const y = ref(props.vzdalenost);
let delayedRecalculation: ReturnType<typeof setTimeout> | undefined;
let resizeObserver: ResizeObserver | undefined;

function recalculate() {
    if (wrapper.value == null || obsah.value == null) return;

    const wrapperRect = wrapper.value.getBoundingClientRect();
    const contentRect = obsah.value.getBoundingClientRect();
    const maxLeft = Math.max(12, document.documentElement.clientWidth - props.sirka - 12);
    const viewportLeft = Math.min(Math.max(contentRect.left + contentRect.width / 2 - props.sirka / 2 + props.vzdalenostX, 12), maxLeft);

    x.value = viewportLeft - wrapperRect.left;
    y.value = contentRect.bottom - wrapperRect.top + props.vzdalenost;
}

onMounted(() => {
    recalculate();
    delayedRecalculation = setTimeout(recalculate, 100);
    window.addEventListener('resize', recalculate);

    if (typeof ResizeObserver !== 'undefined' && obsah.value) {
        resizeObserver = new ResizeObserver(recalculate);
        resizeObserver.observe(obsah.value);
    }
});

onUnmounted(() => {
    if (delayedRecalculation !== undefined) clearTimeout(delayedRecalculation);
    resizeObserver?.disconnect();
    window.removeEventListener('resize', recalculate);
});

watch(
    () => [props.zprava, props.sirka, props.vzdalenost, props.vzdalenostX, props.xOffset, props.yOffset],
    () => void nextTick(recalculate),
);
</script>

<template>
    <div ref="wrapper" class="tooltip-wrapper">
        <div ref="obsah" class="tooltip-trigger" :style="{ top: `${props.yOffset}px`, left: `${props.xOffset}px` }">
            <slot />
        </div>
        <div class="tooltip" :style="{ top: `${y}px`, left: `${x}px`, width: `${props.sirka}px` }">
            <slot name="content">
                <span v-html="props.zprava"></span>
            </slot>
        </div>
    </div>
</template>

<style scoped>
.tooltip {
    opacity: 0%;
    background-color: black;
    color: white;
    text-align: center;
    padding: 5px;
    border-radius: 6px;
    font-size: 15px;

    position: absolute;
    z-index: 100;
    line-height: 16px;
    pointer-events: none;
    transition: 0.1s opacity;
}

.tooltip-trigger:hover ~ .tooltip {
    opacity: 100%;
    transition-delay: v-bind('props.hoverDelay');
}

.tooltip-trigger {
    position: relative;
    cursor: help;
}

.tooltip-wrapper {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
}
</style>
