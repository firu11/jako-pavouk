<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps({
    aktivniPismeno: {
        type: String,
        default: ""
    },
})

let barvy = ["#6ada56", "#81bffc", "#fa5ca1", "#ff8800", "#6f86f7"]
let schema = [
    ["°;", "1+", "2ě", "3š", "4č", "5ř", "6ž", "7ý", "8á", "9í", "0é", "%=", "ˇ´", "⟵"],
    ["TAB", "Q", "W", "E", "R", "T", "Z", "U", "I", "O", "P", "/ú", "()", "'¨"],
    ["CapsLock", "A", "S", "D", "F", "G", "H", "J", "K", "L", '"ů', "!§", "Enter ↵"],
    ["Shift", "Y", "X", "C", "V", "B", "N", "M", "?,", ":.", "_-", "Shift"],
    ["  ", "", "", "", "______", "", "", "", "^v", ""]
]
const delkaKlaves: {[id: string]: number} = { "⟵": 3, "Shift": 1, "Enter ↵": 1, "CapsLock": 1, "TAB": 1, "______": 24, "  ": 2 }
const prstoklad: {[id: string]: string[]} = {
    "P_Ukaz": [barvy[0], "J", "H", "U", "Z", "N", "M", "7ý", "6ž"],
    "L_Ukaz": [barvy[1], "G", "T", "R", "F", "V", "B", "5ř", "4č"],
    "P_Pros": [barvy[2], "K", "I", "?,", "8á"],
    "L_Pros": [barvy[2], "D", "E", "C", "3š"],
    "P_Prs": [barvy[3], "O", "L", ":.", "9í"],
    "L_Prs": [barvy[3], "X", "S", "W", "2ě"],
    "P_Mali": [barvy[4], '"ů', "P", "_-", "0é", '%=', 'ˇ´', '⟵', '()', '/ú', "'¨", '!§', 'Enter ↵', 'Shift'],
    "L_Mali": [barvy[4], "Shift", "A", "Q", "Y", "1+", "°;", "TAB", "CapsLock", "Ctrl"],
    "Palce": ["#bc73ff", "______"]
}
const shiftSviti = ref(false)


function tlacPismeno(cislo: number, tlacitko: string) {
    if (tlacitko.length === 2) return tlacitko.at(cislo)
    else if (tlacitko.length === 1 && cislo === 0) return tlacitko.at(0)
    else if (tlacitko.length >= 2 && cislo === 0) return tlacitko
}

function oznacene(tlacitko: string) {
    let pismeno = props.aktivniPismeno
    if ((pismeno === ' ' && tlacitko === '______') || (tlacitko.length === 1 && tlacitko.toLowerCase() === pismeno.toLowerCase()) || (tlacitko.length === 2 && tlacitko.toLowerCase().includes(pismeno.toLowerCase()) && tlacitko !== '  ')) {
        return true
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
    if (['"', '/', '?', ':', '_', '!', '(', '%', 'ˇ', '°'].includes(pismeno)) {
        return true
    } else if (/^\d$/.test(pismeno)) { // jestli to je cislo
        return true
    }

    return (pismeno === pismeno.toUpperCase() && !["Ů", "Ú", "+", "=", "-"].includes(pismeno))  // Ů a Ú nejdou psat se shiftem

}
</script>

<template>
    <div id="klavesnice" class="pruhledne">
        <div class="radek" v-for="radek in schema">
            <div v-for="tlacitko in radek" class="klavesa"
                :class="{ oznacenaKlavesa: oznacene(tlacitko) || (tlacitko === 'Shift' && shiftSviti), fjPodtrzeni: tlacitko === 'F' || tlacitko === 'J' }"
                :style="{ backgroundColor: barva(tlacitko), flexGrow: delkaTlacitka(tlacitko) }">

                <div v-if="tlacitko !== '^v'" :style="{ color: '#000' }">
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
    transition: background-color 0.2s;
}

.oznacenaKlavesa {
    border: #fff solid 3px;
    filter: brightness(150%);
}

.oznacenaKlavesa div {
    font-weight: 800 !important;
    transition: 0.1s;
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
    margin-top: 40px;
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
}

.fjPodtrzeni div {
    border-bottom: 1px solid black !important;
    width: 10px;
    padding-top: 2px;
}

</style>