<script setup lang="ts">
import { ref } from 'vue';
import TextZadani from '../components/TextZadani.vue';

const sus = ref<InstanceType<typeof TextZadani> | null>(null)

function getText() {
    sus.value!.text = "ssdkfsjdnkfjnk"
}

</script>
<template>
    <div id="pulic">
        <div id="nastaveni">
            cau
        </div>

        <div id="text">
            <div>
                <div>{{ sus?.text.length }} znaků</div>
                <button class="tlacitko" id="genBtn" @click="getText()">Generovat</button>
            </div>

            <TextZadani ref="sus" />
        </div>
    </div>
</template>
<style scoped>
#genBtn {
    margin-top: 0;
    width: 120px;
    background-color: var(--tmave-fialova);
}

#genBtn:hover {
    background-color: var(--fialova);
}

#pulic {
    display: flex;
    justify-content: space-between;
    gap: 40px;
    width: 860px;
}

#text {
    width: 410px;
    display: flex;
    flex-direction: column;
    height: 400px;
    border-radius: 10px;
    gap: 15px;
}

#text>div {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

#text>div>div {
    width: 150px;
    padding: 10px;
    background-color: var(--tmave-fialova);
    border-radius: 5px;
}

#nastaveni {
    width: 410px;
    display: flex;
    gap: 10px;
    flex-direction: column;
    height: 400px;
    background-color: var(--tmave-fialova);
    padding: 10px 15px;
    border-radius: 10px;
}
</style>