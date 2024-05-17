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
        'html': {
          fontSize: "16px", // Default font size, which is typically 16px
          '@screen md': {
            fontSize: "14px", // Font size on medium (md) breakpoint, which is typically 14px
          },
        },
      })
     }),
     require('@tailwindcss/typography'),
     require('daisyui')
  ],
   daisyui: {
      themes:[
         "autumn",
         "sunset",
      ],
      darkTheme: "sunset",
   },
}
