<script setup lang="ts">
import { onMounted, ref } from "vue";
import { nastaveniJmeno, prihlasen } from "../stores";
import { saveNastaveni } from "../utils";

const emit = defineEmits(["restart", "toggle"])

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
        let obj: any
        try {
            obj = JSON.parse(nastaveni)
        } catch {
            saveNastaveni(diakritika.value, velkaPismena.value, typ.value, delka.value, klavModel.value)
            return
        }
        
        diakritika.value = obj.diakritika
        velkaPismena.value = obj.velkaPismena
        typ.value = obj.typ
        console.log(obj.delka)
        delka.value = obj.delka ? obj.delka : 60
        klavModel.value = obj.klavesnice
    }
})

function disabledBtn(e: KeyboardEvent) {
    e.preventDefault()
    saveNastaveni(diakritika.value, velkaPismena.value, typ.value, delka.value, klavModel.value)
}

function d(x: number) {
    delka.value = x
    emit("restart")
}

defineExpose({ diakritika, velkaPismena, typ, delka, klavModel })

</script>
<template>
    <div id="psani-menu">

        <div v-if="vyberTextu && prihlasen" class="kontejner" style="gap: 20px;">
            <input v-if="vyberTextu" v-model="typ" type="checkbox" id="toggle" class="toggleCheckbox"
                @change="$emit('restart'); saveNastaveni(diakritika, velkaPismena, typ, delka, klavModel)" />
            <label v-if="vyberTextu" for="toggle" class="toggleContainer">
                <div>Slova</div>
                <div>Věty</div>
            </label>
            <input v-if="!prihlasen" @change="saveNastaveni(diakritika, velkaPismena, typ, delka, klavModel)" v-model="klavModel" type="checkbox"
                id="toggle1" class="toggleCheckbox" />
            <label v-if="!prihlasen" for="toggle1" class="toggleContainer">
                <div>Qwertz</div>
                <div>Qwerty</div>
            </label>
        </div>

        <div class="kontejner">
            <div id="delka">
                <button @keyup="disabledBtn" :class="{ aktivni: 15 == delka }" @click="d(15)">15s</button>
                <button @keyup="disabledBtn" :class="{ aktivni: 30 == delka }" @click="d(30)">30s</button>
                <button @keyup="disabledBtn" :class="{ aktivni: 60 == delka }" @click="d(60)">1min</button>
                <button @keyup="disabledBtn" :class="{ aktivni: 120 == delka }" @click="d(120)">2min</button>
                <button @keyup="disabledBtn" :class="{ aktivni: 300 == delka }" @click="d(300)">5min</button>
                <button @keyup="disabledBtn" :class="{ aktivni: 600 == delka }" @click="d(600)">10min</button>
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
</template>
<style scoped>
#delka {
    display: flex;
    align-items: center;
    gap: 6px;
    padding-top: 4px;
}

#predel {
    width: 92%;
    border: 1px solid var(--fialova);
}

.kontejner {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
    margin: 0 10px;
    cursor: pointer;
    height: 30px;
}
</style>