module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: ["plugin:vue/essential"],
  rules: {
    // allow async-await
    "generator-star-spacing": "off",
    "no-console": process.env.NODE_ENV === "production" ? "error" : "off",
    "no-debugger": process.env.NODE_ENV === "production" ? "error" : "off",
    "semi": ["error", "never"],
    "no-extra-semi": "error"
  },
  parserOptions: {
    parser: "babel-eslint"
  }
};
