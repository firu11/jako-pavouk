<script setup lang="ts">

defineProps<{
    o: number
}>()

</script>
<template>
    <div class="pilulka">
        <div :class="{ easy: o == 1, medium: o == 2, hard: o == 3 }" />
    </div>
</template>
<style scoped>
.pilulka {
    display: flex;
    align-items: end;
    height: 40px;
    border-radius: 10px;
    background-color: rgba(255, 255, 255, 0.25);
    overflow: hidden;
    width: 8px;
    position: relative;
    right: 40px;
}

.pilulka>div {
    width: 100%;
}

.easy {
    height: 33%;
    background-color: rgb(0, 189, 0);
}

.medium {
    height: 66%;
    background-color: orange;
}

.hard {
    height: 100%;
    background-color: red;
}
</style>