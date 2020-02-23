const { BabelPlugin } = require("jbld/babel");
const { Configuration } = require("jbld");

const babel = new BabelPlugin({
  comments: false,
  presets: [
    "react-app",
    "minify",
  ],
  plugins: [
    "transform-inline-environment-variables",
    "@babel/plugin-transform-modules-commonjs",
  ]
});

module.exports = new Configuration({
  options: {
    baseUrl: "/",
    entrypoints: ["src/index.js"],
    workers: 5,
  },
  rules: [
    {
      test: /\.js$/,
      use: [babel],
    }
  ]
});
