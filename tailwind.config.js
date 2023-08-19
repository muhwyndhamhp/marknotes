const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["./public/**/*.{html,js}"],
  theme: {
    screens: {
        'xs': '320px',
      ...defaultTheme.screens,
    },
    extend: {},
  },
  plugins: [],
}