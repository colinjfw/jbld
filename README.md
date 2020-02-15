

## Arch

- compiler:
  - Traverses all workspace files.
  - Run a set of plugins transforming files.
  - Special plugin possible to output imports, (async / static).
  - Outputs to a new dist/ directory with same file structure.
  - Writes importmap structure to disk for linker. Plugins possible.

- linker:
  - Reads importmap from compiled files to link files together.
  - Optional graph of files is submitted to a plugin to build chunks.
