<script setup lang="ts">
import { onMounted, ref } from 'vue';

let delka = 9
const counter = ref(delka - 1)
const text = ref([["1", "J"], ["2", "a"], ["3", "k"], ["4", "o"], ["5", " "], ["6", "P"], ["7", "a"], ["8", "v"], ["9", "o"], ["10", "u"], ["11", "k"], ["12", " "]])
const viditelny = ref(text.value.slice(0, delka))

onMounted(() => {
    setTimeout(dalsi, 200)
})


function dalsi() {
    counter.value++
    if (counter.value == text.value.length) {
        counter.value = 0
    }
    viditelny.value.shift()
    viditelny.value.push(text.value[counter.value])
    setTimeout(dalsi, Math.floor(Math.random() * 3) * 300 + 300)
}

</script>

<template>
    <div id="box">
        <TransitionGroup name="pismenka">
            <p v-for="(p, i) in viditelny" :class="{ spravnePismeno: i < 2 }" class="pismeno" :key="p.toString()">
                {{ p[1] != " " ? p[1] : "&nbsp" }}
            </p>
        </TransitionGroup>
        <div id="cara"></div>
    </div>
</template>

<style scoped>
.pismenka-move,
.pismenka-enter-active,
.pismenka-leave-active {
    transition: all 0.2s ease-in-out;

}

.pismenka-enter-from {
    transform: translateX(80px);
}

.pismenka-leave-to {
    transform: translateX(-80px);
}

.pismenka-leave-active {
    position: absolute;
}

#cara {
    width: 64px;
    height: 7px;
    border-radius: 1px;
    background-color: var(--bila);
    position: absolute;
    left: 138px;
    top: 125px;
}

#box {
    height: 8.2em;
    display: flex;
    position: relative;
}

.pismeno {
    font-family: 'Red Hat Mono', monospace !important;
    font-size: 7.1em;
    font-weight: 600 !important;
    line-height: 1.2;
    text-decoration: none;
    color: var(--bila);
}

.spravnePismeno {
    color: rgba(100, 100, 100, 20);
}

@media screen and (max-width: 1000px) {
    .pismeno {
        font-size: 3.2em;
    }

    #cara {
        width: 28px;
        height: 4px;
        left: 56px;
        top: 58px;
        border-radius: 1px;
    }

    #box {
        height: 4.6em;
        margin-top: 50px;
    }
}
</style>