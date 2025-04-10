<script setup lang="ts">
import { computed, ref, watch, onMounted } from "vue";
import { useRoute } from "vue-router";
import { os } from "../stores.ts";
import { delkaKlavesMacOS, delkaKlavesWindows, prstokladMacOS, prstokladWindows, schemaMacOS, schemaWindows, specialniZnakyMacOS, specialniZnakyWindowsQWERTZ, specialniZnakyWindowsQWERTY, schemaLinux, prstokladLinux, barvy, specialniZnakyLinuxQWERTZ, specialniZnakyLinuxQWERTY } from "../layouts.ts";

const props = defineProps({
    aktivniPismeno: {
        type: String,
        default: ""
    },
    typ: {
        type: String,
        reqired: true
    },
    rozmazat: Boolean,
    cekame: Boolean,
    fullHide: Boolean
})

const cesta = useRoute().path.split("/")

const znakySeShiftem: Map<string, boolean> = new Map()
const schema = ref([] as string[][])
let delkaKlaves: { [id: string]: number } = {}
let prstoklad: Map<string, string> = new Map()
let specialniZnaky: Map<string, string[]> = new Map()

if (os.value === "windows") {
    schema.value = schemaWindows
    delkaKlaves = delkaKlavesWindows
} else if (os.value === "macos") {
    schema.value = schemaMacOS
    delkaKlaves = delkaKlavesMacOS
} else if (os.value === "linux") {
    schema.value = schemaLinux
    delkaKlaves = delkaKlavesWindows
} else {
    schema.value = schemaWindows
    delkaKlaves = delkaKlavesWindows

    console.log("CO TO MÁŠ ZA SYSTÉM BRO")
}

prstoklad = calcPrstoklad()

const oznacenyPrst = computed(() => {
    prstoklad.forEach((prst, tlacitko) => {
        if (oznacene(tlacitko)) return prst
    })
    return ""
})

onMounted(() => {
    prohoditKlavesnici(props.typ!)
    calculateZnakySeShiftem()
})

watch(() => props.typ, (ted) => {
    prohoditKlavesnici(ted!)
    calculateZnakySeShiftem()
})

function prohoditKlavesnici(rozlozeni: string) {
    if (rozlozeni === "qwerty") {
        schema.value[1][6] = "Y"
        if (os.value === "windows") {
            schema.value[3][1] = "Z"
            specialniZnaky = specialniZnakyWindowsQWERTY
        } else if (os.value === "macos") {
            schema.value[3][2] = "Z"
            specialniZnaky = specialniZnakyMacOS
        } else if (os.value === "linux") {
            schema.value = schemaLinux
            specialniZnaky = specialniZnakyLinuxQWERTY
            schema.value[3][1] = "Z"
        }
        prstoklad.set("Y", "P_Ukaz")
        prstoklad.set("Z", "L_Mali")
    } else {
        schema.value[1][6] = "Z"
        if (os.value === "windows") {
            schema.value[3][1] = "Y"
            specialniZnaky = specialniZnakyWindowsQWERTZ
        } else if (os.value === "macos") {
            schema.value[3][2] = "Y"
            specialniZnaky = specialniZnakyMacOS
        } else if (os.value === "linux") {
            schema.value = schemaLinux
            specialniZnaky = specialniZnakyLinuxQWERTZ
            schema.value[3][1] = "Y"
        }
        prstoklad.set("Z", "P_Ukaz")
        prstoklad.set("Y", "L_Mali")
    }
}

function tlacPismeno(cislo: number, tlacitko: string) {
    try {
        if (tlacitko.length === 2) return tlacitko[cislo]
        else if (tlacitko.length === 1 && cislo === 0) return tlacitko[0]
        else if ((tlacitko === "Shift" || tlacitko === "Shift1" || tlacitko === "Shift2") && cislo === 0) return "Shift"
        else if (tlacitko.length >= 2 && cislo === 0) return tlacitko
    } catch { // kvůli starším prohlížečům (koukám na tebe safari <14.0)
        if (tlacitko.length === 2) return tlacitko.charAt(cislo)
        else if (tlacitko.length === 1 && cislo === 0) return tlacitko.charAt(0)
        else if ((tlacitko === "Shift" || tlacitko === "Shift1" || tlacitko === "Shift2") && cislo === 0) return "Shift"
        else if (tlacitko.length >= 2 && cislo === 0) return tlacitko
    }
    return ""
}

function oznacene(tlacitko: string) {
    if (props.aktivniPismeno == "") return false
    if (tlacitko === "Shift" || tlacitko === "Shift1" || tlacitko === "Shift2") {
        return potrebujeShift(props.aktivniPismeno) // shift když velký písmeno
    }
    const pismeno = props.aktivniPismeno.toLowerCase()
    const normalized = pismeno.normalize("NFD")
    const velky = props.aktivniPismeno.toUpperCase() === props.aktivniPismeno

    // easy
    if (normalized.length === 1 && pismeno === tlacitko.toLowerCase()) return true // klasický písmeno bez diakritiky
    if (tlacitko.length === 2 && (props.aktivniPismeno === tlacitko.charAt(1) || props.aktivniPismeno === tlacitko.charAt(0))) return true // znaky v dvouznakových klávesách
    if (pismeno === " " && tlacitko === "______") return true // spacebar

    // horší
    if ("óťňď".includes(pismeno)) { // znaky s diakritikou
        if (tlacitko === "ˇ´") return true
        if (normalized.replace(/\p{Diacritic}/gu, "") === tlacitko.toLowerCase()) return true
        return false
    }
    if (velky && normalized.length === 2) { // velká písmena s diakritikou
        if (pismeno !== "ů") { // Ů je ve speciálních znacích
            if (tlacitko.toLowerCase() === normalized.replace(/\p{Diacritic}/gu, "")) return true
            if (tlacitko === "ˇ´") return true
        }
    }
    if (specialniZnaky.has(props.aktivniPismeno)) { // speciální znaky, které potřebují speciální kombinaci
        const potrebnaTlacitka = specialniZnaky.get(props.aktivniPismeno)!
        for (let i = 0; i < potrebnaTlacitka.length; i++) {
            if (tlacitko == potrebnaTlacitka[i]) {
                return true
            }
        }
    }

    return false
}

function barva(tlacitko: string) {
    if (tlacitko === "Shift" || tlacitko === "Shift1" || tlacitko === "Shift2") return barvy.get("P_Mali")
    if (!specialniZnaky.has(props.aktivniPismeno) && tlacitko === "⌥") return ""
    if (specialniZnaky.has(props.aktivniPismeno) && (tlacitko === "Ctrl" || tlacitko === "Alt" || tlacitko === "AltGr")) {
        const tlacitka = specialniZnaky.get(props.aktivniPismeno)!
        if (!tlacitka?.includes(tlacitko)) return ""
    }
    if (!specialniZnaky.has(props.aktivniPismeno) && (tlacitko === "Ctrl" || tlacitko === "Alt" || tlacitko === "AltGr")) return ""

    const prst: string = prstoklad.get(tlacitko) || ""
    return barvy.get(prst)
}

function delkaTlacitka(tlacitko: string) {
    if (delkaKlaves[tlacitko] === undefined) {
        return 0
    }
    return delkaKlaves[tlacitko]
}

function calculateZnakySeShiftem() {
    znakySeShiftem.clear()
    for (let r = 0; r < schema.value.length; r++) {
        for (let s = 0; s < schema.value[r].length; s++) {
            const klavesa = schema.value[r][s]
            const prvniZnakVKlavese = klavesa.charAt(0)
            if (/^\d$/.test(prvniZnakVKlavese)) continue
            if (prvniZnakVKlavese === " " || prvniZnakVKlavese === "∧") continue
            if (klavesa.length == 2) znakySeShiftem.set(klavesa.charAt(0), true)
        }
    }

    specialniZnaky.forEach((klavesy, znak) => {
        if (klavesy.includes("Shift")) znakySeShiftem.set(znak, true)
    })
}

function calcPrstoklad(): Map<string, string> {
    const result = new Map<string, string>()
    let x: { [id: string]: string[] }

    if (os.value === "windows") x = prstokladWindows
    else if (os.value === "macos") x = prstokladMacOS
    else if (os.value === "linux") x = prstokladLinux
    else x = prstokladWindows

    for (let prst in x) {
        for (let i = 0; i < x[prst].length; i++) {
            result.set(x[prst][i], prst)
        }
    }
    return result
}

function potrebujeShift(pismeno: string) {
    if (/^\d$/.test(pismeno)) return true // cisla
    const pismenoBezDiakritiky = pismeno.normalize("NFD").replace(/\p{Diacritic}/gu, "")
    if (/^[A-Z]$/.test(pismenoBezDiakritiky)) return true // velká písmena (i ty co mají diakritiku)
    if (znakySeShiftem.has(pismeno)) return true
    if ("óťňď".includes(pismeno)) return true
}
</script>

<template>
    <div :class="{ hide: fullHide }" id="wraper">
        <div id="klavesnice" :class="{ rozmazany: props.rozmazat, 'rozmazany-min': props.cekame && !props.rozmazat }">
            <div :class="{}" class="radek" v-for="radek, i in schema" :key="i">
                <div v-for="tlacitko in radek" :key="tlacitko" class="klavesa"
                    :class="{ oznacenaKlavesa: oznacene(tlacitko), fjPodtrzeni: tlacitko === 'F' || tlacitko === 'J', enterNoha: tlacitko === 'enter-noha' }"
                    :style="{ backgroundColor: barva(tlacitko), flexGrow: delkaTlacitka(tlacitko), flexShrink: tlacitko === 'enter-noha' ? '30' : '0' }">

                    <div v-if="tlacitko === '∧∨'" id="sipky">
                        <div class="klavesa" style="height: 18px"></div>
                        <div class="klavesa" style="height: 18px"></div>
                    </div>
                    <div
                        v-else-if="tlacitko == 'enter-noha' || barva(tlacitko) === ''" />
                    <div v-else :style="{ fontSize: tlacitko === 'AltGr' ? '10px' : ''}">
                        {{ tlacPismeno(0, tlacitko) }} <br>
                        {{ tlacPismeno(1, tlacitko) }}
                    </div>
                </div>
            </div>
        </div>

        <div v-if="cesta[1] != 'prvni-psani'" id="zakladni-pozice"
            :style="{ opacity: props.cekame && !props.rozmazat ? 1 : 0, left: os === 'macos' ? '98px' : '92px' }">
            <div v-for="tlacitko in ['A', 'S', 'D', 'F', 'J', 'K', 'L', '&quot;ů']" :key="tlacitko" class="klavesa"
                :class="{ oznacenaKlavesa: oznacene(tlacitko), fjPodtrzeni: tlacitko === 'F' || tlacitko === 'J' }"
                :style="{ backgroundColor: barva(tlacitko), flexGrow: delkaTlacitka(tlacitko), marginRight: tlacitko == 'F' ? '90px' : 0 }">

                <div>
                    {{ tlacPismeno(0, tlacitko) }} <br>
                    {{ tlacPismeno(1, tlacitko) }}
                </div>
            </div>
        </div>

        <div id="nohy" :class="{ rozmazany: props.rozmazat }">
            <div class="kolecko" style="left: 101.5px; top: 7px;" :style="{ backgroundColor: barvy.get('L_Mali') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'L_Mali' }" />
            <div class="kolecko" style="left: 140.5px; top: -2px;" :style="{ backgroundColor: barvy.get('L_Prs') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'L_Prs' }" />
            <div class="kolecko" style="left: 182.5px; top: -9px;" :style="{ backgroundColor: barvy.get('L_Pros') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'L_Pros' }" />
            <div class="kolecko" style="left: 227px; top: -12px;" :style="{ backgroundColor: barvy.get('L_Ukaz') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'L_Ukaz' }" />
            <div class="kolecko" style="right: 149px; top: 7px;" :style="{ backgroundColor: barvy.get('P_Mali') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'P_Mali' }" />
            <div class="kolecko" style="right: 188px; top: -2px;" :style="{ backgroundColor: barvy.get('P_Prs') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'P_Prs' }" />
            <div class="kolecko" style="right: 229.5px; top: -9px;" :style="{ backgroundColor: barvy.get('P_Pros') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'P_Pros' }" />
            <div class="kolecko" style="right: 272px; top: -12px;" :style="{ backgroundColor: barvy.get('P_Ukaz') }"
                :class="{ koleckoAktivni: oznacenyPrst == 'P_Ukaz' }" />
            <img src="../assets/nohy.svg" alt="Nohy pavouka">
        </div>
    </div>
</template>

<style scoped>
#wraper {
    transition: opacity 0.15s;
}

.hide {
    opacity: 0;
}

#zakladni-pozice {
    display: flex;
    gap: 5px;
    font-size: 13px;
    line-height: 21px;
    height: 0;
    position: relative;
    bottom: 138px;
    transition: opacity 0.2s;
    user-select: none;
}

#nohy {
    margin-top: 25px;
    margin-bottom: -200px;
    position: relative;
    display: flex;
    justify-content: center;
    width: 100%;
    user-select: none;
    transition: filter 0.2s;
}

#nohy img {
    width: 61%;
    margin-right: 4%;
    margin-left: -2%;
}

.kolecko {
    position: absolute;
    width: 36px;
    height: 36px;
    border-radius: 100%;
    z-index: -1;
}

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
    color: black;
}

.oznacenaKlavesa,
.koleckoAktivni {
    border: #fff solid 3.5px;
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
    font-size: 13px;
    line-height: 21px;
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

.enterNoha {
    position: relative;
    top: -10px;
    height: 50px;
    margin-bottom: -10px;
    width: 27px;
}

.enterNoha::before {
    content: "";
    background-color: var(--tmave-fialova);
    z-index: -10;
    width: 5px;
    height: 5px;
    position: absolute;
    top: 6px;
    left: -5px;
    border-radius: 0 5px 0 0;
}

.enterNoha::after {
    content: "";
    background-color: inherit;
    z-index: -11;
    width: 10px;
    height: 10px;
    position: absolute;
    top: 1px;
    left: -5px;
}

.rozmazany {
    filter: blur(2px) brightness(20%) contrast(110%);
    /* blur je trochu heavy */
}

.rozmazany-min {
    filter: brightness(30%) contrast(105%);
    /* blur je trochu heavy */
}
</style>