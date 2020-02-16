const { BabelPlugin } = require("../lib/babel");
const { Configuration } = require("../lib");

const babel = new BabelPlugin({
  presets: [
    "minify",
    [
      "@babel/preset-env",
      {
        corejs: 2,
        useBuiltIns: "entry"
      }
    ]
  ]
});

module.exports = new Configuration({
  rules: [
    {
      test: /\.js$/,
      use: [babel],
    }
  ]
});
