#!/usr/bin/env node

const argv = require("minimist")(process.argv.slice(2));
const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");
const help = `Javascript build tool and bundler. Executing this script with no arguments will compile and bundle based on configuration provided.

  --watch  Configures watch mode
  --serve  Will serve a static file server on the specified interface
  --mode   Sets the NODE_ENV variable (default: development)
  --config Configuration file to load (default: ./config.jbld.js)
  --clean  Clean out the previous output directory

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

if (argv.clean) {
  try {
    fs.rmdirSync(opts.output, { recursive: true });
  } catch (error) {
    // Ignore.
  }
}

spawn(
  path.join(__dirname, "bin", target()),
  [JSON.stringify(options)],
  { stdio: 'inherit' },
).on("exit", (code) => {
  process.exit(code);
});

function abs(n) {
  return path.join(process.cwd(), n);
}

function arch() {
  switch (process.arch) {
    case "arm":
      return "arm";
    case "arm64":
      return "arm64";
    case "ia32":
      return "386";
    case "x32":
      return "amd32";
    case "x64":
      return "amd64";
    case "mipsel":
    case "mips":
    case "ppc":
    case "ppc64":
    case "s390":
    case "s390x":
      throw new Error("Architecture unsupported " + process.arch);
    default:
      throw new Error("Unknown architecture " + process.arch);
  }
}

function os() {
  switch (process.platform) {
    case "darwin":
      return "darwin";
    case "win32":
      return "windows";
    case "linux":
      return "linux";
    case "freebsd":
    case "openbsd":
    case "sunos":
    case "aix":
      throw new Error("OS unsupported " + process.platform);
    default:
      throw new Error("Unknown OS" + process.platform);
  }
}

function target() {
  const supported = [
    "darwin-amd64",
    "windows-amd64",
    "linux-amd64",
    "linux-arm",
    "linux-arm64",
    "linux-386"
  ];
  let t = os() + "-" + arch();
  if (supported.indexOf(t) === -1) {
    throw new Error("Unsupported architecture " + t);
  }
  return t;
}
