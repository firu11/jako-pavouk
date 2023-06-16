<script>
import axios from 'axios';

export default {
    name: "login",
    data() {
        return {
            email: this.email,
            heslo: this.heslo,
            info: "",
            spatnej_uzivatel: false,
            spatny_heslo: false,
        }
    },
    methods: {
        login(e) {
            e.preventDefault(); //aby se nerefreshla stranka

            if (!this.heslo) { //pokud uzivatel nic nenapsal
                this.spatny_heslo = true
            }
            if (!this.email) {
                this.spatnej_uzivatel = true
            }
            if (this.spatnej_uzivatel || this.spatny_heslo) return //nezkoušet ani

            axios.post('/prihlaseni', {
                "email": this.email,
                "heslo": this.heslo
            }).then(response => {
                this.info = response.data
                this.$ls.setItem("token", this.info.token, this.info.expiryTime)
                this.$router.push("/ucet")
            }).catch(e => {
                try {
                    if (e.response.status == 400 || e.response.status == 401) {
                        this.spatny_heslo = true
                        this.spatnej_uzivatel = true
                    } else {
                        alert("Něco se pokazilo na naší straně...")
                    }
                } catch {
                    console.log(e)
                }
            })

        },
        zmenaUziv() { // pokud zacnu znova psat tak zrusim znaceni spatnyho inputu
            this.spatnej_uzivatel = false
            this.spatny_heslo = false
        },
        zmenaHeslo() {
            this.spatny_heslo = false
            this.spatnej_uzivatel = false
        }
    }
}
</script>

<template>
    <h2>Přihlášení</h2>
    <form class="pruhledne">
        <h3 class="nadpis">Email:</h3>
        <input style="margin-bottom: 20px" :class="{ spatnej_input: spatnej_uzivatel }" :oninput="this.zmenaUziv"
            type="text" v-model="email" placeholder="Např: pepa@zdepa.cz">
        <h3 class="nadpis">Heslo:</h3>
        <input class="margin" :class="{ spatnej_input: spatny_heslo }" :oninput="this.zmenaHeslo" type="password"
            v-model="heslo" placeholder='Rozhodně ne "Pepa123"'>
        <button class="tlacitko" @click="login">Přihlásit</button>

    </form>
    <p>Nemáte ještě účet?
        <router-link to="/register">Registrace</router-link>
    </p>
</template>

<style scoped>
@import "@/loginRegisterForma.css";
</style>
