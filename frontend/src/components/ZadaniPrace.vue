<script setup lang="ts">
import { ref } from 'vue';
import TextZadani from '../components/TextZadani.vue';
import axios from 'axios';
import { checkTeapot, getToken, pridatOznameni } from '../utils';
import Tooltip from "../components/Tooltip.vue";

const textovePole = ref<InstanceType<typeof TextZadani> | null>(null)

const delka = ref(1)
const typTextu = ref("")

function getText() {
    axios.post("/skola/text", { "typ": typTextu.value }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        console.log(response)
    }).catch(e => {
        if (checkTeapot(e)) return
        console.log(e)
        pridatOznameni("Chyba serveru")
    })

    puvodniText.value = textovePole.value!.text
}

function d(x: number) {
    delka.value = x
}

const puvodniText = ref("")

function smazatDiakritiku() {
    textovePole.value!.text = textovePole.value!.text.normalize("NFD").replace(/[\u0300-\u036f]/g, "")
}

function smazatVelkaPismena() {
    textovePole.value!.text = textovePole.value!.text.toLocaleLowerCase()
}

function smazatEnterAMezery() {
    textovePole.value!.text = textovePole.value!.text.replace(/\n/g, " ").replace(/ {2,}/g, " ")
}

function resetSmazanych() {
    textovePole.value!.text = puvodniText.value
}

</script>
<template>
    <div id="pulic">
        <div id="nastaveni">
            <h2>Nastavení</h2>

            <div id="moznosti">
                <div id="delka">
                    <Tooltip zprava="Pokud student dopíše text před vypršením časového limitu, bude ho psát znovu." :sirka="210" :vzdalenost="0">
                        <h3>Čas</h3>
                    </Tooltip>
                    <hr id="predel2">
                    <button :class="{ aktivni: 1 == delka }" @click="d(1)">1min</button>
                    <button :class="{ aktivni: 2 == delka }" @click="d(2)">2min</button>
                    <button :class="{ aktivni: 3 == delka }" @click="d(3)">3min</button>
                    <button :class="{ aktivni: 5 == delka }" @click="d(5)">5min</button>
                    <button :class="{ aktivni: 10 == delka }" @click="d(10)">10min</button>
                    <button :class="{ aktivni: 15 == delka }" @click="d(15)">15min</button>
                    <button :class="{ aktivni: 20 == delka }" @click="d(20)">20min</button>
                    <button :class="{ aktivni: 30 == delka }" @click="d(30)">30min</button>
                    <button :class="{ aktivni: -1 == delka }" @click="d(-1)">Do konce</button>
                </div>

                <hr id="predel">

                <div class="vertKontejner">
                    <div class="kontejner">
                        <button @click="smazatVelkaPismena" class="tlacitko">Smazat velká písmena</button>
                        <button @click="smazatDiakritiku" class="tlacitko">Smazat diakritiku</button>
                        <button @click="smazatEnterAMezery" class="tlacitko">Smazat mezery</button>
                    </div>

                    <div class="kontejner">
                        <button @click="resetSmazanych">Reset</button>
                    </div>
                </div>
            </div>
        </div>

        <div id="text">
            <div>
                <div>{{ textovePole?.text.length }} znaků</div>
                <select v-model="typTextu" placeholder="Druh textu..." @select="getText">
                    <option value="" selected>Vlastní text</option>
                    <option value="Věty z pohádek">Věty z pohádek</option>
                    <option value="Zeměpis">Zeměpis</option>
                    <option value="Dějepis">Dějepis</option>
                </select>
            </div>

            <TextZadani ref="textovePole" />
        </div>
    </div>
</template>
<style scoped>
.tlacitko {
    width: 250px;
}

.vertKontejner {
    display: flex;
    flex-direction: column;
    justify-content: space-around;
}

select {
    height: 40px;
    border: none;
    border-radius: 5px;
    padding: 5px;
    font-size: 1rem;
    color: white;
    background-color: var(--tmave-fialova);
    cursor: pointer;
    transition: 0.2s;
}

select:hover {
    background-color: var(--fialova) !important;
}

select option {
    background-color: var(--tmave-fialova) !important;
}

select option:disabled {
    color: white !important;
    opacity: 1;
}

#genBtn {
    margin-top: 0;
    width: 105px;
    background-color: var(--tmave-fialova);
}

#genBtn:hover {
    background-color: var(--fialova);
}

#pulic {
    display: flex;
    justify-content: space-between;
    gap: 40px;
    width: 860px;
}

#text {
    width: 410px;
    display: flex;
    flex-direction: column;
    height: 400px;
    border-radius: 10px;
    gap: 15px;
}

#text>div {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

#text>div>div {
    width: 120px;
    height: 40px;
    padding: 10px;
    background-color: var(--tmave-fialova);
    border-radius: 5px;
}

#nastaveni {
    width: 410px;
    display: flex;
    gap: 10px;
    flex-direction: column;
    height: 400px;
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    border-radius: 10px;
}

#moznosti {
    display: flex;
    height: 100%;
}

.kontejner {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    transition: filter 0.2s;
}

.kontejner>div {
    display: flex;
    gap: 15px;
}

.kontejner>div>button {
    width: auto;
    padding: 0 12px;
}

#delka {
    display: flex;
    gap: 10px;
    flex-direction: column;
    justify-content: center;
    flex-wrap: wrap;
    margin-right: 14px;
}

#delka>button {
    width: 90px;
}

#delka h3 {
    position: relative;
    bottom: 15px;
}

#predel {
    margin: 12px 20px 12px 0;
    border: 1px solid var(--fialova);
}

#predel2 {
    margin-top: -22px;
    margin-bottom: 8px;
    width: 70%;
    border: 1px solid rgb(206, 206, 206);
    align-self: center;
}
</style>