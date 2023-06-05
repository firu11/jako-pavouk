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
            this.$ls.removeItem("token")
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
                console.log(this.info)
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
    <div id="bloky">
        <div class="blok">
            <div id="nacitani-pozadi">
                <div id="nacitani" :style="{ width: info.dokonceno + '%' }"></div>
            </div>
            <h2>Dokonceno: {{ info.dokonceno }}%</h2>
        </div>
        <div class="blok">
            <img src="/icony/rychlost.svg" alt="Rychlost">
            <h2 v-if="info.prumerPreklepu == -1">Zatím nic</h2>
            <h2 v-else>{{ Math.round(info.prumerPreklepu * 10) / 10 }} CPM</h2>
        </div>
        <div class="blok">
            <img src="/icony/terc.svg" alt="Rychlost">
            <h2 v-if="info.prumerRychlosti == -1">Zatím nic</h2>
            <h2 v-else>{{ Math.round(info.prumerRychlosti * 10) / 10 }}%</h2>
        </div>
    </div>

    <button @click="odhlasit" class="tlacitko">Odhlásit</button>
</template>

<style scoped>
#bloky {
    display: flex;
    flex-direction: row;
    gap: 20px;
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
    margin-bottom: 25px;
    padding: 15px 30px 15px 15px;
    border-radius: 10px;
    gap: 15px;
}

.tlacitko {
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