<script setup lang="ts">
import { onMounted, ref } from 'vue';
import MenuLink from './components/MenuLink.vue';
import { prihlasen, tokenJmeno } from './stores';
import { checkTeapot, getToken, oznameni, pridatOznameni } from './utils';
import { useHead } from 'unhead'
import axios from 'axios';

useHead({
    titleTemplate: (title?: string) => title != "" ? `${title} | Jako Pavouk` : "Jako Pavouk | Psaní všemi deseti",
    meta: [
        {
            name: "description",
            content: "Webová aplikace pro výuku psaní všemi deseti",
        },
        {
            name: "keywords",
            content: "psani, psaní, všemi, vsemi, všema, vsema, deseti",
        },
    ]
})

onMounted(() => {
    console.log("%cCo sem koukáš koloušku?", "color: white; font-size: x-large"); // troulin
    if (getToken()) {
        axios.get("/token-expirace", {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }).then(response => {
            if (response.data.jePotrebaVymenit) {
                localStorage.removeItem(tokenJmeno)
                prihlasen.value = false
                pridatOznameni("Z bezpečnostních důvodů jste byli odhlášeni", 8000)
            } else {
                prihlasen.value = true
            }
        }).catch(e => {
            if (!checkTeapot(e)) {
                console.log(e)
                pridatOznameni()
            }
        })
    }
})

const mobilMenu = ref(false)

</script>

<template>
    <header>
        <div id="menuMobilniBtn" @click="mobilMenu = !mobilMenu"><img id="menuIcon" src="./assets/icony/menu.svg"
                alt="Menu"></div>
        <nav :class="{ mobilHidden: !mobilMenu }" @click="mobilMenu = !mobilMenu">
            <MenuLink jmeno="Domů" cesta="/" />
            <MenuLink jmeno="Jak psát" cesta="/jak-psat" />
            <MenuLink jmeno="Lekce" cesta="/lekce" />
            <MenuLink jmeno="O nás" cesta="/o-nas" />
            <MenuLink v-if="!prihlasen" jmeno="Přihlásit se" cesta="/prihlaseni" />
            <MenuLink v-else jmeno="Můj účet" cesta="/ucet" />
        </nav>
    </header>
    <div id="view">
        <RouterView :key="$route.fullPath" />
    </div>

    <div id="alerty">
        <TransitionGroup name="list">
            <div v-for="(o, i) in oznameni" class="alert" :key="i">
                <img src="./assets/icony/alret.svg" alt="Vykřičník">
                {{ o.text }}
            </div>
        </TransitionGroup>
    </div>

    <img id="pavucina1" src="./assets/pavucina.svg" alt="Pavucina">
</template>

<style scoped>
/* na tu animaci oznameni */
.list-enter-active,
.list-leave-active {
    transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
    opacity: 0;
    transform: translateX(30px);
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
    overflow: hidden;
    padding: 20px;
}

.alert {
    min-height: 60px;
    background-color: var(--tmave-fialova);
    min-width: 100px;
    max-width: 390px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 5px;
    padding: 10px 20px 10px 20px;
    gap: 20px;
    box-shadow: 0px 0px 10px 2px rgba(0, 0, 0, 0.75);
}

nav {
    position: fixed;
    left: 10px;
    top: 10px;
    width: var(--sirka-menu);
    height: calc(100vh - 20px);
    flex-shrink: 0;
    border-radius: 10px;
    background-color: var(--tmave-fialova);
    transition: ease-in-out 0.3s;
}

#view {
    padding-top: 30px;
    margin-left: calc(var(--sirka-menu) + 10px);
    margin-bottom: 50px;
    text-align: center;
    width: 720px;
    display: flex;
    flex-direction: column;
    align-items: center;
}

#pavucina1 {
    position: absolute;
    top: 0;
    right: 0;
    transform: rotate(180deg);
    width: 700px;
    z-index: -1000;
    opacity: 0.3;
}

#menuMobilniBtn {
    display: none;
}

@media screen and (max-width: 1000px) {
    .mobilHidden {
        transform: translateX(-250px);
        transition: ease-in-out 0.3s;
    }

    #menuMobilniBtn {
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
    }

    nav {
        border-radius: 10px;
        background-color: var(--tmave-fialova);
        display: flex;
        flex-direction: column;
        z-index: 10;
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

    #pavucina1 {
        position: absolute;
        top: 0;
        right: 0;
        transform: rotate(180deg);
        width: 400px;
        z-index: -1000;
        opacity: 0.3;
    }
}
</style>
