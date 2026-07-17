<script setup lang="ts">
import { onMounted, onUnmounted, ref, useTemplateRef, watch } from 'vue';
import { pridatOznameni, postKlavesnice } from '@/utils';
import { mobil } from '@/stores';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter();
const route = useRoute();

const klavesnice = ref(false);
const mameJi = ref(false);
const rucne = ref(mobil.value);
const img = useTemplateRef<HTMLImageElement>('img');
let zoomTimeout: ReturnType<typeof setTimeout> | undefined;

watch(mobil, (isMobile) => {
    if (isMobile) vybratRucne();
});

onMounted(() => {
    document.addEventListener('keypress', click);
    zoomTimeout = setTimeout(zoomIn, 300);
});

onUnmounted(() => {
    if (zoomTimeout !== undefined) clearTimeout(zoomTimeout);
    document.removeEventListener('keypress', click);
});

function click(e: KeyboardEvent) {
    if (e.key.toLowerCase() == 'z') {
        klavesnice.value = false;
        mameJi.value = true;
        document.removeEventListener('keypress', click);
    } else if (e.key.toLowerCase() == 'y') {
        klavesnice.value = true;
        mameJi.value = true;
        document.removeEventListener('keypress', click);
    } else {
        pridatOznameni('Asi mačkáš špatné tlačítko');
    }
}

function vybratRucne() {
    rucne.value = true;
    if (img.value) img.value.style.transform = 'none';
}

function vybratAutomaticky() {
    rucne.value = false;
    zoomIn();
}

function zoomIn() {
    if (rucne.value || !img.value) return;
    img.value.style.transform = 'scale(1.6) translateY(36px)';
}

function zoomOut() {
    if (rucne.value || !img.value) return;
    img.value.style.transition = '0.8s';
    img.value.style.transform = 'none';
}

function potvrdit() {
    postKlavesnice(klavesnice.value);

    const queryDestination = Array.isArray(route.query.kam) ? route.query.kam[0] : route.query.kam;
    const destination = queryDestination === 'nastaveni' || queryDestination === 'statistiky' ? `/${queryDestination}` : '/';
    void router.push(destination);
}
</script>
<template>
    <h1>Výběr klávesnice</h1>

    <div id="popup">
        <div>
            <img @mouseenter="zoomOut" @mouseleave="zoomIn" src="../assets/vyberKlavesnice.png" alt="Klavesnice" width="600" ref="img" />
        </div>

        <p v-if="!mameJi && !rucne">Zmáčkni prosím toto tlačítko na své klávesnici.</p>
        <p v-else-if="mameJi">
            Super! Tvoje rozložení je <b>{{ klavesnice ? 'Qwerty' : 'Qwertz' }}.</b>
        </p>
        <button v-if="mameJi" class="tlacitko" @click="potvrdit" style="margin-top: -14px">Pokračovat</button>

        <div v-if="rucne && !mameJi" id="tlacitka">
            <button
                class="tlacitko"
                @click="
                    klavesnice = false;
                    mameJi = true;
                "
            >
                Qwertz
            </button>
            <button
                class="tlacitko"
                @click="
                    klavesnice = true;
                    mameJi = true;
                "
            >
                Qwerty
            </button>
        </div>
    </div>

    <p id="dole" v-if="!rucne && !mameJi && route.query['kam'] != 'nastaveni'" @click="vybratRucne">Vybrat rozložení ručně</p>
    <p id="dole" v-else-if="!mameJi && route.query['kam'] != 'nastaveni'" @click="vybratAutomaticky">Zpět k automatickému výběru</p>
</template>
<style scoped>
#tlacitka {
    display: flex;
    gap: 20px;
}

#dole {
    text-decoration: underline;
    margin: 10px;
    cursor: pointer;
}

#popup > div {
    overflow: hidden;
    border-radius: 8px;
}

#popup {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 20px;
    margin-top: 20px;
    background-color: var(--tmave-fialova);
    padding: 20px;
    border-radius: 10px;
}

#popup > div > img {
    transition: 2s ease-in-out;
}

#popup > p {
    font-size: 26px;
    line-height: 40px;
}

h1 {
    margin: 40px 0 0 0;
    font-size: 32px;
}

.tlacitko {
    margin: 0;
}

@media screen and (max-width: 1100px) {
    #popup {
        max-width: 90vw;
    }

    #popup > div {
        max-width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
}
</style>
