import { defineConfig } from "eslint/config";
import tslint from "typescript-eslint";
import json from "@eslint/json";
import eslint from "@eslint/js";
import css from "@eslint/css";
import globals from "globals";

export default defineConfig([
  {
    files: ["**/*.{js,mjs,cjs,ts,mts,cts}"],
    plugins: { eslint },
    extends: ["js/recommended"],
    languageOptions: { globals: { ...globals.browser, ...globals.node } },
  },
  tslint.configs.recommended as any,
  { files: ["**/*.json"], plugins: { json }, language: "json/json", extends: ["json/recommended"] },
  { files: ["**/*.css"], plugins: { css }, language: "css/css", extends: ["css/recommended"] },
  {
    extends: "prettier",
  },
]);
