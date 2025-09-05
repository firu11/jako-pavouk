<script setup lang="ts">
import axios from 'axios';
import { onMounted, ref, computed, onUnmounted, useTemplateRef, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getToken, pridatOznameni, naJednoDesetiny } from '../../utils';
import SipkaZpet from '../../components/SipkaZpet.vue';
import ZadaniPrace from '../../components/ucitel/ZadaniPrace.vue';
import { useHead } from '@unhead/vue';
import Tooltip from '../../components/Tooltip.vue';
import NastaveniTridy from '../../components/ucitel/NastaveniTridy.vue';
import KodTridy from '../../components/KodTridy.vue';
import PrepinacTabu from '../../components/PrepinacTabu.vue';
import { mobil } from '../../stores';
import PraceBlok, { Prace, Zak } from '../../components/ucitel/PraceBlok.vue';
import { Tabulka } from '../../components/ucitel/HodnoticiTabulka.vue';

const id = useRoute().params.id;

const trida = ref({} as { id: number; jmeno: string; ucitelID: number; kod: string; zamknuta: boolean; klavesnice: string; hodnoceni: Tabulka });
const prace = ref([] as Prace[]);
const studentiVPraci = ref(new Map<number, Array<Zak>>());
const studenti = ref([] as { id: number; jmeno: string; email: string; cpm: number }[]);
const vsechnyTridy = ref([] as { id: number; jmeno: string; kod: string; zamknuta: boolean; pocet_studentu: number; klavesnice: string; pocet_praci: number }[]);

const prepinacTabu = useTemplateRef('prepinac-tabu');

const selectnutej = ref(-1);
const studentOznacenej = ref({ jmeno: '...', email: '...@...', dokonceno: 0, daystreak: 0, rychlost: -1, uspesnost: -1, klavesnice: 'QWERTZ', nejcastejsiChyby: new Map() });
const cpmVPracich = ref(new Map<number, number>());
const presnostVPracich = ref(new Map<number, number>());

const selectnutaPraceID = ref(-1);

const upravaStudenta = ref(false);
const jmenoUprava = ref();

const studentTridaZmena = ref();

const nacitamStudenta = ref(false);

const router = useRouter();

const copyPraciIndex = ref(-1);
const copyTrida = ref(0);

const nastaveni = ref();

let source: EventSource | null = null;

const titleJmenoTridy = computed(() => {
    return trida.value.jmeno;
});
useHead({
    title: titleJmenoTridy,
});

onMounted(() => {
    get();

    if (window.location.hostname == 'jakopavouk.cz' || window.location.hostname == 'test.jakopavouk.cz') {
        source = new EventSource('/api/skola/zaci-stream/' + id);
    } else {
        source = new EventSource('http://127.0.0.1:1323/api/skola/zaci-stream/' + id);
    }

    source.onmessage = function () {
        get();
        getVysledkyStudentuVPraci(selectnutaPraceID.value);
    };
});

onUnmounted(() => {
    source?.close();
});

function get() {
    axios
        .get('/skola/trida/' + id, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then((response) => {
            trida.value = response.data.trida;
            studenti.value = response.data.studenti;
            studenti.value.sort((a: { jmeno: string }, b: { jmeno: string }) => a.jmeno.localeCompare(b.jmeno));

            prace.value = [];
            for (let i = 0; i < response.data.prace.length; i++) {
                const prace1 = response.data.prace[i];
                let p: Prace = {
                    cislo: 0,
                    id: prace1.id,
                    text: prace1.text,
                    cas: prace1.cas,
                    datum: new Date(prace1.datum),
                    prumerneCPM: prace1.prumerne_cpm,
                    prumernaPresnost: prace1.prumerna_presnost,
                    StudentuDokoncilo: prace1.studentu_dokoncilo,
                    hodnocena: prace1.hodnoceni.hodnoceni1 !== null,
                };
                prace.value.push(p);
            }
            prace.value.sort((a: { datum: Date }, b: { datum: Date }) => a.datum.getTime() - b.datum.getTime());

            let counter: number[] = [1, 1];
            for (let i = 0; i < prace.value.length; i++) {
                if (prace.value[i].hodnocena) {
                    prace.value[i].cislo = counter[1];
                    counter[1]++;
                } else {
                    prace.value[i].cislo = counter[0];
                    counter[0]++;
                }
            }
            prace.value.reverse();

            vsechnyTridy.value = [];
            for (let i = 0; i < response.data.ostatniTridy.length; i++) {
                const trida = response.data.ostatniTridy[i];
                vsechnyTridy.value.push({
                    id: trida.id,
                    jmeno: trida.jmeno,
                    kod: trida.kod,
                    zamknuta: trida.zamknuta,
                    pocet_studentu: trida.pocet_studentu,
                    klavesnice: trida.klavesnice,
                    pocet_praci: trida.pocet_praci,
                });
            }
            vsechnyTridy.value.sort((a: { jmeno: string }, b: { jmeno: string }) => a.jmeno.localeCompare(b.jmeno));

            studentTridaZmena.value = trida.value.id;
        })
        .catch((e) => {
            console.log(e);
            if (e.response.data.error == 'sql: no rows in result set') {
                pridatOznameni('Taková třída neexistuje');
                router.push('/skola');
                return;
            }
            pridatOznameni('Chyba serveru');
        });
}

function select(id: number) {
    if (mobil.value) return;
    if (selectnutej.value == id) {
        //unselect
        selectnutej.value = -1;
        return;
    }
    upravaStudenta.value = false;
    selectnutej.value = id;
    studentOznacenej.value = { jmeno: '...', email: '...@...', dokonceno: 0, daystreak: 0, rychlost: -1, uspesnost: -1, klavesnice: 'QWERTZ', nejcastejsiChyby: new Map() };
    cpmVPracich.value = new Map<number, number>();
    presnostVPracich.value = new Map<number, number>();

    nacitamStudenta.value = true;

    axios
        .get('/skola/student/' + id, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then((response) => {
            studentOznacenej.value = response.data.student;

            for (const key in response.data.cpmVPracich) {
                cpmVPracich.value.set(+key, response.data.cpmVPracich[key]);
            }
            for (const key in response.data.presnostVPracich) {
                presnostVPracich.value.set(+key, response.data.presnostVPracich[key]);
            }
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
            selectnutej.value = -1;
        })
        .finally(() => {
            nacitamStudenta.value = false;
        });
}

function zmenaJmena(e: Event) {
    e.preventDefault();
    if (jmenoUprava.value == studentOznacenej.value.jmeno) {
        upravaStudenta.value = false;
        return;
    }

    if (jmenoUprava.value == '' || jmenoUprava.value.length > 30) {
        pridatOznameni('Jméno musí být 1-30 znaků dlouhé');
        return;
    }
    axios
        .post(
            '/skola/student',
            { jmeno: jmenoUprava.value, id: selectnutej.value },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            upravaStudenta.value = false;
            studentOznacenej.value.jmeno = jmenoUprava.value;
            get();
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
            selectnutej.value = -1;
        });
}

function zmenaStudentTridy(e: Event) {
    e.preventDefault();

    if (studentTridaZmena.value == trida.value.id) {
        upravaStudenta.value = false;
        return;
    }

    axios
        .post(
            '/skola/student',
            { trida_id: parseInt(studentTridaZmena.value), id: selectnutej.value },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            upravaStudenta.value = false;
            get();
            selectnutej.value = -1;
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
            selectnutej.value = -1;
        });
}

function prejmenovatTridu(e: Event, novyJmeno: string) {
    e.preventDefault();

    let staryJmeno = trida.value.jmeno;
    trida.value.jmeno = novyJmeno;
    axios
        .post(
            '/skola/zmena-tridy',
            { trida_id: trida.value.id, zmena: 'jmeno', hodnota: novyJmeno },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
            trida.value.jmeno = staryJmeno;
        });
}

function zadano() {
    prepinacTabu.value!.tab = 'prace';
    get();
}

function zadatDoJineTridy() {
    if (copyTrida.value == trida.value.id) {
        if (!confirm('Chcete práci zadat znovu do této třídy? Budou tak dvě práce se stejným textem.')) return;
    }

    axios
        .post(
            '/skola/pridat-praci',
            {
                cas: prace.value[copyPraciIndex.value].cas,
                trida_id: copyTrida.value,
                text: prace.value[copyPraciIndex.value].text,
            },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            if (copyTrida.value == trida.value.id) get();
            copyPraciIndex.value = -1;
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

const posledniRychlostPrace = computed(() => {
    if (prace.value.length == 0) return -1;

    for (let i = 0; i < prace.value.length; i++) {
        if (prace.value[i].prumerneCPM != -1) return prace.value[i].prumerneCPM;
    }
    return -1;
});

const dialog1 = useTemplateRef('dialog1');

watch(copyPraciIndex, () => {
    if (copyPraciIndex.value == -1) {
        dialog1.value?.close();
    } else {
        dialog1.value?.showModal();
    }
});

function getVysledkyStudentuVPraci(id: number) {
    axios
        .get('/skola/get-statistiky-prace/' + id, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then((response) => {
            const s: Array<Zak> = response.data.studenti;
            s.sort((a: Zak, b: Zak) => {
                if (a.cpm === -1 && b.cpm !== -1) return 1;
                if (a.cpm !== -1 && b.cpm === -1) return -1;
                if ((a.cpm === -1 && b.cpm === -1) || (a.cpm !== -1 && b.cpm !== -1)) {
                    if (a.jmeno < b.jmeno) {
                        return -1;
                    }
                    if (a.jmeno > b.jmeno) {
                        return 1;
                    }
                }
                return 0;
            });
            studentiVPraci.value.set(id, s);
            selectnutaPraceID.value = id;
        })
        .catch(() => {
            pridatOznameni('Chyba serveru');
            selectnutaPraceID.value = -1;
        });
}
</script>
<template>
    <h1 class="nadpis-se-sipkou" style="margin: 0; direction: ltr">
        <SipkaZpet />
        Třída: {{ trida.jmeno == undefined ? '-.-' : trida.jmeno }}
    </h1>

    <PrepinacTabu
        v-show="prepinacTabu?.tab != 'zadani'"
        :taby="[
            ['zaci', 'Studenti'],
            ['prace', 'Práce'],
            ['nastaveni', 'Nastavení'],
        ]"
        default-tab="zaci"
        ref="prepinac-tabu"
    />

    <div v-if="prepinacTabu?.tab == 'zaci'" id="pulic">
        <div id="kontejner">
            <div v-for="st in studenti" class="blok" @click="select(st.id)" :class="{ oznaceny: selectnutej == st.id }" :key="st.id">
                <div>
                    <h3>{{ st.jmeno }}</h3>
                    <h4>{{ st.email }}</h4>
                </div>
                <span
                    ><b>{{ naJednoDesetiny(st.cpm) }}</b> <span style="font-size: 15px">CPM</span></span
                >
            </div>
            <div v-if="studenti.length == 0 && trida.jmeno" id="text-zaci">
                <KodTridy :id="trida.id" :kod="trida.kod" :zamknuta="trida.zamknuta" />
                <h3>Jak se mohou žáci připojit?</h3>
                <ol>
                    <li>Sdělte žákům tento 4-místný kód</li>
                    <li>
                        Žáci se zapíšou:
                        <ul>
                            <li>V nastavení účtu (ozubené kolečko vlevo dole)</li>
                            <li>Tlačítkem "Zapsat se"</li>
                            <li>Napíšou kód této třídy</li>
                            <li>Zadají svoje "Školní jméno", podle kterého je zde potom uvidíte</li>
                        </ul>
                    </li>
                </ol>
            </div>
        </div>
        <div v-if="selectnutej != -1 && !mobil" class="detail">
            <div id="vrsek">
                <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk" width="200" height="126" />
                <div v-if="!upravaStudenta">
                    <h2>
                        <span>{{ studentOznacenej.jmeno }}</span>
                        <img
                            id="upravit"
                            @click="
                                upravaStudenta = true;
                                jmenoUprava = studentOznacenej.jmeno;
                            "
                            src="../../assets/icony/upravit.svg"
                            alt="Upravit"
                        />
                    </h2>
                    <h3>{{ studentOznacenej.email }}</h3>
                </div>
                <form v-else>
                    <input v-model="jmenoUprava" type="text" placeholder="Doporučeno: Příjmení Jméno" />
                    <div>
                        <select v-model="studentTridaZmena">
                            <option :value="-1">Odebrat</option>
                            <option v-for="v in vsechnyTridy" :value="v.id" :key="v.id">{{ v.jmeno }}</option>
                        </select>
                        <button
                            type="submit"
                            @click="
                                zmenaJmena($event);
                                zmenaStudentTridy($event);
                            "
                            class="tlacitko"
                            id="ulozit"
                        >
                            Potvrdit
                        </button>
                    </div>
                </form>
            </div>

            <hr />

            <div class="udaj">
                <h4>Průměrná rychlost:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.rychlost) == -1 ? '-' : naJednoDesetiny(studentOznacenej.rychlost) }} <span>CPM</span></div>
            </div>

            <div class="udaj">
                <h4>Průměrná přesnost:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.uspesnost) == -1 ? '-' : naJednoDesetiny(studentOznacenej.uspesnost) }} <span>%</span></div>
            </div>

            <div class="udaj">
                <h4>Dní v řadě procvičuje:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.daystreak) }}</div>
            </div>

            <div class="udaj">
                <h4>Postup v kurzu:</h4>
                <div>{{ naJednoDesetiny(studentOznacenej.dokonceno) }} <span>%</span></div>
            </div>

            <hr />

            <div id="posledni-prace">
                <div v-if="prace.length == 0" style="margin-top: 40px; width: 100%; border: none">
                    <span>Zatím žádné práce nebyly zadány</span>
                </div>
                <div v-else v-for="(p, i) in prace.slice(0, 3)" :key="i">
                    <h4>Práce {{ prace.length - i }}</h4>
                    <span>{{ p.datum.toLocaleDateString('cs-CZ') }}</span>

                    <div v-if="cpmVPracich.get(p.id) != undefined || nacitamStudenta" class="udaje">
                        <div>
                            {{ cpmVPracich.get(p.id) == undefined ? '...' : naJednoDesetiny(cpmVPracich.get(p.id)!) }}
                            <span>CPM</span>
                        </div>
                        <div>{{ presnostVPracich.get(p.id) == undefined ? '...' : naJednoDesetiny(presnostVPracich.get(p.id)!) }} <span>%</span></div>
                    </div>
                    <div v-else>
                        <Tooltip zprava="Nedokončil/a" :sirka="120" :vzdalenost="-3"><span>---</span></Tooltip>
                    </div>
                </div>
            </div>
        </div>
        <div v-else-if="!mobil" class="detail" id="pred-kliknutim">
            <img src="../../assets/pavoucekBezPozadi.svg" alt="Pavouk" />
            <h2 v-if="studenti.length == 0" style="font-size: 20px">Tady budou statistiky studentů!</h2>
            <h2 v-else>Vyberte studenta!</h2>
        </div>
    </div>
    <div id="prace-kontejner" v-if="prepinacTabu?.tab == 'prace'">
        <div v-if="prace.length == 0" id="text-prace">
            <span>Zatím nejsou žádné zadané práce. <br />První vytvoříte pomocí tohoto tlačítka.</span>
            <img src="../../assets/icony/sipkaOhnuta.svg" alt="Šipka na tlačítko" width="100" />
        </div>
        <PraceBlok
            v-for="(v, i) in prace"
            :key="i"
            :prace="v"
            :selectnutaPraceID
            :studentiVPraci
            @unselect="selectnutaPraceID = -1"
            @select="selectnutaPraceID = v.id"
            @reload="
                () => {
                    get();
                    getVysledkyStudentuVPraci(v.id);
                }
            "
            @copy="copyPraciIndex = i"
            :pocetStudentu="studenti.length"
        />
    </div>
    <ZadaniPrace v-else-if="prepinacTabu?.tab == 'zadani'" :tridaID="trida.id" @zadano="zadano" :posledniRychlost="posledniRychlostPrace" :hodnoticiTabulka="trida.hodnoceni" />
    <NastaveniTridy
        v-else-if="prepinacTabu?.tab == 'nastaveni'"
        ref="nastaveni"
        :trida="trida"
        :pocetStudentu="vsechnyTridy.find((t) => t.id === trida.id)!.pocet_studentu"
        @prejmenovatTridu="prejmenovatTridu"
        @refresh="get"
    />

    <dialog ref="dialog1">
        <div id="copy-menu">
            <h2>Zkopírovat práci {{ prace.length - copyPraciIndex }} do:</h2>
            <select v-model="copyTrida">
                <option :value="0">Vyberte třídu</option>
                <option v-for="t in vsechnyTridy" :value="t.id" :key="t.id">
                    {{ t.id == trida.id ? t.jmeno + ' (Tato třída)' : t.jmeno }}
                </option>
            </select>
            <div>
                <button class="tlacitko" @click="copyPraciIndex = -1">Zrušit</button>
                <button class="tlacitko" :disabled="copyTrida == 0" @click="zadatDoJineTridy">Zadat práci</button>
            </div>
        </div>
    </dialog>

    <div
        v-if="prepinacTabu?.tab == 'prace' || prepinacTabu?.tab == 'zadani'"
        id="pridat"
        @click="prepinacTabu!.tab = prepinacTabu?.tab == 'prace' ? 'zadani' : 'prace'"
        :style="{ transform: prepinacTabu?.tab == 'zadani' ? 'rotate(-45deg)' : 'rotate(0deg)' }"
    >
        <img src="../../assets/icony/plus.svg" alt="Přidat" />
    </div>
</template>
<style scoped>
dialog {
    width: 320px;
    height: 160px;
    margin-left: -160px;
    margin-top: -80px;
}

#copy-menu {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 15px;
}

#copy-menu > select {
    min-width: 230px;
    max-width: 100%;
}

#copy-menu > div {
    display: flex;
    gap: 10px;
    justify-content: center;
}

#copy-menu .tlacitko {
    margin: 0;
    height: 34px;
}

.udaje {
    display: flex;
    flex-direction: column;
    margin-top: 12px;
}

.udaje > div {
    font-weight: 500;
    font-size: 19px;
}

.udaje > div > span {
    font-size: 13px;
}

hr {
    width: 90%;
    align-self: center;
    margin: 5px 0;
}

form input::placeholder {
    color: rgba(255, 255, 255, 0.42);
}

#posledni-prace {
    display: flex;
    height: 100%;
    margin-top: 5px;
}

#posledni-prace > div:nth-child(1),
#posledni-prace > div:nth-child(2) {
    border-right: 2px solid #afafaf;
}

#posledni-prace > * {
    width: 33%;
    padding-top: 5px;
    height: 100%;
}

#posledni-prace > div > span {
    font-size: 14px;
}

#posledni-prace > div > div:not(.udaje) {
    margin-top: 20px;
    font-size: 24px;
}

#posledni-prace > div > h4 {
    font-weight: 500;
    font-size: 19px;
}

.udaj {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 80%;
    align-self: center;
}

.udaj > h4 {
    font-size: 18px;
}

.udaj > div {
    font-size: 24px;
    font-weight: 500;
}

.udaj > div > span {
    font-size: 18px;
}

#text-zaci {
    height: 380px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    align-items: center;
    justify-content: center;
}

#text-zaci > ol {
    text-align: left;
    display: flex;
    flex-direction: column;
    gap: 6px;
    max-width: 300px;
}

#text-zaci ul {
    margin-left: 24px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-top: 3px;
}

#text-zaci > h3 {
    margin-top: 25px;
    margin-bottom: 4px;
    font-size: 21px;
}

#text-prace {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    position: fixed;
    bottom: 30px;
    right: 120px;
}

#text-prace > img {
    margin-left: 230px;
}

#prace-kontejner {
    display: flex;
    flex-direction: column;
    gap: 15px;
    width: 590px;
}

.nadpis-prace h2 {
    font-size: 22px;
    font-weight: 500;
}

.nadpis-prace h3 {
    font-size: 18px;
    font-weight: 200;
    margin: 0;
}

#pridat {
    background-color: var(--tmave-fialova);
    border-radius: 100%;
    width: 55px;
    height: 55px;
    position: fixed;
    right: 30px;
    bottom: 25px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 13px;
    cursor: pointer;
    transition:
        background-color 0.15s,
        transform 0.3s;
    z-index: 10;
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
    user-select: none;
}

#pridat:hover {
    background-color: var(--fialova);
}

#pred-kliknutim {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px;
    gap: 10px;
}

#pred-kliknutim img {
    width: 200px;
}

#upravit {
    width: 26px !important;
    height: 24px !important;
    cursor: pointer;
    margin: 2px 0 0 2px;
}

.tlacitko {
    width: auto;
    padding: 0 15px;
}

#pulic {
    display: flex;
    justify-content: space-between;
    width: 860px;
    margin-bottom: -50px;
    /* countruju view margin */
}

#kontejner {
    width: 430px;
    display: flex;
    gap: 10px;
    max-height: calc(100vh - 90px - 60px - 40px - 25px - 30px - 5px);
    /* celá obrazovka - všechno co je nad seznamem zaku */
    flex-direction: column;
    overflow-y: scroll;
    padding-right: 10px;

    scrollbar-gutter: stable;
    scrollbar-width: auto;
}

.detail {
    width: 410px;
    height: 380px;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    padding: 10px 15px;
    gap: 3px;
}

.detail h3 {
    font-size: 16px;
    text-align: start;
}

#vrsek {
    margin-left: -8px;
    display: flex;
    align-items: center;
    gap: 5px;
}

#vrsek img {
    width: 110px;
    height: 80px;
    object-fit: cover;
}

#vrsek h2 {
    display: flex;
    text-wrap: nowrap;
}

#vrsek h2 span {
    max-width: 265px;
    overflow: hidden;
    text-overflow: ellipsis;
    font-weight: 400;
}

#vrsek form {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 5px;
    width: 100%;
}

#vrsek form > div {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
}

#vrsek select {
    height: 30px;
    font-size: 16px;
}

#vrsek div {
    max-width: 280px;
}

#vrsek h3 {
    text-overflow: ellipsis;
    overflow: hidden;
}

#vrsek h3:hover {
    overflow: visible;
}

#ulozit {
    font-size: 14px;
    margin-top: 0;
    width: 80px;
    height: 30px;
    padding: 0;
}

#vrsek input {
    width: 100%;
    height: 30px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 16px;
}

#vrsek input:focus {
    outline: none !important;
}

.blok {
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    padding: 6px 12px;
    cursor: pointer;
    transition: 0.1s;
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 100%;
}

.blok:hover {
    background-color: var(--fialova);
}

.blok h4 {
    font-size: 14px;
    text-overflow: ellipsis;
    overflow: hidden;
    text-wrap: nowrap;
}

.blok h4:hover {
    overflow: visible;
}

.blok h3 {
    text-wrap: nowrap;
    font-weight: 500;
}

.blok span {
    font-size: 21px;
}

.blok div {
    text-align: start;
    width: calc(100% - 120px);
}

.oznaceny {
    background-color: var(--fialova);
}

::-webkit-scrollbar {
    width: 10px;
    /* Width of the scrollbar */
}

/* Customizes the track of the scrollbar */
::-webkit-scrollbar-track {
    background: var(--tmave-fialova);
    border-radius: 3px;
    padding: 1px;
}

/* Customizes the thumb of the scrollbar */
::-webkit-scrollbar-thumb {
    background: var(--fialova);
    border-radius: 3px;
}

/* Changes the thumb color on hover */
::-webkit-scrollbar-thumb:hover {
    background: var(--svetle-fialova);
    cursor: ns-resize;
}

@media screen and (max-width: 900px) {
    #pulic {
        width: 400px;
    }

    #kontejner {
        height: auto;
        overflow-y: visible;
        padding-right: 0;
    }

    .nadpis-prace h2 {
        font-size: 18px;
    }

    .nadpis-prace h3 {
        font-size: 14px;
    }

    #prace-kontejner {
        width: auto;
        gap: 10px;
    }

    .blok {
        padding: 4px 8px;
        height: 48px;
    }

    .blok h3 {
        height: 18px;
        position: relative;
        bottom: 2px;
    }

    #text-prace {
        bottom: 40px;
        right: 20px;
    }

    #text-prace > img {
        margin-left: 20px;
        transform: rotate(8deg);
    }
}
</style>
