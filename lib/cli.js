#!/usr/bin/env node

const argv = require("minimist")(process.argv.slice(2));
const path = require("path");
const { target } = require("./target");
const { spawn } = require("child_process");
const help = `Javascript build tool and bundler. Executing this script with no arguments will compile and bundle based on configuration provided.

  --watch  Configures watch mode
  --serve  Will serve a static file server on the specified interface
  --mode   Sets the NODE_ENV variable (default: development)
  --config Configuration file to load (default: ./config.jbld.js)

`;

if (argv.h || argv.help) {
  process.stderr.write(help);
  process.exit(0);
}

const configFile = require.resolve(
  argv.config ? abs(argv.config) : abs("config.jbld.js")
);
const configuration = require(configFile);
const hostJs = require.resolve("./host.js");
const defaults = {
  env: {},
  entrypoints: ["src/index.js"],
  baseUrl: "/",
  assetPath: "static",
  source: ".",
  output: "./dist",
  workers: 5,
  public: {
    dir: "./public",
    html: ["index.html"],
  },
};

const opts = Object.assign({}, defaults, configuration.options());
const options = {
  watch: opts.watch || argv.watch || false,
  serve: opts.serve || argv.serve || "",
  mode: opts.mode || argv.mode || "development",
  env: opts.env || {},
  compiler: {
    hostJs,
    configFile,
    entrypoints: opts.entrypoints,
    sourceDir: abs(opts.source),
    outputDir: path.join(abs(opts.output), "target"),
    workers: opts.workers,
  },
  bundler: {
    hostJs,
    configFile,
    baseUrl: opts.baseUrl,
    outputDir: path.join(abs(opts.output), "bundle"),
    assetPath: opts.assetPath,
    public: opts.public,
  },
};

spawn(
  path.join(__dirname, "bin", target()),
  [JSON.stringify(options)],
  { stdio: 'inherit' },
);

function abs(n) {
  return path.join(process.cwd(), n);
}
