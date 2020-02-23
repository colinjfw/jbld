export interface Config {
  options?: {
    // Env configures the process environment variables.
    env?: { [k: string]: string };
    // Entrypoints are the initial entrypoints.
    entrypoints?: string[];
    // Base url is prepended to all assets.
    baseUrl?: string;
    // Asset path groups all bundled assets under this path.
    assetPath?: string;
    // Source directory to resolve files starting from here.
    source?: string;
    // Output directory is where compiled and bundled files will exist.
    output?: string;
    // Workers is the count of total node processes used.
    workers?: number;
    // Public configures copying public assets to the output dir.
    public?: {
      // Directory to copy assets from.
      dir?: string;
      // HTML configures templates to use for inserting script and link tags.
      html?: string[];
    };
  };
  // Rules configure processing directives.
  rules?: Rule[];
}

/**
 * Configuration should be instantiated and returned from the config.jbld.js
 * file.
 *
 * Example:
 * ```
 * // config.jbld.js
 * const { Configuration } = require("jbld");
 *
 * module.exports = new Configuration({ rules: [] });
 * ```
 */
export class Configuration {
  constructor(c: Config);
}

export interface Source {
  name: string;
  src: string;
  dst: string;
  srcDir: string;
  dstDir: string;
}

export interface Import {
  // Name is the name in the file.
  name: string;
  // Resolved is the actual file on disk relative to the source directory.
  resolved: string;
  // Kind is a metadata field for determining if this is static or dynamic.
  kind: 'static' | 'dynamic';
}

export interface Output {
  output?: string;
  // Imports lists all the resolved imports of a file.
  imports?: Import[];
  // Type optionally declares the file type. Defaults to the file extension.
  type?: string;
}

export interface Plugin {
  run(code: string, source: Source): Output;
}

export interface Rule {
  // Test is a regex to match a filename.
  test?: RegExp;
  // Use will run the below plugins on the file.
  use?: Plugin[];
}
