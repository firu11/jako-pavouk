<script setup lang="ts">
import { onMounted, ref, computed, useTemplateRef } from 'vue';
import { getToken, pridatOznameni } from '../../utils';
import { moznostiRocnik, moznostiTrida, moznostiSkupina } from '../../stores';
import axios from 'axios';
import { useRouter } from 'vue-router';
import KodTridy from '../../components/KodTridy.vue';
import HodnoticiTabulka, { Tabulka } from './HodnoticiTabulka.vue';

const emit = defineEmits(['prejmenovatTridu', 'refresh']);

type Trida = {
    id: number;
    jmeno: string;
    ucitelID: number;
    kod: string;
    zamknuta: boolean;
    klavesnice: string;
    hodnoceni: Tabulka;
};

const props = defineProps<{
    trida: Trida;
    pocetStudentu: number;
}>();

const router = useRouter();

const tridaRocnikUprava = ref();
const tridaJmenoUprava = ref();
const tridaSkupinaUprava = ref();
const klavesniceUprava = ref();

let smazatZprava = ref('');

function smazatTridu(e: Event) {
    e.preventDefault();

    if (props.pocetStudentu != 0) {
        if (props.pocetStudentu == 1) smazatZprava.value = '1 student už je v této třídě.';
        else if (props.pocetStudentu <= 4) smazatZprava.value = `${props.pocetStudentu} studenti už jsou v této třídě.`;
        else smazatZprava.value = `${props.pocetStudentu} studentů už je v této třídě.`;
        dialog1.value?.showModal();
    } else {
        postSmazat(e);
    }
}

function zavritDialog(e: Event) {
    e.preventDefault();
    dialog1.value?.close();
}

function postSmazat(e: Event) {
    e.preventDefault();

    axios
        .post(
            '/skola/zmena-tridy',
            { trida_id: props.trida.id, zmena: 'smazat' },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            router.push('/skola');
            pridatOznameni(`Třída "${props.trida.jmeno}", byla smazána.`);
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

function postKlavesnice(e: Event) {
    e.preventDefault();

    axios
        .post(
            '/skola/zmena-tridy',
            { trida_id: props.trida.id, zmena: 'klavesnice', hodnota: klavesniceUprava.value ? 'qwerty' : 'qwertz' },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            emit('refresh');
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

onMounted(() => {
    let a = props.trida.jmeno.split(/[. ]/);
    tridaJmenoUprava.value = a[1];
    tridaRocnikUprava.value = a[0] + (isNaN(+a[0]) ? ' ' : '.');
    tridaSkupinaUprava.value = a[3] == undefined ? '-' : a[3];

    klavesniceUprava.value = props.trida.klavesnice == 'qwerty';
});

const tridaJmeno = computed(() => {
    return `${tridaRocnikUprava.value}${tridaJmenoUprava.value}${tridaSkupinaUprava.value != '-' ? ' ￨ ' + tridaSkupinaUprava.value : ''}`;
});

const dialog1 = useTemplateRef('dialog1');
</script>
<template>
    <div id="nastaveni">
        <KodTridy :kod="trida.kod" :id="trida.id" :zamknuta="trida.zamknuta" />

        <form>
            <div>
                <h3>Ročník:</h3>
                <select v-model="tridaRocnikUprava">
                    <option v-for="v in moznostiRocnik" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <div>
                <h3>Písmeno:</h3>
                <select v-model="tridaJmenoUprava">
                    <option v-for="v in moznostiTrida" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <div>
                <h3>Skupina:</h3>
                <select v-model="tridaSkupinaUprava">
                    <option v-for="v in moznostiSkupina" :value="v" :key="v">{{ v }}</option>
                </select>
            </div>

            <button class="tlacitko" @click="emit('prejmenovatTridu', $event, tridaJmeno)" :disabled="tridaJmeno == trida.jmeno">Potvrdit</button>

            <span>Podle jména se třídy řadí do ročníků v seznamu tříd.</span>
            <span>
                Skupina je určena pro rozdělení třídy. To se může hodit, pokud vyučujete třídy po menších skupinkách.
                <br />
                Např.: <b>3.B&nbsp;￨&nbsp;1</b> a <b>3.B&nbsp;￨&nbsp;2</b>.
            </span>
        </form>

        <form>
            <h3 style="padding-bottom: 10px">Hodnotící tabulka</h3>
            <HodnoticiTabulka :stavajiciRychlosti="props.trida.hodnoceni" :tridaID="trida.id" />
            <span>U prácí označených jako hodnocené budou žáci automaticky ohodnoceni podle této tabulky.</span>
        </form>

        <form>
            <div>
                <h3>Preference klávesnice:</h3>

                <input type="checkbox" id="toggle1" class="toggle-checkbox" v-model="klavesniceUprava" @change="postKlavesnice" />
                <label for="toggle1" class="toggle-contejner">
                    <div>Qwertz</div>
                    <div>Qwerty</div>
                </label>
            </div>

            <span>
                Klávesnice ovlivňuje pouze to, jak se generují texty typu: <b>Naučená písmena</b> a <b>Slova</b>, v lekcích <b>t,z - y,x,m</b>. Tam se totiž množina již naučených písmen liší
                klávesami Z a Y.
            </span>
            <span>Žákům, kteří mají jinou klácesnici než jaká je nastavena v jejich třídě, se zobrazí varování.</span>
        </form>

        <form>
            <div>
                <h3>Smazání třídy:</h3>
                <button type="button" class="cervene-tlacitko" @click="smazatTridu">Smazat třídu</button>
            </div>
        </form>

        <dialog ref="dialog1">
            <div id="dialog-kontejner">
                <h2>Opravdu chcete smazat třídu?</h2>
                <div>
                    <button class="cervene-tlacitko" @click="postSmazat">Opravdu smazat</button>
                    <button class="tlacitko" @click="zavritDialog">Zrušit</button>
                </div>
                <span>{{ smazatZprava }}</span>
            </div>
        </dialog>
    </div>
</template>
<style scoped>
/* firefox nenenene */
@supports (-webkit-tap-highlight-color: black) {
    select:hover {
        background-color: var(--svetle-fialova) !important;
    }

    select {
        padding-left: 5px;
    }
}

select:hover {
    background-color: var(--fialova);
}

select option {
    background-color: var(--fialova) !important;
    font-weight: 400;
}

#dialog-kontejner {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 1em;
}

#dialog-kontejner > span {
    color: var(--seda);
}

#dialog-kontejner > div {
    display: flex;
    gap: 1em;
}

dialog {
    width: 410px;
    height: 140px;
    margin-left: -205px;
    margin-top: -70px;
}
</style>
