<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter()

function hop() {
    router.push("/jak-psat")
}
</script>

<template>
    <div id="box">
        <div>
            <h2>Nevíš jak na to?</h2>
            <p>Mrkni na trochu teorie! Základem je výchozí pozice.</p>
            <button class="tlacitko" :onclick="hop">Jak psát</button>
        </div>
        <img src="../assets/theThinker.svg" alt="Pavouk the thinker" width="180" height="157">
    </div>
</template>

<style scoped>
#box {
    background-color: var(--tmave-fialova);
    border-radius: 8px;
    padding: 16px 15px 12px 25px;
    display: flex;
    width: 500px;
    align-items: center;
    justify-content: space-between;
    margin-top: 5px;
}

#box div {
    display: flex;
    flex-direction: column;
}

h2 {
    font-size: 1.6rem;
    margin-bottom: 12px;
    font-weight: 300;
    width: 100%;
}

p {
    font-size: 1.15rem;
    margin-bottom: 5px;
}

img {
    width: 180px;
}

@media screen and (max-width: 1100px) {
    #box {
        width: 100%;
        flex-direction: column;
        text-align: center;
        gap: 12px;
        padding-top: 16px;
    }

    #box div {
        align-items: center;
    }
}

</style>