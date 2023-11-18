<script setup lang="ts">
import { useRoute } from 'vue-router';

let routeArray = useRoute().path.split("/")
routeArray.pop() // odstranime posledni
const zpet = routeArray.join("/") //zase spojime a mame to predchozi jakoby

</script>

<template>
    <div>
        <RouterLink :to="zpet">
            <img src="../assets/icony/sipkaL.svg" alt="Zpět">
        </RouterLink>
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

@media screen and (max-width: 1100px) {
    div img {
        height: 38px;
    }
}
</style>