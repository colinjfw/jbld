const { BabelPlugin } = require("./plugin");

module.exports = {
  plugins: {
    babel: new BabelPlugin({
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
    }),
  },
};
