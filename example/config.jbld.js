const babel = require("@babel/core");
const importVisitor = require("babel-plugin-import-visitor");

class BabelPlugin {
  constructor(opts) {
    this.opts = opts;
  }

  run(code) {
    const imports = [];
    const opts = {
      ...this.opts,
      plugins: (this.opts.plugins ? this.opts.plugins : []).concat([
        importVisitor(node => {
          imports.push({
            name: node.value,
            resolved: node.value + ".js",
            kind: 'static',
          });
        }),
      ]),
      ast: false,
    }
    const result = babel.transformSync(code, opts);
    return {
      imports,
      output: result.code,
    };
  }
}

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
