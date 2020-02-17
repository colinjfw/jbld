const main = require(process.argv[2]);
const readline = require("readline");
const host = readline.createInterface({
  input: process.stdin,
  output: new require('stream').Writable(),
  terminal: false
});

host.on("line", function (line) {
  try {
    let input = JSON.parse(line);
    let res = main[input.method](input.req);
    let out = JSON.stringify({ res: res });
    process.stdout.write(out + "\n");
  } catch (err) {
    console.error(err);
    let out = JSON.stringify({ err: err.message });
    process.stdout.write(out + "\n");
  }
});
