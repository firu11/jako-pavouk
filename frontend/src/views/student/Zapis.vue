<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getToken, pridatOznameni } from '../../utils';
import axios from 'axios';
import { prihlasen, role } from '../../stores';

const router = useRouter();

const kodURL = useRoute().params.kod as string;
const state = ref('kod'); // kod, jmeno

const kod = ref('');
const jmeno = ref('');

onMounted(() => {
    if (kodURL != undefined) {
        kod.value = kodURL;
    }
});

function potvrditKod(e: Event) {
    e.preventDefault();
    if (getToken() == null) {
        pridatOznameni('Nejsi přihlášen/a!');
        return;
    }
    if (!/^[a-zA-Z]{4}$/.test(kod.value)) {
        pridatOznameni('Kód je neplatný');
        return;
    }
    axios
        .get('/skola/test-tridy/' + kod.value, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then(() => {
            state.value = 'jmeno';
        })
        .catch((e) => {
            console.log(e.response.data.error);
            if (e.response.data.error == 'Takova trida neexistuje') {
                pridatOznameni('Taková třída neexistuje');
                return;
            }
            if (e.response.data.error == 'Trida je zamcena') {
                pridatOznameni('Tato třída je zamčená');
                return;
            }
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

function zapsatSe(e: Event) {
    e.preventDefault();
    if (jmeno.value == '') {
        pridatOznameni('Zadej prosím jméno');
        return;
    }
    axios
        .post(
            '/skola/zapis',
            { kod: kod.value, jmeno: jmeno.value },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            router.push('/trida');
            role.value = 'student';
        })
        .catch((e) => {
            if (e.response.data.error == 'Uz jsi ve tride') {
                pridatOznameni('Už jsi ve třídě');
                return;
            }
            if (e.response.data.error == 'jako ucitel nemuzete byt ve tride') {
                pridatOznameni('Jako učitel/ka se nemůžete připojit do žádné třídy');
                return;
            }
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

function checkPrihlaseni() {
    if (!prihlasen.value) pridatOznameni('Nejsi přihlášen/a!');
}
</script>
<template>
    <h1>Zápis do třídy</h1>
    <div id="kontejner">
        <img src="../../assets/pavoukSkola.svg" alt="Pavouk před školou" width="250" height="175" />
        <form v-if="state == 'kod'">
            <h2>Zadej kód:</h2>
            <input type="text" v-model="kod" placeholder="ABCD" @focus="checkPrihlaseni" />
            <button class="tlacitko" @click="potvrditKod">Potvrdit</button>
        </form>
        <form v-else>
            <h2>Zadej své jméno:</h2>
            <input id="jmeno" type="text" v-model="jmeno" placeholder="Novák Honza" />
            <p>(Nejlépe příjmení první)</p>
            <button class="tlacitko" @click="zapsatSe">Zapsat se</button>
        </form>
    </div>
    <RouterLink to="/skolni-system">Jsem učitel/ka. </RouterLink>
</template>
<style scoped>
/* eslint-disable-next-line vue-scoped-css/no-unused-selector */
a {
    color: white;
    margin-top: 16px;
    font-size: 1.1em;
}

.tlacitko {
    margin: 5px;
    width: 120px;
    align-self: center;
}

h2 {
    padding-left: 5px;
    font-size: 21px;
}

#kontejner {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    display: flex;
    padding: 40px 25px 20px 25px;
    gap: 10px;
    width: 550px;
    justify-content: space-around;
}

#kontejner img {
    width: 250px;
}

form {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

form p {
    padding-left: 5px;
    font-size: 14px;
    margin-bottom: 8px;
}

input {
    width: 190px;
    height: 50px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 30px;
    transition: all 0.15s cubic-bezier(0.5, 0, 0.5, 1) 0s;
    margin: 5px;
    text-align: center;
}

#jmeno {
    height: 40px;
    font-size: 16px;
}

input:focus {
    outline: none !important;
    width: 200px;
    margin: 5px 0;
}

form input::placeholder {
    color: rgba(255, 255, 255, 0.42);
}

@media screen and (max-width: 1100px) {
    #kontejner > img {
        display: none;
    }

    #kontejner {
        padding: 20px 25px;
        width: auto;
        margin: 0 30px;
    }

    h1 {
        max-width: 60vw;
    }
}
</style>
