package bundler

const runtime = `(function(){var g=typeof window==="undefined"?global:window;if(g.__modules){return}var m={defined:{},cache:{},resolve:{}};m.require=function(name){if(m.resolve[name]){name=m.resolve[name]}if(!m.defined[name]){throw new Error("Module not found "+name)}var def=m.defined[name];if(!m.cache[name]){var mod={exports:{}};def.modFn(mod,mod.exports,m.require);m.cache[name]=mod}return m.cache[name].exports};m.define=function(name,modFn){m.defined[name]={modFn:modFn}};g.__modules=m})();
`
