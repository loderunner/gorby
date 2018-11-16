module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: ['plugin:vue/essential'],
  rules: {
    // allow async-await
    'generator-star-spacing': 'off',
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'semi': ['error', 'never'],
    'no-extra-semi': 'error',
    'quotes': ['error', 'single'],
    'sort-imports': 'error',
    'no-var': 'error',
    'prefer-const': 'error',
    'sort-imports': 'off',
    "space-before-function-paren": ["error", "never"],
  },
  parserOptions: {
    parser: 'babel-eslint'
  }
};
