<script setup lang="ts">
import axios from "axios";
import { checkTeapot, getToken, pridatOznameni } from "../utils";
import { ref } from "vue";
import Tooltip from "./Tooltip.vue"

const props = defineProps({
    id: {
        type: Number,
        required: true
    },
    kod: {
        type: String,
        required: true
    },
    zamknuta: {
        type: Boolean,
        required: true
    }
})

const zamknuta = ref(props.zamknuta)

function copy() {
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
            <span @click="copy" :class="{ 'zamknuty-kod': zamknuta }">{{ props.kod == undefined ? "------" : props.kod }}</span>
            <Tooltip zprava="Uzamknout/Odemknout přístup novým žákům do této třídy" :sirka="192" :vzdalenost="12">
                <img v-if="!zamknuta" src="../assets/icony/zamekOpen.svg" alt="Odemčená třída" @click="zamek()">
                <img v-else src="../assets/icony/zamekClosed.svg" alt="Zamčená třída" @click="zamek()">
            </Tooltip>
        </div>
        <span :class="{ 'zamknuty-kod': zamknuta }">jakopavouk.cz/zapis/{{ props.kod }}</span>
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
    justify-content: space-around;
    width: 250px;
}

#kod div {
    font-size: 2rem;
    display: flex;
    gap: 5px;
    justify-content: center;
}

#kod div span {
    font-weight: 500;
    padding: 0 6px;
    transition: 0.2s;
    border-radius: 8px;
}

#kod>span {
    transition: 0.2s;
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
    color: rgba(255, 255, 255, 0.3);
    text-decoration: line-through;
}
</style>