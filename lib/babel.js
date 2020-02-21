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
  let seen = {};
  return importVisitor(node => {
    if (seen[node.value]) {
      return;
    }
    const resolved = resolve(source.srcDir, source.src, node.value);
    if (!resolved) {
      return;
    }
    node.value = resolved;
    seen[resolved] = true;
    imports.push({ name: resolved, resolved, kind: "static" });
  });
}

/**
 * Returns a preset for the babel plugin.
 *
 * @param {Source} source
 */
function babelPreset(source) {
  let imports = [];
  return {
    imports,
    preset: {
      plugins: [babelResolvePlugin(imports, source)],
    }
  }
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
    let { imports, preset } = babelPreset(source);
    let opts = {
      ...this.opts,
      presets: (this.opts.presets ? this.opts.presets : []).concat([
        preset,
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
