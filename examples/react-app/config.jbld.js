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
    logConfig: true,
    serve: ":3000",
    mode: "production",
    bundler: {
      baseUrl: "/",
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
      workers: 5,
    },
  },
  rules: [
    {
      test: /\.js$/,
      use: [babel],
    }
  ]
});
