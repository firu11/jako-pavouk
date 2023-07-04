<script setup lang="ts">
import axios from 'axios'
import { prihlasen, token_jmeno } from '../stores';
import { useRouter } from 'vue-router';
import { onMounted, ref } from 'vue';
import { get_token } from '../utils';

const router = useRouter()

let info = ref({jmeno: "...", email: "...@...", dokonceno: 0, daystreak: 0, prumer_rychlosti: -1, uspesnost: -1})

function odhlasit() {
    localStorage.removeItem(token_jmeno)
    prihlasen.value = false
    router.push("/")
}

function zaokrouhlit(cislo: number | null) {
    if (cislo == null) {
        return "-"
    }
    return Math.round(cislo * 10) / 10
}

onMounted(() => {
    if (get_token()) {
        axios.get('/ja', {
            headers: {
                Authorization: `Bearer ${get_token()}`
            }
        }).then(response => {
            info.value = response.data
        }).catch(_ => {
            router.push("/prihlaseni")
        })
    } else { //nebudeš tam chodit nemas ucet more
        router.push("/prihlaseni")
    }
})

</script>

<template>
    <div id="ucet">
        <img src="../assets/icony/user.svg" alt="uzivatel">
        <div id="nadpisy">
            <h1>{{ info.jmeno }}</h1>
            <h2>{{ info.email }}</h2>
        </div>
    </div>
    <div id="progres">
        <div id="nacitani-pozadi">
            <div id="nacitani" :style="{ width: info.dokonceno + '%' }"></div>
        </div>
        <span class="popis" style="width: 100%;">Dokončeno: <span class="cislo">{{ zaokrouhlit(info.dokonceno)
        }}%</span></span>
    </div>
    <div id="bloky">
        <div class="blok">
            <img src="../assets/icony/kalendar.svg" alt="Přesnost">
            <span class="popis">Počet dní v řadě: <br><span class="cislo">{{ zaokrouhlit(info.daystreak) }}</span></span>
        </div>
        <div class="blok">
            <img src="../assets/icony/rychlost.svg" alt="Rychlost" width="75">
            <span v-if="info.prumer_rychlosti == -1">Zatím nic</span>
            <span v-else class="popis">Rychlost: <br><span class="cislo">{{ zaokrouhlit(info.prumer_rychlosti) }}</span>
                CPM</span>
        </div>
        <div class="blok">
            <img src="../assets/icony/terc.svg" alt="Přesnost">
            <span v-if="info.uspesnost == -1">Zatím nic</span>
            <span v-else class="popis">Přesnost: <br><span class="cislo">{{ zaokrouhlit(info.uspesnost) }}</span> %</span>
        </div>
    </div>

    <button @click="odhlasit" class="tlacitko">Odhlásit</button>
</template>

<style scoped>
.popis {
    font-size: 15pt;
    width: 60%;
}

.cislo {
    font-size: 28pt;
    font-weight: 500;
}

#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
}

.blok {
    display: flex;
    text-decoration: none;
    border-radius: 10px;
    justify-content: space-evenly;
    align-items: center;
    width: 320px;
    background-color: var(--tmave-fialova);
    height: 120px;
    transition-duration: 0.2s;
    padding: 15px;
    gap: 10px;
}

#progres {
    margin-bottom: 20px;
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 420px;
    background-color: var(--tmave-fialova);
    height: 110px;
    transition-duration: 0.2s;
    padding: 15px;
    gap: 10px;
}

#nadpisy {
    display: flex;
    flex-direction: column;
    justify-content: center;
}

#nadpisy h1 {
    margin-bottom: 0;
    align-self: flex-start;
}

#ucet img {
    height: 100px;
}

#ucet {
    display: flex;
    background-color: var(--tmave-fialova);
    margin-bottom: 40px;
    padding: 15px 30px 15px 15px;
    border-radius: 10px;
    gap: 15px;
}

.tlacitko {
    margin-top: 20px;
    background-color: var(--tmave-fialova);
}

.tlacitko:hover {
    background-color: var(--fialova);
}

#bloky div img {
    height: 65px;
}

#bloky div h2 {
    font-weight: 500;
}

#nacitani-pozadi {
    height: 20px;
    background-color: var(--fialova);
    border-radius: 5px;
    padding: 0;
}

#nacitani {
    background-color: var(--bila);
    height: 20px;
    border-radius: 5px;
    position: relative;
    left: -1px
}
</style>