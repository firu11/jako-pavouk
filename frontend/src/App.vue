<script setup lang="ts">
import { onMounted, ref, useTemplateRef, watch } from "vue";
import MenuLink from "./components/MenuLink.vue";
import { mobil, prihlasen, role, tokenJmeno, uziv } from "./stores";
import { jeToRobot, getToken, oznameni, pridatOznameni } from "./utils";
import { useHead } from "unhead";
import axios from "axios";
import { useRouter } from "vue-router";
import Tooltip from "./components/Tooltip.vue";

useHead({
    titleTemplate: (title?: string) => title == "" || title == undefined ? "Psaní všemi deseti zdarma | Jako Pavouk" : `${title} | Jako Pavouk`
})

const router = useRouter()
const mobilMenu = ref(false)

const jmenoSpan = useTemplateRef('jmenoSpan')
const nadpisyDiv = useTemplateRef('nadpisyDiv')

onMounted(() => {
    if (getToken()) {
        axios.get("/token-expirace", {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }).then(response => {
            role.value = response.data.role
            uziv.value.email = response.data.email
            uziv.value.jmeno = response.data.jmeno

            if (response.data.jePotrebaVymenit) {
                localStorage.removeItem(tokenJmeno)
                prihlasen.value = false
                router.push("/prihlaseni")
                pridatOznameni("Z bezpečnostních důvodů jsme tě odhlásili ze sítě 🕸️", 8000)
            } else {
                prihlasen.value = true
            }
        }).catch(e => {
            if (!(e.response && e.response.status == 418)) {
                console.log(e)
                pridatOznameni("Chyba serveru")
            }
        })
    } else if (!jeToRobot(navigator.userAgent) && window.location.host !== "localhost:5173") { //test jestli to neni bot + počítají se jen na produkci
        axios.post("/navsteva")
    }

    window.addEventListener("resize", function () {
        mobil.value = document.body.clientWidth <= 900
    })
})

async function upravitSirkuJmena() {
    if (jmenoSpan.value == undefined || nadpisyDiv.value == undefined) {
        setTimeout(upravitSirkuJmena, 10)
        return
    }

    let velikost = 24
    jmenoSpan.value.style.fontSize = `${velikost}px`
    while (jmenoSpan.value.clientWidth! > nadpisyDiv.value.clientWidth) {
        velikost -= 0.5
        jmenoSpan.value.style.fontSize = `${velikost}px`
        if (velikost <= 0) break
    }
    console.log("velikost jmena v px:", velikost)
}

function odhlasit(e: Event) {
    zavritDialog(e)
    localStorage.removeItem(tokenJmeno)
    role.value = "basic"
    prihlasen.value = false
    router.push("/prihlaseni")

    uziv.value.email = ""
    uziv.value.jmeno = ""
}

const dialog1 = useTemplateRef("dialog1")
function otevritDialog(e: Event) {
    e.preventDefault()
    dialog1.value?.showModal()
}

function zavritDialog(e: Event) {
    e.preventDefault()
    dialog1.value?.close()
}

watch(() => uziv.value.jmeno, () => {
    setTimeout(upravitSirkuJmena, 1)
})

</script>

<template>
    <div id="menu-mobilni-btn" @click="mobilMenu = !mobilMenu">
        <img id="menuIcon" src="./assets/icony/menu.svg" alt="Menu" width="40" height="40">
    </div>
    <header :class="{ 'mobil-hidden': !mobilMenu }">
        <nav @click="mobilMenu = !mobilMenu">
            <MenuLink jmeno="Domů" cesta="/" />
            <MenuLink jmeno="Jak psát" cesta="/jak-psat" />
            <MenuLink jmeno="Kurz" cesta="/kurz" />
            <MenuLink jmeno="Procvičování" cesta="/procvic" />
            <MenuLink v-if="!mobil" jmeno="Test psaní" cesta="/test-psani" />
            <MenuLink v-if="role == 'student'" jmeno="Škola" cesta="/trida" />
            <MenuLink v-else-if="role == 'ucitel'" jmeno="Škola" cesta="/skola" />
            <MenuLink jmeno="O nás" cesta="/o-nas" />
        </nav>
        <div v-if="prihlasen && uziv.jmeno != ''" id="ucet" @click="mobilMenu = !mobilMenu">
            <div id="kontejner">
                <div id="tlacitka">
                    <Tooltip zprava="Nastavení účtu" :sirka="100" :vzdalenost="-36" :vzdalenostX="75">
                        <div class="kulate-tlacitko" @click="router.push('/nastaveni')">
                            <img src="./assets/icony/nastaveni.svg" alt="" width="22" height="22">
                        </div>
                    </Tooltip>
                    <Tooltip zprava="Statitiky" :sirka="100" :vzdalenost="-29" :vzdalenostX="75">
                        <div class="kulate-tlacitko" @click="router.push('/statistiky')">
                            <img src="./assets/icony/statistiky.svg" alt="" width="22" height="22">
                        </div>
                    </Tooltip>
                    <Tooltip zprava="Odhlásit" :sirka="100" :vzdalenost="-29" :vzdalenostX="75">
                        <div class="kulate-tlacitko" @click="otevritDialog">
                            <img src="./assets/icony/odhlasit.svg" alt="" width="22" height="22">
                        </div>
                    </Tooltip>
                </div>
                <img src="./assets/pavoucekBezPozadi.svg" alt="uzivatel" width="181" height="114">
            </div>
            <hr style="border: white solid 1px;">
            <div id="nadpisy" ref="nadpisyDiv">
                <span id="jmeno" ref="jmenoSpan">{{ uziv.jmeno }}</span>
                <span id="email">{{ uziv.email }}</span>
            </div>
        </div>
        <div v-else id="ucet" class="neprihlasen" @click="mobilMenu = !mobilMenu">
            <img src="./assets/pavoucekBezPozadi.svg" alt="uzivatel" width="181" height="114">
            <span>Nepřihlášný pavouk</span>
            <MenuLink jmeno="Přihlásit se" cesta="/prihlaseni" />
        </div>
    </header>
    <main id="view">
        <RouterView :key="$route.fullPath" />
    </main>

    <div id="alerty">
        <TransitionGroup name="list">
            <div v-for="(o, i) in oznameni" class="alert" :key="i">
                <img v-if="o.typ == 'vykricnik'" src="./assets/icony/alert.svg" alt="Vykřičník">
                <img v-else-if="o.typ == 'copy'" src="./assets/icony/copy.svg" alt="Zkopírováno">
                <span v-html="o.text"></span>
            </div>
        </TransitionGroup>
    </div>

    <dialog ref="dialog1">
        <div id="dialog-kontejner">
            <h2>Opravdu se chceš odhlásit?</h2>
            <div>
                <button class="cervene-tlacitko" @click="odhlasit">Odhlásit se</button>
                <button class="tlacitko" @click="zavritDialog">Zrušit</button>
            </div>
        </div>
    </dialog>
</template>

<style scoped>
#dialog-kontejner {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.2em;
}

#dialog-kontejner>div {
    display: flex;
    justify-content: center;
    gap: 1em;
}

dialog {
    width: 410px;
    height: 140px;
    margin-left: -205px;
    margin-top: -70px;
    padding: 1.4em;
}

#dialog-kontejner>div button {
    margin: 0;
}

.neprihlasen {
    padding: 15px 0 0 0 !important;
    align-items: center;
}

.neprihlasen span {
    font-size: 18px;
    margin-bottom: 10px;
    font-weight: 500;
}

/* eslint-disable-next-line vue-scoped-css/no-unused-selector */
.neprihlasen a {
    width: 100%;
}

#ucet {
    padding: 15px 17px;
    margin-top: auto;
    background-color: var(--tmave-fialova);
    border-radius: 10px;
    aspect-ratio: 1/1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    transition: transform ease-in-out 0.3s;
    width: var(--sirka-menu);
}

#ucet #nadpisy {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
}

#ucet #jmeno {
    font-size: 24px;
    font-weight: 600;
}

#ucet #email {
    max-width: 100%;
    font-size: 16px;
    overflow: hidden;
    text-overflow: ellipsis;
}

#ucet #tlacitka {
    display: flex;
    flex-direction: column;
    gap: 5px;
}

#ucet #tlacitka .kulate-tlacitko {
    border-radius: 4px;
    padding: 5px;
    width: 32px;
    height: 32px;
    transition: 0.3s;
    cursor: pointer;
}

#ucet #tlacitka .kulate-tlacitko:hover {
    background-color: var(--fialova);
}

#kontejner {
    display: flex;
    width: 100%;
    align-items: center;
}

#kontejner>img {
    max-width: calc(100% - 15px);
    margin-right: -15px;
    user-select: none;
}

/* na tu animaci oznameni */
.list-move {
    transition: all 0.2s ease;
}

.list-enter-active,
.list-leave-active {
    transition: all 0.1s ease;
}

.list-enter-from,
.list-leave-to {
    opacity: 0;
    transform: translateX(50px);
}

.list-leave-active {
    position: absolute;
}

#alerty {
    position: fixed;
    right: 0;
    bottom: 0;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: end;
    gap: 10px;
    padding: 20px;
    min-height: 100px;
    pointer-events: none;
    width: 100vw;
    z-index: 1000;
}

.alert {
    min-height: 60px;
    background-color: var(--tmave-fialova);
    min-width: 100px;
    max-width: min(85%, 330px);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 5px;
    padding: 10px 20px 10px 20px;
    gap: 15px;
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
}

.alert img {
    width: 24px;
}

header {
    display: flex;
    flex-direction: column;
    height: 100%;
    position: fixed;
    padding: 10px;
    left: 0;
    gap: 10px;
    z-index: 1000 !important;
    transition: transform ease-in-out 0.3s;
}

nav {
    flex-grow: 10;
    width: var(--sirka-menu);
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

#menu-mobilni-btn {
    display: none;
}

@media screen and (max-width: 1100px) {
    .mobil-hidden {
        transform: translateX(-250px);
        transition: transform ease-in-out 0.3s;
    }

    #menu-mobilni-btn {
        background-color: var(--tmave-fialova);
        border-radius: 100px;
        padding: 10px;
        display: block;
        position: fixed;
        right: 10px;
        top: 10px;
        width: 60px;
        height: 60px;
        box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
        z-index: 1000;
    }

    nav,
    #ucet {
        box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
    }

    #view {
        padding-top: 30px;
        margin-left: 0;
        margin-bottom: 50px;
        text-align: center;
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }

    dialog {
        width: min(80%, 400px);
        margin-top: -70px;
    }

    #dialog-kontejner>div button {
        width: 120px;
    }
}

@media screen and (max-width: 500px) {
    dialog {
        margin-left: -40%;
        height: 170px;
    }
}
</style>
