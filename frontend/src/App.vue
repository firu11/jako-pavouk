<script setup lang="ts">
import { onMounted, ref } from 'vue';
import MenuLink from './components/MenuLink.vue';
import { prihlasen, tokenJmeno } from './stores';
import { checkTeapot, jeToRobot, getToken, oznameni, pridatOznameni } from './utils';
import { useHead } from 'unhead'
import axios from 'axios';
import router from './router';

useHead({
    titleTemplate: (title?: string) => title == "" || title == undefined ? "Psaní všemi deseti | Jako Pavouk" : `${title} | Jako Pavouk`
})

const mobilMenu = ref(false)

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
                router.push("/prihlaseni")
                pridatOznameni("Z bezpečnostních důvodů jsme tě odhlásili ze sítě 🕸️", 8000)
            } else {
                prihlasen.value = true
            }
        }).catch(e => {
            if (!checkTeapot(e)) {
                console.log(e)
                pridatOznameni()
            }
        })
    } else if (!jeToRobot(navigator.userAgent)) { //test jestli to neni bot
        axios.post("/navsteva")
    }
})

</script>

<template>
    <header>
        <div id="menuMobilniBtn" @click="mobilMenu = !mobilMenu"><img id="menuIcon" src="./assets/icony/menu.svg" alt="Menu"
                width="40" height="40"></div>
        <nav :class="{ mobilHidden: !mobilMenu }" @click="mobilMenu = !mobilMenu">
            <MenuLink jmeno="Domů" cesta="/" />
            <MenuLink jmeno="Jak psát" cesta="/jak-psat" />
            <MenuLink jmeno="Lekce" cesta="/lekce" />
            <MenuLink jmeno="Procvičování" cesta="/procvic" />
            <MenuLink jmeno="Test psaní" cesta="/test-psani" />
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
                <img src="./assets/icony/alert.svg" alt="Vykřičník">
                <span v-html="o.text"></span>
            </div>
        </TransitionGroup>
    </div>
</template>

<style scoped>
/* na tu animaci oznameni */
.list-move {
    transition: all 0.2s ease;
}

.list-enter-active,
.list-leave-active {
    transition: all 0.3s ease;
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
    overflow: hidden;
    padding: 20px;
    min-height: 100px;
    min-width: 410px;
    pointer-events: none;
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
    overflow: hidden;
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

#menuMobilniBtn {
    display: none;
}

#vanocni {
    position: absolute;
    top: 3em;
    right: -5em;
    width: 300px;
    transform: rotate(55deg);
    user-select: none;
    pointer-events: none
}

@media screen and (max-width: 1100px) {
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
        z-index: 1000;
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
}
</style>
