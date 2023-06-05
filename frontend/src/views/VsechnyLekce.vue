<script>
import BlokLekce from "@/components/BlokLekce.vue";
import axios from "axios";

export default {
    name: "VsechnyLekce",
    data() {
        return {
            data: {},
        }
    },
    components: {
        BlokLekce
    },
    beforeMount() {
        axios.get('/lekce', {
            headers: {
                "Token": this.$ls.getItem("token").value
            }
        }).then(response => {
            this.data = response.data
        })
    },
}
</script>

<template>
    <h1>Seznam lekcí</h1>
    <div id="lekce">
        <h2>Střední řada</h2>
        <BlokLekce v-if="data.hasOwnProperty('lekce')" v-for="lekce in data.lekce[0]" :pismena="lekce.pismena" :je_dokoncena="data.dokoncene.includes(lekce.id)"/>
        <h2>Spodní řada</h2>
        <BlokLekce v-if="data.hasOwnProperty('lekce')" v-for="lekce in data.lekce[1]" :pismena="lekce.pismena" :je_dokoncena="data.dokoncene.includes(lekce.id)"/>
        <h2>Horní řada</h2>
        <BlokLekce v-if="data.hasOwnProperty('lekce')" v-for="lekce in data.lekce[2]" :pismena="lekce.pismena" :je_dokoncena="data.dokoncene.includes(lekce.id)"/>
        <h2>Diakritika</h2>
        <BlokLekce v-if="data.hasOwnProperty('lekce')" v-for="lekce in data.lekce[3]" :pismena="lekce.pismena" :je_dokoncena="data.dokoncene.includes(lekce.id)"/>
    </div>
</template>

<style scoped>
#lekce {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}
</style>