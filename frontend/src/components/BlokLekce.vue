<script setup lang="ts">
import { prihlasen } from '../stores'
import { formatovanyPismena, format } from '../utils';

defineProps({
    pismena: {
        type: String,
        default: ""
    },
    jeDokoncena: Boolean,
    oznacena: Boolean
})

</script>

<template>
    <RouterLink v-if="pismena !== '...'" class="lekceBlok" :class="{ hotovoBlok: jeDokoncena, oznacene: oznacena }" :to="'/lekce/' + pismena">
        <h2>Lekce: <b>{{ format(pismena) }}</b></h2>
        <img class="fajvka" v-if="prihlasen && jeDokoncena" src="../assets/icony/right.svg" alt="Dokonceno!">
    </RouterLink>
    <a v-else class="lekceBlok"> <!-- aby na to neslo kliknout nez se to nacte -->
        <h2>Lekce: <b>{{ formatovanyPismena(pismena) }}</b></h2>
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
    transition-duration: 0.1s;
    cursor: pointer;
    /* kvuli tomu neprihlasenymu */
}

.lekceBlok:hover,
.oznacene {
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

@media screen and (max-width: 1100px) {
    .lekceBlok {
        min-width: 260px;
        width: 100%;
        background-color: var(--tmave-fialova);
        min-height: 48px;
        max-height: 100px;
        height: auto;
        padding: 10px 20px 10px 18px;

        /* kvuli tomu neprihlasenymu */
        cursor: pointer;
    }

    .lekceBlok h2 {
        font-size: 1.3rem;
    }

    .fajvka {
        height: 25px;
        margin-left: 10px;
    }
}
</style>