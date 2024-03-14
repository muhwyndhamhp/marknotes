const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin');

module.exports = {
   darkMode: 'class',
  content: ["./pub/**/*.templ", "./src/**/*.js"],
  theme: {
    screens: {
        'xs': '320px',
      ...defaultTheme.screens,
    },
     extend: {},
  },
  plugins: [
    plugin(function({ addBase }) {
      addBase({
         'html': { fontSize: "16px" },
       })
     }),
     require('@tailwindcss/typography'),
     require("daisyui")
  ],
   daisyui: {
      themes:[
         "fantasy",
         "light",
         "dark",
         "cupcake",
         "bumblebee",
         "emerald",
         "corporate",
         "synthwave",
         "retro",
         "cyberpunk",
         "valentine",
         "halloween",
         "garden",
         "forest",
         "aqua",
         "lofi",
         "pastel",
         "wireframe",
         "black",
         "luxury",
         "dracula",
         "cmyk",
         "autumn",
         "business",
         "acid",
         "lemonade",
         "night",
         "coffee",
         "winter",
         "dim",
         "nord",
         "sunset",
      ],
      darkTheme: "night",
   },
}
