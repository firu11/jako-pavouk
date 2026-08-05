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

const obsah = useTemplateRef<HTMLElement>('obsah');
const tooltip = useTemplateRef<HTMLElement>('tooltip');
const x = ref(0);
const y = ref(props.vzdalenost);
const visible = ref(false);
let delayedRecalculation: ReturnType<typeof setTimeout> | undefined;
let resizeObserver: ResizeObserver | undefined;

function recalculate() {
    if (obsah.value == null || tooltip.value == null) return;

    const contentRect = obsah.value.getBoundingClientRect();
    const tooltipWidth = tooltip.value.getBoundingClientRect().width;
    const maxLeft = Math.max(12, document.documentElement.clientWidth - tooltipWidth - 12);

    x.value = Math.min(Math.max(contentRect.left + contentRect.width / 2 - tooltipWidth / 2 + props.vzdalenostX, 12), maxLeft);
    y.value = contentRect.bottom + props.vzdalenost;
}

function show() {
    recalculate();
    visible.value = true;
}

function hide() {
    visible.value = false;
}

onMounted(() => {
    recalculate();
    delayedRecalculation = setTimeout(recalculate, 100);
    window.addEventListener('resize', recalculate);
    window.addEventListener('scroll', recalculate, true);

    if (typeof ResizeObserver !== 'undefined' && obsah.value) {
        resizeObserver = new ResizeObserver(recalculate);
        resizeObserver.observe(obsah.value);
    }
});

onUnmounted(() => {
    if (delayedRecalculation !== undefined) clearTimeout(delayedRecalculation);
    resizeObserver?.disconnect();
    window.removeEventListener('resize', recalculate);
    window.removeEventListener('scroll', recalculate, true);
});

watch(
    () => [props.zprava, props.sirka, props.vzdalenost, props.vzdalenostX, props.xOffset, props.yOffset],
    () => void nextTick(recalculate),
);
</script>

<template>
    <div class="tooltip-wrapper" @mouseenter="show" @mouseleave="hide" @focusin="show" @focusout="hide">
        <div ref="obsah" class="tooltip-trigger" :style="{ top: `${props.yOffset}px`, left: `${props.xOffset}px` }">
            <slot />
        </div>
        <Teleport to="body">
            <div ref="tooltip" class="tooltip" :class="{ 'tooltip-visible': visible }" :style="{ top: `${y}px`, left: `${x}px`, width: `${props.sirka}px` }" role="tooltip">
                <slot name="content">
                    <span v-html="props.zprava"></span>
                </slot>
            </div>
        </Teleport>
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

    position: fixed;
    z-index: 100;
    line-height: 16px;
    pointer-events: none;
    transition: 0.1s opacity;
}

.tooltip-visible {
    opacity: 100%;
    transition-delay: v-bind('props.hoverDelay');
}

.tooltip-trigger {
    position: relative;
    cursor: help;
}

/* The repeated class keeps parent scoped selectors such as `.blok div` from overriding the component's layout. */
.tooltip-wrapper.tooltip-wrapper {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
}
</style>
