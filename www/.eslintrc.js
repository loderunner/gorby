module.exports = {
  root: true,
  env: { node: true },
  extends: ['plugin:vue/essential'],
  rules: {
    // allow async-await
    'generator-star-spacing': 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'semi': ['error', 'never'],
    'no-extra-semi': 'error',
    'quotes': ['error', 'single'],
    'brace-style': 'error',
    'no-var': 'error',
    'prefer-const': 'error',
    'sort-imports': 'off',
    'space-before-function-paren': ['error', 'never'],
    'object-curly-spacing': ['error', 'always'],
    'object-property-newline': ['error', { 'allowAllPropertiesOnSameLine': true }],
  },
  parserOptions: { parser: 'babel-eslint' }
};
