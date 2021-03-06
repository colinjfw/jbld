const { BabelPlugin } = require("../../lib/babel");
const { Configuration } = require("../../lib");

const babel = new BabelPlugin({
  plugins: [
    "@babel/plugin-transform-runtime",
    "@babel/plugin-syntax-dynamic-import",
  ],
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
  options: {
    public: null
  },
  rules: [
    {
      test: /\.js$/,
      use: [babel],
    }
  ]
});
