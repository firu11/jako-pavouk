<script>
import axios from 'axios'

export default {
    name: "Lekce",
    data() {
        return {
            pismena: this.$route["params"].pismena,
            info: {},
            dostal_jsem_data: false
        }
    },
    methods: {
        prihlaste_se() {
            alert('Nejprve se prosím přihlašte')
        },
        index1(index) {
            return index + 1
        },
        jeDokoncene(id) {
            return this.info['dokoncene'].includes(id)
        }
    },
    mounted() {
        axios
            .get('/lekce/' + this.pismena, {
                headers: {
                    "Token": this.$ls.getItem("token").value
                }
            })
            .then(response => {
                this.info = response.data
                this.dostal_jsem_data = true
            }).catch(e => {
                setTimeout(() => {this.$router.push('/404')}, 2000);
            });
    },
}
</script>

<template>
    <h1>
        <router-link class="tlacZpet" :to="'/lekce'">
            <img src="/icony/sipkaL.svg" alt="Zpět">
        </router-link>
        Lekce: {{$format(pismena) }}
    </h1>
    <div class="kontejnr" v-if="!info.error">
        <div v-if="dostal_jsem_data && info['cviceni'] !== null && info['cviceni'].length !== 0" v-for="(cviceni, index) in info['cviceni']">
            <h2>
                <router-link class="lekceBlok" :class="{ dokoncenyBlok: jeDokoncene(cviceni.id) }" v-if="cviceni.typ === 'nova'"
                    :to="'/lekce/' + pismena + '/' + index1(index)">
                    <h2>{{ index1(index) }}</h2>
                    <hr>
                    <h3>Nová písmenka</h3>
                    <img class="fajvkaVetsi" v-if="jeDokoncene(cviceni.id)" src="/icony/right.svg" alt="Dokonceno!">
                    <img class="playVetsi" v-else src="/icony/start.svg" alt="Začít lekci">
                </router-link>
                <router-link class="lekceBlok" :class="{ dokoncenyBlok: jeDokoncene(cviceni.id) }"
                    v-else-if="cviceni.typ === 'probrana'" :to="'/lekce/' + pismena + '/' + index1(index)">
                    <h2>{{ index1(index) }}</h2>
                    <hr>
                    <h3>Probraná písmenka</h3>
                    <img class="fajvkaVetsi" v-if="jeDokoncene(cviceni.id)" src="/icony/right.svg" alt="Dokonceno!">
                    <img class="playVetsi" v-else src="/icony/start.svg" alt="Začít lekci">
                </router-link>
                <router-link v-else class="lekceBlok" :class="{ dokoncenyBlok: jeDokoncene(cviceni.id) }"
                    :to="'/lekce/' + pismena + '/' + index1(index)">
                    <h2>{{ index1(index) }}</h2>
                    <hr>
                    <h3>Se slovy</h3>
                    <img class="fajvkaVetsi" v-if="jeDokoncene(cviceni.id)" src="/icony/right.svg" alt="Dokonceno!">
                    <img class="playVetsi" v-else src="/icony/start.svg" alt="Začít lekci">
                </router-link>
            </h2>
        </div>
        <p v-else-if="!dostal_jsem_data && info['cviceni'] == null"></p>
        <p v-else>Tato lekce zatím nemá žádná cvičení</p>
    </div>
    <p v-else>{{ info.error }}</p>
</template>

<style scoped>
.kontejnr {
    display: flex;
    gap: 15px;
    max-width: 700px;
    flex-wrap: wrap;
    justify-content: center;
}

.lekceBlok {
    color: var(--bila);
    display: flex;
    flex-direction: column;
    text-decoration: none;
    border-radius: 10px;
    width: 200px;
    background-color: var(--tmave-fialova);
    height: 240px;
    transition-duration: 0.2s;
    padding: 15px 15px 30px 15px;
}

.lekceBlok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.lekceBlok hr {
    width: 160px;
    align-self: center;
    margin: 5px;
}

.dokoncenyBlok {
    opacity: 80%;
}

.lekceBlok h3 {
    align-self: center;
    font-size: 24px;
    height: 100px;
}

.lekceBlok a {
    text-decoration: none;
    color: var(--bila);
    cursor: pointer;
}

h1 {
    display: inline-flex;
    position: relative;
    right: 25px; /* posunuti o pulku sipky */
}
</style>
