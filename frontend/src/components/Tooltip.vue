<script setup lang="ts">
import { onUnmounted, ref, useTemplateRef, watch } from 'vue';
import { onMounted } from 'vue';

interface Props {
    zprava: string;
    sirka: number;
    xOffset?: number;
    yOffset?: number;
    vzdalenost?: number;
    vzdalenostX?: number;
    hoverDelay?: string;
}
const props = withDefaults(defineProps<Props>(), {
    xOffset: 0,
    yOffset: 0,
    vzdalenost: 15,
    vzdalenostX: 0,
    hoverDelay: '0.4s',
});

const obsah = useTemplateRef('obsah');
const tip = useTemplateRef('tip');
const y = ref(props.vzdalenost + document.documentElement.scrollTop);

onMounted(() => {
    recalc();
    recalcTipY();
    setTimeout(recalcTipY, 100);
    window.addEventListener('resize', recalc);
});

onUnmounted(() => {
    window.removeEventListener('resize', recalc);
});

function getPageTopLeft(el: Element) {
    var rect = el.getBoundingClientRect();
    var docEl = document.documentElement;
    return {
        left: rect.left + (window.scrollX || docEl.scrollLeft || 0),
        top: rect.top + (window.scrollY || docEl.scrollTop || 0),
    };
}

function recalcTipY() {
    if (obsah.value == null) return;
    y.value = obsah.value.getBoundingClientRect().bottom + props.vzdalenost;
}

function recalc() {
    if (tip.value == null) return;
    tip.value.style.removeProperty('left');
    tip.value.style.removeProperty('right');

    let left = getPageTopLeft(tip.value).left + props.vzdalenostX + props.xOffset;
    if (left + props.sirka! > document.body.clientWidth) {
        tip.value.style.right = `12px`;
    } else {
        if (obsah.value == null || typeof obsah.value.getBoundingClientRect !== 'function') {
            tip.value.style.left = `${props.vzdalenostX + document.documentElement.scrollLeft}px`;
        } else if (props.xOffset == 0 && props.vzdalenostX == 0) {
            return; // nevim co to dělá no nic
        } else {
            tip.value.style.left = `${obsah.value.getBoundingClientRect().left + obsah.value.getBoundingClientRect().width / 2 - props.sirka! / 2 + props.vzdalenostX}px`;
        }
    }
}

watch(obsah, recalc);
</script>

<template>
    <div id="wrap">
        <div id="obsah" ref="obsah" :style="{ top: `${props.yOffset}px`, left: `${props.xOffset}px` }">
            <slot />
        </div>
        <div id="tooltip" :style="{ top: `${y}px`, width: `${props.sirka == null ? obsah!.getBoundingClientRect().width * 2.2 : props.sirka}px` }" v-html="zprava" ref="tip" />
    </div>
</template>

<style scoped>
#tooltip {
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

#obsah:hover ~ #tooltip {
    opacity: 100%;
    transition-delay: v-bind('props.hoverDelay');
}

#obsah {
    position: relative;
    cursor: help;
}

#wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
}
</style>
