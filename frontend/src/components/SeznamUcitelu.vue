<script setup lang="ts">
import axios from 'axios';
import { getToken, pridatOznameni } from '../utils';
import { onMounted, ref, useTemplateRef } from 'vue';

const ucitele = ref([] as { email: string, pocet_trid: number, pocet_zaku: number }[])
const email = ref("")
onMounted(get)

function get() {
    axios.get("/skola/ucitele", {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(response => {
        ucitele.value = response.data
    }).catch(e => {
        console.log(e)
        pridatOznameni("Chyba serveru")
    })
}

const dialog1 = useTemplateRef("dialog1")
function otevritDialog(e: Event) {
    e.preventDefault()
    
    let button = e.currentTarget as Element
    let mailos = button.getAttribute("data-email")
    if (mailos === null) return
    email.value = mailos

    dialog1.value?.showModal()
}

function zavritDialog(e: Event) {
    e.preventDefault()
    dialog1.value?.close()
}

function smazatUcitele(e: Event) {
    axios.post("/skola/upravit-ucitele", { email: email.value, akce: "smazat" }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        get()
    }).catch(e => {
        console.log(e)
        pridatOznameni(e)
    }).finally(() => {
        zavritDialog(e)
    })
}
</script>
<template>
    <div id="seznam">
        <div v-for="ucitel in ucitele" :key="ucitel.email">
            <div class="blok">
                <h3>{{ ucitel.email }}</h3>
                <div class="stats">
                    <span>
                        <b>{{ ucitel.pocet_trid }}</b>
                        {{ (ucitel.pocet_trid == 0 || ucitel.pocet_trid > 4) ? "tříd" : (ucitel.pocet_trid == 1) ? "třída" : "třídy" }}
                    </span>
                    <span>
                        <b>{{ ucitel.pocet_zaku }}</b>
                        {{ (ucitel.pocet_zaku == 0 || ucitel.pocet_zaku > 4) ? "studentů" : (ucitel.pocet_zaku == 1) ? "student" : "studenti" }}
                    </span>
                </div>
            </div>
            <button @click="otevritDialog" :data-email="ucitel.email">
                <img src="../assets/icony/trash.svg" alt="Smazat">
            </button>
        </div>
    </div>

    <dialog ref="dialog1">
        <div id="dialog-kontejner">
            <h2>Opravdu chcete odstranit tohoto učitele?</h2>
            <h3></h3>
            <div>
                <button class="cervene-tlacitko" @click="smazatUcitele">Smazat</button>
                <button class="tlacitko" @click="zavritDialog">Zrušit</button>
            </div>
        </div>
    </dialog>
</template>
<style scoped>
#seznam {
    display: flex;
    flex-direction: column;
    gap: 5px;
    width: 500px;
}

#seznam>div {
    display: flex;
    align-items: center;
    gap: 5px;
}

#seznam>div>button {
    background-color: var(--tmave-fialova);
    border: none;
    border-radius: 10px;
    cursor: pointer;
    padding: 10px;
    width: 45px;
    height: 45px;

    display: flex;
    justify-content: center;
    align-items: center;
}

.blok {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    width: 460px;
    height: 55px;
    padding: 0.2em 0.5em 0.2em 1em;

    display: flex;
    align-items: center;
    gap: 7px;
}

.blok>h3 {
    width: 312px;
    text-align: left;
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
}

.blok .stats {
    display: flex;
    flex-direction: column;
    text-align: left;
}
</style>