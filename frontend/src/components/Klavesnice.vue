<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router';

const props = defineProps({
    aktivniPismeno: {
        type: String,
        default: "ß"
    },
    typ: {
        type: String,
        reqired: true
    },
})

const cesta = useRoute().path.split("/")[2]

const barvy = ["#6ada56", "#81bffc", "#fa5ca1", "#ff8800", "#6f86f7"]
let schema = [
    ["°;", "1+", "2ě", "3š", "4č", "5ř", "6ž", "7ý", "8á", "9í", "0é", "%=", "ˇ´", "⟵"],
    ["TAB", "Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P", "/ú", "()", "'¨"],
    ["CapsLock", "A", "S", "D", "F", "G", "H", "J", "K", "L", '"ů', "!§", "Enter ↵"],
    ["Shift", "Y", "X", "C", "V", "B", "N", "M", "?,", ":.", "_-", "Shift"],
    ["  ", "", "", "", "______", "", "", "", "∧∨", ""]
]
const delkaKlaves: { [id: string]: number } = { "⟵": 3, "Shift": 1, "Enter ↵": 1, "CapsLock": 1, "TAB": 1, "______": 24, "  ": 2, "Ctrl": 2 }
const prstoklad: { [id: string]: string[] } = {
    "P_Ukaz": [barvy[0], "J", "H", "U", "Z", "N", "M", "7ý", "8á"],
    "L_Ukaz": [barvy[1], "G", "T", "R", "F", "V", "B", "5ř", "6ž"],
    "P_Pros": [barvy[2], "K", "I", "?,", "9í"],
    "L_Pros": [barvy[2], "D", "E", "C", "4č"],
    "P_Prs": [barvy[3], "O", "L", ":.", "0é"],
    "L_Prs": [barvy[3], "X", "S", "W", "3š"],
    "P_Mali": [barvy[4], '"ů', "P", "_-", '%=', 'ˇ´', '⟵', '()', '/ú', "'¨", '!§', 'Enter ↵', 'Shift'],
    "L_Mali": [barvy[4], "Shift", "A", "Q", "Y", "1+", "°;", "2ě", "TAB", "CapsLock", "Ctrl"],
    "Palce": ["#bc73ff", "______", "Alt"]
}
const shiftSviti = ref(false)

watch(() => props.typ, (ted) => {
    if (ted == "qwerty") {
        schema[1][6] = "Y"
        schema[3][1] = "Z"
        prstoklad.P_Ukaz[4] = "Y"
        prstoklad.L_Mali[4] = "Z"
    } else {
        schema[1][6] = "Z"
        schema[3][1] = "Y"
        prstoklad.P_Ukaz[4] = "Z"
        prstoklad.L_Mali[4] = "Y"
    }
})

if (cesta == "závorky" || cesta == "operátory") {
    schema[4][0] = "Ctrl"
    schema[4][3] = "Alt"
}

function tlacPismeno(cislo: number, tlacitko: string) {
    try {
        if (tlacitko.length === 2) return tlacitko.at(cislo)
        else if (tlacitko.length === 1 && cislo === 0) return tlacitko.at(0)
        else if (tlacitko.length >= 2 && cislo === 0) return tlacitko
    } catch { // kvůli starším prohlížečům (koukám na tebe safari <14.0)
        if (tlacitko.length === 2) return tlacitko.charAt(cislo)
        else if (tlacitko.length === 1 && cislo === 0) return tlacitko.charAt(0)
        else if (tlacitko.length >= 2 && cislo === 0) return tlacitko
    }
}

function oznacene(tlacitko: string) {
    if (props.aktivniPismeno == "") return false
    let velky = props.aktivniPismeno.toUpperCase() === props.aktivniPismeno
    let pismeno = props.aktivniPismeno.toLowerCase()
    if (pismeno.normalize('NFD').length === 2 && velky) {
        if (pismeno.normalize("NFD").replace(/\p{Diacritic}/gu, "") === tlacitko.toLowerCase()) {
            return true
        }
        if (pismeno === "ů" && tlacitko === "°;") return true
        if (pismeno !== "ů" && tlacitko === "ˇ´") return true
        return false
    }
    if ((pismeno === ' ' && tlacitko === '______') || (tlacitko.length === 1 && tlacitko.toLowerCase() === pismeno) || (tlacitko.length === 2 && tlacitko.toLowerCase().includes(pismeno) && tlacitko !== '  ')) {
        return true
    } else if ('óťňď'.includes(pismeno)) {
        if (tlacitko == 'ˇ´') return true
        if ((tlacitko == 'O' && pismeno == 'ó') || (tlacitko == 'T' && pismeno == 'ť') || (tlacitko == 'N' && pismeno == 'ň') || (tlacitko == 'D' && pismeno == 'ď')) return true
    } else if ("[]{}<>*".includes(pismeno)) {
        if (tlacitko === "Ctrl" || tlacitko === "Alt") return true
        if (props.typ === "qwertz") {
            if (pismeno === "[" && tlacitko === "F") return true
            if (pismeno === "]" && tlacitko === "G") return true
            if (pismeno === "{" && tlacitko === "B") return true
            if (pismeno === "}" && tlacitko === "N") return true
            if (pismeno == "*" && tlacitko == "_-") return true
        } else {
            if (pismeno === "[" && tlacitko === "/ú") return true
            if (pismeno === "]" && tlacitko === "()") return true
            if (pismeno === "{" && tlacitko === "/ú") return true
            if (pismeno === "}" && tlacitko === "()") return true
            if (pismeno == "*" && tlacitko == "8á") return true
        }
        if (pismeno === "<" && tlacitko === "?,") return true
        if (pismeno === ">" && tlacitko === ":.") return true
    } else {
        return false
    }
}

function barva(tlacitko: string) {
    let pismeno = props.aktivniPismeno
    if (tlacitko === 'Shift') {
        if (potrebujeShift(pismeno) && pismeno !== " ") {
            shiftSviti.value = true
        } else shiftSviti.value = false

        return prstoklad['P_Mali'][0]

    } else if (tlacitko.length === 2 && tlacitko.toLowerCase().includes(pismeno.toLowerCase())) {
        for (let prst in prstoklad) {
            for (let tla in prstoklad[prst]) {
                if (prstoklad[prst][tla].toLowerCase().charAt(0) === tlacitko.toLowerCase().charAt(0)) {
                    return prstoklad[prst][0]
                } else if (prstoklad[prst][tla].toLowerCase().charAt(1) === tlacitko.toLowerCase().charAt(1)) {
                    return prstoklad[prst][0]
                }
            }
        }
    } else {
        for (let prst in prstoklad) {
            for (let tla in prstoklad[prst]) {
                if (prstoklad[prst][tla].toLowerCase() === tlacitko.toLowerCase()) {
                    return prstoklad[prst][0]
                }
            }
        }
    }
}

function delkaTlacitka(tlacitko: string) {
    if (delkaKlaves[tlacitko] === undefined) {
        return 0
    }
    return delkaKlaves[tlacitko]
}

function potrebujeShift(pismeno: string) {
    if (['"', '/', '?', ':', '_', '!', '(', '%', 'ˇ', '°', 'ť', 'Ť', 'ď', 'Ď', 'ň', 'Ň', 'Ě', 'Š', 'Č', 'Ř', 'Ž', 'Ý', 'Á', 'Í', 'É', 'Ú', 'Ů'].includes(pismeno)) {
        return true
    } else if (/^\d$/.test(pismeno)) { // jestli to je cislo
        return true
    }

    if (props.typ === "qwertz" && "[]{}<>*".includes(pismeno)) return false
    else if (props.typ === "qwerty" && "<>[]*".includes(pismeno)) return false

    return (pismeno === pismeno.toUpperCase() && !["+", "=", "-", ".", ",", "§", ")", "´", ";"].includes(pismeno))
}
</script>

<template>
    <div id="klavesnice">
        <div class="radek" v-for="radek in schema">
            <div v-for="tlacitko in radek" class="klavesa"
                :class="{ oznacenaKlavesa: oznacene(tlacitko) || (tlacitko === 'Shift' && shiftSviti), fjPodtrzeni: tlacitko === 'F' || tlacitko === 'J' }"
                :style="{ backgroundColor: barva(tlacitko), flexGrow: delkaTlacitka(tlacitko) }">

                <div v-if="tlacitko !== '∧∨'" :style="{ color: '#000' }">
                    {{ tlacPismeno(0, tlacitko) }} <br>
                    {{ tlacPismeno(1, tlacitko) }}
                </div>
                <div v-else id="sipky">
                    <div class="klavesa" style="height: 18px"></div>
                    <div class="klavesa" style="height: 18px"></div>
                </div>
            </div>

        </div>
    </div>
</template>

<style scoped>
.klavesa {
    width: 40px;
    height: 40px;
    background-color: var(--fialova);
    border-radius: 5px;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: 0.1s;
    filter: brightness(0.9);
}

.oznacenaKlavesa {
    border: #fff solid 3px;
    transition: 0.1s;
    filter: brightness(1.6) saturate(2);
}

.oznacenaKlavesa div {
    font-weight: 800 !important;
}

.klavesa div {
    font-weight: 500;
    transition: 0.1s;
}

#klavesnice {
    display: flex;
    flex-direction: column;
    gap: 4px;
    background-color: var(--tmave-fialova);
    padding: 10px;
    border-radius: 10px;
    font-size: 0.8em;
    line-height: 1.3em;
    width: 675px;
    margin-top: 25px;
    user-select: none;
    transition: filter 0.2s;
}

.radek {
    display: flex;
    gap: 5px;
}

#sipky {
    display: flex;
    flex-direction: column;
    gap: 4px;
    background-color: var(--tmave-fialova);
    filter: brightness(1.1);
}

.fjPodtrzeni div {
    border-bottom: 1px solid black !important;
    width: 10px;
    padding-top: 2px;
}
</style>