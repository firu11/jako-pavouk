<script setup lang="ts">
import { onMounted } from 'vue';
import MenuLink from './components/MenuLink.vue';
import { prihlasen, token_jmeno } from './stores';

onMounted(() => {
    const token = localStorage.getItem(token_jmeno)
    if (token != null) {
        prihlasen.value = true
    }
})
</script>

<template>
	<header>
		<nav>
			<MenuLink jmeno="Domů" cesta="/"/>
			<MenuLink jmeno="Lekce" cesta="/lekce"/>
			<MenuLink jmeno="O nás" cesta="/o-nas"/>
            <MenuLink v-if="!prihlasen" jmeno="Prihlášení" cesta="/prihlasit"/>
            <MenuLink v-else jmeno="Účet" cesta="/ucet"/>
			<!-- 
			
			<MenuLink v-if="true" jmeno="Prihlásit se" cesta="/o-nas"/> -->
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
