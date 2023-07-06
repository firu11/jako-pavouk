<script setup lang="ts">

</script>

<template>
    <h1>Chtěl bys podpořit tento projekt?</h1>
    
    <p>Vyjádřit podporu tomuto projektu můžete přes platformu Ko-Fi</p>
    <a href="https://ko-fi.com/jakopavouk"><img id="kofi" src="../assets/kofi_button.webp" alt="Podpořit"></a>
</template>

<style scoped>
#kofi {
    width: 15em;
    transition: 0.2s;
    cursor: pointer;
    margin-top: 20px;
}

#kofi:hover {
    filter:contrast(65%);
    transition: 0.2s;
}
</style>