const { BabelPlugin, ImportManglePlugin } = require("../../lib/babel");
const { Configuration } = require("../../lib");

const babel = new BabelPlugin({
  comments: false,
  presets: [
    "react-app",
    "minify",
  ],
  plugins: [
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
