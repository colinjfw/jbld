const path = require('path');
const fs = require('fs');
const Module = require('module');

/**
 * Resolve from handles resolving from a directory.
 *
 * @param {string} fromDirectory
 * @param {string} moduleId
 */
function resolveFrom(fromDirectory, moduleId) {
  try {
    fromDirectory = fs.realpathSync(fromDirectory);
  } catch (error) {
    if (error.code === 'ENOENT') {
      fromDirectory = path.resolve(fromDirectory);
    } else {
      throw error;
    }
  }

  const fromFile = path.join(fromDirectory, 'noop.js');
  const resolveFileName = () => Module._resolveFilename(moduleId, {
    id: fromFile,
    filename: fromFile,
    paths: Module._nodeModulePaths(fromDirectory)
  });
  return resolveFileName();
}

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
 * Configuration handles providing the host with something to call.
 */
class Configuration {
  constructor(config) {
    this.config = config;
    if (!this.config.rules) {
      this.config.rules = [];
    }
  }

  addRule(rule) {
    this.config.rules.push(rule);
    return this;
  }

  run(compiler, input) {
    let src = path.join(compiler.sourceDir, input.name);
    let dst = path.join(compiler.outputDir, input.name);
    let source = {
      name: input.name,
      src: src,
      dst: dst,
      srcDir: compiler.sourceDir,
      dstDir: compiler.outputDir,
    };

    let code = fs.readFileSync(src).toString('utf-8');
    let imports = [];

    this.config.rules.forEach(rule => {
      if (!testRule(rule, source.name)) {
        return;
      }
      rule.use.forEach(plugin => {
        let run = plugin.run(code, source);
        if (run.output) code = run.output;
        if (run.imports) imports = run.imports;
      });
    });

    fs.mkdirSync(path.dirname(dst), { recursive: true });
    fs.writeFileSync(dst, code);
    return { imports };
  }
}

function testRule(rule, name) {
  if (!rule.test) return true;
  return rule.test.test(name);
}

module.exports = { Configuration, resolve };
