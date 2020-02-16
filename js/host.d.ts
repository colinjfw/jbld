export interface Source {
  name: string;
  plugins: string[];
  src: string;
  dst: string;
  srcDir: string;
  dstDir: string;
}

export interface Import {
  name: string;
  resolved: string;
  kind: string;
}

export interface Output {
  output: string;
  imports: Import[];
}

export interface Plugin {
  run(code: string, source: Source): Output;
}
