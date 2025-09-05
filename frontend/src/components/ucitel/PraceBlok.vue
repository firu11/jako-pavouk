<script setup lang="ts">
import axios from 'axios';
import { getToken, pridatOznameni, naJednoDesetiny } from '../../utils';
import { ref } from 'vue';
import Tooltip from '../../components/Tooltip.vue';
import { mobil } from '../../stores';

export type Prace = { cislo: number; id: number; text: string; cas: number; datum: Date; prumerneCPM: number; prumernaPresnost: number; StudentuDokoncilo: number; hodnocena: boolean };
export type Zak = { id: number; jmeno: string; email: string; cpm: number; presnost: number; znamka: number };

const emit = defineEmits(['unselect', 'select', 'reload', 'copy']);
const props = defineProps<{
    prace: Prace;
    selectnutaPraceID: number;
    studentiVPraci: Map<number, Array<Zak>>;
    pocetStudentu: number;
}>();

const smazatPraciID = ref(0);

function smazatPraci() {
    if (props.prace.StudentuDokoncilo != 0) {
        if (!confirm(`Tuto práci už dokončilo ${props.prace.StudentuDokoncilo} studentů! Opravdu ji chcete smazat?`)) return;
    }

    axios
        .delete('/skola/smazat-praci/' + props.prace.id, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then(() => {
            emit('reload');
            pridatOznameni('Práce byla smazána.');
            smazatPraciID.value = 0;
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

function selectPraci(id: number) {
    if (props.selectnutaPraceID == id) {
        //unselect
        emit('unselect');
        return;
    }
    if (props.pocetStudentu === 0) return;

    emit('reload');
}
</script>
<template>
    <div
        class="jedna-prace"
        :style="{
            maxHeight: selectnutaPraceID !== prace.id ? (mobil ? '48px' : '60px') : `${60 + 10 + 20 + (37.5 + 5) * pocetStudentu || 0}px`,
            marginBottom: selectnutaPraceID === prace.id ? '15px' : 0,
        }"
    >
        <div class="radek">
            <div id="prace-uprava-kontejner">
                <div class="uprava-pill" :style="{ opacity: smazatPraciID == prace.id || smazatPraciID == 0 ? 1 : 0.4 }">
                    <div class="copy-btn" @click="emit('copy')">
                        <img src="../../assets/icony/copy.svg" alt="Kopírovat práci" />
                    </div>
                    <div v-if="smazatPraciID != prace.id" class="smazat-btn" @click="smazatPraciID = prace.id">
                        <img src="../../assets/icony/trash.svg" alt="Smazat práci" />
                    </div>
                    <div v-else class="smazat-btn" @click="smazatPraci" @mouseleave="smazatPraciID = 0">
                        <img src="../../assets/icony/right.svg" alt="Smazat práci" />
                    </div>
                </div>
            </div>
            <div class="prace" :class="{ nekliknutelna: pocetStudentu == 0 }" @click="selectPraci(prace.id)">
                <Tooltip :zprava="`<b>${prace.cas / 60} min</b> | ${prace.text.slice(0, 100)}...`" :sirka="300" :vzdalenost="3">
                    <div class="nadpis-prace">
                        <h2 v-if="!prace.hodnocena">Práce {{ prace.cislo }}</h2>
                        <h2 v-else>Písemka {{ prace.cislo }}</h2>
                        <h3>{{ prace.datum.toLocaleDateString('cs-CZ') }}</h3>
                    </div>
                </Tooltip>

                <div class="statistika">
                    <Tooltip v-if="prace.prumerneCPM != -1" zprava="Průměrná rychlost studentů" :sirka="160" :vzdalenost="5">
                        <span>
                            <b>{{ naJednoDesetiny(prace.prumerneCPM) }}</b> CPM
                        </span>
                    </Tooltip>
                    <Tooltip v-if="prace.prumernaPresnost != -1" zprava="Průměrná přesnost studentů" :sirka="160" :vzdalenost="5">
                        <span>
                            <b>{{ naJednoDesetiny(prace.prumernaPresnost) }}</b> %
                        </span>
                    </Tooltip>
                    <Tooltip zprava="Studentů kteří mají hotovo" :sirka="160" :vzdalenost="5">
                        <span class="udaj2" :style="{ 'min-width': pocetStudentu > 10 ? '115px' : '80px' }">
                            <b>{{ prace.StudentuDokoncilo }}</b
                            >/<b>{{ pocetStudentu }}</b>
                        </span>
                    </Tooltip>
                </div>
            </div>
        </div>
        <div class="zaci-v-praci">
            <div v-for="zak in studentiVPraci.get(prace.id)" :key="zak.id" class="zak-v-praci">
                <span>{{ zak.jmeno }}</span>
                <div v-if="zak.cpm !== -1" class="statistika-zaka">
                    <span>
                        <b>{{ naJednoDesetiny(zak.cpm) }}</b> CPM
                    </span>
                    <span>
                        <b>{{ naJednoDesetiny(zak.presnost) }}</b> %
                    </span>
                    <span>
                        <b v-if="prace.hodnocena">{{ naJednoDesetiny(zak.znamka) }}</b>
                        <b v-else>-</b>
                    </span>
                </div>
                <div v-else>Ještě nedokončil</div>
            </div>
        </div>
    </div>
</template>
<style scoped>
.zaci-v-praci {
    display: flex;
    flex-direction: column;
    gap: 5px;

    width: 70%;
    align-self: center;
}

.zak-v-praci {
    display: flex;
    justify-content: space-between;

    background-color: var(--tmave-fialova);
    padding: 2px 12px;
    border-radius: 6px;

    font-size: 1.1em;

    min-height: 33px;
}

.uprava-pill {
    height: 60px;
    width: 28px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    transition: 0.2s;
}

.smazat-btn,
.copy-btn {
    width: 100%;
    padding: 5px;
    height: 28px;
    cursor: pointer;
    border-radius: 6px;
    background-color: var(--tmave-fialova);
    transition: 0.15s;
}

.smazat-btn:hover,
.copy-btn:hover {
    background-color: var(--fialova);
}

.copy-btn > img,
.smazat-btn > img {
    position: relative;
}

.udaj2 {
    font-size: 26px !important;
}

.statistika,
.statistika-zaka {
    display: flex;
    align-self: center;
}

.jedna-prace {
    display: flex;
    flex-direction: column;
    gap: 10px;
    overflow: hidden;
    max-height: 60px;

    transition:
        max-height 0.3s,
        margin-bottom 0.3s;
}

.prace {
    width: 100%;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 6px 12px;
    transition: 0.1s;
    display: flex;
    justify-content: space-between;
    align-items: center;
    max-width: 100%;
    height: 60px;

    transition: background-color 0.15s;
}

.prace:not(.nekliknutelna):hover {
    background-color: var(--fialova);
    cursor: pointer;
}

.nadpis-prace {
    display: flex;
    flex-direction: column;
    align-items: start;
}

.radek {
    display: flex;
    gap: 8px;
}

.statistika span b {
    font-family: 'Red Hat Mono';
    font-size: 29px;
}

.statistika-zaka span b {
    font-family: 'Red Hat Mono';
    font-size: 22px;
}

.statistika-zaka span {
    font-size: 0.7em;
    min-width: 107px;
    text-align: end;
}

.statistika-zaka:has(> :last-child:nth-child(3)) {
    span {
        font-size: 0.7em;
        min-width: 80px;
        text-align: end;
    }
    span:last-of-type {
        min-width: 40px;
    }
}

.zak-v-praci > span,
.zak-v-praci > div:not(:has(span)) {
    height: 22px;
    align-self: center;
}

.statistika span {
    font-size: 19px;
    display: flex;
    align-items: baseline;
    gap: 5px;
    justify-content: end;
    height: 34px;
    position: relative;
    top: -2px;

    min-width: 115px;
}

@media screen and (max-width: 900px) {
    .udaj2 {
        min-width: 80px !important;
    }

    .uprava-pill {
        display: none;
    }

    .prace {
        height: 48px;
        width: 400px;
        padding: 8px;
    }

    .jedna-prace {
        max-height: 48px;
    }

    .statistika *:not(#tooltip) {
        font-size: 21px !important;
        height: 22px !important;
    }

    .statistika span {
        min-width: 90px;
    }

    .nadpis-prace h2 {
        font-size: 18px;
    }

    .nadpis-prace h3 {
        font-size: 14px;
    }

    .zaci-v-praci {
        width: 85%;
    }

    .statistika-zaka span b {
        font-size: 18px;
    }

    .statistika-zaka span {
        min-width: 80px;
    }

    .radek {
        gap: 0;
    }
}
</style>
