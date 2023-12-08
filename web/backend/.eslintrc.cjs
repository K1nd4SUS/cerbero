module.exports = {
  "root": true,
  "env": {
    "node": true
  },
  "extends": [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended-type-checked"
  ],
  "ignorePatterns": [
    ".eslintrc.cjs"
  ],
  "parser": "@typescript-eslint/parser",
  "parserOptions": {
    "ecmaVersion": "latest",
    "sourceType": "module",
    "project": ["tsconfig.json"],
    "tsconfigRootDir": __dirname
  },
  "rules": {
    "indent": ["error", 2, { "SwitchCase": 1 }],
    "linebreak-style": ["error", "unix"],
    "quotes": ["error", "double"],
    "semi": ["error", "never"],
    "eqeqeq": "error",
    "no-trailing-spaces": "error",
    "object-curly-spacing": ["error", "always"],
    "arrow-spacing": ["error", { "before": true, "after": true }],
    "eol-last": 1
  }
}
