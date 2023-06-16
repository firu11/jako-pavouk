<script>
import axios from 'axios'

export default {
    name: "ucet",
    data() {
        return {
            info: {},
        }
    },
    methods: {
        odhlasit() {
            this.$ls.clear()
            this.$router.push("/")
        }
    },
    mounted() {
        if (this.$ls.getItem("token").value) {
            axios.get('/ja', {
                headers: {
                    "Token": this.$ls.getItem("token").value
                }
            }).then(response => {
                this.info = response.data
            }).catch(function (e) { //nebudeš tam chodit nemas ucet more
                this.$router.push("/login")
            })
        } else {
            this.$router.push("/login")
        }

    }
}
</script>

<template>
    <div id="ucet">
        <img src="/icony/user.svg" alt="uzivatel">
        <div id="nadpisy">
            <h1>{{ info.jmeno }}</h1>
            <h2>{{ info.email }}</h2>
        </div>
    </div>
    <div id="progres">
        <div id="nacitani-pozadi">
            <div id="nacitani" :style="{ width: info.dokonceno + '%' }"></div>
        </div>
        <span class="popis" style="width: 100%;">Dokončeno: <span class="cislo">{{ Math.round(info.dokonceno * 10) / 10 }}%</span></span>
    </div>
    <div id="bloky">
        <div class="blok">
            <img src="/icony/rychlost.svg" alt="Rychlost" width="75">
            <span v-if="info.prumerRychlosti == -1">Zatím nic</span>
            <span v-else class="popis">Rychlost: <br><span class="cislo">{{ Math.round(info.prumerRychlosti * 10) / 10 }}</span> CPM</span>
        </div>
        <div class="blok">
            <img src="/icony/terc.svg" alt="Přesnost">
            <span v-if="info.uspesnost == -1">Zatím nic</span>
            <span v-else class="popis">Přesnost: <br><span class="cislo">{{ Math.round(info.uspesnost * 10) / 10 }}</span> %</span>
        </div>
    </div>

    <button @click="odhlasit" class="tlacitko">Odhlásit</button>
</template>

<style scoped>
.popis {
    font-size: 15pt;
    width: 60%;
}

.cislo {
    font-size: 28pt;
    font-weight: 500;
}

#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
}

.blok {
    display: flex;
    text-decoration: none;
    border-radius: 10px;
    justify-content: space-evenly;
    align-items: center;
    width: 320px;
    background-color: var(--tmave-fialova);
    height: 120px;
    transition-duration: 0.2s;
    padding: 15px;
    gap: 10px;
}

#progres {
    margin-bottom: 20px;
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 420px;
    background-color: var(--tmave-fialova);
    height: 110px;
    transition-duration: 0.2s;
    padding: 15px;
    gap: 10px;
}

#nadpisy {
    display: flex;
    flex-direction: column;
    justify-content: center;
}

#nadpisy h1 {
    margin-bottom: 0;
    align-self: flex-start;
}

#ucet img {
    height: 100px;
}

#ucet {
    display: flex;
    background-color: var(--tmave-fialova);
    margin-bottom: 40px;
    padding: 15px 30px 15px 15px;
    border-radius: 10px;
    gap: 15px;
}

.tlacitko {
    margin-top: 20px;
    background-color: var(--tmave-fialova);
}

.tlacitko:hover {
    background-color: var(--fialova);
}

#bloky div img {
    height: 65px;
}

#bloky div h2 {
    font-weight: 500;
}

#nacitani-pozadi {
    height: 20px;
    background-color: var(--fialova);
    border-radius: 5px;
    padding: 0;
}

#nacitani {
    background-color: var(--bila);
    height: 20px;
    border-radius: 5px;
    position: relative;
    left: -1px
}
</style>