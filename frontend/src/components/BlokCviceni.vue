<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { levelyPresnosti, levelyRychlosti, prihlasen } from '../stores';
import { pridatOznameni } from '../utils';

const props = defineProps({
    dokonceno: Boolean,
    index: Number,
    pismena: {
        type: String,
        default: ""
    },
    typ: String,
    rychlost: {
        type: Number,
        default: 0
    },
    presnost: {
        type: Number,
        default: 0
    },
    oznacena: Boolean
})

const mobil = ref(document.body.clientWidth <= 1000)
const hvezdy = ref(0)

onMounted(() => {
    if (props.rychlost >= levelyRychlosti[1] && props.presnost >= levelyPresnosti[1]) { // paradni
        hvezdy.value = 3
    } else if (props.rychlost >= levelyRychlosti[0] && props.rychlost < levelyRychlosti[1] && props.presnost >= levelyPresnosti[1]) { // rychlost muze byt lepsi
        hvezdy.value = 2
    } else if (props.presnost >= levelyPresnosti[0] && props.presnost < levelyPresnosti[1] && props.rychlost >= levelyRychlosti[1]) { // presnost muze byt lepsi
        hvezdy.value = 2
    } else if (props.presnost >= levelyPresnosti[0] && props.presnost < levelyPresnosti[1] && props.rychlost >= levelyRychlosti[0] && props.rychlost <= levelyRychlosti[1]) { // oboje muze byt lepsi
        hvezdy.value = 1
    } else if (props.rychlost < levelyRychlosti[0] && props.presnost < levelyPresnosti[0]) { // oboje bad
        hvezdy.value = 0
    } else if (props.rychlost < levelyRychlosti[0]) { // rychlost bad
        hvezdy.value = 0
    } else if (props.presnost < levelyPresnosti[0]) { // presnost bad
        hvezdy.value = 0
    }
})

</script>

<template>
    <router-link v-if="prihlasen && typ !== '...' && !mobil" class="cvicBlok" :class="{ dokoncenyBlok: dokonceno, oznacene: oznacena }"
        :to="'/lekce/' + pismena + '/' + index">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova' || typ === 'programator'">Se slovy</h3>
        <h3 v-else>...</h3>
        <div v-if="dokonceno" id="hvezdy">
            <img v-if="hvezdy >= 1" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
            <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
            <img v-if="hvezdy >= 2" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
            <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
            <img v-if="hvezdy == 3" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
            <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
        </div>
        <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
    </router-link>
    <a v-else-if="typ === '...' && !mobil" class="cvicBlok"> <!-- aby na to ńeslo kliknout nez se to nacte -->
        <h2>{{ index }}</h2>
        <hr>
        <h3>...</h3>
        <img class="playVetsi" src="../assets/icony/start.svg" alt="Začít lekci">
    </a>
    <a v-else-if="!mobil" class="cvicBlok" @click="pridatOznameni('Bez přihlášení si můžeš psaní vyzkoušet v sekci Procvičování')">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova'">Se slovy</h3>
        <h3 v-else>...</h3>
        <img class="playVetsi" src="../assets/icony/start.svg" alt="Začít lekci">
    </a>
    <a v-else class="cvicBlok" :class="{ dokoncenyBlok: dokonceno }"
        @click="pridatOznameni('Psaní na telefonech zatím neučíme...')">
        <h2>{{ index }}</h2>
        <hr>
        <h3 v-if="typ === 'nova'">Nová písmenka</h3>
        <h3 v-else-if="typ === 'naucena'">Probraná písmenka</h3>
        <h3 v-else-if="typ === 'slova'">Se slovy</h3>
        <h3 v-else>...</h3>
        <div v-if="dokonceno" id="hvezdy">
            <img v-if="hvezdy >= 1" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
            <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
            <img v-if="hvezdy >= 2" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
            <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
            <img v-if="hvezdy == 3" src="../assets/icony/hvezda.svg" alt="Hvezda" class="hvezda">
            <img v-else src="../assets/icony/hvezdaPrazdna.svg" alt="Hvezda" class="hvezda">
        </div>
        <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
    </a>
</template>

<style scoped>
.hvezda {
    width: 45px;
    height: 45px;
}

#hvezdy :nth-child(2) {
    position: relative;
    top: -10px;
}

#hvezdy {
}

.cvicBlok {
    color: var(--bila);
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 170px;
    background-color: var(--tmave-fialova);
    height: 220px;
    transition-duration: 0.2s;
    padding: 15px 15px 15px 15px;
}

.cvicBlok:hover, .oznacene {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.cvicBlok hr {
    width: 135px;
    align-self: center;
    margin: 5px;
    margin-top: 2px;
    color: gray;
}

.dokoncenyBlok {
    opacity: 70%;
}

.cvicBlok h3 {
    align-self: center;
    font-size: 23px;
    height: 60px;
    margin-bottom: 22px
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

@media screen and (max-width: 1100px) {
    .cvicBlok {
        max-width: 155px;
        width: auto;
        background-color: var(--tmave-fialova);
        height: 180px;
        transition-duration: 0.2s;
        padding: 15px;
        font-size: 0.8em;
    }

    .fajvkaVetsi,
    .playVetsi {
        width: 100px;
        height: 30px;
        align-self: center;
        margin-bottom: 10px;
    }

    h3 {
        font-size: 18px !important;
    }

    .cvicBlok hr {
        width: 130px;
    }

    .hvezda {
        width: 35px;
        height: 35px;
    }
}
</style>