<script setup lang="ts">
import axios from 'axios';
import { onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getToken } from '../utils';

const emit = defineEmits(["restart"])

const props = defineProps({
    preklepy: {
        type: Number,
        default: 0
    },
    delkaTextu: {
        type: Number,
        default: 1
    },
    cas: {
        type: Number,
        default: 1
    },
    casF: String,
    pismena: {
        type: String,
        default: ""
    },
    cislo: String,
    posledni: Boolean
})

let rychlost = Math.round((props.delkaTextu / props.cas) * 60 * 10) / 10
const route = useRoute()
const router = useRouter()

function reset() {
    emit("restart")
}

function dalsi() {
    if (props.cislo == undefined) return
    let r = route.path.split("/")
    r.pop()
    let c = r.join("/") 
    if (props.posledni) router.push(c) // /lekce/pismena
    else router.push(c + "/" + (parseInt(props.cislo) + 1).toString()) // /lekce/pismena/cislo
}

onMounted(() => {
    axios.post('/dokonceno/' + encodeURIComponent(props.pismena) + '/' + props.cislo, {
        "cpm": rychlost,
        "preklepy": props.preklepy
    }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).catch(function (e) {
        console.log(e)
    })
})

</script>

<template>
    <div id="bloky">
        <div class="blok">
            <h2>{{ rychlost }}</h2>
            <hr>
            <p class="jednotka">CPM</p>
            <h3>Rychlost</h3>
        </div>
        <div class="blok">
            <div>
                <h2>{{ preklepy }}</h2>
                <h3 class="procento">({{ Math.round(preklepy / delkaTextu * 1000) / 10 }}%)</h3>
            </div>
            <hr>
            <p class="jednotka">&zwnj;</p>
            <h3>Překlepy</h3>
        </div>
        <div class="blok">
            <h2>{{ casF }}s</h2>
            <hr>
            <p class="jednotka">&zwnj;</p>
            <h3>Čas</h3>
        </div>
    </div>
    <div id="tlacitka_kontainer">
        <button class="tlacitko" @click="reset">Zkusit znovu</button>
        <button class="tlacitko" @click="dalsi()">Pokračovat</button>
    </div>
</template>

<style scoped>
#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
    margin-top: 25px;
}

.blok div {
    display: flex;
    align-items: baseline;
    justify-content: center;
    gap: 5px;
}

.procento {
    font-size: 20px;
    position: relative;
    top: -3px;
}

.blok {
    color: var(--bila);
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 200px;
    background-color: var(--tmave-fialova);
    height: 140px;
    transition-duration: 0.2s;
    padding: 15px 15px 30px 15px;
}

.blok h2 {
    font-size: 40px;
    font-weight: 500;
}

.blok p {
    font-size: 14px;
    margin-bottom: 8px;
    margin-top: 4px;
}

.blok hr {
    width: 160px;
    align-self: center;
}

#tlacitka_kontainer {
    display: inline-flex;
    gap: 20px;
    margin-top: 10px;
}
</style>