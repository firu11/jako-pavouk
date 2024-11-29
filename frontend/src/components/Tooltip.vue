<script setup lang="ts">
import { computed, onUnmounted } from "vue";
import { onMounted, ref } from "vue";

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
    },
    vzdalenostX: {
        type: Number,
        default: 0
    }
})

const obsah = ref({} as HTMLElement)
const tip = ref({} as HTMLElement)

const y = computed(() => {
    if (typeof obsah.value.getBoundingClientRect !== 'function') return props.vzdalenost + document.documentElement.scrollTop
    return obsah.value.getBoundingClientRect().bottom + props.vzdalenost + document.documentElement.scrollTop
})

onMounted(() => {
    if (props.xOffset !== 0) {
        let rect = tip.value.getBoundingClientRect()
        tip.value.style.left = `${props.xOffset + rect.left}px`
        obsah.value.style.left = `${props.xOffset}px`
    }

    recalc()
    window.addEventListener('resize', recalc)
})

onUnmounted(() => {
    window.removeEventListener('resize', recalc)
})

function getPageTopLeft(el: Element) {
    var rect = el.getBoundingClientRect()
    var docEl = document.documentElement
    return {
        left: rect.left + (window.scrollX || docEl.scrollLeft || 0),
        top: rect.top + (window.scrollY || docEl.scrollTop || 0)
    }
}

function recalc() {
    if (props.vzdalenostX !== 0) {
        tip.value.style.removeProperty('left')
        tip.value.style.removeProperty('right')
        let left = getPageTopLeft(tip.value).left + props.vzdalenostX

        if (left + props.sirka! > document.body.clientWidth) {
            tip.value.style.right = `12px`
        } else {
            tip.value.style.left = `${left}px`
        }
    }
}

</script>

<template>
    <div id="wrap">
        <div id="obsah" ref="obsah">
            <slot />
        </div>
        <div id="tooltip" :style="{ top: `${y}px`, width: `${props.sirka == null ? obsah.getBoundingClientRect().width * 2.2 : props.sirka}px` }"
            v-html="zprava" ref="tip" />
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

#obsah:hover~#tooltip {
    opacity: 100%;
    transition-delay: 0.4s;
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