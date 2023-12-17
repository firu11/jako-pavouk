<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter()
let routeArray = useRoute().path.split("/")
routeArray.pop() // odstranime posledni
const zpet = routeArray.join("/") //zase spojime a mame to predchozi jakoby

onMounted(() => {
    document.addEventListener('keydown', e1)
})

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
})

function e1(e: KeyboardEvent) { // Escape = šipka zpět
    if (e.key == "Escape") {
        router.push(zpet)
    }
}

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