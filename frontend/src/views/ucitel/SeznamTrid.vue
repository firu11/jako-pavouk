<script setup lang="ts">
import axios from 'axios';
import { onMounted, onUnmounted, ref, useTemplateRef } from 'vue';
import { getToken, pridatOznameni, oznameni } from '../../utils';
import { moznostiRocnik, moznostiTrida, moznostiSkupina, prihlasen } from '../../stores';
import { useHead } from '@unhead/vue';
import { useRouter } from 'vue-router';
import PrepinacTabu from '../../components/PrepinacTabu.vue';
import SeznamUcitelu from '../../components/ucitel/SeznamUcitelu.vue';

const router = useRouter();

const skola = ref({} as { jmeno: string; id: number; aktivni: boolean });
type Trida = { id: number; jmeno: string; ucitel_id: number; kod: string; zamknuta: boolean; pocet_studentu: number; pocet_praci: number };
const rocniky = ref(new Map<string, Trida[]>());
const pridavani = ref(false);

const rocnik = ref('1.');
const trida = ref('A');
const skupina = ref('-');

const nacitam = ref(true);

const sources: EventSource[] = [];

useHead({
    title: 'Škola',
});

onMounted(() => {
    if (!prihlasen.value && getToken() == null) {
        router.push('/');
    }
    get();
});

onUnmounted(() => {
    sources.forEach((s) => {
        s.close();
    });
});

function get() {
    axios
        .get('/skola/tridy', {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        })
        .then((response) => {
            if (response.data.skola == undefined) {
                pridatOznameni('Něco se pokazilo');
                router.push('/');
                return;
            }
            skola.value = response.data.skola;

            Object.keys(response.data.tridy).forEach((key) => {
                rocniky.value.set(
                    key,
                    response.data.tridy[key].sort((a: { jmeno: string }, b: { jmeno: string }) => a.jmeno.localeCompare(b.jmeno)),
                );
            });

            sources.splice(0, sources.length - 1);
            if (sources.length == 0) {
                rocniky.value.forEach((r) => {
                    r.forEach((t) => {
                        let s: EventSource;
                        if (window.location.hostname == 'jakopavouk.cz' || window.location.hostname == 'test.jakopavouk.cz') {
                            s = new EventSource('/api/skola/zaci-stream/' + t.id);
                        } else {
                            s = new EventSource('http://127.0.0.1:1323/api/skola/zaci-stream/' + t.id);
                        }
                        s.onmessage = function () {
                            get();
                        };
                        sources.push(s);
                    });
                });
            }

            if (!skola.value.aktivni && oznameni.value.length == 0) {
                pridatOznameni('Vaše škola nemá platnou licenci a školní systém se brzy uzamkne. Pro více informací pište na e-mail: firu@jakopavouk.cz.', 15000);
            }
        })
        .catch((e) => {
            pridatOznameni('Chyba serveru');
            if (e.response === undefined) router.push('/');
            if (e.response.status == 401) {
                pridatOznameni('Sem nemáš přístup');
                router.push('/');
                return;
            }
        })
        .finally(() => {
            nacitam.value = false;
        });
}

function vytvorit(e: Event) {
    e.preventDefault();

    axios
        .post(
            '/skola/create-trida',
            { jmeno: `${rocnik.value}${trida.value}${skupina.value != '-' ? ' ￨ ' + skupina.value : ''}` },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            pridavani.value = false;
            get();
        })
        .catch((e) => {
            console.log(e);
            pridatOznameni('Chyba serveru');
        });
}

const prepinacTabu = useTemplateRef('prepinac-tabu');

const pridatEmail = ref('');
function pridatUcitele(e: Event) {
    e.preventDefault();
    axios
        .post(
            '/skola/upravit-ucitele',
            { email: pridatEmail.value, akce: 'pridat' },
            {
                headers: {
                    Authorization: `Bearer ${getToken()}`,
                },
            },
        )
        .then(() => {
            pridavani.value = false;
            pridatEmail.value = '';
        })
        .catch((e) => {
            pridatOznameni(e.response.data.error);
        });
}
</script>
<template>
    <h1 style="margin: 0">{{ skola.jmeno ? skola.jmeno : 'Škola' }}</h1>
    <PrepinacTabu
        v-show="!pridavani"
        :taby="[
            ['tridy', 'Moje Třídy'],
            ['ucitele', 'Učitelé'],
        ]"
        default-tab="tridy"
        sirka="120px"
        ref="prepinac-tabu"
    />

    <div id="tridy" v-if="prepinacTabu?.tab === 'tridy'">
        <div id="rocniky" v-if="!pridavani && rocniky.size !== 0">
            <div v-for="[rocnik, tridy] in rocniky" :key="rocnik" class="rocnik">
                <h2>{{ rocnik }}{{ isNaN(+rocnik) ? '' : '. ročník' }}</h2>
                <div id="kontejner">
                    <div class="blok" v-for="t in tridy" :key="t.id" @click="$router.push('/skola/' + t.id)">
                        <h3>{{ t.jmeno }}</h3>

                        <hr style="margin: 0 8px 8px 8px; border: #c0c0c0 1px solid" />

                        <div style="display: flex; justify-content: space-around">
                            <div class="statistiky">
                                <span v-if="t.pocet_studentu == 1"
                                    ><b>{{ t.pocet_studentu }}</b> student</span
                                >
                                <span v-else-if="t.pocet_studentu >= 2 && t.pocet_studentu <= 4"
                                    ><b>{{ t.pocet_studentu }}</b> studenti</span
                                >
                                <span v-else
                                    ><b>{{ t.pocet_studentu }}</b> studentů</span
                                >

                                <span v-if="t.pocet_praci == 0 || t.pocet_praci > 4"
                                    ><b>{{ t.pocet_praci }}</b> prací</span
                                >
                                <span v-else
                                    ><b>{{ t.pocet_praci }}</b> práce</span
                                >
                            </div>

                            <img v-if="!t.zamknuta" src="../../assets/icony/zamekOpen.svg" alt="Odemčená třída" />
                            <img v-else src="../../assets/icony/zamekClosed.svg" alt="Zamčená třída" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div v-else-if="pridavani">
            <form>
                <h2>Vytvořit třídu</h2>
                <div>
                    <h3>Ročník:</h3>
                    <select v-model="rocnik">
                        <option v-for="v in moznostiRocnik" :value="v" :key="v">{{ v }}</option>
                    </select>
                </div>

                <div>
                    <h3>Písmeno:</h3>
                    <select v-model="trida">
                        <option v-for="v in moznostiTrida" :value="v" :key="v">{{ v }}</option>
                    </select>
                </div>

                <div>
                    <h3>Skupina:</h3>
                    <select v-model="skupina">
                        <option v-for="v in moznostiSkupina" :value="v" :key="v">{{ v }}</option>
                    </select>
                </div>

                <button class="tlacitko" @click="vytvorit">Vytvořit</button>

                <span>Podle jména se třídy řadí do ročníků v seznamu tříd.</span>
                <span>
                    Skupina je určena pro rozdělení třídy. To se může hodit, pokud vyučujete třídy po menších skupinkách.
                    <br />
                    Např.: <b>3.B&nbsp;￨&nbsp;1</b> a <b>3.B&nbsp;￨&nbsp;2</b>.
                </span>
            </form>
        </div>
        <div v-else-if="rocniky.size === 0 && nacitam">Načítám...</div>
        <div v-else style="background-color: var(--tmave-fialova); padding: 20px; border-radius: 10px; max-width: 450px; margin: 0 5vw">
            <h2 style="font-size: 21px">Vítejte v rozhraní pro učitele!</h2>
            <br />
            Zde se vám budou třídy řadit do ročníků.
            <br />
            Žáci se do nich mohou připojit pomocí 4místného kódu, který bude pro každou třídu automaticky vytvořen.
        </div>

        <div v-if="rocniky.size == 0 && !pridavani && !nacitam" id="text-prace">
            <span>Zatím tu nejsou žádné třídy. <br />První vytvoříte pomocí tohoto tlačítka.</span>
            <img src="../../assets/icony/sipkaOhnuta.svg" alt="Šipka na tlačítko" width="100" />
        </div>
    </div>
    <SeznamUcitelu v-else-if="prepinacTabu?.tab === 'ucitele' && !pridavani" />
    <div v-if="pridavani && prepinacTabu?.tab === 'ucitele'">
        <form>
            <h2>Přidat učitele</h2>
            <div>
                <h3>E-mail:</h3>
                <input type="email" placeholder="Např: pavoukova@jakopavouk.cz" v-model="pridatEmail" />
            </div>

            <button class="tlacitko" @click="pridatUcitele">Přidat</button>
        </form>
    </div>
    <div id="pridat" @click="pridavani = !pridavani" :style="{ transform: pridavani ? 'rotate(-45deg)' : 'rotate(0deg)' }">
        <img src="../../assets/icony/plus.svg" alt="Přidat" />
    </div>
</template>
<style scoped>
h1 {
    max-width: 90%;
}

.tlacitko {
    width: 100px;
    align-self: center;
}

form {
    display: flex;
    flex-direction: column;
    max-width: 450px;
    align-items: center;
    gap: 8px;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    padding: 15px 25px;
    margin: 5.8em 5vw 0 5vw;
}

form > h2 {
    margin-bottom: 15px;
}

form > div {
    display: flex;
    width: 100%;
    align-items: center;
    justify-content: space-between;
    gap: 10px;
}

form > span {
    opacity: 0.5;
    font-size: 14px;
    margin-bottom: -4px;
}

form span:first-of-type {
    margin-top: 8px;
}

form span:last-of-type {
    margin-bottom: 0;
}

form input[type='text'],
form input[type='email'] {
    width: 73%;
    height: 36px;
    background-color: var(--fialova);
    border: 0;
    border-radius: 5px;
    transition: all 0.2s cubic-bezier(0.5, 0, 0.5, 1) 0s;
    color: var(--bila);
    padding: 10px;
    font-weight: normal;
    font-size: 20px;
}

form input[type='text']:focus,
form input[type='email']:focus {
    width: 75%;
}

::placeholder {
    color: rgba(255, 255, 255, 0.42);
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
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
    user-select: none;
}

#pridat:hover {
    background-color: var(--fialova);
}

#rocniky {
    display: flex;
    flex-direction: column;
    gap: 25px;
    width: 640px;
}

.rocnik {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: start;
    gap: 12px;
}

.rocnik h2 {
    margin-left: 5px;
}

.blok {
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    cursor: pointer;
    transition: 0.1s;
    width: max(30%, 200px);
    position: relative;
}

.blok img {
    width: 30px;
    height: 30px;
    position: relative;
    top: 5px;
    right: 3px;
}

.blok h3 {
    margin-bottom: 6px;
    font-size: 22px;
    font-weight: 500;
}

.blok .statistiky {
    display: flex;
    flex-direction: column;
}

.blok:hover {
    background-color: var(--fialova);
}

#kontejner {
    display: flex;
    width: 100%;
    flex-wrap: wrap;
    gap: 20px;
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

@media screen and (max-width: 800px) {
    #rocniky {
        width: 420px;
    }
}

@media screen and (max-width: 450px) {
    #tridy {
        transform: scale(0.85);
        position: absolute;
        top: 130px;
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

@media screen and (max-width: 380px) {
    #tridy {
        max-width: 200px;
        top: 130px;
    }

    #rocniky {
        width: 100%;
    }
}
</style>
