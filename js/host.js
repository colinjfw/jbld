var compiler = JSON.parse(process.argv[2]);
var readline = require("readline");
var fs = require("fs");
var path = require("path");
var config = require(compiler.configFile);
var host = readline.createInterface({
  input: process.stdin,
  output: new require('stream').Writable(),
  terminal: false
});

host.on("line", function (line) {
  try {
    var input = JSON.parse(line);
    var imports = run(input);
    var out = JSON.stringify({ err: null, imports: imports });
    process.stdout.write(out + "\n");
  } catch (err) {
    console.error(err);
    var out = JSON.stringify({ err: err.message, imports: [] });
    process.stdout.write(out + "\n");
  }
});

function run(input) {
  var src = path.join(compiler.sourceDir, input.name);
  var dst = path.join(compiler.outputDir, input.name);
  var opts = {
    name: input.name,
    plugins: input.plugins,
    src: src,
    dst: dst,
    srcDir: compiler.sourceDir,
    dstDir: compiler.outputDir,
  };

  var code = fs.readFileSync(src).toString('utf-8');
  var imports = [];

  input.plugins.forEach(function (plugin) {
    var run = config.plugins[plugin]
    if (!run) {
      throw new Error('plugin not found ' + plugin);
    }
    var resp = run.run(code, opts);
    code = resp.output;
    imports = resp.imports;
  });

  fs.mkdirSync(path.dirname(dst), { recursive: true });
  fs.writeFileSync(dst, code);
  return imports;
}
