import { FlatCompat } from "@eslint/eslintrc";
import globals from "globals";

import pluginJs from "@eslint/js";
import typescriptParser from "@typescript-eslint/parser";
import pluginReact from "eslint-plugin-react";
import reactRefresh from "eslint-plugin-react-refresh";
import tseslint from "typescript-eslint";

const __dirname = import.meta.dirname;

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const reactRecommended = pluginReact.configs.flat.recommended;

/** @type {import('eslint').Linter.Config[]} */
export default [
  {
    ignores: [],
  },
  {
    files: ["src/**/*.{js,jsx,mjs,cjs,ts,tsx}"],

    ignores: ["dist/**/*"],
    ...reactRecommended,
    settings: {
      react: {
        version: "detect",
      },
    },
    languageOptions: {
      ...reactRecommended.languageOptions,
      ecmaVersion: "latest",
      sourceType: "module",
      parser: typescriptParser,
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
      globals: {
        ...globals.serviceworker,
        ...globals.browser,
      },
    },
  },
  pluginJs.configs.recommended,
  ...tseslint.configs.recommended,
  pluginReact.configs.flat.recommended,
  reactRefresh.configs.recommended,
  ...compat.extends("plugin:react-hooks/recommended"),
  // prevent `'React' must be in scope when using JSX` (`import React` not needed in React 17+)
  pluginReact.configs.flat?.["jsx-runtime"] ?? {},

  // custom rules
  // TODO: clean these up and examine if we should enforce some subset of the rules being turned to warning and off
  {
    files: ["src/**/*.{js,jsx,mjs,cjs,ts,tsx}", "eslint.config.js"],
    ignores: ["*.config.ts", ".next/**/*"],
    rules:{
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-unused-vars": "warn",
      "@typescript-eslint/no-empty-object-type": "warn",
      "no-useless-escape": "off",
      "react-refresh/only-export-components": "off",
      "react/prop-types": "off",
      "@typescript-eslint/no-unsafe-function-type": "warn",
      "react/display-name": "off",
      "@typescript-eslint/ban-ts-comment": "off",
      "no-case-declarations": "off",
      "no-useless-catch": "off",
      "no-empty-pattern": "off",
      "react-hooks/exhaustive-deps": "off"
    },
  },
];