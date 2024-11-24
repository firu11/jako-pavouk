<script setup lang="ts">
import { computed, ref } from 'vue';

const props = defineProps({
    taby: {
        type: Array<Array<string>>,
        default: function () {
            return [["sus", "Nějaký tab"]]
        }
    },
    defaultTab: String
})
const emit = defineEmits(["zmena"])

const tab = ref(props.defaultTab)
const index = computed(() => {
    for (let i = 0; i < props.taby.length; i++) {
        if (props.taby[i][0] == tab.value) return i
    }
    return 0
})

defineExpose({ tab })

</script>
<template>
    <div id="prepinac-tabu">
        <label v-for="x, i in taby" :key="i" :class="{ oznaceny: tab == x[0] }">
            {{ x[1] }}
            <input type="radio" :value="x[0]" v-model="tab" @change="emit('zmena')">
        </label>

        <span :style="{ transform: `translateX(${100 * index}px)` }"></span>
    </div>
</template>
<style scoped>
#prepinac-tabu {
    display: flex;
    padding: 8px;
    border-radius: 100px;
    background-color: var(--tmave-fialova);
    margin: 20px 0;
    height: 46px;
}

#prepinac-tabu input {
    display: none;
}

#prepinac-tabu label {
    padding: 5px;
    width: 100px;
    z-index: 1;
    cursor: pointer;
    color: #c5c5c5;
    transition: 0.15s;
    font-weight: 400;
}

.oznaceny {
    font-weight: 500 !important;
    color: var(--bila) !important;
}

#prepinac-tabu span {
    position: absolute;
    background-color: var(--fialova);
    width: 100px;
    height: 30px;
    border-radius: 100px;
    transition: 0.15s ease-out;
}
</style>