<script>

export default {
    name: 'app',
    mounted() {
        // tady bych chtěl checknout jestli token je furt funkcni
    },
}
</script>

<template>
    <header>
        <nav class="pruhledne">
            <RouterLink to="/">
                <div class="hoverBoard">
                    <svg class="dot" height="10" width="15">
                        <circle cx="5" cy="5" r="5" fill="white" />
                    </svg>
                    <p>Domů</p>
                </div>
            </RouterLink>
            <RouterLink to="/lekce">
                <div class="hoverBoard">
                    <svg class="dot" height="10" width="15">
                        <circle cx="5" cy="5" r="5" fill="white" />
                    </svg>
                    <p>Lekce</p>
                </div>
            </RouterLink>
            <RouterLink to="/about">
                <div class="hoverBoard">
                    <svg class="dot" height="10" width="15">
                        <circle cx="5" cy="5" r="5" fill="white" />
                    </svg>
                    <p>O nás</p>
                </div>
            </RouterLink>
            <RouterLink v-if="$ls.getItem('token').value" to="/ucet">
                <div class="hoverBoard">
                    <svg class="dot" height="10" width="15">
                        <circle cx="5" cy="5" r="5" fill="white" />
                    </svg>
                    <p>Můj účet</p>
                </div>
            </RouterLink>
            <RouterLink v-else to="/login">
                <div class="hoverBoard">
                    <svg class="dot" height="10" width="15">
                        <circle cx="5" cy="5" r="5" fill="white" />
                    </svg>
                    <p>Přihlásit se</p>
                </div>
            </RouterLink>
        </nav>
    </header>

    <img id="pavouk" src="/icony/pavouk.png" alt="Pavouk">

    <div id="kontent">
        <RouterView />
    </div>

    <img id="pavucina1" src="/pavucina.svg" alt="Pavucina">
</template>

<style scoped>
#pavouk {
    position: absolute;
    left: 120px;
    width: 100px;
    pointer-events: none;
}

#kontent {
    padding-top: 30px;
    margin-left: calc(var(--sirka-menu) + 10px);
    margin-bottom: 50px;
    text-align: center;
    width: 720px;
    display: flex;
    flex-direction: column;
    align-items: center;
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
}

nav a {
    color: var(--bila);
    display: block;
    padding: 5px;
    text-decoration: none;
    cursor: pointer;
    font-size: 17px;
}

nav a p {
    display: inline-block;
    margin-left: 3px;
    position: relative;
    left: -15px;
    transition: 0.3s;
}

.router-link-active p {
    left: 0;
    transition: 0.3s;
}

nav a:hover .hoverBoard {
    background-color: var(--fialova);
    width: calc(var(--sirka-menu) - 10px);
    transition: background-color 0.3s;
}

.hoverBoard {
    padding: 10px;
    border-radius: 8px;
    transition: background-color 0.3s;
}

.dot {
    opacity: 0;
    transition: 0.2s;
}

.router-link-active .dot {
    opacity: 100;
    margin-bottom: 1px;
    transition: 0.2s;
}

div {
    /* google button */
    padding: 10px 15px 15px 15px;
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
