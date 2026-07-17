import { ref } from 'vue';
import { getOperatingSystem } from '@/platform';

export const uziv = ref({ jmeno: '', email: '' });
export const prihlasen = ref(false);
export const role = ref('basic');
export const tokenJmeno = 'pavouk_token';
export const nastaveniJmeno = 'pavouk_nastaveni_psani';
export const cislaProcvicJmeno = 'pavouk_procvic_';
export const levelyRychlosti = [30, 60, 90] as const;
export const levelyPresnosti = [92.5, 97.5] as const; // jen pro message uzivateli, ne pro hvezdy

export const moznostiRocnik = ['1.', '2.', '3.', '4.', '5.', '6.', '7.', '8.', '9.', 'Prima ', 'Sekunda ', 'Tercie ', 'Kvarta ', 'Kvinta ', 'Sexta ', 'Septima ', 'Oktáva '] as const;
export const moznostiTrida = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'] as const;
export const moznostiSkupina = ['-', '1', '2', '3', '4'] as const;

export const delkyCviceni = new Map<string, number>([
    ['nova', 3 * 60],
    ['naucena', 1 * 60],
    ['slova', 5 * 60],
    ['programator', 3 * 60],
]);

export function getCas(key: string) {
    return delkyCviceni.get(key) ?? 60;
}

export const mobil = ref(document.body.clientWidth <= 900);
export const os = ref(getOperatingSystem());

export const okZnaky = /([^A-Za-z0-9ěščřžýáíéůúťďňóĚŠČŘŽÝÁÍÉŮÚŤĎŇÓ ,.!?;:_=+\-*/%()[\]{}<>"])/;
