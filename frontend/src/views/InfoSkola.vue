<script setup lang="ts">
import { useHead } from "unhead";
import { ref } from "vue";
import { pridatOznameni } from "../utils";

useHead({
    title: "Systém pro školy"
})

const jmeno = ref("")
const prijmeni = ref("")
const email = ref("")
const telefon = ref("+420")
const skola = ref("")

const odeslano = ref(false)

function potvrdit(e: Event) {
    e.preventDefault()
    if (chekujUdaje("") == false) return
    if (jmeno.value.length == 0 || prijmeni.value.length == 0 || email.value.length == 0 || telefon.value.length <= 4 || skola.value.length == 0) {
        pridatOznameni("Vyplň prosím všechna pole!")
        return
    }
    odeslano.value = true
}

function chekujUdaje(jaky: string) {
    if ((jaky === "email" || jaky == "") && email.value && !/^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/g.test(email.value)) pridatOznameni("Špatný formát e-mailu.")
    else if ((jaky === "telefon" || jaky == "") && telefon.value != "+420" && !/^\+[0-9]{6,15}$/g.test(telefon.value)) pridatOznameni("Telefon musí být ve formátu: +420123456789 <br>(předčíslí a číslo nez mezer)")
    else return true
    return false
}

</script>

<template>
    <h1>Systém pro školy</h1>

    <div class="bloky">
        <video controls controlslist="nodownload" width="100%" poster="/OGnahled.png">
            <source src="../assets/video.mp4">

            Váš prohlížeč neumí přehrát toto video. <!-- to snad už ani nejde aby to někdo nepodporoval -->
            <a href="../assets/video.mp4" download="../assets/video.mp4">Tady</a> si ho můžete zkusit stáhnout.
        </video>
    </div>
    <div class="bloky">
        <div class="text">
            <h2>Co systém nabízí?</h2>
            <ul>
                <li>
                    Správa tříd
                    <ul>
                        <li>Řazení tříd podle ročníků</li>
                        <li></li>
                    </ul>
                </li>
                <li>
                    Zadávání prací
                    <ul>
                        <li>Řazení tříd podle ročníků</li>
                        <li></li>
                    </ul>
                </li>
            </ul>
        </div>
        <img src="../assets/pavoukStudent.svg" alt="Pavouk s klávesnicí" width="300" height="223">
    </div>
    <div class="bloky" id="formular">
        <h2>Formulář pro zařazení školy</h2>
        <hr id="predel">
        <form v-if="!odeslano">
            <div>
                <label for="skola">Jméno školy</label>
                <input type="text" id="skola" placeholder="Např: Gymnázium pana Pavouka" v-model="skola">
            </div>
            <div>
                <h2>Kontaktní osoba</h2>
                <label for="email">E-mail</label>
                <input type="email" id="email" placeholder="Např: pavoukova@jakopavouk.cz" v-model="email" @focusout="chekujUdaje('email')">
                <label for="tel">Telefonní číslo</label>
                <input type="tel" id="tel" placeholder="Např: +420123456789" v-model="telefon" @focusout="chekujUdaje('telefon')">
            </div>

            <button class="tlacitko" type="submit" @click="potvrdit">Odeslat</button>
        </form>
        <p v-else>
            Formulář odeslán! <br> Během pár hodin se vám ozveme.
        </p>
    </div>
</template>

<style scoped>
.bloky:first-of-type {
    flex-direction: column;
    height: auto;
    padding: 30px;
}

#predel {
    margin-bottom: 10px;
    width: 91%;
    border: 1px solid var(--fialova);
}

#formular {
    flex-direction: column;
    height: auto;
    gap: 8px;
}

form {
    display: flex;
    flex-direction: column;
    width: 100%;
    align-items: center;
}

form h2 {
    font-weight: 500;
    align-self: center;
    margin-bottom: 10px;
    margin-top: 5px;
    font-size: 19px;
}

form>div {
    display: flex;
    flex-direction: column;
}

form input {
    width: 300px;
    height: 30px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    transition: all 0.15s cubic-bezier(0.5, 0, 0.5, 1) 0s;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    margin: 0 5px 20px 5px;
}

form input:focus {
    outline: none !important;
    width: 310px;
    margin-right: 0;
    margin-left: 0;
    transition: all 0.15s cubic-bezier(0.5, 0, 0.5, 1) 0s;
}

form input::placeholder {
    color: rgba(255, 255, 255, 0.42);
}

form label {
    align-self: flex-start;
    font-size: 18px;
    margin: 0 5px 5px 5px;
    display: inline-block;
}

h1 {
    margin-bottom: 0;
}

@media screen and (max-width: 1100px) {
    form input {
        width: 100%;
    }

    form input:focus {
        width: 101%;
    }

    form input {
        width: 97%;
    }
}

@media screen and (max-width: 800px) {
    form input {
        width: 98%;
    }

    .bloky:first-of-type {
        padding: 10px;
    }
}
</style>