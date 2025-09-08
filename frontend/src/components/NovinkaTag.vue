<script setup lang="ts"></script>
<template>
    <span>Novinka</span>
</template>
<style scoped>
span {
    font-size: 0.9rem;
    padding: 3px 6px;
    border-radius: 10px;
    background-color: var(--fialova);
    color: white;
    display: inline-block;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.3);
    animation: pulse-scale 2s ease-in-out infinite;
}

@keyframes pulse-scale {
    0%,
    100% {
        transform: scale(1);
    }
    50% {
        transform: scale(1.1);
    }
}
</style>
