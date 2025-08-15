import js from "@eslint/js";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import globals from "globals";
import reactX from "eslint-plugin-react-x";
import reactDom from "eslint-plugin-react-dom";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";
import { importX } from "eslint-plugin-import-x";
import tseslint from "typescript-eslint";
import { globalIgnores } from "eslint/config";
import tsParser from "@typescript-eslint/parser";

/**
 * @type {import("eslint").Config}
 */
export default tseslint.config([
  globalIgnores(["dist", ".yarn"]),
  {
    /* All React and TS / TSX file specific rules */
    files: ["**/*.{ts,tsx}"],
    extends: [
      // basic JS / TS related best practices
      js.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked,
      // import related best practices
      importX.flatConfigs.recommended,
      importX.flatConfigs.typescript,
      // react related lint rules
      reactX.configs["recommended-typescript"],
      reactDom.configs.recommended,
      reactHooks.configs["recommended-latest"],
      reactRefresh.configs.vite,
      // formatting comes last, should be applied
      // after all other lint autofixes have beenapplied
      eslintPluginPrettierRecommended,
    ],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: [
          "./tsconfig.node.json",
          "./tsconfig.app.json",
          "./tsconfig.playwright.json",
        ],
        tsconfigRootDir: import.meta.dirname,
      },
      ecmaVersion: "latest",
      globals: globals.browser,
    },
    rules: {
      // prefer types over interfaces
      // interfaces can be overridden with no warning from compiler
      // there are places where overriding is desired behavior
      // such as adding a field to `window` or Node's process.env
      // but these are exceptions, not the norm
      // For the most part, we want to use types over interfaces
      "@typescript-eslint/consistent-type-definitions": ["error", "type"],
      // we want some predictable ordering and grouping of imports
      "import-x/order": ["error"],
    },
  },
  {
    /* JS specific rules, apply to config files that are still in JS */
    files: ["**/*.{js,cjs}"],
    extends: [
      tseslint.configs.disableTypeChecked,
      importX.flatConfigs.recommended,
      eslintPluginPrettierRecommended,
    ],
    languageOptions: {
      globals: globals.node,
      ecmaVersion: "latest",
    },
    rules: {
      // we want some predictable ordering and grouping of imports
      "import-x/order": ["error"],
    },
  },
]);
