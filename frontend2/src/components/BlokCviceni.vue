<script setup lang="ts">
import { prihlasen } from '../stores';

defineProps({
    dokonceno: Boolean,
    index: Number,
    pismena: String,
    typ: String
})

function prihlaste_se() {
    alert("Nejprve se prosím přihlašte.")
}

</script>

<template>
    <router-link v-if="prihlasen" class="cvicBlok" :class="{ dokoncenyBlok: dokonceno }" :to="'/lekce/' + pismena + '/' + index">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova'">Se slovy</h3>
        <h3 v-else>...</h3>
        <img class="fajvkaVetsi" v-if="dokonceno" src="../assets/icony/right.svg" alt="Dokonceno!">
        <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
    </router-link>
    <a v-else class="cvicBlok" @click="prihlaste_se">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova'">Se slovy</h3>
        <h3 v-else>...</h3>
        <img class="playVetsi" src="../assets/icony/start.svg" alt="Začít lekci">
    </a>
</template>

<style scoped>
.cvicBlok {
    color: var(--bila);
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 200px;
    background-color: var(--tmave-fialova);
    height: 240px;
    transition-duration: 0.2s;
    padding: 15px 15px 30px 15px;
}

.cvicBlok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.cvicBlok hr {
    width: 160px;
    align-self: center;
    margin: 5px;
    color: gray;
}

.dokoncenyBlok {
    opacity: 50%;
}

.cvicBlok h3 {
    align-self: center;
    font-size: 24px;
    height: 100px;
}

.cvicBlok a {
    text-decoration: none;
    color: var(--bila);
    cursor: pointer;
}

h2 {
    font-size: 2em;
    font-weight: bolder;
}
</style>