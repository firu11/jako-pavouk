<script setup lang="ts">
import { ref, watch } from 'vue';

const text = ref("")

const textarea = ref<HTMLElement | null>(null)
const div = ref<HTMLElement>()

function scrollDiv() {
    div.value?.scrollTo(0, textarea.value!.scrollTop)
}

watch(text, () => {
    let t = text.value
        .replace(/\n$/g, '\n\n')
        .replace(/([^A-Za-z0-9ěščřžýáíéůúťďňó ,.!?;:_=+\-*/%()[\]{}<>]|.(\r\n|\r|\n))/g, "<m style='background-color: rgba(255, 0, 0, 0.4); border-radius: 2px;'>$&</m>")
    // krejzy https://codersblock.com/blog/highlight-text-inside-a-textarea/
    div.value!.innerHTML = t
})

defineExpose({ text })

</script>
<template>
    <div ref="div"></div>
    <textarea ref="textarea" placeholder="Text který budou žáci psát..." v-model="text" @scroll="scrollDiv" />
</template>
<style scoped>
div {
    position: absolute;
    padding: 10px;
    width: 410px;
    height: 345px;
    text-align: start;
    top: 300px;
    word-wrap: break-word;
    z-index: 1;
    background-color: var(--tmave-fialova);
    border-radius: 8px;
    font-size: 1rem;
    overflow-y: auto;

    color: green;

    white-space: pre-wrap;
    word-wrap: break-word;
}

textarea {
    background-color: transparent;
    color: blue;
    border: 0;
    height: 100%;
    width: 100%;
    font-size: 1rem;
    color: var(--bila);
    padding: 10px;
    resize: none;
    z-index: 2;
}

::placeholder {
    color: rgb(175, 175, 175);
    opacity: 1;
}

textarea::-webkit-scrollbar {
    width: 10px;
}

/* Customizes the track of the scrollbar */
textarea::-webkit-scrollbar-track {
    background: var(--tmave-fialova);
    /* Color of the track */
    border-radius: 3px;
    padding: 1px;
    cursor: default;
}

/* Customizes the thumb of the scrollbar */
textarea::-webkit-scrollbar-thumb {
    background: var(--fialova);
    /* Color of the thumb */
    border-radius: 3px;
    cursor: default;
}

/* Changes the thumb color on hover */
textarea::-webkit-scrollbar-thumb:hover {
    background: var(--svetle-fialova);
    /* Darker color on hover */
}

div::-webkit-scrollbar {
    width: 10px;
    background: transparent;
}
</style>