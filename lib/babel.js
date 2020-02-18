const babel = require("@babel/core");
const { resolve } = require("./index");

function importVisitor(callback) {
  function isRequire(path) {
    return (
      path.node.callee.type === "Identifier" &&
      path.node.callee.name === "require"
    );
  }
  function isImportCall(path) {
    return path.node.callee.type === "Import";
  }
  function isLiteralArg(path) {
    return (
      path.node.arguments.length > 0 &&
      path.node.arguments[0].type === "StringLiteral"
    );
  }
  return {
    manipulateOptions(opts, parserOpts) {
      parserOpts.plugins.push("dynamicImport", "importMeta");
    },
    visitor: {
      CallExpression(path) {
        if ((isRequire(path) || isImportCall(path)) && isLiteralArg(path)) {
          callback(path.node.arguments[0]);
        }
      },
      ImportDeclaration(path) {
        callback(path.node.source);
      }
    }
  };
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
    imports.push({ name: node.value, resolved, kind: "static" });
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
        babelResolvePlugin(imports, source)
      ]),
      filename: source.name,
      ast: false
    };
    let result = babel.transformSync(code, opts);
    return {
      imports,
      output: result.code
    };
  }
}

module.exports = { BabelPlugin };
