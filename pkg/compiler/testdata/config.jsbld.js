class TestPlugin {
  run(input) {
    return {
      imports: [{
        kind: 'static',
        name: 'file2',
        resolved: 'file2.js'
      }],
      output: input
    };
  }
}


module.exports = {
  plugins: {
    test: new TestPlugin(),
  },
};
