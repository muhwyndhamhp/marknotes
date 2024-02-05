const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ["./pub/**/*.templ"],
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
         'html': { fontSize: "16px" },
       })
     }),
     require('@tailwindcss/typography'),
     require("daisyui")
  ],
   daisyui: {
      themes:[
         {
            garden: {
               ...require("daisyui/src/theming/themes")["garden"],
               "primary": "#e11d48",
               "secondary": "#db2777",
               "accent": "#93c5fd",
               "neutral": "#fecdd3",
               "neutral-content": "#4c4949",
               "base-100": "#f5f5f4",
               "info": "#06b6d4",
               "success": "#34d399",
               "warning": "#fcd34d",
               "error": "#f87171",            
            }
         },         
         "dark",
      ],
      darkTheme: "dark",
   },
}
