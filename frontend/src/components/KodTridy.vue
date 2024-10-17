<script setup lang="ts">
import axios from "axios";
import { checkTeapot, getToken, pridatOznameni } from "../utils";
import { ref } from "vue";
import Tooltip from "./Tooltip.vue"

const props = defineProps({
    id: Number,
    kod: String,
    zamknuta: Boolean,
})

const zamknuta = ref(props.zamknuta)

function copy() {
    if (!props.kod) return
    if (zamknuta.value) {
        pridatOznameni("Třída je zamknutá, a tak se do ní nikdo další nemůže připojit.")
        return
    }
    navigator.clipboard.writeText(props.kod)
    pridatOznameni("Zkopírováno!", undefined, "copy")
}

function zamek() {
    zamknuta.value = !zamknuta.value

    axios.post("/skola/zmena-tridy", { trida_id: props.id, zmena: "zamek" }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).catch(e => {
        if (!checkTeapot(e)) {
            console.log(e)
            pridatOznameni("Chyba serveru")
        }
        zamknuta.value = !zamknuta.value
    })
}

</script>
<template>
    <div id="kod">
        <div>
            <Tooltip zprava="Kliknutím zkopírujete" :sirka="192" :vzdalenost="6">
                <span @click="copy" :class="{ 'zamknuty-kod': zamknuta }">{{ props.kod == undefined ? "----" : props.kod }}</span>
            </Tooltip>
            <Tooltip zprava="Uzamknout/Odemknout přístup novým žákům do této třídy" :sirka="192" :vzdalenost="6">
                <img v-if="!zamknuta" src="../assets/icony/zamekOpen.svg" alt="Odemčená třída" @click="zamek()">
                <img v-else src="../assets/icony/zamekClosed.svg" alt="Zamčená třída" @click="zamek()">
            </Tooltip>
        </div>
        <a :class="{ 'zamknuty-kod': zamknuta }" :href="'https://jakopavouk.cz/zapis/' + props.kod" target="_blank">jakopavouk.cz/zapis/{{ props.kod }}</a>
    </div>
</template>
<style scoped>
#kod {
    background-color: var(--tmave-fialova);
    padding: 10px;
    border-radius: 8px;
    transition: 0.1s;
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 2px;
    width: 250px;
}

#kod>div {
    font-size: 32px;
    display: flex;
    gap: 5px;
    justify-content: center;
}

#kod div span {
    font-weight: 500;
    padding: 0 6px 0 calc(6px + 6px);
    border-radius: 8px;
    letter-spacing: 6px;
}

#kod>div * {
    transition: 0.2s;
}

#kod a {
    text-decoration: none;
    color: var(--bila);
    font-weight: 300;
    opacity: 0.9;
    width: auto;
}

#kod div span:hover {
    background-color: var(--fialova);
    cursor: grab;
}

#kod img {
    width: 30px;
    height: 30px;
    position: relative;
    top: 4px;
    transition: 0.2s;
    cursor: pointer;
}

#kod img:hover {
    transform: scale(1.1);
}

.zamknuty-kod {
    color: rgba(255, 255, 255, 0.3) !important;
    text-decoration: line-through !important;
}
</style>