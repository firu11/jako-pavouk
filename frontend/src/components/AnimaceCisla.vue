<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue';

interface Props {
    cislo: number;
    desetinaMista?: number;
}
const props = withDefaults(defineProps<Props>(), {
    desetinaMista: 1,
});

const zobrazeneCislo = ref('0');
const dobaTrvani = 1400;
let animationFrame: number | null = null;

watch([() => props.cislo, () => props.desetinaMista], animace);

onMounted(animace);

onUnmounted(() => {
    zrusAnimaci();
});

function animace() {
    zrusAnimaci();

    const zobrazenaHodnota = Number.parseFloat(zobrazeneCislo.value);
    const puvodniCislo = Number.isNaN(zobrazenaHodnota) ? 0 : zobrazenaHodnota;
    const ciloveCislo = props.cislo;
    const desetinaMista = props.desetinaMista;
    let zacatek: number | null = null;

    const vykresliFrame = (cas: number) => {
        zacatek ??= cas;
        const t = Math.min((cas - zacatek) / dobaTrvani, 1);
        const prubeh = Math.sqrt(1 - Math.pow(t - 1, 6));

        zobrazeneCislo.value = transform(prubeh, puvodniCislo, ciloveCislo).toFixed(desetinaMista);

        if (t < 1) {
            animationFrame = requestAnimationFrame(vykresliFrame);
        } else {
            animationFrame = null;
        }
    };

    animationFrame = requestAnimationFrame(vykresliFrame);
}

function zrusAnimaci() {
    if (animationFrame == null) return;

    cancelAnimationFrame(animationFrame);
    animationFrame = null;
}

function transform(x: number, a: number, b: number) {
    if (a == b) return a;
    else if (a == 0) return b * x;
    else if (b == 0) return 0;
    else return a + (b - a) * x;
}
</script>
<template>
    <span>{{ zobrazeneCislo }}</span>
</template>
<style scoped></style>
