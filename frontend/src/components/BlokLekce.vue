<script>
export default {
    name: "blokLekce",
    props: {
        pismena: String,
        je_dokoncena: Boolean
    },
    methods: {
        prihlaste_se() {
            alert('Nejprve se prosím přihlašte')
        }
    },
    computed: {
        formatovany_pismena() {
            let vratit = "";
            for (let i = 0; i < this.pismena.length; i++) {
                vratit += i < this.pismena.length - 1 ? this.pismena.at(i) + ", " : this.pismena.at(i);
            }
            return vratit;
        }
    }
}
</script>

<template>
    <router-link class="lekceBlok" :class="{hotovoBlok: je_dokoncena}" v-if="$ls.getItem('token')" :to="'/lekce/' + pismena">
        <h2>Lekce: {{ formatovany_pismena }}</h2>
        <img class="fajvka" v-if="je_dokoncena" src="/icony/right.svg" alt="Dokonceno!">
    </router-link>
    <a v-else class="lekceBlok" @click="prihlaste_se"><h2>Lekce: {{ formatovany_pismena }}</h2></a>
</template>

<style scoped>
.lekceBlok {
    color: var(--bila);
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    padding: 12px 20px 12px 25px;
    text-decoration: none;
    border-radius: 10px;
    width: 500px;
    background-color: var(--tmave-fialova);
    height: 64px;
    transition-duration: 0.2s;
}

.lekceBlok:hover {
    background-color: var(--fialova);
    transition-duration: 0.2s;
}

.hotovoBlok {
    color: var(--seda);
    opacity: 80%;
}

.lekceBlok h2 {
    align-self: center;
    font-size: 24px;
}
</style>
