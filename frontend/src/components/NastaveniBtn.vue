<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

const rotacePocet = ref(0)
function animace() {
    rotacePocet.value += 1
}

const rotace = computed(() => {
    return `rotate(${rotacePocet.value * 45}deg)`
})

const ready = ref(false)

onMounted(() => {
    ready.value = true
})

</script>
<template>
    <Transition>
        <div v-if="ready" id="nastaveniBtn" @click="$emit('klik'); animace()">
            <img :style="{ transform: rotace }" src="../assets/icony/nastaveni.svg" alt="Nastavení">
        </div>
    </Transition>
</template>
<style scoped>
.v-enter-active {
    transition: 0.3s !important;
}

.v-enter-from {
    opacity: 0;
}

#nastaveniBtn {
    position: relative;
    width: 55px;
    height: 55px;
    background-color: var(--tmave-fialova);
    border-radius: 100px;
    display: flex;
    align-items: center;
    justify-content: center;
    left: 385px;
    top: -140px;
    cursor: pointer;
    transition: background-color 0.1s;
    user-select: none;
}

#nastaveniBtn img {
    width: 30px;
    transition-duration: 0.4s;
}

#nastaveniBtn:hover {
    background-color: var(--fialova);
}
</style>