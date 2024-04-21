import { type Config } from "tailwindcss";
import { fontFamily } from "tailwindcss/defaultTheme";

export default {
  content: ["./src/**/*.tsx"],

  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-sans)", ...fontFamily.sans],
      },

      colors: {
        primary: "#cc00ff",
        primaryGrad1: "#201f24",
        primaryGrad2: "#201f24",
        secondary: "#FFD700",
        black: "#000000",
        background: "#121212",
        background3:"#352F44",
        background4:"#5C5470",
        border: "#462c6e",
        purple: "#cc00ff",
        purple3:"#14001a",
        card:"#17191a"



      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
      }
    },

  },

  plugins: [require("daisyui")],
  daisyui: {

    themes: [
      {
      mytheme: {
        ...require("daisyui/src/theming/themes")["dark"],
        primary: "#cc00ff",
        primaryGrad1: "#201f24",
        primaryGrad2: "#201f24",
        secondary: "#FFD700",
        // black: "#000000",
        background: "#121212",
        background3:"#352F44",
        background4:"#5C5470",
        border: "#462c6e",
        purple: "#cc00ff",
        purple3:"#14001a",
        card:"#17191a"
      }
    }],
  },


} satisfies Config;
