<script setup lang="ts">
import { useTemplateRef } from 'vue';
import TextZadani from './ucitel/TextZadani.vue';

const textovePole = useTemplateRef("textove-pole")

function getZnakyASlova() {
    if (textovePole.value == null) return "0 / 0"
    let vys = ""

    vys += textovePole.value.text.length
    if (textovePole.value.text.length == 0 || textovePole.value.text.length >= 5) vys += " znaků"
    else if (textovePole.value.text.length == 1) vys += " znak"
    else vys += " znaky"

    vys += " / "
    let slova = 0
    if (textovePole.value.text.length != 0) slova = textovePole.value.text.trim().split(/[ \n]+/).length
    vys += slova
    if (slova == 0 || slova >= 5) vys += " slov"
    else if (slova == 1) vys += " slovo"
    else vys += " slova"

    return vys
}

</script>
<template>
    <h1>Schvalování textu</h1>

    <div id="kontejner">
        <TextZadani ref="textove-pole" />

        <span>{{ getZnakyASlova() }} sus</span>
    </div>
</template>
<style scoped>
#kontejner {
    position: relative;
    top: -60px;
    padding-top: 55px;
    width: 410px;
    height: 80vh;
}

h1 {
    z-index: 1;
}
</style>