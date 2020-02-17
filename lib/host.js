const compiler = JSON.parse(process.argv[2]);
const config = require(compiler.configFile);
const readline = require("readline");
const host = readline.createInterface({
  input: process.stdin,
  output: new require('stream').Writable(),
  terminal: false
});

host.on("line", function (line) {
  try {
    let input = JSON.parse(line);
    let resp = run(input);
    let out = JSON.stringify({ type: resp.type, imports: resp.imports });
    process.stdout.write(out + "\n");
  } catch (err) {
    console.error(err);
    let out = JSON.stringify({ err: err.message });
    process.stdout.write(out + "\n");
  }
});

function run(input) {
  return config.run(compiler, input);
}
