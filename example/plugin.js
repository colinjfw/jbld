const importVisitor = require("babel-plugin-import-visitor");
const resolveFrom = require('resolve-from');
const path = require('path');
const babel = require("@babel/core");

/**
 * Resolve implements the required jbld algorithm for finding a file:
 * - If the path is inside srcDir return relative path.
 * - If the path is outside of srcDir return null.
 * - If resolve returns null this is not a valid, knowable import. Don't emit.
 *
 * TODO: Want to implement a typescript style paths resolution so that we can
 *       copy files outside of srcDir into a node_modules folder for example
 *       so that the files can be optionally processed.
 *
 * @param {string} srcDir Source directory
 * @param {string} src    Relative source path
 * @param {string} name   Name to resolve
 * @returns {string}
 */
function resolve(srcDir, src, name) {
  const resolved = path.relative(srcDir,
    resolveFrom(path.dirname(src), name),
  );
  if (resolved.startsWith("../") || resolved === "..") {
    return null;
  }
  return resolved;
}

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

class BabelPlugin {
  constructor(opts) {
    this.opts = opts;
  }

  run(code, source) {
    const imports = [];
    const opts = {
      ...this.opts,
      plugins: (this.opts.plugins ? this.opts.plugins : []).concat([
        "@babel/plugin-transform-runtime",
        babelResolvePlugin(imports, source), // Must be last.
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
  resolve,
  BabelPlugin,
};
