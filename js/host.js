var filename = process.argv[2];
var readline = require("readline");
var fs = require("fs");
var path = require("path");
var config = require(filename);
var host = readline.createInterface({
  input: process.stdin,
  output: new require('stream').Writable(),
  terminal: false
});

host.on("line", function (line) {
  try {
    var input = JSON.parse(line);
    var imports = run(input.plugins, input.src, input.dst);
    var out = JSON.stringify({ err: null, imports: imports });
    process.stdout.write(out + "\n");
  } catch (err) {
    console.error(err);
    var out = JSON.stringify({ err: err.message, imports: [] });
    process.stdout.write(out + "\n");
  }
});

function run(plugins, src, dst) {
  var code = fs.readFileSync(src).toString('utf-8');
  var imports = [];

  plugins.forEach(function (plugin) {
    var run = config.plugins[plugin]
    if (!run) {
      throw new Error('plugin not found ' + plugin);
    }
    var resp = run.run(code, src);
    code = resp.output;
    imports = resp.imports;
  });

  fs.mkdirSync(path.dirname(dst), { recursive: true });
  fs.writeFileSync(dst, code);
  return imports;
}
