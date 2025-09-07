<script setup lang="ts">
import axios from 'axios';
import { getToken, pridatOznameni, naJednoDesetiny } from '../../utils';
import { computed, onMounted, ref } from 'vue';
import { useHead } from '@unhead/vue';
import { mobil, role } from '../../stores';
import Tooltip from '../../components/Tooltip.vue';
import { useRouter } from 'vue-router';

useHead({
    title: 'Třída',
});

const router = useRouter();

const trida = ref({} as { id: number; jmeno: string; ucitel_id: number; kod: string; zamknuta: boolean; pocet_studentu: number; klavesnice: string });
const praceNove = ref([] as { id: number; cislo: number; datum: string; cpm: number; presnost: number; hodnocena: boolean; znamka: number }[]);
const praceDoko = ref([] as { id: number; cislo: number; datum: string; cpm: number; presnost: number; hodnocena: boolean; znamka: number }[]);

const nacitam = ref(false);

const klavesniceStudenta = ref('');

onMounted(() => {
    if (mobil.value) {
        history.back();
        pridatOznameni('Psaní na telefonech zatím neučíme...');
    }
    get();
});

function get() {
    nacitam.value = true;
    axios
        .get('/skola/trida', {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then((response) => {
            trida.value = response.data.trida;
            klavesniceStudenta.value = response.data.klavesnice;

            praceNove.value = [];
            praceDoko.value = [];

            response.data.prace.sort((a: { datum: string }, b: { datum: string }) => b.datum.localeCompare(a.datum));

            let counter: number[] = [1, 1];
            for (let i = response.data.prace.length - 1; i >= 0; i--) {
                const prace1 = response.data.prace[i];
                if (prace1.hodnocena) {
                    prace1.cislo = counter[1];
                    counter[1]++;
                } else {
                    prace1.cislo = counter[0];
                    counter[0]++;
                }

                const p = {
                    id: prace1.id,
                    cislo: prace1.cislo,
                    datum: new Date(prace1.datum).toLocaleDateString('cs-CZ'),
                    cpm: prace1.cpm,
                    presnost: prace1.presnost,
                    hodnocena: prace1.hodnocena,
                    znamka: prace1.znamka,
                };

                if (prace1.cpm != -1) praceDoko.value.push(p);
                else praceNove.value.push(p);
            }

            praceDoko.value.reverse();
            praceNove.value.reverse();
        })
        .catch(() => {
            pridatOznameni('Chyba serveru');
            router.back();
            role.value = 'basic';
        })
        .finally(() => {
            nacitam.value = false;
        });
}

function prvniVelky(s: string) {
    if (!s) return '---------';
    return s[0].toUpperCase() + s.slice(1);
}

const dobraKlavesnice = computed(() => {
    if (klavesniceStudenta.value == '') return true;
    return klavesniceStudenta.value == trida.value.klavesnice;
});

const zprava = computed(() => {
    if (dobraKlavesnice.value) return 'Používáš správnou klávesnici!';
    return `Třída používá <b>${prvniVelky(trida.value.klavesnice)}</b>, zatímco ty používáš <b>${prvniVelky(klavesniceStudenta.value)}</b>. Nějaké texty by ti nemuseli sedět! Kliknutím si ji změň.`;
});

function zmenaKlavesnice() {
    if (dobraKlavesnice.value) return;
    axios
        .post('/ucet-zmena', { zmena: 'klavesnice', hodnota: trida.value.klavesnice }, { headers: { Authorization: `Bearer ${getToken()}` } })
        .then(() => {
            praceDoko.value = [];
            praceNove.value = [];
            get();
        })
        .catch(() => {
            pridatOznameni();
        });
}
</script>
<template>
    <h1>
        Třída: {{ trida.jmeno == undefined ? '-.-' : trida.jmeno }}
        <Tooltip @click="zmenaKlavesnice" :zprava="zprava" :sirka="230" :vzdalenost="2" style="cursor: pointer">
            <span :class="{ blba: !dobraKlavesnice, dobra: dobraKlavesnice }">
                {{ prvniVelky(trida.klavesnice) }}
            </span>
        </Tooltip>
    </h1>

    <div v-if="!nacitam && (praceNove.length != 0 || praceDoko.length != 0)" id="kontejner">
        <h2>Čeká na dokončení</h2>
        <div v-if="praceNove.length != 0" class="prace-kontejner">
            <RouterLink :to="`/prace/${v.id}`" v-for="v in praceNove" :key="v.id" class="prace">
                <div class="nadpis-prace">
                    <h3 v-if="!v.hodnocena">Práce {{ v.cislo }}</h3>
                    <h3 v-else>Zkouška {{ v.cislo }}</h3>
                    <h4>{{ v.datum }}</h4>
                </div>
                <img class="play" src="../../assets/icony/start.svg" alt="Dokonceno!" />
            </RouterLink>
        </div>
        <span v-else>Žádné</span>

        <h2>Dokončené</h2>
        <div v-if="praceDoko.length != 0" class="prace-kontejner">
            <div v-for="v in praceDoko" :key="v.id" class="prace hotova">
                <div class="nadpis-prace">
                    <h3 v-if="!v.hodnocena">Práce {{ v.cislo }}</h3>
                    <h3 v-else>Zkouška {{ v.cislo }}</h3>
                    <h4>{{ v.datum }}</h4>
                </div>
                <div class="statistika">
                    <Tooltip zprava="Rychlost" :sirka="90" :vzdalenost="0">
                        <span
                            ><b>{{ naJednoDesetiny(v.cpm) }}</b> CPM
                        </span>
                    </Tooltip>

                    <Tooltip zprava="Přesnost" :sirka="90" :vzdalenost="0">
                        <span>
                            <b>{{ naJednoDesetiny(v.presnost) }}</b> %
                        </span>
                    </Tooltip>

                    <Tooltip v-if="v.hodnocena" zprava="Známka" :sirka="90" :vzdalenost="0">
                        <span>
                            <b>{{ v.znamka }}</b>
                        </span>
                    </Tooltip>
                    <Tooltip v-else zprava="Známka <br> (není hodnoceno)" :sirka="150" :vzdalenost="0">
                        <span>
                            <b>-</b>
                        </span>
                    </Tooltip>
                </div>
            </div>
        </div>
        <span v-else>Žádné</span>
    </div>

    <span v-if="praceNove.length == 0 && praceDoko.length == 0 && !nacitam" id="textZaci">Zatím tu nejsou žádné zadané práce.</span>
</template>
<style scoped>
h1 {
    margin-left: calc(65px + 10px);
    display: flex;
    gap: 10px;
}

h1 span {
    font-size: 14px;
    border-radius: 100px;
    padding: 5px;
    position: relative;
    top: -5px;
    user-select: none;
    transition: 0.15s;
}

h1 span.dobra {
    background-color: rgba(0, 128, 0, 0.5);
    border: 2px solid rgb(0, 128, 0);
}

h1 span.blba {
    background-color: rgba(180, 0, 0, 0.5);
    border: 2px solid rgb(170, 0, 0);
    animation: klavesniceSpatnaSviti 1s infinite 0.4s;
}

@keyframes klavesniceSpatnaSviti {
    0% {
        background-color: rgba(180, 0, 0, 0.5);
        border: 2px solid rgb(170, 0, 0);
    }

    50% {
        background-color: rgba(218, 0, 0, 0.5);
        border: 2px solid rgb(255, 0, 0);
    }

    100% {
        background-color: rgba(180, 0, 0, 0.5);
        border: 2px solid rgb(170, 0, 0);
    }
}

.statistika div span b {
    font-family: 'Red Hat Mono';
    font-size: 29px;
}

.statistika div span {
    font-size: 19px;
    display: flex;
    align-items: baseline;
    gap: 5px;
    justify-content: end;
    height: 34px;
    align-self: center;
    position: relative;
    top: -2px;

    min-width: 115px;
}

.statistika > div:last-of-type span {
    min-width: 43px;
    top: -2px;

    b {
        font-size: 29px;
    }
}

.statistika {
    display: flex;
}

.play {
    height: 32px;
    margin-right: 8px;
    margin-bottom: 1px;
}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

.hotova {
    color: var(--seda) !important;
    background-color: rgba(63, 51, 81, 0.8) !important;
    cursor: default !important;
}

#kontejner {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}

#kontejner > span {
    align-self: center;
    margin: 10px 0 10px 0;
}

.prace-kontejner {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.prace {
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 6px 12px;
    cursor: pointer;
    transition: 0.1s;
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 500px;
    height: 64px;
    text-decoration: none;
    color: var(--bila);
}

.prace:not(.hotova):hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.nadpis-prace {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.nadpis-prace h3 {
    font-size: 22px;
    font-weight: 500;
}

.nadpis-prace h4 {
    font-size: 18px;
    font-weight: 200;
    margin: 0;
}

@media screen and (max-width: 1100px) {
    h1 {
        direction: ltr;
    }
}
</style>
