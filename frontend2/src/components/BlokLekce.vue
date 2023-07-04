<script setup lang="ts">
import {prihlasen} from '../stores'
import { formatovany_pismena } from '../utils';

defineProps({
    pismena: String,
    je_dokoncena: Boolean,
})

</script>

<template>
    <RouterLink class="lekceBlok" :class="{ hotovoBlok: je_dokoncena }" v-if="prihlasen" :to="'/lekce/' + pismena">
        <h2>Lekce: {{ formatovany_pismena(pismena) }}</h2>
        <img class="fajvka" v-if="je_dokoncena" src="../assets/icony/right.svg" alt="Dokonceno!">
    </RouterLink>
    <a v-else class="lekceBlok">
        <h2>Lekce: {{ formatovany_pismena(pismena) }}</h2>
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
    cursor: pointer; /* kvuli tomu neprihlasenymu */
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
}</style>