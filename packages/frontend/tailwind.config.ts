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
        primary: "#311744",
        primaryGrad1: "#201f24",
        primaryGrad2: "#201f24",
        secondary: "#FFD700",
        red: "#FF0000",
        black: "#000000",
        background: "#121212",
        border: "#462c6e",
        purple: "#cc00ff",
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
      "light","dark",{
      mytheme: {
        "base-100": "#80658c",
      }
    }],
  },


} satisfies Config;
