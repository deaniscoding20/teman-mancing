import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/features/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        brand: {
          primary: "#123B50",    // Deep Ocean Blue
          secondary: "#35A7B5",  // Turquoise / Lagoon
          accent: "#7BAE7F",     // Sage / Foliage Green
          background: "#F5F7F7", // Pale clean background
          surface: "#FFFFFF",    // White card / surface
          text: "#172B36",       // Deep Slate readable text
          muted: "#627B87",      // Secondary text
          border: "#E2E8EA",     // Clean subtle border
        },
      },
      fontFamily: {
        sans: ["var(--font-sans)", "system-ui", "-apple-system", "sans-serif"],
      },
    },
  },
  plugins: [],
};

export default config;
