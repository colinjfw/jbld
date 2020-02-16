const { BabelPlugin } = require("./plugin");

module.exports = {
  plugins: {
    babel: new BabelPlugin({
      presets: [
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
