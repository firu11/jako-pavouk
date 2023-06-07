<script>
import Klavesnice from "@/components/Klavesnice.vue";
import VysledekCvic from "@/components/VysledekCvic.vue";
import axios from 'axios'

export default {
    name: "cviceni",
    components: {Klavesnice, VysledekCvic},
    data() {
        return {
            info: {},
            cas: [0, 0],
            pismena: this.$route["params"].pismena,
            cislo: this.$route["params"].id,
            counter: 0,     //         0             1                2
            list_textu: {}, // {str pismeno, bool aktivni, bool je/bylo_spatne}
            timer_zacatek: null,
            capslock: false,
            text_pripraven: false,
            dokonceno: false,
            pocitadloCasu: null,
            audio: []
        }
    },
    computed: {
        cas_format() {
            return this.cas[0] === 0 ? this.cas[1] : `${this.cas[0]}:${this.cas[1] > 9 ? this.cas[1] : "0" + this.cas[1]}`; // MM:SS
        },
        progress() {
            return Math.floor(((this.counter) / Object.keys(this.list_textu).length) * 100)
        },
        preklepy() {
            let pocet = 0
            for (const [key, value] of Object.entries(this.list_textu)) {
                if (value[2]) pocet++
            }
            return pocet
        },
    },
    mounted() {
        this.get()
        document.addEventListener("keydown", this.klik);
        this.audio = [new Audio('/klik1.mp3'), new Audio('/klik2.mp3'), new Audio('/klik3.mp3'), new Audio('/miss.mp3')]
    },
    unmounted() {
        document.removeEventListener("keydown", this.klik);
    },
    methods: {
        get() {
            axios
                .get("/cvic/" + this.pismena + "/" + this.cislo, {
                    headers: {
                        "Token": this.$ls.getItem("token").value
                    }
                })
                .then(response => {
                    this.info = response.data
                    let pismenoCount = 0
                    for (const slovo in this.info.text) {
                        for (const pismeno in this.info.text[slovo]) {
                            this.list_textu[pismenoCount] = [this.info.text[slovo][pismeno], false, false];
                            pismenoCount++
                        }
                    }
                    this.list_textu[0][1] = true //prvni pismeno podtrhu
                    this.text_pripraven = true
                }).catch(e => {
                    setTimeout(() => {this.$router.push('/404')}, 2000);
                });
        },
        klik(e) {
            if (e.key === " ") {
                e.preventDefault() //ať to nescrolluje
            }
            this.capslock = e.getModifierState('CapsLock')
            if (["Shift", "CapsLock"].includes(e.key)) return

            if (!this.timer_zacatek) {
                this.timer_zacatek = new Date().getTime();
                this.pocitadloCasu = setInterval(() => {
                    this.cas[1]++
                    if (this.cas[1] === 60) {
                        this.cas[1] = 0
                        this.cas[0]++
                    }
                }, 1000);
            }

            if (e.key === this.list_textu[this.counter][0]) {
                if (this.counter === Object.keys(this.list_textu).length - 1) {
                    this.dokonceno = true
                    clearInterval(this.pocitadloCasu)
                    document.removeEventListener("keydown", this.klik);
                    this.audio[Math.floor(Math.random() * 2.4)].cloneNode(true).play()

                } else {
                    this.audio[Math.floor(Math.random() * 2.4)].cloneNode(true).play()
                    this.dalsi()
                }
            } else {
                this.audio[3].cloneNode(true).play()
                this.list_textu[this.counter][2] = true
            }
        },
        dalsi() {
            this.list_textu[this.counter][1] = false
            this.counter++
            this.list_textu[this.counter][1] = true
        },
        get_index_pismena(i, j) {
            let index = j
            for (let k = 0; k < i; k++) {
                index += this.info.text[k].length
            }
            return index
        },
        reset() {
            this.cas = [0, 0]
            this.counter = 0
            this.dokonceno = false
            this.timer_zacatek = null
            this.get()
            document.addEventListener("keydown", this.klik);
        }
    }
}
</script>

<template>
    <h1 v-if="!dokonceno" style="margin: 0"><router-link class="tlacZpet" :to="'/lekce/' + this.pismena"><img src="/icony/sipkaL.svg" alt="Zpět"></router-link>Lekce: {{ $format(pismena) }}</h1>
    <h2 v-if="!dokonceno">Cviceni: {{cislo}}</h2>
    <h1 v-else>Výsledky lekce: {{ $format(pismena) }}</h1>
    <div id="obsah" v-if="!info.error && !dokonceno">

        <div id="cviceniNabidka">
            <h3 id="cas">{{ cas_format }}s</h3>
            <h3 :style="{visibility: capslock  ? 'visible' : 'hidden'}" id="capslock">CapsLock</h3>
            <h3 id="preklepy">Chyby: {{ preklepy }}</h3>
        </div>
        <div id="pozadi_ramecku">
            <div id="ramecek">
                <div id="text">
                    <div class="slovo" v-for="(slovo, i) in info.text">
                        <div v-for="(pismeno, j) in slovo">
                            <p class="pismeno"
                               :class="{podtrzenePismeno: this.list_textu[get_index_pismena(i,j)][1],
                                    spatnePismeno: this.list_textu[get_index_pismena(i,j)][2] && this.counter <= (get_index_pismena(i,j)),
                                    opravenePismeno: this.list_textu[get_index_pismena(i,j)][2] && this.counter > (get_index_pismena(i,j)),
                                    pismenoCoBylo: this.counter > (get_index_pismena(i,j))}"
                               :id="'p' + (i * slovo.length + j)">{{
                                    (pismeno !== ' ' ? pismeno : "&nbsp")
                                }}</p>
                        </div>
                    </div>
                </div>
            </div>
            <div :style="'width:' + progress + '%; border-bottom-right-radius:' + (progress === 100 ? '10px':'0')"
                 id="progress_bar">&nbsp{{ progress }}%&nbsp
            </div>
        </div>
        <Klavesnice id="klavesnice" v-if="!dokonceno && this.list_textu[this.counter] !== undefined" :aktivniPismeno="this.list_textu[this.counter]"></Klavesnice>
    </div>
    <VysledekCvic v-else-if="dokonceno" :rychlost="Object.keys(list_textu).length / parseFloat(cas[1] * 100 / 60 < 10 ? cas[0] + '.0' + cas[1] * 100 / 60 : cas[0] + '.' + cas[1] * 100 / 60)" :preklepy="preklepy" :cas="cas" :delka_textu="Object.keys(list_textu).length"></VysledekCvic>
    <p v-else>{{ info.error }}</p>

</template>

<style scoped>
#obsah {
    display: flex;
    flex-direction: column;
}

#klavesnice {
    align-self: center;
}

#cviceniNabidka {
    margin: 20px 0 6px 0;
}

#text {
    display: flex;
    flex-wrap: wrap;
}

#cas {
    float: left;
    width: 150px;
    display: inline-block;
    text-align: left;
}

#preklepy {
    float: right;
    display: inline-block;
    width: 150px;
    text-align: right;
}

#capslock {
    display: inline-block;
    color: red;
    font-weight: bold;
}

.slovo {
    display: flex;
    flex-wrap: nowrap;
}

.pismeno {
    border-radius: 3px;
    display: inline-flex;
    font-family: Menlo, Monaco, Consolas, Courier New, monospace;
    font-size: 25px;
    line-height: 1.2;
    text-decoration: none;
    padding: 0 1px;
    margin-right: 1px;
    border-bottom: 3px solid rgba(255, 255, 255, 0); /* aby se nedojebala vyska liny když jdu na dalsi radek*/
    color: #8c8c8c;
}

.pismenoCoBylo {
    color: white;
}

#ramecek {
    padding: 10px 10px 0 10px;
    min-height: 190px;
    display: flex;
    justify-content: center;
    align-items: flex-start;
}

#pozadi_ramecku {
    border-radius: 10px;
    background-color: var(--tmave-fialova);
}

#progress_bar {
    height: 20px;
    background-color: var(--fialova);
    width: 0;
    border-bottom-left-radius: 10px;
    transition: ease 0.5s;
    text-align: right;
}

.podtrzenePismeno {
    border-bottom: 3px solid var(--bila);
    border-radius: 0;
}

.spatnePismeno {
    background-color: var(--cervena);
    border-radius: 3px 3px 0 0;
}

.opravenePismeno {
    background-color: var(--fialova);
}

h1 {
    display: inline-flex;
    position: relative;
    right: 25px; /* posunuti o pulku sipky */
}
</style>
