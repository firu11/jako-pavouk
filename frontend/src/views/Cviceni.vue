<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router';
import { formatovanyPismena, getToken } from '../utils';
import SipkaZpet from '../components/SipkaZpet.vue';
import { computed, onMounted, ref } from 'vue';
import axios from 'axios';
import { onUnmounted } from 'vue';
import Vysledek from '../components/Vysledek.vue';


const router = useRouter()
const route = useRoute()
const pismena: string = Array.isArray(route.params.pismena) ? route.params.pismena[0] : route.params.pismena
const cislo: string = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id

const text = ref([[]] as { id: number, znak: string, spatne: boolean, }[][])
const delkaTextu = ref(0)
const counter = ref(0)
const preklepy = ref(0)
const timerZacatek = ref(0)
const cas = ref(0)

const capslock = ref(false)
let interval: number
const konec = ref(false)

const casFormat = computed(() => {
    return cas.value < 60 ? Math.floor(cas.value).toString() : `${Math.floor(cas.value / 60)}:${cas.value % 60 < 10 ? "0" + Math.floor(cas.value % 60).toString() : Math.floor(cas.value % 60)}`
})

const progress = computed(() => {
    return delkaTextu.value !== 0 ? Math.floor(((counter.value) / delkaTextu.value) * 100) : 0
})

const aktivniPismeno = computed(() => {
    for (const slovo of text.value) {
        for (const pismenoObj of slovo) {
            if (pismenoObj.id === counter.value) {
                return pismenoObj
            }
        }
    }
    return { id: -1, znak: "", spatne: 0 }
})

onMounted(() => {
    axios.get("/cvic/" + pismena + "/" + cislo, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        response.data.text.forEach((slovo: string, i: number) => {
            text.value.push([])
            const slovoArr = [...slovo]
            slovoArr.forEach(pismeno => {
                text.value[i].push({ id: delkaTextu.value, znak: pismeno, spatne: false })
                delkaTextu.value++
            })
        })
    }).catch(_ => {
        router.push('/404')
    });

    document.addEventListener("keypress", klik)
    document.addEventListener("keydown", capslockCheck)
})

onUnmounted(() => {
    document.removeEventListener("keypress", klik)
    document.removeEventListener("keydown", capslockCheck)
})

function capslockCheck(e: KeyboardEvent) { // TODO chtelo by to checknout hned po nacteni stranky ale nevim jestli to jde
    capslock.value = e.getModifierState("CapsLock")
}

function klik(e: KeyboardEvent) {
    if (e.key === " ") {
        e.preventDefault() // ať to nescrolluje
    }

    startTimer()

    if (e.key === aktivniPismeno.value.znak) {
        counter.value++
    } else {
        aktivniPismeno.value.spatne = true
        preklepy.value++
    }

    if (aktivniPismeno.value.id === -1) { // konec
        clearInterval(interval)
        calcCas() // naposledy
        konec.value = true
        document.removeEventListener("keypress", klik)
        document.removeEventListener("keydown", capslockCheck)
    }
}

function startTimer() {
    if (timerZacatek.value === 0) {
        timerZacatek.value = Date.now()
        interval = setInterval(calcCas, 100)
    }
}

function calcCas() {
    cas.value = (Date.now() - timerZacatek.value) / 1000
}

function restart() {
    timerZacatek.value = 0
    cas.value = 0
    counter.value = 0
    preklepy.value = 0
    konec.value = false
    document.addEventListener("keypress", klik)
    document.addEventListener("keydown", capslockCheck)
}

</script>

<template>
    <h1 class="nadpisSeSipkou" style="margin: 0;">
        <SipkaZpet />
        Lekce: {{ formatovanyPismena(pismena) }}
    </h1>
    <h2>Cviceni: {{ cislo }}</h2>

    <div v-if="!konec">
        <div id="nabidka">
            <h3 id="cas">{{ casFormat }}s</h3>
            <h3 :style="{ visibility: capslock ? 'visible' : 'hidden' }" id="capslock">CapsLock</h3>
            <h3 id="preklepy">Chyby: {{ preklepy }}</h3>
        </div>

        <div id="ramecek">
            <div id="text">
                <div class="slovo" v-for="s in text">
                    <div v-for="p in s" class="pismeno"
                        :class="{ podtrzeni: p.id === counter, spatnePismeno: p.spatne && counter === p.id, opravenePismeno: p.spatne && counter > p.id, spravnePismeno: !p.spatne && counter > p.id }">
                        {{ (p.znak !== " " ? p.znak : p.spatne ? "_" : "&nbsp") }}
                    </div>
                </div>
            </div>
        </div>
        <div id="bar">
            <div :style="'width:' + progress + '%; border-bottom-right-radius:' + (progress === 100 ? '10px' : '0')"
                id="progress">&nbsp{{ progress }}%&nbsp
            </div>
        </div>
    </div>

    <Vysledek v-else @restart="restart" :preklepy="preklepy" :delkaTextu="delkaTextu" :casF="casFormat" :cas="cas" :pismena="pismena" :cislo="cislo"></Vysledek>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Red+Hat+Mono&display=swap');

#nabidka {
    margin: 20px 0 6px 0;
    width: var(--sirka-textoveho-pole);
}

#cas {
    float: left;
    width: 150px;
    display: block;
    text-align: left;
}

#preklepy {
    float: right;
    display: block;
    width: 150px;
    text-align: right;
}

#capslock {
    display: inline-block;
    color: red;
    font-weight: bold;
}

#ramecek {
    padding: 10px;
    min-height: 190px;
    border-radius: 10px 10px 0 0;
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
}

#text {
    display: flex;
    flex-wrap: wrap;
}

.slovo {
    display: flex;
    flex-wrap: nowrap;
}

.pismeno {
    border-radius: 3px;
    display: inline-flex;
    font-family: 'Red Hat Mono', monospace;
    font-size: 25px;
    line-height: 1.2;
    text-decoration: none;
    padding: 0 1px;
    margin-right: 1px;
    border-bottom: 3px solid rgba(255, 255, 255, 0);
    /* aby se nedojebala vyska lajny když jdu na dalsi radek*/
    color: var(--bila);
    transition: 60ms;
}

#progress {
    height: 20px;
    background-color: var(--fialova);
    width: 0;
    border-bottom-left-radius: 10px;
    transition: ease 0.5s;
    text-align: right;
}

#bar {
    background-color: var(--tmave-fialova);
    width: var(--sirka-textoveho-pole);
    border-radius: 0 0 10px 10px;
}

.spravnePismeno {
    color: #afafaf;
}

.podtrzeni {
    border-bottom: 3px solid var(--bila);
    border-radius: 0;
    transition: 60ms;
}

.spatnePismeno {
    color: #ff0000;
}

.opravenePismeno {
    color: #b1529c;
}
</style>