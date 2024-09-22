import globals from "globals";
import pluginJs from "@eslint/js";
import tseslint from "typescript-eslint";
import pluginVue from "eslint-plugin-vue";
import eslintPluginVueScopedCSS from "eslint-plugin-vue-scoped-css";

export default [
    { files: ["**/*.{js,mjs,cjs,ts,vue}"] },
    { files: ["**/*.js"], languageOptions: { sourceType: "script" } },
    { languageOptions: { globals: globals.browser } },
    pluginJs.configs.recommended,
    ...tseslint.configs.recommended,
    ...pluginVue.configs["flat/essential"],
    ...eslintPluginVueScopedCSS.configs['flat/recommended'],
    {
        rules: {
            "vue/multi-word-component-names": "off",
        },
    },
    { files: ["**/*.vue"], languageOptions: { parserOptions: { parser: tseslint.parser } } },
];