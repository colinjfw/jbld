package bundler

const runtime = `
var __modules = {
  defined: {},
  cache: {},
  resolve: {},
  define: function(name, modFn) {
    this.defined[name] = { modFn: modFn };
  }.bind(__modules),
  require: function(name) {
    if (this.resolve[name]) {
      name = this.resolve[name];
    }
    if (!this.defined[name]) {
      throw new Error('Module not found ' + name);
    }
    var def = this.defined[name];
    if (!this.cache[name]) {
      var mod = { exports: {} };
      def.modFn(mod, mod.exports, this.require);
      this.cache[name] = mod;
    }
    return this.cache[name].exports;
  }.bind(__modules),
};
`
