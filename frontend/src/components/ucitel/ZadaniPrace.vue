<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from 'vue';
import TextZadani from './TextZadani.vue';
import axios from 'axios';
import { getToken, pridatOznameni, format } from '../../utils';
import Tooltip from '../../components/Tooltip.vue';
import HodnoticiTabulka, { Tabulka } from './HodnoticiTabulka.vue';

const props = defineProps<{
    tridaID: number;
    posledniRychlost: number;
    hodnoticiTabulka: Tabulka;
}>();

const emit = defineEmits(['zadano']);

const textovePole = useTemplateRef('textove-pole');

const delka = ref(5 * 60);
const typTextu = ref('');
const lekceTextu = ref();
const hodnocena = ref(false);

const texty = ref([] as { jmeno: string; obtiznost: number }[]);
const lekce = ref([] as { id: number; lekce_id: number; pismena: string }[]);
let mapa: Map<string, { id: number; lekce_id: number; pismena: string }[]> = new Map<string, { id: number; lekce_id: number; pismena: string }[]>();

onMounted(() => {
    axios
        .get('/procvic')
        .then((response) => {
            response.data.texty.forEach((el: { jmeno: string; obtiznost: number }) => {
                texty.value.push({ jmeno: el.jmeno, obtiznost: el.obtiznost });
            });
            texty.value.sort((a: { obtiznost: number }, b: { obtiznost: number }) => {
                return a.obtiznost - b.obtiznost;
            });
        })
        .catch(() => {
            pridatOznameni('Chyba serveru');
        });

    axios
        .get('/skola/typy-cviceni', { params: { trida_id: props.tridaID } })
        .then((response) => {
            for (const k in response.data) {
                mapa.set(
                    k,
                    response.data[k].sort((a: { id: number; lekce_id: number; pismena: string }, b: { id: number; lekce_id: number; pismena: string }) => a.lekce_id - b.lekce_id),
                );
            }
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
});

function getText() {
    if (typTextu.value == '') return;

    axios
        .post(
            '/skola/text',
            {
                typ: typTextu.value,
                z_lekce: lekceTextu.value.pismena,
                trida_id: props.tridaID,
                delka: odhadovanaDelkaTextu.value == -1 ? 239 : odhadovanaDelkaTextu.value,
            },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then((response) => {
            if (textovePole.value?.text.length != 0 && textovePole.value?.text[textovePole.value?.text.length - 1] != ' ') textovePole.value!.text += ' ';
            textovePole.value!.text += response.data.text;
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

function pridatPraci() {
    if (textovePole.value!.text.length <= 10) {
        pridatOznameni('Není text nějak krátký?');
        return;
    }

    if (!textovePole.value!.ready) {
        pridatOznameni('Text obsahuje červeně označené znaky, které nelze jednoduše napsat na české klávesnici.');
        return;
    }

    axios
        .post(
            '/skola/pridat-praci',
            {
                cas: delka.value,
                trida_id: props.tridaID,
                text: textovePole.value!.text,
                hodnocena: hodnocena.value,
            },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            emit('zadano');
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

function d(x: number) {
    delka.value = x;
}

const puvodniText = ref('');

function smazatDiakritiku() {
    if (textovePole.value!.text.length == 0) return;
    puvodniText.value = textovePole.value!.text;
    textovePole.value!.text = textovePole.value!.text.normalize('NFD').replace(/[\u0300-\u036f]/g, '');
    if (puvodniText.value === textovePole.value!.text) {
        puvodniText.value = '';
    }
}

function smazatVelkaPismena() {
    if (textovePole.value!.text.length == 0) return;
    puvodniText.value = textovePole.value!.text;
    textovePole.value!.text = textovePole.value!.text.toLocaleLowerCase();
    if (puvodniText.value === textovePole.value!.text) {
        puvodniText.value = '';
    }
}

function smazatEnterAMezery() {
    if (textovePole.value!.text.length == 0) return;
    puvodniText.value = textovePole.value!.text;
    textovePole.value!.text = textovePole.value!.text.replace(/\n/g, ' ').replace(/ {2,}/g, ' ').trim();
    if (puvodniText.value === textovePole.value!.text) {
        puvodniText.value = '';
    }
}

function vymenitUvozovky() {
    if (textovePole.value!.text.length == 0) return;
    puvodniText.value = textovePole.value!.text;
    textovePole.value!.text = textovePole.value!.text.replace(/[„“”‟❞❝＂⹂]/g, '"').replace(/[‚’‘‛❛❜]/g, "'");
    if (puvodniText.value === textovePole.value!.text) {
        puvodniText.value = '';
    }
}

function zrusitPosledniUpravu() {
    textovePole.value!.text = puvodniText.value;
    puvodniText.value = '';
}

function getZnakyASlova() {
    if (textovePole.value == null) return '0 / 0';
    let vys = '';

    vys += textovePole.value.text.length;
    if (textovePole.value.text.length == 0 || textovePole.value.text.length >= 5) vys += ' znaků';
    else if (textovePole.value.text.length == 1) vys += ' znak';
    else vys += ' znaky';

    vys += ' / ';
    let slova = 0;
    if (textovePole.value.text.length != 0) slova = textovePole.value.text.trim().split(/[ \n]+/).length;
    vys += slova;
    if (slova == 0 || slova >= 5) vys += ' slov';
    else if (slova == 1) vys += ' slovo';
    else vys += ' slova';

    return vys;
}

function upravaSelectuLekci() {
    lekce.value = [];
    switch (typTextu.value) {
        case '1':
            lekce.value = mapa.get('nova')!;
            break;
        case '2':
            lekce.value = mapa.get('naucena')!;
            break;
        case '3':
            lekce.value = mapa.get('slova')!;
            break;
        case '4':
            lekce.value = mapa.get('programator')!;
            break;
        default:
            lekce.value = [];
    }

    lekceTextu.value = '';
}

const odhadovanaDelkaTextu = computed(() => {
    if (props.posledniRychlost == -1) return -1;
    return Math.ceil((props.posledniRychlost + 10) * (delka.value / 60));
});
</script>
<template>
    <div id="pulic-zadani">
        <div id="leva-pulka">
            <div id="parametry">
                <h2>Nastavení</h2>

                <div id="moznosti">
                    <div id="delka">
                        <Tooltip zprava="Pokud student dopíše text před vypršením časového limitu, bude ho psát znovu." :sirka="210" :vzdalenost="0">
                            <h3>Čas</h3>
                        </Tooltip>
                        <hr id="predel2" />
                        <button :class="{ aktivni: 60 == delka }" @click="d(60)">1min</button>
                        <button :class="{ aktivni: 120 == delka }" @click="d(120)">2min</button>
                        <button :class="{ aktivni: 180 == delka }" @click="d(180)">3min</button>
                        <button :class="{ aktivni: 300 == delka }" @click="d(300)">5min</button>
                        <button :class="{ aktivni: 600 == delka }" @click="d(600)">10min</button>
                        <button :class="{ aktivni: 900 == delka }" @click="d(900)">15min</button>
                        <button :class="{ aktivni: 1200 == delka }" @click="d(1200)">20min</button>
                        <button :class="{ aktivni: 1800 == delka }" @click="d(1800)">30min</button>
                    </div>

                    <hr id="predel" />

                    <div class="vert-kontejner">
                        <div class="kontejner">
                            <button @click="smazatVelkaPismena" class="tlacitko">Smazat velká písmena</button>
                            <button @click="smazatDiakritiku" class="tlacitko">Smazat diakritiku</button>
                            <button @click="smazatEnterAMezery" class="tlacitko">Smazat mezery navíc</button>
                            <button @click="vymenitUvozovky" class="tlacitko">Vyměnit uvozovky</button>
                        </div>

                        <div class="kontejner">
                            <button @click="zrusitPosledniUpravu" class="cervene-tlacitko" :disabled="puvodniText.length == 0">Zrušit poslední úpravu</button>
                        </div>

                        <button @click="pridatPraci" class="tlacitko">Zadat práci</button>
                    </div>
                </div>
            </div>

            <div>
                <label for="toggle" class="kontejner">
                    Hodnocená práce:
                    <input v-model="hodnocena" type="checkbox" id="toggle" class="radio" />
                </label>
                <HodnoticiTabulka :stavajiciRychlosti="props.hodnoticiTabulka" :tridaID="props.tridaID" v-show="hodnocena" />
            </div>
        </div>

        <div id="text">
            <div style="display: flex; gap: 5px; width: 100%">
                <select v-model="typTextu" @change="upravaSelectuLekci">
                    <option value="" style="color: #a1a1a1">Generovat text</option>
                    <option value="1">Nová písmena z lekce</option>
                    <option value="2">Naučená písmena z lekce</option>
                    <option value="3">Slova z lekce</option>
                    <option value="4">Programátorské</option>
                    <option v-for="(t, i) in texty" :value="t.jmeno" :class="{ lehkaObtiznost: t.obtiznost == 1, stredniObtiznost: t.obtiznost == 2, tezkaObtiznost: t.obtiznost == 3 }" :key="i">
                        {{ t.jmeno }}
                    </option>
                </select>
                <select v-model="lekceTextu" :disabled="typTextu != '1' && typTextu != '2' && typTextu != '3' && typTextu != '4'" style="width: 121px">
                    <option value="" style="color: #a1a1a1">Vybrat lekci</option>
                    <option v-for="l in lekce" :value="l" :key="l.pismena">{{ format(l.pismena) }}</option>
                </select>
                <button class="tlacitko" @click="getText" :disabled="typTextu == '' || ((typTextu == '1' || typTextu == '2' || typTextu == '3' || typTextu == '4') && lekceTextu == '')">
                    <img src="../../assets/icony/plus.svg" alt="Prodloužit" />
                </button>
            </div>

            <TextZadani ref="textove-pole" />

            <div>
                <span>{{ getZnakyASlova() }}</span>
                <Tooltip
                    v-show="odhadovanaDelkaTextu != -1"
                    :zprava="`Odhad, jak má být text dlouhý, aby ho studenti nestihli napsat až do konce. Počítá s rychlostí z poslední práce.`"
                    :sirka="350"
                    :vzdalenost="-190"
                    :vzdalenostX="120"
                >
                    <span>~ {{ odhadovanaDelkaTextu }} znaků</span>
                </Tooltip>
            </div>
        </div>
    </div>
</template>
<style scoped>
#leva-pulka {
    display: flex;
    flex-direction: column;
    gap: 15px;
    padding-bottom: 30px;
    justify-content: stretch;
}

#leva-pulka div:has(#hodnotici-tabulka) {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    background-color: var(--tmave-fialova);
    padding: 15px 10px;
    border-radius: 10px;
    flex-grow: 3;
}

.lehkaObtiznost {
    color: rgb(0, 185, 0);
}

.stredniObtiznost {
    color: rgb(255, 184, 52);
}

.tezkaObtiznost {
    color: rgb(255, 51, 51);
}

.kontejner .tlacitko,
.cervene-tlacitko {
    width: 225px;
}

.vert-kontejner {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-around;
}

select {
    height: 40px;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    color: white;
    background-color: var(--tmave-fialova);
    cursor: pointer;
    transition: 0.2s;
    font-family: inherit;
}

/* firefox nenenene */
@supports (-webkit-tap-highlight-color: black) {
    select:hover {
        background-color: var(--fialova) !important;
    }

    select {
        padding-left: 5px;
    }
}

select:hover {
    background-color: var(--tmave-fialova);
}

select option {
    background-color: var(--tmave-fialova) !important;
    font-family: 'Montserrat', Candara !important;
    font-weight: 400;
}

select option:disabled {
    color: white !important;
    opacity: 1;
}

#pulic-zadani {
    display: flex;
    justify-content: center;
    gap: 15px;
    width: 860px;
    margin-top: 20px;
    margin-bottom: -50px;
}

#text {
    position: relative;
    width: 410px;
    display: flex;
    flex-direction: column;
    height: calc(100vh - 60px - 25px - 30px - 5px);
    /* celá obrazovka - všechno co je nad tím */
    border-radius: 10px;
    gap: 15px;
}

#text > div {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

#text span {
    align-self: flex-start;
    display: flex;
    padding: 8px;
    height: 10px;
    margin: 0;
    position: relative;
    top: -10px;
}

#text .tlacitko {
    width: 40px;
    margin: 0;
    background-color: var(--tmave-fialova);
    display: flex;
    align-items: center;
    justify-content: center;
}

#text .tlacitko:hover {
    background-color: var(--fialova);
}

#text .tlacitko img {
    width: 25px;
}

#parametry {
    width: 410px;
    display: flex;
    gap: 10px;
    flex-direction: column;
    height: 380px;
    background-color: var(--tmave-fialova);
    padding: 15px 15px 20px 15px;
    border-radius: 10px;
}

#moznosti {
    display: flex;
    height: 100%;
}

.kontejner {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    transition: filter 0.2s;
}

.kontejner:has(input[id='toggle']) {
    flex-direction: row;
    font-size: 1.2em;
    gap: 1em;
}

.kontejner > div {
    display: flex;
    gap: 15px;
}

#delka {
    display: flex;
    gap: 10px;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 100px !important;
}

#delka > button {
    background-color: transparent;
    border: none;
    color: rgba(240, 240, 240, 0.7);
    font-size: 16px;
    cursor: pointer;
    font-weight: 400;
    transition: 0.1s;
}

#delka button:hover {
    font-weight: 600;
    color: var(--bila);
    cursor: pointer;
}

.aktivni {
    color: var(--svetle-fialova) !important;
    font-weight: 600 !important;
}

#delka h3 {
    position: relative;
    bottom: 15px;
}

#predel {
    margin: 12px 20px 12px 15px;
    border: 1px solid var(--fialova);
}

#predel2 {
    margin-top: -22px;
    margin-bottom: 8px;
    width: 70%;
    border: 1px solid rgb(206, 206, 206);
    align-self: center;
}

@media screen and (max-width: 900px) {
    #pulic-zadani {
        flex-direction: column;
        width: auto;
    }
}
</style>
