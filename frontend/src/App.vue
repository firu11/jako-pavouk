<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue';
import MenuLink from '@/components/MenuLink.vue';
import NotificationHost from '@/components/NotificationHost.vue';
import Tooltip from '@/components/Tooltip.vue';
import { mobil, prihlasen, role, tokenJmeno, uziv } from '@/stores';
import { oznameni, pridatOznameni } from '@/utils';
import { useHead } from '@unhead/vue';
import api, { ApiError } from '@/api';
import { useRouter } from 'vue-router';

useHead({
    titleTemplate: (title?: string) => (!title ? 'Psaní všemi deseti zdarma | Jako Pavouk' : `${title} | Jako Pavouk`),
});

const router = useRouter();
const mobilMenu = ref(false);
const sessionReady = ref(false);

const jmenoSpan = useTemplateRef<HTMLElement>('jmenoSpan');
const nadpisyDiv = useTemplateRef<HTMLElement>('nadpisyDiv');

interface SessionResponse {
    role: string;
    email: string;
    jmeno: string;
    jePotrebaVymenit: boolean;
}

async function nacistRelaci() {
    try {
        const response = await api.get<SessionResponse>('/token-expirace');

        if (response.data.jePotrebaVymenit) {
            localStorage.removeItem(tokenJmeno);
            role.value = 'basic';
            uziv.value = { email: '', jmeno: '' };
            prihlasen.value = false;
            await router.push('/prihlaseni');
            pridatOznameni('Z bezpečnostních důvodů jsme tě odhlásili ze sítě 🕸️', 8000);
            return;
        }

        role.value = response.data.role;
        uziv.value = { email: response.data.email, jmeno: response.data.jmeno };
        prihlasen.value = true;
    } catch (error) {
        if (error instanceof ApiError && (error.response.status === 401 || error.response.status === 418)) return;
        console.error(error);
        pridatOznameni('Chyba serveru');
    } finally {
        sessionReady.value = true;
    }
}

function upravitSirkuJmena() {
    const nameElement = jmenoSpan.value;
    const containerElement = nadpisyDiv.value;
    if (!nameElement || !containerElement) return;

    let fontSize = 24;
    nameElement.style.fontSize = `${fontSize}px`;
    while (nameElement.scrollWidth > containerElement.clientWidth && fontSize > 12) {
        fontSize -= 0.5;
        nameElement.style.fontSize = `${fontSize}px`;
    }
}

function naplanovatUpravuJmena() {
    void nextTick(upravitSirkuJmena);
}

function aktualizovatViewport() {
    mobil.value = document.body.clientWidth <= 900;
    naplanovatUpravuJmena();
}

function zavritMobilniMenuPriScrollu() {
    if (mobil.value) mobilMenu.value = false;
}

onMounted(() => {
    void nacistRelaci();
    aktualizovatViewport();
    window.addEventListener('resize', aktualizovatViewport);
    window.addEventListener('scroll', zavritMobilniMenuPriScrollu, { passive: true });
});

onUnmounted(() => {
    window.removeEventListener('resize', aktualizovatViewport);
    window.removeEventListener('scroll', zavritMobilniMenuPriScrollu);
});

function odhlasit(e: Event) {
    zavritDialog(e);
    localStorage.removeItem(tokenJmeno);
    api.post('/odhlaseni').catch(console.error);
    role.value = 'basic';
    prihlasen.value = false;
    router.push('/prihlaseni');

    uziv.value.email = '';
    uziv.value.jmeno = '';
}

const dialog1 = useTemplateRef('dialog1');
function otevritDialog(e: Event) {
    e.preventDefault();
    dialog1.value?.showModal();
}

function zavritDialog(e: Event) {
    e.preventDefault();
    dialog1.value?.close();
}

watch(() => uziv.value.jmeno, naplanovatUpravuJmena, { flush: 'post' });
</script>
<template>
    <div id="menu-mobilni-btn" @click="mobilMenu = !mobilMenu">
        <img id="menuIcon" src="./assets/icony/menu.svg" alt="Menu" width="40" height="40" />
    </div>
    <header :class="{ 'mobil-hidden': !mobilMenu }">
        <nav @click="mobilMenu = false">
            <!-- <img id="vanocni" src="./assets/vanocni.svg" alt="Vanoční světélka" /> -->
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
                            <img src="./assets/icony/nastaveni.svg" alt="" width="22" height="22" />
                        </div>
                    </Tooltip>
                    <Tooltip zprava="Statistiky" :sirka="100" :vzdalenost="-29" :vzdalenostX="75">
                        <div class="kulate-tlacitko" @click="router.push('/statistiky')">
                            <img src="./assets/icony/statistiky.svg" alt="" width="22" height="22" />
                        </div>
                    </Tooltip>
                    <Tooltip zprava="Odhlásit" :sirka="100" :vzdalenost="-29" :vzdalenostX="75">
                        <div class="kulate-tlacitko" @click="otevritDialog">
                            <img src="./assets/icony/odhlasit.svg" alt="" width="22" height="22" />
                        </div>
                    </Tooltip>
                </div>
                <img id="pavouk" src="./assets/pavoucekBezPozadi.svg" alt="uzivatel" width="181" height="114" />
            </div>
            <hr style="border: white solid 1px" />
            <div id="nadpisy" ref="nadpisyDiv">
                <span id="jmeno" ref="jmenoSpan">{{ uziv.jmeno }}</span>
                <span id="email">{{ uziv.email }}</span>
            </div>
        </div>
        <div v-else id="ucet" class="neprihlasen" @click="mobilMenu = !mobilMenu">
            <img id="pavouk" src="./assets/pavoucekBezPozadi.svg" alt="uzivatel" width="181" height="114" />
            <span>Nepřihlášený pavouk</span>
            <MenuLink jmeno="Přihlásit se" cesta="/prihlaseni" />
        </div>
    </header>
    <main id="view">
        <RouterView v-if="sessionReady" :key="$route.fullPath" />
    </main>

    <NotificationHost :notifications="oznameni" />

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

#dialog-kontejner > div {
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

#dialog-kontejner > div button {
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

#ucet #pavouk {
    max-width: calc(100% - 15px);
    margin-right: -15px;
    user-select: none;
}

#kontejner {
    display: flex;
    width: 100%;
    align-items: center;
}



header {
    display: flex;
    flex-direction: column;
    height: 100vh;
    height: 100dvh;
    position: fixed;
    padding: 10px;
    left: 0;
    gap: 10px;
    z-index: 1000 !important;
    transition: transform ease-in-out 0.3s;
}

nav {
    position: relative;
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

    #dialog-kontejner > div button {
        width: 120px;
    }
}

@media screen and (max-width: 500px) {
    dialog {
        margin-left: -40%;
        height: 170px;
    }
}

/*#vanocni {
    position: absolute;
    top: 3em;
    right: -5em;
    width: 300px;
    transform: rotate(55deg);
    user-select: none;
    pointer-events: none;
}*/
</style>
