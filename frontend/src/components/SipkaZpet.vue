<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';

const router = useRouter()
let routeArray = useRoute().path.split("/")
routeArray.pop() // odstranime posledni
const zpet = routeArray.join("/") //zase spojime a mame to predchozi jakoby

</script>

<template>
    <div>
        <img src="../assets/icony/sipkaL.svg" alt="Zpět" @click="router.push(zpet)">
    </div>
</template>

<style scoped>
div {
    height: 39px;
    margin-right: 10px;
}

div img {
    height: 42px;
    cursor: pointer;
}
</style>