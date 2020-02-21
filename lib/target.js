const fs = require("fs");
const path = require("path");

const supported = [
  "darwin-amd64",
  "windows-amd64",
  "linux-amd64",
  "linux-arm",
  "linux-arm64",
  "linux-386"
];

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
  let t = os() + "-" + arch();
  if (supported.indexOf(t) === -1) {
    throw new Error("Unsupported architecture " + t);
  }
  return t;
}

module.exports = { target };
