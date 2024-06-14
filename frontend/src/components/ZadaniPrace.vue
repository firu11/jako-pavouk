<script setup lang="ts">
import { ref } from 'vue';
import TextZadani from '../components/TextZadani.vue';

const textovePole = ref<InstanceType<typeof TextZadani> | null>(null)

const delka = ref(1)
const velkaPismena = ref(true)
const diakritika = ref(true)

function getText() {
    textovePole.value!.text = "ssdkfsjdnkfjnk"

    puvodniText.value = textovePole.value!.text
}

function d(x: number) {
    delka.value = x
}

const puvodniText = ref("")
function toggleDiakritikaAVelkaPismena() {
    puvodniText.value = textovePole.value!.text
    let arr = Array.from(textovePole.value!.text)

    if (!diakritika.value && !velkaPismena.value) {
        for (let i = 0; i < arr.length; i++) {
            arr[i] = arr[i].normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLocaleLowerCase()
        }
    } else if (!diakritika.value) {
        for (let i = 0; i < arr.length; i++) {
            arr[i] = arr[i].normalize("NFD").replace(/[\u0300-\u036f]/g, "")
        }
    } else if (!velkaPismena.value) {
        for (let i = 0; i < arr.length; i++) {
            arr[i] = arr[i].toLocaleLowerCase()
        }
    } else {
        textovePole.value!.text = puvodniText.value
        return
    }

    textovePole.value!.text = arr.join()
}

</script>
<template>
    <div id="pulic">
        <div id="nastaveni">
            <h2>Nastavení</h2>

            <div id="moznosti">

                <div id="delka">
                    <button :class="{ aktivni: 1 == delka }" @click="d(1)">1min</button>
                    <button :class="{ aktivni: 2 == delka }" @click="d(2)">2min</button>
                    <button :class="{ aktivni: 3 == delka }" @click="d(3)">3min</button>
                    <button :class="{ aktivni: 5 == delka }" @click="d(5)">5min</button>
                    <button :class="{ aktivni: 10 == delka }" @click="d(10)">10min</button>
                    <button :class="{ aktivni: 15 == delka }" @click="d(15)">15min</button>
                    <button :class="{ aktivni: 20 == delka }" @click="d(20)">20min</button>
                    <button :class="{ aktivni: 30 == delka }" @click="d(30)">30min</button>
                    <button :class="{ aktivni: -1 == delka }" @click="d(-1)">Do konce</button>
                </div>

                <hr id="predel">

                <div class="kontejner">
                    <button>Velká písmena</button>
                    <button>Diakritika</button>
                </div>
            </div>
        </div>

        <div id="text">
            <div>
                <div>{{ textovePole?.text.length }} znaků</div>
                <button class="tlacitko" id="genBtn" @click="getText()">Generovat</button>
            </div>

            <TextZadani ref="textovePole" />
        </div>
    </div>
</template>
<style scoped>
#genBtn {
    margin-top: 0;
    width: 120px;
    background-color: var(--tmave-fialova);
}

#genBtn:hover {
    background-color: var(--fialova);
}

#pulic {
    display: flex;
    justify-content: space-between;
    gap: 40px;
    width: 860px;
}

#text {
    width: 410px;
    display: flex;
    flex-direction: column;
    height: 400px;
    border-radius: 10px;
    gap: 15px;
}

#text>div {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

#text>div>div {
    width: 120px;
    height: 40px;
    padding: 10px;
    background-color: var(--tmave-fialova);
    border-radius: 5px;
}

#nastaveni {
    width: 410px;
    display: flex;
    gap: 10px;
    flex-direction: column;
    height: 400px;
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    border-radius: 10px;
}

#moznosti {
    display: flex;
    height: 100%;
}

.kontejner {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
    margin: 0 10px;
    cursor: pointer;
    transition: filter 0.2s;
}

#delka {
    display: flex;
    gap: 10px;
    flex-direction: column;
    justify-content: center;
    flex-wrap: wrap;
    padding: 5px;
    margin-right: 15px;
    width: 94px;
}

#predel {
    margin: 12px 0 12px 0;
    border: 1px solid var(--fialova);
}
</style>