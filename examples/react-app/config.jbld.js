const { BabelPlugin } = require("../../lib/babel");
const { Configuration } = require("../../lib");

const babel = new BabelPlugin({
  presets: ["react-app"],
  plugins: ["transform-node-env-inline", "@babel/plugin-transform-modules-commonjs"]
});

module.exports = new Configuration({
  options: {
    logConfig: true,
    serve: ":3000",
    bundler: {
      baseUrl: "",
      outputDir: "./dist/bundle",
      assetPath: "static",
      public: {
        dir: "./public",
        html: ["index.html"],
      }
    },
    compiler: {
      sourceDir: ".",
      outputDir: "./dist/target",
      entrypoints: ["src/index.js"],
      workers: 50,
    },
  },
  rules: [
    {
      test: /\.js$/,
      use: [babel],
    }
  ]
});
