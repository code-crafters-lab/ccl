import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => ({
  entry: ["src/index.ts", "src/node.ts", "src/web.ts"],
  format: ["esm", "cjs", "iife"],
  treeshake: false,
  splitting: true,
  // dts: true,
  // minify: false,
  // target: "es6",
  clean: true,
  // legacyOutput: true,
  // sourcemap: true,
  ...options,
}));
