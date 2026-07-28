module.exports = {
  root: true,
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:vue/vue3-recommended",
    "plugin:vue/vue3-essential",
  ],
  parserOptions: {
    ecmaVersion: "latest",
    sourceType: "module",
  },
  plugins: ["vue"],
  rules: {
    "vue/multi-word-component-names": "off", // karena komponen seperti AppLayout
    "vue/no-v-html": "off",
    "no-console": process.env.NODE_ENV === "production" ? "warn" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "warn" : "off",
    indent: ["error", 2],
    quotes: ["error", "single"],
    semi: ["error", "never"],
  },
  overrides: [
    {
      files: ["**/*.vue"],
      rules: {
        "vue/require-default-prop": "off",
        "vue/require-prop-types": "off",
      },
    },
  ],
};
