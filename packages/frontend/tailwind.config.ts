import { type Config } from "tailwindcss";
import { fontFamily } from "tailwindcss/defaultTheme";

export default {
  content: ["./src/**/*.tsx"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-sans)", ...fontFamily.sans],
      },


    },
    colors:{
      primary:"#311744",
      secondary:"#FFD700",
      
      
    }
  },
  plugins: [require("daisyui")],

} satisfies Config;
