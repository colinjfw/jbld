package bundler

const runtime = `(function(){var g=typeof window==="undefined"?global:window;if(g.__modules){return}var m={defined:{},cache:{},resolve:{}};m.require=function(name){if(m.resolve[name]){name=m.resolve[name]}if(!m.defined[name]){throw new Error("Module not found "+name)}var def=m.defined[name];if(!m.cache[name]){var mod={exports:{}};def.modFn(mod,mod.exports,m.require);m.cache[name]=mod}return m.cache[name].exports};m.define=function(name,modFn){m.defined[name]={modFn:modFn}};m.main=function(chunks,main){if(chunks.length===0){m.require(main);return}var l=0;function run(){l++;if(l===chunks.length){m.require(main)}}function ls(src){var s=document.createElement("script");s.src=src;s.onload=run;document.head.appendChild(s)}chunks.forEach(ls)};g.__modules=m})();
`
