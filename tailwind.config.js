const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ["./pub/**/*.templ", "./src/**/*.js"],
  theme: {
    screens: {
        'xs': '320px',
      ...defaultTheme.screens,
    },
    fontFamily: {
      'sans': ['Open Sauce Sans', 'system-ui', 'sans-serif'],
      'serif': ['Spicy Kebab', 'system-ui', 'serif'],
      'mono': ['Overpass Mono', 'system-ui', 'mono'],
    },
    extend: {},
  },
  plugins: [
    plugin(function({ addBase }) {
      addBase({
         'html': { fontSize: "14px" },
       })
     }),
     require('@tailwindcss/typography'),
     require("daisyui")
  ],
   daisyui: {
      themes:[
         "fantasy",
         "light",
         "cupcake",
         "wireframe",
         "autumn",
         "emerald",
         "retro",
         "cyberpunk",
         "pastel",
         "lemonade",
         "cmyk",
         {

            mytheme: {

               "primary": "#0e7490",
               "secondary": "#0f766e",
               "accent": "#e11d48",
               "neutral": "#1f2937",
               "base-100": "#f5efef",
               "info": "#38bdf8",
               "success": "#34d399",
               "warning": "#fb923c",
               "error": "#f87171",
            }

         },
         "dark",
         "night",
         "forest",
      ],
      darkTheme: "night",
   },
}
