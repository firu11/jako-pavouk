<script setup lang="ts">

defineProps<{
    o: number
}>()

</script>
<template>
    <div v-if="o == 1" class="pilulka">
        <div />
        <div />
        <div class="easy" />
    </div>
    <div v-else-if="o == 2" class="pilulka">
        <div />
        <div class="medium" />
        <div class="medium" />
    </div>
    <div v-else-if="o == 3" class="pilulka">
        <div v-for="_ in 3" class="hard" />
    </div>
</template>
<style scoped>
.pilulka {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    align-items: end;
    height: 40px;
    border-radius: 10px;
    background-color: var(--cerna);
    overflow: hidden;
    width: 8px;
}

.pilulka>div {
    width: 100%;
    background-color: rgba(255, 255, 255, 0.25);
    height: 11.666px;
    transition: 0.3s;
}

.easy {
    background-color: rgb(0, 189, 0) !important;
}

.medium {
    background-color: orange !important;
}

.hard {
    background-color: red !important;
}

.pilulka:hover>div:nth-child(3) {
    height: 100%;
}

.pilulka:hover>div:nth-child(1),
.pilulka:hover>div:nth-child(2) {
    height: 0%;
}

@media screen and (max-width: 1100px) {
    .pilulka {
        right: 30px;
    }
}
</style>