<script setup lang="ts">
import { onMounted, ref } from "vue";
import { nastaveniJmeno, prihlasen } from "../stores";
import { saveNastaveni } from "../utils";

defineEmits(["restart", "toggle", "toggleInterpunkce"])

defineProps({
    vyberTextu: {
        type: Boolean,
        default: true
    },
})

const typ = ref(true) // false = slova, true = vety
const delka = ref(60)
const diakritika = ref(true)
const velkaPismena = ref(false)
const klavModel = ref(false)

onMounted(() => {
    let nastaveni = localStorage.getItem(nastaveniJmeno)
    if (nastaveni !== null) {
        let obj: { diakritika: boolean, velkaPismena: boolean, vetySlova: boolean, delka: number, klavesnice: boolean }
        try {
            obj = JSON.parse(nastaveni)
        } catch {
            console.log("nepodařilo se načíst nastavení")
            saveNastaveni(diakritika.value, velkaPismena.value, typ.value, delka.value, klavModel.value)
            return
        }

        if (obj.diakritika == undefined || obj.velkaPismena == undefined || obj.vetySlova == undefined || obj.delka == undefined || obj.klavesnice == undefined) {
            saveNastaveni(diakritika.value, velkaPismena.value, typ.value, delka.value, klavModel.value)
            return
        }

        diakritika.value = obj.diakritika
        velkaPismena.value = obj.velkaPismena
        typ.value = obj.vetySlova
        delka.value = obj.delka ? obj.delka : 60
        klavModel.value = obj.klavesnice
    }
})

function disabledBtn(e: KeyboardEvent) {
    e.preventDefault()
}

function d(x: number) {
    delka.value = x
    saveNastaveni(diakritika.value, velkaPismena.value, typ.value, delka.value, klavModel.value)
}

defineExpose({ diakritika, velkaPismena, typ, delka, klavModel })

</script>
<template>
    <div id="kontejneros">
        <div id="psani-menu">
            <div v-if="!(prihlasen && !vyberTextu)" class="kontejner" style="gap: 20px;">
                <input v-if="vyberTextu" v-model="typ" type="checkbox" id="toggle" class="toggle-checkbox"
                    @change="$emit('restart'); saveNastaveni(diakritika, velkaPismena, typ, delka, klavModel)" />
                <label v-if="vyberTextu" for="toggle" class="toggle-contejner">
                    <div>Slova</div>
                    <div>Věty</div>
                </label>
                <input v-if="!prihlasen" @change="saveNastaveni(diakritika, velkaPismena, typ, delka, klavModel)" v-model="klavModel" type="checkbox"
                    id="toggle1" class="toggle-checkbox" />
                <label v-if="!prihlasen" for="toggle1" class="toggle-contejner">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </div>

            <div class="kontejner">
                <div id="delka">
                    <button @keyup="disabledBtn" :class="{ aktivni: 15 == delka }" @click="d(15)" in="15s">15s</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 30 == delka }" @click="d(30)" in="30s">30s</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 60 == delka }" @click="d(60)" in="1min">1min</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 120 == delka }" @click="d(120)" in="2min">2min</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 300 == delka }" @click="d(300)" in="5min">5min</button>
                    <button @keyup="disabledBtn" :class="{ aktivni: 600 == delka }" @click="d(600)" in="10min">10min</button>
                    <!-- "in" je kvuli css ::after -->
                </div>
            </div>

            <hr id="predel">

            <div class="kontejner">
                <label for="toggle2" class="kontejner">
                    <input v-model="velkaPismena" @change="$emit('toggle')" type="checkbox" id="toggle2" class="radio" />
                    Velká písmena
                </label>

                <label for="toggle3" class="kontejner">
                    <input v-model="diakritika" @change="$emit('toggle')" type="checkbox" id="toggle3" class="radio" />
                    Diakritika
                </label>
            </div>
        </div>
    </div>
</template>
<style scoped>
#delka {
    display: flex;
    align-items: center;
    gap: 6px;
    padding-top: 4px;
}

#delka button::after {
    display: block;
    content: attr(in);
    font-weight: 500;
    height: 1px;
    color: transparent;
    overflow: hidden;
    visibility: hidden;
}

#predel {
    width: 92%;
    border: 1px solid var(--fialova);
}

.kontejner {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 8px;
    margin: 0 10px;
    cursor: pointer;
    height: 30px;
}

#kontejneros {
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: calc(370px);
    height: 236px;
    width: 675px;
}

#psani-menu {
    background-color: var(--tmave-fialova);
    padding: 10px;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    position: absolute;
    flex-wrap: wrap;
    align-items: center;
    box-shadow: 0 0 50px 0 rgba(0, 0, 0, 10);
    transition: opacity 0.2s;
}

#psani-menu button,
#delka button {
    background-color: transparent;
    border: none;
    color: var(--bila);
    transition: 0.1s;
    font-size: 16px;
    border-radius: 5px;
    padding: 0 2px;
}

#psani-menu button:hover,
#delka button:hover {
    color: var(--svetle-fialova);
    font-weight: 500;
    color: var(--bila);
    cursor: pointer;
}

.aktivni {
    color: var(--svetle-fialova) !important;
    font-weight: 500 !important;
}
</style>