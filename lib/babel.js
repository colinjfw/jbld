const importVisitor = require("babel-plugin-import-visitor");
const transformRuntime = require("@babel/plugin-transform-runtime");
const dynamicImport = require("@babel/plugin-syntax-dynamic-import");
const babel = require("@babel/core");
const { resolve } = require("./index");

/**
 * Resolve plugin implements a babel plugin to resolve paths based on the jbld
 * resolve function specification.
 *
 * @param {Array} imports
 * @param {Source} source
 */
function babelResolvePlugin(imports, source) {
  return importVisitor(node => {
    const resolved = resolve(source.srcDir, source.src, node.value);
    if (!resolved) {
      return;
    }
    imports.push({ name: node.value, resolved, kind: 'static' });
  });
}

/**
 * Initializes a babel plugin. Takes any option like 'presets' or 'plugins' that
 * can be passed to a babel transform.
 */
class BabelPlugin {
  constructor(opts) {
    this.opts = opts;
  }

  run(code, source) {
    let imports = [];
    let opts = {
      ...this.opts,
      plugins: (this.opts.plugins ? this.opts.plugins : []).concat([
        dynamicImport,
        transformRuntime,
        babelResolvePlugin(imports, source),
      ]),
      ast: false,
    }
    let result = babel.transformSync(code, opts);
    return {
      imports,
      output: result.code,
    };
  }
}

module.exports = { BabelPlugin };
