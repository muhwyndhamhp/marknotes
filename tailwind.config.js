const defaultTheme = require('tailwindcss/defaultTheme')
const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ["./public/**/*.{html,js}"],
  theme: {
    screens: {
        'xs': '320px',
      ...defaultTheme.screens,
    },
    fontFamily: {
      'sans': ['Open Sauce Sans', 'sans-serif'],
      'serif': ['Spicy Kebab', 'serif'],
      'mono': ['Overpass Mono', 'mono'],
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
  ],
}