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
    height: 40px;
    width: 40px;
    margin-right: 10px;
    padding: 4px;
}

div img {
    height: 32px;
    width: 32px;
    cursor: pointer;
}

@media screen and (max-width: 1100px) {
    div img {
        height: 35px;
        position: relative;
        top: -5px;
    }
}
</style>