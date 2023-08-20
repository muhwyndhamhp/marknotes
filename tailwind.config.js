const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["./public/**/*.{html,js}"],
  theme: {
    screens: {
        'xs': '320px',
      ...defaultTheme.screens,
    },
    fontFamily: {
      'sans': ['Nunito Sans', 'sans-serif'],
      'serif': ['Merriweather', 'serif'],
      'mono': ['PT Mono', 'mono'],
    },
    extend: {},
  },
  plugins: [],
}