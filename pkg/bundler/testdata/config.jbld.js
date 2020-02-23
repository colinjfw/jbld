const { Configuration } = require("../../../lib");

class TestPlugin {
  run(input) {
    if (input.name === "index.js") {
      return {
        type: "js",
        imports: [{
          kind: 'static',
          name: 'file2',
          resolved: 'file2.js'
        }],
        output: input
      };
    } else {
      return {
        type: "js",
        output: input,
      }
    }
  }
}

module.exports = new Configuration({
  rules: [{
    use: [new TestPlugin()],
  }],
});
