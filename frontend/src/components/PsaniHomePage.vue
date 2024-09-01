<script setup lang="ts">
import { onMounted, ref } from "vue";

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
            <span v-for="(p, i) in viditelny" :class="{ 'spravne-pismeno': i < 2 }" class="pismeno" :key="p.toString()">
                {{ p[1] != " " ? p[1] : "&nbsp" }}
            </span>
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
    opacity: 0;
}

.pismenka-leave-to {
    transform: translateX(-80px);
    opacity: 0;
}

.pismenka-leave-active {
    position: absolute;
}

#cara {
    width: 65px;
    height: 7px;
    border-radius: 1px;
    background-color: var(--bila);
    position: absolute;
    left: 8.46rem;
    top: 125px;
}

#box {
    height: 8.2rem;
    display: flex;
    position: relative;
    margin: 0;
    margin-bottom: 2rem;
    direction: ltr;
}

.pismeno {
    font-family: "Red Hat Mono", monospace;
    line-height: 1.2;
    text-decoration: none;
    color: var(--bila);
    font-size: 7.1rem;
    font-weight: 600;
    max-width: 67px;
}

.spravne-pismeno {
    color: rgba(100, 100, 100, 20);
}

@media screen and (max-width: 1100px) {
    .pismeno {
        font-size: 3.6rem;
        max-width: 34px;
    }

    #cara {
        width: 34px;
        height: 4px;
        left: 4.3rem;
        top: 62px;
        border-radius: 1px;
    }

    #box {
        height: 4rem;
        margin-top: 50px;
    }

    .pismenka-enter-from {
        transform: translateX(30px); /* nevim proc ale vzdycky to odsoucalo to hamburget tlacitko */
        opacity: 0;
    }

    .pismenka-leave-to {
        transform: translateX(-30px);
        opacity: 0;
    }
}

@media screen and (max-width: 350px) {
    .pismeno {
        font-size: 2.6rem;
        max-width: 25px;
    }
    #cara {
        width: 23px;
        height: 3px;
        left: 3.2rem;
        top: 45px;
        border-radius: 1px;
    }
}
</style>