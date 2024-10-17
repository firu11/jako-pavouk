<script setup lang="ts">
import { prihlasen } from "../stores";
import { format } from "../utils";

defineProps({
    pismena: {
        type: String,
        default: ""
    },
    jeDokoncena: Boolean,
    oznacena: Boolean,
    cislo: Number
})

</script>

<template>
    <RouterLink v-if="pismena !== '. . .'" class="lekce-blok" :class="{ hotovoBlok: jeDokoncena, oznacene: oznacena }" :to="'/lekce/' + pismena">
        <h2>Lekce {{ cislo }}: <b>{{ format(pismena) }}</b></h2>
        <img class="fajvka" v-if="prihlasen && jeDokoncena" src="../assets/icony/right.svg" alt="Dokonceno!">
    </RouterLink>
    <a v-else class="lekce-blok"> <!-- aby na to neslo kliknout nez se to nacte -->
        <h2>Lekce {{ cislo }}: <b>{{ pismena }}</b></h2>
        <img class="fajvka" v-if="prihlasen && jeDokoncena" src="../assets/icony/right.svg" alt="Dokonceno!">
    </a>
</template>

<style scoped>
.lekce-blok {
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

.lekce-blok:hover,
.oznacene {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.hotovoBlok {
    color: var(--seda);
    opacity: 80%;
}

.lekce-blok h2 {
    align-self: center;
    font-size: 24px;
    font-weight: 300;
}

@media screen and (max-width: 1100px) {
    .lekce-blok {
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

    .lekce-blok h2 {
        font-size: 21px;
    }

    .fajvka {
        height: 25px;
        margin-left: 10px;
    }
}
</style>