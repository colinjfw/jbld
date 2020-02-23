# jbld

Very fast and efficient Javascript compilation tool and bundler. A replacement
for webpack.

Current bundling solutions like Webpack have performance problems on large
projects. Jbld has performance critical components written in Golang while still
calling out to Javascript so that you can take advantage of plugins like Babel
without sacrificing performance.

Try it out on a project built with create-react-app:

    git clone https://github.com/colinjfw/jbld
    cd examples/react-app
    yarn install
    yarn start

## Status

This project is not production ready as it's missing features like source maps
or chunk optimizations.

It's currently built as an example of how we could build a faster build tool
for Javascript.

## Getting started

1. `npm install --save-dev jbld` or `yarn add -D jbld`.
2. Add a script your your package.json `"build": "jbld"`.
2. Create a configuration file `./config.jbld.js`.
3. Run `yarn|npm run build`.

The `config.jbld.js` file should contain configuration rules for setting up your
project. A basic example which runs no plugins across a file is as follows:

```javascript
// config.jbld.js
const { Configuration } = require("jbld");

module.exports = new Configuration({ rules: [] });
```

A babel plugin is provided out of the box which is probably one of the most
commonly used plugins:

```javascript
const { BabelPlugin } = require("jbld/babel");
const { Configuration } = require("jbld");

const babel = new BabelPlugin({ ... }); // Options: https://babeljs.io/docs/en/options
module.exports = new Configuration({
  options: { },
  rules: [{
    test: /\.js$/,
    use: [babel],
  }]
});
```

[View the full set of configuration options](lib/index.d.ts).

## How it works

Jbld is logically separated into two components:

- Compiler
- Bundler

### Compiler

The compiler works on a single file at a time by traversing the tree of your
project starting at the entrypoints. Files are processed by running them through
a set of Javascript plugins which transform the files and discover imports
returning the information to the Go process. The process looks like:

1. Compute the compilation hash of the file. If we have already compiled this
   and nothing has changed continue.
2. If not cached, pass the Node process the filename to compile.
3. Node process returns imports after running all plugins.

The Go process first reads an entrypoint and checks if it's hash indicates that
this file needs to be processed. If it does it calls out to a Node process via
a stdin/stdout interface and requests the necessary plugins to be run over the
file. The Node process returns imports and a success message telling the Go
process that it may continue.

Every compiled file has an associated `.o` file that lives beside it which
indicates the hash of the file and allows robust caching the next time we run
through the compilation.

### Bundler

The bundler is written entirely in Golang and takes all compiled files as a
manifest and then writes them out to a set of single files. The included
Javascript runtime has the capability to fetch chunks and import files as
needed.
