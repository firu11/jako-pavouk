<script>
import axios from "axios";

export default {
    name: "VysledekCvic",
    methods: {
        znovu() {
            this.$parent.reset();
        }
    },
    props: ["rychlost", "preklepy", "cas", "delka_textu"],
    data() {
        return {
            pismena: this.$route["params"].pismena,
            cislo: this.$route["params"].id,
        }
    },
    mounted() {
        axios.post('/dokonceno/' + this.pismena + '/' + this.cislo, {
            "cpm": this.rychlost,
            "preklepy": this.preklepy
        }, {
            headers: {
                "Token": this.$ls.getItem("token").value
            }
        }).catch(function (e) {
            console.log(e)
        })
    }
}
</script>


<template>
    <div id="bloky">
        <div class="blok">
            <h2>{{ Math.round(rychlost * 10) / 10 }}</h2>
            <hr>
            <p class="jednotka">CPM</p>
            <h3>Rychlost</h3>
        </div>
        <div class="blok">
            <div>
                <h2>{{ preklepy }}</h2>
                <h3 class="procento">({{ Math.round(preklepy / delka_textu * 1000) / 10 }}%)</h3>
            </div>
            <hr>
            <p class="jednotka">&zwnj;</p>
            <h3>Překlepy</h3>
        </div>
        <div class="blok">
            <h2>{{ cas[0] === 0 ? cas[1] : cas[0] + ':' + cas[1] }}s</h2>
            <hr>
            <p class="jednotka">&zwnj;</p>
            <h3>Čas</h3>
        </div>
    </div>
    <div id="tlacitka_kontainer">
        <button class="tlacitko" @click="znovu">Zkusit znovu</button>
        <button class="tlacitko" @click="$router.push('/lekce/' + this.pismena)">Pokračovat</button>
    </div>
</template>

<style scoped>
#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
}

.blok div {
    display: flex;
    align-items: baseline;
    justify-content: center;
    gap: 5px;
}

.procento {
    font-size: 20px;
    position: relative;
    top: -3px;
}

.blok {
    color: var(--bila);
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 200px;
    background-color: var(--tmave-fialova);
    height: 140px;
    transition-duration: 0.2s;
    padding: 15px 15px 30px 15px;
}

.blok h2 {
    font-size: 40px;
    font-weight: 500;
}

.blok p {
    font-size: 14px;
    margin-bottom: 8px;
    margin-top: 4px;
}

.blok hr {
    width: 160px;
    align-self: center;
}

#tlacitka_kontainer {
    display: inline-flex;
    gap: 20px;
    margin-top: 10px;
}
</style>