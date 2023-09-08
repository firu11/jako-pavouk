<script setup lang="ts">
import { prihlasen } from '../stores'
import { formatovanyPismena } from '../utils';

defineProps({
    pismena: {
        type: String,
        default: ""
    },
    jeDokoncena: Boolean,
})


</script>

<template>
    <RouterLink v-if="pismena !== '...'" class="lekceBlok" :class="{ hotovoBlok: jeDokoncena }"
        :to="'/lekce/' + encodeURIComponent(pismena)">
        <h2>Lekce: <span class="tlusty">{{ formatovanyPismena(pismena) }}</span></h2>
        <img class="fajvka" v-if="prihlasen && jeDokoncena" src="../assets/icony/right.svg" alt="Dokonceno!">
    </RouterLink>
    <a v-else class="lekceBlok"> <!-- aby na to ńeslo kliknout nez se to nacte -->
        <h2>Lekce: <span class="tlusty">{{ formatovanyPismena(pismena) }}</span></h2>
        <img class="fajvka" v-if="prihlasen && jeDokoncena" src="../assets/icony/right.svg" alt="Dokonceno!">
    </a>
</template>

<style scoped>
.lekceBlok {
    color: var(--bila);
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 12px 20px 12px 25px;
    text-decoration: none;
    border-radius: 10px;
    width: 500px;
    background-color: var(--tmave-fialova);
    height: 64px;
    transition-duration: 0.2s;
    cursor: pointer;
    /* kvuli tomu neprihlasenymu */
}

.lekceBlok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.hotovoBlok {
    color: var(--seda);
    opacity: 80%;
}

.lekceBlok h2 {
    align-self: center;
    font-size: 24px;
    font-weight: 300;
}

@media screen and (max-width: 1000px) {
    .lekceBlok {
        min-width: 270px;
        width: 100%;
        background-color: var(--tmave-fialova);
        height: 64px;
        transition-duration: 0.2s;

        /* kvuli tomu neprihlasenymu */
        cursor: pointer; 
    }
}
</style>