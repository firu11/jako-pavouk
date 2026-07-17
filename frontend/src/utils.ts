import { ref } from 'vue';
import { cislaProcvicJmeno, levelyPresnosti, levelyRychlosti, nastaveniJmeno, prihlasen, tokenJmeno } from '@/stores';
import api, { ApiError } from '@/api';

export function formatovanyPismena(pismena: string | string[] | undefined): string {
    if (pismena === undefined || pismena === '...') return pismena ?? '';
    return (typeof pismena === 'string' ? [...pismena] : pismena).join(', ');
}

const formatovaneKategorie: Readonly<Record<string, string>> = {
    'zbylá diakritika': 'Zbylá diakritika',
    'velká písmena (shift)': 'Velká písmena (Shift)',
    závorky: 'Závorky',
    operátory: 'Operátory',
    čísla: 'Číslovky',
    interpunkce: 'Interpunkce',
};

export function format(pismena: string): string {
    return formatovaneKategorie[pismena] ?? formatovanyPismena(pismena);
}

export function getToken() {
    return localStorage.getItem(tokenJmeno);
}

const formatovacDataPraha = new Intl.DateTimeFormat('cs-CZ', { timeZone: 'Europe/Prague' });
const formatovacDneMesicePraha = new Intl.DateTimeFormat('cs-CZ', { timeZone: 'Europe/Prague', day: 'numeric', month: 'numeric' });

export function formatDatumPraha(value: string | Date): string {
    const datum = typeof value == 'string' ? new Date(value) : value;
    return formatovacDataPraha.format(datum);
}

export function formatDenMesicPraha(value: string | Date): string {
    const datum = typeof value == 'string' ? new Date(value) : value;
    const parts = formatovacDneMesicePraha.formatToParts(datum);
    const den = parts.find((part) => part.type == 'day')?.value ?? '';
    const mesic = parts.find((part) => part.type == 'month')?.value ?? '';
    return `${den}.${mesic}.`;
}

export type NotificationType = 'warning' | 'copy' | 'info';

export interface AppNotification {
    id: number;
    text: string;
    type: NotificationType;
}

const notificationTypes = {
    vykricnik: 'warning',
    copy: 'copy',
    'svisla-cara': 'info',
} as const;

let nextNotificationId = 0;
export const oznameni = ref<AppNotification[]>([]);

export function pridatOznameni(text = 'Něco se pokazilo', cas = 4000, typ: keyof typeof notificationTypes = 'vykricnik') {
    const notification: AppNotification = {
        id: nextNotificationId++,
        text,
        type: notificationTypes[typ],
    };
    oznameni.value.push(notification);

    window.setTimeout(() => {
        const index = oznameni.value.findIndex(({ id }) => id === notification.id);
        if (index !== -1) oznameni.value.splice(index, 1);
    }, cas);
}

export function napovedaKNavigaci() {
    pridatOznameni('Pro nápovědu k navigaci se podívej do záložky Jak psát.');
}

export function checkTeapot(e: unknown): boolean {
    if (e instanceof ApiError && e.response.status == 418) {
        if (oznameni.value.length < 3) {
            pridatOznameni('Dej si čajík a vydýchej se...');
        }
        return true;
    }
    return false;
}

export class Oznacene {
    index = ref(0);
    max: number = 4;
    bezOznaceni: boolean = false;
    mensi() {
        if (this.index.value > 1) {
            this.index.value--;
        }
    }
    vetsi() {
        if (this.index.value < this.max) {
            this.index.value++;
        }
    }
    setMax(max: number) {
        this.max = max;
    }
    is(n: number) {
        if (n < 6 && n == this.index.value) return true;
        else if (n >= 6 && 14 > n && this.index.value + 1 == n) return true;
        else if (n >= 14 && this.index.value + 2 == n) return true;
        return false;
    }
}

export class MojeMapa extends Map<string, number> {
    put(znak: string) {
        const normalized = znak.toLocaleLowerCase();
        this.set(normalized, (this.get(normalized) ?? 0) + 1);
    }

    top(n: number) {
        return [...this.entries()]
            .sort(([, firstCount], [, secondCount]) => secondCount - firstCount)
            .slice(0, n)
            .map(([znak, pocet]) => ({ znak, pocet }));
    }
}

export function getCisloPochvaly(rychlost: number, presnost: number): number {
    if (rychlost >= levelyRychlosti[2] && presnost >= levelyPresnosti[1]) {
        // paradni
        return 0;
    } else if (rychlost >= levelyRychlosti[1] && rychlost < levelyRychlosti[2] && presnost >= levelyPresnosti[1]) {
        // rychlost muze byt lepsi
        return 1;
    } else if (presnost >= levelyPresnosti[0] && presnost < levelyPresnosti[1] && rychlost >= levelyRychlosti[2]) {
        // presnost muze byt lepsi
        return 2;
    } else if (presnost >= levelyPresnosti[0] && presnost < levelyPresnosti[1] && rychlost >= levelyRychlosti[1] && rychlost < levelyRychlosti[2]) {
        // oboje muze byt lepsi
        return 3;
    } else if (rychlost < levelyRychlosti[1] && presnost < levelyPresnosti[0]) {
        // oboje bad
        return 6;
    } else if (rychlost < levelyRychlosti[1]) {
        // rychlost bad
        return 4;
    } else if (presnost < levelyPresnosti[0]) {
        // presnost bad
        return 5;
    }
    return 0; // nestane se
}

export function clone<typ>(obj: typ): typ {
    // kvůli starším prohlížečům (koukám na tebe safari <14.0)
    let x: typ;
    try {
        x = structuredClone(obj);
    } catch {
        x = JSON.parse(JSON.stringify(obj));
    }
    return x;
}

export function saveNastaveni(diakritika: boolean, velkaPismena: boolean, vetySlova: boolean, delka: number, klavesnice: boolean) {
    localStorage.setItem(nastaveniJmeno, JSON.stringify({ diakritika, velkaPismena, vetySlova, delka, klavesnice }));
}

export function naJednoDesetiny(cpm: number): number {
    return Math.round(cpm * 10) / 10;
}

export async function getCisloProcvic(id: string): Promise<number[]> {
    if (prihlasen.value) {
        const cisla = await getCisloProcvicFromServer(id);
        if (cisla.length != 0) {
            console.log('Progress načten ze serveru!');
            return cisla;
        }
        console.log('Na serveru byl progress prázdný.');
    }

    const v = localStorage.getItem(cislaProcvicJmeno + id);
    if (v === null) {
        return [1, 0];
    }
    const cisla = v.split(',');
    if (cisla.length == 1) {
        return [Number(cisla[0]), 0];
    } else if (cisla.length < 1) {
        return [1, 0];
    }
    return [Number(cisla[0]), Number(cisla[1])];
}

export function setCisloProcvic(id: string, cisla: number[]) {
    localStorage.setItem(cislaProcvicJmeno + id, cisla.join(','));

    if (prihlasen.value) {
        saveCisloProcvicToServer(id, cisla);
    }

    console.log(`Progres uložen: cisloTextu=${id}, localstorage/server=[${cisla}]`);
}

export function saveCisloProcvicToServer(id: string, cisla: number[]) {
    api.post('/uloz-procvic-postup', { cislo_textu: Number(id), cislo_kapitoly: cisla[0], cislo_slova: cisla[1] }, { headers: { Authorization: `Bearer ${getToken()}` } }).catch((e) => {
        console.log(e);
    });
}

export async function getCisloProcvicFromServer(id: string): Promise<number[]> {
    return api
        .get('/procvic-postup/' + id, {
            headers: { Authorization: `Bearer ${getToken()}` },
        })
        .then((response) => {
            return response.data as number[];
        })
        .catch((error) => {
            console.error(error);
            return [];
        });
}

export function postKlavesnice(klavesnice: boolean) {
    const k = klavesnice ? 'qwerty' : 'qwertz';
    api.post('/ucet-zmena', { zmena: 'klavesnice', hodnota: k }, { headers: { Authorization: `Bearer ${getToken()}` } }).catch((e) => {
        console.log(e);
    });
}
