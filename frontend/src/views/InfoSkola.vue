<script setup lang="ts">
import { useHead } from "unhead";
import { ref } from "vue";
import { getToken, pridatOznameni } from "../utils";
import axios from "axios";
import { prihlasen, role } from "../stores";

useHead({
    title: "Systém pro školy"
})

const email = ref("")
const telefon = ref("+420")
const skola = ref("")

const odeslano = ref(false)
const odesilame = ref(false)

function potvrdit(e: Event) {
    e.preventDefault()
    if (chekujUdaje() == false || role.value == "student") return
    if (email.value.length == 0 || telefon.value.length <= 4 || skola.value.length == 0) {
        pridatOznameni("Vyplň prosím všechna pole!")
        return
    }
    odesilame.value = true
    axios.post("/skola/zapis-skoly", {
        "jmeno_skoly": skola.value,
        "kontaktni_email": email.value,
        "kontaktni_telefon": telefon.value
    }, {
        headers: {
            Authorization: `Bearer ${getToken()}`
        }
    }).then(() => {
        odeslano.value = true
        role.value = "ucitel"
    }).catch(e => {
        pridatOznameni(e.response.data.error)
    }).finally(() => {
        odesilame.value = false
    })
}

function chekujUdaje() {
    if (email.value && !/^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$/g.test(email.value)) pridatOznameni("Špatný formát e-mailu.")
    else if (telefon.value != "+420" && !/^\+[0-9]{6,15}$/g.test(telefon.value)) pridatOznameni("Telefon musí být ve formátu: +420123456789 <br>(předčíslí a číslo nez mezer)")
    else return true
    return false
}

</script>

<template>
    <h1>Systém pro školy</h1>

    <div class="bloky">
        <video controls controlslist="nodownload" width="100%">
            <source src="/skolni_system.mp4">
            Váš prohlížeč neumí přehrát toto video.
            <a href="/skolni_system.mp4" download="/skolni_system.mp4">Tady</a> si ho můžete zkusit stáhnout.
        </video>
    </div>
    <div class="specialni-blok bloky">
        <div class="obrazek">
            <img src="../assets/seznamZaku.svg" alt="Seznam žáků">
        </div>
        <div class="text">
            <h2>Řazení žáků do tříd</h2>
            <div>
                <p>
                    Studenti se do tříd připojují pomocí <b>4-místného kódu</b>, který se s třídou automaticky vytvoří. <br>
                    Vyučujícímu jsou ke každému žákovi dostupné <b>statistiky</b> a výsledky jeho předchozích prací. <br>
                    Studenty lze také <b>přesouvat mezi třídami</b>, nebo je z tříd odebírat.
                </p>
            </div>
        </div>
    </div>
    <div class="specialni-blok bloky">
        <div class="text">
            <h2>Zadávání prací</h2>
            <div>
                <p>
                    Práce jsou způsob, jak žákům zadat <b>specifický text</b> podle vašeho uvážení.
                    Je tedy možné zadat nejen automaticky <b>vygenerovaný text</b> (z nabídky), ale i <b>libovolný text</b> (třeba zkopírovaný z
                    internetu).
                    Možnostem se meze nekladou.
                </p>
            </div>
        </div>
        <div class="obrazek">
            <img src="../assets/praceZadani.svg" alt="Textové pole na zadání práce">
        </div>
    </div>
    <div class="bloky">
        <div class="text">
            <h2>Cena, platba a zkušební období</h2>
            <p>
                <b>Zdarma</b> můžete po dobu <b>tří měsíců</b> vyzkoušet systém s libovolným počtem žáků.
                Po jeho uplynutí je cena na žáka jednorázová. Od <b>99Kč</b> za nového žáka.
                Fakturace probíhá <b>jednou ročně bankovním převodem</b>.
            </p>
        </div>
    </div>
    <div class="bloky" id="formular">
        <h2>Formulář pro zařazení školy</h2>
        <hr id="predel">
        <div class="container">
            <div id="ilustrace">
                <img src="../assets/pavoukSkola.svg" alt="Pavouk a škola">
                <p>Po odeslání tohoto formuláře si budete moct systém okamžitě vyzkoušet!</p>
            </div>
            <div v-if="!odeslano && !prihlasen">
                <p>
                    Před vyplněním formuláře se prosím přihlašte, a nebo pokud účet ještě nemáte, založte si ho prosím.
                </p>
            </div>
            <form v-if="!odeslano && prihlasen" @submit="potvrdit">
                <div>
                    <label for="skola">Jméno školy</label>
                    <input type="text" id="skola" placeholder="Např: Gymnázium pana Pavouka" v-model="skola">
                </div>
                <div>
                    <h2>Kontaktní osoba</h2>
                    <label for="email">E-mail</label>
                    <input type="email" id="email" placeholder="Např: pavoukova@jakopavouk.cz" v-model="email">
                    <label for="tel">Telefonní číslo</label>
                    <input type="tel" id="tel" placeholder="Např: +420123456789" v-model="telefon">
                </div>
                <button class="tlacitko" type="submit" :disabled="odesilame">Odeslat</button>
            </form>
            <div v-if="odeslano">
                <p>
                    Díky za registraci školy! <br> Vlevo v navigaci se vám objevila záložka jménem "Škola". Začněte tam!
                </p>
            </div>
        </div>
    </div>
</template>

<style scoped>
.specialni-blok {
    display: flex;
    gap: 6pt;
    margin-top: 50px;
    background-color: transparent;
    align-items: normal;
    padding: 0;
}

.specialni-blok>.obrazek {
    align-self: center;
    display: flex;
    height: 85%;
    width: 80%;
}

.specialni-blok>.obrazek:first-of-type {
    overflow: hidden !important;
    width: auto;
}

.specialni-blok>.obrazek:first-of-type>img {
    width: 150%;
    height: 100%;
}

.specialni-blok>.obrazek>img {
    width: 100%;
}

.specialni-blok>.text {
    display: flex;
    flex-direction: column;
    justify-content: center;
    background-color: var(--tmave-fialova);
    padding: 25px 15px;
    border-radius: 10px;
}

.text {
    width: 100%;
}

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

#formular>div {
    display: flex;
    align-items: center;
}

#formular>div #ilustrace {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5em;
}

#formular>div #ilustrace>img {
    width: 240px;
}

#formular>div>div:not(#ilustrace) {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 320px;
}

form {
    display: flex;
    flex-direction: column;
    width: 100%;
    align-items: center;
    height: 320px;
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
    .container {
        flex-direction: column;
        gap: 2em;
    }

    #formular>div #ilustrace>img {
        display: none;
    }

    #formular>div #ilustrace>p {
        max-width: 90%;
    }

    form input {
        width: 100%;
    }

    form input:focus {
        width: 101%;
    }

    form input {
        width: 97%;
    }

    form>div {
        width: 60%;
    }

    .obrazek {
        overflow: visible;
        width: 80% !important;
        max-width: 400px;
    }

    .obrazek>img {
        max-width: none;
        width: 100% !important;
    }

    .specialni-blok {
        flex-direction: column-reverse !important;
    }

    div:nth-of-type(odd).specialni-blok {
        flex-direction: column !important;
    }

    .specialni-blok {
        margin-top: 25px;
    }
}

@media screen and (max-width: 800px) {
    form input {
        width: 98%;
    }

    form>div {
        width: 85%;
    }

    .bloky:first-of-type {
        padding: 10px;
    }
}
</style>