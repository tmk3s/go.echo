import type { Config } from "tailwindcss";

// トークンの実体は app/globals.css の CSS 変数。
// ここではそれを Tailwind のユーティリティ名に割り当てるだけにして、
// 色の定義がコードの2箇所に散らばらないようにしている。
const token = (name: string) => `rgb(var(--color-${name}) / <alpha-value>)`;

const config: Config = {
  // ThemeToggle が html に付け外しする dark クラスで切り替える
  darkMode: "class",
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./consts/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        canvas: token("canvas"),
        surface: {
          DEFAULT: token("surface"),
          muted: token("surface-muted"),
        },
        body: token("body"),
        muted: token("muted"),
        line: token("line"),
        brand: {
          DEFAULT: token("brand"),
          hover: token("brand-hover"),
          contrast: token("brand-contrast"),
          subtle: token("brand-subtle"),
          text: token("brand-text"),
        },
        success: {
          DEFAULT: token("success"),
          hover: token("success-hover"),
        },
        danger: {
          DEFAULT: token("danger"),
          hover: token("danger-hover"),
        },
      },
    },
  },
  plugins: [],
};
export default config;
