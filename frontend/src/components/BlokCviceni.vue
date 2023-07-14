<script setup lang="ts">
import { ref } from 'vue';
import { prihlasen } from '../stores';
import { pridatOznameni } from '../utils';

defineProps({
    dokonceno: Boolean,
    index: Number,
    pismena: {
        type: String,
        default: ""
    },
    typ: String
})

const mobil = ref(document.body.clientWidth <= 1000)

</script>

<template>
    <router-link v-if="prihlasen && typ !== '...' && !mobil" class="cvicBlok" :class="{ dokoncenyBlok: dokonceno }"
        :to="'/lekce/' + encodeURIComponent(pismena) + '/' + index">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova'">Se slovy</h3>
        <h3 v-else>...</h3>
        <img class="fajvkaVetsi" v-if="dokonceno" src="../assets/icony/right.svg" alt="Dokonceno!">
        <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
    </router-link>
    <a v-else-if="typ === '...' && !mobil" class="cvicBlok"> <!-- aby na to ńeslo kliknout nez se to nacte -->
        <h2>{{ index }}</h2>
        <hr>
        <h3>...</h3>
        <img class="playVetsi" src="../assets/icony/start.svg" alt="Začít lekci">
    </a>
    <a v-else-if="!mobil" class="cvicBlok" @click="pridatOznameni('Nejprve se prosím přihlašte')">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova'">Se slovy</h3>
        <h3 v-else>...</h3>
        <img class="playVetsi" src="../assets/icony/start.svg" alt="Začít lekci">
    </a>
    <a v-else class="cvicBlok" @click="console.log('ssss')"> <!-- TODO -->
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

@media screen and (max-width: 1000px) {
    .cvicBlok {
        width: 180px;
        background-color: var(--tmave-fialova);
        height: 180px;
        transition-duration: 0.2s;
        padding: 15px 15px 30px 15px;
        font-size: 0.8em;
    }

    .fajvkaVetsi, .playVetsi {
        width: 100px;
        height: 30px;
        align-self: center;
    }

    h3 {
        font-size: 20px !important;
    }
}
</style>