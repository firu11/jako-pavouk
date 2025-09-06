<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue';

const text = ref('');
const vyska = 'calc(100vh - 60px - 40px - 25px - 30px - 40px - 11px)';

const textarea = ref<HTMLTextAreaElement | null>(null);
const div = ref<HTMLElement>();
const highlightText = ref('');

const isComposing = ref(false);

const lineBreak = /(\r\n|\r|\n)/g;
const red = /([^A-Za-z0-9ěščřžýáíéůúťďňóĚŠČŘŽÝÁÍÉŮÚŤĎŇÓ ,.!?;:_=+\-*/%#$~§()[\]{}<>"'@&|\r\n\\]+)/g;
const orange = /(([_=+\-*%()[\]{}"'@&|#$~§\\]+)|(?<!\/?mark)([<>])(?!\/?mark))|(\/)(?!mark>)/g;
const mezeraNaZacatku = /(^ )/g;
const viceMezer = /( {2,})/g;

function scrollDiv() {
    div.value?.scrollTo(0, textarea.value?.scrollTop ?? 0);
}

function updateHighlighting(fromInput = false) {
    const val = textarea.value?.value ?? '';
    if (!fromInput) text.value = val.replace(lineBreak, ''); // clean breaks if not from input

    let t = val.replace(red, `<mark>$&</mark>`);
    t = t.replace(orange, '<mark2>$&</mark2>');
    t = t.replace(lineBreak, '<mark>↵\n</mark>');
    t = t.replace(viceMezer, '<mark>$&</mark>');
    t = t.replace(mezeraNaZacatku, '<mark>$&</mark>');

    highlightText.value = t;
}

function onInput(e: Event) {
    const val = (e.target as HTMLTextAreaElement).value;
    text.value = val;
    updateHighlighting(true); // Update highlights on input
}

function onCompositionStart() {
    isComposing.value = true;
}

function onCompositionUpdate() {
    // If composing, update highlights as you type composed characters
    updateHighlighting(true);
}

function onCompositionEnd() {
    isComposing.value = false;
    updateHighlighting(true); // Re-apply final highlight after composition ends
}

onMounted(() => {
    updateHighlighting(); // Initial highlight setup
});

watch(text, () => {
    nextTick(() => {
        updateHighlighting(true);
    });
});

const ready = computed(() => !text.value.match(red));

defineExpose({ text, ready });
// krejzy https://codersblock.com/blog/highlight-text-inside-a-textarea/
</script>

<template>
    <div ref="div" v-html="highlightText"></div>
    <textarea
        name="text"
        ref="textarea"
        placeholder="Text, který budou žáci psát..."
        v-model="text"
        @scroll="scrollDiv"
        @compositionstart="onCompositionStart"
        @input="onInput"
        @compositionupdate="onCompositionUpdate"
        @compositionend="onCompositionEnd"
    />
</template>
<style scoped>
div {
    position: absolute;
    padding: 10px;
    width: 410px;
    height: v-bind(vyska);
    /* celá obrazovka - všechno co je nad tím */
    text-align: start;
    top: 55px;
    word-wrap: break-word;
    z-index: 1;
    background-color: var(--tmave-fialova);
    border-radius: 5px;
    font-size: 16px !important;
    overflow-y: auto;
    text-decoration: none;

    color: rgb(156, 156, 156);
    font-family: 'Red Hat Mono', monospace;

    white-space: pre-wrap;
    word-wrap: break-word;
    line-height: 21px;
}

textarea {
    position: relative;
    background-color: transparent;
    border: 0;
    height: v-bind(vyska);
    width: 100%;
    font-size: 16px !important;
    text-decoration: none;
    color: var(--bila);
    padding: 10px;
    resize: none;
    z-index: 2;
    line-height: 21px;
    font-family: 'Red Hat Mono', monospace;
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
    border-radius: 5px;
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
    cursor: ns-resize;
    /* Darker color on hover */
}

div::-webkit-scrollbar {
    width: 10px;
    background: transparent;
}
</style>

<!-- eslint-disable-next-line vue-scoped-css/enforce-style-type -->
<style>
div > mark,
div > mark2 {
    font-family: 'Red Hat Mono', monospace;
    border-radius: 3px;
    color: rgb(156, 156, 156) !important;
}

div > mark {
    background-color: rgba(255, 0, 0, 0.4);
}

div > mark2 {
    background-color: rgba(255, 136, 0, 0.4);
}
</style>
