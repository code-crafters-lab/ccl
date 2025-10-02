import { defineConfig, Options } from "tsup";

export default defineConfig((options: Options) => [
  {
    entry: ["src/index.ts", "src/node.ts", "src/web.ts"],
    format: ["esm", "cjs", "iife"],
    treeshake: true,
    splitting: true,
    // dts: true,
    // minify: false,
    // target: "es6",
    clean: true,
    // legacyOutput: true,
    outExtension: ({ format, pkgType }) => {
      // if (format === "iife") {
      //   return { js: `.min.js` };
      // }
      return { js: `.${format}.js` };
    },
    // sourcemap: true,
    // minify: true,
    globalName: "MyTool", // IIFE 格式的全局变量名
    ...options,
  },
]);
