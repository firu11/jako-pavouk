<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import Tooltip from "../components/Tooltip.vue";

onMounted(() => {
    document.addEventListener("keydown", e1)
})

onUnmounted(() => {
    document.removeEventListener("keydown", e1)
})

function e1(e: KeyboardEvent) { // Escape = šipka zpět
    if (e.key == "Escape") {
        zpatky()
    }
}

const router = useRouter()

function zpatky() {
    if (history.state.current == "/test-psani") {
        router.push("/procvic")
        return
    }
    if (history.state.back == null) {
        history.back()
        return
    }

    const current: Array<string> = history.state.current.split("/")
    const back: Array<string> = history.state.back.split("/")
    current.pop()
    current.shift()
    back.pop()
    back.shift()
    if (current.toString() === back.toString()) router.push("/" + current.join("/")) // pokud píšu cvičení kam jsem šel přes pokračovat, aby to neskočilo zase na to první ale zpátky na lekci
    else history.back() // používám history, aby to scrollovalo tam kde jsem skoncil
}

</script>

<template>
    <Tooltip zprava="Vrátit se zpět <span class='klavesa-v-textu-mensi'>Escape</span>" :sirka="110" :vzdalenost="-28">
        <div @click="zpatky">
            <img src="../assets/icony/sipkaL.svg" alt="Zpět">
        </div>
    </Tooltip>
</template>

<style scoped>
div {
    height: 40px;
    width: 40px;
    margin-right: 10px;
    padding: 4px;
}

div img {
    height: 32px;
    width: 32px;
    cursor: pointer;
    position: relative;
    top: -3px;
}

@media screen and (max-width: 1100px) {
    div img {
        height: 35px;
        top: -7px;
    }
}
</style>