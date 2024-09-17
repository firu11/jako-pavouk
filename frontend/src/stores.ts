import { ref } from "vue";

export let prihlasen = ref(false)
export const role = ref("basic")
export const tokenJmeno = "pavouk_token"
export const nastaveniJmeno = "pavouk_nastaveni_psani"
export const cislaProcvicJmeno = "pavouk_procvic_"
export const levelyRychlosti = [60, 100, 140]
export const levelyPresnosti = [92.5, 97.5] // jen pro message uzivateli, ne pro hvezdy
export const maxPismenNaRadek = 629 / 19 // sirka ramecku / sirka pismene

export const moznostiRocnik = ["1.", "2.", "3.", "4.", "5.", "6.", "7.", "8.", "9.", "Prima ", "Sekunda ", "Tercie ", "Kvarta ", "Kvinta ", "Sexta ", "Septima ", "Oktáva "]
export const moznostiTrida = ["A", "B", "C", "D", "E", "F", "G", "H"]
export const moznostiSkupina = ["-", "1", "2", "3", "4"]

export const delkyCviceni = new Map<string, number>([
    ["nova", 2*60],
    ["naucena", 1*60],
    ["slova", 3*60],
    ["programator", 2*60],
])

export function getCas(key: string) {
    return delkyCviceni.get(key) || 60 // default
}

export const mobil = ref(document.body.clientWidth <= 900)

export const okZnaky = /([^A-Za-z0-9ěščřžýáíéůúťďňóĚŠČŘŽÝÁÍÉŮÚŤĎŇÓ ,.!?;:_=+\-*/%()[\]{}<>"])/
