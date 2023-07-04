<script setup lang="ts">
import axios from 'axios'
import { onMounted, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router'
import { formatovany_pismena, get_token } from '../utils';

const pismena = useRoute().params.pismena
const router = useRouter()

const cviceni = ref([] as {id: number, typ: string}[])
const dokoncene = ref([])

onMounted(() => {
    axios.get('/lekce/' + pismena, {
        headers: {
            Authorization: `Bearer ${get_token()}`
        }
    }).then(response => {
        if (response.data.cviceni === null) {
            router.push('/404')
        }
        cviceni.value = response.data.cviceni
        dokoncene.value = response.data.cviceni

    }).catch(_ => {
        router.push('/404')
    })
})

function prihlaste_se() {
    alert('Nejprve se prosím přihlašte')
}

function index1(index: number) {
    return index + 1
}

function je_dokoncene(id: number) {
    return dokoncene.value.includes(id)
}

</script>

<template>
    <h1>
        <router-link class="tlacZpet" :to="'/lekce'">
            <img src="../assets/icony/sipkaL.svg" alt="Zpět">
        </router-link>
        Lekce: {{ formatovany_pismena(pismena) }}
    </h1>
    <div class="kontejnr">
        <div v-if="cviceni.length !== 0" v-for="({id, typ}, index) in cviceni">
            <h2>
                <router-link class="lekceBlok" :class="{ dokoncenyBlok: je_dokoncene(id) }" v-if="typ === 'nova'"
                    :to="'/lekce/' + pismena + '/' + index1(index)">
                    <h2>{{ index1(index) }}</h2>
                    <hr>
                    <h3>Nová písmenka</h3>
                    <img class="fajvkaVetsi" v-if="je_dokoncene(id)" src="../assets/icony/right.svg" alt="Dokonceno!">
                    <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
                </router-link>
                <router-link class="lekceBlok" :class="{ dokoncenyBlok: je_dokoncene(id) }"
                    v-else-if="typ === 'naucena'" :to="'/lekce/' + pismena + '/' + index1(index)">
                    <h2>{{ index1(index) }}</h2>
                    <hr>
                    <h3>Probraná písmenka</h3>
                    <img class="fajvkaVetsi" v-if="je_dokoncene(id)" src="../assets/icony/right.svg" alt="Dokonceno!">
                    <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
                </router-link>
                <router-link v-else class="lekceBlok" :class="{ dokoncenyBlok: je_dokoncene(id) }"
                    :to="'/lekce/' + pismena + '/' + index1(index)">
                    <h2>{{ index1(index) }}</h2>
                    <hr>
                    <h3>Se slovy</h3>
                    <img class="fajvkaVetsi" v-if="je_dokoncene(id)" src="../assets/icony/right.svg" alt="Dokonceno!">
                    <img class="playVetsi" v-else src="../assets/icony/start.svg" alt="Začít lekci">
                </router-link>
            </h2>
        </div>
        <p v-else>Tato lekce zatím nemá žádná cvičení</p>
    </div>
    <!-- <p v-else>Něco se pokazilo</p> -->
</template>

<style scoped>
.kontejnr {
    display: flex;
    gap: 15px;
    max-width: 700px;
    flex-wrap: wrap;
    justify-content: center;
}

.lekceBlok {
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

.lekceBlok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.lekceBlok hr {
    width: 160px;
    align-self: center;
    margin: 5px;
}

.dokoncenyBlok {
    opacity: 50%;
}

.lekceBlok h3 {
    align-self: center;
    font-size: 24px;
    height: 100px;
}

.lekceBlok a {
    text-decoration: none;
    color: var(--bila);
    cursor: pointer;
}

h1 {
    display: inline-flex;
    position: relative;
    right: 25px;
    /* posunuti o pulku sipky */
}
</style>
