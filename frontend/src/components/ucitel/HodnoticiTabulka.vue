<script setup lang="ts">
import { ref } from 'vue';
import { pridatOznameni } from '../../utils';

const rychlosti = ref(new Array<number>(4));
const spatne = ref(new Array<boolean>(4));
const minMax = [1, 999];

function zmena() {
    let notifikace = false;
    for (let i = 0; i < rychlosti.value.length; i++) {
        const r = rychlosti.value[i];
        if (r == null || r > minMax[1] || r < minMax[0] || r <= rychlosti.value[i + 1]) {
            notifikace = true;
            spatne.value[i] = true;
        } else {
            spatne.value[i] = false;
        }
    }
    if (notifikace) {
        pridatOznameni('Všechny rychlosti musí být vyplněné, <1000 a musí odpovídat pořadí...');
    }
}
</script>
<template>
    <div id="hodnotici-tabulka">
        <table>
            <colgroup>
                <col style="width: 150px" />
                <col style="width: 150px" />
                <col />
            </colgroup>
            <thead>
                <tr>
                    <th scope="col">Známka</th>
                    <th scope="col">Rychlost</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(_, i) in rychlosti">
                    <td>{{ i + 1 }}</td>
                    <td>
                        &ge;
                        <input
                            type="number"
                            :min="minMax[0]"
                            :max="minMax[1]"
                            :name="String(i + 1)"
                            :placeholder="String(150 - i * 20)"
                            v-model="rychlosti[i]"
                            @change="zmena"
                            :class="{ spatne: spatne[i] }"
                        />
                        CPM
                    </td>
                </tr>
                <tr>
                    <td>5</td>
                    <td>&lt; {{ rychlosti[3] ? rychlosti[3] : '???' }} CPM</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<style scoped>
#hodnotici-tabulka {
    display: flex;
    justify-content: center;
}

#hodnotici-tabulka table {
    border-collapse: collapse;
    border-spacing: 10px 0;
    width: 300px;
}

#hodnotici-tabulka table td {
    height: 34px;
}

#hodnotici-tabulka table thead tr th {
    border-bottom: 1px solid rgb(210, 210, 210) !important;
    font-size: 1.1em;
    font-weight: 500;
}

#hodnotici-tabulka table tbody tr:first-of-type td {
    padding-top: 6px;
    height: 40px;
}

#hodnotici-tabulka table tbody tr td {
    padding-bottom: 4px;
}

#hodnotici-tabulka table tbody tr td:first-of-type {
    font-size: 1.7em;
    font-weight: bold;
}

input[type='number'] {
    height: 25px;
    width: 3.2em;
    border: none;
    border-radius: 5px;
    font-size: 16px;
    color: white;
    background-color: var(--fialova);
    transition: 0.2s;
    font-family: inherit;
    padding: 0 10px 0 0;
    box-sizing: border-box;
    text-align: right;
}

input[type='number'].spatne {
    border: 1px solid red;
}

input[type='number']:hover {
    transition: 0.2s;
}

input[type='number']:focus {
    outline: none;
}

input[type='number']::placeholder {
    color: rgb(160, 160, 160);
}

input[type='number']::-webkit-inner-spin-button,
input[type='number']::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}
</style>
