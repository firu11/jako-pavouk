<script setup lang="ts">
defineProps(["jmeno", "cesta"])
</script>

<template>
    <RouterLink :to="cesta" >
        <div class="hover">

            <svg class="tecka" height="10" width="15">
                <circle cx="5" cy="5" r="5" fill="white" />
            </svg>

            <p>{{ jmeno }}</p>
        </div>
    </RouterLink>
</template>

<style scoped>
a {
    color: var(--bila);
    display: block;
    padding: 5px;
    text-decoration: none;
    cursor: pointer;
}

a p {
    display: inline-block;
    margin-left: 3px;
    position: relative;
    left: -15px;
    transition: 0.3s;
    font-size: 18px;
}

.router-link-active p {
    left: 0;
    transition: 0.3s;
}

a:hover .hover {
    background-color: var(--fialova);
    width: calc(var(--sirka-menu) - 10px);
    transition: background-color 0.3s;
}

.hover {
    padding: 10px;
    border-radius: 8px;
    transition: background-color 0.3s;
}

.tecka {
    opacity: 0;
    transition: 0.2s;
}

.router-link-active .tecka {
    opacity: 100;
    margin-bottom: 1px;
    transition: 0.2s;
}

@media screen and (max-width: 1000px) {
    a {
        color: var(--bila);
        display: block;
        padding: 5px;
        text-decoration: none;
        cursor: pointer;
        font-size: 17px;
    }

    a p {
        display: inline-block;
        margin-left: 3px;
        transition: 0.3s;
    }

    a:hover .hover {
        background-color: var(--tmave-fialova);
        width: calc(var(--sirka-menu) - 10px);
        transition: background-color 0.3s;
    }
}
</style>