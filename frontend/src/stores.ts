import { ref } from "vue";

export let prihlasen = ref(false)
export const role = ref("basic")
export const tokenJmeno = "pavouk_token"
export const cislaProcvicJmeno = "pavouk_procvic_"
export const nastaveniJmeno = "pavouk_nastaveni_psani"
export const levelyRychlosti = [50, 90, 130]
export const levelyPresnosti = [92.5, 97.5] // jen pro message uzivateli, ne pro hvezdy
export const maxPismenNaRadek = 629 / 17 // sirka ramecku / sirka pismene

export const moznostiRocnik = ["1.", "2.", "3.", "4.", "5.", "6.", "7.", "8.", "9.", "Prima ", "Sekunda ", "Tercie ", "Kvarta ", "Kvinta ", "Sexta ", "Septima ", "Oktáva "]
export const moznostiTrida = ["A", "B", "C", "D", "E", "F", "G"]

export const mobil = ref(document.body.clientWidth <= 900)
