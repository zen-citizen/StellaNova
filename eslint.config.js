import js from "@eslint/js";
import globals from "globals";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import reactX from "eslint-plugin-react-x";
import reactDom from "eslint-plugin-react-dom";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";
import importPlugin from "eslint-plugin-import";
import tseslint from "typescript-eslint";
import { globalIgnores } from "eslint/config";

/**
 * @type {import("eslint").Config}
 */
export default tseslint.config([
  globalIgnores(["dist", ".yarn"]),
  {
    files: ["**/*.{ts,tsx}"],
    extends: [
      js.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      ...tseslint.configs.stylisticTypeChecked,
      importPlugin.flatConfigs.recommended,
      importPlugin.flatConfigs.typescript,
      reactX.configs["recommended-typescript"],
      reactDom.configs.recommended,
      reactHooks.configs["recommended-latest"],
      reactRefresh.configs.vite,
      eslintPluginPrettierRecommended,
    ],
    languageOptions: {
      parserOptions: {
        project: ["./tsconfig.node.json", "./tsconfig.app.json"],
        tsconfigRootDir: import.meta.dirname,
      },
      ecmaVersion: 2020,
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
    },
  },
  {
    files: ["**/*.{js,cjs}"],
    extends: [
      tseslint.configs.disableTypeChecked,
      importPlugin.flatConfigs.recommended,
      eslintPluginPrettierRecommended,
    ],
    languageOptions: {
      globals: globals.node,
      ecmaVersion: "latest",
    },
  },
]);
