<script setup lang="ts">
import { onMounted } from 'vue';
import MenuLink from './components/MenuLink.vue';
import { prihlasen, tokenJmeno } from './stores';
import { getToken } from './utils';
import axios from 'axios';

onMounted(() => {
    if (getToken()) {
        axios.get("/token-expirace", {
            headers: {
                Authorization: `Bearer ${getToken()}`
            }
        }).then(response => {
            if (response.data.je_potreba_vymenit) {
                localStorage.removeItem(tokenJmeno)
                prihlasen.value = false
            } else {
                prihlasen.value = true
            }
        }).catch(e => {
            console.log(e)
        })
    }

})
</script>

<template>
    <header>
        <nav>
            <MenuLink jmeno="Domů" cesta="/" />
            <MenuLink jmeno="Lekce" cesta="/lekce" />
            <MenuLink jmeno="O nás" cesta="/o-nas" />
            <MenuLink jmeno="Podpořit projekt" cesta="/podporit" />
            <MenuLink v-if="!prihlasen" jmeno="Přihlásit se" cesta="/prihlaseni" />
            <MenuLink v-else jmeno="Můj účet" cesta="/ucet" />
        </nav>
    </header>
    <div id="view">
        <RouterView />
    </div>

    <img id="pavucina1" src="./assets/pavucina.svg" alt="Pavucina">
</template>

<style scoped>
nav {
    position: fixed;
    left: 10px;
    top: 10px;
    width: var(--sirka-menu);
    height: calc(100vh - 20px);
    flex-shrink: 0;
    border-radius: 10px;
    background-color: var(--tmave-fialova);
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
    width: 800px;
    z-index: -1000;
    opacity: 0.3;
}
</style>
