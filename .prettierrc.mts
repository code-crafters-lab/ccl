import { type Config } from "prettier";

/**
 * 定义包含 extends 的扩展配置类型
 * @see https://prettier.io/docs/configuration
 * @type {import("prettier").Config}
 */
interface ExtendedConfig extends Config {
  extends?: string | string[];
}

const config: ExtendedConfig = {
  // 继承 standard 预设规则
  extends: ["prettier-config-standard"],
  // 加载 Astro 和 Tailwind 插件
  plugins: ["prettier-plugin-sort-imports", "prettier-plugin-astro", "prettier-plugin-tailwindcss"],
  // 可选：自定义规则（会覆盖预设中的同名规则）
  printWidth: 120,
  tabWidth: 2,
  // 针对 Astro 文件的特殊配置（可选）
  overrides: [
    {
      files: "*.astro",
      options: {
        parser: "astro",
        // Astro 特有的规则（如组件标签换行等）
        bracketSameLine: false,
      },
    },
  ],
};

export default config;
